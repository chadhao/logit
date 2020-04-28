package model

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Txt 电话短信
type Txt struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

func (t *Txt) valid() error {
	t.Number = strings.TrimLeft(t.Number, "0")
	if !strings.HasPrefix(t.Number, "+64") {
		t.Number = "+64" + t.Number
	}
	return nil
}

// Send 发送短信
func (t *Txt) Send() error {
	if err := t.valid(); err != nil {
		return err
	}

	// log.Println(t)
	params := &sns.PublishInput{
		PhoneNumber: aws.String(t.Number),
		Message:     aws.String(t.Message),
	}

	_, err := awsSNS.Publish(params)
	return err
}
