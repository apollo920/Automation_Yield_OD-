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
	app.Post("/process_excel", func(c *fiber.Ctx) error {
		handlers.ProcessExcelHandler(c.Context().Response().BodyWriter(), c.Context().Request())
		return nil
	})

	// Rota para gerar o relatório automaticamente
	app.Get("/gerar_relatorio", func(c *fiber.Ctx) error {
		handlers.GerarRelatorioHandler(c.Context().Response().BodyWriter(), c.Context().Request())
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
