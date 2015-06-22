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
  exec "go run api_wsc.go"
end

desc "Starts API in development env with autorefreshing"
task :rego do
  exec "fresh"
end


namespace :build do
  desc "Build linux release"
  task :linux do
    system "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/api-wsc-linux api_wsc.go"
  end

  desc "Build mac release"
  task :osx do
    system "CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o release/api-wsc-osx api_wsc.go"
  end

  desc "Cleanup release"
  task :clean do
    system "rm -Rf release/*"
  end

  desc "Copy project files into release"
  task :provision do
    mkdir "release", :noop => true
    %w|doc public config README.md CHANGELOG.md|.map do |f|
      cp_r f, "release"
    end
  end

  desc " * * * Build, create and copy project structure * * *"
  task :release => [:clean, :linux, :provision]

end
