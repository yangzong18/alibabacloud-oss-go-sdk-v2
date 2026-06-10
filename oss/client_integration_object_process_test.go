//go:build integration

package oss

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessObject(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6) + ".jpg"
	objectDestName := objectNamePrefix + randLowStr(6) + "dest.jpg"

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}

	filePath := "../sample/example.jpg"
	_, err = client.PutObjectFromFile(context.TODO(), putObjRequest, filePath)
	assert.Nil(t, err)

	request := &ProcessObjectRequest{
		Bucket:  Ptr(bucketName),
		Key:     Ptr(objectName),
		Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte(objectDestName)))),
	}
	result, err := client.ProcessObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, result.Bucket, "")
	assert.NotEmpty(t, result.FileSize)
	assert.Equal(t, result.Object, objectDestName)
	assert.Equal(t, result.ProcessStatus, "OK")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &ProcessObjectRequest{
		Bucket:  Ptr(bucketNameNotExist),
		Key:     Ptr(objectName),
		Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte(objectDestName)))),
	}
	result, err = client.ProcessObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAsyncProcessObject(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6) + ".mp4"
	objectDestName := objectNamePrefix + randLowStr(6) + "dest.mp4"

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	putObjrequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	videoUrl := "https://oss-console-img-demo-cn-hangzhou.oss-cn-hangzhou.aliyuncs.com/video.mp4?spm=a2c4g.64555.0.0.515675979u4B8w&file=video.mp4"
	fileName := "video.mp4"
	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = http.Get(videoUrl)
		if err != nil {
			continue
		}
	}
	assert.Nil(t, err)
	defer resp.Body.Close()
	defer os.Remove(fileName)

	file, err := os.Create(fileName)
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	assert.Nil(t, err)
	_, err = client.PutObjectFromFile(context.TODO(), putObjrequest, fileName)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	style := "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1"
	process := fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", style, strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(bucketName)), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(objectDestName)), "="))
	request := &AsyncProcessObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		AsyncProcess: Ptr(process),
	}
	var serr *ServiceError
	_, err = client.AsyncProcessObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, "Imm Client", serr.Code)
	assert.Contains(t, serr.Message, "ResourceNotFound, The specified resource Attachment is not found")
	assert.NotEmpty(t, serr.RequestID)

	time.Sleep(1 * time.Second)
	bucketNameNotExist := bucketName + "-not-exist"
	request = &AsyncProcessObjectRequest{
		Bucket:       Ptr(bucketNameNotExist),
		Key:          Ptr(objectName),
		AsyncProcess: Ptr(process),
	}
	_, err = client.AsyncProcessObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectWithProcess(t *testing.T) {
	after := before(t)
	defer after(t)

	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	objectName := objectNamePrefix + randLowStr(6) + ".jpg"

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}

	filePath := "../sample/example.jpg"
	_, err = client.PutObjectFromFile(context.TODO(), putObjRequest, filePath)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	style := "image/resize,m_fixed,w_100,h_100/rotate,90"
	getObjRequest := &GetObjectRequest{
		Bucket:  Ptr(bucketName),
		Key:     Ptr(objectName),
		Process: Ptr(style),
	}

	downloadFile := "example-download.jpg"
	_, err = client.GetObjectToFile(context.TODO(), getObjRequest, downloadFile)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	content, err := os.ReadFile(downloadFile)
	assert.Nil(t, err)

	result, err := client.GetObject(context.TODO(), getObjRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	content2, err := io.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, content2, content)

	sign, err := client.Presign(context.TODO(), getObjRequest)
	req, err := http.NewRequest(sign.Method, sign.URL, nil)
	assert.Nil(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	time.Sleep(1 * time.Second)

	content3, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, content3, content)

	os.Remove(downloadFile)
}
