package main

import "github.com/slack-go/slack"

// SlackNotify...Slackに通知する
func SlackNotify(slackToken, message string) error {
	// xoxb-xxxxxxxの値
	api := slack.New(slackToken)

	// red表示させる
	attachment := slack.Attachment{
		Color: ColorRED,
		Title: SlackMsgTitle,
		Text:  message,
	}

	_, _, err := api.PostMessage(
		SlackChannelName,
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		return nil
	}
	return nil
}
