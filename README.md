## Введение

Данная программа представляет собой http-сервер для записи данных в mongodb поступающих
от браузерных плагинов клиентов по http запросам

## Описание

base_url=http://apiwholesale.alibaba24.ru

CustomerAPI точки входа (endpoints)

    base_url + "/v1/session"
    base_url + "/v1/order"

Операции:

Create session POST /v1/session
Delete session DELETE /v1/session

    LIST /v1/orders
    GET /v1/order/:id
    CREATE POST /v1/order

## SIGNIN
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



## Orders

Headers

Accept: application/json
GET /v1/orders
возвращает список моих заказов и товаров в этом заказе

POST /v1/orders

[
{
.. order object ..
},
{
.. order object ..
}
]



## Быстрый старт работы над этим проектом

Вкратце описать шаги или примеры команд для запуска

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

Имя+Фамилия, емейл, гитхаб ссылка


## README presentation
[Readme Driven Development](http://www.slideshare.net/maetl/readme-driven-development-12783652)
