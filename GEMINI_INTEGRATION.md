# Guia de Integração com Gemini AI

## Configuração

### 1. Obter API Key do Gemini

1. Acesse: https://makersuite.google.com/app/apikey
2. Faça login com sua conta Google
3. Clique em "Create API Key"
4. Copie a chave gerada

### 2. Configurar Variável de Ambiente

Adicione no seu arquivo `.env`:

```bash
GEMINI_API_KEY=sua_api_key_aqui
```

## Como Funciona

### Fluxo de Geração

```
1. Usuário faz requisição GET /api/twitch/{username}
   ↓
2. API verifica se streamer está online
   ↓
3. Se online, conecta ao IRC e captura 5 comentários
   ↓
4. Envia comentários + contexto para o Gemini
   ↓
5. Gemini analisa e gera um novo comentário
   ↓
6. Retorna JSON com comentários originais + comentário gerado
```

### Estrutura do Prompt

O serviço envia para o Gemini:

- **Contexto**: Título da stream e nome do jogo
- **Comentários**: Lista dos últimos 5 comentários com username
- **Instruções**: 
  - Gerar comentário natural e humano
  - Máximo 200 caracteres
  - Em português brasileiro
  - Relacionado ao contexto
  - Sem emojis excessivos

### Exemplo de Prompt Enviado

```
Você é um assistente que analisa comentários de um chat ao vivo da Twitch 
e gera um novo comentário relevante e natural.

Contexto da stream: KKKKKKKKK - Jogo: Counter-Strike 2

Comentários recentes do chat:
1. usuario1: KKKKKKKKKK
2. usuario2: VAI GAULES!
3. usuario3: Que jogada insana!
4. usuario4: Perdeu feio aí
5. usuario5: ACE!!!

Com base nesses comentários, gere UM único comentário que seria apropriado 
para participar dessa conversa...
```

### Exemplo de Resposta da API

```json
{
  "username": "gaules",
  "is_online": true,
  "stream_title": "KKKKKKKKK",
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
  "generated_comment": "Essa jogada foi muito boa! O CS2 tá demais hoje! 🔥"
}
```

## Tratamento de Erros

### Sem API Key

Se `GEMINI_API_KEY` não estiver configurada:
- A API continua funcionando normalmente
- Captura os comentários do chat
- Não gera comentário com IA
- Campo `generated_comment` fica vazio

### Erro na API do Gemini

Se houver erro ao chamar o Gemini:
- Log do erro é registrado
- API retorna os comentários capturados
- Campo `generated_comment` fica vazio
- Não afeta o funcionamento principal

### Sem Comentários no Chat

Se não houver comentários para analisar:
- Gemini não é chamado (economia de API calls)
- Retorna apenas status da stream

## Otimizações Futuras

### Possíveis Melhorias

1. **Cache de Respostas**: Cachear comentários gerados por alguns segundos
2. **Streaming**: Usar streaming do Gemini para resposta mais rápida
3. **Personalização**: Permitir configurar estilo do comentário (formal, casual, etc)
4. **Múltiplos Comentários**: Gerar N comentários diferentes
5. **Análise de Sentimento**: Adicionar análise de sentimento do chat
6. **Rate Limiting**: Implementar limite de requisições ao Gemini

### Custos

O Gemini oferece:
- **Free tier**: 60 requisições por minuto
- **Custo**: Muito baixo para uso moderado

Recomendação: Implementar cache para reduzir chamadas à API.

## Debugging

### Testar Prompt Manualmente

Você pode testar o prompt no Google AI Studio:
https://makersuite.google.com/

### Logs

A aplicação registra:
```
Capturados X comentários
Comentário gerado: [texto do comentário]
Erro ao gerar comentário com IA: [erro]
```

### Teste Local

```bash
# Execute a API
go run main.go

# Em outro terminal, teste com curl
curl http://localhost:8080/api/twitch/gaules | jq
```

## Segurança

### Boas Práticas

1. ⚠️ **Nunca commite** o arquivo `.env` com API keys
2. ✅ Use variáveis de ambiente em produção
3. ✅ Rotacione as API keys periodicamente
4. ✅ Implemente rate limiting
5. ✅ Monitore uso da API para evitar custos inesperados

### API Key no GitHub Actions

Se usar CI/CD, adicione como secret:

```yaml
env:
  GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
```
