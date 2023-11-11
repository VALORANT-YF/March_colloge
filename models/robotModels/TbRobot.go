package robotModels

import "gorm.io/gorm"

type TbRobot struct {
	gorm.Model
	WebhookURL string `json:"webhook_url"` //群聊机器人的webhookURL
}

func (t TbRobot) TableName() string {
	return "tb_robot"
}
