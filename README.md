# Math-IA 🤖📊

Uma plataforma inteligente de IA que utiliza **Ollama** para gestão de modelos locais e **Milvus** como banco de dados vetorial, oferecendo capacidades avançadas de processamento de linguagem natural com contexto semântico.

## 🚀 Características

- **Gestão Multi-Modelo**: Seleção inteligente de modelos baseada no conteúdo da pergunta
- **Busca Semântica**: Utiliza embeddings para encontrar contexto relevante
- **RAG (Retrieval-Augmented Generation)**: Combina busca vetorial com geração de texto
- **API REST**: Interface HTTP simples e eficiente
- **Arquitetura Modular**: Código organizado em camadas bem definidas

## 🏗️ Arquitetura

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│   API Routes    │───▶│   AI Handler    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Model Selector │◀───┤   Ollama Client │
                       └─────────────────┘    └─────────────────┘
                                                        │
                       ┌─────────────────┐    ┌─────────────────┐
                       │ Vector Search   │◀───┤ Milvus Database │
                       └─────────────────┘    └─────────────────┘
```

## 📂 Estrutura do Projeto

```
math-ia/
├── main.go                     # Ponto de entrada da aplicação
├── go.mod                      # Dependências do projeto
├── internal/
│   ├── api/
│   │   └── ask.go              # Handlers da API REST
│   ├── ia/
│   │   ├── ollama/             # Cliente Ollama
│   │   │   ├── cliente.go      # Cliente HTTP base
│   │   │   ├── generate.go     # Geração de texto
│   │   │   └── embeding.go     # Geração de embeddings
│   │   ├── selector/
│   │   │   └── selector.go     # Seleção inteligente de modelos
│   │   └── vectorstore/
│   │       ├── milvus.go       # Cliente Milvus
│   │       ├── config/
│   │       │   └── milvus_config.go
│   │       └── operations/
│   ├── router/
│   │   └── routes.go           # Configuração de rotas
│   └── tools/
│       └── ingest.go           # Ingestão de documentos
├── embedEtcd.yaml              # Configuração etcd
├── standalone_embed.sh         # Script de inicialização
└── user.yaml                   # Configuração de usuário
```

## 🛠️ Tecnologias

- **Go 1.24+**: Linguagem principal
- **Ollama**: Servidor de modelos de IA local
- **Milvus**: Banco de dados vetorial
- **Chi Router**: Framework HTTP minimalista
- **godotenv**: Gerenciamento de variáveis de ambiente

## 📋 Pré-requisitos

- Go 1.24 ou superior
- Ollama instalado e rodando
- Milvus rodando (local ou Docker)
- Modelos baixados no Ollama:
  - `llama3.1:8b`
  - `nomic-embed-text`
  - `codellama`
  - `hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF`

## ⚙️ Configuração

1. Clone o repositório:
```bash
git clone <repository-url>
cd math-ia
```

2. Instale as dependências:
```bash
go mod tidy
```

3. Configure as variáveis de ambiente (`.env`):
```env
OLLAMA_HOST=http://localhost:11434
MILVUS_HOST=localhost
MILVUS_PORT=19530
```

4. Execute a aplicação:
```bash
go run main.go
```

## 🔌 Endpoints da API

### POST `/ask`
Faz uma pergunta simples sem contexto adicional.

**Request:**
```json
{
  "prompt": "Qual é a derivada de x²?"
}
```

**Response:**
```json
{
  "model": "hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF",
  "response": "A derivada de x² é 2x."
}
```

### POST `/ask-with-context`
Faz uma pergunta utilizando busca semântica para contexto adicional (RAG).

**Request:**
```json
{
  "prompt": "Como calcular uma integral definida?"
}
```

**Response:**
```json
{
  "model": "hf.co/iamcoder18/AceMath-7B-Instruct-Q4_K_M-GGUF",
  "response": "Para calcular uma integral definida..."
}
```

## 🧠 Modelos Disponíveis

| Modelo | Especialização | Palavras-chave |
|--------|----------------|----------------|
| `AceMath-7B` | Matemática | integral, derivada, cálculo, equação |
| `codellama` | Programação | python, golang, typescript, api |
| `llama3.1:8b` | Uso geral | Modelo padrão |
| `nomic-embed-text` | Embeddings | Geração de vetores |

## 🎯 Seleção Inteligente de Modelos

O sistema analisa o prompt e seleciona automaticamente o modelo mais adequado baseado em:

- **Palavras-chave específicas** (com pesos diferentes)
- **Contexto semântico** da pergunta
- **Especialização** de cada modelo

Exemplo:
- "Como resolver esta integral?" → `AceMath-7B-Instruct`
- "Escreva uma função em Python" → `codellama`
- "Explique sobre história" → `llama3.1:8b`

## 📊 Fluxo RAG (Retrieval-Augmented Generation)

1. **Recebe pergunta** do usuário
2. **Gera embedding** da pergunta usando `nomic-embed-text`
3. **Busca similaridade** no Milvus (top-3 documentos)
4. **Seleciona modelo** baseado no conteúdo
5. **Gera resposta** usando contexto encontrado
6. **Retorna resposta** enriquecida

## 🚀 Deploy

### Docker Compose (recomendado)
```yaml
version: '3.8'
services:
  math-ia:
    build: .
    ports:
      - "8081:8081"
    environment:
      - OLLAMA_HOST=http://ollama:11434
      - MILVUS_HOST=milvus
      - MILVUS_PORT=19530
    depends_on:
      - ollama
      - milvus
```

### Standalone
```bash
# Build
go build -o math-ia main.go

# Run
./math-ia
```

## 🔧 Desenvolvimento

### Estrutura de Dados
```go
type Document struct {
    ID      int64  `json:"id"`
    Text    string `json:"text"`
    Source  string `json:"source"`
    Content string `json:"content"`
}
```

### Adicionando Novos Modelos
1. Baixe o modelo no Ollama
2. Adicione em `selector.go`:
```go
WeightedKeywordsByModel["novo-modelo"] = map[string]int{
    "palavra-chave": 3,
}
```

## 📈 Monitoramento

A aplicação inclui:
- **Logs estruturados** para debugging
- **CORS configurado** para frontend
- **Tratamento de erros** robusto
- **Timeouts configuráveis**

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 🆘 Suporte

- **Documentação**: Este README
---

Feito com ❤️ usando Go e tecnologias de IA modernas.
