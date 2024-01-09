package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type EcsCredentialsFetcher struct {
	Count    int64
	RoleName string
}

type EcsRoleCredentials struct {
	AccessKeyId     string `json:"AccessKeyId,omitempty"`
	AccessKeySecret string `json:"AccessKeySecret,omitempty"`
	SecurityToken   string `json:"SecurityToken,omitempty"`
	Expiration      string `json:"Expiration,omitempty"`
	LastUpDated     string `json:"LastUpDated,omitempty"`
	Code            string `json:"Code,omitempty"`
}

func (s *EcsCredentialsFetcher) fetchEcsCredentials() (EcsRoleCredentials, error) {
	var ecsCred EcsRoleCredentials
	var url string
	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	ecsUrlMetaData := "http://100.100.100.200/latest/meta-data/ram/security-credentials/"
	if s.RoleName != "" {
		url = ecsUrlMetaData + s.RoleName
	}else{
		resp, err := c.Get(ecsUrlMetaData)
		if err != nil {
			return ecsCred, err
		}
		if resp.StatusCode != http.StatusOK {
			return ecsCred, errors.New("failed to fetch ecs role name")
		}
		defer resp.Body.Close()
		roleName, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ecsCred, err
		}
		url = ecsUrlMetaData + string(roleName)
	}
	resp, err := c.Get(url)
	if err != nil {
		return ecsCred, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ecsCred, err
	}
	err = json.Unmarshal(body, &ecsCred)
	if err != nil {
		return ecsCred, err
	}
	if ecsCred.Code != "" && strings.ToUpper(ecsCred.Code) != "SUCCESS" {
		return ecsCred, fmt.Errorf("get sts ak error,code:%s", ecsCred.Code)
	}

	if ecsCred.AccessKeyId == "" || ecsCred.AccessKeySecret == "" {
		return ecsCred, fmt.Errorf("parsar http json body error:\n%s\n", string(body))
	}
	return ecsCred, nil
}

func (s *EcsCredentialsFetcher) Fetch(ctx context.Context) (Credentials, error) {
	credentials, err := s.fetchEcsCredentials()
	if err != nil {
		return Credentials{}, err
	}
	expires, err := time.Parse(time.RFC3339, credentials.Expiration)
	if err != nil {
		return Credentials{}, err
	}
	s.Count++

	return Credentials{
		AccessKeyID:     credentials.AccessKeyId,
		AccessKeySecret: credentials.AccessKeySecret,
		SessionToken:    credentials.SecurityToken,
		Expires:         &expires,
	}, nil
}


