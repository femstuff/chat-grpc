# gRPC Chat Service

Этот проект представляет собой gRPC-сервис для общения в чатах. Состоит из двух микросервисов:

    - Auth Service – аутентификация пользователей через JWT.

    - Chat Service – создание чатов, отправка и получение сообщений.

    - CLI-клиент – взаимодействие с сервисами через терминал.

## Функционал

### Auth Service:

* Регистрация и авторизация пользователей.
* Выдача access_token и refresh_token.
* Проверка прав доступа к методам чата.

### Chat Service:

* Создание и удаление чатов.
* Отправка сообщений и их сохранение в базе.
* Подключение к чату с загрузкой истории и стримингом новых сообщений.

## Технологии

* Go – язык программирования.


* gRPC – для обмена данными между сервисами.


* PostgreSQL – база данных для хранения пользователей и сообщений.


* bcrypt – хеширование паролей.


* JWT – авторизация через токены.


* zap – логирование.

## Запуск проекта

### 1. Клонирование репозитория

    git clone git@github.com:femstuff/chat-grpc.git
    cd chat-grpc

### 2. Запуск PostgreSQL

    docker-compose up -d

### 3. Миграции (создание таблиц)

    make migrate-up

    Используй, если появляются ошибки при попытке создания таблиц:
    make install-migrate

### 4. Запуск Auth Service

    go run auth-service/cmd/main.go

### 5. Запуск Chat Service

    go run chat-service/cmd/main.go

### 6. Запуск CLI-клиента

    go run chat-cli/main.go

_Использование CLI
Доступные команды:_

    login <username> <password>         - Войти в систему  
    create_chat <user1,user2,...>        - Создать чат  
    send_message <chat_id> <from> <text> - Отправить сообщение  
    connect <chat_id>                    - Подключиться к чату  
    exit                                  - Выйти

### Пример работы:

> login user1 password123
Успешный вход!

> create_chat user1,user2
Чат создан, ID: 10

> send_message 10 user1 "Привет, как дела?"
Сообщение отправлено

> connect 10
[12:30] user1: Привет, как дела?
