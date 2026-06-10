package oss

import (
	"testing"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutBucketReplicationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketReplicationRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketReplicationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=add&replication", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationConfiguration><Rule><Destination><Bucket>destbucket</Bucket><Location>oss-cn-beijing</Location><TransferType>oss_acc</TransferType></Destination><SyncRole>aliyunramrole</SyncRole><SourceSelectionCriteria><SseKmsEncryptedObjects><Status>Enabled</Status></SseKmsEncryptedObjects></SourceSelectionCriteria><EncryptionConfiguration><ReplicaKmsKeyID>c4d49f85-ee30-426b-a5ed-95e9139d****</ReplicaKmsKeyID></EncryptionConfiguration><HistoricalObjectReplication>enabled</HistoricalObjectReplication><RTC><Status>enabled</Status></RTC></Rule></ReplicationConfiguration>")
		},
		&PutBucketReplicationRequest{
			Bucket: Ptr("bucket"),
			ReplicationConfiguration: &ReplicationConfiguration{
				[]ReplicationRule{
					{
						RTC: &ReplicationTimeControl{
							Status: Ptr("enabled"),
						},
						Destination: &ReplicationDestination{
							Bucket:       Ptr("destbucket"),
							Location:     Ptr("oss-cn-beijing"),
							TransferType: TransferTypeOssAcc,
						},
						HistoricalObjectReplication: HistoricalObjectReplicationEnabled,
						SyncRole:                    Ptr("aliyunramrole"),
						SourceSelectionCriteria: &ReplicationSourceSelectionCriteria{
							&SseKmsEncryptedObjects{
								Status: StatusEnabled,
							},
						},
						EncryptionConfiguration: &ReplicationEncryptionConfiguration{
							Ptr("c4d49f85-ee30-426b-a5ed-95e9139d****"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketReplicationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=add&replication", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationConfiguration><Rule><Destination><Bucket>destbucket</Bucket><Location>oss-cn-beijing</Location><TransferType>oss_acc</TransferType></Destination><HistoricalObjectReplication>enabled</HistoricalObjectReplication></Rule></ReplicationConfiguration>")
		},
		&PutBucketReplicationRequest{
			Bucket: Ptr("bucket"),
			ReplicationConfiguration: &ReplicationConfiguration{
				[]ReplicationRule{
					{
						Destination: &ReplicationDestination{
							Bucket:       Ptr("destbucket"),
							Location:     Ptr("oss-cn-beijing"),
							TransferType: TransferTypeOssAcc,
						},
						HistoricalObjectReplication: HistoricalObjectReplicationEnabled,
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketReplicationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketReplication_Success(t *testing.T) {
	for _, c := range testMockPutBucketReplicationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketReplication(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketReplicationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketReplicationRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketReplicationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=add&replication", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationConfiguration><Rule><Destination><Bucket>destbucket</Bucket><Location>oss-cn-beijing</Location><TransferType>oss_acc</TransferType></Destination><HistoricalObjectReplication>enabled</HistoricalObjectReplication></Rule></ReplicationConfiguration>")
		},
		&PutBucketReplicationRequest{
			Bucket: Ptr("bucket"),
			ReplicationConfiguration: &ReplicationConfiguration{
				[]ReplicationRule{
					{
						Destination: &ReplicationDestination{
							Bucket:       Ptr("destbucket"),
							Location:     Ptr("oss-cn-beijing"),
							TransferType: TransferTypeOssAcc,
						},
						HistoricalObjectReplication: HistoricalObjectReplicationEnabled,
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketReplicationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=add&replication", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationConfiguration><Rule><Destination><Bucket>destbucket</Bucket><Location>oss-cn-beijing</Location><TransferType>oss_acc</TransferType></Destination><HistoricalObjectReplication>enabled</HistoricalObjectReplication></Rule></ReplicationConfiguration>")
		},
		&PutBucketReplicationRequest{
			Bucket: Ptr("bucket"),
			ReplicationConfiguration: &ReplicationConfiguration{
				[]ReplicationRule{
					{
						Destination: &ReplicationDestination{
							Bucket:       Ptr("destbucket"),
							Location:     Ptr("oss-cn-beijing"),
							TransferType: TransferTypeOssAcc,
						},
						HistoricalObjectReplication: HistoricalObjectReplicationEnabled,
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketReplicationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutBucketReplication_Error(t *testing.T) {
	for _, c := range testMockPutBucketReplicationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketReplication(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketRtcSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRtcRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketRtcResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?rtc", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationRule><RTC><Status>disabled</Status></RTC><ID>test_replication_rule_1</ID></ReplicationRule>")
		},
		&PutBucketRtcRequest{
			Bucket: Ptr("bucket"),
			RtcConfiguration: &RtcConfiguration{
				RTC: &ReplicationTimeControl{
					Status: Ptr("disabled"),
				},
				ID: Ptr("test_replication_rule_1"),
			},
		},
		func(t *testing.T, o *PutBucketRtcResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?rtc", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationRule><RTC><Status>enabled</Status></RTC><ID>test_replication_rule_1</ID></ReplicationRule>")
		},
		&PutBucketRtcRequest{
			Bucket: Ptr("bucket"),
			RtcConfiguration: &RtcConfiguration{
				RTC: &ReplicationTimeControl{
					Status: Ptr("enabled"),
				},
				ID: Ptr("test_replication_rule_1"),
			},
		},
		func(t *testing.T, o *PutBucketRtcResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketRtc_Success(t *testing.T) {
	for _, c := range testMockPutBucketRtcSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketRtc(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketRtcErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRtcRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketRtcResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?rtc", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationRule><RTC><Status>enabled</Status></RTC><ID>test_replication_rule_1</ID></ReplicationRule>")
		},
		&PutBucketRtcRequest{
			Bucket: Ptr("bucket"),
			RtcConfiguration: &RtcConfiguration{
				RTC: &ReplicationTimeControl{
					Status: Ptr("enabled"),
				},
				ID: Ptr("test_replication_rule_1"),
			},
		},
		func(t *testing.T, o *PutBucketRtcResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?rtc", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationRule><RTC><Status>enabled</Status></RTC><ID>test_replication_rule_1</ID></ReplicationRule>")
		},
		&PutBucketRtcRequest{
			Bucket: Ptr("bucket"),
			RtcConfiguration: &RtcConfiguration{
				RTC: &ReplicationTimeControl{
					Status: Ptr("enabled"),
				},
				ID: Ptr("test_replication_rule_1"),
			},
		},
		func(t *testing.T, o *PutBucketRtcResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutBucketRtc_Error(t *testing.T) {
	for _, c := range testMockPutBucketRtcErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketRtc(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketReplicationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketReplicationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketReplicationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<ReplicationConfiguration>
  <Rule>
    <ID>test_replication_1</ID>
    <PrefixSet>
      <Prefix>source1</Prefix>
      <Prefix>video</Prefix>
    </PrefixSet>
    <Action>PUT</Action>
    <Destination>
      <Bucket>destbucket</Bucket>
      <Location>oss-cn-beijing</Location>
      <TransferType>oss_acc</TransferType>
    </Destination>
    <Status>doing</Status>
    <HistoricalObjectReplication>enabled</HistoricalObjectReplication>
    <SyncRole>aliyunramrole</SyncRole>
    <RTC>
      <Status>enabled</Status>
    </RTC>
  </Rule>
</ReplicationConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replication", urlStr)
		},
		&GetBucketReplicationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketReplicationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.ReplicationConfiguration.Rules), 1)
			assert.Equal(t, *o.ReplicationConfiguration.Rules[0].RTC.Status, "enabled")
			assert.Equal(t, *o.ReplicationConfiguration.Rules[0].ID, "test_replication_1")
			assert.Equal(t, o.ReplicationConfiguration.Rules[0].PrefixSet.Prefixs[0], "source1")
			assert.Equal(t, o.ReplicationConfiguration.Rules[0].PrefixSet.Prefixs[1], "video")
			assert.Equal(t, *o.ReplicationConfiguration.Rules[0].Action, "PUT")
			assert.Equal(t, o.ReplicationConfiguration.Rules[0].Destination.TransferType, TransferTypeOssAcc)
			assert.Equal(t, *o.ReplicationConfiguration.Rules[0].Destination.Bucket, "destbucket")
			assert.Equal(t, *o.ReplicationConfiguration.Rules[0].Destination.Location, "oss-cn-beijing")
			assert.Equal(t, *o.ReplicationConfiguration.Rules[0].Status, "doing")
			assert.Equal(t, o.ReplicationConfiguration.Rules[0].HistoricalObjectReplication, HistoricalObjectReplicationEnabled)
			assert.Equal(t, *o.ReplicationConfiguration.Rules[0].SyncRole, "aliyunramrole")
		},
	},
}

func TestMockGetBucketReplication_Success(t *testing.T) {
	for _, c := range testMockGetBucketReplicationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketReplication(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketReplicationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketReplicationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketReplicationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replication", urlStr)
		},
		&GetBucketReplicationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketReplicationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replication", urlStr)
		},
		&GetBucketReplicationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketReplicationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replication", urlStr)
		},
		&GetBucketReplicationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketReplicationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketReplication fail")
		},
	},
}

func TestMockGetBucketReplication_Error(t *testing.T) {
	for _, c := range testMockGetBucketReplicationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketReplication(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketReplicationLocationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketReplicationLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketReplicationLocationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<ReplicationLocation>
  <Location>oss-cn-beijing</Location>
  <Location>oss-cn-qingdao</Location>
  <Location>oss-cn-shenzhen</Location>
  <Location>oss-cn-hongkong</Location>
  <Location>oss-us-west-1</Location>
  <LocationTransferTypeConstraint>
    <LocationTransferType>
      <Location>oss-cn-hongkong</Location>
        <TransferTypes>
          <Type>oss_acc</Type>          
        </TransferTypes>
      </LocationTransferType>
      <LocationTransferType>
        <Location>oss-us-west-1</Location>
        <TransferTypes>
          <Type>oss_acc</Type>
        </TransferTypes>
      </LocationTransferType>
    </LocationTransferTypeConstraint>
  </ReplicationLocation>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replicationLocation", urlStr)
		},
		&GetBucketReplicationLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketReplicationLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.ReplicationLocation.Locations), 5)
			assert.Equal(t, o.ReplicationLocation.Locations[0], "oss-cn-beijing")
			assert.Equal(t, o.ReplicationLocation.Locations[1], "oss-cn-qingdao")
			assert.Equal(t, o.ReplicationLocation.Locations[2], "oss-cn-shenzhen")
			assert.Equal(t, o.ReplicationLocation.Locations[3], "oss-cn-hongkong")
			assert.Equal(t, o.ReplicationLocation.Locations[4], "oss-us-west-1")
			assert.Equal(t, len(o.ReplicationLocation.LocationTransferTypeConstraint.LocationTransferTypes), 2)
			assert.Equal(t, *o.ReplicationLocation.LocationTransferTypeConstraint.LocationTransferTypes[0].Location, "oss-cn-hongkong")
			assert.Equal(t, o.ReplicationLocation.LocationTransferTypeConstraint.LocationTransferTypes[0].TransferTypes.Types[0], "oss_acc")
			assert.Equal(t, *o.ReplicationLocation.LocationTransferTypeConstraint.LocationTransferTypes[1].Location, "oss-us-west-1")
			assert.Equal(t, o.ReplicationLocation.LocationTransferTypeConstraint.LocationTransferTypes[1].TransferTypes.Types[0], "oss_acc")
		},
	},
}

func TestMockGetBucketReplicationLocation_Success(t *testing.T) {
	for _, c := range testMockGetBucketReplicationLocationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketReplicationLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketReplicationLocationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketReplicationLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketReplicationLocationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replicationLocation", urlStr)
		},
		&GetBucketReplicationLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketReplicationLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replicationLocation", urlStr)
		},
		&GetBucketReplicationLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketReplicationLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replicationLocation", urlStr)
		},
		&GetBucketReplicationLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketReplicationLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketReplicationLocation fail")
		},
	},
}

func TestMockGetBucketReplicationLocation_Error(t *testing.T) {
	for _, c := range testMockGetBucketReplicationLocationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketReplicationLocation(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketReplicationProgressSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketReplicationProgressRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketReplicationProgressResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<ReplicationProgress>
 <Rule>
   <ID>test_replication_1</ID>
   <PrefixSet>
    <Prefix>source_image</Prefix>
    <Prefix>video</Prefix>
   </PrefixSet>
   <Action>PUT</Action>
   <Destination>
    <Bucket>target-bucket</Bucket>
    <Location>oss-cn-beijing</Location>
    <TransferType>oss_acc</TransferType>
   </Destination>
   <Status>doing</Status>
   <HistoricalObjectReplication>enabled</HistoricalObjectReplication>
   <Progress>
    <HistoricalObject>0.85</HistoricalObject>
    <NewObject>2015-09-24T15:28:14.000Z</NewObject>
   </Progress>
 </Rule>
</ReplicationProgress>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replicationProgress&rule-id=test_replication_1", urlStr)
		},
		&GetBucketReplicationProgressRequest{
			Bucket: Ptr("bucket"),
			RuleId: Ptr("test_replication_1"),
		},
		func(t *testing.T, o *GetBucketReplicationProgressResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, len(o.ReplicationProgress.Rules), 1)
			assert.Equal(t, *o.ReplicationProgress.Rules[0].ID, "test_replication_1")
			assert.Equal(t, o.ReplicationProgress.Rules[0].PrefixSet.Prefixs[0], "source_image")
			assert.Equal(t, o.ReplicationProgress.Rules[0].PrefixSet.Prefixs[1], "video")
			assert.Equal(t, *o.ReplicationProgress.Rules[0].Action, "PUT")
			assert.Equal(t, o.ReplicationProgress.Rules[0].Destination.TransferType, TransferTypeOssAcc)
			assert.Equal(t, *o.ReplicationProgress.Rules[0].Destination.Bucket, "target-bucket")
			assert.Equal(t, *o.ReplicationProgress.Rules[0].Destination.Location, "oss-cn-beijing")
			assert.Equal(t, *o.ReplicationProgress.Rules[0].Status, "doing")
			assert.Equal(t, *o.ReplicationProgress.Rules[0].HistoricalObjectReplication, "enabled")
			assert.Equal(t, *o.ReplicationProgress.Rules[0].Progress.HistoricalObject, "0.85")
			assert.Equal(t, *o.ReplicationProgress.Rules[0].Progress.NewObject, "2015-09-24T15:28:14.000Z")
		},
	},
}

func TestMockGetBucketReplicationProgress_Success(t *testing.T) {
	for _, c := range testMockGetBucketReplicationProgressSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketReplicationProgress(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketReplicationProgressErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketReplicationProgressRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketReplicationProgressResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replicationProgress&rule-id=test_replication_1", urlStr)
		},
		&GetBucketReplicationProgressRequest{
			Bucket: Ptr("bucket"),
			RuleId: Ptr("test_replication_1"),
		},
		func(t *testing.T, o *GetBucketReplicationProgressResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replicationProgress&rule-id=test_replication_1", urlStr)
		},
		&GetBucketReplicationProgressRequest{
			Bucket: Ptr("bucket"),
			RuleId: Ptr("test_replication_1"),
		},
		func(t *testing.T, o *GetBucketReplicationProgressResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?replicationProgress&rule-id=test_replication_1", urlStr)
		},
		&GetBucketReplicationProgressRequest{
			Bucket: Ptr("bucket"),
			RuleId: Ptr("test_replication_1"),
		},
		func(t *testing.T, o *GetBucketReplicationProgressResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketReplicationProgress fail")
		},
	},
}

func TestMockGetBucketReplicationProgress_Error(t *testing.T) {
	for _, c := range testMockGetBucketReplicationProgressErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketReplicationProgress(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketReplicationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketReplicationRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketReplicationResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=delete&replication", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationRules><ID>test_replication_1</ID></ReplicationRules>")
		},
		&DeleteBucketReplicationRequest{
			Bucket: Ptr("bucket"),
			ReplicationRules: &ReplicationRules{
				[]string{"test_replication_1"},
			},
		},
		func(t *testing.T, o *DeleteBucketReplicationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteBucketReplication_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketReplicationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketReplication(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketReplicationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketReplicationRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketReplicationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=delete&replication", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationRules><ID>test_replication_1</ID></ReplicationRules>")
		},
		&DeleteBucketReplicationRequest{
			Bucket: Ptr("bucket"),
			ReplicationRules: &ReplicationRules{
				[]string{"test_replication_1"},
			},
		},
		func(t *testing.T, o *DeleteBucketReplicationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=delete&replication", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<ReplicationRules><ID>test_replication_1</ID></ReplicationRules>")
		},
		&DeleteBucketReplicationRequest{
			Bucket: Ptr("bucket"),
			ReplicationRules: &ReplicationRules{
				[]string{"test_replication_1"},
			},
		},
		func(t *testing.T, o *DeleteBucketReplicationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteBucketReplication_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketReplicationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketReplication(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


