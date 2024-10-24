// mongo localhost:27017/ali24_prod db/indexes.js

print("Add indexes to 'orders' collection");
db.orders.createIndex({customer_id: 1});

print("Add indexes to 'currencyRates' collection");
db.currency_rates.createIndex({created_at: 1});

db.search_translations.createIndex(
   { rus : "text" },
   { default_language: "ru" }
);
