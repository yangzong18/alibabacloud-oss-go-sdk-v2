package oss

import (
	"context"
	"io"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

type LiveChannelPlayUrls struct {
	// The playback URL.
	Url *string `xml:"Url"`
}

type CreateLiveChannelResult struct {
	// The container that stores the URL used to ingest streams to the LiveChannel.
	PublishUrls *LiveChannelPublishUrls `xml:"PublishUrls"`

	// The container that stores the URL used to play the streams ingested to the LiveChannel.
	PlayUrls *LiveChannelPlayUrls `xml:"PlayUrls"`
}

type LiveChannelHistory struct {
	// The container that stores a list of stream pushing records.
	LiveRecords []LiveRecord `xml:"LiveRecord"`
}

type LiveChannelAudio struct {
	// The bitrate of the current audio stream.  Bandwidth indicates the average bitrate of the audio stream or video stream in the recent period. When LiveChannel is switched to the Live state, the returned value of Bandwidth may be 0. Unit: B/s.
	Bandwidth *int64 `xml:"Bandwidth"`

	// The sample rate of the current audio stream.
	SampleRate *int64 `xml:"SampleRate"`

	// The encoding format of the current audio stream.
	Codec *string `xml:"Codec"`
}

type LiveChannelVideo struct {
	// The width of the video stream. Unit: pixels.
	Width *int64 `xml:"Width"`

	// The height of the video stream. Unit: pixels.
	Height *int64 `xml:"Height"`

	// The frame rate of the video stream.
	FrameRate *int64 `xml:"FrameRate"`

	// The bitrate of the video stream. Unit: bit/s.
	Bandwidth *int64 `xml:"Bandwidth"`

	// The encoding format of the video stream.
	Codec *string `xml:"Codec"`
}

type LiveChannelStat struct {
	// The container that stores audio stream information if Status is set to Live.Video and audio containers can be returned only if Status is set to Live. However, these two containers may not necessarily be returned if Status is set to Live. For example, if the client has connected to the LiveChannel but no audio or video stream is sent, these two containers are not returned.
	Audio *LiveChannelAudio `xml:"Audio"`

	// The current stream ingestion status of the LiveChannel. Valid value: Disabled、Live、Idle。
	Status *string `xml:"Status"`

	// If Status is set to Live, this element indicates the time when the current client starts to ingest streams. The value of the element is in the ISO 8601 format.
	ConnectedTime *string `xml:"ConnectedTime"`

	// If Status is set to Live, this element indicates the IP address of the current client that ingests streams.
	RemoteAddr *string `xml:"RemoteAddr"`

	// The container that stores video stream information if Status is set to Live.Video and audio containers can be returned only if Status is set to Live. However, these two containers may not necessarily be returned if Status is set to Live. For example, if the client has connected to the LiveChannel but no audio or video stream is sent, these two containers are not returned.
	Video *LiveChannelVideo `xml:"Video"`
}

type LiveChannelTarget struct {
	// The duration of each TS file when the value of Type is HLS. Unit: seconds.Valid values: \[1, 100]. Default value: **5**  If you do not specify values for the FragDuration and FragCount parameters, the default values of the two parameters are used. If you specify one of the parameters, you must also specify the other.
	FragDuration *int64 `xml:"FragDuration"`

	// The number of TS files included in the M3U8 file when the value of Type is HLS.Valid values: \[1, 100] Default value: **3**  If you do not specify values for the FragDuration and FragCount parameters, the default values of the two parameters are used. If you specify one of the parameters, you must also specify the other.
	FragCount *int64 `xml:"FragCount"`

	// The name of the generated M3U8 file when the value of Type is HLS. The name must be 6 to 128 bytes in length and end with .m3u8.Default value: **playlist.m3u8** Valid values: \[6, 128]
	PlaylistName *string `xml:"PlaylistName"`

	// The format in which the LiveChannel stores uploaded data.Valid value: **HLS** *   When you set the value of Type to HLS, Object Storage Service (OSS) updates the M3U8 file each time when a TS file is generated. The maximum number of the latest .ts files that can be included in the M3U8 file is specified by the FragCount parameter.*   If you set the value of Type to HLS and the duration of the audio and video data written to the current TS file exceeds the duration specified by FragDuration, OSS switches to the next TS file when the next key frame is received. If OSS does not receive the next key frame after 60 seconds or twice the duration specified by FragDuration (whichever is greater), OSS forcibly switches to the next TS file. In this case, stuttering may occur during the playback of the stream.
	Type *string `xml:"Type"`
}

type LiveChannelSnapshot struct {
	// The name of the role used to perform high-frequency snapshot operations. The role must have the write permissions on DestBucket and the permissions to send messages to NotifyTopic.
	RoleName *string `xml:"RoleName"`

	// The bucket that stores the results of high-frequency snapshot operations. The bucket must belong to the same owner as the current bucket.
	DestBucket *string `xml:"DestBucket"`

	// The Message Service (MNS) topic used to notify users of the results of high-frequency snapshot operations.
	NotifyTopic *string `xml:"NotifyTopic"`

	// The interval of high-frequency snapshot operations. If no key frame (inline frame) exists within the interval, no snapshot is captured. Unit: seconds. Valid values: [1, 100].
	Interval *int64 `xml:"Interval"`
}

type LiveChannelPublishUrls struct {
	// The URL used to ingest streams to the LiveChannel. *   The URL used to ingest streams is not signed. If the ACL of the bucket is not public-read-write, you must add a signature to the URL before you use the URL to access the bucket.*   The URL used to play streams is not signed. If the ACL of the bucket is private, you must add a signature to the URL before you use the URL to access the bucket.
	Url *string `xml:"Url"`
}

type LiveRecord struct {
	// The start time of stream ingest, which is in the ISO8601 format.
	StartTime *string `xml:"StartTime"`

	// The end time of stream ingest, which is in the ISO8601 format.
	EndTime *string `xml:"EndTime"`

	// The IP address of the stream ingest client.
	RemoteAddr *string `xml:"RemoteAddr"`
}

type LiveChannelConfiguration struct {
	// The description of the LiveChannel.
	Description *string `xml:"Description"`

	// The status of the LiveChannel.Valid values:- enabled: indicates that the LiveChannel is enabled.- disabled: indicates that the LiveChannel is disabled.
	Status *string `xml:"Status"`

	// The container that stores the configurations used by the LiveChannel to store uploaded data. FragDuration, FragCount, and PlaylistName are returned only when the value of Type is HLS.
	Target *LiveChannelTarget `xml:"Target"`

	// The container that stores the options of the high-frequency snapshot operations.
	Snapshot *LiveChannelSnapshot `xml:"Snapshot"`
}

type LiveChannel struct {
	// The container that stores the URL used to ingest a stream to the LiveChannel.
	PublishUrls *LiveChannelPublishUrls `xml:"PublishUrls"`

	// The container that stores the URL used to play a stream ingested to the LiveChannel.
	PlayUrls *LiveChannelPlayUrls `xml:"PlayUrls"`

	// The name of the LiveChannel.
	Name *string `xml:"Name"`

	// The description of the LiveChannel.
	Description *string `xml:"Description"`

	// The status of the LiveChannel. Valid values:*   disabled*   enabled
	Status *string `xml:"Status"`

	// The time when the LiveChannel configuration is last modified. Standard: ISO 8601.
	LastModified *string `xml:"LastModified"`
}

type PutLiveChannelRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the LiveChannel. The name cannot contain forward slashes (/).
	Channel *string `input:"path,channel,required"`

	// The request body schema.
	LiveChannelConfiguration *LiveChannelConfiguration `input:"body,LiveChannelConfiguration,xml,required"`

	RequestCommon
}

type PutLiveChannelResult struct {
	// The container that stores the result of the CreateLiveChannel request.
	CreateLiveChannelResult *CreateLiveChannelResult `output:"body,CreateLiveChannelResult,xml"`

	ResultCommon
}

// PutLiveChannel You can call this operation to create a LiveChannel. Before you can upload audio and video data by using the Real-Time Messaging Protocol (RTMP), you must call the PutLiveChannel operation to create a LiveChannel.
func (c *Client) PutLiveChannel(ctx context.Context, request *PutLiveChannelRequest, optFns ...func(*Options)) (*PutLiveChannelResult, error) {
	var err error
	if request == nil {
		request = &PutLiveChannelRequest{}
	}
	input := &OperationInput{
		OpName: "PutLiveChannel",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"live": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"live"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &PutLiveChannelResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type ListLiveChannelRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the LiveChannel from which the list operation starts. LiveChannels whose names are alphabetically after the value of the marker parameter are returned.
	Marker *string `input:"query,marker"`

	// The maximum number of LiveChannels that can be returned for the current request. The value of max-keys cannot exceed 1000. Default value: 100.
	MaxKeys int64 `input:"query,max-keys"`

	// The prefix that the names of the LiveChannels that you want to return must contain. If you specify a prefix in the request, the specified prefix is included in the response.
	Prefix *string `input:"query,prefix"`

	RequestCommon
}

type ListLiveChannelResult struct {
	// If not all results are returned, the NextMarker parameter is included in the response to indicate the Marker value of the next request.
	NextMarker *string `xml:"NextMarker"`

	// The container that stores the information about each returned LiveChannel.
	LiveChannels []LiveChannel `xml:"LiveChannel"`

	// The prefix that the names of the returned LiveChannels contain.
	Prefix *string `xml:"Prefix"`

	// The name of the LiveChannel after which the ListLiveChannel operation starts.
	Marker *string `xml:"Marker"`

	// The maximum number of returned LiveChannels in the response.
	MaxKeys *int64 `xml:"MaxKeys"`

	// Indicates whether all results are returned.- true: All results are returned.- false: Not all results are returned.
	IsTruncated *bool `xml:"IsTruncated"`

	ResultCommon
}

// ListLiveChannel You can call this operation to list specified LiveChannels.
func (c *Client) ListLiveChannel(ctx context.Context, request *ListLiveChannelRequest, optFns ...func(*Options)) (*ListLiveChannelResult, error) {
	var err error
	if request == nil {
		request = &ListLiveChannelRequest{}
	}
	input := &OperationInput{
		OpName: "ListLiveChannel",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"live": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"live"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &ListLiveChannelResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type DeleteLiveChannelRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of an existing LiveChannel.
	Channel *string `input:"path,channel,required"`

	RequestCommon
}

type DeleteLiveChannelResult struct {
	ResultCommon
}

// DeleteLiveChannel You can call this operation to delete a specified LiveChannel.
func (c *Client) DeleteLiveChannel(ctx context.Context, request *DeleteLiveChannelRequest, optFns ...func(*Options)) (*DeleteLiveChannelResult, error) {
	var err error
	if request == nil {
		request = &DeleteLiveChannelRequest{}
	}
	input := &OperationInput{
		OpName: "DeleteLiveChannel",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"live": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"live"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &DeleteLiveChannelResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type PutLiveChannelStatusRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of an existing LiveChannel.
	Channel *string `input:"path,channel,required"`

	// The status of the LiveChannel. Valid values:- enabled: enables the LiveChannel.- disabled: disables the LiveChannel.
	Status *string `input:"query,status,required"`

	RequestCommon
}

type PutLiveChannelStatusResult struct {
	ResultCommon
}

// PutLiveChannelStatus You can call this operation to switch the status of a LiveChannel. A LiveChannel can be in one of the following states: enabled or disabled.
func (c *Client) PutLiveChannelStatus(ctx context.Context, request *PutLiveChannelStatusRequest, optFns ...func(*Options)) (*PutLiveChannelStatusResult, error) {
	var err error
	if request == nil {
		request = &PutLiveChannelStatusRequest{}
	}
	input := &OperationInput{
		OpName: "PutLiveChannelStatus",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"live": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"live"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &PutLiveChannelStatusResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetLiveChannelInfoRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the LiveChannel about which you want to query configuration information. The name cannot contain forward slashes (/).
	Channel *string `input:"path,channel,required"`

	RequestCommon
}

type GetLiveChannelInfoResult struct {
	// The container that stores the returned results of the GetLiveChannelInfo request.
	LiveChannelConfiguration *LiveChannelConfiguration `output:"body,LiveChannelConfiguration,xml"`

	ResultCommon
}

// GetLiveChannelInfo You can call this operation to query the configuration information about a LiveChannel.
func (c *Client) GetLiveChannelInfo(ctx context.Context, request *GetLiveChannelInfoRequest, optFns ...func(*Options)) (*GetLiveChannelInfoResult, error) {
	var err error
	if request == nil {
		request = &GetLiveChannelInfoRequest{}
	}
	input := &OperationInput{
		OpName: "GetLiveChannelInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"live": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"live"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &GetLiveChannelInfoResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetLiveChannelHistoryRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the LiveChannel.
	Channel *string `input:"path,channel,required"`

	RequestCommon
}

type GetLiveChannelHistoryResult struct {
	// The container that stores the returned results of the GetLiveChannelHistory request.
	LiveChannelHistory *LiveChannelHistory `output:"body,LiveChannelHistory,xml"`

	ResultCommon
}

// GetLiveChannelHistory You can call this operation to query the stream ingestion records of a LiveChannel.
func (c *Client) GetLiveChannelHistory(ctx context.Context, request *GetLiveChannelHistoryRequest, optFns ...func(*Options)) (*GetLiveChannelHistoryResult, error) {
	var err error
	if request == nil {
		request = &GetLiveChannelHistoryRequest{}
	}
	input := &OperationInput{
		OpName: "GetLiveChannelHistory",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp": "history",
			"live": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"live", "comp"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &GetLiveChannelHistoryResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetLiveChannelStatRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the LiveChannel.
	Channel *string `input:"path,channel,required"`

	RequestCommon
}

type GetLiveChannelStatResult struct {
	// The container that stores the returned results of the GetLiveChannelStat request.
	LiveChannelStat *LiveChannelStat `output:"body,LiveChannelStat,xml"`

	ResultCommon
}

// GetLiveChannelStat You can call this operation to query the stream ingestion status of a LiveChannel.
func (c *Client) GetLiveChannelStat(ctx context.Context, request *GetLiveChannelStatRequest, optFns ...func(*Options)) (*GetLiveChannelStatResult, error) {
	var err error
	if request == nil {
		request = &GetLiveChannelStatRequest{}
	}
	input := &OperationInput{
		OpName: "GetLiveChannelStat",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"comp": "stat",
			"live": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"live", "comp"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &GetLiveChannelStatResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetVodPlaylistRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of an existing LiveChannel.
	Channel *string `input:"path,channel,required"`

	// The end time of the time range during which the TS files that you want to query are generated in the Unix timestamp format.  The value of EndTime must be greater than the value of StartTime. The duration between EndTime and StartTime must be less than one day.
	EndTime *string `input:"query,endTime,required"`

	// The start time of the time range during which the TS files that you want to query are generated in the Unix timestamp format.
	StartTime *string `input:"query,startTime,required"`

	RequestCommon
}

type GetVodPlaylistResult struct {
	Body io.ReadCloser

	ResultCommon
}

// GetVodPlaylist You can call this operation to query the playlist that is generated by the streams ingested to the specified LiveChannel within the specified time range.
func (c *Client) GetVodPlaylist(ctx context.Context, request *GetVodPlaylistRequest, optFns ...func(*Options)) (*GetVodPlaylistResult, error) {
	var err error
	if request == nil {
		request = &GetVodPlaylistRequest{}
	}
	input := &OperationInput{
		OpName: "GetVodPlaylist",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"vod": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"vod"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &GetVodPlaylistResult{
		Body: output.Body,
	}

	if err = c.unmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type PostVodPlaylistRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of an existing LiveChannel.
	Channel *string `input:"path,channel,required"`

	// The name of the generated VOD playlist, which must end with ".m3u8".
	Playlist *string `input:"path,playlist,required"`

	// The end time of the time range during which the TS files that you want to query are generated, which is a Unix timestamp. The value of EndTime must be later than the value of StartTime. The duration between EndTime and StartTime must be shorter than one day.
	EndTime *string `input:"query,endTime,required"`

	// The start time of the time range during which the TS files that you want to query are generated, which is a Unix timestamp.
	StartTime *string `input:"query,startTime,required"`

	RequestCommon
}

type PostVodPlaylistResult struct {
	ResultCommon
}

// PostVodPlaylist You can call this operation to generate a VOD playlist for the specified LiveChannel.
func (c *Client) PostVodPlaylist(ctx context.Context, request *PostVodPlaylistRequest, optFns ...func(*Options)) (*PostVodPlaylistResult, error) {
	var err error
	if request == nil {
		request = &PostVodPlaylistRequest{}
	}
	input := &OperationInput{
		OpName: "PostVodPlaylist",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"vod": "",
		},
		Bucket: request.Bucket,
	}

	input.OpMetadata.Set(signer.SubResource, []string{"vod"})

	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}
	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &PostVodPlaylistResult{}

	if err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
