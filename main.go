package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
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

	notionToken := os.Getenv("NOTION_TOKEN")
	if notionToken == "" {
		fmt.Println("Please create a Notion integration, generate a secret and provide it in the 'NOTION_TOKEN' environment variable")
		os.Exit(1)
	}

	contentDir := os.Getenv("CONTENT_DIR")
	if contentDir == "" {
		contentDir, _ = filepath.Abs("./content")
	}

	var contentNotionUrl string
	contentNotionUrl = os.Args[len(os.Args)-1]

	if !isValidUrl(contentNotionUrl) && os.Getenv("CONTENT_NOTION_URL") != "" {
		contentNotionUrl = os.Getenv("CONTENT_NOTION_URL")
	}
	if contentNotionUrl == "" {
		fmt.Println("Please provide the URL of the Notion page you'd like to sync in the `CONTENT_NOTION_URL` environment variable or as the first argument")
		os.Exit(1)
	}
	parsedContentNotionUrl, err := url.ParseRequestURI(contentNotionUrl)
	if err != nil {
		fmt.Println("The Notion URL you've provided is not a valid URL. Please provide a valid URL.")
		os.Exit(1)
	}
	pathFragments := strings.Split(parsedContentNotionUrl.Path, "-")
	contentNotionPageId := pathFragments[len(pathFragments)-1]

	jomeiNotionApiClient := notionapi.NewClient(notionapi.Token(notionToken))
	fmt.Println("Syncing content from Notion...")
	syncNotionPage(jomeiNotionApiClient, contentNotionPageId, contentDir)
	fmt.Println("Done.")
}

func isValidUrl(url string) bool {
	return govalidator.IsRequestURL(url)
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
