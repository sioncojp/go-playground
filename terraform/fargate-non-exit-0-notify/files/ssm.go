package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// SsmClient...SSMに関わるクライアント情報をストア
type SsmClient struct {
	Session *session.Session
	svc     ssmiface.SSMAPI
}

// NewClient...クライアント情報を初期化する
func (c *SsmClient) NewClient() error {
	c.Session = session.Must(session.NewSessionWithOptions(session.Options{
		Config: *aws.NewConfig().WithCredentialsChainVerboseErrors(true),
	}))

	c.svc = ssm.New(c.Session, aws.NewConfig().WithMaxRetries(10).WithCredentialsChainVerboseErrors(true))
	return nil
}

// Decrypt...復号化する
func (c *SsmClient) Decrypt(ssmName string) (string, error) {
	params := &ssm.GetParameterInput{
		Name:           aws.String(ssmName),
		WithDecryption: aws.Bool(true),
	}
	resp, err := c.svc.GetParameter(params)
	if err != nil {
		return "", err
	}
	return *resp.Parameter.Value, nil
}
