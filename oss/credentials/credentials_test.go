package credentials

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCredentials(t *testing.T) {
	cred := &Credentials{}
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.False(t, cred.HasKeys())

	expires := time.Now().Add(10 * time.Second)
	cred = &Credentials{
		AccessKeyID:     "ak",
		AccessKeySecret: "sk",
		Expires:         &expires,
	}
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())

	expires = time.Now().Add(-10 * time.Second)
	cred = &Credentials{
		AccessKeyID:     "ak",
		AccessKeySecret: "sk",
		Expires:         &expires,
	}
	assert.NotNil(t, cred)
	assert.True(t, cred.Expired())
	assert.True(t, cred.HasKeys())
}

func TestStaticCredentialsProvider(t *testing.T) {
	provider := NewStaticCredentialsProvider("ak", "sk")
	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())
	assert.Equal(t, "ak", cred.AccessKeyID)
	assert.Equal(t, "sk", cred.AccessKeySecret)
	assert.Equal(t, "", cred.SessionToken)

	provider = NewStaticCredentialsProvider("ak1", "sk1", "token1")
	cred, err = provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())
	assert.Equal(t, "ak1", cred.AccessKeyID)
	assert.Equal(t, "sk1", cred.AccessKeySecret)
	assert.Equal(t, "token1", cred.SessionToken)
}

func TestEnvironmentVariableCredentialsProvider(t *testing.T) {
	provider := NewEnvironmentVariableCredentialsProvider()
	assert.NotNil(t, provider)

	oriak := os.Getenv("OSS_ACCESS_KEY_ID")
	orisk := os.Getenv("OSS_ACCESS_KEY_SECRET")
	oritk := os.Getenv("OSS_SESSION_TOKEN")

	defer func() {
		if oriak == "" {
			os.Clearenv()
		} else {
			os.Setenv("OSS_ACCESS_KEY_ID", oriak)
		}
		if orisk == "" {
			os.Clearenv()
		} else {
			os.Setenv("OSS_ACCESS_KEY_SECRET", orisk)
		}
		if oritk == "" {
			os.Clearenv()
		} else {
			os.Setenv("OSS_SESSION_TOKEN", oritk)
		}
	}()

	os.Setenv("OSS_ACCESS_KEY_ID", "myak")
	os.Setenv("OSS_ACCESS_KEY_SECRET", "mysk")
	provider = NewEnvironmentVariableCredentialsProvider()
	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())
	assert.Equal(t, "myak", cred.AccessKeyID)
	assert.Equal(t, "mysk", cred.AccessKeySecret)
	assert.Equal(t, "", cred.SessionToken)

	err = os.Setenv("OSS_SESSION_TOKEN", "mytoken")

	provider = NewEnvironmentVariableCredentialsProvider()
	cred, err = provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.Expired())
	assert.True(t, cred.HasKeys())
	assert.Equal(t, "myak", cred.AccessKeyID)
	assert.Equal(t, "mysk", cred.AccessKeySecret)
	assert.Equal(t, "mytoken", cred.SessionToken)
}

func TestAnonymousCredentialsProvider(t *testing.T) {
	provider := NewAnonymousCredentialsProvider()
	assert.NotNil(t, provider)

	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	assert.False(t, cred.HasKeys())
	assert.False(t, cred.Expired())
}
