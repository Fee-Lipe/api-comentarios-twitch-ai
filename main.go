package main

import (
	"log"
	"net/http"
	"os"

	"apiComentariosAI/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	// Valida variáveis obrigatórias
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("ERRO: TWITCH_CLIENT_ID e TWITCH_CLIENT_SECRET são obrigatórios!\n" +
			"Obtenha suas credenciais em: https://dev.twitch.tv/console/apps")
	}

	// Configura porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicializa handlers
	twitchHandler := handlers.NewTwitchHandler()

	// Configura rotas
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", twitchHandler.HealthCheck).Methods("GET")

	// Endpoint principal
	router.HandleFunc("/api/twitch/{username}", twitchHandler.GetUserStatus).Methods("GET")

	// Middleware de logging
	router.Use(loggingMiddleware)

	// Inicia servidor
	log.Printf("🚀 Servidor iniciado na porta %s", port)
	log.Printf("📡 Endpoints disponíveis:")
	log.Printf("   GET /health")
	log.Printf("   GET /api/twitch/{username}")
	log.Printf("\n💡 Exemplo: curl http://localhost:%s/api/twitch/gaules\n", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

// loggingMiddleware registra todas as requisições HTTP
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
