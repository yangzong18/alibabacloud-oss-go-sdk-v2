package signer

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var requiredSignedParametersMap = map[string]struct{}{
	"acl":                          {},
	"bucketInfo":                   {},
	"location":                     {},
	"stat":                         {},
	"delete":                       {},
	"append":                       {},
	"tagging":                      {},
	"objectMeta":                   {},
	"uploads":                      {},
	"uploadId":                     {},
	"partNumber":                   {},
	"security-token":               {},
	"position":                     {},
	"response-content-type":        {},
	"response-content-language":    {},
	"response-expires":             {},
	"response-cache-control":       {},
	"response-content-disposition": {},
	"response-content-encoding":    {},
	"restore":                      {},
	"callback":                     {},
	"callback-var":                 {},
	"versions":                     {},
	"versioning":                   {},
	"versionId":                    {},
	"sequential":                   {},
	"continuation-token":           {},
	"regionList":                   {},
	"cloudboxes":                   {},
	"symlink":                      {},
}

const (
	// headers
	authorizationHeader = "Authorization"
	securityTokenHeader = "x-oss-security-token"

	dateHeader        = "Date"
	contentTypeHeader = "Content-Type"
	contentMd5Header  = "Content-MD5"
	ossHeaderPreifx   = "x-oss-"
	ossDateHeader     = "x-oss-date"

	//Query
	securityTokenQuery = "security-token"
	expiresQuery       = "Expires"
	accessKeyIdQuery   = "OSSAccessKeyId"
	signatureQuery     = "Signature"

	defaultExpiresDuration = 15 * time.Minute
)

type SignerV1 struct {
}

func isSubResource(list []string, key string) bool {
	for _, k := range list {
		if key == k {
			return true
		}
	}
	return false
}

func formatDate(time time.Time, unix bool) string {
	if unix {
		return fmt.Sprintf("%v", time.Unix())
	}
	return time.Format(http.TimeFormat)
}

func (*SignerV1) calcStringToSign(signingCtx *SigningContext) string {
	/*
		SignToString =
			VERB + "\n"
			+ Content-MD5 + "\n"
			+ Content-Type + "\n"
			+ Date + "\n"
			+ CanonicalizedOSSHeaders
			+ CanonicalizedResource
		Signature = base64(hmac-sha1(AccessKeySecret, SignToString))
	*/
	request := signingCtx.Request

	contentMd5 := request.Header.Get(contentMd5Header)
	contentType := request.Header.Get(contentTypeHeader)
	date := formatDate(signingCtx.Time, signingCtx.AuthMethodQuery)

	//CanonicalizedOSSHeaders
	var headers []string
	for k := range request.Header {
		lowerCaseKey := strings.ToLower(k)
		if strings.HasPrefix(lowerCaseKey, ossHeaderPreifx) {
			headers = append(headers, lowerCaseKey)
		}
	}
	sort.Strings(headers)
	headerItems := make([]string, len(headers))
	for i, k := range headers {
		headerValues := make([]string, len(request.Header.Values(k)))
		for i, v := range request.Header.Values(k) {
			headerValues[i] = strings.TrimSpace(v)
		}
		headerItems[i] = k + ":" + strings.Join(headerValues, ",") + "\n"
	}
	canonicalizedOSSHeaders := strings.Join(headerItems, "")

	//CanonicalizedResource
	query := request.URL.Query()
	var params []string
	for k := range query {
		if _, ok := requiredSignedParametersMap[k]; ok {
			params = append(params, k)
		} else if strings.HasPrefix(k, ossHeaderPreifx) {
			params = append(params, k)
		} else if isSubResource(signingCtx.SubResource, k) {
			params = append(params, k)
		}
	}
	sort.Strings(params)
	paramItems := make([]string, len(params))
	for i, k := range params {
		v := query.Get(k)
		if len(v) > 0 {
			paramItems[i] = k + "=" + v
		} else {
			paramItems[i] = k
		}
	}
	subResource := strings.Join(paramItems, "&")
	canonicalizedResource := "/"
	if signingCtx.Bucket != nil {
		canonicalizedResource += *signingCtx.Bucket + "/"
	}
	if signingCtx.Key != nil {
		canonicalizedResource += *signingCtx.Key
	}
	if subResource != "" {
		canonicalizedResource += "?" + subResource
	}

	// string to Sign
	stringToSign :=
		request.Method + "\n" +
			contentMd5 + "\n" +
			contentType + "\n" +
			date + "\n" +
			canonicalizedOSSHeaders +
			canonicalizedResource

	//fmt.Printf("stringToSign:%s\n", stringToSign)
	return stringToSign
}

func (s *SignerV1) authHeader(ctx context.Context, signingCtx *SigningContext) error {
	request := signingCtx.Request
	cred := signingCtx.Credentials

	// Date
	time := time.Now().UTC()
	if signingCtx.ClockOffset != 0 {
		time = time.Add(signingCtx.ClockOffset)
	}
	date := request.Header.Get(ossDateHeader)
	if len(date) > 0 {
		time, _ = http.ParseTime(date)
	}
	request.Header.Set(dateHeader, formatDate(time, false))
	signingCtx.Time = time

	// Credentials information
	if cred.SessionToken != "" {
		request.Header.Set(securityTokenHeader, cred.SessionToken)
	}

	// StringToSign
	stringToSign := s.calcStringToSign(signingCtx)
	signingCtx.StringToSign = stringToSign

	// Signature
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(cred.AccessKeySecret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Authorization header
	request.Header.Set(authorizationHeader, fmt.Sprintf("OSS %s:%s", cred.AccessKeyID, signature))

	return nil
}

func (s *SignerV1) authQuery(ctx context.Context, signingCtx *SigningContext) error {
	request := signingCtx.Request
	cred := signingCtx.Credentials

	// Date
	if signingCtx.Time.IsZero() {
		signingCtx.Time = time.Now().Add(defaultExpiresDuration)
	}

	// Credentials information
	query, _ := url.ParseQuery(request.URL.RawQuery)
	if cred.SessionToken != "" {
		query.Add(securityTokenQuery, cred.SessionToken)
		request.URL.RawQuery = query.Encode()
	}

	// StringToSign
	stringToSign := s.calcStringToSign(signingCtx)
	signingCtx.StringToSign = stringToSign

	// Signature
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(cred.AccessKeySecret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Authorization query
	query.Add(expiresQuery, fmt.Sprintf("%v", signingCtx.Time.Unix()))
	query.Add(accessKeyIdQuery, cred.AccessKeyID)
	query.Add(signatureQuery, signature)
	request.URL.RawQuery = query.Encode()

	return nil
}

func (s *SignerV1) Sign(ctx context.Context, signingCtx *SigningContext) error {
	if signingCtx == nil {
		return fmt.Errorf("SigningContext is null.")
	}

	if signingCtx.Credentials == nil || !signingCtx.Credentials.HasKeys() {
		return fmt.Errorf("SigningContext.Credentials is null or empty.")
	}

	if signingCtx.Request == nil {
		return fmt.Errorf("SigningContext.Request is null.")
	}

	if signingCtx.AuthMethodQuery {
		return s.authQuery(ctx, signingCtx)
	}

	return s.authHeader(ctx, signingCtx)
}

func (*SignerV1) IsSignedHeader(h string) bool {
	lowerCaseKey := strings.ToLower(h)
	if strings.HasPrefix(lowerCaseKey, ossHeaderPreifx) ||
		lowerCaseKey == "date" ||
		lowerCaseKey == "content-type" ||
		lowerCaseKey == "content-md5" {
		return true
	}
	return false
}
