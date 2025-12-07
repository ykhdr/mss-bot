# MSS-Bot (Minecraft Server Status Bot)

Telegram бот для проверки статуса Minecraft серверов.

## Возможности

- 📊 Просмотр статуса сервера (онлайн/оффлайн)
- 👥 Список игроков онлайн
- ⚙️ Настройка сервера для каждого чата
- 💾 Сохранение конфигурации между перезапусками

## Требования

- Go 1.23+
- Docker & Docker Compose (для контейнеризации)

## Установка

### Локальный запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/ykhdr/mss-bot.git
cd mss-bot
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте конфигурацию:
```bash
cp configs/config.kdl configs/config.local.kdl
```

4. Отредактируйте `configs/config.local.kdl` и укажите токен бота:
```kdl
bot {
    token "YOUR_TELEGRAM_BOT_TOKEN"
}

database {
    path "./data/mss-bot.db"
}

minecraft {
    timeout 5
}
```

5. Запустите бота:
```bash
go run ./cmd/bot -config configs/config.local.kdl
```

### Docker

1. Создайте конфигурацию:
```bash
cp configs/config.kdl configs/config.local.kdl
# Отредактируйте configs/config.local.kdl
```

2. Запустите:
```bash
docker-compose up -d
```

## Использование

### Команды

- `/mss` - Открыть главное меню
- `/set <ip:port> <name>` - Настроить сервер (только из меню настроек)
- `/help` - Справка

### Пример

1. Отправьте `/mss` для открытия меню
2. Нажмите "Настройки"
3. Отправьте `/set mc.hypixel.net:25565 Hypixel`
4. Нажмите "Назад", затем "Статус" для просмотра информации

## Разработка

### Структура проекта

```
mss-bot/
├── cmd/bot/              # Точка входа
├── internal/
│   ├── app/              # Инициализация приложения
│   ├── bot/              # Telegram бот и обработчики
│   ├── config/           # Парсинг конфигурации
│   ├── minecraft/        # Клиент для MC серверов
│   ├── service/          # Бизнес-логика
│   └── storage/          # Работа с БД
├── configs/              # Файлы конфигурации
├── Dockerfile
└── docker-compose.yml
```

### Тестирование

```bash
go test ./...
```

### Сборка

```bash
go build -o mss-bot ./cmd/bot
```

## Технологии

- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) - Telegram Bot API
- [minequery](https://github.com/dreamscached/minequery) - Minecraft Server Query
- [kdl-go](https://github.com/ykhdr/kdl-go) - KDL парсер
- [goose](https://github.com/pressly/goose) - Миграции БД
- [squirrel](https://github.com/Masterminds/squirrel) - SQL Builder
- [SQLite](https://www.sqlite.org/) - База данных
- [testify](https://github.com/stretchr/testify) - Тестирование

## Лицензия

MIT
