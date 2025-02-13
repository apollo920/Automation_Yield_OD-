package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"odonto-reports/handlers"
)

func main() {
	app := fiber.New()

	// Rota para upload do arquivo Excel
	app.Post("/upload", handlers.UploadExcelHandler)

	log.Println("Servidor rodando na porta 3000...")
	log.Fatal(app.Listen(":3000"))
}
