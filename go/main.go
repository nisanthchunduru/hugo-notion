package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		fmt.Println("Please create a Notion integration, generate a secret and provide it in the 'NOTION_TOKEN' environment variable")
		os.Exit(1)
	}

	version := "2022-02-22"
	url := "https://api.notion.com/v1/blocks/0f1b55769779411a95df1ee9b4b070c9/children?page_size=100"

	// response, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// responseBody, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(responseBody))

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		printErrorAndExit(err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Notion-Version", version)
	response, err := client.Do(request)
	if err != nil {
		printErrorAndExit(err)
	}
	defer response.Body.Close()
	responseJson, err := io.ReadAll(response.Body)
	if err != nil {
		printErrorAndExit(err)
	}

	// fmt.Println(string(responseBody))

	// prettyPrintResponseJson(string(responseJson))

	var responseMap map[string]interface{}
	err = json.Unmarshal(responseJson, &responseMap)
	if err != nil {
		printErrorAndExit(err)
	}
}

func printErrorAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func prettyPrintResponseJson(responseJson string) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(responseJson), "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(prettyJSON.Bytes()))
}
