package backend

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
)

type GoogleInfo struct {
	Info GoogleCredential `json:"web"`
}

type GoogleCredential struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func LoadGoogleCredential() *GoogleCredential {
	bytes, err := ioutil.ReadFile(filepath.Join(BasePath(), "../credential/google.json"))
	if err != nil {
		log.Fatal("cannot load google.json")
	}

	info := &GoogleInfo{}

	json.Unmarshal(bytes, info)

	return &info.Info
}

func LoadSecret() string {
	bytes, err := ioutil.ReadFile(filepath.Join(BasePath(), "../credential/secret.json"))
	if err != nil {
		log.Fatal("cannot load secret.json")
	}

	return string(bytes)
}

func BasePath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
