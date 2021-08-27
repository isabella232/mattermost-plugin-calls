package main

import (
	"encoding/json"
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

type clusterMessage struct {
	UserID        string        `json:"user_id"`
	ChannelID     string        `json:"channel_id"`
	SenderID      string        `json:"sender_id"`
	ClientMessage clientMessage `json:"client_message"`
}

type clusterMessageType string

const (
	clusterMessageTypeConnect    clusterMessageType = "connect"
	clusterMessageTypeDisconnect clusterMessageType = "disconnect"
	clusterMessageTypeSignaling  clusterMessageType = "signaling"
)

func (m *clusterMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *clusterMessage) FromJSON(data []byte) error {
	return json.Unmarshal(data, &m)
}

func (p *Plugin) sendClusterMessage(msg clusterMessage, msgType clusterMessageType, targetID string) error {
	msgData, err := msg.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to encode to JSON: %w", err)
	}

	ev := model.PluginClusterEvent{
		Id:   string(msgType),
		Data: msgData,
	}

	opts := model.PluginClusterEventSendOptions{
		SendType: model.PluginClusterEventSendTypeReliable,
		TargetId: targetID,
	}

	if appErr := p.API.PublishPluginClusterEvent(ev, opts); appErr != nil {
		return fmt.Errorf("failed to publish cluster event: %w", appErr)
	}

	return nil
}