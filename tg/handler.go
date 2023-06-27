package tg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"kupi-dom-admin/dicts"
	"log"
	"math/rand"
	"strings"
	"time"
)

const msgTTL = 30 * time.Second

var (
	Instance      instance
	State         = make(map[int64]int)
	msgs          = make(map[int64][]*tgbotapi.Message)
	msgsToPublish = make(map[string]ToPublish)
	codeToAdd     = make(map[string]int64)
)

type ToPublish struct {
	Message    *tgbotapi.MessageConfig
	MediaGroup *tgbotapi.MediaGroupConfig
	Photo      *tgbotapi.PhotoConfig
}

type instance struct {
	Bot *tgbotapi.BotAPI
	Log *log.Logger
}

func NewBot(token string, debug bool) (err error) {
	Instance.Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}
	Instance.Bot.Debug = debug
	return nil
}

//for chatID := range dicts.Chats.Map {
//	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID("AgACAgIAAxkBAAIH82RwylOIHHL9Dg5BL3ztTL14vIfqAAKYwzEbbECJS8MEqQ0ZnotfAQADAgADeQADLwQ"))
//	msg.Caption = fmt.Sprintf("Мечтаете проводить больше времени у моря? Для Вас есть отличное предложение!\n\n" +
//		"Апартаменты премиум класса в квартале «Прайм» курортного комплекса «Прибрежный» на западном берегу Крыма. 1 шаг до пляжа, подогреваемый бассейн, SPA, медцентр, парк, набережная. Старт от 3,6 млн !!!\n\n" +
//		"Получи лучшее предложение доходной недвижимости одним из первых!!!\n\n" +
//		"https://кк-прибрежный.рф")
//	i.Bot.Send(msg)
//	time.Sleep(time.Millisecond * 500)
//}
//os.Exit(0)

func (i *instance) HandleUpdates(adminChatId int64) {
	//for chatID := range dicts.Chats.Map {
	//	msg := tgbotapi.NewVideo(chatID, tgbotapi.FileID("BAACAgIAAxkBAAIH-GR3QQABBX49mnjqryiQREcHTceR-AAC1y8AAnRQuEudttxtr1WZ7S8E"))
	//	msg.Caption = "*Как получить за 2 дня больше опыта, чем за 10 лет работы в недвижимости?*\n\n" +
	//		"3\\-4 июня в Москве *пройдет самая крупная конференция по недвижимости* от Private Money Forum и инвест\\-клуба «Деньги»\\.\n\n" +
	//		"Более 600 инвесторов и предпринимателей, с которыми вы сможете познакомиться, перенять опыт и первыми получить информацию о закрытых прибыльных проектах\\. \n\n" +
	//		"Вот лишь малая часть того, что будет раскрыто на конференции:\n\n" +
	//		"▫️ Как зарабатывать на недвижимости, не покупая ее; \n" +
	//		"▫️ Факапы инвестирования в 22\\-ом и как перевести их в плюсы;\n" +
	//		"▫️ Самые стабильные проекты в недвижимости \n\n" +
	//		"👉 *Получить подробную программу и забронировать билет – [money\\-event](https://money-event.ru/3-4june?utm_source=telegram&utm_medium=ads&utm_campaign=Investportal&utm_content=3-4-june-2023)*\n\n" +
	//		"Хотите узнать о других инвестиционных мероприятиях? — [подпишитесь на канал](https://t.me/+u9dxX3y8nlw4ZDgy), чтобы быть в числе первых\\."
	//	msg.ParseMode = "MarkdownV2"
	//	i.Bot.Send(msg)
	//	time.Sleep(time.Millisecond * 500)
	//}
	//os.Exit(0)
	u := tgbotapi.NewUpdate(0)
	updates := i.Bot.GetUpdatesChan(u)
	for update := range updates {
		fmt.Println(dicts.Estate.Check(update.Message.Text))
		//go i.haldleUpdates(update, adminChatId)
	}
}

func (i instance) haldleUpdates(update tgbotapi.Update, adminChatId int64) {
	if update.CallbackQuery != nil {
		i.handleCallbackQuery(update)
		return
	}
	if update.Message == nil {
		return
	}

	if update.Message.Chat.IsPrivate() {
		i.handlePrivateMessage(update)
		return
	}

	if _, ok := codeToAdd[update.Message.Text]; ok {
		if dicts.Chats.Check(update.Message.Chat.ID) {
			err := i.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
			if err != nil {
				return
			}
			msg := tgbotapi.NewMessage(codeToAdd[update.Message.Text], "Данная группа уже находится в списке рассылки")
			_, err = i.Bot.Send(msg)
			if err != nil {
				return
			}
			delete(codeToAdd, update.Message.Text)
			return
		}

		err := AddChat(update.Message.Chat.ID)
		if err != nil {
			return
		}
		err = i.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		if err != nil {
			return
		}
		msg := tgbotapi.NewMessage(codeToAdd[update.Message.Text], "Группа была успешно добавлена в список рассылки")
		_, err = i.Bot.Send(msg)
		if err != nil {
			return
		}
		delete(codeToAdd, update.Message.Text)
		return
	}

	if i.IsAdmin(update.Message.Chat.ID, update.Message.From.ID) {
		return
	}

	if !dicts.Chats.Check(update.Message.Chat.ID) {
		return
	}

	if !i.IsGroupMember(-1001843717362, update.Message.From.ID) {
		err := i.ForwardMessage(adminChatId, update.Message.Chat.ID, update.Message.MessageID)
		if err != nil {
			log.Println(err)
		}

		msg, err := i.SendInlineKeyboard(update.Message.Chat.ID,
			fmt.Sprintf("%s %s для того, что бы опубликовать объявление в группе Вам необходимо подписаться на наш канал 👇️",
				update.Message.From.FirstName, update.Message.From.LastName),
			tgbotapi.NewInlineKeyboardButtonURL("Инвестиции без границ", "https://t.me/investbezgranic"))
		if err == nil {
			go func() {
				time.Sleep(msgTTL)
				err = i.DeleteMessage(msg.Chat.ID, msg.MessageID)
				if err != nil {
					log.Println(err)
				}
			}()
		}
		err = i.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if dicts.Estate.Check(update.Message.Text + update.Message.Caption) {
		err := i.ForwardMessage(adminChatId, update.Message.Chat.ID, update.Message.MessageID)
		if err != nil {
			log.Println(err)
		}

		msg, err := i.SendInlineKeyboard(update.Message.Chat.ID,
			fmt.Sprintf("%s Вы можете разместить объявление по недвижимости только через публикацию объявления в нашем боте \"МойДом\"",
				update.Message.From.FirstName),
			tgbotapi.NewInlineKeyboardButtonURL("МойДом", "https://t.me/MoyDom_Rielty_bot"))
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			go func() {
				time.Sleep(msgTTL)
				err = i.DeleteMessage(msg.Chat.ID, msg.MessageID)
				if err != nil {
					log.Println(err)
				}
			}()
		}

		err = i.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if dicts.Swears.Check(update.Message.Text) {
		err := i.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if IsNightTime(update.Message.Time()) {
		err := i.ForwardMessage(adminChatId, update.Message.Chat.ID, update.Message.MessageID)
		if err != nil {
			log.Println(err)
		}

		msg, err := i.SendMessage(update.Message.Chat.ID, "Публиковать объявления можно только с 7:00 до 22:00", true)
		if err != nil {
			log.Println(err)
		}
		if err == nil {
			go func() {
				time.Sleep(msgTTL)
				err = i.DeleteMessage(msg.Chat.ID, msg.MessageID)
				if err != nil {
					log.Println(err)
				}
			}()
		}

		err = i.DeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		if err != nil {
			log.Println(err)
		}
		return
	}
}

func (i *instance) DeleteMessage(chatID int64, messageID int) error {
	_, err := i.Bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
	return err
}

func (i *instance) ForwardMessage(chatID int64, fromChatID int64, messageID int) error {
	_, err := i.Bot.Request(tgbotapi.NewForward(chatID, fromChatID, messageID))
	return err
}

func (i *instance) SendMessage(chatID int64, text string, disableNotification bool) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.DisableNotification = disableNotification
	res, err := i.Bot.Send(msg)
	return res, err
}

func (i *instance) SendInlineKeyboard(chatID int64, text string, buttons ...tgbotapi.InlineKeyboardButton) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
	return i.Bot.Send(msg)
}

func (i *instance) RestrictUserUntilSevenAM(chatID int64, userID int64) error {
	currentTime := time.Now()
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	if currentTime.Hour() >= 7 {
		currentTime = currentTime.AddDate(0, 0, 1)
	}
	until := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
		7, 0, 0, 0, location).Unix()

	_, err = i.Bot.Send(tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: until,
	})
	return err
}

func (i *instance) IsGroupMember(chatID int64, userID int64) bool {
	member, err := i.Bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		}})
	if err != nil {
		return false
	}
	return member.IsMember || member.Status == "member" || member.IsCreator() || member.IsAdministrator()
}

func (i instance) IsAdmin(chatID, userID int64) bool {
	member, err := i.Bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		fmt.Printf("\n\n\n\n%s\n\n\n\n", err.Error())
	}
	return err == nil && (member.IsCreator() || member.IsAdministrator())
}

func (i instance) handlePrivateMessage(u tgbotapi.Update) {
	if !dicts.Admins.Check(u.Message.From.ID) {
		return
	}

	if u.Message.Command() == "start" {
		replyKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Опубликовать объявление")),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Добавить группу"),
				tgbotapi.NewKeyboardButton("Удалить группу"),
			),
		)
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Привет, "+u.Message.From.FirstName)
		msg.ReplyMarkup = replyKeyboard
		_, err := i.Bot.Send(msg)
		if err != nil {
			return
		}
	}
	if u.Message.Text == "Опубликовать объявление" {
		State[u.Message.From.ID] = 1
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Перешлите объявление")
		_, err := i.Bot.Send(msg)
		if err != nil {
			return
		}
		return
	} else if u.Message.Text == "Добавить группу" {
		State[u.Message.From.ID] = 2
		code := generateCode()
		codeToAdd[code] = u.Message.Chat.ID
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, fmt.Sprintf("Сначала добавьте бота в группу, затем напишите код `%s` в чате группы", code))
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		_, err := i.Bot.Send(msg)
		if err != nil {
			return
		}
		return
	} else if u.Message.Text == "Удалить группу" {
		State[u.Message.From.ID] = 3
		return
	}
	if State[u.Message.From.ID] == 1 {
		if u.Message.Photo != nil {
			if msgs[u.Message.Chat.ID] == nil {
				msgs[u.Message.Chat.ID] = []*tgbotapi.Message{}
				go i.startHandle(u.Message.Chat.ID)
			}
			msgs[u.Message.Chat.ID] = append(msgs[u.Message.Chat.ID], u.Message)
			return
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
		msg.Entities = u.Message.Entities
		msg.DisableWebPagePreview = true

		_, err := i.Bot.Send(msg)
		if err != nil {
			log.Println(err)
		}

		uid := uuid.New().String()
		msgsToPublish[uid] = ToPublish{Message: &msg}
		sendConifrmationMessage(i.Bot, u.Message.Chat.ID, uid)
	}
	State[u.Message.From.ID] = 0
	if msgs[u.Message.Chat.ID] != nil {
		msgs[u.Message.Chat.ID] = nil
	}

}

func (i instance) startHandle(chatID int64) {
	time.Sleep(2 * time.Second)
	msg := msgs[chatID]
	if len(msg) == 1 {
		mesg := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(msg[0].Photo[len(msg[0].Photo)-1].FileID))
		mesg.Caption = msg[0].Caption
		mesg.CaptionEntities = msg[0].CaptionEntities

		_, err := i.Bot.Send(mesg)
		if err != nil {
			log.Println(err)
		}

		uid := uuid.New().String()
		msgsToPublish[uid] = ToPublish{Photo: &mesg}
		sendConifrmationMessage(i.Bot, chatID, uid)
	} else if len(msgs[chatID]) > 1 {
		var photos []interface{}
		for _, photo := range msgs[chatID] {
			photos = append(photos, tgbotapi.InputMediaPhoto{BaseInputMedia: tgbotapi.BaseInputMedia{
				Type:            "photo",
				Media:           tgbotapi.FileID(photo.Photo[len(photo.Photo)-1].FileID),
				Caption:         photo.Caption,
				CaptionEntities: photo.CaptionEntities,
			}})
		}

		mesg := tgbotapi.NewMediaGroup(chatID, photos)

		_, err := i.Bot.Send(mesg)
		if err != nil {
			log.Println(err)
		}
		uid := uuid.New().String()
		msgsToPublish[uid] = ToPublish{MediaGroup: &mesg}
		sendConifrmationMessage(i.Bot, chatID, uid)
	}
	if msgs[chatID] != nil {
		msgs[chatID] = nil
	}
	State[chatID] = 0
}

func (i instance) handleCallbackQuery(u tgbotapi.Update) {
	if u.CallbackQuery != nil {
		callback := strings.Split(u.CallbackQuery.Data, ":")
		if len(callback) != 2 {
			return
		}
		if callback[0] == "delete" {
			delete(msgsToPublish, callback[1])
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Публикация отменена")
			_, err := i.Bot.Send(msg)
			if err != nil {
				return
			}
			State[u.CallbackQuery.From.ID] = 0
			return
		}
		if callback[0] == "publish" {
			if msgsToPublish[callback[1]].MediaGroup == nil && msgsToPublish[callback[1]].Photo == nil && msgsToPublish[callback[1]].Message == nil {
				return
			}
			if msgsToPublish[callback[1]].Message != nil {
				msg := msgsToPublish[callback[1]].Message
				for group := range dicts.Chats.Map {
					msg.ChatID = group
					_, err := i.Bot.Send(msg)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			} else if msgsToPublish[callback[1]].Photo != nil {
				msg := msgsToPublish[callback[1]].Photo
				for group := range dicts.Chats.Map {
					msg.ChatID = group
					_, err := i.Bot.Send(msg)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			} else if msgsToPublish[callback[1]].MediaGroup != nil {
				msg := msgsToPublish[callback[1]].MediaGroup
				for group := range dicts.Chats.Map {
					msg.ChatID = group
					_, err := i.Bot.Send(msg)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
			State[u.CallbackQuery.From.ID] = 0
			delete(msgsToPublish, callback[1])
			msg := tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Публикация успешно завершена")
			_, err := i.Bot.Send(msg)
			if err != nil {
				return
			}
			return
		}
	}
}

func sendConifrmationMessage(bot *tgbotapi.BotAPI, chatID int64, messageID string) {
	msg := tgbotapi.NewMessage(chatID, "Проверьте правильность объявления и нажмите кнопку \"Опубликовать\"")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Опубликовать во всех группах", fmt.Sprintf("publish:%s", messageID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить", fmt.Sprintf("delete:%s", messageID))),
	)
	bot.Send(msg)
}

func IsNightTime(time time.Time) bool {
	return time.Hour() >= 22 || time.Hour() < 7
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func AddChat(chatID int64) error {
	err := dicts.AddInt64(&dicts.Chats.Map, chatID, "./dicts/files/chats.txt")
	if err != nil {
		return err
	}
	return nil
}

func DeleteChat(chatID int64) error {
	err := dicts.RemoveInt64(&dicts.Chats.Map, chatID, "./dicts/files/chats.txt")
	if err != nil {
		return err
	}
	return nil
}
