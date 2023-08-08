package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CreateMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keys := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Написать"),
		tgbotapi.NewKeyboardButton("Информация"),
		tgbotapi.NewKeyboardButton("Предложить пост"),
	}
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(keys...),
	)
}

// CreateSupportMenuKeyboard создает клавиатуру для меню обращения
func CreateSupportMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	// Создайте и верните клавиатуру для меню обращения
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Реклама"),
			tgbotapi.NewKeyboardButton("Предложить пост"),
		),
	)
	return keyboard
}
