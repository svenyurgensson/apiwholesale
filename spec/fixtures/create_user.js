// mongo localhost:27017/wsapi_test spec/fixtures/create_user.js

print("Drop 'customers' collection");
db.customers.drop();

print("Create 'customers' collection");
db.createCollection("customers");

print("Create 'vasya' customer");
db.customers.insert({
  login: "vasya",
  password: "12345"
});
