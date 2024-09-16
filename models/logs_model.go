package models

import "time"

type LogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	ServerID    string    `json:"server_id"`
	ServerName  string    `json:"server_name"`
	ChannelID   string    `json:"channel_id"`
	ChannelName string    `json:"channel_name"`
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Message     string    `json:"message"`
}
