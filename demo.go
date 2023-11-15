package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/signer"
)

func main() {

	var body []byte
	BucketName := "bucket-test"
	Key := "key-test"
	provider := credentials.NewEnvironmentVariableCredentialsProvider()

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(provider).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := oss.NewClient(cfg)

	// 子资源在缺省子资源列表
	input := &oss.OperationInput{
		OpName: "GetBucketAcl",
		Bucket: oss.Ptr(BucketName),
		Method: "GET",
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
		OpName: "GetBucketResourceGroup",
		Bucket: oss.Ptr(BucketName),
		Method: "GET",
		Parameters: map[string]string{
			"resourceGroup": "",
		}}
	input.OpMetadata.Set(signer.SubResource, []string{"resourceGroup"})
	output, err = client.InvokeOperation(context.TODO(), input)
	body = nil
	if output != nil && output.Body != nil {
		defer output.Body.Close()
		body, _ = io.ReadAll(output.Body)
	}
	fmt.Printf("client.InvokeOperation \ninput:%+v, \noutput: %+v\nbody: %s\nerr:%+v\n", input, output, string(body), err)

	// 通过PutObject上传
	input = &oss.OperationInput{
		OpName: "PutObject",
		Bucket: oss.Ptr(BucketName),
		Key:    oss.Ptr("test-key.txt"),
		Method: "PUT",
		Body:   strings.NewReader("hello world"),
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
		OpName: "GetObject",
		Bucket: oss.Ptr(BucketName),
		Key:    oss.Ptr("test-key.txt"),
		Method: "GET",
	}
	output, err = client.InvokeOperation(context.TODO(), input)
	body = nil
	if output != nil && output.Body != nil {
		defer output.Body.Close()
		body, _ = io.ReadAll(output.Body)
	}
	fmt.Printf("client.InvokeOperation \ninput:%+v, \noutput: %+v\nbody: %s\nerr:%+v\n", input, output, string(body), err)

	// 使用基础接口
	listObjectsRequest := &oss.ListObjectsRequest{
		Bucket: oss.Ptr(BucketName),
		Prefix: oss.Ptr("skyranch"),
	}
	listObjectsResult, err := client.ListObjects(context.TODO(), listObjectsRequest)
	fmt.Printf("client.ListObjects \nrequest:%+v, \nresult: %+v \nerr:%+v\n", listObjectsRequest, listObjectsResult, err)

	// 使用Pageinators
	pageinators := client.NewListObjectsPaginator(
		&oss.ListObjectsRequest{
			Bucket: oss.Ptr(BucketName),
		},
		func(o *oss.PaginatorOptions) {
			o.Limit = 1
		},
	)

	for pageinators.HasNext() {
		result, err := pageinators.NextPage(context.TODO())

		if err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}

		for _, o := range result.Contents {
			fmt.Printf("Key:%v\n", *o.Key)
		}
	}

	//使用File-Like 接口
	file, err := client.OpenFile(BucketName, Key)
	if err != nil {
		fmt.Printf("client.OpenFile fail:%v\n", err)
		return
	}

	wfile, err := os.Create("dump.dat")
	if err != nil {
		fmt.Printf("os.Create fail:%v\n", err)
		return
	}

	stat, _ := file.Stat()
	fmt.Printf("file name:%s, modTime:%v, size:%v\n", stat.Name(), stat.ModTime(), stat.Size())

	offset, _ := file.Seek(128, io.SeekStart)
	wfile.Seek(128, io.SeekStart)
	fmt.Printf("new offset:%v\n", offset)
	io.Copy(wfile, file)

	file.Seek(0, io.SeekStart)
	wfile.Seek(0, io.SeekStart)
	io.CopyN(wfile, file, 128)

	file.Close()
	wfile.Close()
}
