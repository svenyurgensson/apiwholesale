require "mongo"

test_db =  ENV["API_DB"] ||  "apiwholesale_dev"

Mongo::Logger.logger.level =  3

$client = Mongo::Client.new([ '127.0.0.1:27017' ], database: test_db, connect: :direct)


module Seed extend self

  def orders_coll
    $client[:orders]
  end

  def clear_orders
    puts "... clear orders"
    orders_coll.drop
  end

  def customers_coll
    $client[:customers]
  end

  def clear_customers
    puts "... clear customers"
    customers_coll.drop
  end

  def clear_all
    clear_customers
    clear_orders
  end


  def insert_customers
    puts "... inserting customers"
    maybe_create_collection(:customers)

    customers_coll.insert_many(
      [
        {
          email: "email@mail.ru",
          password: "something_simple",
          token: "simple-token",
          tokenTTL: Time.now + 10 * 24 * 3600,

          created_at: Time.now - 10000,
          updated_at: Time.now - 8000,
          lastSeen_at: Time.now,

          raw_data: {},
        },
        {
          email: "email@google.com",
          password: "qUIteHArd_2Crack",
          token: nil,
          tokenTTL: nil,

          created_at: Time.now - 6000,
          updated_at: Time.now - 3000,
          lastSeen_at: Time.now - 1000 ,

          raw_data: {},
        },
      ])
  end

  def insert_orders
    puts "... insert orders"
    maybe_create_collection(:orders)

    c = customers_coll.find().first
    orders_coll.insert_many(
      [
        {
          created_at: Time.now,
          updated_at: Time.now,
          customer_id: c["_id"],
          raw_data: "huevert"
        },
        {
          created_at: Time.now - 10000,
          updated_at: Time.now - 8000,
          customer_id: c["_id"],
          raw_data: "BADABOOOM!"
        }
      ])
  end

  def fill_docs
    insert_customers
    insert_orders
  end

  def maybe_create_collection(name)
    unless $client.database.collection_names.include?(name.to_s)
      $client[name.to_sym].create
    end
  end


  def order_for(customer = customer_with_orders)
    orders_coll.find(customer_id: customer["_id"]).first
  end

  def orders_for(customer = customer_with_orders)
    orders_coll.find(customer_id: customer["_id"]).to_a
  end

  def customer_with_orders
    customers_coll.find().first
  end


end
