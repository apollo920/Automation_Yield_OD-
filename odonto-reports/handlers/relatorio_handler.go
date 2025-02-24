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
	if relatorioGerado.DiasUteis == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nenhum dado processado ainda",
		})
	}

	// Gera o relatório
	relatorio := services.GerarRelatorio(relatorioGerado)

	// Retorna o relatório
	return c.JSON(fiber.Map{
		"relatorio": relatorio,
	})
}
