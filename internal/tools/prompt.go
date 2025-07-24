package tools

import (
	"fmt"
	"strings"
)

const instruction = `
Você é um assistente que responde **exclusivamente** sobre assuntos relacionados à empresa "Produzindo Certo". 
Use **apenas** as informações fornecidas no contexto abaixo para gerar sua resposta. 
Se a pergunta não estiver relacionada à Produzindo Certo ou não puder ser respondida com base no contexto, diga educadamente que não pode responder.
`

func BuildPrompt(context []string, question string) string {
	return fmt.Sprintf("%s\n\nContexto:\n%s\n\nPergunta:\n%s\n\nResposta:",
		strings.TrimSpace(instruction),
		strings.Join(context, "\n"),
		strings.TrimSpace(question),
	)
}
