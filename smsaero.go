package go_sms_sender

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
)

// В заголовке указываем accept=application/json

type SmsaeroClient struct {
	signature string
	// Канал отправки: FREE SIGN — бесплатное имя, PAY SIGN — платное имя, SERVICE — сервисные сообщения.
	channel  string
	url      string
	template string
}

type SmsaeroMessage struct {
	Numbers []string `json:"numbers"`
	Sign    string   `json:"sign"`
	Text    string   `json:"text"`
	Channel string   `json:"channel"`
}

type SmsaeroResult struct {
	success bool
	data    list.List
	message string
}

func GetSmsaeroClient(email string, apikey string, signature string, template string) (*SmsaeroClient, error) {
	url := fmt.Sprintf("https://%s:%s@gate.smsaero.ru/v2", email, apikey)

	smsaeroClient := &SmsaeroClient{
		signature: signature,
		channel:   "FREE SIGN",
		url:       url,
		template:  template,
	}
	return smsaeroClient, nil
}

func buildSmsaeroMessage(message string, signature string, numbers []string) (*SmsaeroMessage, error) {
	smsaeroMessage := &SmsaeroMessage{
		Numbers: numbers,
		Sign:    signature,
		Text:    message,
		Channel: "FREE SIGN",
	}
	return smsaeroMessage, nil
}

func (c *SmsaeroClient) SendMessage(param map[string]string, numbers ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: msg code")
	}
	if len(numbers) < 1 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	smsContent := fmt.Sprintf(c.template, code)

	smsaeroMessage, _ := buildSmsaeroMessage(smsContent, c.signature, numbers)
	fmt.Println(smsaeroMessage)
	requestBody, err := json.Marshal(smsaeroMessage)
	fmt.Println("Тело запроса")
	fmt.Println(requestBody)
	if err != nil {
		return fmt.Errorf("error creating request body: %w", err)
	}
	fmt.Println(c.url + "/sms/send")
	response, err := http.Post(c.url+"/sms/send", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	fmt.Println(response)

	response.Body.Close()

	return nil
}
