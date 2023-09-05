package credentials

import (
	"context"
)

type Credentials struct {
	AccessKeyID     string
	AccessKeySecret string
	SessionToken    string
}

type CredentialsProvider interface {
	GetCredentials(ctx context.Context) (Credentials, error)
}

type AnonymousCredentialsProvider struct{}

func (AnonymousCredentialsProvider) GetCredentials(ctx context.Context) (Credentials, error) {
	return Credentials{AccessKeyID: "", AccessKeySecret: ""}, nil
}
