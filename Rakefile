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

  Seed.insert_customers
  Seed.insert_orders

end
