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

	// Rota para processar o Excel
	app.Post("/process_excel", handlers.ProcessExcelHandler)

	// Rota para gerar o relatório automaticamente
	app.Get("/gerar_relatorio", handlers.GerarRelatorioHandler)

	log.Fatal(app.Listen(":3000"))
}

