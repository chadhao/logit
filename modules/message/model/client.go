package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	awsSession *session.Session
	awsSNS     *sns.SNS
	awsSES     *ses.SES
	config     map[string]string
)

func awsInit() (err error) {
	awsAccessKeyID := config["message.aws.accesskeyid"]
	awsSecrectAccessKey := config["message.aws.secrectaccesskey"]
	awsRegion := config["message.aws.region"]

	awsSession, err = session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecrectAccessKey, ""),
	})
	awsSNS = sns.New(awsSession)
	awsSES = ses.New(awsSession)
	return
}

// New 传入config
func New(c map[string]string) error {
	config = c

	if err := awsInit(); err != nil {
		return err
	}

	return nil
}

// Close 关闭
func Close() {
}
