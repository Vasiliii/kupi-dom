package main

import (
	"kupi-dom-admin/dicts"
	"kupi-dom-admin/tg"
	"log"
	"os"
	"time"
)

func main() {
	err := InitDicts()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	dicts.InitPattern()

	err = tg.NewBot("5927726248:AAFQJ0ENrXU_NRX5gjKVdvEyO5bC6RSC3dw", true)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go tg.Instance.HandleUpdates(-1001573131520)

	targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 22, 0, 0, 0, time.Local)

	if time.Now().Hour() >= 22 {
		targetTime = targetTime.AddDate(0, 0, 1)
	}

	duration := time.Until(targetTime)
	timer := time.NewTimer(duration)
	for {
		<-timer.C
		for chatID := range dicts.Chats.Map {
			message, err := tg.Instance.SendMessage(chatID, "Доступ к отправке сообщений будет ограничен с 22:00 до 7:00! Спокойной ночи 😴", true)
			if err != nil {
				log.Println(err)
				continue
			}
			go func() {
				time.Sleep(9 * time.Hour)
				err = tg.Instance.DeleteMessage(message.Chat.ID, message.MessageID)
				if err != nil {
					log.Println(err)
				}
			}()
			time.Sleep(500 * time.Millisecond)
		}
		timer.Reset(24 * time.Hour) // Установите таймер на следующий день
	}

}

func InitDicts() error {
	err := dicts.Load(&dicts.Estate.Map, "./dicts/files/estate.txt")
	if err != nil {
		return err
	}
	err = dicts.Load(&dicts.Swears.Map, "./dicts/files/swears.txt")
	if err != nil {
		return err
	}
	err = dicts.LoadInt64(&dicts.Chats.Map, "./dicts/files/chats.txt")
	if err != nil {
		return err
	}
	err = dicts.LoadInt64(&dicts.Admins.Map, "./dicts/files/admins.txt")
	if err != nil {
		return err
	}
	return nil
}
