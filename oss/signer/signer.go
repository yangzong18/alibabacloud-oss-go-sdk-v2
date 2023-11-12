package signer

import (
	"context"
	"net/http"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
)

const (
	SubResource = "SubResource"
	SignTime    = "SignTime"
)

type SigningContext struct {
	//input
	Product *string
	Region  *string
	Bucket  *string
	Key     *string
	Request *http.Request

	SubResource []string

	Credentials *credentials.Credentials

	AuthMethodQuery bool

	// input and output
	Time time.Time

	// output
	SignedHeaders map[string]string
	StringToSign  string
}

type Signer interface {
	Sign(ctx context.Context, signingCtx *SigningContext) error
}

type NopSigner struct{}

func (*NopSigner) Sign(ctx context.Context, signingCtx *SigningContext) error {
	return nil
}
