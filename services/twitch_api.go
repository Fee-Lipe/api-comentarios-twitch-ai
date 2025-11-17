package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"apiComentariosAI/models"
)

// TwitchAPI gerencia as chamadas à API da Twitch
type TwitchAPI struct {
	clientID     string
	clientSecret string
	accessToken  string
}

// NewTwitchAPI cria uma nova instância do serviço
func NewTwitchAPI() *TwitchAPI {
	return &TwitchAPI{
		clientID:     os.Getenv("TWITCH_CLIENT_ID"),
		clientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
	}
}

// GetAccessToken obtém um token de acesso da Twitch
func (t *TwitchAPI) GetAccessToken() error {
	tokenURL := "https://id.twitch.tv/oauth2/token"

	data := url.Values{}
	data.Set("client_id", t.clientID)
	data.Set("client_secret", t.clientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return fmt.Errorf("erro ao obter token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao obter token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp models.TwitchTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("erro ao decodificar resposta do token: %w", err)
	}

	t.accessToken = tokenResp.AccessToken
	return nil
}

// CheckStreamStatus verifica se um usuário está online e retorna informações da stream
func (t *TwitchAPI) CheckStreamStatus(username string) (*models.TwitchStreamData, error) {
	if t.accessToken == "" {
		if err := t.GetAccessToken(); err != nil {
			return nil, err
		}
	}

	apiURL := fmt.Sprintf("https://api.twitch.tv/helix/streams?user_login=%s", username)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Client-ID", t.clientID)
	req.Header.Set("Authorization", "Bearer "+t.accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		// Token expirado, tenta renovar
		if err := t.GetAccessToken(); err != nil {
			return nil, err
		}
		// Tenta novamente com o novo token
		return t.CheckStreamStatus(username)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na API: status %d, body: %s", resp.StatusCode, string(body))
	}

	var streamResp models.TwitchStreamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&streamResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Se não há dados, o usuário está offline
	if len(streamResp.Data) == 0 {
		return nil, nil
	}

	return &streamResp.Data[0], nil
}
