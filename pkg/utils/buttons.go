package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CreateMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keys := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Обратиться"),
		tgbotapi.NewKeyboardButton("Помощь"),
	}
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(keys...),
	)
}

func CreateSupportMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Реклама"),
			tgbotapi.NewKeyboardButton("Предложить пост"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ЛС"),
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	return keyboard
}

func CreateHelpMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Команды"),
			tgbotapi.NewKeyboardButton("Контакты"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	return keyboard
}
