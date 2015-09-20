// mongo localhost:27017/ali24_prod db/create.js

print("Create 'orders' collection");
db.createCollection("orders");

print("Create 'customers' collection");
db.createCollection("customers");

print("Create 'customersBalance' collection");
db.createCollection("customersBalance");

print("Create 'searchTranslations' collection");
db.createCollection("searchTranslations");

print("Create 'currencyRates' collection");
db.createCollection("currencyRates");

print("Create 'messages' collection");
db.createCollection("messages");

print("Create 'mq_jobs' collection");
db.createCollection("mq_jobs", {capped: true, size: 16384});
