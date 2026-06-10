package dataprocess

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func sortQuery(r *http.Request) string {
	u := r.URL
	var buf strings.Builder
	keys := make([]string, 0, len(u.Query()))
	for k := range u.Query() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := u.Query()[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			if len(v) > 0 {
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(v))
			}
		}
	}
	u.RawQuery = buf.String()
	return u.String()
}

func testSetupMockServer(t *testing.T, statusCode int, headers map[string]string, body []byte,
	chkfunc func(t *testing.T, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check request
		chkfunc(t, r)

		// header s
		for k, v := range headers {
			w.Header().Set(k, v)
		}

		// status code
		w.WriteHeader(statusCode)

		// body
		w.Write(body)
	}))
}

var testMockSimpleQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SimpleQueryRequest
	CheckOutputFn  func(t *testing.T, o *SimpleQueryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
  <NextToken>MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****</NextToken>
  <TotalHits>150</TotalHits>
  <Files>
    <File>
      <Filename>docs/report.pdf</Filename>
      <Size>5242880</Size>
      <URI>oss://examplebucket/docs/report.pdf</URI>
      <OSSURI>oss://examplebucket/docs/report.pdf</OSSURI>
      <MediaType>document</MediaType>
      <ContentType>application/pdf</ContentType>
      <FileModifiedTime>2025-12-01T10:30:00Z</FileModifiedTime>
      <PageCount>20</PageCount>
    </File>
  </Files>
  <Aggregations>
    <Aggregation>
      <Field>MediaType</Field>
      <Operation>group</Operation>
      <Groups>
        <AggregationGroup>
          <Value>document</Value>
          <Count>80</Count>
        </AggregationGroup>
        <AggregationGroup>
          <Value>image</Value>
          <Count>70</Count>
        </AggregationGroup>
      </Groups>
    </Aggregation>
  </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=simpleQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SimpleQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SimpleQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.NextToken, "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****")
			assert.Equal(t, *o.TotalHits, int64(150))
			assert.Equal(t, len(o.Files), 1)
			assert.Equal(t, *o.Files[0].Filename, "docs/report.pdf")
			assert.Equal(t, *o.Files[0].Size, int64(5242880))
			assert.Equal(t, *o.Files[0].URI, "oss://examplebucket/docs/report.pdf")
			assert.Equal(t, *o.Files[0].OSSURI, "oss://examplebucket/docs/report.pdf")
			assert.Equal(t, *o.Files[0].MediaType, "document")
			assert.Equal(t, *o.Files[0].ContentType, "application/pdf")
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2025-12-01T10:30:00Z")
			assert.Equal(t, *o.Files[0].PageCount, int64(20))
			assert.Equal(t, len(o.Aggregations), 1)
			assert.Equal(t, *o.Aggregations[0].Field, "MediaType")
			assert.Equal(t, *o.Aggregations[0].Operation, "group")
			assert.Equal(t, len(o.Aggregations[0].AggregationGroups), 2)
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Value, "document")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Count, int64(80))
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Value, "image")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Count, int64(70))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
  <NextToken>MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****</NextToken>
  <TotalHits>258</TotalHits>
  <Files>
    <File>
      <Filename>photos/sunset.jpg</Filename>
      <Size>2048000</Size>
      <URI>oss://examplebucket/photos/sunset.jpg</URI>
      <OSSURI>oss://examplebucket/photos/sunset.jpg</OSSURI>
      <MediaType>image</MediaType>
      <ContentType>image/jpeg</ContentType>
      <FileModifiedTime>2025-12-01T10:30:00Z</FileModifiedTime>
      <ImageWidth>3840</ImageWidth>
      <ImageHeight>2160</ImageHeight>
      <Orientation>1</Orientation>
    </File>
    <File>
      <Filename>photos/mountain.png</Filename>
      <Size>5120000</Size>
      <URI>oss://examplebucket/photos/mountain.png</URI>
      <OSSURI>oss://examplebucket/photos/mountain.png</OSSURI>
      <MediaType>image</MediaType>
      <ContentType>image/png</ContentType>
      <FileModifiedTime>2025-11-20T14:00:00Z</FileModifiedTime>
      <ImageWidth>1920</ImageWidth>
      <ImageHeight>1080</ImageHeight>
      <Orientation>1</Orientation>
    </File>
  </Files>
  <Aggregations>
    <Aggregation>
      <Field>MediaType</Field>
      <Operation>group</Operation>
      <Groups>
        <AggregationGroup>
          <Value>image</Value>
          <Count>200</Count>
        </AggregationGroup>
        <AggregationGroup>
          <Value>video</Value>
          <Count>58</Count>
        </AggregationGroup>
      </Groups>
    </Aggregation>
  </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=simpleQuery&aggregations=Size&datasetName=your_dataset&maxResults=99&metaQuery&nextToken=MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw%2A%2A%2A%2A&order=acs&query=%7B%22Field%22%3A+%22Size%22%2C%22Value%22%3A+%221%22%2C%22Operation%22%3A+%22gt%22%7D&sort=Size&withFields=%5B%22Filename%22%2C%22Size%22%5D&withoutTotalHits=true", strUrl)
		},
		&SimpleQueryRequest{
			Bucket:           oss.Ptr("bucket"),
			DatasetName:      oss.Ptr("your_dataset"),
			NextToken:        oss.Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
			MaxResults:       oss.Ptr(int32(99)),
			Query:            oss.Ptr("{\"Field\": \"Size\",\"Value\": \"1\",\"Operation\": \"gt\"}"),
			Sort:             oss.Ptr("Size"),
			Order:            oss.Ptr("acs"),
			Aggregations:     oss.Ptr("Size"),
			WithFields:       oss.Ptr(`["Filename","Size"]`),
			WithoutTotalHits: oss.Ptr(true),
		},
		func(t *testing.T, o *SimpleQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.NextToken, "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAx****")
			assert.Equal(t, *o.TotalHits, int64(258))
			assert.Equal(t, len(o.Files), 2)
			assert.Equal(t, *o.Files[0].Filename, "photos/sunset.jpg")
			assert.Equal(t, *o.Files[0].Size, int64(2048000))
			assert.Equal(t, *o.Files[0].URI, "oss://examplebucket/photos/sunset.jpg")
			assert.Equal(t, *o.Files[0].OSSURI, "oss://examplebucket/photos/sunset.jpg")
			assert.Equal(t, *o.Files[0].MediaType, "image")
			assert.Equal(t, *o.Files[0].ContentType, "image/jpeg")
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2025-12-01T10:30:00Z")
			assert.Equal(t, *o.Files[0].ImageWidth, int64(3840))
			assert.Equal(t, *o.Files[0].ImageHeight, int64(2160))
			assert.Equal(t, *o.Files[0].Orientation, int64(1))

			assert.Equal(t, *o.Files[1].Filename, "photos/mountain.png")
			assert.Equal(t, *o.Files[1].Size, int64(5120000))
			assert.Equal(t, *o.Files[1].URI, "oss://examplebucket/photos/mountain.png")
			assert.Equal(t, *o.Files[1].OSSURI, "oss://examplebucket/photos/mountain.png")
			assert.Equal(t, *o.Files[1].MediaType, "image")
			assert.Equal(t, *o.Files[1].ContentType, "image/png")
			assert.Equal(t, *o.Files[1].FileModifiedTime, "2025-11-20T14:00:00Z")
			assert.Equal(t, *o.Files[1].ImageWidth, int64(1920))
			assert.Equal(t, *o.Files[1].ImageHeight, int64(1080))
			assert.Equal(t, *o.Files[1].Orientation, int64(1))
			assert.Equal(t, len(o.Aggregations), 1)
			assert.Equal(t, *o.Aggregations[0].Field, "MediaType")
			assert.Equal(t, *o.Aggregations[0].Operation, "group")
			assert.Equal(t, len(o.Aggregations[0].AggregationGroups), 2)
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Value, "image")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Count, int64(200))
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Value, "video")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Count, int64(58))
		},
	},
}

func TestMockSimpleQuery_Success(t *testing.T) {
	for _, c := range testMockSimpleQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SimpleQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSimpleQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SimpleQueryRequest
	CheckOutputFn  func(t *testing.T, o *SimpleQueryResult, err error)
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
			assert.Equal(t, "/bucket/?action=simpleQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SimpleQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SimpleQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=simpleQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SimpleQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SimpleQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockSimpleQuery_Error(t *testing.T) {
	for _, c := range testMockSimpleQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SimpleQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSemanticQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SemanticQueryRequest
	CheckOutputFn  func(t *testing.T, o *SemanticQueryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
<Files>
      <File>
          <Addresses/>
          <AudioCovers/>
          <AudioStreams>
              <AudioStream>
                  <Bitrate>128000</Bitrate>
                  <ChannelLayout>stereo</ChannelLayout>
                  <Channels>2</Channels>
                  <CodecLongName>AAC (Advanced Audio Coding)</CodecLongName>
                  <CodecName>aac</CodecName>
                  <CodecTag>0x6134706d</CodecTag>
                  <CodecTagString>mp4a</CodecTagString>
                  <Duration>16.021769</Duration>
                  <FrameCount>690</FrameCount>
                  <Index>1</Index>
                  <SampleFormat>fltp</SampleFormat>
                  <SampleRate>44100</SampleRate>
                  <TimeBase>1/44100</TimeBase>
              </AudioStream>
          </AudioStreams>
          <Bitrate>1656706</Bitrate>
          <ContentMd5>5oJccWuBoqVXS8zrzckPlg==</ContentMd5>
          <ContentType>video/mp4</ContentType>
          <CreateTime>2026-04-21T20:28:17.018858947+08:00</CreateTime>
          <CroppingSuggestions/>
          <DatasetName>test-dataset-sem-vid-1776774492</DatasetName>
          <Duration>16.034</Duration>
          <ETag>\"E6825C716B81A2A5574BCCEBCDC90F96\"</ETag>
          <Elements/>
          <Figures/>
          <FileHash>E6825C716B81A2A5574BCCEBCDC90F96</FileHash>
          <FileModifiedTime>2026-04-21T20:28:13+08:00</FileModifiedTime>
          <Filename>test-temp/sem-vid-1776774492774503000.mp4</Filename>
          <FormatLongName>QuickTime / MOV</FormatLongName>
          <FormatName>mov,mp4,m4a,3gp,3g2,mj2</FormatName>
          <Insights>
              <Video>
                  <Caption>蓝衣男走向餐桌</Caption>
                  <Description>这是一段室内高角度监控录像，场景为一个客厅。</Description>
              </Video>
          </Insights>
          <Labels/>
          <MediaType>video</MediaType>
          <OCRContents/>
          <OSSCRC64>2327801188977127298</OSSCRC64>
          <OSSObjectType>Normal</OSSObjectType>
          <OSSStorageClass>Standard</OSSStorageClass>
          <OSSTagging>
              <routing-dataset>test-dataset-sem-vid-1776774492</routing-dataset>
          </OSSTagging>
          <OSSTaggingCount>1</OSSTaggingCount>
          <ObjectACL>default</ObjectACL>
          <SequenceNumber>2</SequenceNumber>
          <SemanticSimilarity>0.5583347777557373</SemanticSimilarity>
          <Size>3320455</Size>
          <SmartClusters/>
          <StreamCount>2</StreamCount>
          <Subtitles/>
          <URI>oss://oss-metaquery-dataset-test/test-temp/sem-vid-1776774492774503000.mp4</URI>
          <UpdateTime>2026-04-21T20:28:27.359034257+08:00</UpdateTime>
          <VideoHeight>1080</VideoHeight>
          <VideoStreams>
              <VideoStream>
                  <AverageFrameRate>21645000/721493</AverageFrameRate>
                  <BitDepth>8</BitDepth>
                  <Bitrate>1521221</Bitrate>
                  <CodecLongName>H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10</CodecLongName>
                  <CodecName>h264</CodecName>
                  <CodecTag>0x31637661</CodecTag>
                  <CodecTagString>avc1</CodecTagString>
                  <ColorPrimaries>bt709</ColorPrimaries>
                  <ColorRange>tv</ColorRange>
                  <ColorSpace>bt709</ColorSpace>
                  <ColorTransfer>bt709</ColorTransfer>
                  <DisplayAspectRatio>16:9</DisplayAspectRatio>
                  <Duration>16.033178</Duration>
                  <FrameCount>481</FrameCount>
                  <FrameRate>90000/2999</FrameRate>
                  <Height>1080</Height>
                  <Level>31</Level>
                  <PixelFormat>yuv420p</PixelFormat>
                  <Profile>High</Profile>
                  <SampleAspectRatio>1:1</SampleAspectRatio>
                  <TimeBase>1/90000</TimeBase>
                  <Width>1920</Width>
              </VideoStream>
          </VideoStreams>
          <VideoWidth>1920</VideoWidth>
      </File>
  </Files>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=semanticQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SemanticQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SemanticQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, len(o.Files), 1)
			assert.Equal(t, *o.Files[0].Bitrate, int64(1656706))
			assert.Equal(t, *o.Files[0].ContentMd5, "5oJccWuBoqVXS8zrzckPlg==")
			assert.Equal(t, *o.Files[0].ContentType, "video/mp4")
			assert.Equal(t, *o.Files[0].CreateTime, "2026-04-21T20:28:17.018858947+08:00")
			assert.Equal(t, *o.Files[0].DatasetName, "test-dataset-sem-vid-1776774492")
			assert.Equal(t, *o.Files[0].Duration, float64(16.034))
			assert.Equal(t, *o.Files[0].ETag, "\\\"E6825C716B81A2A5574BCCEBCDC90F96\\\"")
			assert.Equal(t, *o.Files[0].FileHash, "E6825C716B81A2A5574BCCEBCDC90F96")
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2026-04-21T20:28:13+08:00")
			assert.Equal(t, *o.Files[0].Filename, "test-temp/sem-vid-1776774492774503000.mp4")
			assert.Equal(t, *o.Files[0].FormatLongName, "QuickTime / MOV")
			assert.Equal(t, *o.Files[0].FormatName, "mov,mp4,m4a,3gp,3g2,mj2")
			assert.Equal(t, *o.Files[0].MediaType, "video")
			assert.Equal(t, *o.Files[0].Size, int64(3320455))
			assert.Equal(t, *o.Files[0].VideoWidth, int64(1920))
			assert.Equal(t, *o.Files[0].VideoHeight, int64(1080))
			assert.Equal(t, *o.Files[0].StreamCount, int64(2))
			assert.Equal(t, *o.Files[0].OSSObjectType, "Normal")
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, *o.Files[0].OSSTaggingCount, int64(1))
			assert.Equal(t, *o.Files[0].ObjectACL, "default")
			assert.Equal(t, *o.Files[0].SequenceNumber, int64(2))
			assert.Equal(t, *o.Files[0].SemanticSimilarity, float64(0.5583347777557373))
			assert.Equal(t, *o.Files[0].Size, int64(3320455))
			assert.Equal(t, *o.Files[0].URI, "oss://oss-metaquery-dataset-test/test-temp/sem-vid-1776774492774503000.mp4")
			assert.Equal(t, *o.Files[0].UpdateTime, "2026-04-21T20:28:27.359034257+08:00")

			assert.Equal(t, len(o.Files[0].AudioStreams), 1)
			assert.Equal(t, *o.Files[0].AudioStreams[0].Bitrate, int64(128000))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Channels, int64(2))
			assert.Equal(t, *o.Files[0].AudioStreams[0].ChannelLayout, "stereo")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecLongName, "AAC (Advanced Audio Coding)")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecName, "aac")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecTag, "0x6134706d")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecTagString, "mp4a")
			assert.Equal(t, *o.Files[0].AudioStreams[0].Duration, float64(16.021769))
			assert.Equal(t, *o.Files[0].AudioStreams[0].FrameCount, int64(690))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Index, int64(1))
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleFormat, "fltp")
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleRate, int64(44100))
			assert.Equal(t, *o.Files[0].AudioStreams[0].TimeBase, "1/44100")
			assert.Equal(t, *o.Files[0].Insights.Video.Description, "这是一段室内高角度监控录像，场景为一个客厅。")
			assert.Equal(t, *o.Files[0].Insights.Video.Caption, "蓝衣男走向餐桌")
			assert.Equal(t, len(o.Files[0].VideoStreams), 1)
			assert.Equal(t, *o.Files[0].VideoStreams[0].AverageFrameRate, "21645000/721493")
			assert.Equal(t, *o.Files[0].VideoStreams[0].BitDepth, int64(8))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Bitrate, int64(1521221))
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecLongName, "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecName, "h264")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecTag, "0x31637661")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecTagString, "avc1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorPrimaries, "bt709")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorRange, "tv")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorTransfer, "bt709")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorSpace, "bt709")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Duration, float64(16.033178))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameCount, int64(481))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameRate, "90000/2999")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Height, int64(1080))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Width, int64(1920))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Level, int64(31))
			assert.Equal(t, *o.Files[0].VideoStreams[0].PixelFormat, "yuv420p")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Profile, "High")
			assert.Equal(t, *o.Files[0].VideoStreams[0].PixelFormat, "yuv420p")
			assert.Equal(t, *o.Files[0].VideoStreams[0].SampleAspectRatio, "1:1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].TimeBase, "1/90000")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
    <Files>
        <File>
            <Addresses/>
            <AudioCovers/>
            <AudioStreams>
                <AudioStream>
                    <Bitrate>14983</Bitrate>
                    <ChannelLayout>mono</ChannelLayout>
                    <Channels>1</Channels>
                    <CodecLongName>AAC (Advanced Audio Coding)</CodecLongName>
                    <CodecName>aac</CodecName>
                    <CodecTag>0x6134706d</CodecTag>
                    <CodecTagString>mp4a</CodecTagString>
                    <Duration>7.936</Duration>
                    <FrameCount>62</FrameCount>
                    <Index>1</Index>
                    <SampleFormat>fltp</SampleFormat>
                    <SampleRate>8000</SampleRate>
                    <TimeBase>1/8000</TimeBase>
                </AudioStream>
            </AudioStreams>
            <Bitrate>196284</Bitrate>
            <ContentMd5>5/ZLrWYXpuQfDfxEf4+lyA==</ContentMd5>
            <ContentType>video/mp4</ContentType>
            <CreateTime>2026-04-21T10:51:38.264045621+08:00</CreateTime>
            <CroppingSuggestions/>
            <DatasetName>dataset-aianalysis-walk</DatasetName>
            <Duration>8</Duration>
            <ETag>\"E7F64BAD6617A6E41F0DFC447F8FA5C8\"</ETag>
            <Elements/>
            <Figures/>
            <FileHash>E7F64BAD6617A6E41F0DFC447F8FA5C8</FileHash>
            <FileModifiedTime>2026-04-21T10:51:25+08:00</FileModifiedTime>
            <Filename>mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4</Filename>
            <FormatLongName>QuickTime / MOV</FormatLongName>
            <FormatName>mov,mp4,m4a,3gp,3g2,mj2</FormatName>
            <Labels>
                <Label>
                    <LabelConfidence>1</LabelConfidence>
                    <LabelName>有人走过</LabelName>
                    <ParentLabelName>自定义标签</ParentLabelName>
                    <Clips>
                        <Clip>
                            <TimeRange>200</TimeRange>
                            <TimeRange>5533</TimeRange>
                        </Clip>
                    </Clips>
                </Label>
            </Labels>
            <MediaType>video</MediaType>
            <OCRContents/>
            <OSSCRC64>16628192875747293357</OSSCRC64>
            <OSSObjectType>Normal</OSSObjectType>
            <OSSStorageClass>Standard</OSSStorageClass>
            <OSSTagging>
                <alarmId>AE09411YAG0008117767395421908241</alarmId>
                <test-routing-dataset>dataset-aianalysis-walk</test-routing-dataset>
            </OSSTagging>
            <OSSTaggingCount>2</OSSTaggingCount>
            <OSSUserMeta>
                <X-Oss-Meta-Author>oss</X-Oss-Meta-Author>
            </OSSUserMeta>
            <ObjectACL>default</ObjectACL>
            <ProduceTime>2026-04-21T10:46:10+08:00</ProduceTime>
            <SceneElements>
                <SceneElement>
                    <FrameTimes>6000</FrameTimes>
                    <TimeRange>4133</TimeRange>
                    <TimeRange>8533</TimeRange>
                    <VideoStreamIndex>0</VideoStreamIndex>
                    <Labels/>
                </SceneElement>
            </SceneElements>
            <SemanticSimilarity>0.2536</SemanticSimilarity>
            <SequenceNumber>5</SequenceNumber>
            <Size>196284</Size>
            <SmartClusters/>
            <StreamCount>2</StreamCount>
            <Subtitles/>
            <URI>oss://paas-smart-cloud-test/mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4</URI>
            <UpdateTime>2026-04-21T10:52:39.412605575+08:00</UpdateTime>
            <VideoHeight>360</VideoHeight>
            <VideoStreams>
                <VideoStream>
                    <AverageFrameRate>15/1</AverageFrameRate>
                    <BitDepth>8</BitDepth>
                    <Bitrate>178202</Bitrate>
                    <CodecLongName>H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10</CodecLongName>
                    <CodecName>h264</CodecName>
                    <CodecTag>0x31637661</CodecTag>
                    <CodecTagString>avc1</CodecTagString>
                    <Duration>8</Duration>
                    <FrameCount>120</FrameCount>
                    <FrameRate>500/33</FrameRate>
                    <Height>360</Height>
                    <Level>22</Level>
                    <PixelFormat>yuv420p</PixelFormat>
                    <Profile>Main</Profile>
                    <TimeBase>1/1000</TimeBase>
                    <Width>640</Width>
                </VideoStream>
            </VideoStreams>
            <VideoWidth>640</VideoWidth>
        </File>
    </Files>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=semanticQuery&datasetName=your_dataset&maxResults=10&mediaTypes=%5B%22video%22%2C%22image%22%5D&metaQuery&nextToken=MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw%2A%2A%2A%2A&query=%7B%22Field%22%3A+%22Size%22%2C%22Value%22%3A+%221%22%2C%22Operation%22%3A+%22gt%22%7D&sourceURI=oss%3A%2F%2Fbucket%2Fprefix&withFields=%5B%22Filename%22%2C%22Size%22%5D", strUrl)
		},
		&SemanticQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			NextToken:   oss.Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
			MaxResults:  oss.Ptr(int32(10)),
			Query:       oss.Ptr("{\"Field\": \"Size\",\"Value\": \"1\",\"Operation\": \"gt\"}"),
			WithFields:  oss.Ptr(`["Filename","Size"]`),
			MediaTypes:  oss.Ptr(`["video","image"]`),
			SourceUri:   oss.Ptr("oss://bucket/prefix"),
		},
		func(t *testing.T, o *SemanticQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Files), 1)
			assert.Equal(t, *o.Files[0].Bitrate, int64(196284))
			assert.Equal(t, *o.Files[0].ContentMd5, "5/ZLrWYXpuQfDfxEf4+lyA==")
			assert.Equal(t, *o.Files[0].ContentType, "video/mp4")
			assert.Equal(t, *o.Files[0].CreateTime, "2026-04-21T10:51:38.264045621+08:00")
			assert.Equal(t, *o.Files[0].DatasetName, "dataset-aianalysis-walk")
			assert.Equal(t, *o.Files[0].Duration, float64(8))
			assert.Equal(t, *o.Files[0].ETag, "\\\"E7F64BAD6617A6E41F0DFC447F8FA5C8\\\"")
			assert.Equal(t, *o.Files[0].FileHash, "E7F64BAD6617A6E41F0DFC447F8FA5C8")
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2026-04-21T10:51:25+08:00")
			assert.Equal(t, *o.Files[0].Filename, "mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4")
			assert.Equal(t, *o.Files[0].FormatLongName, "QuickTime / MOV")
			assert.Equal(t, *o.Files[0].FormatName, "mov,mp4,m4a,3gp,3g2,mj2")
			assert.Equal(t, *o.Files[0].MediaType, "video")
			assert.Equal(t, *o.Files[0].OSSCRC64, "16628192875747293357")
			assert.Equal(t, *o.Files[0].Size, int64(196284))
			assert.Equal(t, *o.Files[0].VideoWidth, int64(640))
			assert.Equal(t, *o.Files[0].VideoHeight, int64(360))
			assert.Equal(t, *o.Files[0].StreamCount, int64(2))
			assert.Equal(t, *o.Files[0].OSSObjectType, "Normal")
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, len(o.Files[0].OSSTagging), 2)
			assert.Equal(t, o.Files[0].OSSTagging["alarmId"], "AE09411YAG0008117767395421908241")
			assert.Equal(t, o.Files[0].OSSTagging["test-routing-dataset"], "dataset-aianalysis-walk")
			assert.Equal(t, *o.Files[0].OSSTaggingCount, int64(2))
			assert.Equal(t, *o.Files[0].ObjectACL, "default")
			assert.Equal(t, len(o.Files[0].OSSUserMeta), 1)
			assert.Equal(t, o.Files[0].OSSUserMeta["X-Oss-Meta-Author"], "oss")
			assert.Equal(t, *o.Files[0].ProduceTime, "2026-04-21T10:46:10+08:00")
			assert.Equal(t, *o.Files[0].SequenceNumber, int64(5))
			assert.Equal(t, *o.Files[0].SemanticSimilarity, float64(0.2536))
			assert.Equal(t, *o.Files[0].Size, int64(196284))
			assert.Equal(t, *o.Files[0].StreamCount, int64(2))
			assert.Equal(t, *o.Files[0].URI, "oss://paas-smart-cloud-test/mp4file/AE09411YAG00081_AE09411YAG00081-0_e723c79f850047458a3e0c0115c4b108_20260421104610825sf0-203372.mp4")
			assert.Equal(t, *o.Files[0].UpdateTime, "2026-04-21T10:52:39.412605575+08:00")

			assert.Equal(t, len(o.Files[0].Labels), 1)
			assert.Equal(t, *o.Files[0].Labels[0].LabelConfidence, float64(1))
			assert.Equal(t, *o.Files[0].Labels[0].LabelName, "有人走过")
			assert.Equal(t, *o.Files[0].Labels[0].ParentLabelName, "自定义标签")
			assert.Equal(t, len(o.Files[0].Labels[0].Clips), 1)
			assert.Equal(t, o.Files[0].Labels[0].Clips[0].TimeRange[0], int64(200))
			assert.Equal(t, o.Files[0].Labels[0].Clips[0].TimeRange[1], int64(5533))

			assert.Equal(t, len(o.Files[0].AudioStreams), 1)
			assert.Equal(t, *o.Files[0].AudioStreams[0].Bitrate, int64(14983))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Channels, int64(1))
			assert.Equal(t, *o.Files[0].AudioStreams[0].ChannelLayout, "mono")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecLongName, "AAC (Advanced Audio Coding)")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecName, "aac")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecTag, "0x6134706d")
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecTagString, "mp4a")
			assert.Equal(t, *o.Files[0].AudioStreams[0].Duration, float64(7.936))
			assert.Equal(t, *o.Files[0].AudioStreams[0].FrameCount, int64(62))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Index, int64(1))
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleFormat, "fltp")
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleRate, int64(8000))
			assert.Equal(t, *o.Files[0].AudioStreams[0].TimeBase, "1/8000")

			assert.Equal(t, len(o.Files[0].VideoStreams), 1)
			assert.Equal(t, *o.Files[0].VideoStreams[0].AverageFrameRate, "15/1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].BitDepth, int64(8))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Bitrate, int64(178202))
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecLongName, "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecName, "h264")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecTag, "0x31637661")
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecTagString, "avc1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Duration, float64(8))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameCount, int64(120))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameRate, "500/33")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Height, int64(360))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Width, int64(640))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Level, int64(22))
			assert.Equal(t, *o.Files[0].VideoStreams[0].TimeBase, "1/1000")
		},
	},
}

func TestMockSemanticQuery_Success(t *testing.T) {
	for _, c := range testMockSemanticQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SemanticQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSemanticQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SemanticQueryRequest
	CheckOutputFn  func(t *testing.T, o *SemanticQueryResult, err error)
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
			assert.Equal(t, "/bucket/?action=semanticQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SemanticQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SemanticQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=semanticQuery&datasetName=your_dataset&metaQuery", strUrl)
		},
		&SemanticQueryRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *SemanticQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockSemanticQuery_Error(t *testing.T) {
	for _, c := range testMockSemanticQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SemanticQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockOpenMetaQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *OpenMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *OpenMetaQueryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=openMetaQuery&metaQuery&mode=basic&role=my-role", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><Filters><Filter>Size &gt; 0</Filter><Filter>ContentType = &#39;image/jpeg&#39;</Filter></Filters></MetaQuery>")
		},
		&OpenMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			Role:   oss.Ptr("my-role"),
			MetaQuery: &OpenMetaQuery{
				Filters: &Filters{
					Filter: []string{
						"Size > 0",
						"ContentType = 'image/jpeg'",
					},
				},
			},
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
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
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=openMetaQuery&metaQuery&mode=semantic&role=AliyunMetaQueryDefaultRole", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><WorkflowParameters><WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter><WorkflowParameter><Name>VideoInsightEnable</Name><Value>True</Value></WorkflowParameter></WorkflowParameters><NotificationAttributes><NotificationAttribute><Notifications><Notification><MNS>imm-index-notification</MNS></Notification></Notifications><WithFields><WithField>Insights</WithField><WithField>Labels</WithField></WithFields></NotificationAttribute></NotificationAttributes><DatasetConfig><Insights><Language>en</Language></Insights></DatasetConfig><IndexOptions><IgnoreObjectDelete>true</IgnoreObjectDelete></IndexOptions><RouteRule><Type>OSSTag</Type><AutoCreateDataset>true</AutoCreateDataset><OSSTagKey>routing-dataset</OSSTagKey></RouteRule></MetaQuery>")
		},
		&OpenMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Role:   oss.Ptr("AliyunMetaQueryDefaultRole"),
			Mode:   oss.Ptr("semantic"),
			MetaQuery: &OpenMetaQuery{
				WorkflowParameters: &WorkflowParameters{
					WorkflowParameter: []WorkflowParameter{
						{
							Name:  oss.Ptr("ImageInsightEnable"),
							Value: oss.Ptr("True"),
						},
						{
							Name:  oss.Ptr("VideoInsightEnable"),
							Value: oss.Ptr("True"),
						},
					},
				},
				NotificationAttributes: &NotificationAttributes{
					Notifications: &Notifications{
						Notification: []Notification{
							{
								MNS: oss.Ptr("imm-index-notification"),
							},
						},
					},
					WithFields: &WithFields{
						[]string{
							"Insights",
							"Labels",
						},
					},
				},
				IndexOptions: &IndexOptions{
					IgnoreObjectDelete: oss.Ptr(true),
				},
				RouteRule: &RouteRule{
					Type:              oss.Ptr("OSSTag"),
					AutoCreateDataset: oss.Ptr(true),
					OSSTagKey:         oss.Ptr("routing-dataset"),
				},
				DatasetConfig: &DatasetConfig{
					Insights: &InsightsConfig{
						Language: oss.Ptr("en"),
					},
				},
			},
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockOpenMetaQuery_Success(t *testing.T) {
	for _, c := range testMockOpenMetaQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.OpenMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockOpenMetaQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *OpenMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *OpenMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?action=openMetaQuery&metaQuery&mode=basic", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=openMetaQuery&metaQuery&mode=basic", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockOpenMetaQuery_Error(t *testing.T) {
	for _, c := range testMockOpenMetaQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.OpenMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetMetaQueryStatusSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetMetaQueryStatusRequest
	CheckOutputFn  func(t *testing.T, o *GetMetaQueryStatusResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQueryStatus>
  <State>Running</State>
  <Phase>IncrementalScanning</Phase>
  <CreateTime>2026-05-20T08:00:00.000+08:00</CreateTime>
  <UpdateTime>2026-05-20T08:30:00.000+08:00</UpdateTime>
  <MetaQueryMode>semantic</MetaQueryMode>
  <WorkflowParameters>
    <WorkflowParameter>
      <Name>ImageInsightEnable</Name>
      <Value>True</Value>
    </WorkflowParameter>
  </WorkflowParameters>
  <Filters>
    <Filter>Size > 1024,FileModifiedTime > 2025-06-03T09:20:47.999Z</Filter>
    <Filter>Filename prefix (YWEvYmIv)</Filter>
  </Filters>
  <IndexOptions>
    <IgnoreObjectDelete>True</IgnoreObjectDelete>
  </IndexOptions>
  <RouteRule>
    <Type>OSSTag</Type>
    <AutoCreateDataset>True</AutoCreateDataset>
    <OSSTagKey>routing-dataset</OSSTagKey>
  </RouteRule>
  <NotificationAttributes>
    <Notifications>
      <Notification>
        <MNS>imm-index-notification</MNS>
      </Notification>
    </Notifications>
    <WithFields>
      <WithField>Insights</WithField>
    </WithFields>
  </NotificationAttributes>
  <DatasetConfig>
    <Insights>
      <Language>en</Language>
    </Insights>
  </DatasetConfig>
</MetaQueryStatus>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getMetaQueryStatus&metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.State, "Running")
			assert.Equal(t, *o.Phase, "IncrementalScanning")
			assert.Equal(t, *o.CreateTime, "2026-05-20T08:00:00.000+08:00")
			assert.Equal(t, *o.UpdateTime, "2026-05-20T08:30:00.000+08:00")
			assert.Equal(t, *o.MetaQueryMode, "semantic")
			assert.Equal(t, *o.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
			assert.Equal(t, *o.WorkflowParameters.WorkflowParameter[0].Value, "True")
			assert.Equal(t, o.Filters.Filter[0], "Size > 1024,FileModifiedTime > 2025-06-03T09:20:47.999Z")
			assert.Equal(t, o.Filters.Filter[1], "Filename prefix (YWEvYmIv)")
			assert.Equal(t, *o.IndexOptions.IgnoreObjectDelete, true)
			assert.Equal(t, *o.RouteRule.Type, "OSSTag")
			assert.Equal(t, *o.RouteRule.AutoCreateDataset, true)
			assert.Equal(t, *o.RouteRule.OSSTagKey, "routing-dataset")
			assert.Equal(t, *o.NotificationAttributes.Notifications.Notification[0].MNS, "imm-index-notification")
			assert.Equal(t, o.NotificationAttributes.WithFields.WithField[0], "Insights")
			assert.Equal(t, *o.DatasetConfig.Insights.Language, "en")

		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<MetaQueryStatus>
    <State>Running</State>
    <Phase>IncrementalScanning</Phase>
    <CreateTime>2026-06-05T10:49:57.109295646+08:00</CreateTime>
    <UpdateTime>2026-06-05T10:49:59.411578836+08:00</UpdateTime>
    <MetaQueryMode>basic</MetaQueryMode>
    <RouteRule>
        <Type>default</Type>
        <AutoCreateDataset>True</AutoCreateDataset>
    </RouteRule>
    <WorkflowParameters>
    </WorkflowParameters>
    <Filters>
    </Filters>
    <DatasetConfig>
        <ReverseImage>
            <Video>
                <Enable>false</Enable>
            </Video>
            <Image>
                <Enable>false</Enable>
            </Image>
        </ReverseImage>
        <Insights>
            <Language>zh-Hans</Language>
            <Image>
                <Caption>
                    <Enable>false</Enable>
                    <Prompt>
                    </Prompt>
                </Caption>
            </Image>
            <Video>
                <Caption>
                    <Enable>false</Enable>
                    <Prompt>
                    </Prompt>
                    <PersonReference>
                        <Enable>false</Enable>
                    </PersonReference>
                </Caption>
                <Label>
                    <System>
                        <Enable>false</Enable>
                    </System>
                    <UserDefined>
                        <Enable>false</Enable>
                        <Labels>
                        </Labels>
                    </UserDefined>
                    <Highlight>
                        <Enable>false</Enable>
                        <Labels>
                        </Labels>
                    </Highlight>
                </Label>
                <MultiStream>
                    <Enable>false</Enable>
                </MultiStream>
            </Video>
        </Insights>
        <SmartCluster>
            <Figure>
                <AutoGenerate>false</AutoGenerate>
                <AutoClustering>false</AutoClustering>
                <MinEntityCount>3</MinEntityCount>
                <EnabledFeatures>face</EnabledFeatures>
            </Figure>
        </SmartCluster>
    </DatasetConfig>
</MetaQueryStatus>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getMetaQueryStatus&metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.State, "Running")
			assert.Equal(t, *o.Phase, "IncrementalScanning")
			assert.Equal(t, *o.CreateTime, "2026-06-05T10:49:57.109295646+08:00")
			assert.Equal(t, *o.UpdateTime, "2026-06-05T10:49:59.411578836+08:00")
			assert.Equal(t, *o.MetaQueryMode, "basic")
			assert.Equal(t, *o.RouteRule.Type, "default")
			assert.Equal(t, *o.RouteRule.AutoCreateDataset, true)
			assert.Equal(t, *o.DatasetConfig.ReverseImage.Video.Enable, false)
			assert.Equal(t, *o.DatasetConfig.ReverseImage.Image.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Language, "zh-Hans")
			assert.Equal(t, *o.DatasetConfig.Insights.Image.Caption.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Caption.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Caption.PersonReference.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Label.System.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Label.UserDefined.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Label.Highlight.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.MultiStream.Enable, false)
			assert.Equal(t, *o.DatasetConfig.SmartCluster.Figure.AutoGenerate, false)
			assert.Equal(t, *o.DatasetConfig.SmartCluster.Figure.AutoClustering, false)
			assert.Equal(t, *o.DatasetConfig.SmartCluster.Figure.MinEntityCount, int64(3))
			assert.Equal(t, o.DatasetConfig.SmartCluster.Figure.EnabledFeatures[0], "face")
		},
	},
}

func TestMockGetMetaQueryStatus_Success(t *testing.T) {
	for _, c := range testMockGetMetaQueryStatusSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetMetaQueryStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetMetaQueryStatusErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetMetaQueryStatusRequest
	CheckOutputFn  func(t *testing.T, o *GetMetaQueryStatusResult, err error)
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
			assert.Equal(t, "/bucket/?action=getMetaQueryStatus&metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=getMetaQueryStatus&metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockGetMetaQueryStatus_Error(t *testing.T) {
	for _, c := range testMockGetMetaQueryStatusErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetMetaQueryStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDoMetaQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DoMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *DoMetaQueryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
    <NextToken>next-page-token-abc</NextToken>
    <TotalHits>123</TotalHits>
    <Files>
        <File>
            <Filename>photos/sunset.jpg</Filename>
            <Size>2097152</Size>
            <FileModifiedTime>2026-05-19T15:30:00.000+08:00</FileModifiedTime>
            <ContentType>image/jpeg</ContentType>
            <ObjectACL>default</ObjectACL>
            <OSSStorageClass>Standard</OSSStorageClass>
        </File>
        <File>
            <Filename>photos/mountain.png</Filename>
            <Size>5242880</Size>
        </File>
    </Files>
    <Aggregations>
        <Aggregation>
            <Field>Size</Field>
            <Operation>sum</Operation>
            <Value>12345678</Value>
        </Aggregation>
    </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.NextToken, "next-page-token-abc")
			assert.Equal(t, *o.TotalHits, int64(123))
			assert.Equal(t, *o.Files[0].Filename, "photos/sunset.jpg")
			assert.Equal(t, *o.Files[0].Size, int64(2097152))
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2026-05-19T15:30:00.000+08:00")
			assert.Equal(t, *o.Files[0].ContentType, "image/jpeg")
			assert.Equal(t, *o.Files[0].ObjectACL, "default")
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, *o.Files[1].Filename, "photos/mountain.png")
			assert.Equal(t, *o.Files[1].Size, int64(5242880))
			assert.Equal(t, *o.Aggregations[0].Field, "Size")
			assert.Equal(t, *o.Aggregations[0].Operation, "sum")
			assert.Equal(t, *o.Aggregations[0].Value, "12345678")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
    <TotalHits>50</TotalHits>
    <Aggregations>
        <Aggregation>
            <Field>StorageClass</Field>
            <Operation>group_by</Operation>
            <Groups>
                <AggregationGroup>
                    <Value>Standard</Value>
                    <Count>30</Count>
                </AggregationGroup>
                <AggregationGroup>
                    <Value>IA</Value>
                    <Count>20</Count>
                </AggregationGroup>
            </Groups>
        </Aggregation>
    </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.TotalHits, int64(50))
			assert.Equal(t, len(o.Aggregations), 1)
			assert.Equal(t, *o.Aggregations[0].Field, "StorageClass")
			assert.Equal(t, *o.Aggregations[0].Operation, "group_by")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Value, "Standard")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Count, int64(30))
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Value, "IA")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Count, int64(20))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
    <TotalHits>2</TotalHits>
    <Files>
        <File>
            <Filename>photos/cat-in-living-room.jpg</Filename>
            <Size>3145728</Size>
            <OSSStorageClass>Standard</OSSStorageClass>
            <Labels>
                <Label>
                    <LabelName>cat</LabelName>
                    <LabelConfidence>0.98</LabelConfidence>
                </Label>
            </Labels>
        </File>
        <File>
            <Filename>photos/kitten-sofa.jpg</Filename>
            <Size>2621440</Size>
        </File>
    </Files>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.TotalHits, int64(2))
			assert.Equal(t, *o.Files[0].Filename, "photos/cat-in-living-room.jpg")
			assert.Equal(t, *o.Files[0].Size, int64(3145728))
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, *o.Files[0].Labels[0].LabelName, "cat")
			assert.Equal(t, *o.Files[0].Labels[0].LabelConfidence, float64(0.98))
			assert.Equal(t, *o.Files[1].Filename, "photos/kitten-sofa.jpg")
			assert.Equal(t, *o.Files[1].Size, int64(2621440))
		},
	},
}

func TestMockDoMetaQuery_Success(t *testing.T) {
	for _, c := range testMockDoMetaQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DoMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDoMetaQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DoMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *DoMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockDoMetaQuery_Error(t *testing.T) {
	for _, c := range testMockDoMetaQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DoMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCloseMetaQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CloseMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *CloseMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?action=closeMetaQuery&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *CloseMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockCloseMetaQuery_Success(t *testing.T) {
	for _, c := range testMockCloseMetaQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CloseMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCloseMetaQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CloseMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *CloseMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?action=closeMetaQuery&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *CloseMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=closeMetaQuery&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *CloseMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockCloseMetaQuery_Error(t *testing.T) {
	for _, c := range testMockCloseMetaQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CloseMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateDatasetSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateDatasetRequest
	CheckOutputFn  func(t *testing.T, o *CreateDatasetResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CreateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters></WorkflowParameters>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
<DatasetMaxBindCount>10</DatasetMaxBindCount>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh</Language></Insights></DatasetConfig>
</Dataset>
</CreateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createDataset&datasetName=your_dataset&metaQuery", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Dataset.DatasetName, "test-dataset")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-22T11:39:28.148283473+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-22T11:39:28.148283473+08:00")
			assert.Equal(t, *o.Dataset.DatasetMaxBindCount, int64(10))
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CreateDatasetResponse>
    <Dataset>
        <DatasetName>test_dataset</DatasetName>
        <WorkflowParameters>
        </WorkflowParameters>
        <CreateTime>2026-06-08T17:36:59.774068044+08:00</CreateTime>
        <UpdateTime>2026-06-08T17:36:59.774068044+08:00</UpdateTime>
        <Description>this is a demo</Description>
        <DatasetMaxBindCount>10</DatasetMaxBindCount>
        <DatasetMaxFileCount>100000000</DatasetMaxFileCount>
        <DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
        <DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
        <DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
        <DatasetConfig>
            <ReverseImage>
                <Video>
                    <Enable>false</Enable>
                </Video>
                <Image>
                    <Enable>false</Enable>
                </Image>
            </ReverseImage>
            <Insights>
                <Language>zh</Language>
                <Image>
                    <Caption>
                        <Enable>false</Enable>
                        <Prompt></Prompt>
                    </Caption>
                </Image>
                <Video>
                    <Caption>
                        <Enable>false</Enable>
                        <Prompt></Prompt>
                        <PersonReference>
                            <Enable>false</Enable>
                        </PersonReference>
                    </Caption>
                    <Label>
                        <System>
                            <Enable>false</Enable>
                        </System>
                        <UserDefined>
                            <Enable>false</Enable>
                            <Labels>
                            </Labels>
                        </UserDefined>
                        <Highlight>
                            <Enable>false</Enable>
                            <Labels>
                            </Labels>
                        </Highlight>
                    </Label>
                    <MultiStream>
                        <Enable>false</Enable>
                    </MultiStream>
                </Video>
            </Insights>
            <SmartCluster>
                <Figure>
                    <AutoGenerate>false</AutoGenerate>
                    <AutoClustering>false</AutoClustering>
                    <MinEntityCount>3</MinEntityCount>
                    <EnabledFeatures>face</EnabledFeatures>
                </Figure>
            </SmartCluster>
        </DatasetConfig>
    </Dataset>
</CreateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createDataset&datasetConfig=%7B%0A%09%09%09%09++%22Insights%22%3A+%7B%0A%09%09%09%09%09%22EnableLabel%22%3A+true%2C%0A%09%09%09%09%09%22EnableOCR%22%3A+true%2C%0A%09%09%09%09%09%22EnableFace%22%3A+true%2C%0A%09%09%09%09%09%22EnableImage%22%3A+true%2C%0A%09%09%09%09%09%22EnableVideo%22%3A+true%2C%0A%09%09%09%09%09%22EnableAudio%22%3A+true%2C%0A%09%09%09%09%09%22Language%22%3A+%22zh%22%0A%09%09%09%09++%7D%0A%09%09%09%09%7D&datasetName=test_dataset&description=this+is+a+demo&metaQuery", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
			Description: oss.Ptr("this is a demo"),
			DatasetConfig: oss.Ptr(`{
				  "Insights": {
					"EnableLabel": true,
					"EnableOCR": true,
					"EnableFace": true,
					"EnableImage": true,
					"EnableVideo": true,
					"EnableAudio": true,
					"Language": "zh"
				  }
				}`),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Dataset.DatasetName, "test_dataset")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-06-08T17:36:59.774068044+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-06-08T17:36:59.774068044+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxBindCount, int64(10))
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Video.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Image.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Video.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Image.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Prompt, "")
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.PersonReference.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.System.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.UserDefined.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.Highlight.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.MultiStream.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh")

			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.AutoGenerate, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.AutoClustering, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.MinEntityCount, int64(3))
			assert.Equal(t, o.Dataset.DatasetConfig.SmartCluster.Figure.EnabledFeatures[0], "face")
		},
	},
}

func TestMockCreateDataset_Success(t *testing.T) {
	for _, c := range testMockCreateDatasetSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateDatasetErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateDatasetRequest
	CheckOutputFn  func(t *testing.T, o *CreateDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=createDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=createDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockCreateDataset_Error(t *testing.T) {
	for _, c := range testMockCreateDatasetErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetDatasetSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetDatasetRequest
	CheckOutputFn  func(t *testing.T, o *GetDatasetResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<GetDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters>
      <WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter>
</WorkflowParameters>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxBindCount>10</DatasetMaxBindCount>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh-Hans</Language></Insights></DatasetConfig>
<FileCount>3456</FileCount>
<TotalFileSize>10737418240</TotalFileSize>
</Dataset>
</GetDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getDataset&datasetName=your_dataset&metaQuery&withStatistics=true", strUrl)
		},
		&GetDatasetRequest{
			Bucket:         oss.Ptr("bucket"),
			DatasetName:    oss.Ptr("your_dataset"),
			WithStatistics: oss.Ptr(true),
		},
		func(t *testing.T, o *GetDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Dataset.DatasetName, "test-dataset")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Value, "True")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxBindCount, int64(10))
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh-Hans")
			assert.Equal(t, *o.Dataset.FileCount, int64(3456))
			assert.Equal(t, *o.Dataset.TotalFileSize, int64(10737418240))

		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<GetDatasetResponse>
<Dataset>
<DatasetName>photos-2026</DatasetName>
<WorkflowParameters></WorkflowParameters>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxBindCount>10</DatasetMaxBindCount>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh-Hans</Language></Insights></DatasetConfig>
</Dataset>
</GetDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&GetDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *GetDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Dataset.DatasetName, "photos-2026")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxBindCount, int64(10))
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh-Hans")
		},
	},
}

func TestMockGetDataset_Success(t *testing.T) {
	for _, c := range testMockGetDatasetSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetDatasetErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetDatasetRequest
	CheckOutputFn  func(t *testing.T, o *GetDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=getDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&GetDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *GetDatasetResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=getDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&GetDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *GetDatasetResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockGetDataset_Error(t *testing.T) {
	for _, c := range testMockGetDatasetErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateDatasetSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateDatasetRequest
	CheckOutputFn  func(t *testing.T, o *UpdateDatasetResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<UpdateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters>
      <WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter>
</WorkflowParameters>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxBindCount>10</DatasetMaxBindCount>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig>
<Insights>
	<Language>zh</Language>
	<Image>
		<Caption>
			<Enable>true</Enable>
		</Caption>
	</Image>
	<Video>
		<Caption>
			<Enable>true</Enable>
		</Caption>
		<Label>
			<System>
				<Enable>true</Enable>
			</System>
			<UserDefined>
				<Enable>true</Enable>
			</UserDefined>
			<Highlight>
				<Enable>true</Enable>
			</Highlight>
		</Label>
		<MultiStream>
			<Enable>true</Enable>
		</MultiStream>
	</Video>
</Insights>
</DatasetConfig>
</Dataset>
</UpdateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=updateDataset&datasetConfig=%7B%22Insights%22%3A%7B%22Language%22%3A%22zh%22%2C%22Image%22%3A%7B%22Caption%22%3A%7B%22Enable%22%3A%22true%22%7D%7D%2C%22Video%22%3A%7B%22Caption%22%3A%7B%22Enable%22%3A%22true%22%7D%2C%22Label%22%3A%7B%22System%22%3A%7B%22Enable%22%3A%22true%22%7D%2C%22UserDefined%22%3A%7B%22Enable%22%3A%22true%22%7D%2C%22Highlight%22%3A%7B%22Enable%22%3A%22true%22%7D%7D%2C%22MultiStream%22%3A%7B%22Enable%22%3A%22true%22%7D%7D%7D%7D&datasetName=your_dataset&description=this+is+a+demo&metaQuery&workflowParameters=%5B%7B%22Name%22%3A+%22ImageInsightEnable%22%2C%22Value%22%3A+%22True%22%2C%22Description%22%3A+%22The+source+bucket+for+data+processing%22%7D%5D", strUrl)

		},
		&UpdateDatasetRequest{
			Bucket:             oss.Ptr("bucket"),
			DatasetName:        oss.Ptr("your_dataset"),
			Description:        oss.Ptr("this is a demo"),
			WorkflowParameters: oss.Ptr(`[{"Name": "ImageInsightEnable","Value": "True","Description": "The source bucket for data processing"}]`),
			DatasetConfig:      oss.Ptr(`{"Insights":{"Language":"zh","Image":{"Caption":{"Enable":"true"}},"Video":{"Caption":{"Enable":"true"},"Label":{"System":{"Enable":"true"},"UserDefined":{"Enable":"true"},"Highlight":{"Enable":"true"}},"MultiStream":{"Enable":"true"}}}}`),
		},
		func(t *testing.T, o *UpdateDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Dataset.DatasetName, "test-dataset")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Value, "True")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxBindCount, int64(10))
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.System.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.UserDefined.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.Highlight.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.MultiStream.Enable)
		},
	},
}

func TestMockUpdateDataset_Success(t *testing.T) {
	for _, c := range testMockUpdateDatasetSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateDatasetErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateDatasetRequest
	CheckOutputFn  func(t *testing.T, o *UpdateDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=updateDataset&datasetName=test_dataset&description=this+is+a+demo&metaQuery", strUrl)
		},
		&UpdateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
			Description: oss.Ptr("this is a demo"),
		},
		func(t *testing.T, o *UpdateDatasetResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=updateDataset&datasetName=test_dataset&description=this+is+a+demo&metaQuery", strUrl)
		},
		&UpdateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
			Description: oss.Ptr("this is a demo"),
		},
		func(t *testing.T, o *UpdateDatasetResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockUpdateDataset_Error(t *testing.T) {
	for _, c := range testMockUpdateDatasetErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListDatasetsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListDatasetsRequest
	CheckOutputFn  func(t *testing.T, o *ListDatasetsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<ListDatasetsResponse>
<NextToken>1986505809429276:oss_1234567890_demo-bucket:test-dataset</NextToken>
<Datasets>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
</Dataset>
</Datasets>
</ListDatasetsResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", strUrl)
		},
		&ListDatasetsRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListDatasetsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, len(o.Datasets), 1)
			assert.Equal(t, *o.NextToken, "1986505809429276:oss_1234567890_demo-bucket:test-dataset")
			assert.Equal(t, *o.Datasets[0].DatasetName, "test-dataset")
			assert.Equal(t, *o.Datasets[0].CreateTime, "2026-04-22T11:39:28.148283473+08:00")
			assert.Equal(t, *o.Datasets[0].UpdateTime, "2026-04-22T11:39:28.148283473+08:00")

		},
	},
}

func TestMockListDatasets_Success(t *testing.T) {
	for _, c := range testMockListDatasetsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListDatasets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListDatasetsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListDatasetsRequest
	CheckOutputFn  func(t *testing.T, o *ListDatasetsResult, err error)
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
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", strUrl)
		},
		&ListDatasetsRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListDatasetsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", strUrl)
		},
		&ListDatasetsRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListDatasetsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockListDatasets_Error(t *testing.T) {
	for _, c := range testMockListDatasetsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListDatasets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteDatasetSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteDatasetRequest
	CheckOutputFn  func(t *testing.T, o *DeleteDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteDataset&datasetName=test-dataset&metaQuery", strUrl)
		},
		&DeleteDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
		},
		func(t *testing.T, o *DeleteDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteDataset_Success(t *testing.T) {
	for _, c := range testMockDeleteDatasetSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteDatasetErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteDatasetRequest
	CheckOutputFn  func(t *testing.T, o *DeleteDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteDataset&datasetName=test-dataset&metaQuery", strUrl)
		},
		&DeleteDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
		},
		func(t *testing.T, o *DeleteDatasetResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=deleteDataset&datasetName=test-dataset&metaQuery", strUrl)
		},
		&DeleteDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
		},
		func(t *testing.T, o *DeleteDatasetResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockDeleteDataset_Error(t *testing.T) {
	for _, c := range testMockDeleteDatasetErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteFileMetaSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteFileMetaRequest
	CheckOutputFn  func(t *testing.T, o *DeleteFileMetaResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteFileMeta&datasetName=test-dataset&metaQuery&uri=oss%3A%2F%2Fbucket%2Fobject", strUrl)
		},
		&DeleteFileMetaRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Uri:         oss.Ptr("oss://bucket/object"),
		},
		func(t *testing.T, o *DeleteFileMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteFileMeta_Success(t *testing.T) {
	for _, c := range testMockDeleteFileMetaSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteFileMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteFileMetaErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteFileMetaRequest
	CheckOutputFn  func(t *testing.T, o *DeleteFileMetaResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteFileMeta&datasetName=test-dataset&metaQuery&uri=oss%3A%2F%2Fbucket%2Fobject", strUrl)
		},
		&DeleteFileMetaRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Uri:         oss.Ptr("oss://bucket/object"),
		},
		func(t *testing.T, o *DeleteFileMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=deleteFileMeta&datasetName=test-dataset&metaQuery&uri=oss%3A%2F%2Fbucket%2Fobject", strUrl)
		},
		&DeleteFileMetaRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Uri:         oss.Ptr("oss://bucket/object"),
		},
		func(t *testing.T, o *DeleteFileMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockDeleteFileMeta_Error(t *testing.T) {
	for _, c := range testMockDeleteFileMetaErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteFileMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateSmartClusterSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *CreateSmartClusterResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CreateSmartClusterResponse>
  <ObjectId>cluster-abc123def456</ObjectId>
</CreateSmartClusterResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createSmartCluster&clusterType=knowledge&datasetName=test-dataset&metaQuery&name=your_name&rules=%5B%7B%22RuleType%22%3A+%22keywords%22%2C%22Keywords%22%3A+%5B%22car%22%5D%7D%5D", strUrl)
		},
		&CreateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Name:        oss.Ptr("your_name"),
			ClusterType: SmartClusterTypeKnowledge,
			Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["car"]}]`),
		},
		func(t *testing.T, o *CreateSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectId, "cluster-abc123def456")
		},
	},
}

func TestMockCreateSmartCluster_Success(t *testing.T) {
	for _, c := range testMockCreateSmartClusterSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateSmartClusterErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *CreateSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=createSmartCluster&clusterType=knowledge&datasetName=test-dataset&metaQuery&name=your_name&rules=%5B%7B%22RuleType%22%3A+%22keywords%22%2C%22Keywords%22%3A+%5B%22car%22%5D%7D%5D", strUrl)
		},
		&CreateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Name:        oss.Ptr("your_name"),
			ClusterType: SmartClusterTypeKnowledge,
			Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["car"]}]`),
		},
		func(t *testing.T, o *CreateSmartClusterResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=createSmartCluster&clusterType=knowledge&datasetName=test-dataset&metaQuery&name=your_name&rules=%5B%7B%22RuleType%22%3A+%22keywords%22%2C%22Keywords%22%3A+%5B%22car%22%5D%7D%5D", strUrl)
		},
		&CreateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Name:        oss.Ptr("your_name"),
			ClusterType: SmartClusterTypeKnowledge,
			Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["car"]}]`),
		},
		func(t *testing.T, o *CreateSmartClusterResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockCreateSmartCluster_Error(t *testing.T) {
	for _, c := range testMockCreateSmartClusterErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetSmartClusterSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *GetSmartClusterResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<GetSmartClusterResponse>
  <SmartCluster>
    <ObjectId>cluster-abc123def456</ObjectId>
    <ClusterType>figure</ClusterType>
    <Name>face-cluster-alice</Name>
    <Description>this is a demo</Description>
    <Rules>
      <Rule>
        <RuleType>face</RuleType>
        <BaseURIs>oss://examplebucket/refs/alice.jpg</BaseURIs>
        <Sensitivity>0.7</Sensitivity>
      </Rule>
    </Rules>
    <Reason></Reason>
    <Notification>
      <MNS><TopicName>imm-cluster-notification</TopicName></MNS>
    </Notification>
    <CreateTime>2026-05-20T11:00:00.000+08:00</CreateTime>
    <UpdateTime>2026-05-20T11:08:00.000+08:00</UpdateTime>
  </SmartCluster>
</GetSmartClusterResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&GetSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *GetSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.SmartCluster.ObjectId, "cluster-abc123def456")
			assert.Equal(t, *o.SmartCluster.ClusterType, "figure")
			assert.Equal(t, *o.SmartCluster.CreateTime, "2026-05-20T11:00:00.000+08:00")
			assert.Equal(t, *o.SmartCluster.UpdateTime, "2026-05-20T11:08:00.000+08:00")
			assert.Equal(t, *o.SmartCluster.Name, "face-cluster-alice")
			assert.Equal(t, *o.SmartCluster.Description, "this is a demo")
			assert.Equal(t, *o.SmartCluster.Rules[0].RuleType, "face")
			assert.Equal(t, o.SmartCluster.Rules[0].BaseURIs[0], "oss://examplebucket/refs/alice.jpg")
			assert.Equal(t, *o.SmartCluster.Rules[0].Sensitivity, 0.7)
			assert.Equal(t, *o.SmartCluster.Notification.MNS.TopicName, "imm-cluster-notification")
		},
	},
}

func TestMockGetSmartCluster_Success(t *testing.T) {
	for _, c := range testMockGetSmartClusterSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetSmartClusterErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *GetSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=getSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&GetSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *GetSmartClusterResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=getSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&GetSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *GetSmartClusterResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockGetSmartCluster_Error(t *testing.T) {
	for _, c := range testMockGetSmartClusterErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateSmartClusterSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *UpdateSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=updateSmartCluster&datasetName=your_dataset&description=this+is+a+demo&metaQuery&name=face-cluster-alice&objectId=cluster-abc123def456&rules=%5B%7B%22RuleType%22%3A%22face%22%2C%22Sensitivity%22%3A0.7%7D%5D", strUrl)
		},
		&UpdateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
			Description: oss.Ptr("this is a demo"),
			Name:        oss.Ptr("face-cluster-alice"),
			Rules:       oss.Ptr("[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]"),
		},
		func(t *testing.T, o *UpdateSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockUpdateSmartCluster_Success(t *testing.T) {
	for _, c := range testMockUpdateSmartClusterSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateSmartClusterErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *UpdateSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=updateSmartCluster&datasetName=your_dataset&description=this+is+a+demo&metaQuery&name=face-cluster-alice&objectId=cluster-abc123def456&rules=%5B%7B%22RuleType%22%3A%22face%22%2C%22Sensitivity%22%3A0.7%7D%5D", strUrl)
		},
		&UpdateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
			Description: oss.Ptr("this is a demo"),
			Name:        oss.Ptr("face-cluster-alice"),
			Rules:       oss.Ptr("[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]"),
		},
		func(t *testing.T, o *UpdateSmartClusterResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=updateSmartCluster&datasetName=your_dataset&description=this+is+a+demo&metaQuery&name=face-cluster-alice&objectId=cluster-abc123def456&rules=%5B%7B%22RuleType%22%3A%22face%22%2C%22Sensitivity%22%3A0.7%7D%5D", strUrl)
		},
		&UpdateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
			Description: oss.Ptr("this is a demo"),
			Name:        oss.Ptr("face-cluster-alice"),
			Rules:       oss.Ptr("[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]"),
		},
		func(t *testing.T, o *UpdateSmartClusterResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockUpdateSmartCluster_Error(t *testing.T) {
	for _, c := range testMockUpdateSmartClusterErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListSmartClustersSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListSmartClustersRequest
	CheckOutputFn  func(t *testing.T, o *ListSmartClustersResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<ListSmartClustersResponse>
    <SmartClusters>
        <SmartCluster>
            <CreateTime>2026-06-10T09:54:27.484585901+08:00</CreateTime>
            <ObjectId>SmartCluster-cb9f8c95-281f-490b-b677-eca2f6ff0c19</ObjectId>
            <UpdateTime>2026-06-10T09:54:27.484585901+08:00</UpdateTime>
            <ClusterType>knowledge</ClusterType>
            <Name>demo-2</Name>
            <Rules>
                <Rule>
                    <Keywords>cat</Keywords>
                    <Sensitivity>0.5</Sensitivity>
                    <RuleType>keywords</RuleType>
                </Rule>
            </Rules>
            <Reason></Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T16:18:20.998039823+08:00</CreateTime>
            <ObjectId>SmartCluster-c30d039b-0b55-4a42-ae19-cae00a53735a</ObjectId>
            <UpdateTime>2026-06-09T18:08:12.90394089+08:00</UpdateTime>
            <ClusterType>knowledge</ClusterType>
            <Name>new-demo</Name>
            <Description>this is a demo</Description>
            <Rules>
                <Rule>
                    <Keywords>hello</Keywords>
                    <Keywords>world</Keywords>
                    <Sensitivity>0.7</Sensitivity>
                    <RuleType>keywords</RuleType>
                </Rule>
            </Rules>
            <Reason></Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T15:21:41.288773397+08:00</CreateTime>
            <ObjectId>SmartCluster-fed69d7a-683a-4452-9d70-54148bd458e9</ObjectId>
            <UpdateTime>2026-06-09T15:21:41.288773397+08:00</UpdateTime>
            <ClusterType>knowledge</ClusterType>
            <Name>demo</Name>
            <Rules>
                <Rule>
                    <Keywords>hello</Keywords>
                    <Keywords>world</Keywords>
                    <Sensitivity>0.5</Sensitivity>
                    <RuleType>keywords</RuleType>
                </Rule>
            </Rules>
            <Reason></Reason>
        </SmartCluster>
    </SmartClusters>
</ListSmartClustersResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listSmartClusters&datasetName=your_dataset&metaQuery", strUrl)
		},
		&ListSmartClustersRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *ListSmartClustersResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.SmartClusters), 3)
			assert.Equal(t, *o.SmartClusters[0].Name, "demo-2")
			assert.Equal(t, *o.SmartClusters[0].ClusterType, "knowledge")
			assert.Equal(t, *o.SmartClusters[0].CreateTime, "2026-06-10T09:54:27.484585901+08:00")
			assert.Equal(t, *o.SmartClusters[0].UpdateTime, "2026-06-10T09:54:27.484585901+08:00")
			assert.Equal(t, *o.SmartClusters[0].ObjectId, "SmartCluster-cb9f8c95-281f-490b-b677-eca2f6ff0c19")
			assert.Equal(t, *o.SmartClusters[0].Reason, "")
			assert.Equal(t, o.SmartClusters[0].Rules[0].Keywords[0], "cat")
			assert.Equal(t, *o.SmartClusters[0].Rules[0].RuleType, "keywords")
			assert.Equal(t, *o.SmartClusters[0].Rules[0].Sensitivity, float64(0.5))
			assert.Equal(t, *o.SmartClusters[1].Name, "new-demo")
			assert.Equal(t, *o.SmartClusters[1].ClusterType, "knowledge")
			assert.Equal(t, *o.SmartClusters[1].CreateTime, "2026-06-09T16:18:20.998039823+08:00")
			assert.Equal(t, *o.SmartClusters[1].UpdateTime, "2026-06-09T18:08:12.90394089+08:00")
			assert.Equal(t, *o.SmartClusters[1].ObjectId, "SmartCluster-c30d039b-0b55-4a42-ae19-cae00a53735a")
			assert.Equal(t, *o.SmartClusters[1].Reason, "")
			assert.Equal(t, o.SmartClusters[1].Rules[0].Keywords[0], "hello")
			assert.Equal(t, *o.SmartClusters[1].Rules[0].RuleType, "keywords")
			assert.Equal(t, *o.SmartClusters[1].Rules[0].Sensitivity, float64(0.7))
			assert.Equal(t, *o.SmartClusters[2].Name, "demo")
			assert.Equal(t, *o.SmartClusters[2].ClusterType, "knowledge")
			assert.Equal(t, *o.SmartClusters[2].CreateTime, "2026-06-09T15:21:41.288773397+08:00")
			assert.Equal(t, *o.SmartClusters[2].UpdateTime, "2026-06-09T15:21:41.288773397+08:00")
			assert.Equal(t, *o.SmartClusters[2].ObjectId, "SmartCluster-fed69d7a-683a-4452-9d70-54148bd458e9")
			assert.Equal(t, *o.SmartClusters[2].Reason, "")
			assert.Equal(t, o.SmartClusters[2].Rules[0].Keywords[0], "hello")
			assert.Equal(t, *o.SmartClusters[2].Rules[0].RuleType, "keywords")
			assert.Equal(t, *o.SmartClusters[2].Rules[0].Sensitivity, float64(0.5))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<ListSmartClustersResponse>
    <SmartClusters>
        <SmartCluster>
            <CreateTime>2026-06-10T09:32:24.217552788+08:00</CreateTime>
            <ObjectId>FigureCluster-cd9da94f-feed-45cb-94df-fdcf3e21d712</ObjectId>
            <UpdateTime>2026-06-10T09:32:25.133798431+08:00</UpdateTime>
            <ClusterType>figure</ClusterType>
            <Name>demo-5</Name>
            <Rules>
                <Rule>
                    <BaseURIs>oss://demo-1889/alice.jpg</BaseURIs>
                    <RuleType>face</RuleType>
                </Rule>
            </Rules>
            <Reason></Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T18:32:46.2522507+08:00</CreateTime>
            <ObjectId>FigureCluster-bd06d605-d469-4339-b256-aa0b832be643</ObjectId>
            <UpdateTime>2026-06-09T18:32:46.2522507+08:00</UpdateTime>
            <ClusterType>figure</ClusterType>
            <Name>demo-1</Name>
            <Rules>
                <Rule>
                    <BaseURIs>oss://demo-1889/OIP-C (1).jpg</BaseURIs>
                    <RuleType>face</RuleType>
                </Rule>
            </Rules>
            <Reason>[InvalidArgument.BaseURIs] The face quality is too low. status: 400, requestId: </Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T18:30:39.050200483+08:00</CreateTime>
            <ObjectId>FigureCluster-faaef950-81c6-4280-9fa4-bade5e94b3d1</ObjectId>
            <UpdateTime>2026-06-09T18:30:39.050200483+08:00</UpdateTime>
            <ClusterType>figure</ClusterType>
            <Name>demo-1</Name>
            <Rules>
                <Rule>
                    <BaseURIs>oss://demo-1889/local-v1.txt</BaseURIs>
                    <RuleType>face</RuleType>
                </Rule>
            </Rules>
            <Reason>*error.OpError : InvalidArgument | File corrupt.</Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T18:24:00.460823026+08:00</CreateTime>
            <ObjectId>FigureCluster-af67c554-10a4-4191-91be-5414e431ec43</ObjectId>
            <UpdateTime>2026-06-09T18:24:00.460823026+08:00</UpdateTime>
            <ClusterType>figure</ClusterType>
            <Name>demo-1</Name>
            <Rules>
                <Rule>
                    <BaseURIs>oss://demo-1889/local.txt</BaseURIs>
                    <RuleType>face</RuleType>
                </Rule>
            </Rules>
            <Reason>*error.OpError : InvalidArgument | File does not exist.</Reason>
        </SmartCluster>
    </SmartClusters>
</ListSmartClustersResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listSmartClusters&clusterType=figure&datasetName=your_dataset&maxResults=10&metaQuery&nextToken=nextToken&ruleTypes=%5B%22face%22%5D", strUrl)
		},
		&ListSmartClustersRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ClusterType: SmartClusterTypeFigure,
			MaxResults:  oss.Ptr(int64(10)),
			RuleTypes:   oss.Ptr(`["face"]`),
			NextToken:   oss.Ptr("nextToken"),
		},
		func(t *testing.T, o *ListSmartClustersResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.SmartClusters), 4)
			assert.Equal(t, *o.SmartClusters[0].Name, "demo-5")
			assert.Equal(t, *o.SmartClusters[0].ClusterType, "figure")
			assert.Equal(t, *o.SmartClusters[0].CreateTime, "2026-06-10T09:32:24.217552788+08:00")
			assert.Equal(t, *o.SmartClusters[0].UpdateTime, "2026-06-10T09:32:25.133798431+08:00")
			assert.Equal(t, *o.SmartClusters[0].ObjectId, "FigureCluster-cd9da94f-feed-45cb-94df-fdcf3e21d712")
			assert.Equal(t, *o.SmartClusters[0].Reason, "")
			assert.Equal(t, o.SmartClusters[0].Rules[0].BaseURIs[0], "oss://demo-1889/alice.jpg")
			assert.Equal(t, *o.SmartClusters[0].Rules[0].RuleType, "face")
			assert.Equal(t, *o.SmartClusters[1].Name, "demo-1")
			assert.Equal(t, *o.SmartClusters[1].ClusterType, "figure")
			assert.Equal(t, *o.SmartClusters[1].CreateTime, "2026-06-09T18:32:46.2522507+08:00")
			assert.Equal(t, *o.SmartClusters[1].UpdateTime, "2026-06-09T18:32:46.2522507+08:00")
			assert.Equal(t, *o.SmartClusters[1].ObjectId, "FigureCluster-bd06d605-d469-4339-b256-aa0b832be643")
			assert.Equal(t, *o.SmartClusters[1].Reason, "[InvalidArgument.BaseURIs] The face quality is too low. status: 400, requestId: ")
			assert.Equal(t, o.SmartClusters[1].Rules[0].BaseURIs[0], "oss://demo-1889/OIP-C (1).jpg")
			assert.Equal(t, *o.SmartClusters[1].Rules[0].RuleType, "face")
			assert.Equal(t, *o.SmartClusters[2].Name, "demo-1")
			assert.Equal(t, *o.SmartClusters[2].ClusterType, "figure")
			assert.Equal(t, *o.SmartClusters[2].CreateTime, "2026-06-09T18:30:39.050200483+08:00")
			assert.Equal(t, *o.SmartClusters[2].UpdateTime, "2026-06-09T18:30:39.050200483+08:00")
			assert.Equal(t, *o.SmartClusters[2].ObjectId, "FigureCluster-faaef950-81c6-4280-9fa4-bade5e94b3d1")
			assert.Equal(t, *o.SmartClusters[2].Reason, "*error.OpError : InvalidArgument | File corrupt.")
			assert.Equal(t, o.SmartClusters[2].Rules[0].BaseURIs[0], "oss://demo-1889/local-v1.txt")
			assert.Equal(t, *o.SmartClusters[2].Rules[0].RuleType, "face")
			assert.Equal(t, *o.SmartClusters[3].Name, "demo-1")
			assert.Equal(t, *o.SmartClusters[3].ClusterType, "figure")
			assert.Equal(t, *o.SmartClusters[3].CreateTime, "2026-06-09T18:24:00.460823026+08:00")
			assert.Equal(t, *o.SmartClusters[3].UpdateTime, "2026-06-09T18:24:00.460823026+08:00")
			assert.Equal(t, *o.SmartClusters[3].ObjectId, "FigureCluster-af67c554-10a4-4191-91be-5414e431ec43")
			assert.Equal(t, *o.SmartClusters[3].Reason, "*error.OpError : InvalidArgument | File does not exist.")
			assert.Equal(t, o.SmartClusters[3].Rules[0].BaseURIs[0], "oss://demo-1889/local.txt")
			assert.Equal(t, *o.SmartClusters[3].Rules[0].RuleType, "face")
		},
	},
}

func TestMockListSmartClusters_Success(t *testing.T) {
	for _, c := range testMockListSmartClustersSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListSmartClusters(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListSmartClustersErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListSmartClustersRequest
	CheckOutputFn  func(t *testing.T, o *ListSmartClustersResult, err error)
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
			assert.Equal(t, "/bucket/?action=listSmartClusters&datasetName=your_dataset&metaQuery", strUrl)
		},
		&ListSmartClustersRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *ListSmartClustersResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=listSmartClusters&datasetName=your_dataset&metaQuery", strUrl)
		},
		&ListSmartClustersRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *ListSmartClustersResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockListSmartClusters_Error(t *testing.T) {
	for _, c := range testMockListSmartClustersErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListSmartClusters(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteSmartClusterSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *DeleteSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&DeleteSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *DeleteSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteSmartCluster_Success(t *testing.T) {
	for _, c := range testMockDeleteSmartClusterSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteSmartClusterErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *DeleteSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&DeleteSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *DeleteSmartClusterResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=deleteSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&DeleteSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *DeleteSmartClusterResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

func TestMockDeleteSmartCluster_Error(t *testing.T) {
	for _, c := range testMockDeleteSmartClusterErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}