package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

const ecs_ram_cred_url = "http://100.100.100.200/latest/meta-data/ram/security-credentials/"

type ecsRoleCredentialsProvider struct {
	ramCredUrl string
	ramRole    string
	timeout    time.Duration
}

type ecsRoleCredentials struct {
	AccessKeyId     string    `json:"AccessKeyId,omitempty"`
	AccessKeySecret string    `json:"AccessKeySecret,omitempty"`
	SecurityToken   string    `json:"SecurityToken,omitempty"`
	Expiration      time.Time `json:"Expiration,omitempty"`
	LastUpDated     time.Time `json:"LastUpDated,omitempty"`
	Code            string    `json:"Code,omitempty"`
}

func (p *ecsRoleCredentialsProvider) getRoleFromMetaData(ctx context.Context) (string, error) {
	c := &http.Client{
		Timeout: p.timeout,
	}

	resp, err := c.Get(p.ramCredUrl)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch ecs role name, resp.StatusCode:%v", resp.StatusCode)
	}
	defer resp.Body.Close()
	roleName, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(roleName) == 0 {
		return "", errors.New("ecs role name is empty")
	}

	return string(roleName), nil
}

func (p *ecsRoleCredentialsProvider) getCredentialsFromMetaData(ctx context.Context) (ecsRoleCredentials, error) {
	var ecsCred ecsRoleCredentials
	c := &http.Client{
		Timeout: p.timeout,
	}
	url := path.Join(p.ramCredUrl, p.ramRole)
	resp, err := c.Get(url)
	if err != nil {
		return ecsCred, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ecsCred, err
	}
	err = json.Unmarshal(body, &ecsCred)
	if err != nil {
		return ecsCred, err
	}

	if ecsCred.Code != "" && strings.ToUpper(ecsCred.Code) != "SUCCESS" {
		return ecsCred, fmt.Errorf("failed to fetch credentials, return code:%s", ecsCred.Code)
	}

	if ecsCred.AccessKeyId == "" || ecsCred.AccessKeySecret == "" {
		return ecsCred, fmt.Errorf("AccessKeyId or AccessKeySecret is empty, response body is '%s'", string(body))
	}

	return ecsCred, nil
}

func (p *ecsRoleCredentialsProvider) GetCredentials(ctx context.Context) (cred Credentials, err error) {
	if len(p.ramRole) == 0 {
		if name, err := p.getRoleFromMetaData(ctx); err != nil {
			return cred, err
		} else {
			p.ramRole = name
		}
	}

	ecsCred, err := p.getCredentialsFromMetaData(ctx)
	if err != nil {
		return cred, err
	}

	cred.AccessKeyID = ecsCred.AccessKeyId
	cred.AccessKeySecret = ecsCred.AccessKeySecret
	cred.SecurityToken = ecsCred.SecurityToken
	if !ecsCred.Expiration.IsZero() {
		cred.Expires = &ecsCred.Expiration
	}

	return cred, nil
}

func NewEcsRoleCredentialsProviderWithoutRefresh(roleNmae string) CredentialsProvider {
	return &ecsRoleCredentialsProvider{
		ramCredUrl: ecs_ram_cred_url,
		ramRole:    roleNmae,
		timeout:    15 * time.Second,
	}
}

func NewEcsRoleCredentialsProvider(roleNmae string) CredentialsProvider {
	p := NewEcsRoleCredentialsProviderWithoutRefresh(roleNmae)
	provider := NewCredentialsFetcherProvider(CredentialsFetcherFunc(func(ctx context.Context) (Credentials, error) {
		return p.GetCredentials(ctx)
	}))
	return provider
}
