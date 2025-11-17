# API de Comentários da Twitch

API em Go que verifica se um streamer está online e captura os últimos 5 comentários do chat em tempo real.

## Funcionalidades

- ✅ Verificar se um usuário está online/offline
- ✅ Conectar ao chat da Twitch via IRC
- ✅ Capturar os últimos 5 comentários em tempo real
- ✅ Gerar comentários inteligentes usando Gemini AI
- ✅ API REST simples e fácil de usar

## Pré-requisitos

1. **Twitch Client ID e Client Secret**
   - Acesse: https://dev.twitch.tv/console/apps
   - Crie uma nova aplicação
   - Copie o Client ID e Client Secret

2. **Gemini API Key**
   - Acesse: https://makersuite.google.com/app/apikey
   - Crie uma nova API Key
   - Copie a chave gerada

3. **(Opcional) OAuth Token para IRC**
   - Para leitura básica do chat, não é necessário
   - Para chat autenticado: https://twitchapps.com/tmi/

## Instalação

1. Clone o repositório
```bash
cd /home/juninho/apiComentariosAI
```

2. Configure as variáveis de ambiente
```bash
cp .env.example .env
# Edite o arquivo .env com suas credenciais
```

3. Instale as dependências
```bash
go mod download
```

4. Execute a aplicação
```bash
go run main.go
```

## Endpoints da API

### 1. Verificar status e obter comentários

**GET** `/api/twitch/:username`

Verifica se o usuário está online e retorna os últimos 5 comentários do chat.

**Exemplo:**
```bash
curl http://localhost:8080/api/twitch/gaules
```

**Resposta (Online):**
```json
{
  "username": "gaules",
  "is_online": true,
  "stream_title": "KKKKKKKKKKKKKKKKK",
  "game": "Counter-Strike 2",
  "viewer_count": 45000,
  "comments": [
    {
      "username": "usuario1",
      "message": "KKKKKKKKKK",
      "timestamp": "2025-11-16T10:30:45Z"
    },
    {
      "username": "usuario2",
      "message": "VAI GAULES!",
      "timestamp": "2025-11-16T10:30:46Z"
    }
  ],
  "generated_comment": "Essa partida tá muito boa! O CS2 tá demais hoje!"
}
```

**Resposta (Offline):**
```json
{
  "username": "usuariooffline",
  "is_online": false,
  "comments": []
}
```

### 2. Health Check

**GET** `/health`

Verifica se a API está funcionando.

```bash
curl http://localhost:8080/health
```

## Estrutura do Projeto

```
.
├── main.go              # Entrada da aplicação
├── handlers/            # Handlers HTTP
│   └── twitch.go
├── services/            # Lógica de negócio
│   ├── twitch_api.go    # Integração com Twitch API
│   ├── twitch_irc.go    # Integração com IRC/Chat
│   └── gemini.go        # Integração com Gemini AI
├── models/              # Estruturas de dados
│   └── models.go
├── go.mod
├── .env
└── README.md
```

## Como funciona a IA?

A API utiliza o **Gemini AI** do Google para gerar comentários inteligentes baseados no contexto do chat:

1. **Captura**: Os últimos 5 comentários do chat são capturados em tempo real
2. **Contexto**: Informações da stream (título e jogo) são adicionadas ao contexto
3. **Geração**: O Gemini analisa os comentários e gera um novo comentário natural
4. **Características do comentário gerado**:
   - Natural e humano
   - Relacionado ao contexto da conversa
   - Máximo de 200 caracteres
   - Em português brasileiro
   - Sem emojis excessivos

### Exemplo de Prompt

```
Contexto da stream: KKKKKKKKKKKKKKKKK - Jogo: Counter-Strike 2

Comentários recentes:
1. usuario1: KKKKKKKKKK
2. usuario2: VAI GAULES!
3. usuario3: Que jogada!

Gere um comentário apropriado para essa conversa...
```

## Observações Técnicas

### Como funciona?

1. **Verificação de Status**: Usa a Twitch Helix API para verificar se o streamer está online
2. **Chat IRC**: Conecta ao servidor IRC da Twitch (irc.chat.twitch.tv:6667)
3. **Captura de Mensagens**: Monitora o chat por alguns segundos e captura as mensagens
4. **Geração por IA**: Envia os comentários para o Gemini gerar um novo comentário
5. **Timeout**: Se o usuário estiver offline ou sem mensagens, retorna após 10 segundos

### Limitações

- A captura de comentários leva alguns segundos (5-10s) pois precisa conectar ao IRC
- A geração de comentário depende da disponibilidade da API do Gemini
- Em chats muito ativos, pode capturar mais de 5 mensagens rapidamente
- Requer credenciais válidas da Twitch API

## Licença

MIT
