package go_sms_sender

import (
	"container/list"
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
	numbers []string
	sign    string
	text    string
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
		channel:   "PAY SIGN",
		url:       url,
		template:  template,
	}
	return smsaeroClient, nil
}

func buildSmsaeroMessage(message string, signature string, numbers []string) (*SmsaeroMessage, error) {
	smsaeroMessage := &SmsaeroMessage{
		numbers: numbers,
		sign:    signature,
		text:    message,
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

	client := &http.Client{}

	smsaeroMessage, _ := buildSmsaeroMessage(smsContent, c.signature, numbers)
	url := fmt.Sprintf(
		c.url+"/sms/send?numbers=%s&text=%s&sign=%s", smsaeroMessage.numbers, smsaeroMessage.text, smsaeroMessage.sign)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	response.Body.Close()

	return nil
}
