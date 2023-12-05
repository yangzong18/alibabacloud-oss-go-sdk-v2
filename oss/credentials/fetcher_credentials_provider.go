package credentials

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// Default expiration time adjustment factor
	defaultExpiredFactor = 0.8

	// backoff of refresh time
	defaultRefreshDuration = 120 * time.Second
)

type CredentialsFetcherOptions struct {
	ExpiredFactor   float64
	RefreshDuration time.Duration
}

type CredentialsFetcher interface {
	Fetch(ctx context.Context) (Credentials, error)
}

type CredentialsFetcherProvider struct {
	m sync.Mutex

	//credentials *Credentials
	credentials atomic.Value

	fetcher CredentialsFetcher

	expiredFactor   float64
	refreshDuration time.Duration
	nextRefreshTime *time.Time
}

func NewCredentialsFetcherProvider(fetcher CredentialsFetcher, optFns ...func(*CredentialsFetcherOptions)) CredentialsProvider {
	options := CredentialsFetcherOptions{
		ExpiredFactor:   defaultExpiredFactor,
		RefreshDuration: defaultRefreshDuration,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &CredentialsFetcherProvider{
		fetcher:         fetcher,
		expiredFactor:   options.ExpiredFactor,
		refreshDuration: options.RefreshDuration,
	}
}

func (c *CredentialsFetcherProvider) GetCredentials(ctx context.Context) (Credentials, error) {
	var curCreds *Credentials
	if v := c.credentials.Load(); v != nil {
		curCreds, _ = v.(*Credentials)
	}
	if c.isExpired(curCreds) {
		c.m.Lock()
		defer c.m.Unlock()
		creds, err := c.fetch(ctx)
		if err == nil {
			c.update(&creds)
		}
		return creds, err
	} else {
		if c.isSoonExpire(curCreds) && c.m.TryLock() {
			defer c.m.Unlock()
			curCreds1 := c.credentials.Load().(*Credentials)
			if curCreds1 != curCreds {
				curCreds = curCreds1
			} else {
				creds, err := c.fetch(ctx)
				if err == nil {
					c.update(&creds)
					curCreds = &creds
				} else {
					c.updateNextRefreshTime()
					err = nil
				}
			}
		}
		return *curCreds, nil
	}
}

type asyncFetchResult struct {
	val Credentials
	err error
}

func (c *CredentialsFetcherProvider) asyncFetch(ctx context.Context) <-chan asyncFetchResult {
	doChan := func() <-chan asyncFetchResult {
		ch := make(chan asyncFetchResult, 1)

		go func() {
			cred, err := c.fetcher.Fetch(ctx)
			ch <- asyncFetchResult{cred, err}
		}()

		return ch
	}

	return doChan()
}

func (c *CredentialsFetcherProvider) fetch(ctx context.Context) (Credentials, error) {
	if c.fetcher == nil {
		return Credentials{}, fmt.Errorf("fetcher is null.")
	}

	select {
	case result, _ := <-c.asyncFetch(ctx):
		return result.val, result.err
	case <-ctx.Done():
		return Credentials{}, fmt.Errorf("FetchCredentialsCanceled")
	}
}

func (c *CredentialsFetcherProvider) update(cred *Credentials) {
	c.credentials.Store(cred)
	c.nextRefreshTime = nil
	if cred.Expires != nil {
		curr := time.Now().Round(0)
		durationS := c.expiredFactor * float64(cred.Expires.Sub(curr).Seconds())
		duration := time.Duration(durationS * float64(time.Second))
		if duration > c.refreshDuration {
			curr = curr.Add(duration)
			c.nextRefreshTime = &curr
		}
	}
}

func (c *CredentialsFetcherProvider) updateNextRefreshTime() {
	if c.nextRefreshTime != nil {
		//next := c.nextRefreshTime.Add(c.refreshDuration)
		next := time.Now().Round(0).Add(c.refreshDuration)
		c.nextRefreshTime = &next
	}
}

func (c *CredentialsFetcherProvider) isExpired(creds *Credentials) bool {
	return creds == nil || creds.Expired()
}

func (c *CredentialsFetcherProvider) isSoonExpire(creds *Credentials) bool {
	if creds == nil || creds.Expired() {
		return true
	}

	if c.nextRefreshTime != nil && !c.nextRefreshTime.After(time.Now().Round(0)) {
		return true
	}

	return false
}
