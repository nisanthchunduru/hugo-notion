package notion_markdown_exporter

import (
	"github.com/jomei/notionapi"
	"strings"
	"fmt"
)

func ConvertBlocksToMarkdown(blocks []notionapi.Block) string {
	var markdowns []string
	for i, block := range blocks {
		var markdown string

		if block.GetType() == "heading_1" {
			heading1Block := block.(*notionapi.Heading1Block)
			markdown = ConvertHeading1ToMarkdown(heading1Block.Heading1)
		} else if block.GetType() == "heading_2" {
			heading2Block := block.(*notionapi.Heading2Block)
			markdown = ConvertHeading2ToMarkdown(heading2Block.Heading2)
		} else if block.GetType() == "heading_3" {
			heading3Block := block.(*notionapi.Heading3Block)
			markdown = ConvertHeading3ToMarkdown(heading3Block.Heading3)
		} else if block.GetType() == "paragraph" {
			paragraphBlock := block.(*notionapi.ParagraphBlock)
			markdown = ConvertParagraphToMarkdown(paragraphBlock.Paragraph)
		} else if block.GetType() == "bulleted_list_item" {
			bulletedListItemBlock := block.(*notionapi.BulletedListItemBlock)
			markdown = ConvertBulletedListItemToMarkdown(bulletedListItemBlock.BulletedListItem)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i + 1]
				if nextBlock.GetType() != "bulleted_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.GetType() == "numbered_list_item" {
			numberedListItemBlock := block.(*notionapi.NumberedListItemBlock)
			markdown = ConvertNumberedListItemToMarkdown(numberedListItemBlock.NumberedListItem)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i + 1]
				if nextBlock.GetType() != "numbered_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.GetType() == "code" {
			codeBlock := block.(*notionapi.CodeBlock)
			markdown = ConvertCodeToMarkdown(codeBlock.Code)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i + 1]
				if nextBlock.GetType() != "numbered_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.GetType() == "image" {
			imageBlock := block.(*notionapi.ImageBlock)
			markdown = ConvertImageToMarkdown(imageBlock.Image)
		}

		if markdown != "" {
			markdowns = append(markdowns, markdown)
		}
	}
	return strings.Join(markdowns, "")
}

func ConvertHeading1ToMarkdown(heading1 notionapi.Heading) string {
	markdown := ConvertRichTextsToMarkdown(heading1.RichText)
	markdown = "# " + markdown + "\n\n"
	return markdown
}

func ConvertHeading2ToMarkdown(heading2 notionapi.Heading) string {
	markdown := ConvertRichTextsToMarkdown(heading2.RichText)
	markdown = "## " + markdown + "\n\n"
	return markdown
}

func ConvertHeading3ToMarkdown(heading3 notionapi.Heading) string {
	markdown := ConvertRichTextsToMarkdown(heading3.RichText)
	markdown = "### " + markdown + "\n\n"
	return markdown
}

func ConvertParagraphToMarkdown(paragraph notionapi.Paragraph) string {
	markdown := ConvertRichTextsToMarkdown(paragraph.RichText)
	markdown = markdown + "\n"
	return markdown
}

func ConvertBulletedListItemToMarkdown(bulleted_list_item notionapi.ListItem) string {
	markdown := ConvertRichTextsToMarkdown(bulleted_list_item.RichText)
	markdown = "- " + markdown + "\n"
	return markdown
}

func ConvertNumberedListItemToMarkdown(numbered_list_item notionapi.ListItem) string {
	markdown := ConvertRichTextsToMarkdown(numbered_list_item.RichText)
	markdown = "1. " + markdown + "\n"
	return markdown
}

func ConvertCodeToMarkdown(code notionapi.Code) string {
	markdown := ConvertRichTextsToMarkdown(code.RichText)
	markdown = "```\n" + markdown + "```\n"
	return markdown
}

func ConvertImageToMarkdown(image notionapi.Image) string {
	markdown := fmt.Sprintf("![Untitled](%s)", image.File.URL)
	return markdown
}

func ConvertRichTextsToMarkdown(richTexts []notionapi.RichText) string {
	var markdowns []string
	for _, block := range richTexts {
		var markdown string;
		if block.Href == "" {
			markdown = block.PlainText
		} else {
			markdown = fmt.Sprintf("[%s](%s)", block.PlainText, block.Href)
		}
		markdowns = append(markdowns, markdown)
	}
	return strings.Join(markdowns, "")
}
