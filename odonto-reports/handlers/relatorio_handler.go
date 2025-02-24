package handlers

import (
	"github.com/gofiber/fiber/v2"
	"odonto-reports/services"
)

// Gera o relatório a partir dos dados processados
func GerarRelatorioHandler(c *fiber.Ctx) error {
	mu.Lock()
	defer mu.Unlock()

	// Verifica se há dados processados
	if relatorioGerado.Pilares == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nenhum dado processado ainda",
		})
	}

	// Gera o relatório
	relatorio := services.GerarRelatorio(relatorioGerado)

	// Envia o relatório por e-mail
	err := services.EnviarEmail(relatorio)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao enviar relatório por e-mail",
		})
	}

	// Retorna o relatório
	return c.JSON(fiber.Map{
		"relatorio": relatorio,
	})
}
