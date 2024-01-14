package oss

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/aliyun-oss-go-sdk-v2/oss/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	matDesc = make(map[string]string)

	rsaPublicKey string = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCokfiAVXXf5ImFzKDw+XO/UByW
6mse2QsIgz3ZwBtMNu59fR5zttSx+8fB7vR4CN3bTztrP9A6bjoN0FFnhlQ3vNJC
5MFO1PByrE/MNd5AAfSVba93I6sx8NSk5MzUCA4NJzAUqYOEWGtGBcom6kEF6MmR
1EKib1Id8hpooY5xaQIDAQAB
-----END PUBLIC KEY-----`

	rsaPrivateKey string = `-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAKiR+IBVdd/kiYXM
oPD5c79QHJbqax7ZCwiDPdnAG0w27n19HnO21LH7x8Hu9HgI3dtPO2s/0DpuOg3Q
UWeGVDe80kLkwU7U8HKsT8w13kAB9JVtr3cjqzHw1KTkzNQIDg0nMBSpg4RYa0YF
yibqQQXoyZHUQqJvUh3yGmihjnFpAgMBAAECgYA49RmCQ14QyKevDfVTdvYlLmx6
kbqgMbYIqk+7w611kxoCTMR9VMmJWgmk/Zic9mIAOEVbd7RkCdqT0E+xKzJJFpI2
ZHjrlwb21uqlcUqH1Gn+wI+jgmrafrnKih0kGucavr/GFi81rXixDrGON9KBE0FJ
cPVdc0XiQAvCBnIIAQJBANXu3htPH0VsSznfqcDE+w8zpoAJdo6S/p30tcjsDQnx
l/jYV4FXpErSrtAbmI013VYkdJcghNSLNUXppfk2e8UCQQDJt5c07BS9i2SDEXiz
byzqCfXVzkdnDj9ry9mba1dcr9B9NCslVelXDGZKvQUBqNYCVxg398aRfWlYDTjU
IoVVAkAbTyjPN6R4SkC4HJMg5oReBmvkwFCAFsemBk0GXwuzD0IlJAjXnAZ+/rIO
ItewfwXIL1Mqz53lO/gK+q6TR585AkB304KUIoWzjyF3JqLP3IQOxzns92u9EV6l
V2P+CkbMPXiZV6sls6I4XppJXX2i3bu7iidN3/dqJ9izQK94fMU9AkBZvgsIPCot
y1/POIbv9LtnviDKrmpkXgVQSU4BmTPvXwTJm8APC7P/horSh3SVf1zgmnsyjm9D
hO92gGc+4ajL
-----END PRIVATE KEY-----`

	rsaPublicKeyPks1 string = `-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBALpUiB+w+r3v2Fgw0SgMbWl8bnzUVc3t3YbA89H13lrw7v6RUbL8+HGl
s5YGoqD4lObG/sCQyaWd0B/XzOhjlSc1b53nyZhms84MGJ6nF2NQP+1gjY1ByDMK
zeyVFFFvl9prlr6XpuJQlY0F/W4pbjLsk8Px4Qix5AoJbShElUu1AgMBAAE=
-----END RSA PUBLIC KEY-----`

	rsaPrivateKeyPks1 string = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQC6VIgfsPq979hYMNEoDG1pfG581FXN7d2GwPPR9d5a8O7+kVGy
/PhxpbOWBqKg+JTmxv7AkMmlndAf18zoY5UnNW+d58mYZrPODBiepxdjUD/tYI2N
QcgzCs3slRRRb5faa5a+l6biUJWNBf1uKW4y7JPD8eEIseQKCW0oRJVLtQIDAQAB
AoGBAJrzWRAhuSLipeMRFZ5cV1B1rdwZKBHMUYCSTTC5amPuIJGKf4p9XI4F4kZM
1klO72TK72dsAIS9rCoO59QJnCpG4CvLYlJ37wA2UbhQ1rBH5dpBD/tv3CUyfdtI
9CLUsZR3DGBWXYwGG0KGMYPExe5Hq3PUH9+QmuO+lXqJO4IBAkEA6iLee6oBzu6v
90zrr4YA9NNr+JvtplpISOiL/XzsU6WmdXjzsFLSsZCeaJKsfdzijYEceXY7zUNa
0/qQh2BKoQJBAMu61rQ5wKtql2oR4ePTSm00/iHoIfdFnBNU+b8uuPXlfwU80OwJ
Gbs0xBHe+dt4uT53QLci4KgnNkHS5lu4XJUCQQCisCvrvcuX4B6BNf+mbPSJKcci
biaJqr4DeyKatoz36mhpw+uAH2yrWRPZEeGtayg4rvf8Jf2TuTOJi9eVWYFBAkEA
uIPzyS81TQsxL6QajpjjI52HPXZcrPOis++Wco0Cf9LnA/tczSpA38iefAETEq94
NxcSycsQ5br97QfyEsgbMQJANTZ/HyMowmDPIC+n9ExdLSrf4JydARSfntFbPsy1
4oC6ciKpRdtAtAtiU8s9eAUSWi7xoaPJzjAHWbmGSHHckg==
-----END RSA PRIVATE KEY-----`
)

type encryptionMockTracker struct {
	lastModified string
	savedata     []byte
	saveHeaders  http.Header
}

func testSetupEncryptionMockServer(t *testing.T, tracker *encryptionMockTracker) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		length := len(tracker.savedata)
		data := tracker.savedata
		switch r.Method {
		case "PUT":
			in, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			hash := NewCRC64(0)
			hash.Write(in)
			crc64ecma := fmt.Sprint(hash.Sum64())

			//save data & headers
			tracker.savedata = in
			tracker.saveHeaders = r.Header

			md5hash := md5.New()
			md5hash.Write(in)
			etag := fmt.Sprintf("\"%s\"", strings.ToUpper(hex.EncodeToString(md5hash.Sum(nil))))

			//PutObject
			w.Header().Set(HeaderOssCRC64, crc64ecma)
			w.Header().Set(HTTPHeaderETag, etag)

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "GET":
			// header
			var httpRange *HTTPRange
			if r.Header.Get("Range") != "" {
				httpRange, _ = ParseRange(r.Header.Get("Range"))
			}

			offset := int64(0)
			statusCode := 200
			sendLen := int64(length)
			if httpRange != nil {
				offset = httpRange.Offset
				sendLen = int64(length) - httpRange.Offset
				if httpRange.Count > 0 {
					sendLen = minInt64(httpRange.Count, sendLen)
				}
				cr := httpContentRange{
					Offset: httpRange.Offset,
					Count:  sendLen,
					Total:  int64(length),
				}
				w.Header().Set("Content-Range", ToString(cr.FormatHTTPContentRange()))
				statusCode = 206
			}

			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(sendLen))
			w.Header().Set(HTTPHeaderLastModified, tracker.lastModified)
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")
			for k, vv := range tracker.saveHeaders {
				lk := strings.ToLower(k)
				if strings.HasPrefix(lk, "x-oss-meta-client-side-encryption-") {
					w.Header().Set(k, vv[0])
				}
			}

			//status code
			w.WriteHeader(statusCode)

			//body
			sendData := data[int(offset):int(offset+sendLen)]
			//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
			w.Write(sendData)
		}
	}))
	return server
}

func TestMockEncryptionPks8(t *testing.T) {
	//length := 123
	data := []byte("hello world")
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &encryptionMockTracker{
		lastModified: gmtTime,
	}
	server := testSetupEncryptionMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	assert.NotNil(t, client)

	mc, err := crypto.CreateMasterRsa(map[string]string{"tag": "value"}, rsaPublicKey, rsaPrivateKey)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)

	result, err := eclient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   bytes.NewReader(data),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, datasum)
	assert.Len(t, tracker.savedata, len(data))
	assert.NotEqualValues(t, data, tracker.savedata)

	gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err := io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.EqualValues(t, data, gData)
	assert.NotEmpty(t, tracker.saveHeaders.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, tracker.saveHeaders.Get(OssClientSideEncryptionStart))
	assert.Equal(t, crypto.AesCtrAlgorithm, tracker.saveHeaders.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, crypto.RsaCryptoWrap, tracker.saveHeaders.Get(OssClientSideEncryptionWrapAlg))
	assert.Equal(t, "{\"tag\":\"value\"}", tracker.saveHeaders.Get(OssClientSideEncryptionMatDesc))
	assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionUnencryptedContentMD5))
	assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionPartSize))

	result, err = eclient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   bytes.NewReader(nil),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, tracker.savedata, 0)
	assert.NotEqualValues(t, data, tracker.savedata)

	gResult, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err = io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.Len(t, gData, 0)
}

func TestMockEncryptionPks1(t *testing.T) {
	var data []byte
	gmtTime := getNowGMT()
	tracker := &encryptionMockTracker{
		lastModified: gmtTime,
	}
	server := testSetupEncryptionMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	assert.NotNil(t, client)

	mc, err := crypto.CreateMasterRsa(nil, rsaPublicKeyPks1, rsaPrivateKeyPks1)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)

	i := 0
	for i = 1; i < 1234; i++ {
		data = []byte(randStr(i))
		result, err := eclient.PutObject(context.TODO(), &PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Body:   bytes.NewReader(data),
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, tracker.savedata, len(data))
		assert.NotEqualValues(t, data, tracker.savedata)
		assert.NotEmpty(t, tracker.saveHeaders.Get(OssClientSideEncryptionKey))
		assert.NotEmpty(t, tracker.saveHeaders.Get(OssClientSideEncryptionStart))
		assert.Equal(t, crypto.AesCtrAlgorithm, tracker.saveHeaders.Get(OssClientSideEncryptionCekAlg))
		assert.Equal(t, crypto.RsaCryptoWrap, tracker.saveHeaders.Get(OssClientSideEncryptionWrapAlg))
		assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionMatDesc))

		gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
		assert.Nil(t, err)
		assert.NotNil(t, gResult)
		gData, err := io.ReadAll(gResult.Body)
		assert.Nil(t, err)
		assert.EqualValues(t, data, gData)
	}
	assert.Equal(t, 1234, i)
}

func TestMockEncryptionGetObjecRange(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	tracker := &encryptionMockTracker{
		lastModified: gmtTime,
	}
	server := testSetupEncryptionMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	assert.NotNil(t, client)

	mc, err := crypto.CreateMasterRsa(nil, rsaPublicKeyPks1, rsaPrivateKeyPks1)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)

	result, err := eclient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   bytes.NewReader(data),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, tracker.savedata, length)
	assert.NotEqualValues(t, data, tracker.savedata)

	gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err := io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.EqualValues(t, data, gData)

	for i := 0; i < 123; i++ {
		for len := 345; len < 456; len++ {
			httpRange := HTTPRange{
				Offset: int64(i),
				Count:  int64(len),
			}
			gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
				Bucket: Ptr("bucket"),
				Key:    Ptr("key"),
				Range:  httpRange.FormatHTTPRange(),
			})
			assert.Nil(t, err)
			assert.NotNil(t, gResult)
			gData, err := io.ReadAll(gResult.Body)
			assert.EqualValues(t, data[i:i+len], gData)
			assert.Equal(t, int64(len), gResult.ContentLength)
			assert.Equal(t, fmt.Sprintf("bytes %v-%v/%v", i, (i+len-1), length), ToString(gResult.ContentRange))
		}
	}
}

func TestMockEncryptionMatDescTest(t *testing.T) {
	data := []byte("hello world")
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &encryptionMockTracker{
		lastModified: gmtTime,
	}
	server := testSetupEncryptionMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	assert.NotNil(t, client)

	mc, err := crypto.CreateMasterRsa(map[string]string{"tag": "rsa"}, rsaPublicKey, rsaPrivateKey)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)

	result, err := eclient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   bytes.NewReader(data),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, datasum)
	assert.Len(t, tracker.savedata, len(data))
	assert.NotEqualValues(t, data, tracker.savedata)

	gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err := io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.EqualValues(t, gData, data)

	//Use other mc, not match
	mc1, err := crypto.CreateMasterRsa(map[string]string{"tag": "rsa1"}, rsaPublicKeyPks1, rsaPrivateKeyPks1)
	eclient1, err := NewEncryptionClient(client, mc1)
	assert.Nil(t, err)
	gResult1, err := eclient1.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.NotNil(t, err)
	assert.Nil(t, gResult1)
	assert.Contains(t, err.Error(), "crypto/rsa: decryption error,object:key")

	//Use mixed mc
	eclient2, err := NewEncryptionClient(client, mc1,
		func(eco *EncryptionClientOptions) {
			eco.MasterCiphers = []crypto.MasterCipher{mc}
		},
	)
	assert.Nil(t, err)
	gResult2, err := eclient2.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData2, err := io.ReadAll(gResult2.Body)
	assert.EqualValues(t, data, gData2)
}

type fakeEncryptionContentCipher struct {
}

func (b fakeEncryptionContentCipher) EncryptContent(io.Reader) (io.ReadCloser, error) {
	return nil, fmt.Errorf("EncryptContent fail")
}

func (b fakeEncryptionContentCipher) DecryptContent(io.Reader) (io.ReadCloser, error) {
	return nil, fmt.Errorf("DecryptContent fail")
}

func (b fakeEncryptionContentCipher) Clone(cd crypto.CipherData) (crypto.ContentCipher, error) {
	return nil, fmt.Errorf("DecryptContent fail")
}

func (b fakeEncryptionContentCipher) GetEncryptedLen(int64) int64 {
	return 0
}

func (b fakeEncryptionContentCipher) GetCipherData() *crypto.CipherData {
	return nil
}

func (b fakeEncryptionContentCipher) GetAlignLen() int {
	return 0
}

type fakeEncryptionContentCipherBuilder struct {
	cc *fakeEncryptionContentCipher
}

func (b fakeEncryptionContentCipherBuilder) ContentCipher() (crypto.ContentCipher, error) {
	if b.cc != nil {
		return *b.cc, nil
	}
	return nil, fmt.Errorf("ContentCipher fail")
}

func (b fakeEncryptionContentCipherBuilder) ContentCipherEnv(crypto.Envelope) (crypto.ContentCipher, error) {
	return nil, fmt.Errorf("ContentCipherEnv fail")
}

func (b fakeEncryptionContentCipherBuilder) GetMatDesc() string {
	return ""
}

func TestEncryptionClientErrorTest(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	assert.NotNil(t, client)

	_, err := NewEncryptionClient(client, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field, masterCipher")

	mc, err := crypto.CreateMasterRsa(map[string]string{"tag": "rsa"}, rsaPublicKey, rsaPrivateKey)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)
	assert.Equal(t, client, eclient.Unwrap())

	_, err = eclient.PutObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field, request")

	_, err = eclient.GetObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field, request")

	_, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Range:  Ptr("invalid-range"),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "range: header invalid")

	_, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "operation error GetObject: invalid field, Endpoint")

	eclient.defualtCCBuilder = fakeEncryptionContentCipherBuilder{}
	_, err = eclient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   strings.NewReader(""),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ContentCipher fail")

	eclient.defualtCCBuilder = fakeEncryptionContentCipherBuilder{
		cc: &fakeEncryptionContentCipher{},
	}
	_, err = eclient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   strings.NewReader(""),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "EncryptContent fail")
}

func TestGetEnvelopeFromHeader(t *testing.T) {
	header := http.Header{}
	header.Set(OssClientSideEncryptionCekAlg, "AES/CTR/NoPadding")
	header.Set(OssClientSideEncryptionMatDesc, "{\"tag\":\"value\"}")
	header.Set(OssClientSideEncryptionStart, "eIR2+1Ufw8/Boh+Kr1Rc2p0dAfluQUi0FltnZY69HfkNpjzgVbynNw76HOY+spqDelQCLheYph/mG02DjILV/yIpPDlicMM/oJFus2Dkj1ug9/4XPX3zO1rGFXfFh6/QKOfuPVreni4b5COrfAV1I+xnKoyJJ+j145o9S/v9yI0=")
	header.Set(OssClientSideEncryptionWrapAlg, "RSA/NONE/PKCS1Padding")
	header.Set(OssClientSideEncryptionKey, "iIps8SuNOCoR5GK5JGuv3RIdPbCB99HY4bB3vqbN9E0AZHUUFXILNEU3P4rD5zWQFFtmLxJwns6pKAG1j40QYw78vsvV4MrSXPrpWVOHwMfPmWywpvs21rYU1e/9GC0ZuJufZRFf504fC1vYtAjbJLl8oqA9zQRvXzvrr9mTHqE=")
	env, err := getEnvelopeFromHeader(header)
	assert.Nil(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "AES/CTR/NoPadding", env.CEKAlg)
	assert.Equal(t, "{\"tag\":\"value\"}", env.MatDesc)
	assert.Equal(t, "RSA/NONE/PKCS1Padding", env.WrapAlg)
	assert.Equal(t, "RSA/NONE/PKCS1Padding", env.WrapAlg)
	assert.Equal(t, "", env.UnencryptedContentLen)
	assert.Equal(t, "", env.UnencryptedMD5)

	//Invalid Key
	header.Set(OssClientSideEncryptionKey, "***InvalidKey/9GC0ZuJufZRFf504fC1vYtAjbJLl8oqA9zQRvXzvrr9mTHqE=")
	env, err = getEnvelopeFromHeader(header)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "illegal base64 data")

	//Invalid Key
	header.Set(OssClientSideEncryptionStart, "***")
	header.Set(OssClientSideEncryptionKey, "iIps8SuNOCoR5GK5JGuv3RIdPbCB99HY4bB3vqbN9E0AZHUUFXILNEU3P4rD5zWQFFtmLxJwns6pKAG1j40QYw78vsvV4MrSXPrpWVOHwMfPmWywpvs21rYU1e/9GC0ZuJufZRFf504fC1vYtAjbJLl8oqA9zQRvXzvrr9mTHqE=")
	env, err = getEnvelopeFromHeader(header)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "illegal base64 data")
}

func TestMockEncryptionInvalidHeader(t *testing.T) {
	//length := 123
	data := []byte("hello world")
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &encryptionMockTracker{
		lastModified: gmtTime,
	}
	server := testSetupEncryptionMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	assert.NotNil(t, client)

	mc, err := crypto.CreateMasterRsa(map[string]string{"tag": "value"}, rsaPublicKey, rsaPrivateKey)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)

	result, err := eclient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   bytes.NewReader(data),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, datasum)
	assert.Len(t, tracker.savedata, len(data))
	assert.NotEqualValues(t, data, tracker.savedata)

	gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err := io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.EqualValues(t, data, gData)
	assert.NotEmpty(t, tracker.saveHeaders.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, tracker.saveHeaders.Get(OssClientSideEncryptionStart))
	assert.Equal(t, crypto.AesCtrAlgorithm, tracker.saveHeaders.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, crypto.RsaCryptoWrap, tracker.saveHeaders.Get(OssClientSideEncryptionWrapAlg))
	assert.Equal(t, "{\"tag\":\"value\"}", tracker.saveHeaders.Get(OssClientSideEncryptionMatDesc))
	assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionUnencryptedContentLength))
	assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionUnencryptedContentMD5))
	assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionDataSize))
	assert.Empty(t, tracker.saveHeaders.Get(OssClientSideEncryptionPartSize))

	result, err = eclient.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   bytes.NewReader(nil),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, tracker.savedata, 0)
	assert.NotEqualValues(t, data, tracker.savedata)

	gResult, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err = io.ReadAll(gResult.Body)
	assert.Nil(t, err)
	assert.Len(t, gData, 0)

	//Set Invaid Encryption-Start Header, shoud fail
	EncryptionStart := tracker.saveHeaders.Get(OssClientSideEncryptionStart)
	tracker.saveHeaders.Set(OssClientSideEncryptionStart, "***invalid")
	gResult, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "illegal base64 data")

	tracker.saveHeaders.Set(OssClientSideEncryptionStart, EncryptionStart)

	//Set Invaid Encryption-Start Header, shoud fail
	CekAlg := tracker.saveHeaders.Get(OssClientSideEncryptionCekAlg)
	tracker.saveHeaders.Set(OssClientSideEncryptionCekAlg, "unsupport-cek")
	gResult, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not supported content algorithm")

	tracker.saveHeaders.Set(OssClientSideEncryptionCekAlg, CekAlg)

	tracker.saveHeaders.Del(OssClientSideEncryptionWrapAlg)
	gResult, err = eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "getEnvelopeFromHeader error")
}

var rsaPrivateKeyCompatibility string = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCokfiAVXXf5ImFzKDw+XO/UByW6mse2QsIgz3ZwBtMNu59fR5z
ttSx+8fB7vR4CN3bTztrP9A6bjoN0FFnhlQ3vNJC5MFO1PByrE/MNd5AAfSVba93
I6sx8NSk5MzUCA4NJzAUqYOEWGtGBcom6kEF6MmR1EKib1Id8hpooY5xaQIDAQAB
AoGAOPUZgkNeEMinrw31U3b2JS5sepG6oDG2CKpPu8OtdZMaAkzEfVTJiVoJpP2Y
nPZiADhFW3e0ZAnak9BPsSsySRaSNmR465cG9tbqpXFKh9Rp/sCPo4Jq2n65yood
JBrnGr6/xhYvNa14sQ6xjjfSgRNBSXD1XXNF4kALwgZyCAECQQDV7t4bTx9FbEs5
36nAxPsPM6aACXaOkv6d9LXI7A0J8Zf42FeBV6RK0q7QG5iNNd1WJHSXIITUizVF
6aX5NnvFAkEAybeXNOwUvYtkgxF4s28s6gn11c5HZw4/a8vZm2tXXK/QfTQrJVXp
VwxmSr0FAajWAlcYN/fGkX1pWA041CKFVQJAG08ozzekeEpAuByTIOaEXgZr5MBQ
gBbHpgZNBl8Lsw9CJSQI15wGfv6yDiLXsH8FyC9TKs+d5Tv4Cvquk0efOQJAd9OC
lCKFs48hdyaiz9yEDsc57PdrvRFepVdj/gpGzD14mVerJbOiOF6aSV19ot27u4on
Td/3aifYs0CveHzFPQJAWb4LCDwqLctfzziG7/S7Z74gyq5qZF4FUElOAZkz718E
yZvADwuz/4aK0od0lX9c4Jp7Mo5vQ4TvdoBnPuGoyw==
-----END RSA PRIVATE KEY-----`

func TestMockEncryptionCompatibility(t *testing.T) {
	//length := 123
	gmtTime := getNowGMT()
	tracker := &encryptionMockTracker{
		lastModified: gmtTime,
	}
	server := testSetupEncryptionMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	assert.NotNil(t, client)

	mc, err := crypto.CreateMasterRsa(map[string]string{"tag": "value"}, "", rsaPrivateKeyCompatibility)
	assert.Nil(t, err)
	eclient, err := NewEncryptionClient(client, mc)
	assert.Nil(t, err)

	file, err := os.Open("../test/testdata/cpp-enc-example.jpg")
	assert.Nil(t, err)
	defer file.Close()
	result, err := client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Body:   file,
		Metadata: map[string]string{
			"client-side-encryption-key":      "nyXOp7delQ/MQLjKQMhHLaT0w7u2yQoDLkSnK8MFg/MwYdh4na4/LS8LLbLcM18m8I/ObWUHU775I50sJCpdv+f4e0jLeVRRiDFWe+uo7Puc9j4xHj8YB3QlcIOFQiTxHIB6q+C+RA6lGwqqYVa+n3aV5uWhygyv1MWmESurppg=",
			"client-side-encryption-start":    "De/S3T8wFjx7QPxAAFl7h7TeI2EsZlfCwox4WhLGng5DK2vNXxULmulMUUpYkdc9umqmDilgSy5Z3Foafw+v4JJThfw68T/9G2gxZLrQTbAlvFPFfPM9Ehk6cY4+8WpY32uN8w5vrHyoSZGr343NxCUGIp6fQ9sSuOLMoJg7hNw=",
			"client-side-encryption-cek-alg":  "AES/CTR/NoPadding",
			"client-side-encryption-wrap-alg": "RSA/NONE/PKCS1Padding",
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.NotEmpty(t, tracker.saveHeaders.Get(OssClientSideEncryptionKey))
	assert.NotEmpty(t, tracker.saveHeaders.Get(OssClientSideEncryptionStart))
	assert.Equal(t, "AES/CTR/NoPadding", tracker.saveHeaders.Get(OssClientSideEncryptionCekAlg))
	assert.Equal(t, "RSA/NONE/PKCS1Padding", tracker.saveHeaders.Get(OssClientSideEncryptionWrapAlg))

	gResult, err := eclient.GetObject(context.TODO(), &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, gResult)
	gData, err := io.ReadAll(gResult.Body)

	ghash := NewCRC64(0)
	ghash.Write(gData)

	file1, err := os.Open("../test/testdata/example.jpg")
	assert.Nil(t, err)
	defer file1.Close()
	fData, err := io.ReadAll(file1)

	fhash := NewCRC64(0)
	fhash.Write(fData)

	assert.Equal(t, fhash.Sum64(), ghash.Sum64())
}
