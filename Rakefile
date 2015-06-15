#!/usr/bin/env ruby

begin
  gem "mongo"
rescue Gem::LoadError
  puts ".. installing mongo gem"
  `gem i mongo`
end


task default: :seed


desc "Clean db and make seed of some data"
task :seed do
  require_relative 'spec/fixtures/seed'

  Seed.clear_all

  Seed.fill_docs

end

desc "Starts API in development env"
task :go do
  system "go run api_wsc.go"
end
