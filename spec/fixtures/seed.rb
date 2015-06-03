require "mongo"

$client = Mongo::Client.new([ '127.0.0.1:27017' ], database: 'apiwholesale_dev', connect: :direct)


module Seed extend self

  def clear_orders
    puts "... clear orders"
    $client[:orders].drop
  end

  def clear_customers
    puts "... clear customers"
    $client[:customers].drop
  end

  def clear_all
    clear_customers
    clear_orders
  end


  def insert_customers
    puts "... inserting customers"
    $client[:customers].insert_many(
      [
        {
          email: "email@mail.ru",
          password: "something_simple",
          token: "simple-token",
          tokenTTL: Time.now + 10 * 24 * 3600
        },
        {
          email: "email@google.com",
          password: "qUIteHArd_2Crack",
          token: nil,
          tokenTTL: nil
        },
      ])
  end

  def insert_orders
    puts ".. insert orders"
    c = $client[:customers].find().first
    $client[:orders].insert_many(
      [
        {
          created_at: Time.now,
          customer_id: c["_id"],
          raw_data: "huevert"
        },
        {
          created_at: Time.now - 10000,
          customer_id: c["_id"],
          raw_data: "BADABOOOM!"
        }
      ])
  end

end
