require_relative '../lib/hugo_notion/synchronizer'

def do_exit
  puts 'Exiting...'
  exit
end
Signal.trap("TERM") { do_exit }
Signal.trap("INT") { do_exit }

if !ARGV[0]
  puts "Please provide the URL of the Notion page you'd like to sync as the first argument"
end
notion_page_url = ARGV[0]
notion_page_id = URI(notion_page_url).path.match(/-(?<page_id>.*)\z/)['page_id']

destination_dir = Dir.pwd
if ARGV[1]
  if Pathname.new(ARGV[1]).absolute?
    destination_dir = ARGV[1]
  else
    destination_dir = File.join(Dir.pwd, ARGV[1])
  end
end

loop do
  puts "Syncing posts from Notion..."
  Synchronizer.run(
    notion_page_id: notion_page_id,
    destination_dir: destination_dir
  )
  puts "Done."
  sleep 10
end