require "mongo"

test_db =  ENV["API_DB"] ||  "apiwholesale_dev"

Mongo::Logger.logger.level =  3

$client = Mongo::Client.new([ '127.0.0.1:27017' ], database: test_db, connect: :direct)


module Seed extend self

  def currencyRates_coll
    $client[:currencyRates]
  end

  def messages_coll
    $client[:messages]
  end


  def orders_coll
    $client[:orders]
  end

  def clear_orders
    log "... clear orders"
    orders_coll.drop
  end

  def customers_coll
    $client[:customers]
  end

  def clear_customers
    log "... clear customers"
    customers_coll.drop
  end

  def clear_rates
    log "... clear currencyRates"
    currencyRates_coll.drop
  end

  def clear_messages
    log "... clear messages"
    messages_coll.drop
  end

  def clear_all
    clear_customers
    clear_orders
    clear_rates
    clear_messages
  end


  def insert_customers
    log "... inserting customers"
    maybe_create_collection(:customers)

    customers_coll.insert_many(
      [
        {
          email: "email@mail.ru",
          password: "something_simple",
          token: "simple-token",
          tokenTTL: Time.now + 10 * 24 * 3600,

          createdAt: Time.now - 10000,
          updatedAt: Time.now - 8000,
          lastSeenAt: Time.now,

          rawData: {},
        },
        {
          email: "email@google.com",
          password: "qUIteHArd_2Crack",
          token: nil,
          tokenTTL: nil,

          createdAt: Time.now - 6000,
          updatedAt: Time.now - 3000,
          lastSeenAt: Time.now - 1000 ,

          rawData: {},
        },
      ])
  end

  def insert_orders
    log "... insert orders"
    maybe_create_collection(:orders)

    c = customers_coll.find().first
    orders_coll.insert_many(
      [
        {
          createdAt: Time.now,
          updatedAt: Time.now,
          customerId: c["_id"],
          rawData: {some: "complex", data: "to store"}
        },
        {
          createdAt: Time.now - 10000,
          updatedAt: Time.now - 8000,
          customerId: c["_id"],
          rawData: "BADABOOOM!"
        }
      ])
  end


  def insert_currencies
    log "... inserting currencies"
    maybe_create_collection(:currencyRates)
    currencyRates_coll.insert_many(
      [
        { rate: 8.92, createdAt: Time.now - 1000 },
        { rate: 9.02, createdAt: Time.now - 12000 },
      ])
  end

  def insert_messages
    log "... inserting messages"
    maybe_create_collection(:messages)

    custId = customer_with_orders["_id"]

    messages_coll.insert_many(
      [
        {type: 'multicast', message: "Hello, all!", createdAt: Time.now, producerId: nil},
        {type: 'direct', message: "Hello, direct!", createdAt: Time.now - 10000, producerId: nil, recipientId: custId}
      ]
    )

  end


  def fill_docs
    insert_customers
    insert_orders
    insert_currencies
    insert_messages
  end

  def maybe_create_collection(name)
    unless $client.database.collection_names.include?(name.to_s)
      $client[name.to_sym].create
    end
  end

  def order_for(customer = customer_with_orders)
    orders_coll.find(customerId: customer["_id"]).first
  end

  def orders_for(customer = customer_with_orders)
    orders_coll.find(customerId: customer["_id"]).to_a
  end

  def last_order
    orders_coll.find.to_a.last
  end

  def customer_with_orders
    customers_coll.find().first
  end


  def log(txt)
    puts txt if ENV['LOG_SEED']
  end


end
