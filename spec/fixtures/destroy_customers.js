// mongo localhost:27017/wsapi_test spec/fixtures/destroy_customers.js

print("Drop 'customers' collection");
db.customers.drop();
