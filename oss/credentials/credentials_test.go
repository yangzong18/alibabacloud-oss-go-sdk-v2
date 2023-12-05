package credentials

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

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

type stubCredentialsFetcher struct {
	delay time.Duration
	token string
	count int64
}

func (s *stubCredentialsFetcher) Fetch(ctx context.Context) (Credentials, error) {
	var expires *time.Time
	if s.delay > 0 {
		now := time.Now()
		new := now.Add(s.delay)
		expires = &new
	}

	s.count++

	return Credentials{
		AccessKeyID:     "ak",
		AccessKeySecret: "sk",
		SessionToken:    s.token,
		Expires:         expires,
	}, nil
}

func TestCredentialsFetcherProvider(t *testing.T) {
	provider := NewCredentialsFetcherProvider(nil)
	assert.NotNil(t, provider)
	fetcherProvider, ok := provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Equal(t, defaultExpiredFactor, fetcherProvider.expiredFactor)
	assert.Equal(t, defaultRefreshDuration, fetcherProvider.refreshDuration)
	assert.Nil(t, fetcherProvider.nextRefreshTime)
	assert.Nil(t, fetcherProvider.fetcher)
	//assert.Nil(t, fetcherProvider.credentials)

	_, err := provider.GetCredentials(context.TODO())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "fetcher is null.")

	provider = NewCredentialsFetcherProvider(nil, func(o *CredentialsFetcherOptions) {
		o.ExpiredFactor = 0.7
		o.RefreshDuration = 1 * time.Second
	})
	assert.NotNil(t, provider)
	fetcherProvider, ok = provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Equal(t, 0.7, fetcherProvider.expiredFactor)
	assert.Equal(t, 1*time.Second, fetcherProvider.refreshDuration)

	provider = NewCredentialsFetcherProvider(&stubCredentialsFetcher{})
	assert.NotNil(t, provider)
	cred, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred.AccessKeyID)
	assert.Equal(t, "sk", cred.AccessKeySecret)
	assert.False(t, cred.Expired())

	// with Expired
	provider = NewCredentialsFetcherProvider(&stubCredentialsFetcher{
		token: "token",
		delay: 2 * time.Second,
	})
	assert.NotNil(t, provider)
	fetcherProvider, ok = provider.(*CredentialsFetcherProvider)
	assert.NotNil(t, fetcherProvider)
	assert.Nil(t, fetcherProvider.nextRefreshTime)

	// 1st
	cred1, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred1.AccessKeyID)
	assert.Equal(t, "sk", cred1.AccessKeySecret)
	assert.Equal(t, "token", cred1.SessionToken)
	assert.NotNil(t, cred1.Expires)
	assert.False(t, cred1.Expired())
	assert.Nil(t, fetcherProvider.nextRefreshTime)

	// 2st
	cred2, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred2.AccessKeyID)
	assert.Equal(t, "sk", cred2.AccessKeySecret)
	assert.Equal(t, "token", cred2.SessionToken)
	assert.Equal(t, cred1.Expires, cred2.Expires)

	time.Sleep(3 * time.Second)
	assert.True(t, cred1.Expired())

	// 3st
	cred3, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred3.AccessKeyID)
	assert.Equal(t, "sk", cred3.AccessKeySecret)
	assert.Equal(t, "token", cred3.SessionToken)
	assert.False(t, cred3.Expired())

	assert.True(t, cred3.Expires.After(*cred1.Expires))

}

func TestCredentialsFetcherProvider_Soon(t *testing.T) {
	// with Expired
	provider := NewCredentialsFetcherProvider(
		&stubCredentialsFetcher{
			token: "token",
			delay: 10 * time.Second,
		},
		func(o *CredentialsFetcherOptions) {
			o.ExpiredFactor = 0.4
			o.RefreshDuration = 1 * time.Second
		},
	)
	assert.NotNil(t, provider)
	fetcherProvider, ok := provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Nil(t, fetcherProvider.nextRefreshTime)

	// 1st
	cred1, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred1.AccessKeyID)
	assert.Equal(t, "sk", cred1.AccessKeySecret)
	assert.Equal(t, "token", cred1.SessionToken)
	assert.NotNil(t, cred1.Expires)
	assert.False(t, cred1.Expired())
	assert.NotNil(t, fetcherProvider.nextRefreshTime)

	// 2st
	time.Sleep(6 * time.Second)
	assert.False(t, cred1.Expired())
	cred2, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred2.AccessKeyID)
	assert.Equal(t, "sk", cred2.AccessKeySecret)
	assert.Equal(t, "token", cred2.SessionToken)
	assert.True(t, cred2.Expires.After(*cred1.Expires))
}

func TestCredentialsFetcherProvider_MultiJobs(t *testing.T) {
	// with Expired
	fetcher := &stubCredentialsFetcher{
		token: "token",
		delay: 10 * time.Second,
	}

	provider := NewCredentialsFetcherProvider(
		fetcher,
		func(o *CredentialsFetcherOptions) {
			o.ExpiredFactor = 0.4
			o.RefreshDuration = 1 * time.Second
		},
	)
	assert.NotNil(t, provider)
	fetcherProvider, ok := provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Nil(t, fetcherProvider.nextRefreshTime)

	run := true
	testFn := func() {
		count := int64(0)
		for run {
			cred, err := provider.GetCredentials(context.TODO())
			assert.Nil(t, err)
			assert.Equal(t, "ak", cred.AccessKeyID)
			assert.Equal(t, "sk", cred.AccessKeySecret)
			assert.Equal(t, "token", cred.SessionToken)
			assert.NotNil(t, cred.Expires)
			assert.False(t, cred.Expired())
			count++
		}
		assert.Greater(t, count, int64(1000))
	}

	for i := 0; i < 20; i++ {
		go testFn()
	}

	time.Sleep(1 * time.Second)
	assert.NotNil(t, fetcherProvider.nextRefreshTime)

	time.Sleep(8 * time.Second)
	run = false
	assert.Less(t, fetcher.count, int64(4))
}

type stubCredentialsFetcher2 struct {
	delay        time.Duration
	token        string
	returnErr    bool
	returnTimout bool
}

func (s *stubCredentialsFetcher2) Fetch(ctx context.Context) (Credentials, error) {
	var expires *time.Time
	if s.delay > 0 {
		now := time.Now()
		new := now.Add(s.delay)
		expires = &new
	}

	if s.returnTimout {
		time.Sleep(10 * time.Second)
		return Credentials{}, fmt.Errorf("returnTimout")
	} else if s.returnErr {
		return Credentials{}, fmt.Errorf("returnErr")
	} else {
		return Credentials{
			AccessKeyID:     "ak",
			AccessKeySecret: "sk",
			SessionToken:    s.token,
			Expires:         expires,
		}, nil
	}
}

func TestCredentialsFetcherProvider_Error(t *testing.T) {
	fetcher := &stubCredentialsFetcher2{
		token:     "token",
		delay:     10 * time.Second,
		returnErr: true,
	}

	provider := NewCredentialsFetcherProvider(
		fetcher,
		func(o *CredentialsFetcherOptions) {
			o.ExpiredFactor = 0.4
			o.RefreshDuration = 1 * time.Second
		},
	)
	assert.NotNil(t, provider)
	fetcherProvider, ok := provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	assert.Nil(t, fetcherProvider.nextRefreshTime)

	// Get Fail
	_, err := provider.GetCredentials(context.TODO())
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "returnErr")

	// Get OK
	fetcher.returnErr = false
	cred1, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred1.AccessKeyID)
	assert.Equal(t, "sk", cred1.AccessKeySecret)
	assert.Equal(t, "token", cred1.SessionToken)
	assert.NotNil(t, cred1.Expires)
	assert.False(t, cred1.Expired())

	// 2st Fail
	fetcher.returnErr = true
	time.Sleep(6 * time.Second)
	assert.False(t, cred1.Expired())
	cred2, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred2.AccessKeyID)
	assert.Equal(t, "sk", cred2.AccessKeySecret)
	assert.Equal(t, "token", cred2.SessionToken)
	assert.Equal(t, *cred1.Expires, *cred2.Expires)
	assert.True(t, fetcherProvider.nextRefreshTime.After(time.Now()))

	// Fetch Timeout
	fetcher = &stubCredentialsFetcher2{
		token: "token",
		delay: 6 * time.Second,
	}

	provider = NewCredentialsFetcherProvider(
		fetcher,
		func(o *CredentialsFetcherOptions) {
			o.ExpiredFactor = 0.4
			o.RefreshDuration = 1 * time.Second
		},
	)
	assert.NotNil(t, provider)
	fetcherProvider, ok = provider.(*CredentialsFetcherProvider)
	assert.True(t, ok)
	assert.NotNil(t, fetcherProvider)
	ctxt1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel1()
	fetcher.returnTimout = true
	_, err = provider.GetCredentials(ctxt1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "FetchCredentialsCanceled")

	fetcher.returnTimout = false
	cred3, err := provider.GetCredentials(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred3.AccessKeyID)
	assert.Equal(t, "sk", cred3.AccessKeySecret)
	assert.Equal(t, "token", cred3.SessionToken)
	assert.NotNil(t, cred3.Expires)
	assert.False(t, cred3.Expired())

	time.Sleep(4 * time.Second)
	ctxt2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()
	fetcher.returnTimout = true
	cred4, err := provider.GetCredentials(ctxt2)
	assert.Nil(t, err)
	assert.Equal(t, "ak", cred4.AccessKeyID)
	assert.Equal(t, "sk", cred4.AccessKeySecret)
	assert.Equal(t, "token", cred4.SessionToken)
	assert.NotNil(t, cred4.Expires)
	assert.Equal(t, *cred3.Expires, *cred4.Expires)
}
