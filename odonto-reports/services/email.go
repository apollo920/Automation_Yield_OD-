package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func EnviarEmail(relatorio string) error {
	// Configurações do Gmail
	from := os.Getenv("GMAIL_USER")
	password := os.Getenv("GMAIL_PASSWORD") // Senha de App (sem espaços)
	to := os.Getenv("DESTINATION_EMAIL")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Mensagem no formato MIME
	subject := "Relatório Financeiro Odonto Company"
	body := relatorio
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

	// Autenticação
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Envia o e-mail
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Printf("Erro detalhado ao enviar e-mail: %v", err)
		return fmt.Errorf("erro ao enviar e-mail: %v", err)
	}

	return nil
}