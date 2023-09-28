package repositories

import (
	"bytes"
	"encoding/json"
	"follooow-be/configs"
	"net/http"
)

type PayloadSendMessage struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

func TelegramSendMessage(text string) error {
	payload := PayloadSendMessage{
		Text:   text,
		ChatId: configs.TELEGRAM_FOLLOOOW_CHANNEL,
	}
	jsonPayload, errJson := json.Marshal(payload)
	if errJson != nil {
		return errJson
	}
	url := "https://api.telegram.org/bot" + configs.TELEGRAM_FOLLOOOW_TOKEN + "/sendMessage"
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	return err
}
