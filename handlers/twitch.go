package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"apiComentariosAI/models"
	"apiComentariosAI/services"

	"github.com/gorilla/mux"
)

// TwitchHandler gerencia as requisições relacionadas à Twitch
type TwitchHandler struct {
	twitchAPI     *services.TwitchAPI
	geminiService *services.GeminiService
}

// NewTwitchHandler cria uma nova instância do handler
func NewTwitchHandler() *TwitchHandler {
	return &TwitchHandler{
		twitchAPI:     services.NewTwitchAPI(),
		geminiService: services.NewGeminiService(),
	}
}

// GetUserStatus verifica o status de um usuário e retorna comentários do chat
func (h *TwitchHandler) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		http.Error(w, "Username é obrigatório", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 Verificando status do usuário: %s", username)

	// Verifica se o usuário está online
	streamCheckStart := time.Now()
	streamData, err := h.twitchAPI.CheckStreamStatus(username)
	streamCheckDuration := time.Since(streamCheckStart)
	log.Printf("⏱️  [TEMPO] Verificação de stream: %v", streamCheckDuration)

	if err != nil {
		log.Printf("❌ Erro ao verificar status: %v", err)
		http.Error(w, "Erro ao verificar status do usuário", http.StatusInternalServerError)
		return
	}

	response := models.TwitchResponse{
		Username: username,
		IsOnline: streamData != nil,
		Comments: make([]models.ChatComment, 0),
	}

	// Se estiver online, preenche os dados da stream
	if streamData != nil {
		response.StreamTitle = streamData.Title
		response.Game = streamData.GameName
		response.ViewerCount = streamData.ViewerCount

		log.Printf("✅ Usuário %s está online. Conectando ao chat...", username)

		// Captura comentários do chat
		commentsStart := time.Now()
		comments, err := services.GetRecentComments(username, 5)
		commentsDuration := time.Since(commentsStart)
		log.Printf("⏱️  [TEMPO] Captura de comentários: %v", commentsDuration)

		if err != nil {
			log.Printf("⚠️  Erro ao capturar comentários: %v", err)
			// Não retorna erro, apenas não inclui comentários
		} else {
			response.Comments = comments
			log.Printf("💬 Capturados %d comentários", len(comments))

			// Gera um comentário usando IA se houver comentários
			if len(comments) > 0 {
				context := fmt.Sprintf("Stream: %s - Jogo: %s", streamData.Title, streamData.GameName)

				aiStart := time.Now()
				generatedComment, err := h.geminiService.GenerateComment(comments, context)
				aiDuration := time.Since(aiStart)
				log.Printf("⏱️  [TEMPO] Geração de comentário IA: %v", aiDuration)

				if err != nil {
					log.Printf("❌ Erro ao gerar comentário com IA: %v", err)
				} else {
					response.GeneratedComment = generatedComment
					log.Printf("🤖 Comentário gerado: %s", generatedComment)
				}
			}
		}
	} else {
		log.Printf("💤 Usuário %s está offline", username)
	}

	totalDuration := time.Since(startTime)
	log.Printf("⏱️  [TEMPO TOTAL] Requisição completa: %v", totalDuration)

	// Retorna resposta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HealthCheck verifica se a API está funcionando
func (h *TwitchHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "API de Comentários da Twitch está funcionando",
	})
}
