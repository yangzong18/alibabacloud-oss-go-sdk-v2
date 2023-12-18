package signer

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var noEscape [256]bool

const (
	// headers
	contentSha256Header   = "x-oss-content-sha256"
	iso8601DatetimeFormat = "20060102T150405Z"
	iso8601DateFormat     = "20060102"
	terminator            = "aliyun_v4_request"
	secretKeyPrefix       = "aliyun_v4"
	signingAlgorithmV4    = "OSS4-HMAC-SHA256"

	ossSecurityTokenQuery    = "x-oss-security-token"
	ossCredentialQuery       = "x-oss-credential"
	ossExpiresQuery          = "x-oss-expires"
	ossSignatureQuery        = "x-oss-signature"
	ossSignatureVersionQuery = "x-oss-signature-version"
	ossAdditionalQuery       = "x-oss-additional-headers"
)

type PayloadType string

const (
	PayloadUnsignedType PayloadType = "UNSIGNED-PAYLOAD"
)

type SignerV4 struct {
}

func (s *SignerV4) calcStringToSign(signingCtx *SigningContext) string {
	/**
	StringToSign
	"OSS4-HMAC-SHA256" + "\n" +
	TimeStamp + "\n" +
	Scope + "\n" +
	Hex(SHA256Hash(Canonical Request))
	*/
	canonicalRequest := s.calcCanonicalRequest(signingCtx)
	//fmt.Printf("canonicalReuqest:%s\n", canonicalRequest)
	hash256 := sha256.New()
	hash256.Write([]byte(canonicalRequest))
	hashValue := hash256.Sum(nil)
	hexCanonicalRequest := hex.EncodeToString(hashValue)
	signDate := ""
	if signingCtx.AuthMethodQuery {
		signDate = signingCtx.currentTime.Format(iso8601DatetimeFormat)
	} else {
		signDate = signingCtx.Time.Format(http.TimeFormat)
	}
	return signingAlgorithmV4 + "\n" + signDate + "\n" + s.buildScope(signingCtx) + "\n" + hexCanonicalRequest
}

func (s *SignerV4) calcCanonicalRequest(signingCtx *SigningContext) string {
	request := signingCtx.Request
	//Canonical Uri
	uri := "/"
	if signingCtx.Bucket != nil {
		uri += *signingCtx.Bucket + "/"
	}
	if signingCtx.Key != nil {
		uri += *signingCtx.Key
	}
	canonicalUri := escapePath(uri, false)
	//Canonical Query
	canonicalQuery := s.getCanonicalQuery(signingCtx)
	additionalList, additionalMap := s.getAdditionalHeaderKeys(signingCtx)
	additionalOSSHeaders := strings.Join(additionalList, ";")
	//Canonical OSS Headers
	var headers []string
	for k := range request.Header {
		lowerCaseKey := strings.ToLower(k)
		if strings.HasPrefix(lowerCaseKey, ossHeaderPreifx) || lowerCaseKey == "content-type" || lowerCaseKey == "content-md5" {
			headers = append(headers, lowerCaseKey)
		} else if len(additionalList) > 0 {
			if _, ok := additionalMap[lowerCaseKey]; ok {
				headers = append(headers, lowerCaseKey)
			}
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
	canonicalOSSHeaders := strings.Join(headerItems, "")
	hashPayload := string(PayloadUnsignedType)
	if val := request.Header.Get(contentSha256Header); val != "" {
		hashPayload = val
	}
	/*
		Canonical Request
		HTTP Verb + "\n" +
		Canonical URI + "\n" +
		Canonical Query String + "\n" +
		Canonical Headers + "\n" +
		Additional Headers + "\n" +
		Hashed PayLoad
	*/
	return request.Method + "\n" +
		canonicalUri + "\n" +
		canonicalQuery + "\n" +
		canonicalOSSHeaders + "\n" +
		additionalOSSHeaders + "\n" +
		hashPayload
}

func (s *SignerV4) authHeader(ctx context.Context, signingCtx *SigningContext) error {
	request := signingCtx.Request
	cred := signingCtx.Credentials
	// Date
	timeUtc := time.Now().UTC()
	if signingCtx.ClockOffset != 0 {
		timeUtc = timeUtc.Add(signingCtx.ClockOffset)
	}
	date := request.Header.Get(ossDateHeader)
	if len(date) > 0 {
		timeUtc, _ = http.ParseTime(date)
	}
	request.Header.Set(dateHeader, formatDate(timeUtc, false))
	signingCtx.Time = timeUtc

	// Credentials information
	if cred.SessionToken != "" {
		request.Header.Set(securityTokenHeader, cred.SessionToken)
	}
	contentSha256 := request.Header.Get(contentSha256Header)
	if len(contentSha256) == 0 {
		request.Header.Set(contentSha256Header, string(PayloadUnsignedType))
	}
	// StringToSign
	stringToSign := s.calcStringToSign(signingCtx)
	//fmt.Printf("stringToSign:%s\n", stringToSign)
	signingCtx.StringToSign = stringToSign
	// credential
	credential := fmt.Sprintf("%s/%s", cred.AccessKeyID, s.buildScope(signingCtx))
	signature := s.generateSignature(signingCtx)
	additionalList, _ := s.getAdditionalHeaderKeys(signingCtx)
	// Authorization header
	if len(additionalList) > 0 {
		additionalOSSHeaders := strings.Join(additionalList, ";")
		request.Header.Set(authorizationHeader, fmt.Sprintf(signingAlgorithmV4+" Credential=%s,AdditionalHeaders=%s,Signature=%s", credential, additionalOSSHeaders, signature))
	} else {
		request.Header.Set(authorizationHeader, fmt.Sprintf(signingAlgorithmV4+" Credential=%s,Signature=%s", credential, signature))
	}
	return nil
}

func (s *SignerV4) authQuery(ctx context.Context, signingCtx *SigningContext) error {
	request := signingCtx.Request
	cred := signingCtx.Credentials

	// Date
	if signingCtx.Time.IsZero() {
		signingCtx.Time = time.Now().Add(defaultExpiresDuration)
	}
	// Credentials information
	query, _ := url.ParseQuery(request.URL.RawQuery)
	if cred.SessionToken != "" {
		query.Add(ossSecurityTokenQuery, cred.SessionToken)
	}
	query.Add(ossSignatureVersionQuery, signingAlgorithmV4)
	query.Add(ossDateHeader, signingCtx.currentTime.Format(iso8601DatetimeFormat))
	diff := signingCtx.Time.Unix() - signingCtx.currentTime.Unix()
	query.Add(ossExpiresQuery, fmt.Sprintf("%v", diff))
	query.Add(ossCredentialQuery, fmt.Sprintf("%s/%s", cred.AccessKeyID, s.buildScope(signingCtx)))
	additionalList, _ := s.getAdditionalHeaderKeys(signingCtx)
	if len(additionalList) > 0 {
		additionalOSSHeaders := strings.Join(additionalList, ";")
		query.Add(ossAdditionalQuery, additionalOSSHeaders)
	}
	request.URL.RawQuery = query.Encode()
	// StringToSign
	stringToSign := s.calcStringToSign(signingCtx)
	//fmt.Printf("stringToSign:%s\n", stringToSign)
	signingCtx.StringToSign = stringToSign
	signature := s.generateSignature(signingCtx)
	query.Add(ossSignatureQuery, signature)
	request.URL.RawQuery = query.Encode()
	request.URL.RawQuery = s.getCanonicalQuery(signingCtx)
	return nil
}

func (s *SignerV4) Sign(ctx context.Context, signingCtx *SigningContext) error {
	if signingCtx == nil {
		return fmt.Errorf("SigningContext is null.")
	}

	if signingCtx.Credentials == nil || !signingCtx.Credentials.HasKeys() {
		return fmt.Errorf("SigningContext.Credentials is null or empty.")
	}

	if signingCtx.Request == nil {
		return fmt.Errorf("SigningContext.Request is null.")
	}
	signingCtx.currentTime = time.Now().UTC()
	if signingCtx.AuthMethodQuery {
		return s.authQuery(ctx, signingCtx)
	}
	return s.authHeader(ctx, signingCtx)
}

func (s *SignerV4) buildScope(signingCtx *SigningContext) string {
	return fmt.Sprintf("%s/%s/%s/%s", signingCtx.currentTime.Format(iso8601DateFormat), *signingCtx.Region, *signingCtx.Product, terminator)
}

func (s *SignerV4) getAdditionalHeaderKeys(signingCtx *SigningContext) ([]string, map[string]string) {
	request := signingCtx.Request
	var keysList []string
	keysMap := make(map[string]string)
	srcKeys := make(map[string]string)

	for k := range request.Header {
		srcKeys[strings.ToLower(k)] = ""
	}
	for _, v := range signingCtx.AdditionalHeaders {
		if _, ok := srcKeys[strings.ToLower(v)]; ok {
			if !strings.EqualFold(v, "content-type") && !strings.EqualFold(v, "content-md5") && !strings.HasPrefix(strings.ToLower(v), ossHeaderPreifx) {
				keysMap[strings.ToLower(v)] = ""
			}
		}
	}
	for k := range keysMap {
		keysList = append(keysList, k)
	}
	sort.Strings(keysList)
	return keysList, keysMap
}

func (*SignerV4) IsSignedHeader(h string) bool {
	lowerCaseKey := strings.ToLower(h)
	if strings.HasPrefix(lowerCaseKey, ossHeaderPreifx) || lowerCaseKey == "content-type" || lowerCaseKey == "content-md5" {
		return true
	}
	return false
}

func (s *SignerV4) generateSignature(signingCtx *SigningContext) string {
	hmacHash := func() hash.Hash { return sha256.New() }
	cred := signingCtx.Credentials
	stringToSign := signingCtx.StringToSign
	h1 := hmac.New(func() hash.Hash { return sha256.New() }, []byte(secretKeyPrefix+cred.AccessKeySecret))
	io.WriteString(h1, signingCtx.Time.Format(iso8601DateFormat))
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	io.WriteString(h2, *signingCtx.Region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	io.WriteString(h3, *signingCtx.Product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	io.WriteString(h4, terminator)
	h4Key := h4.Sum(nil)

	h := hmac.New(hmacHash, h4Key)
	io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}

func (s *SignerV4) getCanonicalQuery(signingCtx *SigningContext) string {
	request := signingCtx.Request
	query := request.URL.Query()
	var params []string
	for k := range query {
		params = append(params, k)
	}
	sort.Strings(params)
	paramItems := make([]string, len(params))
	for i, k := range params {
		v := query.Get(k)
		if len(v) > 0 {
			paramItems[i] = url.QueryEscape(k) + "=" + url.QueryEscape(v)
		} else {
			paramItems[i] = url.QueryEscape(k)
		}
	}
	return strings.Join(paramItems, "&")
}

func escapePath(path string, encodeSep bool) string {
	var buf bytes.Buffer
	for i := 0; i < len(path); i++ {
		c := path[i]
		if noEscape[c] || (c == '/' && !encodeSep) {
			buf.WriteByte(c)
		} else {
			fmt.Fprintf(&buf, "%%%02X", c)
		}
	}
	return buf.String()
}
