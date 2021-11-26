package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	SsmSlackTokenName = "/ssm/token_name"
	ColorRED          = "#F08080"
	SlackMsgTitle     = "Fargate: error to running container"
)

var (
	SlackChannelName = os.Getenv("SlackChannelName")
	eventDetails     CloudWatchEventDetails
	AwsAccountIDs    = map[string]string{
		"123456789012": "aws-account-name01",
		"111111111111": "aws-account-name02",
	}
)

// handleRequest...メイン処理
func handleRequest(event events.CloudWatchEvent) error {
	// SlackTokenのdecrypt
	var ssmClient SsmClient
	if err := ssmClient.NewClient(); err != nil {
		return err
	}
	slackToken, err := ssmClient.Decrypt(SsmSlackTokenName)
	if err != nil {
		return err
	}

	// Cloudwatchのeventを拾って、exit 0 以外をslackに通知させる
	c, err := ToCloudWatchEventDetailsStruct(event)
	if err != nil {
		return err
	}

	// Containersの中はサイドカー含めるコンテナの数sliceがあるので
	// exit 0以外のコンテナ全て取得
	containerStatus := make(map[string]int)
	for _, v := range c.Containers {
		if v.ExitCode != 0 {
			containerStatus[v.Name] = v.ExitCode
		}
	}

	// exit 0以外があるなら通知
	if len(containerStatus) != 0 {
		return SlackNotify(slackToken, c.FormatToSlackMessage(containerStatus))
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
