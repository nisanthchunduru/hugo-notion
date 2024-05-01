package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/nisanthchunduru/hugo-notion/notion_api_client"
	// "github.com/kjk/notionapi"
	// "github.com/kjk/notionapi/tomarkdown"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		printErrorAndExit(err)
	}

	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		fmt.Println("Please create a Notion integration, generate a secret and provide it in the 'NOTION_TOKEN' environment variable")
		os.Exit(1)
	}

	notionApiClient := &notion_api_client.NotionApiClient{
		Token: token,
	}
	if _, err := os.Stat("content"); os.IsNotExist(err) {
		os.Mkdir("content", 0755)
	}

	// err, responseMap := notionApiClient.Get("/v1/blocks/0f1b55769779411a95df1ee9b4b070c9/children?page_size=100")
	// if err != nil {
	// 	printErrorAndExit(err)
	// }

	// fmt.Println(responseMap)

	// for _, _block := range responseMap["results"].([]interface{}) {
	// 	block := _block.(map[string]interface{})
	// 	if block["type"] == "child_page" {
	// 		childPageId := block["id"].(string)

	// 		childPageTitle := block["child_page"].(map[string]interface{})["title"].(string)
	// 		err, getChildPageChildrenResponseMap := notionApiClient.GetBlockChildren(childPageId)
	// 		if err != nil {
	// 			printErrorAndExit(err)
	// 		}
	// 		fmt.Println(getChildPageChildrenResponseMap)

	// 		tokenV2CookieString := "dummyTokenV2CookieString"
	// 		kjkNotionApiClient := &notionapi.Client{
	// 			AuthToken: tokenV2CookieString,
	// 		}
	// 		markdown, err := kjkNotionApiClient.ExportPages(childPageId, "markdown", false)
	// 		if err != nil {
	// 			printErrorAndExit(fmt.Errorf("MarkdownExportFailed: %w", err))
	// 		}
	// 		fmt.Println(markdown)
	// 		childPage, err := kjkNotionApiClient.DownloadPage(childPageId)
	// 		if err != nil {
	// 			printErrorAndExit(err)
	// 		}
	// 		markdown := tomarkdown.NewConverter(childPage).ToMarkdown()
	// 		fmt.Println(string(markdown))
	// 	}
	// }

	response, err := notionApiClient.Get("/v1/blocks/0f1b55769779411a95df1ee9b4b070c9/children?page_size=100")
	if err != nil {
		printErrorAndExit(err)
	}
	for _, block := range response.Results {
		if block.Type == "child_page" {
			childPageId := block.Id
			getBlockChildrenResponse, err := notionApiClient.GetBlockChildren(childPageId)
			if err != nil {
				printErrorAndExit(err)
			}
			fmt.Println(notion_api_client.ConvertBlocksToMarkdown(getBlockChildrenResponse.Results))
			// for _, childPageBlock := range getBlockChildrenResponse.Results {
			// 	if childPageBlock.Type == "heading_1" {
			// 		for _, richTextBlock := range childPageBlock.Heading1.RichText {
			// 			fmt.Println(richTextBlock.PlainText)
			// 		}
			// 	}
			// }
		}
	}
}

func printErrorAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
