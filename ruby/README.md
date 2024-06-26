# ⚠️ Deprecation Notice

This directory contains hugo-nation's older Ruby implementation which is functional but not actively maintained anymore

# hugo-notion

Write in Notion. Publish with Hugo.

Use Notion as a CMS (Content Management System) for your Hugo site/blog.

`hugo-notion` is a command line (CLI) tool that syncs your Notion page url to your Hugo site's/blog's `content` directory.

## Installation

`hugo-notion` is a Ruby gem. Use the `gem` command to install it

```
gem install hugo-notion --prerelease
```

Installing the `hugo-notion` ruby gem will install the `huno` command

## Usage

First, create a Notion integration, generate a secret and connect that integration to the Notion page https://developers.notion.com/docs/create-a-notion-integration#getting-started

Go to your Hugo site directory and run

```
NOTION_TOKEN=your_notion_secret huno your_notion_page_url
```

`huno` will sync your Notion page and its children pages to the `content` directory.

`huno` can also sync your Notion page periodically every 10 seconds. To do so, run

```
huno -r
```

If you'd like to sync at a different frequency (say, 5 seconds), run

```
huno -r 5
```
To avoid the hassle of providing your Notion token and your Notion page url to `huno` every time you run it, create an .env file

```
echo 'NOTION_TOKEN=your_notion_secret' > .env
echo 'CONTENT_NOTION_URL=your_notion_page_url >> .env'
```

### Migration

For an easy migration to Notion, you can use my "blog_content" Notion page as a template https://www.notion.so/ blog_content-0f1b55769779411a95df1ee9b4b070c9

I recommend that you move one page from Notion to Hugo first, try hugo-notion to sync that page and once you're happy with hugo-notion, move your other Hugo pages to Notion one by one.

## Bug Reports

If you hit a bug, please do report it by creating a GitHub issue
