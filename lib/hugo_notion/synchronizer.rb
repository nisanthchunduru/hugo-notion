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
      if options[:notion_page_id]
        notion_page_id = options[:notion_page_id]
      elsif options[:notion_database_id]
        notion_database_id = options[:notion_database_id]
      end
      destination_dir = options[:destination_dir]
      Dir.mkdir(destination_dir) unless Dir.exist?(destination_dir)

      notion = NotionApi.new(ENV.fetch('NOTION_TOKEN'))
      response = if notion_database_id
        notion.post("/databases/#{notion_database_id}/query")
      elsif notion_page_id
        notion.get("/blocks/#{notion_page_id}/children")
      end
      notion_blocks = response['results']

      existing_page_file_names = Dir.entries(destination_dir).select { |f| !File.directory? f }
      page_file_names = []
      notion_blocks.each do |notion_block|
        notion_block_id = notion_block['id']

        if notion_block['type'] == 'child_database'
          notion_child_database_title = notion_block['child_database']['title']
          run(
            notion_database_id: notion_block_id,
            destination_dir: File.join(destination_dir, notion_child_database_title)
          )
          next
        end

        block_last_edited_at = Time.parse(notion_block['last_edited_time'])

        page_front_matter = {
          'date' => Time.parse(notion_block['created_time'])
        }
        if notion_block['properties']
          if notion_block.dig('properties', 'date', 'date', 'start')
            page_front_matter['date'] = Time.parse(notion_block['properties']['date']['date']['start'])
          end

          if notion_block.dig('properties', 'Name', 'title', 0, 'plain_text')
            page_front_matter['title'] = notion_block['properties']['Name']['title'][0]['plain_text']
          end
        end
        if notion_block['type'] == 'child_page'
          page_front_matter['title'] = notion_block['child_page']['title']
          page_front_matter['type'] = notion_block['child_page']['title']
        end
        page_title = page_front_matter['title']

        # Ignore incomplete notion pages
        next unless page_front_matter['date']
        next unless page_front_matter['title']

        page_file_name_without_extension = page_title.gsub(' ', '-')
        page_file_name = "#{page_file_name_without_extension}.md"
        page_file_names << page_file_name

        page_path = File.join(destination_dir, page_file_name)
        unless ENV['NOTION_SYNCHRONIZER_SKIP_OPTIMIZATIONS']
          next if File.exist?(page_path) && File.mtime(page_path) >= block_last_edited_at
        end

        page_front_matter_yaml = page_front_matter.to_yaml.chomp
        page_markdown = NotionToMd.convert(page_id: notion_block_id, token: ENV.fetch('NOTION_TOKEN'))
        page_content = <<-page_CONTENT
#{page_front_matter_yaml}
---

#{page_markdown}
page_CONTENT
        File.write(page_path, page_content)
      end
      old_page_file_names = existing_page_file_names - page_file_names
      old_page_file_names.each do |file_name|
        file_path = File.join(destination_dir, file_name)
        File.delete(file_path) if File.exist?(file_path) && !File.directory?(file_path)
      end
    end
  end
end
