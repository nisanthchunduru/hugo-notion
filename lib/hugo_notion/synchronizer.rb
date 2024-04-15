env_file_path = File.expand_path('../../.env', File.dirname(__FILE__))
if File.exist?(env_file_path)
  require 'dotenv'
  Dotenv.load(env_file_path)
end
notion_token = ENV['NOTION_TOKEN']
unless notion_token
  throw "Please create a Notion integration, generate a secret and provide it in the 'NOTION_TOKEN' environment variable"
end

NOTION_TOKEN = ENV['NOTION_TOKEN']

require 'httparty'
require 'notion-ruby-client'
require 'notion_to_md'
require 'yaml'
require 'pry'

class NotionApi
  BASE_URL = "https://api.notion.com/v1"

  def initialize(api_secret)
    @api_secret = api_secret
  end

  def get(path)
    url = File.join(BASE_URL, path)
    HTTParty.get(url, headers: {
      "Notion-Version" => '2022-02-22',
      "Authorization" => "Bearer #{@api_secret}"
    })
  end

  def post(path)
    url = File.join(BASE_URL, path)
    HTTParty.post(url, headers: {
      "Notion-Version" => '2022-02-22',
      "Authorization" => "Bearer #{@api_secret}"
    })
  end
end

class Synchronizer
  class << self
    def run(options = {})
      if options[:page_id]
        page_id = options[:page_id]
      elsif options[:database_id]
        database_id = options[:database_id]
      end
      destination_dir = options[:destination_dir]
      Dir.mkdir(destination_dir) unless Dir.exist?(destination_dir)

      notion = NotionApi.new(NOTION_TOKEN)
      response = if database_id
        notion.post("/databases/#{database_id}/query")
      elsif page_id
        notion.get("/blocks/#{page_id}/children")
      end
      blocks = response['results']

      existing_blog_page_file_names = Dir.entries(destination_dir).select { |f| !File.directory? f }
      blog_page_file_names = []
      blocks.each do |block|
        block_id = block['id']

        if block['type'] == 'child_database'
          child_database_title = block['child_database']['title']
          run(
            database_id: block_id,
            destination_dir: File.join(destination_dir, child_database_title)
          )
          next
        end

        block_last_edited_at = Time.parse(block['last_edited_time'])

        blog_page_front_matter = {
          'date' => Time.parse(block['created_time'])
        }
        if block['properties']
          if block.dig('properties', 'date', 'date', 'start')
            blog_page_front_matter['date'] = Time.parse(block['properties']['date']['date']['start'])
          end

          if block.dig('properties', 'Name', 'title', 0, 'plain_text')
            blog_page_front_matter['title'] = block['properties']['Name']['title'][0]['plain_text']
          end
        end
        if block['type'] == 'child_page'
          blog_page_front_matter['title'] = block['child_page']['title']
          blog_page_front_matter['type'] = block['child_page']['title']
        end
        blog_page_title = blog_page_front_matter['title']

        # Ignore incomplete pages
        next unless blog_page_front_matter['date']
        next unless blog_page_front_matter['title']

        blog_page_file_name_without_extension = blog_page_title.gsub(' ', '-')
        blog_page_file_name = "#{blog_page_file_name_without_extension}.md"
        blog_page_file_names << blog_page_file_name

        blog_page_path = File.join(destination_dir, blog_page_file_name)
        unless ENV['NOTION_SYNCHRONIZER_SKIP_OPTIMIZATIONS']
          next if File.exist?(blog_page_path) && File.mtime(blog_page_path) >= block_last_edited_at
        end

        blog_page_front_matter_yaml = blog_page_front_matter.to_yaml.chomp
        blog_page_markdown = NotionToMd.convert(page_id: block_id, token: NOTION_TOKEN)
        blog_page_content = <<-BLOG_PAGE_CONTENT
#{blog_page_front_matter_yaml}
---

#{blog_page_markdown}
BLOG_PAGE_CONTENT
        File.write(blog_page_path, blog_page_content)
      end
      old_blog_page_file_names = existing_blog_page_file_names - blog_page_file_names
      old_blog_page_file_names.each do |file_name|
        file_path = File.join(destination_dir, file_name)
        File.delete(file_path) if File.exist?(file_path) && !File.directory?(file_path)
      end
    end
  end
end
