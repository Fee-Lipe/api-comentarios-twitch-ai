package services

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"apiComentariosAI/models"
)

const (
	ircServer = "irc.chat.twitch.tv:6667"
	ircNick   = "justinfan12345" // Nick anônimo padrão da Twitch
)

// TwitchIRC gerencia a conexão com o chat da Twitch via IRC
type TwitchIRC struct {
	conn net.Conn
}

// NewTwitchIRC cria uma nova instância do cliente IRC
func NewTwitchIRC() *TwitchIRC {
	return &TwitchIRC{}
}

// Connect estabelece conexão com o servidor IRC da Twitch
func (t *TwitchIRC) Connect() error {
	conn, err := net.DialTimeout("tcp", ircServer, 10*time.Second)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao IRC: %w", err)
	}
	t.conn = conn

	// Autentica como usuário anônimo
	fmt.Fprintf(t.conn, "NICK %s\r\n", ircNick)
	fmt.Fprintf(t.conn, "USER %s 8 * :%s\r\n", ircNick, ircNick)

	return nil
}

// JoinChannel entra em um canal específico
func (t *TwitchIRC) JoinChannel(channel string) error {
	if t.conn == nil {
		return fmt.Errorf("não conectado ao IRC")
	}

	// Remove # do início se existir
	channel = strings.TrimPrefix(channel, "#")

	// Solicita capacidades para obter tags (nome do usuário, etc)
	fmt.Fprintf(t.conn, "CAP REQ :twitch.tv/tags twitch.tv/commands\r\n")

	// Entra no canal
	fmt.Fprintf(t.conn, "JOIN #%s\r\n", channel)

	return nil
}

// ReadComments lê mensagens do chat e retorna os últimos N comentários
func (t *TwitchIRC) ReadComments(maxComments int, timeout time.Duration) ([]models.ChatComment, error) {
	if t.conn == nil {
		return nil, fmt.Errorf("não conectado ao IRC")
	}

	comments := make([]models.ChatComment, 0, maxComments)
	reader := bufio.NewReader(t.conn)

	// Define timeout para a leitura
	deadline := time.Now().Add(timeout)
	t.conn.SetReadDeadline(deadline)

	for len(comments) < maxComments && time.Now().Before(deadline) {
		line, err := reader.ReadString('\n')
		if err != nil {
			// Se atingiu o timeout, retorna o que foi coletado
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			return comments, fmt.Errorf("erro ao ler do IRC: %w", err)
		}

		line = strings.TrimSpace(line)

		// Responde aos PINGs para manter a conexão viva
		if strings.HasPrefix(line, "PING") {
			pong := strings.Replace(line, "PING", "PONG", 1)
			fmt.Fprintf(t.conn, "%s\r\n", pong)
			continue
		}

		// Parse de mensagens PRIVMSG (mensagens do chat)
		if strings.Contains(line, "PRIVMSG") {
			comment := t.parseMessage(line)
			if comment != nil {
				comments = append(comments, *comment)
			}
		}
	}

	return comments, nil
}

// parseMessage faz o parse de uma mensagem IRC e extrai username e mensagem
func (t *TwitchIRC) parseMessage(line string) *models.ChatComment {
	// Formato IRC: :username!username@username.tmi.twitch.tv PRIVMSG #channel :message
	// ou com tags: @badges=... :username!username@username.tmi.twitch.tv PRIVMSG #channel :message

	// Remove tags se existirem
	if strings.HasPrefix(line, "@") {
		parts := strings.SplitN(line, " :", 2)
		if len(parts) >= 2 {
			line = ":" + parts[1]
		}
	}

	// Extrai username
	parts := strings.Split(line, " ")
	if len(parts) < 4 {
		return nil
	}

	userPart := parts[0]
	if !strings.HasPrefix(userPart, ":") {
		return nil
	}

	username := strings.Split(strings.TrimPrefix(userPart, ":"), "!")[0]

	// Extrai mensagem (tudo após PRIVMSG #channel :)
	msgStartIdx := strings.Index(line, "PRIVMSG")
	if msgStartIdx == -1 {
		return nil
	}

	remainder := line[msgStartIdx:]
	colonIdx := strings.Index(remainder, ":")
	if colonIdx == -1 {
		return nil
	}

	message := remainder[colonIdx+1:]

	return &models.ChatComment{
		Username:  username,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// Close fecha a conexão IRC
func (t *TwitchIRC) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

// GetRecentComments conecta ao canal e captura os últimos N comentários
func GetRecentComments(channel string, maxComments int) ([]models.ChatComment, error) {
	irc := NewTwitchIRC()

	if err := irc.Connect(); err != nil {
		return nil, err
	}
	defer irc.Close()

	if err := irc.JoinChannel(channel); err != nil {
		return nil, err
	}

	// Aguarda um pouco para garantir que entrou no canal
	time.Sleep(1 * time.Second)

	// Lê comentários por até 10 segundos ou até obter maxComments
	comments, err := irc.ReadComments(maxComments, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
