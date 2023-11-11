package dingOfficialControllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func SendTextMessage(message string, webhookURL string) error {

	messageDing := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": "三月软件\n" + message,
		},
		"at": map[string]interface{}{
			"atMobiles": []string{},
			"isAtAll":   true,
		},
	}

	return sendMessage(webhookURL, messageDing)
}

func sendMessage(webhookURL string, message map[string]interface{}) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	r := struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return err
	}
	if r.Errcode != 0 {
		return errors.New(r.Errmsg)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}
	return nil
}
