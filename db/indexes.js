// mongo localhost:27017/ali24_prod db/indexes.js

print("Add indexes to 'orders' collection");
db.orders.createIndex({customerId: 1});

print("Add indexes to 'currencyRates' collection");
db.currencyRates.createIndex({createdAt: 1});

print("Add indexes to 'messages' collection");
db.messages.createIndex({createdAt: 1});
db.messages.createIndex({createdAt: 1, type: 1});
db.messages.createIndex({createdAt: 1, type: 1, recipientId: 1});
// will be removed in a 200 days
db.messages.createIndex( { "createdAt": 1 }, { expireAfterSeconds: (24 * 3600 * 200) });
