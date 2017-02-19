1.7.0 2017-02-19
=================
* Новая схема запросов Microsoft Azure Cognitive Translator

1.6.0 2017-02-18
=================
* Поменял схему работы с сессиями/токенами
* Отключил ненужные контроллеры
* Дополнительные переменные в эндпоинте /ping

1.4.0 2016-10-12
=================
* Обновил тесты
* Новая схема для mongodb
* Обновил драйвер mgo

1.3.5 2015-09-22
=================
* Добавили эскейпинг GDB ответа

1.3.4 2015-09-22
=================
* Возвращаем результат дополнительно в GDB кодировке

1.3.4 2015-09-21
=================
* Fix require expired token

1.3 2015-09-20
=================
* Add /search?q=... to return translated ru->zh
* Searches saves to collection `searchTranslations`

1.2.1 2015-09-08
=================
* Creation new order adds `MoMq` message

1.2.0 2015-08-22
=================
* Go 1.5
* Implement cleaning cache for customer on placing new order

1.01.4 2015-08-18
=================
* Implement uuid for orders

1.01.3 2015-07-31
=================
* Now /v1/me work for unauthorized customers


1.01.2 2015-07-29
=================
* Implement Id for /v1/me
* Implement Balance for /v1/me
