package main

import (
	"fmt"
	"strings"

	"itdobro/config"
	"itdobro/pkg/logger"
	"itdobro/pkg/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	// Init logger...
	log := logger.GetLogger()
	// Load dotenv variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// Set bot token
	botToken := utils.GetBotToken()
	bot, err := utils.CreateBot(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Info(fmt.Sprintf("Бот %s запущен!", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Text {
		case "/start":
			msg.Text = "Привет! Я бот с кнопками. Выберите действие:"
			msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
		case "Обратиться":
			msg.Text = "Выберите действие:"
			msg.ReplyMarkup = utils.CreateSupportMenuKeyboard()
		case "Реклама":
			msg.Text = "Введите текст рекламного сообщения:"
			utils.SetUserStatus(update.Message.From.ID, "waiting_for_ad_message")
		case "Предложить пост":
			utils.SetUserStatus(update.Message.From.ID, "waiting_for_post_title")
			msg.Text = "Введите заголовок поста:"
		default:
			status := utils.GetUserStatus(update.Message.From.ID)
			if status == "waiting_for_ad_message" {
				// Обработка рекламного сообщения
				user := update.Message.From
				username := "@" + user.UserName
				userLink := fmt.Sprintf("Пользователь %s", username)

				adMessage := fmt.Sprintf("📢 Реклама от %s:\n%s", userLink, update.Message.Text)

				utils.SendToGroup(bot, config.GetGroupID(), adMessage)

				msg.Text = "Рекламное сообщение отправлено в группу!"
				utils.SetUserStatus(update.Message.From.ID, "")
				msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
			} else if status == "waiting_for_post_title" {
				utils.SetUserStatus(update.Message.From.ID, "waiting_for_post_text:"+update.Message.Text)
				msg.Text = "Введите текст поста:"
			} else if strings.HasPrefix(status, "waiting_for_post_text:") {
				utils.SetUserStatus(update.Message.From.ID, "waiting_for_media_links:"+status[len("waiting_for_post_text:"):])
				msg.Text = "Введите ссылки на медиафайлы через запятую (URL1, URL2, ...)." +
					"\nВы также можете нажать 'Пропустить', если не хотите добавлять медиа материалы."
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Пропустить"),
					),
				)
			} else if strings.HasPrefix(status, "waiting_for_media_links:") {
				parts := strings.Split(status, ":")
				if len(parts) != 2 {
					continue
				}
				postTitle := parts[1]

				mediaLinks := update.Message.Text
				if mediaLinks != "Пропустить" {
					mediaLinks = strings.TrimSpace(mediaLinks)

					media := []interface{}{}

					for _, link := range strings.Split(mediaLinks, ",") {
						trimmedLink := strings.TrimSpace(link)
						if isPhotoLink(trimmedLink) {
							media = append(media, tgbotapi.NewInputMediaPhoto(trimmedLink))
						} else if isVideoLink(trimmedLink) {
							video := tgbotapi.NewInputMediaVideo(trimmedLink)
							media = append(media, video)
						}
					}

					if len(media) > 0 {
						mediaMsg := tgbotapi.NewMediaGroup(int64(update.Message.From.ID), media)
						mediaMsg.ChatID = int64(config.GetGroupID())
						bot.Send(mediaMsg)
					}
				}

				groupMessage := fmt.Sprintf("Предложенный пост от %s:\n\nЗаголовок: %s\nТекст: %s",
					utils.GetUserLink(update.Message.From), postTitle, status[len("waiting_for_post_text:"):])

				if mediaLinks != "Пропустить" {
					mediaMessage := "Медиа:\n" + strings.Join(strings.Split(mediaLinks, ","), "\n")
					groupMessage = fmt.Sprintf("%s\n\n%s", groupMessage, mediaMessage)
				}

				utils.SendToGroup(bot, config.GetGroupID(), groupMessage)

				msg.Text = "Пост отправлен в предложку!"
				utils.SetUserStatus(update.Message.From.ID, "")
				msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
			} else {
				msg.Text = "Выберите действие:"
				msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
			}
		}

		bot.Send(msg)
	}
}

func isPhotoLink(link string) bool {
	return strings.HasSuffix(link, ".jpg") || strings.HasSuffix(link, ".png")
}

func isVideoLink(link string) bool {
	return strings.HasSuffix(link, ".mp4") || strings.HasSuffix(link, ".avi")
}
