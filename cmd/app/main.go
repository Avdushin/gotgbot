package main

import (
	"fmt"
	"os"
	"strings"

	"itdobro/config"
	"itdobro/pkg/logger"
	"itdobro/pkg/templates"
	"itdobro/pkg/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	// Init logger...
	log := logger.GetLogger()
	// Load dotenv variables
	if err := godotenv.Load(); err != nil {
		log.Error("Error loading .env file")
	}
	// Set bot token
	bot, err := utils.CreateBot(utils.GetBotToken())
	if err != nil {
		log.Panic(err)
	}

	// @ Set false to the Production
	dbug := os.Getenv("DEBUG")
	// ? Get Debug mode from .env file
	switch dbug {
	case "true", "TRUE":
		bot.Debug = true
	default:
		bot.Debug = false
	}

	log.Printf("Бот %s запущен!", bot.Self.UserName)

	fmt.Printf("\nDEBUG=%s\n\n", dbug)

	// @ Getting updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // set update's Timeout...

	// TODO: Make refatoring!

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// get messages
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		// Set parse mode to Markdown
		msg.ParseMode = "markdown"

		switch update.Message.Text {
		case "/start":
			msg.Text = templates.Welcome
			msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
		case "Помощь", "/help":
			msg.Text = templates.Welcome
			msg.ReplyMarkup = utils.CreateHelpMenuKeyboard()
		case "Команды", "команды", "/commands":
			msg.Text = templates.Commands
			msg.ReplyMarkup = utils.CreateHelpMenuKeyboard()
		case "Обратиться", "/treatment":
			msg.Text = fmt.Sprintf(`Вы можете обратиться к %s со своим вопросом/предложением`, config.TgUserName)
			msg.ReplyMarkup = utils.CreateSupportMenuKeyboard()
			if msg.Text == "Назад" {
				// msg.Text = "Вы вернулись назад"
				msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
			}
		case "Реклама", "/ads":
			msg.Text = templates.ADS
			utils.SetUserStatus(update.Message.From.ID, "waiting_for_ad_message")
		case "ЛС", "/ls":
			msg.Text = templates.LS
			utils.SetUserStatus(update.Message.From.ID, "waiting_for_ls_message")
		case "Предложить пост", "/SuggestPost":
			utils.SetUserStatus(update.Message.From.ID, "waiting_for_post_title")
			msg.Text = "Введите заголовок поста:"
		case "Контакты", "/contacts":
			utils.SetUserStatus(update.Message.From.ID, "")
			msg.Text = templates.Contacts
			if msg.Text == "Назад" {
				msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
			}
		case "Назад":
			utils.SetUserStatus(update.Message.From.ID, "")
			msg.Text = "Выберете действие..."
			msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
		default:
			status := utils.GetUserStatus(update.Message.From.ID)
			if status == "waiting_for_ad_message" {
				// Обработка рекламного сообщения
				user := update.Message.From
				username := "@" + user.UserName

				adMessage := fmt.Sprintf("📢 Рекламное предложение от %s:\n%s\n\n@%s", username, update.Message.Text, config.TgUserName)

				utils.SendToGroup(bot, config.GetGroupID(), adMessage)

				msg.Text = "Предложение отправлено!"
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
			} else if status == "waiting_for_ls_message" {
				// Обработка личного сообщения
				user := update.Message.From
				username := "@" + user.UserName

				adMessage := fmt.Sprintf("Новое обращение от %s:\n%s", username, update.Message.Text)

				utils.SendToGroup(bot, config.GetGroupID(), adMessage)

				msg.Text = "Обращение отправлено!"
				utils.SetUserStatus(update.Message.From.ID, "")
				msg.ReplyMarkup = utils.CreateMainMenuKeyboard()
			} else {
				utils.SetUserStatus(update.Message.From.ID, "")
				msg.Text = "Выберете действие..."
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
