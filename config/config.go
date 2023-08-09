package config

import (
	"itdobro/pkg/utils"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	_          = godotenv.Load()
	bot, err   = utils.CreateBot(utils.GetBotToken())
	TgUserName = os.Getenv("TG_USERNAME")
	BotName    = bot.Self.UserName
	AuthorName = "ITDOBRO"
)

var SocialMedia = map[string]string{
	"Youtube":        "https://youtube.com/@itdobr0",
	"Канал Telegram": "https://t.me/itdobr0",
	"Boosty":         "https://boosty.to/itdobro",
	"Discord":        "https://discord.gg/xJ58eVZjxu",
	"GitHub":         "https://github.com/Avdushin",
	"Хабр":           "https://habr.com/ru/users/Avdushin",
	"Pinterest":      "https://ru.pinterest.com/itdobro",
}

var SMLinks string = utils.GetSocialMedia(SocialMedia)

func GetGroupID() int64 {
	groupIDStr := os.Getenv("GROUP_ID")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return groupID
}
