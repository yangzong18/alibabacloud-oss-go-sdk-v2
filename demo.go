package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/signer"
)

func main() {

	var body []byte
	BucketName := "bucket-test"
	Key := "dump.tif"
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

	//使用File-Like 顺序读模式
	file, err := client.OpenFile(context.TODO(), BucketName, Key)
	if err != nil {
		fmt.Printf("client.OpenFile fail:%v\n", err)
		return
	}

	wfile, err := os.Create("dump-seq.dat")
	if err != nil {
		fmt.Printf("os.Create fail:%v\n", err)
		return
	}

	stat, _ := file.Stat()
	fmt.Printf("file name:%s, modTime:%v, size:%v\n", stat.Name(), stat.ModTime(), stat.Size())

	start := time.Now()

	hash := oss.NewCRC64(0)
	offset, _ := file.Seek(128, io.SeekStart)
	fmt.Printf("new offset:%v\n", offset)
	wfile.Seek(128, io.SeekStart)
	io.Copy(io.MultiWriter(wfile, hash), file)

	hash1 := oss.NewCRC64(0)
	file.Seek(0, io.SeekStart)
	wfile.Seek(0, io.SeekStart)
	io.CopyN(io.MultiWriter(wfile, hash1), file, 128)

	file.Close()
	wfile.Close()

	crc64 := oss.CRC64Combine(hash1.Sum64(), hash.Sum64(), uint64(stat.Size())-128)

	duration := time.Now().Sub(start)
	averSpeed := float64(stat.Size()/1024) / float64(duration.Seconds())

	fmt.Printf("averSpeed :%.2f(KB/s), duration:%v\n", averSpeed, duration)
	fmt.Printf("File-Like seq read done, src file crc64:%v, dest file crc64:%v\n", (stat.Sys().(http.Header)).Get(oss.HeaderOssCRC64), crc64)

	//使用Async Reader 接口
	getFn := func(ctx context.Context, httpRange oss.HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		request := &oss.GetObjectRequest{
			Bucket: oss.Ptr(BucketName),
			Key:    oss.Ptr(Key),
		}
		rangeStr := httpRange.FormatHTTPRange()
		if rangeStr != nil {
			request.Range = rangeStr
			request.RangeBehavior = oss.Ptr("standard")
		}
		result, err := client.GetObject(ctx, request)
		if err != nil {
			return nil, 0, "", err
		}
		offset, _ = oss.ParseOffsetAndSizeFromHeaders(result.Headers)
		return result.Body, offset, result.Headers.Get("ETag"), nil
	}

	//part 1
	hash = oss.NewCRC64(0)
	reader, err := oss.NewAsyncRangeReader(context.TODO(), getFn, &oss.HTTPRange{Offset: 0, Count: 15597568}, "", 4)
	if err != nil {
		fmt.Printf("error")
	}

	wfile, err = os.Create("dump-async-reader.dat")
	io.Copy(io.MultiWriter(wfile, hash), reader)

	//part 2
	hash1 = oss.NewCRC64(hash.Sum64())
	reader1, err := oss.NewAsyncRangeReader(context.TODO(), getFn, &oss.HTTPRange{Offset: 15597568, Count: 8362724}, "", 4)
	if err != nil {
		fmt.Printf("error")
	}
	io.Copy(io.MultiWriter(wfile, hash1), reader1)
	reader.Close()
	wfile.Close()

	fmt.Printf("Async Reader done, dest file crc64:%v", hash1.Sum64())

	//使用File-Like异步并发读模式
	file, err = client.OpenFile(context.TODO(), BucketName, Key, func(oo *oss.OpenOptions) {
		oo.EnablePrefetch = true
		oo.PrefetchThreshold = int64(0)
	})
	if err != nil {
		fmt.Printf("client.OpenFile fail:%v\n", err)
		return
	}

	wfile, err = os.Create("dump-parallel.dat")
	if err != nil {
		fmt.Printf("os.Create fail:%v\n", err)
		return
	}

	stat, _ = file.Stat()
	fmt.Printf("file name:%s, modTime:%v, size:%v\n", stat.Name(), stat.ModTime(), stat.Size())

	start = time.Now()

	hash = oss.NewCRC64(0)
	io.Copy(io.MultiWriter(wfile, hash), file)

	file.Close()
	wfile.Close()
	duration = time.Now().Sub(start)
	averSpeed = float64(stat.Size()/1024) / float64(duration.Seconds())
	fmt.Printf("averSpeed :%.2f(KB/s), duration:%v\n", averSpeed, duration)
	fmt.Printf("File-Like parallel read done, src file crc64:%v, dest file crc64:%v\n", (stat.Sys().(http.Header)).Get(oss.HeaderOssCRC64), hash.Sum64())

	//使用File-Like Append 模式
	afile, err := client.AppendFile(context.TODO(), BucketName, Key+"-append")
	n, err := afile.Write([]byte("hello"))
	fmt.Printf("append file: wirte reulst :%d\n", n)
	n, err = afile.Write([]byte("world"))
	fmt.Printf("append file: wirte reulst :%d\n", n)
	afile.Close()
}
