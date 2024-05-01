package main

import (
	"fmt"
	"os"
	"context"

	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/nisanthchunduru/hugo-notion/notion_markdown_exporter"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		printErrorAndExit(err)
	}

	tokenString := os.Getenv("NOTION_TOKEN")
	if tokenString == "" {
		fmt.Println("Please create a Notion integration, generate a secret and provide it in the 'NOTION_TOKEN' environment variable")
		os.Exit(1)
	}
	var token notionapi.Token
	token = notionapi.Token(tokenString)

	// contentNotionPageId := "0f1b55769779411a95df1ee9b4b070c9"
	var contentNotionPageId notionapi.BlockID
	contentNotionPageId = "0f1b55769779411a95df1ee9b4b070c9"

	jomeiNotionApiClient := notionapi.NewClient(token)
	pagination := notionapi.Pagination{
		PageSize: 100,
	}
	getChildrenResponse, err := jomeiNotionApiClient.Block.GetChildren(context.Background(), contentNotionPageId, &pagination)
	if err != nil {
		printErrorAndExit(err)
	}
	for _, block := range getChildrenResponse.Results {
		if (block.GetType() == "child_page") {
			childPageId := block.GetID()
			getChildPageChildrenResponse, err := jomeiNotionApiClient.Block.GetChildren(context.Background(), childPageId, &pagination)
			if err != nil {
				printErrorAndExit(err)
			}
			childPageBlocks := getChildPageChildrenResponse.Results
			markdown := notion_markdown_exporter.ConvertBlocksToMarkdown(childPageBlocks)
			fmt.Println(markdown)
		}
	}
}

func printErrorAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
