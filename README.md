# Telegram бот помошник для ITDOBRO

![](https://cdn.discordapp.com/attachments/650681889308278785/1138632403074691154/telegram_bot_golang.png)


## Команды

```
/start - Запустить бота
/help - Знакомство с ботом
/commands - Список команд
/contacts - контакты
/ads - Заказать рекламу
/ls - Личное сообщение
/SuggestPost - предложить пост в канал
/treatment - написать обращение (выбор пост, личное обращение или реклама)
```

## Установка

### Зависиости

`go`
`make`

### Deploy

```
git clone https://github.com/Avdushin/gotgbot
cd gotgbot
make
```

или ручками

```
git clone https://github.com/Avdushin/gotgbot
cd gotgbot
go build -o build/tgbot ./cmd/app/main.go
./build/tgbot
```

<p align="center">2023 © <a href="https://github.com/Avdushin" target="_blank">AVDUSHIN</a></p>