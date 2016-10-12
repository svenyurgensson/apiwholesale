// mongo localhost:27017/ali24_prod db/create.js

print("Create 'orders' collection");
db.createCollection("raw_orders");

print("Create 'customers' collection");
db.createCollection("customers");

print("Create 'customersBalance' collection");
db.createCollection("customers_balance");

print("Create 'searchTranslations' collection");
db.createCollection("search_translations");

print("Create 'currencyRates' collection");
db.createCollection("currency_rates");

print("Create 'messages' collection");
db.createCollection("messages");
