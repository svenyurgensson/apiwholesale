# coding: utf-8
require 'mongo'

db = ENV['API_DB'] || 'ali24_test'

Mongo::Logger.logger.level = 3
$client = Mongo::Client.new(['127.0.0.1:27017'], database: db, connect: :direct)

module Seed
  extend Seed

  def currency_rates_coll
    $client[:currency_rates]
  end

  def clear_rates
    log '... clear currencyRates'
    currency_rates_coll.drop
  end

  def messages_coll
    $client[:messages]
  end

  def clear_messages
    log "... clear messages"
    messages_coll.drop
  end

  def orders_coll
    $client[:raw_orders]
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

  def search_coll
    $client[:search_translations]
  end

  def clear_search
    log "... clear search translations"
    search_coll.drop
  end

  def clear_all
    clear_customers
    clear_orders
    clear_rates
    clear_messages
    clear_search
  end


  def insert_customers
    log "... inserting customers"
    maybe_create_collection(:customers)

    customers_coll.insert_many(
      [
        {
          last_name: 'Мишина',
          first_name: 'Евдокия',
          middle_name: 'Сергеевна',
          email: 'email@mail.ru',
          gender: 'f',
          color: 'FFEEAA',
          type: 'ordinary',
          description: '',

          mobile_phone: '79079129912',
          additional_phone: '79079129900',

          facebook_url: 'https://facebook.com/evmishina',
          vk_url: 'https://vk.com/evmishina',
          skype: 'evmishina',

          country: 'Россия',
          post_index: '660069',
          region: 'Красноярский край',
          city: 'Красноярск',
          street: 'Красноярский рабочий 66-12',
          additional_address: '',

          birthday_at: Date.parse('2000-1-13'),

          created_at: Time.now - 10000,
          updated_at: Time.now - 8000,
          last_seen_at: Time.now,

          balance_total: rand(10000),

          password: 'something_simple',
          token: 'simple-token',
          token_ttl: Time.now + 10 * 24 * 3600
        },
        {
          last_name: 'Велетень',
          first_name: 'Евгения',
          middle_name: 'Сергеевна',
          email: 'email@google.com',
          gender: 'f',
          color: 'FFEEAA',
          type: 'ordinary',
          description: '',

          mobile_phone: '79079129912',
          additional_phone: '79079129900',

          facebook_url: 'https://facebook.com/eveleten',
          vk_url: 'https://vk.com/eveleten',
          skype: 'eveleten',

          country: 'Россия',
          post_index: '660122',
          region: 'Красноярский край',
          city: 'Красноярск',
          street: 'Академика Павлова 66-12',
          additional_address: '',

          birthday_at: Date.parse('2000-1-13'),

          created_at: Time.now - 10000,
          updated_at: Time.now - 8000,
          last_seen_at: nil,

          balance_total: rand(10000),

          password: 'qUIteHArd_2Crack',
          token: nil,
          token_ttl: nil
        }
      ])
  end

  def insert_orders
    log '... insert orders'
    maybe_create_collection(:raw_orders)

    c = customers_coll.find().first
    orders_coll.insert_many([
                              {
                                created_at: Time.now,
                                updated_at: Time.now,
                                customer_id: c['_id'],
                                rawData: {some: 'complex', data: 'to store'}
                              },
                              {
                                created_at: Time.now - 10000,
                                updated_at: Time.now - 8000,
                                customer_id: c['_id'],
                                rawData: 'BADABOOOM!'
                              }
                            ])
  end

  def insert_currencies
    log '... inserting currencies'
    maybe_create_collection(:currency_rates)
    currency_rates_coll.insert_many([
                                     { rate: 8.92, created_at: Time.now - 1000 },
                                     { rate: 9.02, created_at: Time.now - 12000 }
                                    ])
  end

  def insert_messages
    log '... inserting messages'
    maybe_create_collection(:messages)

    cust_id = customer_with_orders['_id']
    messages_coll.insert_many([
        { type: 'multicast',
          message: 'Hello, all!',
          createdAt: Time.now,
          producerId: nil },
        { type: 'direct',
          message: 'Hello, direct!',
          createdAt: Time.now - 10000,
          producerId: nil,
          recipientId: cust_id }
      ])
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
    orders_coll.find(customer_id: customer['_id']).first
  end

  def orders_for(customer = customer_with_orders)
    orders_coll.find(customer_id: customer['_id']).to_a
  end

  def last_order
    orders_coll.find.to_a.last
  end

  def customer_with_orders
    customers_coll.find.first
  end

  def log(txt)
    puts txt if ENV['LOG_SEED']
  end
end

Seed.log("Connected to: #{db}")
