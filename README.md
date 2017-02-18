## Введение

Данная программа представляет собой http-сервер для записи данных в mongodb поступающих
от браузерных плагинов клиентов по http запросам
Также данный сервис представляет специальные точки входа под паролем /admin
для работы над заказами менеджеров и админов

## Описание

base_url=http://apiwholesale.alibaba24.ru

Мониторинговая точка входа:

    /v1/ping

CustomerAPI точки входа (endpoints)

    /v1/session
    /v1/order
    /v1/orders

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
    email: email,
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

Сервис написан на языке программирования golang ver >= 1.3
Использует базу данных mongodb для сохранения результатов.
Использован веб фреймворк https://github.com/zenazn/goji
Драйвер для mongodb http://godoc.org/labix.org/v2/mgo

Для запусков тестов необходимо, чтобы в PATH был доступен файл `disco.rb` (disco)
Все спецификации для тестов написаны в `spec/api.yml` файле

тесты запускаются:

    go run api_wsc.go
или для уже скомпилированного проекта, находясь в директории проекта

    ./api_wsc
и затем в директории проекта запустить

    disco

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
который расположен в `spec/fixtures/seed.rb`

Для быстрой разработки очень удобно использовать проект [fresh](https://github.com/pilu/fresh)
Он будет компилировать проект после каждого изменения и перестартовать его

    fresh api_wsc.go

## Инсталляция и деплой

Как данный код запустить локально на машине разработчика?
Что нужно установить, чтобы этот код мог быть запущен
Какие действия еще необходимо сделать?
Как запускать тесты?
Ключи запуска, переменные окружения, конфигурационные файлы

## Предварительная подготовка базы данных

После деплоя проекта, в случае обновления файла `db/*.js` необходимо запустить
миграцию базы данных:

    mongo localhost:27017/apiwholesale_prod db/indexes.js

... и другие файлы `*.js`

## Downloads

Где (если есть) можно загрузить версии данного приложения

## FAQs

Вопросы и ответы

## Примеры

Если применимо

## Разработчики и участники

* Батенко Юрий, jurbat@gmail.com, https://github.com/svenyurgensson


## README Дополнительная информация для прочтения

* [Readme Driven Development](http://www.slideshare.net/maetl/readme-driven-development-12783652)


### Mgo
I have yet to write a service that would keep a connection alive and
handle reconnecting. What's the best practice for this in Go and mgo / other engines?

In mgo, what I think is recommended (at least it seems to work for me)
is that you use mgo.Dial to connect to MongoDB on startup.
Then you use Session.Copy on that original session for every time you need
to concurrently use the DB (such as per incoming http request
if you are writing a web service)


http://denis.papathanasiou.org/2012/10/14/go-golang-and-mongodb-using-mgo/comment-page-1/


### Токены и авторизация

Вначале клиент обращается к эндпоинту `/session` с POST запросом - в тебе запроса должны быть
параметры `login` и `password` - соответственно логин и пароль данного клиента.
API в случае если клиент найден возвращает в теле ответа `json` в виде:

    token: 'XXX-XXX.........' // клиентский токен

Далее обращение клиента происходит отправкой этого токена в заголовке запроса `Authorization: token: XXX-XXX.....`
Данный токен имеет срок действия, который может истечь, вследствии чего при очередном запросе от клиента
с просроченным токеном, клиент получает запрошенный ответ и новый токе, который он должен сохранить и использовать в
дальнейшем обмене, пока этот токен также не обновится.
