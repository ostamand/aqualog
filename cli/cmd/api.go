package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ostamand/aqualog/api"
)

const apiAddress = "http://localhost:8080" // TODO how to not hard code this and get at build time

var ErrUserNotFound = errors.New("user does not exists")
var ErrWrongPassword = errors.New("wrong password")

func getAutPath() string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".aqualog")
	os.MkdirAll(path, os.ModePerm)
	return path
}

func saveLoginResp(resp api.LoginResponse) {
	path := getAutPath()
	bodyData, _ := json.Marshal(resp)
	_ = ioutil.WriteFile(filepath.Join(path, "auth.json"), bodyData, 0644)
}

type aqualogAPI struct {
	apiAddress  string
	accessToken string
	username    string
	email       string
}

func NewAqualogAPI() aqualogAPI {
	return aqualogAPI{
		apiAddress: apiAddress,
	}
}

func (aqualog *aqualogAPI) LoadAuth() error {
	path := getAutPath()

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
	aqualog.accessToken = data.AccessToken
	aqualog.username = data.User.Username
	aqualog.email = data.User.Email

	return nil
}

func (aqualog *aqualogAPI) Login(username string, password string) error {
	req := api.LoginRequest{
		Username: username,
		Password: password,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpResp, err := http.Post(aqualog.apiAddress+"/login", "application/json", bytes.NewBuffer(data))
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
		return fmt.Errorf("api request error, try again")
	}

	// decode response body
	var resp api.LoginResponse
	json.NewDecoder(httpResp.Body).Decode(&resp)

	aqualog.accessToken = resp.AccessToken
	aqualog.username = resp.User.Username
	aqualog.email = resp.User.Email

	// save login response to home folder
	saveLoginResp(resp)
	return nil
}
