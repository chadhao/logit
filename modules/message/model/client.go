package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	awsSNSSession *session.Session
	awsSESSession *session.Session
	awsSNS        *sns.SNS
	awsSES        *ses.SES
	config        map[string]string
)

func awsInit() (err error) {
	awsSNSAccessKeyID := config["message.aws.snsaccesskeyid"]
	awsSNSSecrectAccessKey := config["message.aws.snssecrectaccesskey"]
	awsSESAccessKeyID := config["message.aws.sesaccesskeyid"]
	awsSESSecrectAccessKey := config["message.aws.sessecrectaccesskey"]
	awsRegion := config["message.aws.region"]

	awsSNSSession, err = session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsSNSAccessKeyID, awsSNSSecrectAccessKey, ""),
	})
	awsSESSession, err = session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsSESAccessKeyID, awsSESSecrectAccessKey, ""),
	})
	awsSNS = sns.New(awsSNSSession)
	awsSES = ses.New(awsSESSession)
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
