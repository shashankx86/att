package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"att/utils"
)

const BASE_URL = "https://hackhour.hackclub.com"

// FetchAndPrintData fetches and prints data from the API
func FetchAndPrintData(endpoint string) {
	configData := utils.LoadConfigData()
	apiToken := configData["api-token"]
	slackID := configData["slack-id"]

	if apiToken == "" || slackID == "" {
		fmt.Println("Please set your API token and Slack ID using the configure command.")
		return
	}

	url := fmt.Sprintf("%s/api/%s/%s", BASE_URL, endpoint, slackID)
	resp, err := utils.MakeAPIRequest("GET", url, nil, apiToken)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: received status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	utils.HandleError("Unable to read response body", err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	utils.HandleError("Unable to unmarshal response", err)

	// Print the response data
	if result["ok"].(bool) {
		data, _ := json.MarshalIndent(result["data"], "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Error: unable to fetch data")
	}
}

// StartNewSession starts a new session
func StartNewSession(work string) {
	configData := utils.LoadConfigData()
	apiToken := configData["api-token"]
	slackID := configData["slack-id"]

	if apiToken == "" || slackID == "" {
		fmt.Println("Please set your API token and Slack ID using the configure command.")
		return
	}

	url := fmt.Sprintf("%s/api/start/%s", BASE_URL, slackID)
	payload := map[string]string{"work": work}
	payloadBytes, err := json.Marshal(payload)
	utils.HandleError("Unable to marshal request payload", err)

	resp, err := utils.MakeAPIRequest("POST", url, payloadBytes, apiToken)
	utils.HandleError("Unable to make request", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	utils.HandleError("Unable to read response body", err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	utils.HandleError("Unable to unmarshal response", err)

	if resp.StatusCode != http.StatusOK {
		if errorMessage, exists := result["error"].(string); exists && errorMessage == "You already have an active session" {
			fmt.Println("You already have an active session")
		} else {
			fmt.Printf("Error: received status code %d with message: %s\n", resp.StatusCode, result["error"])
		}
		return
	}

	if ok, exists := result["ok"].(bool); exists && ok {
		data, _ := json.MarshalIndent(result["data"], "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Error:", result["error"])
	}
}

// PauseOrResumeSession pauses or resumes the current session
func PauseOrResumeSession() {
	configData := utils.LoadConfigData()
	apiToken := configData["api-token"]
	slackID := configData["slack-id"]

	if apiToken == "" || slackID == "" {
		fmt.Println("Please set your API token and Slack ID using the configure command.")
		return
	}

	url := fmt.Sprintf("%s/api/pause/%s", BASE_URL, slackID)
	resp, err := utils.MakeAPIRequest("POST", url, nil, apiToken)
	utils.HandleError("Unable to make request", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	utils.HandleError("Unable to read response body", err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	utils.HandleError("Unable to unmarshal response", err)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d with message: %s\n", resp.StatusCode, result["error"])
		return
	}

	if ok, exists := result["ok"].(bool); exists && ok {
		data, _ := json.MarshalIndent(result["data"], "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Error:", result["error"])
	}
}

// CancelSession cancels the current session
func CancelSession() {
	configData := utils.LoadConfigData()
	apiToken := configData["api-token"]
	slackID := configData["slack-id"]

	if apiToken == "" || slackID == "" {
		fmt.Println("Please set your API token and Slack ID using the configure command.")
		return
	}

	url := fmt.Sprintf("%s/api/cancel/%s", BASE_URL, slackID)
	resp, err := utils.MakeAPIRequest("POST", url, nil, apiToken)
	utils.HandleError("Unable to make request", err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	utils.HandleError("Unable to read response body", err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	utils.HandleError("Unable to unmarshal response", err)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d with message: %s\n", resp.StatusCode, result["error"])
		return
	}

	if ok, exists := result["ok"].(bool); exists && ok {
		data, _ := json.MarshalIndent(result["data"], "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Error:", result["error"])
	}
}
