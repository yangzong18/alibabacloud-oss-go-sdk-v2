package oss

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/readers"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/retry"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/signer"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/transport"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/types"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/util"
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

func New(cfg *Config, optFns ...func(*Options)) *Client {
	options := Options{
		Region:              cfg.Region,
		RetryMaxAttempts:    cfg.RetryMaxAttempts,
		Retryer:             cfg.Retryer,
		CredentialsProvider: cfg.CredentialsProvider,
		HttpClient:          cfg.HTTPClient,
	}
	resolveEndpoint(cfg, &options)
	resolveRetryer(cfg, &options)
	resolveHTTPClient(cfg, &options)
	resolveSigner(cfg, &options)

	for _, fn := range optFns {
		fn(&options)
	}

	client := &Client{
		options: options,
	}

	return client
}

func resolveEndpoint(cfg *Config, o *Options) {
	var scheme string
	var endpoint string
	if strings.HasPrefix(cfg.Endpoint, "http://") {
		scheme = "http"
		endpoint = cfg.Endpoint[len("http://"):]
	} else if strings.HasPrefix(endpoint, "https://") {
		scheme = "https"
		endpoint = cfg.Endpoint[len("https://"):]
	} else {
		scheme = "http"
		endpoint = cfg.Endpoint
	}

	strUrl := fmt.Sprintf("%s://%s", scheme, endpoint)
	o.Endpoint, _ = url.Parse(strUrl)
}

func resolveRetryer(cfg *Config, o *Options) {
	if o.Retryer != nil {
		return
	}
	retryMode := cfg.RetryMode
	if len(retryMode) == 0 {
		retryMode = retry.RetryModeStandard
	}
	switch retryMode {
	case retry.RetryModeAdaptive:
		o.Retryer = retry.NopRetryer{}
	default:
		o.Retryer = retry.NopRetryer{}
	}
}

func resolveHTTPClient(cfg *Config, o *Options) {
	if o.HttpClient != nil {
		return
	}

	//TODO timeouts from config

	o.HttpClient = &http.Client{
		Transport: transport.NewTransportCustom(),
	}
}

func resolveSigner(cfg *Config, o *Options) {
	if o.Signer != nil {
		return
	}

	o.Signer = signer.SignerV1{}
}

type OperationInput struct {
	OperationName string

	Bucket string
	Key    string

	Method     string
	Headers    map[string]string
	Parameters map[string]string
	Body       io.Reader

	Metadata types.ApiMetadata
}

type OperationOutput struct {
	input *OperationInput

	Status     string
	StatusCode int
	Headers    http.Header
	Body       io.ReadCloser

	Metadata types.ApiMetadata
}

func (c *Client) InvokeOperation(ctx context.Context, input *OperationInput, optFns ...func(*Options)) (output *OperationOutput, err error) {
	options := c.options.Copy()

	for _, fn := range optFns {
		fn(&options)
	}

	//finalize retry

	//finalize endpoint

	//default reponse handler

	output, err = c.sendRequest(ctx, input, &options)

	if err != nil {
		return output, &types.OperationError{
			OperationName: input.OperationName,
			Err:           err}
	}

	return output, err
}

func (c *Client) sendRequest(ctx context.Context, input *OperationInput, opts *Options) (output *OperationOutput, err error) {
	// covert input into httpRequest
	if opts.Endpoint == nil {
		return output, errors.New("Endpoint is nil")
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
	subResource, _ := input.Metadata.Get(signer.SubResource).([]string)
	signingCtx := &signer.SigningContext{
		Product:     "oss",
		Region:      opts.Region,
		Bucket:      input.Bucket,
		Key:         input.Key,
		Request:     request,
		SubResource: subResource,
	}

	// send request
	response, err := c.signAndSendRequest(ctx, signingCtx, opts)

	if err != nil {
		return output, err
	}

	// covert http response into output context
	output = &OperationOutput{
		input:      input,
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Body:       response.Body,
		Headers:    response.Header,
	}

	// save other info by Metadata filed, ex. retry detail info
	//output.Metadata.Set()

	return output, err
}

func (c *Client) signAndSendRequest(ctx context.Context, signingCtx *signer.SigningContext, opts *Options) (response *http.Response, err error) {
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
			if err = util.SleepWithContext(ctx, delay); err != nil {
				err = &types.CanceledError{Err: err}
				break
			}
		}

		response, err = c.signAndSendRequestOnce(ctx, signingCtx, opts)

		if err == nil {
			break
		}

		if types.ContextError(ctx, &err) {
			err = &types.CanceledError{Err: err}
			break
		}

		if !retryer.IsErrorRetryable(err) {
			break
		}

		if !readers.IsReaderSeekable(request.Body) {
			break
		}

		_, err = body.Seek(bodyStart, io.SeekStart)
		if err != nil {
			break
		}
	}
	return response, err
}

func (c *Client) signAndSendRequestOnce(ctx context.Context, signingCtx *signer.SigningContext, opts *Options) (
	response *http.Response, err error,
) {
	cred, err := opts.CredentialsProvider.GetCredentials(ctx)
	if err != nil {
		return response, err
	}
	signingCtx.Credentials = &cred
	err = c.options.Signer.Sign(ctx, signingCtx)
	if err != nil {
		return response, err
	}
	response, err = c.options.HttpClient.Do(signingCtx.Request)

	if err != nil {
		return response, err
	}

	err = handleResponseServiceError(response)

	if err != nil {
		return response, err
	}

	for _, fn := range opts.ResponseHandlers {
		err = fn(response)
		if err != nil {
			return response, err
		}
	}

	return response, err
}

func buildURL(input *OperationInput, opts *Options) (string, string) {
	var host = ""
	var path = ""

	if input == nil || opts == nil || opts.Endpoint == nil {
		return host, path
	}

	bucket := input.Bucket
	object := util.EscapePath(input.Key, false)

	if bucket == "" {
		host = opts.Endpoint.Host
		path = "/"
	} else {
		host = bucket + "." + opts.Endpoint.Host
		path = "/" + object
	}

	return host, path
}

func handleResponseServiceError(response *http.Response) error {
	if response.StatusCode/100 == 2 {
		return nil
	}

	timestamp, err := time.Parse(http.TimeFormat, response.Header.Get("Date"))
	if err != nil {
		timestamp = util.NowTime()
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	se := types.ServiceError{
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

func defaultUserAgent() string {
	return fmt.Sprintf("aliyun-sdk-go/%s (%s/%s/%s;%s)", Version(), runtime.GOOS,
		"-", runtime.GOARCH, runtime.Version())
}
