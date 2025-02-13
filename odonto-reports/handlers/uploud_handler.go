package handlers

import (
	"fmt"
	"log"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
)

// UploadExcelHandler recebe o arquivo Excel e o processa
func UploadExcelHandler(c *fiber.Ctx) error {
	// Obtém o arquivo enviado
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Falha ao receber o arquivo",
		})
	}

	// Define o caminho para salvar o arquivo temporariamente
	savePath := filepath.Join("uploads", file.Filename)

	// Cria o diretório de uploads se não existir
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	// Salva o arquivo no servidor
	err = c.SaveFile(file, savePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Falha ao salvar o arquivo",
		})
	}

	// Processar o Excel após o upload
	reportData, err := ProcessExcel(savePath)
	if err != nil {
		log.Println("Erro ao processar o arquivo Excel:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao processar o arquivo Excel",
		})
	}

	// Retornar os dados extraídos
	return c.JSON(fiber.Map{
		"message":     "Arquivo processado com sucesso",
		"dias_uteis":  reportData.DiasUteis,
		"dias_corridos": reportData.DiasCorridos,
		"dias_faltam": reportData.DiasFaltam,
		"pilares":     reportData.Pilares,
	})
}
