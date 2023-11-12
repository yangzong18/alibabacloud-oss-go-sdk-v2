package oss

import (
	"net/http"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/retry"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Config struct {
	// The region in which the bucket is located.
	Region string

	// The domain names that other services can use to access OSS.
	Endpoint *string

	// RetryMaxAttempts specifies the maximum number attempts an API client will call
	// an operation that fails with a retryable error.
	RetryMaxAttempts int

	// Retryer guides how HTTP requests should be retried in case of recoverable failures.
	Retryer retry.Retryer

	// The HTTP client to invoke API calls with. Defaults to client's default HTTP
	// implementation if nil.
	HttpClient HTTPClient

	// The credentials provider to use when signing requests.
	CredentialsProvider credentials.CredentialsProvider

	// Allows you to enable the client to use path-style addressing, i.e., https://oss-cn-hangzhou.aliyuncs.com/bucket/key.
	// By default, the oss client will use virtual hosted addressing i.e., https://bucket.oss-cn-hangzhou.aliyuncs.com/key.
	UsePathStyle *bool

	// If the endpoint is s CName, set this flag to true
	UseCName *bool

	// Connect timeout
	ConnectTimeout *time.Duration

	// read & write timeout
	ReadWriteTimeout *time.Duration

	// Skip server certificate verification
	InsecureSkipVerify *bool

	// Enable http redirect or not. Default is disable
	EnabledRedirect *bool

	// Flag of using proxy host.
	ProxyHost *string

	// Read the proxy setting from the environment variables.
	// HTTP_PROXY, HTTPS_PROXY and NO_PROXY (or the lowercase versions thereof).
	// HTTPS_PROXY takes precedence over HTTP_PROXY for https requests.
	ProxyFromEnvironment *bool

	// Upload bandwidth limit in kBytes/s for all request
	UploadBandwidthlimit *int64

	// Download bandwidth limit in kBytes/s for all request
	DownloadBandwidthlimit *int64

	// Authentication with OSS Signature Version
	SignatureVersion SignatureVersionType
}

func NewConfig() *Config {
	return &Config{}
}

func (c Config) Copy() Config {
	cp := c
	return cp
}

func LoadDefaultConfig() *Config {
	config := &Config{
		RetryMaxAttempts: 3,
		SignatureVersion: SignatureVersionV1,
	}
	return config
}

func (c *Config) WithRegion(region string) *Config {
	c.Region = region
	return c
}

func (c *Config) WithEndpoint(endpoint string) *Config {
	c.Endpoint = Ptr(endpoint)
	return c
}

func (c *Config) WithRetryMaxAttempts(value int) *Config {
	c.RetryMaxAttempts = value
	return c
}

func (c *Config) WithRetryer(retryer retry.Retryer) *Config {
	c.Retryer = retryer
	return c
}

func (c *Config) WithHttpClient(client *http.Client) *Config {
	c.HttpClient = client
	return c
}

func (c *Config) WithCredentialsProvider(provider credentials.CredentialsProvider) *Config {
	c.CredentialsProvider = provider
	return c
}

func (c *Config) WithUsePathStyle(enable bool) *Config {
	c.UsePathStyle = Ptr(enable)
	return c
}

func (c *Config) WithUseCName(enable bool) *Config {
	c.UseCName = Ptr(enable)
	return c
}

func (c *Config) WithConnectTimeout(value time.Duration) *Config {
	c.ConnectTimeout = Ptr(value)
	return c
}

func (c *Config) WithReadWriteTimeout(value time.Duration) *Config {
	c.ReadWriteTimeout = Ptr(value)
	return c
}

func (c *Config) WithInsecureSkipVerify(value bool) *Config {
	c.InsecureSkipVerify = Ptr(value)
	return c
}

func (c *Config) WithEnabledRedirect(value bool) *Config {
	c.EnabledRedirect = Ptr(value)
	return c
}

func (c *Config) WithProxyHost(value string) *Config {
	c.ProxyHost = Ptr(value)
	return c
}

func (c *Config) WithProxyFromEnvironment(value bool) *Config {
	c.ProxyFromEnvironment = Ptr(value)
	return c
}

func (c *Config) WithUploadBandwidthlimit(value int64) *Config {
	c.UploadBandwidthlimit = Ptr(value)
	return c
}

func (c *Config) WithDownloadBandwidthlimit(value int64) *Config {
	c.DownloadBandwidthlimit = Ptr(value)
	return c
}

func (c *Config) WithSignatureVersion(value SignatureVersionType) *Config {
	c.SignatureVersion = value
	return c
}
