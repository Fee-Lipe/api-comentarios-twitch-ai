package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"apiComentariosAI/models"
)

const (
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-lite:generateContent"
)

// GeminiService gerencia interações com a API do Gemini
type GeminiService struct {
	apiKey string
	client *http.Client
}

// NewGeminiService cria uma nova instância do serviço Gemini
func NewGeminiService() *GeminiService {
	return &GeminiService{
		apiKey: os.Getenv("GEMINI_API_KEY"),
		client: &http.Client{},
	}
}

// GeminiRequest representa a estrutura de requisição para o Gemini
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse representa a resposta da API do Gemini
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// GenerateComment gera um novo comentário baseado nos comentários capturados
func (g *GeminiService) GenerateComment(comments []models.ChatComment, context string) (string, error) {
	if g.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY não configurada")
	}

	// Constrói o prompt com os comentários
	prompt := g.buildPrompt(comments, context)

	// Prepara a requisição
	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Faz a requisição para o Gemini
	url := fmt.Sprintf("%s?key=%s", geminiAPIURL, g.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro na API do Gemini (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse da resposta
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Extrai o texto gerado
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado pela IA")
	}

	generatedText := geminiResp.Candidates[0].Content.Parts[0].Text
	return strings.TrimSpace(generatedText), nil
}

// buildPrompt constrói o prompt para o Gemini baseado nos comentários
func (g *GeminiService) buildPrompt(comments []models.ChatComment, context string) string {
	var sb strings.Builder

	sb.WriteString("Você é um assistente que analisa comentários de um chat ao vivo da Twitch e gera um novo comentário relevante e natural.\n\n")

	if context != "" {
		sb.WriteString(fmt.Sprintf("Contexto da stream: %s\n\n", context))
	}

	sb.WriteString("Comentários recentes do chat:\n")
	for i, comment := range comments {
		sb.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, comment.Username, comment.Message))
	}

	sb.WriteString("\nCom base nesses comentários, gere UM único comentário que seria apropriado para participar dessa conversa. ")
	sb.WriteString("O comentário deve:\n")
	sb.WriteString("- Ser natural e parecer humano\n")
	sb.WriteString("- Ter no máximo 200 caracteres\n")
	sb.WriteString("- Estar relacionado ao contexto dos comentários\n")
	sb.WriteString("- Não usar emojis excessivos\n")
	sb.WriteString("- Ser em português brasileiro\n")
	sb.WriteString("- Não incluir aspas ou marcações especiais\n\n")
	sb.WriteString("Responda APENAS com o comentário gerado, sem explicações adicionais.")

	return sb.String()
}
