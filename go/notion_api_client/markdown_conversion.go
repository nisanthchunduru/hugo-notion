package notion_api_client

import (
	"strings"
	"fmt"
)

func ConvertBlocksToMarkdown(blocks []Block) string {
	var markdowns []string
	for i, block := range blocks {
		var markdown string

		if block.Type == "heading_1" {
			markdown = ConvertHeading1ToMarkdown(block.Heading1)
		} else if block.Type == "heading_2" {
			markdown = ConvertHeading2ToMarkdown(block.Heading2)
		} else if block.Type == "heading_3" {
			markdown = ConvertHeading3ToMarkdown(block.Heading3)
		} else if block.Type == "paragraph" {
			markdown = ConvertParagraphToMarkdown(block.Paragraph)
		} else if block.Type == "bulleted_list_item" {
			markdown = ConvertBulletedListItemToMarkdown(block.BulletedListItem)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i + 1]
				if nextBlock.Type != "bulleted_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.Type == "numbered_list_item" {
			markdown = ConvertNumberedListItemToMarkdown(block.NumberedListItem)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i + 1]
				if nextBlock.Type != "numbered_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.Type == "code" {
			markdown = ConvertCodeToMarkdown(block.Code)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i + 1]
				if nextBlock.Type != "numbered_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.Type == "image" {
			markdown = ConvertImageToMarkdown(block.Image)
		}

		if markdown != "" {
			markdowns = append(markdowns, markdown)
		}
	}
	return strings.Join(markdowns, "")
}

func ConvertHeading1ToMarkdown(content RichTextable) string {
	markdown := ConvertTextBlocksToMarkdown(content.RichText)
	markdown = "# " + markdown + "\n\n"
	return markdown
}

func ConvertHeading2ToMarkdown(content RichTextable) string {
	markdown := ConvertTextBlocksToMarkdown(content.RichText)
	markdown = "## " + markdown + "\n\n"
	return markdown
}

func ConvertHeading3ToMarkdown(content RichTextable) string {
	markdown := ConvertTextBlocksToMarkdown(content.RichText)
	markdown = "### " + markdown + "\n\n"
	return markdown
}

func ConvertParagraphToMarkdown(content RichTextable) string {
	markdown := ConvertTextBlocksToMarkdown(content.RichText)
	markdown = markdown + "\n"
	return markdown
}

func ConvertBulletedListItemToMarkdown(content RichTextable) string {
	markdown := ConvertTextBlocksToMarkdown(content.RichText)
	markdown = "- " + markdown + "\n"
	return markdown
}

func ConvertNumberedListItemToMarkdown(content RichTextable) string {
	markdown := ConvertTextBlocksToMarkdown(content.RichText)
	markdown = "1. " + markdown + "\n"
	return markdown
}

func ConvertCodeToMarkdown(content RichTextable) string {
	markdown := ConvertTextBlocksToMarkdown(content.RichText)
	markdown = "```\n" + markdown + "```\n"
	return markdown
}

func ConvertImageToMarkdown(image Image) string {
	markdown := fmt.Sprintf("![Untitled](%s)", image.File.Url)
	return markdown
}

func ConvertTextBlocksToMarkdown(textBlocks []TextBlock) string {
	var markdowns []string
	for _, block := range textBlocks {
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
