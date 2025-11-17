package models

import "time"

// ChatComment representa um comentário do chat
type ChatComment struct {
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// TwitchResponse representa a resposta da API
type TwitchResponse struct {
	Username         string        `json:"username"`
	IsOnline         bool          `json:"is_online"`
	StreamTitle      string        `json:"stream_title,omitempty"`
	Game             string        `json:"game,omitempty"`
	ViewerCount      int           `json:"viewer_count,omitempty"`
	Comments         []ChatComment `json:"comments"`
	GeneratedComment string        `json:"generated_comment,omitempty"`
}

// TwitchStreamData representa dados da stream da API Twitch
type TwitchStreamData struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	UserLogin    string `json:"user_login"`
	UserName     string `json:"user_name"`
	GameID       string `json:"game_id"`
	GameName     string `json:"game_name"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	ViewerCount  int    `json:"viewer_count"`
	StartedAt    string `json:"started_at"`
	Language     string `json:"language"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// TwitchStreamsResponse representa a resposta da API de streams
type TwitchStreamsResponse struct {
	Data []TwitchStreamData `json:"data"`
}

// TwitchTokenResponse representa a resposta do OAuth token
type TwitchTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
