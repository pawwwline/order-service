![Coverage](https://img.shields.io/badge/dynamic/json?color=success&label=Coverage&query=%24.coverage&suffix=%25&url=https%3A%2F%2F<ваш-юзер>.github.io%2F<название-репо>%2Fcoverage.json&logo=go&style=for-the-badge)

# Order processing service

> Микросервис на Go для обработки заказов с интеграцией Kafka, PostgreSQL и внутреннего кэша. Получает заказы из очереди сообщений, сохраняет их в базу данных и кэширует для ускоренного доступа. Поддерживает повторные попытки обработки сообщений (retry) и отправку неуспешных сообщений в очередь ошибок (DLQ)

## Требования
Убедитесь, что у вас установлен docker и docker compose, go v1.24.5
```bash
docker --version
docker compose version
go version
```

## Установка и запуск

1. Клонировать репозиторий
```bash
git clone https://github.com/pawwwline/order-service
cd order-service
```

2. Создать .env на основе .env.example (при необходимости поставить нужные значения)
```bash
cp env.example .env
```
> **Note:** Доступные APP_ENV `local` `test` `dev` `prod`

3. Запуск сервиса

```bash
make run
```

4. Интеграционные тесты
```bash
make integration
```
> **Note:** Останавливает рабочие контейнеры при запуске тестов; для продолжения работы выполните `make run`.

5. Юнит тесты
```bash
make test
```

## Демонстрация работы с Kafka

1. Создать JSON-файл с заказом в одну строку, например order1.json
2. Положить в директорию ./demo-producer/msgs
3. Запустить продюсера

```bash
make demo
```
### Документация

> **Note:** Документация API доступна на `/swagger/index.html` после запуска сервиса.

### Интерфейс
> **Note:** Интерфейс проекта доступен по пути `/` после запуска сервиса. Статические файлы находятся в папке public.



