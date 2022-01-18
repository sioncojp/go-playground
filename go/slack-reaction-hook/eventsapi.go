package main

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackListener struct {
	api    *slack.Client
	client *socketmode.Client
}

// NewSlack ...Websocketで受け付けるための、SlackListenerの初期化
func NewSlack() *SlackListener {
	api := slack.New(
		BotToken,
		slack.OptionAppLevelToken(AppToken),
	)

	client := socketmode.New(api)

	return &SlackListener{
		api:    api,
		client: client,
	}
}

// ListenAndResponse ...Websocketの立ち上げ
func (s *SlackListener) ListenAndResponse() {
	// Handle slack events
	for e := range s.client.Events {

		switch e.Type {
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent := e.Data.(slackevents.EventsAPIEvent)
			s.client.Ack(*e.Request)

			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
				case *slackevents.ReactionAddedEvent:
					if err := s.handleReactionAddedEvent(ev); err != nil {
						log.sugar.Errorf("Failed to reaction added event: %s", err)
					}
				}
			}
		}
	}
}

// handleMesageEvent ...handles message events
func (s *SlackListener) handleReactionAddedEvent(ev *slackevents.ReactionAddedEvent) error {
	if ev.Reaction == TriggerReaction {
		if _, _, err := s.api.PostMessage(
			ev.Item.Channel,
			slack.MsgOptionText("test dayo", false),
		); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
	}
	return nil
}
