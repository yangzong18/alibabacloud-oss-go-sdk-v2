//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketWebsite(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketWebsiteRequest{
		Bucket: Ptr(bucketName),
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{
				Suffix:        Ptr("index.html"),
				SupportSubDir: Ptr(true),
				Type:          Ptr(int64(0)),
			},
			ErrorDocument: &ErrorDocument{
				Key:        Ptr("error.html"),
				HttpStatus: Ptr(int64(404)),
			},
		},
	}
	putResult, err := client.PutBucketWebsite(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketWebsiteRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketWebsite(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, getResult.WebsiteConfiguration, putRequest.WebsiteConfiguration)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketWebsiteRequest{
		Bucket: Ptr(bucketName),
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{
				Suffix:        Ptr("index.html"),
				SupportSubDir: Ptr(true),
				Type:          Ptr(int64(0)),
			},
			ErrorDocument: &ErrorDocument{
				Key:        Ptr("error.html"),
				HttpStatus: Ptr(int64(404)),
			},
			RoutingRules: &RoutingRules{
				[]RoutingRule{
					{
						Redirect: &RoutingRuleRedirect{
							MirrorPassOriginalSlashes: Ptr(false),
							RedirectType:              Ptr("Mirror"),
							MirrorURL:                 Ptr("http://example.com/"),
							MirrorPassQueryString:     Ptr(true),
							MirrorSNI:                 Ptr(true),
							ReplaceKeyPrefixWith:      Ptr("def/"),
							MirrorFollowRedirect:      Ptr(true),
							HostName:                  Ptr("example.com"),
							MirrorHeaders: &MirrorHeaders{
								Passes: []string{"myheader-key1", "myheader-key2"},
								Sets: []MirrorHeadersSet{
									{
										Key:   Ptr("myheader-key5"),
										Value: Ptr("myheader-value5"),
									},
								},
								PassAll: Ptr(true),
							},
							PassQueryString:                Ptr(true),
							EnableReplacePrefix:            Ptr(true),
							HttpRedirectCode:               Ptr(int64(301)),
							MirrorURLSlave:                 Ptr("http://example.com/"),
							MirrorSaveOssMeta:              Ptr(true),
							MirrorProxyPass:                Ptr(false),
							MirrorAllowGetImageInfo:        Ptr(true),
							MirrorAllowVideoSnapshot:       Ptr(false),
							MirrorIsExpressTunnel:          Ptr(true),
							MirrorDstRegion:                Ptr("cn-hangzhou"),
							MirrorUserLastModified:         Ptr(false),
							MirrorUsingRole:                Ptr(true),
							MirrorRole:                     Ptr("aliyun-test-role"),
							MirrorAllowHeadObject:          Ptr(true),
							TransparentMirrorResponseCodes: Ptr("400"),
							MirrorTaggings: &MirrorTaggings{
								Taggings: []MirrorTagging{
									{
										Key:   Ptr("k"),
										Value: Ptr("v"),
									},
								},
							},
							MirrorReturnHeaders: &MirrorReturnHeaders{
								ReturnHeaders: []ReturnHeader{
									{
										Key:   Ptr("k"),
										Value: Ptr("v"),
									},
								},
							},
							MirrorAuth: &MirrorAuth{
								AuthType:        Ptr("S3V4"),
								Region:          Ptr("ap-southeast-1"),
								AccessKeyId:     Ptr("TESTAK"),
								AccessKeySecret: Ptr("TESTSK"),
							},
						},
						RuleNumber: Ptr(int64(1)),
						Condition: &RoutingRuleCondition{
							KeySuffixEquals:             Ptr(".txt"),
							KeyPrefixEquals:             Ptr("abc/"),
							HttpErrorCodeReturnedEquals: Ptr(int64(404)),
						},
						LuaConfig: &RoutingRuleLuaConfig{
							Script: Ptr("test.lua"),
						},
					},
				},
			},
		},
	}
	putResult, err = client.PutBucketWebsite(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketWebsiteRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketWebsite(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketWebsiteRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketWebsite(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketWebsiteRequest{
		Bucket: Ptr(bucketNameNotExist),
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{
				Suffix:        Ptr("index.html"),
				SupportSubDir: Ptr(true),
				Type:          Ptr(int64(0)),
			},
			ErrorDocument: &ErrorDocument{
				Key:        Ptr("error.html"),
				HttpStatus: Ptr(int64(404)),
			},
		},
	}
	serr = &ServiceError{}
	putResult, err = client.PutBucketWebsite(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketWebsiteRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delResult, err = client.DeleteBucketWebsite(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
