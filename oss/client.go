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
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
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

	HttpClient HTTPClient

	ResponseHandlers []func(*http.Response) error

	UrlStyle UrlStyleType

	FeatureFlags FeatureFlagsType

	OpReadWriteTimeout *time.Duration

	AuthMethod *AuthMethodType

	AdditionalHeaders []string
}

func (c Options) Copy() Options {
	to := c
	to.ResponseHandlers = make([]func(*http.Response) error, len(c.ResponseHandlers))
	copy(to.ResponseHandlers, c.ResponseHandlers)
	return to
}

type innerOptions struct {
	BwTokenBuckets BwTokenBuckets

	// A clock offset that how much client time is different from server time
	ClockOffset time.Duration
}

type Client struct {
	options Options
	inner   innerOptions
}

func NewClient(cfg *Config, optFns ...func(*Options)) *Client {
	options := Options{
		Region:              cfg.Region,
		RetryMaxAttempts:    cfg.RetryMaxAttempts,
		Retryer:             cfg.Retryer,
		CredentialsProvider: cfg.CredentialsProvider,
		HttpClient:          cfg.HttpClient,
		FeatureFlags:        FeatureFlagsDefault,
	}
	inner := innerOptions{}

	resolveEndpoint(cfg, &options)
	resolveRetryer(cfg, &options)
	resolveHTTPClient(cfg, &options, &inner)
	resolveSigner(cfg, &options)
	resolveUrlStyle(cfg, &options)
	resolveFeatureFlags(cfg, &options)

	for _, fn := range optFns {
		fn(&options)
	}

	client := &Client{
		options: options,
		inner:   inner,
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

func resolveHTTPClient(cfg *Config, o *Options, inner *innerOptions) {
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
	if cfg.UploadBandwidthlimit != nil {
		value := *cfg.UploadBandwidthlimit * 1024
		tb := newBwTokenBucket(value)
		tcfg.PostWrite = append(tcfg.PostWrite, func(n int, _ error) {
			tb.LimitBandwidth(n)
		})
		inner.BwTokenBuckets[BwTokenBucketSlotTx] = tb
	}
	if cfg.DownloadBandwidthlimit != nil {
		value := *cfg.DownloadBandwidthlimit * 1024
		tb := newBwTokenBucket(value)
		tcfg.PostRead = append(tcfg.PostRead, func(n int, _ error) {
			tb.LimitBandwidth(n)
		})
		inner.BwTokenBuckets[BwTokenBucketSlotRx] = tb
	}

	o.HttpClient = transport.NewHttpClient(tcfg, custom...)
}

func resolveSigner(cfg *Config, o *Options) {
	if o.Signer != nil {
		return
	}

	switch cfg.SignatureVersion {
	case SignatureVersionV1:
		o.Signer = &signer.SignerV1{}
	default:
		o.Signer = &signer.SignerV4{}
	}
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

	applyOperationMetadata(input, &options)

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

	var writers []io.Writer
	// tracker in OpertionMetaData
	if trackers, ok := input.OpMetadata.Get(OpMetaKeyRequestBodyTracker).([]io.Writer); ok {
		writers = append(writers, trackers...)
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
	var body io.Reader
	if input.Body == nil {
		body = strings.NewReader("")
	} else {
		body = input.Body
	}
	length := GetReaderLen(body)
	if length >= 0 && request.Header.Get("Content-Length") == "" {
		request.ContentLength = length
	}
	request.Body = TeeReadNopCloser(body, writers...)

	//signing context
	subResource, _ := input.OpMetadata.Get(signer.SubResource).([]string)
	clockOffset := c.inner.ClockOffset
	signingCtx := &signer.SigningContext{
		Product:         Ptr("oss"),
		Region:          Ptr(opts.Region),
		Bucket:          input.Bucket,
		Key:             input.Key,
		Request:         request,
		SubResource:     subResource,
		AuthMethodQuery: opts.AuthMethod != nil && *opts.AuthMethod == AuthMethodQuery,
		ClockOffset:     clockOffset,
	}

	if date := request.Header.Get(HeaderOssDate); date != "" {
		signingCtx.Time, _ = http.ParseTime(date)
	} else if signTime, ok := input.OpMetadata.Get(signer.SignTime).(time.Time); ok {
		signingCtx.Time = signTime
	}

	// send http request
	response, err := c.sendHttpRequest(ctx, signingCtx, opts)

	if err != nil {
		return output, err
	}

	// covert http response into output context
	output = &OperationOutput{
		Input:       input,
		Status:      response.Status,
		StatusCode:  response.StatusCode,
		Body:        response.Body,
		Headers:     response.Header,
		httpRequest: request,
	}

	// save other info by Metadata filed, ex. retry detail info
	//output.OpMetadata.Set(...)
	if signingCtx.AuthMethodQuery {
		output.OpMetadata.Set(signer.SignTime, signingCtx.Time)
	}

	if signingCtx.ClockOffset != clockOffset {
		c.inner.ClockOffset = signingCtx.ClockOffset
	}

	return output, err
}

func (c *Client) sendHttpRequest(ctx context.Context, signingCtx *signer.SigningContext, opts *Options) (response *http.Response, err error) {
	request := signingCtx.Request
	retryer := opts.Retryer
	body, _ := request.Body.(*teeReadNopCloser)
	body.Mark()
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

			if err = body.Reset(); err != nil {
				break
			}
		}

		if response, err = c.sendHttpRequestOnce(ctx, signingCtx, opts); err == nil {
			break
		}

		c.postSendHttpRequestOnce(signingCtx, response, err)

		if isContextError(ctx, &err) {
			err = &CanceledError{Err: err}
			break
		}

		if !body.IsSeekable() {
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

	if response, err = opts.HttpClient.Do(signingCtx.Request); err != nil {
		return response, err
	}

	for _, fn := range opts.ResponseHandlers {
		if err = fn(response); err != nil {
			return response, err
		}
	}

	return response, err
}

func (c *Client) postSendHttpRequestOnce(signingCtx *signer.SigningContext, response *http.Response, err error) {
	if err != nil {
		switch e := err.(type) {
		case *ServiceError:
			if c.hasFeature(FeatureCorrectClockSkew) &&
				e.Code == "RequestTimeTooSkewed" &&
				!e.Timestamp.IsZero() {
				signingCtx.ClockOffset = e.Timestamp.Sub(signingCtx.Time)
			}
		}
	}
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

func callbackResponseHandler(opts *Options) {
	opts.ResponseHandlers = []func(*http.Response) error{
		callbackErrorResponseHandler,
	}
}

func serviceErrorResponseHandler(response *http.Response) error {
	if response.StatusCode/100 == 2 {
		return nil
	}
	return tryConvertServiceError(response)
}

func callbackErrorResponseHandler(response *http.Response) error {
	if response.StatusCode == 200 {
		return nil
	}
	return tryConvertServiceError(response)
}

func tryConvertServiceError(response *http.Response) (err error) {
	var respBody []byte
	var body []byte
	timestamp, err := time.Parse(http.TimeFormat, response.Header.Get("Date"))
	if err != nil {
		timestamp = time.Now()
	}

	defer response.Body.Close()
	respBody, err = io.ReadAll(response.Body)
	body = respBody
	if len(respBody) == 0 && len(response.Header.Get(HeaderOssERR)) > 0 {
		body, err = base64.StdEncoding.DecodeString(response.Header.Get(HeaderOssERR))
		if err != nil {
			body = respBody
		}
	}
	se := &ServiceError{
		StatusCode:    response.StatusCode,
		Code:          "BadErrorResponse",
		RequestID:     response.Header.Get(HeaderOssRequestID),
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

	if op.HttpClient != nil {
		c.HttpClient = op.HttpClient
	}

	if op.AuthMethod != nil {
		c.AuthMethod = op.AuthMethod
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

func applyOperationMetadata(input *OperationInput, c *Options) {
	if handles, ok := input.OpMetadata.Get(OpMetaKeyResponsHandler).([]func(*http.Response) error); ok {
		c.ResponseHandlers = append(c.ResponseHandlers, handles...)
	}
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
	fTypeReader
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
		case "reader":
			flags |= fTypeReader
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

func (c *Client) marshalInput(request any, input *OperationInput, handlers ...func(any, *OperationInput) error) error {
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
				} else {
					if r, ok := v.Interface().(io.Reader); ok {
						input.Body = r
					} else {
						return NewErrParamTypeNotSupport(t.Field(k).Name)
					}
				}
			}
		}
	}

	if err := validateInput(input); err != nil {
		return err
	}

	for _, h := range handlers {
		if err := h(request, input); err != nil {
			return err
		}
	}

	return nil
}

func marshalDeleteObjects(request any, input *OperationInput) error {
	var builder strings.Builder
	delRequest := request.(*DeleteMultipleObjectsRequest)
	builder.WriteString("<Delete>")
	builder.WriteString("<Quiet>")
	builder.WriteString(strconv.FormatBool(delRequest.Quiet))
	builder.WriteString("</Quiet>")
	if len(delRequest.Objects) > 0 {
		for _, object := range delRequest.Objects {
			builder.WriteString("<Object>")
			if object.Key != nil {
				builder.WriteString("<Key>")
				builder.WriteString(escapeXml(*object.Key))
				builder.WriteString("</Key>")
			}
			if object.VersionId != nil {
				builder.WriteString("<VersionId>")
				builder.WriteString(*object.VersionId)
				builder.WriteString("</VersionId>")
			}
			builder.WriteString("</Object>")
		}
	} else {
		return NewErrParamInvalid("Objects")
	}
	builder.WriteString("</Delete>")
	input.Body = strings.NewReader(builder.String())
	return nil
}

func discardBody(result any, output *OperationOutput) error {
	var err error
	if output.Body != nil {
		defer output.Body.Close()
		_, err = io.Copy(io.Discard, output.Body)
	}
	return err
}

func unmarshalBodyXml(result any, output *OperationOutput) error {
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

func unmarshalBodyDefault(result any, output *OperationOutput) error {
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

func unmarshalCallbackBody(result any, output *OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}
	if len(body) > 0 {
		switch r := result.(type) {
		case *PutObjectResult:
			if err = json.Unmarshal(body, &r.CallbackResult); err != nil {
				return err
			}
		case *CompleteMultipartUploadResult:
			if err = json.Unmarshal(body, &r.CallbackResult); err != nil {
				return err
			}
		}
	}
	return err
}

func unmarshalHeader(result any, output *OperationOutput) error {
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
	var usermetaKeys []string
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
				if flags&fTypeUsermeta != 0 {
					usermetaKeys = append(usermetaKeys, lowkey)
				}
			}
		}
	}
	var err error
	for key, vv := range output.Headers {
		lkey := strings.ToLower(key)
		for _, prefix := range usermetaKeys {
			if strings.HasPrefix(lkey, prefix) {
				if field, ok := filedInfos[prefix]; ok {
					if field.flags&fTypeUsermeta != 0 {
						mapKey := strings.TrimPrefix(lkey, prefix)
						err = setMapStringReflectValue(val.Field(field.idx), mapKey, vv[0])
					}
				}
			}
		}
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

func unmarshalHeaderLite(result any, output *OperationOutput) error {
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

func (c *Client) unmarshalOutput(result any, output *OperationOutput, handlers ...func(any, *OperationOutput) error) error {
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

func updateContentMd5(_ any, input *OperationInput) error {
	var err error
	var contentMd5 string
	if input.Body != nil {
		var r io.ReadSeeker
		var ok bool
		if r, ok = input.Body.(io.ReadSeeker); !ok {
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
		input.Headers["Content-MD5"] = contentMd5
	}

	return err
}

func updateContentType(_ any, input *OperationInput) error {
	if input.Headers == nil {
		input.Headers = map[string]string{}
	}

	if _, ok := input.Headers[HTTPHeaderContentType]; !ok {
		value := TypeByExtension(ToString(input.Key))
		if value == "" {
			value = contentTypeDefault
		}
		input.Headers[HTTPHeaderContentType] = value
	}
	return nil
}

func encodeSourceObject(request any) string {
	var bucket, key, versionId string
	switch req := request.(type) {
	case *CopyObjectRequest:
		key = ToString(req.SourceKey)
		if req.SourceBucket != nil {
			bucket = *req.SourceBucket
		} else {
			bucket = ToString(req.Bucket)
		}
		versionId = ToString(req.SourceVersionId)
	case *UploadPartCopyRequest:
		key = ToString(req.SourceKey)
		if req.SourceBucket != nil {
			bucket = *req.SourceBucket
		} else {
			bucket = ToString(req.Bucket)
		}
		versionId = ToString(req.SourceVersionId)
	}

	source := fmt.Sprintf("/%s/%s", bucket, escapePath(key, false))
	if versionId != "" {
		source += "?versionId=" + versionId
	}

	return source
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

func (c *Client) hasFeature(flag FeatureFlagsType) bool {
	return (c.options.FeatureFlags & flag) > 0
}

// Content-Type
const (
	contentTypeDefault string = "application/octet-stream"
	contentTypeXML            = "application/xml"
)
