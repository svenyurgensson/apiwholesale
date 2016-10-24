// mongo localhost:27017/ali24_prod test_seed.js

db.customers.insert({
    "last_name" : "BATENKO",
    "first_name" : "YURY",
    "middle_name" : "",
    "email" : "jurbat@gmail.com",
    "gender" : "m",
    "password" : "fantom242",
    "token" : "301377d6f91551c76bdeceb505896fd2d31b918e",
    "token_ttl" : null,
    "mobile_phone" : "+79069125317",
    "additional_phone" : "",
    "color" : null,
    "type" : "ordinary",
    "description" : null,
    "birthday_at" : ISODate("1974-02-02T00:00:00.000Z"),
    "skype" : "jurbat",
    "vk_url" : null,
    "facebook_url" : null,
    "web_url" : null,
    "country" : "Russia",
    "post_index" : "660046",
    "region" : "",
    "city" : "Krasnoyarsk",
    "street" : "Sportivnaya st, 190-63",
    "additional_address" : "Sportivnaya st, 190-63",
    "updated_at" : ISODate("2016-10-14T10:55:13.021Z"),
    "created_at" : ISODate("2016-10-14T10:55:13.021Z")
});

db.currency_rates.insert({
  "rate": 10.0,
  "created_at" : ISODate("2016-10-23T10:55:13.021Z")
});
