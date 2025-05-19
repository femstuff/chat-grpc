# gRPC Chat Service

Этот проект представляет собой микросервисную систему для общения в чатах с надёжной доставкой сообщений и email-уведомлений. Архитектура построена на gRPC и включает сагу оркестрации для согласованности между сервисами.

## Компоненты

- **Auth Service** – аутентификация и авторизация пользователей через JWT.
- **Chat Service** – управление чатами и сообщениями.
- **Notification Service** – отправка email-уведомлений при новых сообщениях.
- **Saga Orchestrator** – оркестрация процесса доставки сообщений и уведомлений.
- **CLI-клиент** – терминальный интерфейс для взаимодействия с системой.

## Функциональность

### Auth Service:
- Регистрация и авторизация пользователей.
- Выдача `access_token` и `refresh_token`.
- Проверка прав доступа.
- Получение email-ов участников чата.
- Получение user_id по email (для саги).

### Chat Service:
- Создание и удаление чатов.
- Отправка и хранение сообщений.
- Подключение к чату: история + стриминг новых сообщений.
- Масштабируемость через NATS pub-sub.

### Notification Service:
- Подписка на `chat.*` из NATS.
- Получение email-ов участников.
- Отправка уведомлений.
- Поддержка заглушки `StubEmailSender` для отладки.

### Saga Orchestrator:
- Сценарий: user → message → уведомление.
- Откат, если уведомление не доставлено.
- Подключение к `auth-db` и `users-db`.

## Технологии

- **Go**
- **gRPC**
- **PostgreSQL** (2 базы: users и auth/messages)
- **JWT**
- **bcrypt**
- **NATS**
- **zap**
- **Docker & Docker Compose**
- **Saga Orchestration Pattern**

## Запуск проекта

### 1. Клонирование

```sh
git clone git@github.com:femstuff/chat-grpc.git
cd chat-grpc
```

### 2. Сборка и запуск

```sh
make build
make up
```

### 3. Миграции

```sh
make migrate-up
# При ошибке:
make install-migrate
```

### 4. Запуск по сервисам

```sh
make auth          # Auth Service
make chat          # Chat Service
make notification  # Notification Service
make saga          # Saga Orchestrator
make cli           # CLI
```

## Использование CLI

```sh
login <email> <password>               # Авторизация
create_chat <user1,user2,...>         # Создание чата
send_message <chat_id> <from> <text>  # Отправка (через сагу)
connect <chat_id>                     # Присоединиться к чату
exit                                  # Завершение
```

### Пример

```sh
> login admin@example.com password123
> create_chat user1,user2
Чат создан, ID: 3
> send_message 3 user1 "Привет всем!"
> connect 3
[13:45] user1: Привет всем!
```

## Очистка

```sh
make down     # Остановка
make clean    # Удаление контейнеров, томов и образов
```

## Email-уведомления

- Отправляются при каждом новом сообщении в чат.
- Email-отправка настраивается через `.env`.
- Используется отдельный сервис Notification.
- Возможна подмена отправщика на `StubEmailSender`.

## Особенности

- Две БД: `auth-db` (чаты/связи), `users-db` (пользователи).
- Оркестрация по паттерну Saga.
- Полная поддержка масштабируемости через NATS.