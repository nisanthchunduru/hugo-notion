package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
	"github.com/nisanthchunduru/hugo-notion/notion_markdown_exporter"
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

	hugoCONTENTDIR := os.Getenv("HUGO_CONTENT_DIR")
	if hugoContentDir == "" {
		hugoContentDir, _ = filepath.Abs("./content")
	}

	pageId := "0f1b55769779411a95df1ee9b4b070c9"
	jomeiNotionApiClient := notionapi.NewClient(notionapi.Token(token))
	syncNotionPage(jomeiNotionApiClient, pageId, hugoContentDir)
}

func syncNotionPage(jomeiNotionApiClient *notionapi.Client, pageIdString string, destinationDir string) {
	pageId := notionapi.BlockID(pageIdString)
	pagination := notionapi.Pagination{
		PageSize: 100,
	}
	getChildrenResponse, err := jomeiNotionApiClient.Block.GetChildren(context.Background(), pageId, &pagination)
	if err != nil {
		printErrorAndExit(err)
	}
	for _, _block := range getChildrenResponse.Results {
		if _block.GetType() == "child_page" {
			block := _block.(*notionapi.ChildPageBlock)
			childPageTitle := block.ChildPage.Title
			hugoPageFrontMatterMap := make(map[string]string)
			hugoPageFrontMatterMap["title"] = childPageTitle
			hugoPageFrontMatterMap["type"] = childPageTitle
			hugoPageFrontMatterMap["date"] = block.GetLastEditedTime().Format(time.RFC3339)
			hugoFrontMatterYaml, err := yaml.Marshal(hugoPageFrontMatterMap)
			if err != nil {
				printErrorAndExit(err)
			}

			childPageId := block.GetID()
			getChildPageChildrenResponse, err := jomeiNotionApiClient.Block.GetChildren(context.Background(), childPageId, &pagination)
			if err != nil {
				printErrorAndExit(err)
			}
			childPageBlocks := getChildPageChildrenResponse.Results
			markdown := notion_markdown_exporter.ConvertBlocksToMarkdown(childPageBlocks)

			hugoPageText := fmt.Sprintf("---\n%s\n---\n\n%s", hugoFrontMatterYaml, markdown)
			hugoPageFileName := strings.ReplaceAll(childPageTitle, " ", "-") + ".md"
			hugoPageDir := destinationDir
			err = os.MkdirAll(hugoPageDir, 0755)
			if err != nil {
				printErrorAndExit(err)
			}
			hugoPageFilePath := filepath.Join(destinationDir, hugoPageFileName)
			err = os.WriteFile(hugoPageFilePath, []byte(hugoPageText), 0644)
			if err != nil {
				printErrorAndExit(err)
			}
		} else if _block.GetType() == "child_database" {
			block := _block.(*notionapi.ChildDatabaseBlock)
			childDatabaseId := notionapi.DatabaseID(block.GetID())
			childDatabaseTitle := block.ChildDatabase.Title
			databaseQueryRequest := notionapi.DatabaseQueryRequest{
				PageSize: 100,
			}
			databaseQueryResponse, err := jomeiNotionApiClient.Database.Query(context.Background(), childDatabaseId, &databaseQueryRequest)
			if err != nil {
				printErrorAndExit(err)
			}
			for _, block := range databaseQueryResponse.Results {
				hugoPageFrontMatterMap := make(map[string]string)
				childPageTitleProperty := block.Properties["Name"].(*notionapi.TitleProperty)
				childPageTitle := childPageTitleProperty.Title[0].PlainText
				hugoPageFrontMatterMap["title"] = childPageTitle
				childPageDateProperty := block.Properties["date"].(*notionapi.DateProperty)
				hugoPageFrontMatterMap["date"] = childPageDateProperty.Date.Start.String()

				hugoFrontMatterYaml, err := yaml.Marshal(hugoPageFrontMatterMap)
				if err != nil {
					printErrorAndExit(err)
				}

				childPageId := notionapi.BlockID(block.ID)
				getChildPageChildrenResponse, err := jomeiNotionApiClient.Block.GetChildren(context.Background(), childPageId, &pagination)
				if err != nil {
					printErrorAndExit(err)
				}
				childPageBlocks := getChildPageChildrenResponse.Results
				markdown := notion_markdown_exporter.ConvertBlocksToMarkdown(childPageBlocks)

				hugoPageText := fmt.Sprintf("---\n%s\n---\n\n%s", hugoFrontMatterYaml, markdown)
				hugoPageFileName := strings.ReplaceAll(childPageTitle, " ", "-") + ".md"
				hugoPageDir := filepath.Join(destinationDir, childDatabaseTitle)
				err = os.MkdirAll(hugoPageDir, 0755)
				if err != nil {
					printErrorAndExit(err)
				}
				hugoPageFilePath := filepath.Join(hugoPageDir, hugoPageFileName)
				err = os.WriteFile(hugoPageFilePath, []byte(hugoPageText), 0644)
				if err != nil {
					printErrorAndExit(err)
				}
			}
		}
	}
}

func printErrorAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
