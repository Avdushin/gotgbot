package utils

import (
	"fmt"
	"itdobro/pkg/logger"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var userStatuses = make(map[int]string)

func GetBotToken() string {
	return os.Getenv("BOT_TOKEN")
}

func CreateBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.GetLogger().Panic(err)
	}

	return bot, nil
}

// SetUserStatus устанавливает статус пользователя
func SetUserStatus(userID int, status string) {
	// Здесь вы можете сохранить статус пользователя в базе данных или каким-либо другим способом
	// В данном примере мы будем сохранять статус в переменной
	userStatuses[userID] = status
}

// GetUserStatus получает статус пользователя
func GetUserStatus(userID int) string {
	// Здесь вы можете получить статус пользователя из базы данных или другого источника
	// В данном примере мы будем получать статус из переменной
	status, exists := userStatuses[userID]
	if !exists {
		return ""
	}
	return status
}

// SendToGroup отправляет сообщение в группу
func SendToGroup(bot *tgbotapi.BotAPI, groupChatID int64, message string) {
	msg := tgbotapi.NewMessage(groupChatID, message)
	bot.Send(msg)
}

func GetUserLink(user *tgbotapi.User) string {
	return fmt.Sprintf("@%s", user.UserName)
}
