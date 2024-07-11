package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// HandleError prints and exits on error
func HandleError(message string, err error) {
	if err != nil {
		fmt.Println(message, err)
		os.Exit(1)
	}
}

// LoadConfigData loads the configuration data from a file
func LoadConfigData() map[string]string {
	configDir, err := os.UserConfigDir()
	HandleError("Unable to get user config directory", err)

	configFilePath := filepath.Join(configDir, "att_config.json")
	configBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return make(map[string]string)
	}

	var configData map[string]string
	err = json.Unmarshal(configBytes, &configData)
	HandleError("Unable to unmarshal config data", err)

	return configData
}

// MakeAPIRequest makes an API request and returns the response
func MakeAPIRequest(method string, url string, payload []byte, apiToken string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	return client.Do(req)
}
