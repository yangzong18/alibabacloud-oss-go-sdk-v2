package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/signer"
)

func main() {

	var body []byte
	BucketName := "your bucket name"
	provider := credentials.NewStaticCredentialsProvider("ak", "sk", "")

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := oss.New(cfg)

	// 子资源在缺省子资源列表
	input := &oss.OperationInput{
		OperationName: "GetBucketAcl",
		Bucket:        BucketName,
		Method:        "GET",
		Parameters: map[string]string{
			"acl": "",
		}}
	output, err := client.InvokeOperation(context.TODO(), input)
	body = nil
	if output != nil && output.Body != nil {
		defer output.Body.Close()
		body, _ = io.ReadAll(output.Body)
	}
	fmt.Printf("client.InvokeOperation \ninput:%+v, \noutput: %+v\nbody: %s\nerr:%+v\n", input, output, string(body), err)

	// 子资源不在缺省子资源列表
	input = &oss.OperationInput{
		OperationName: "GetBucketResourceGroup",
		Bucket:        BucketName,
		Method:        "GET",
		Parameters: map[string]string{
			"resourceGroup": "",
		}}
	input.Metadata.Set(signer.SubResource, []string{"resourceGroup"})
	output, err = client.InvokeOperation(context.TODO(), input)
	body = nil
	if output != nil && output.Body != nil {
		defer output.Body.Close()
		body, _ = io.ReadAll(output.Body)
	}
	fmt.Printf("client.InvokeOperation \ninput:%+v, \noutput: %+v\nbody: %s\nerr:%+v\n", input, output, string(body), err)

	// 通过PutObject上传
	input = &oss.OperationInput{
		OperationName: "PutObject",
		Bucket:        BucketName,
		Key:           "test-key.txt",
		Method:        "PUT",
		Body:          strings.NewReader("hello world"),
	}
	output, err = client.InvokeOperation(context.TODO(), input)
	body = nil
	if output != nil && output.Body != nil {
		defer output.Body.Close()
		body, _ = io.ReadAll(output.Body)
	}
	fmt.Printf("client.InvokeOperation \ninput:%+v, \noutput: %+v\nbody: %s\nerr:%+v\n", input, output, string(body), err)

	// 通过GetObject获取数据
	input = &oss.OperationInput{
		OperationName: "GetObject",
		Bucket:        BucketName,
		Key:           "test-key.txt",
		Method:        "GET",
	}
	output, err = client.InvokeOperation(context.TODO(), input)
	body = nil
	if output != nil && output.Body != nil {
		defer output.Body.Close()
		body, _ = io.ReadAll(output.Body)
	}
	fmt.Printf("client.InvokeOperation \ninput:%+v, \noutput: %+v\nbody: %s\nerr:%+v\n", input, output, string(body), err)
}
