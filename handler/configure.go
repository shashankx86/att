package handler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"att/utils"
)

// UpdateConfigData updates the configuration data
func UpdateConfigData(key string, value string) {
	configData := utils.LoadConfigData()
	configData[key] = value
	saveConfigData(configData)
}

// saveConfigData saves the configuration data to a file
func saveConfigData(configData map[string]string) {
	configDir, err := os.UserConfigDir()
	utils.HandleError("Unable to get user config directory", err)

	configFilePath := filepath.Join(configDir, "att_config.json")
	configFile, err := os.Create(configFilePath)
	utils.HandleError("Unable to create config file", err)
	defer configFile.Close()

	configBytes, err := json.MarshalIndent(configData, "", "  ")
	utils.HandleError("Unable to marshal config data", err)

	_, err = configFile.Write(configBytes)
	utils.HandleError("Unable to write config data to file", err)
}
