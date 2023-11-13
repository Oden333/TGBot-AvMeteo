package telegram

import (
	"log"

	meteo "github.com/Oden333/TGBot_AvMet/pkg/meteoAPI"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart = "start"
	commandMetar = "metar"
	commandTAF   = "taf"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandMetar:
		return b.handleMetar_latest_Command(message)
	case commandTAF:
		return b.handleTaf_latest_Command(message)
	default:
		return b.handleUnknownCommand(message)
	}
}
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	return nil
}
func (b *Bot) handleMetar_latest_Command(message *tgbotapi.Message) error {
	//TODO: Запросить у пользвателя ICAO идентефикатор (и мб другие параметры запроса) и передавать это всё в функцию
	resp, err := meteo.MeteoRequest()
	if err != nil {
		b.handleError(message.Chat.ID, err)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, resp)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleTaf_latest_Command(message *tgbotapi.Message) error {
	//TODO: Запросить у пользвателя ICAO идентефикатор (и мб другие параметры запроса) и передавать это всё в функцию
	resp, err := meteo.TAFRequest()
	if err != nil {
		b.handleError(message.Chat.ID, err)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, resp)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда")
	_, err := b.bot.Send(msg)
	return err
}
