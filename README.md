## Введение

Данная программа представляет собой http-сервер для записи данных в mongodb поступающих
от браузерных плагинов клиентов по http запросам
Также данный сервис представляет специальные точки входа под паролем /admin
для работы над заказами менеджеров и админов

## Описание

base_url=http://apiwholesale.alibaba24.ru

CustomerAPI точки входа (endpoints)

    base_url + "/v1/session"
    base_url + "/v1/order"
    base_url + "/v1/orders"
    base_url + "/v1/ping"

Операции:

Create session POST /v1/session
Delete session DELETE /v1/session

    LIST /v1/orders
    GET /v1/order/:id
    CREATE POST /v1/order

### SIGNIN
Создает сессию, при успехе возвращает временный ключ
Успех означает то, что данный пользователь известен системе и ему разрешен вход.

Временный ключ используется при дальнейшей работе с API,
каждый запрос подписывается этим ключем


Headers:
   Conetent-Type: application/json

POST /v1/session
{
login: email,
password: password
}

-->
Headers:
{
Content-Type: application/json
X-Token-TTL: 3600 (1 hour)
}

Body:
{
token: "...long..token string..."
}



### Orders

Headers

Accept: application/json
GET /v1/orders
возвращает список моих заказов и товаров в этом заказе

[
{
.. order object ..
},
{
.. order object ..
}
]

Headers
Content-Type: application/json

POST /v1/orders
Создает новую запись в таблице заказов, перед сохранением валидирует структуру входных данных.
Body
{
 .. object to save ..
}

### Ping endpoint

Accept: application/json
GET /v1/ping

Предназначен для контроля и мониторинга работы сервиса, получения статистической
информации по состоянии, запросам и ошибкам.

{
status: "OK",
memory: "",
....
}


## Быстрый старт работы над этим проектом

Сервис написан на языке программировани golang ver >= 1.3
Использует базу данных mongodb для сохранения результатов.
Использован веб фреймворк https://github.com/zenazn/goji
Драйвер для mongodb http://godoc.org/labix.org/v2/mgo

Для запусков тестов необходимо, чтобы в PATH был доступен файл disco.rb (disco)
тесты запускаются:

%application root%> go run api_wsc.go

или для уже скомпилированного проекта

%application root%> ./api_wsc
%application root%> disco

результатом выполнения тестов будет файл doc/index.html который содержит примеры обращения к
API и ответы сервиса

ПО умолчанию тесты запускаются в окружении development, также можно запускать в режиме
production, test
Для этого необходимо указать базу данных нужного окружения перед запуском disco

API_DB=apiwholesale_test disco

Сам сервис API тоже необходимо запустить в нужном окружении:

./api_wsc -e test

В процессе отладки и изменения спецификаций API а также для добавления новых тестов необходимо
вносить изменения в файл спецификации spec/api.yml поскольку он используется при работе
программы disco

Возможно потребуется изучить и внести изменения в скрипт для первоначального засевания базы данных
который расположен в spec/fixtures/seed.rb


## Инсталляция и деплой

Как данный код запустить локально на машине разработчика?
Что нужно установить, чтобы этот код мог быть запущен
Какие действия еще необходимо сделать?
Как запускать тесты?
Ключи запуска, переменные окружения, конфигурационные файлы

## Downloads

Где (если есть) можно загрузить версии данного приложения

## FAQs

Вопросы и ответы

## Примеры

Если применимо

## Разработчики и участники

Батенко Юрий, jurbat@gmail.com, https://github.com/svenyurgensson


## README presentation
[Readme Driven Development](http://www.slideshare.net/maetl/readme-driven-development-12783652)


### mgo
I have yet to write a service that would keep a connection alive and
handle reconnecting. What's the best practice for this in Go and mgo / other engines?

In mgo, what I think is recommended (at least it seems to work for me)
is that you use mgo.Dial to connect to MongoDB on startup.
Then you use Session.Copy on that original session for every time you need
to concurrently use the DB (such as per incoming http request
if you are writing a web service)


http://denis.papathanasiou.org/2012/10/14/go-golang-and-mongodb-using-mgo/comment-page-1/
