package oss

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/readers"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/retry"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/signer"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/transport"
)

type Options struct {
	Region string

	Endpoint *url.URL

	RetryMaxAttempts int

	Retryer retry.Retryer

	Signer signer.Signer

	CredentialsProvider credentials.CredentialsProvider

	HttpClient *http.Client

	ResponseHandlers []func(*http.Response) error

	UrlStyle UrlStyleType

	FeatureFlags FeatureFlagsType

	OpReadWriteTimeout *time.Duration
}

func (c Options) Copy() Options {
	to := c
	to.ResponseHandlers = make([]func(*http.Response) error, len(c.ResponseHandlers))
	copy(to.ResponseHandlers, c.ResponseHandlers)
	return to
}

type Client struct {
	options Options
}

func NewClient(cfg *Config, optFns ...func(*Options)) *Client {
	options := Options{
		Region:              cfg.Region,
		RetryMaxAttempts:    cfg.RetryMaxAttempts,
		Retryer:             cfg.Retryer,
		CredentialsProvider: cfg.CredentialsProvider,
		HttpClient:          cfg.HttpClient,
	}
	resolveEndpoint(cfg, &options)
	resolveRetryer(cfg, &options)
	resolveHTTPClient(cfg, &options)
	resolveSigner(cfg, &options)
	resolveUrlStyle(cfg, &options)
	resolveFeatureFlags(cfg, &options)

	for _, fn := range optFns {
		fn(&options)
	}

	client := &Client{
		options: options,
	}

	return client
}

func resolveEndpoint(cfg *Config, o *Options) {
	if cfg.Endpoint == nil {
		return
	}
	scheme := "http"
	endpoint := *cfg.Endpoint
	if strings.HasPrefix(endpoint, "http://") {
		scheme = "http"
		endpoint = endpoint[len("http://"):]
	} else if strings.HasPrefix(endpoint, "https://") {
		scheme = "https"
		endpoint = endpoint[len("https://"):]
	}
	o.Endpoint, _ = url.Parse(fmt.Sprintf("%s://%s", scheme, endpoint))
}

func resolveRetryer(cfg *Config, o *Options) {
	if o.Retryer != nil {
		return
	}

	o.Retryer = retry.NewStandard()
}

func resolveHTTPClient(cfg *Config, o *Options) {
	if o.HttpClient != nil {
		return
	}

	//config in http.Transport
	custom := []func(*http.Transport){}
	if cfg.InsecureSkipVerify != nil {
		custom = append(custom, transport.InsecureSkipVerify(*cfg.InsecureSkipVerify))
	}
	if cfg.ProxyFromEnvironment != nil && *cfg.ProxyFromEnvironment {
		custom = append(custom, transport.ProxyFromEnvironment())
	}
	if cfg.ProxyHost != nil {
		if url, err := url.Parse(*cfg.ProxyHost); err == nil {
			custom = append(custom, transport.HttpProxy(url))
		}
	}

	//config in transport  package
	tcfg := &transport.Config{}
	if cfg.ConnectTimeout != nil {
		tcfg.ConnectTimeout = cfg.ConnectTimeout
	}
	if cfg.ReadWriteTimeout != nil {
		tcfg.ReadWriteTimeout = cfg.ReadWriteTimeout
	}
	if cfg.EnabledRedirect != nil {
		tcfg.EnabledRedirect = cfg.EnabledRedirect
	}

	o.HttpClient = transport.NewHttpClient(tcfg, custom...)
}

func resolveSigner(cfg *Config, o *Options) {
	if o.Signer != nil {
		return
	}

	o.Signer = signer.SignerV1{}
}

func resolveUrlStyle(cfg *Config, o *Options) {
	if cfg.UseCName != nil && *cfg.UseCName {
		o.UrlStyle = UrlStyleCName
	} else if cfg.UsePathStyle != nil && *cfg.UsePathStyle {
		o.UrlStyle = UrlStylePath
	} else {
		o.UrlStyle = UrlStyleVirtualHosted
	}

	// if the endpoint is ip, set to path-style
	if o.Endpoint != nil {
		if ip := net.ParseIP(o.Endpoint.Hostname()); ip != nil {
			o.UrlStyle = UrlStylePath
		}
	}
}

func resolveFeatureFlags(cfg *Config, o *Options) {
	//TODO
}

func (c *Client) invokeOperation(ctx context.Context, input *OperationInput, optFns []func(*Options)) (output *OperationOutput, err error) {
	options := c.options.Copy()
	opOpt := Options{}

	for _, fn := range optFns {
		fn(&opOpt)
	}

	applyOperationOpt(&options, &opOpt)

	ctx = applyOperationContext(ctx, &options)

	output, err = c.sendRequest(ctx, input, &options)

	if err != nil {
		return output, &OperationError{
			name: input.OpName,
			err:  err}
	}

	return output, err
}

func (c *Client) sendRequest(ctx context.Context, input *OperationInput, opts *Options) (output *OperationOutput, err error) {
	// covert input into httpRequest
	if !isValidEndpoint(opts.Endpoint) {
		return output, NewErrParamInvalid("Endpoint")
	}

	// host & path
	host, path := buildURL(input, opts)
	strUrl := fmt.Sprintf("%s://%s%s", opts.Endpoint.Scheme, host, path)

	// querys
	if len(input.Parameters) > 0 {
		var buf bytes.Buffer
		for k, v := range input.Parameters {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(url.QueryEscape(k))
			if len(v) > 0 {
				buf.WriteString("=" + strings.Replace(url.QueryEscape(v), "+", "%20", -1))
			}
		}
		strUrl += "?" + buf.String()
	}

	request, err := http.NewRequestWithContext(ctx, input.Method, strUrl, nil)
	if err != nil {
		return output, err
	}

	// headers
	for k, v := range input.Headers {
		if len(k) > 0 && len(v) > 0 {
			request.Header.Add(k, v)
		}
	}
	request.Header.Set("User-Agent", defaultUserAgent())

	// body
	var body readers.ReadSeekerNopClose
	if input.Body == nil {
		body = readers.ReadSeekNopCloser(strings.NewReader(""))
	} else {
		body = readers.ReadSeekNopCloser(input.Body)
	}
	len, _ := body.GetLen()
	if len >= 0 && request.Header.Get("Content-Length") == "" {
		request.ContentLength = len
	}
	request.Body = body

	//signing context
	subResource, _ := input.OpMetadata.Get(signer.SubResource).([]string)
	signingCtx := &signer.SigningContext{
		Product:     Ptr("oss"),
		Region:      Ptr(opts.Region),
		Bucket:      input.Bucket,
		Key:         input.Key,
		Request:     request,
		SubResource: subResource,
	}

	// send http request
	response, err := c.sendHttpRequest(ctx, signingCtx, opts)

	if err != nil {
		return output, err
	}

	// covert http response into output context
	output = &OperationOutput{
		Input:      input,
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Body:       response.Body,
		Headers:    response.Header,
	}

	// save other info by Metadata filed, ex. retry detail info
	//output.Metadata.Set()

	return output, err
}

func (c *Client) sendHttpRequest(ctx context.Context, signingCtx *signer.SigningContext, opts *Options) (response *http.Response, err error) {
	request := signingCtx.Request
	retryer := opts.Retryer
	body, _ := request.Body.(readers.ReadSeekerNopClose)
	bodyStart, _ := body.Seek(0, io.SeekCurrent)
	for tries := 1; tries <= retryer.MaxAttempts(); tries++ {
		if tries > 1 {
			delay, err := retryer.RetryDelay(tries, err)
			if err != nil {
				break
			}
			if err = sleepWithContext(ctx, delay); err != nil {
				err = &CanceledError{Err: err}
				break
			}

			if _, err = body.Seek(bodyStart, io.SeekStart); err != nil {
				break
			}
		}

		if response, err = c.sendHttpRequestOnce(ctx, signingCtx, opts); err == nil {
			break
		}

		if isContextError(ctx, &err) {
			err = &CanceledError{Err: err}
			break
		}

		if !readers.IsReaderSeekable(request.Body) {
			break
		}

		if !retryer.IsErrorRetryable(err) {
			break
		}
	}
	return response, err
}

func (c *Client) sendHttpRequestOnce(ctx context.Context, signingCtx *signer.SigningContext, opts *Options) (
	response *http.Response, err error,
) {
	if _, anonymous := opts.CredentialsProvider.(*credentials.AnonymousCredentialsProvider); !anonymous {
		cred, err := opts.CredentialsProvider.GetCredentials(ctx)
		if err != nil {
			return response, err
		}

		signingCtx.Credentials = &cred
		if err = c.options.Signer.Sign(ctx, signingCtx); err != nil {
			return response, err
		}
	}

	if response, err = c.options.HttpClient.Do(signingCtx.Request); err != nil {
		return response, err
	}

	for _, fn := range opts.ResponseHandlers {
		if err = fn(response); err != nil {
			return response, err
		}
	}

	return response, err
}

func buildURL(input *OperationInput, opts *Options) (host string, path string) {
	if input == nil || opts == nil || opts.Endpoint == nil {
		return host, path
	}

	var paths []string
	if input.Bucket == nil {
		host = opts.Endpoint.Host
	} else {
		switch opts.UrlStyle {
		default: // UrlStyleVirtualHosted
			host = fmt.Sprintf("%s.%s", *input.Bucket, opts.Endpoint.Host)
		case UrlStylePath:
			host = opts.Endpoint.Host
			paths = append(paths, *input.Bucket)
		case UrlStyleCName:
			host = opts.Endpoint.Host
		}
	}

	if input.Key != nil {
		paths = append(paths, escapePath(*input.Key, false))
	}

	return host, ("/" + strings.Join(paths, "/"))
}

func serviceErrorResponseHandler(response *http.Response) error {
	if response.StatusCode/100 == 2 {
		return nil
	}

	timestamp, err := time.Parse(http.TimeFormat, response.Header.Get("Date"))
	if err != nil {
		timestamp = time.Now()
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	se := &ServiceError{
		StatusCode:    response.StatusCode,
		Code:          "BadErrorResponse",
		RequestID:     response.Header.Get("x-oss-request-id"),
		Timestamp:     timestamp,
		RequestTarget: fmt.Sprintf("%s %s", response.Request.Method, response.Request.URL),
		Snapshot:      body,
	}

	if err != nil {
		se.Message = fmt.Sprintf("The body of the response was not readable, due to :%s", err.Error())
		return se
	}

	err = xml.Unmarshal(body, &se)
	if err != nil {
		len := len(body)
		if len > 256 {
			len = 256
		}
		se.Message = fmt.Sprintf("Failed to parse xml from response body due to: %s. With part response body %s.", err.Error(), string(body[:len]))
		return se
	}
	return se
}

func applyOperationOpt(c *Options, op *Options) {
	if c == nil || op == nil {
		return
	}

	if op.Endpoint != nil {
		c.Endpoint = op.Endpoint
	}

	if op.RetryMaxAttempts > 0 {
		c.RetryMaxAttempts = op.RetryMaxAttempts
	}

	if op.Retryer != nil {
		c.Retryer = op.Retryer
	}

	if c.Retryer == nil {
		c.Retryer = retry.NopRetryer{}
	}

	if op.OpReadWriteTimeout != nil {
		c.OpReadWriteTimeout = op.OpReadWriteTimeout
	}

	//response handler
	handlers := []func(*http.Response) error{
		serviceErrorResponseHandler,
	}
	handlers = append(handlers, c.ResponseHandlers...)
	handlers = append(handlers, op.ResponseHandlers...)
	c.ResponseHandlers = handlers
}

func applyOperationContext(ctx context.Context, c *Options) context.Context {
	if ctx == nil || c.OpReadWriteTimeout == nil {
		return ctx
	}
	return context.WithValue(ctx, "OpReadWriteTimeout", c.OpReadWriteTimeout)
}

// fieldInfo holds details for the input/output of a single field.
type fieldInfo struct {
	idx   int
	flags int
}

const (
	fRequire int = 1 << iota

	fTypeUsermeta
	fTypeXml
	fTypeTime
)

func parseFiledFlags(tokens []string) int {
	var flags int = 0
	for _, token := range tokens {
		switch token {
		case "required":
			flags |= fRequire
		case "time":
			flags |= fTypeTime
		case "xml":
			flags |= fTypeXml
		case "usermeta":
			flags |= fTypeUsermeta
		}
	}
	return flags
}

func validateInput(input *OperationInput) error {
	if input == nil {
		return NewErrParamNull("OperationInput")
	}

	if input.Bucket != nil && !isValidBucketName(input.Bucket) {
		return NewErrParamInvalid("OperationInput.Bucket")
	}

	if input.Key != nil && !isValidObjectName(input.Key) {
		return NewErrParamInvalid("OperationInput.Key")
	}

	if !isValidMethod(input.Method) {
		return NewErrParamInvalid("OperationInput.Method")
	}

	return nil
}

func (c *Client) marshalInput(request interface{}, input *OperationInput, handlers ...func(*OperationInput) error) error {
	// merge common fields
	if cm, ok := request.(RequestCommonInterface); ok {
		h, p, b := cm.GetCommonFileds()
		// headers
		if len(h) > 0 {
			if input.Headers == nil {
				input.Headers = map[string]string{}
			}
			for k, v := range h {
				input.Headers[k] = v
			}
		}

		// parameters
		if len(p) > 0 {
			if input.Parameters == nil {
				input.Parameters = map[string]string{}
			}
			for k, v := range p {
				input.Parameters[k] = v
			}
		}

		// body
		input.Body = b
	}

	val := reflect.ValueOf(request)
	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || input == nil {
		return nil
	}

	t := val.Type()
	for k := 0; k < t.NumField(); k++ {
		if tag, ok := t.Field(k).Tag.Lookup("input"); ok {
			// header|query|body,filed_name,[required,time,usermeta...]
			v := val.Field(k)
			var flags int = 0
			tokens := strings.Split(tag, ",")
			if len(tokens) < 2 {
				continue
			}

			// parse field flags
			if len(tokens) > 2 {
				flags = parseFiledFlags(tokens[2:])
			}

			// check required flag
			if isEmptyValue(v) {
				if flags&fRequire != 0 {
					return NewErrParamRequired(t.Field(k).Name)
				}
				continue
			}

			switch tokens[0] {
			case "query":
				if input.Parameters == nil {
					input.Parameters = map[string]string{}
				}
				if v.Kind() == reflect.Pointer {
					v = v.Elem()
				}
				input.Parameters[tokens[1]] = fmt.Sprintf("%v", v.Interface())
			case "header":
				if input.Headers == nil {
					input.Headers = map[string]string{}
				}
				if v.Kind() == reflect.Pointer {
					v = v.Elem()
				}
				if flags&fTypeUsermeta != 0 {
					if m, ok := v.Interface().(map[string]string); ok {
						for k, v := range m {
							input.Headers[tokens[1]+k] = v
						}
					}
				} else {
					input.Headers[tokens[1]] = fmt.Sprintf("%v", v.Interface())
				}
			case "body":
				if flags&fTypeXml != 0 {
					var b bytes.Buffer
					if err := xml.NewEncoder(&b).EncodeElement(
						v.Interface(),
						xml.StartElement{Name: xml.Name{Local: tokens[1]}}); err != nil {
						return &SerializationError{
							Err: err,
						}
					}
					input.Body = bytes.NewReader(b.Bytes())
				}
			}
		}
	}

	if err := validateInput(input); err != nil {
		return err
	}

	for _, h := range handlers {
		if err := h(input); err != nil {
			return err
		}
	}

	return nil
}

func discardBody(result interface{}, output *OperationOutput) error {
	var err error
	if output.Body != nil {
		defer output.Body.Close()
		_, err = io.Copy(io.Discard, output.Body)
	}
	return err
}

func unmarshalBodyXml(result interface{}, output *OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}
	if len(body) > 0 {
		if err = xml.Unmarshal(body, result); err != nil {
			err = &DeserializationError{
				Err:      err,
				Snapshot: body,
			}
		}
	}
	return err
}

func unmarshalBodyDefault(result interface{}, output *OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}

	// extract body
	if len(body) > 0 {
		contentType := output.Headers.Get("Content-Type")
		switch contentType {
		case "application/xml":
			err = xml.Unmarshal(body, result)
		case "application/json":
			err = json.Unmarshal(body, result)
		default:
			err = fmt.Errorf("unsupport contentType:%s", contentType)
		}

		if err != nil {
			err = &DeserializationError{
				Err:      err,
				Snapshot: body,
			}
		}
	}
	return err
}

func unmarshalHeader(result interface{}, output *OperationOutput) error {
	val := reflect.ValueOf(result)
	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || output == nil {
		return nil
	}

	filedInfos := map[string]fieldInfo{}

	t := val.Type()
	for k := 0; k < t.NumField(); k++ {
		if tag, ok := t.Field(k).Tag.Lookup("output"); ok {
			tokens := strings.Split(tag, ",")
			if len(tokens) < 2 {
				continue
			}
			// header|query|body,filed_name,[required,time,usermeta...]
			switch tokens[0] {
			case "header":
				lowkey := strings.ToLower(tokens[1])
				var flags int = 0
				if len(tokens) >= 3 {
					flags = parseFiledFlags(tokens[2:])
				}
				filedInfos[lowkey] = fieldInfo{idx: k, flags: flags}
			}
		}
	}

	var err error
	for key, vv := range output.Headers {
		lkey := strings.ToLower(key)
		if field, ok := filedInfos[lkey]; ok {
			if field.flags&fTypeTime != 0 {
				if t, err := http.ParseTime(vv[0]); err == nil {
					err = setTimeReflectValue(val.Field(field.idx), t)
				}
			} else {
				err = setReflectValue(val.Field(field.idx), vv[0])
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func unmarshalHeaderLite(result interface{}, output *OperationOutput) error {
	val := reflect.ValueOf(result)
	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || output == nil {
		return nil
	}

	t := val.Type()
	for k := 0; k < t.NumField(); k++ {
		if tag := t.Field(k).Tag.Get("output"); tag != "" {
			tokens := strings.Split(tag, ",")
			if len(tokens) != 2 {
				continue
			}
			switch tokens[0] {
			case "header":
				if src := output.Headers.Get(tokens[1]); src != "" {
					if err := setReflectValue(val.Field(k), src); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (c *Client) unmarshalOutput(result interface{}, output *OperationOutput, handlers ...func(interface{}, *OperationOutput) error) error {
	// Common
	if cm, ok := result.(ResultCommonInterface); ok {
		cm.CopyIn(output.Status, output.StatusCode, output.Headers, output.OpMetadata)
	}

	var err error
	for _, h := range handlers {
		if err = h(result, output); err != nil {
			break
		}
	}
	return err
}

func updateContentMd5(input *OperationInput) error {
	var err error
	var contentMd5 string
	if input.Body != nil {
		var r io.ReadSeeker
		if r, ok := input.Body.(io.ReadSeeker); !ok {
			buf, _ := io.ReadAll(input.Body)
			r = bytes.NewReader(buf)
			input.Body = r
		}
		h := md5.New()
		if _, err = copySeekableBody(h, r); err != nil {
			// error
		} else {
			contentMd5 = base64.StdEncoding.EncodeToString(h.Sum(nil))
		}
	} else {
		contentMd5 = "1B2M2Y8AsgTpgAmY7PhCfg=="
	}

	// set content-md5 and content-type
	if err == nil {
		if input.Headers == nil {
			input.Headers = map[string]string{}
		}
		input.Headers["Content-Md5"] = contentMd5
	}

	return err
}

func (c *Client) toClientError(err error, code string, output *OperationOutput) error {
	if err == nil {
		return nil
	}

	return &ClientError{
		Code: code,
		Message: fmt.Sprintf("execute %s fail, error code is %s, request id:%s",
			output.Input.OpName,
			code,
			output.Headers.Get(HeaderOssRequestID),
		),
		Err: err}
}

// Content-Type
const (
	contentTypeDefault string = "application/octet-stream"
	contentTypeXML            = "application/xml"
)
