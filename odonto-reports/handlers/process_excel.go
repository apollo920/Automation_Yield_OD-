package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"odonto-reports/models"
	"odonto-reports/services"
	"sync"
)

// Estrutura para armazenar os dados extraídos
type ReportData struct {
	DiasUteis    int
	DiasCorridos int
	DiasFaltam   int
	Pilares      map[string]PilarData
}

// Estrutura para armazenar os dados de cada Pilar
type PilarData struct {
	Real        float64
	Meta        float64
	PercentReal float64
	Projecao    float64
	PercentProj float64
}
// Função auxiliar para converter string em float
func parseFloat(value string) (float64, error) {
	if value == "" {
		return 0, nil
	}
	var num float64
	_, err := fmt.Sscanf(value, "%f", &num)
	return num, err
}

var (
	relatorioGerado models.Relatorio // Armazena os dados processados
	mu              sync.Mutex       // Mutex para evitar concorrência
)

// Processa o Excel e armazena os dados processados
func ProcessExcelHandler(c *fiber.Ctx) error {
	filePath := "uploads/relatorio.xlsx"
	dadosProcessados, err := services.ProcessarExcel(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao processar o Excel",
		})
	}

	// Salvar os dados no cache global
	mu.Lock()
	relatorioGerado = dadosProcessados
	mu.Unlock()

	// Retorna os dados processados
	return c.JSON(fiber.Map{
		"message": "Arquivo processado com sucesso",
		"data":    dadosProcessados,
	})
}
