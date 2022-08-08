package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ostamand/aqualog/api"
)

const apiEndpoint = "http://localhost:8080" // TODO how to not hard code this and get at build time
const reqContentType = "application/json"
const apiTokenType = "Bearer"
const apiEndpointEnv = "AQUALOG_ENDPOINT"

var ErrUserNotFound = errors.New("user does not exists")
var ErrWrongPassword = errors.New("wrong password")
var ErrNeedToLogin = errors.New("new login required")
var ErrAPI = errors.New("issue with Aqualog API")

type aqualogAPI struct {
	endpoint string
	auth     api.LoginResponse
}

func NewAqualogAPI() aqualogAPI {
	// setup API endpoint
	var aqualog aqualogAPI

	endpoint, present := os.LookupEnv(apiEndpointEnv)
	if present {
		aqualog.endpoint = endpoint
	} else {
		aqualog.endpoint = apiEndpoint
	}

	return aqualog
}

// CreateParam saves a new parameter on Aqualog
func (aqualog *aqualogAPI) CreateParam(args api.CreateParamRequest) error {
	data, err := json.Marshal(args)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, aqualog.endpoint+"/params", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", reqContentType)
	req.Header.Add("Authorization", apiTokenType+" "+aqualog.auth.AccessToken)

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	code := httpResp.StatusCode
	if code != http.StatusOK {
		if code == http.StatusUnauthorized {
			return ErrNeedToLogin
		}
		return fmt.Errorf("api request error")
	}

	return nil
}

// LoadAuth will get access and refresh token from local storage.
func (aqualog *aqualogAPI) LoadAuth() error {
	path := GetAutPath()

	// check if auth file exists
	if _, err := os.Stat(path); err != nil {
		return err
	}

	// read auth file
	jsonFile, err := os.Open(filepath.Join(path, "auth.json"))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	data := api.LoginResponse{}
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return err
	}
	aqualog.auth = data

	return nil
}

// Login call Aqualog API to get access and renew tokens.
func (aqualog *aqualogAPI) Login(username string, password string) error {
	req := api.LoginRequest{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpResp, err := http.Post(aqualog.endpoint+"/login", reqContentType, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	code := httpResp.StatusCode
	if code != http.StatusOK {
		if code == http.StatusNotFound {
			return ErrUserNotFound
		} else if code == http.StatusUnauthorized {
			return ErrWrongPassword
		}
		// TODO add bad request password validation
		return fmt.Errorf("api request error")
	}

	// decode response body
	var resp api.LoginResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return err
	}

	aqualog.auth = resp
	// save login response to home folder
	SaveLoginResp(resp)
	return nil
}

// RenewTokenIf will renew the access token only if needed i.e. access token is expired.
func (aqualog *aqualogAPI) RenewTokenIf() error {
	if (aqualog.auth == api.LoginResponse{}) {
		return ErrNeedToLogin
	}
	// check if access token is expired
	if time.Now().After(aqualog.auth.AccessTokenExpiresAt) {
		// need to renew acess token
		if time.Now().Before(aqualog.auth.RenewTokenExpiresAt) {
			// can renew using the refresh token
			req := api.RenewTokenRequest{
				RenewToken: aqualog.auth.RenewToken,
			}
			data, err := json.Marshal(req)
			if err != nil {
				return ErrNeedToLogin
			}
			httpResp, err := http.Post(aqualog.endpoint+"/renew_token", reqContentType, bytes.NewBuffer(data))
			if err != nil {
				return err
			}
			code := httpResp.StatusCode
			if code != http.StatusOK {
				return ErrAPI
			}

			// decode response body
			var resp api.RenewTokenResponse
			err = json.NewDecoder(httpResp.Body).Decode(&resp)
			if err != nil {
				return err
			}

			aqualog.auth.AccessToken = resp.AccessToken
			aqualog.auth.AccessTokenExpiresAt = resp.AccessTokenExpiresAt

			// save with new access token
			SaveLoginResp(aqualog.auth)
		} else {
			return ErrNeedToLogin
		}
	}

	return nil
}
