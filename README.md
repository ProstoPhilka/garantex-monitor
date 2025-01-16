# garantex-monitor

Garantex Monitor — это приложение для мониторинга курса валютной пары USDT-RUB. В качестве API используется биржа Garantex.
Взаимодействие с сервисом происходит посредством gRPC методов:
 GetRates() - для получения курса;
 Healthcheck() - для проверки работоспособности сервиса.

## Требования

Перед началом работы убедитесь, что у вас установлены следующие инструменты:

- [Go](https://golang.org/dl/) >= 1.23
- [Docker](https://www.docker.com/get-started) и [Docker Compose](https://docs.docker.com/compose/install/)
- [GolangCI-Lint](https://golangci-lint.run/) для анализа кода

## Команды `Makefile`

Проект содержит `Makefile` для упрощения работы с основными задачами:

### 1. Локальная сборка приложения
Для сборки Go-приложения используйте команду:

```bash
make build
```
### 2. Запуск приложения с использование Docker
Для запуска Go-приложения в Docker используйте команду:

```bash
make run
```
### 3. Сборка Docker образа
Для сборки Go-приложения используйте команду:

```bash
make docker-build
```
### 4. Запуск линтера
Для запуска линтера golangci-lint используйте команду:

```bash
make lint
```