package notion_api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type NotionApiClient struct {
	Token string
}

type Titleable struct {
	Title string `json:"title"`
}

type TextBlock struct {
	PlainText string `json:"plain_text,omitempty"`
	Href string `json:"href,omitempty"`
}

type RichTextable struct {
	RichText []TextBlock `json:"rich_text"`
}

type ImageFile struct {
	Url string `json:"url"`
}

type Image struct {
	File ImageFile `json:"file"`
}

type Block struct {
	Id string `json:"id"`
	Type string `json:"type"`

	ChildDatabase Titleable `json:"child_database,omitempty"`
	ChildPage Titleable `json:"child_page,omitempty"`
	Heading1 RichTextable `json:"heading_1,omitempty"`
	Heading2 RichTextable `json:"heading_2,omitempty"`
	Heading3 RichTextable `json:"heading_3,omitempty"`
	Paragraph RichTextable `json:"paragraph,omitempty"`
	BulletedListItem RichTextable `json:"bulleted_list_item,omitempty"`
	NumberedListItem RichTextable `json:"numbered_list_item,omitempty"`
	Code RichTextable `json:"code,omitempty"`
	Image Image `json:"image,omitempty"`
}

type ListResponse struct {
	Results []Block `json:"results"`
	Type string `json:"type"`
}

// func (notionApiClient *NotionApiClient) GetBlockChildren(blockId string) (error, map[string]interface{}) {
func (notionApiClient *NotionApiClient) GetBlockChildren(blockId string) (*ListResponse, error) {
	path := fmt.Sprintf("/v1/blocks/%s/children?page_size=100", blockId)
	return notionApiClient.Get(path)
}

// func (notionApiClient *NotionApiClient) Get(path string) (error, map[string]interface{}) {
func (notionApiClient *NotionApiClient) Get(path string) (*ListResponse, error) {
	version := "2022-02-22"
	base_url := "https://api.notion.com"
	url := base_url + path

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+ notionApiClient.Token)
	request.Header.Set("Notion-Version", version)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseJson, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// var responseMap map[string]interface{}
	// err = json.Unmarshal(responseJson, &responseMap)
	// if err != nil {
	// 	return nil, err
	// }
	// return nil, responseMap
	var parsedResponse ListResponse
	err = json.Unmarshal(responseJson, &parsedResponse)
	if err != nil {
		return nil, err
	}
	return &parsedResponse, nil
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
