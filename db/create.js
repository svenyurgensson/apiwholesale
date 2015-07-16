// mongo localhost:27017/ali24_prod db/create.js

print("Create 'orders' collection");
db.createCollection("orders");

print("Create 'customers' collection");
db.createCollection("customers");

print("Create 'currencyRates' collection");
db.createCollection("currencyRates");

print("Create 'messages' collection");
db.createCollection("messages");
