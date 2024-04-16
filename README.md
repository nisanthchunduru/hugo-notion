# hugo-notion

Write in Notion. Publish with Hugo.

## Installation

```
gem install hugo-notion --prerelease
```

Installing the `hugo-notion` ruby gem will install the `huno` command.

## Usage

First, create a Notion integration, generate a secret and connect that integration to the Notion page https://developers.notion.com/docs/create-a-notion-integration#getting-started

Go to your Hugo site directory and run

```
NOTION_TOKEN=your_notion_secret huno your_notion_page_url
```

`huno` will sync your Notion page and its children pages to the `content` directory.

If you're yet to move your Hugo pages to Notion, you can use my "blog_content" Notion page as a template https://www.notion.so/blog_content-0f1b55769779411a95df1ee9b4b070c9

## Tips

If you'd like `huno` to sync Notion pages to a different directory, you can do that too

```
NOTION_TOKEN=your_notion_secret huno your_notion_page_url site_content/
```

To avoid having to provide the `NOTION_TOKEN` env var again and again, you can create a `.env` file

```
echo 'NOTION_TOKEN=your_notion_secret' > .env
```

To run `huno` say, every 15 seconds, use the `watch` command

```
watch -n15 huno your_notion_page_url
```

If you're on MacOS and don't have the `watch` command installed, you can use Homebrew to install it

```
brew install watch
```

## Bug Reports

If you'd like to report a bug (if there are any, please do), please create a GitHub issue
