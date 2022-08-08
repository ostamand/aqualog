package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ostamand/aqualog/api"
)

func GetAutPath() string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".aqualog")
	os.MkdirAll(path, os.ModePerm)
	return path
}

func SaveLoginResp(resp api.LoginResponse) {
	path := GetAutPath()
	bodyData, _ := json.Marshal(resp)
	_ = ioutil.WriteFile(filepath.Join(path, "auth.json"), bodyData, 0644)
}
