package handlers

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

// UploadExcelHandler recebe o arquivo Excel e o armazena temporariamente
func UploadExcelHandler(c *fiber.Ctx) error {

	// Obtém o arquivo enviado
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Falha ao receber o arquivo",
		})
	}

	// Criando diretório se não existir
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	// Salvando o arquivo com um nome fixo
	filePath := "uploads/relatorio.xlsx"
	err = c.SaveFile(file, filePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao salvar arquivo"})
	}
	

	return c.JSON(fiber.Map{"message": "Arquivo enviado com sucesso!"})
}
