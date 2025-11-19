// main.go
package main

import (
	"log"
	"uas-pelaporan-prestasi-backend/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Load()

	app := config.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "Server UAS Backend Lanjut hidup bro! 🚀",
			"author":  "Nama Kamu",
			"nim":     "NIM Kamu",
		})
	})

	port := config.Env("APP_PORT", "3000")
	log.Println("Server berjalan di http://localhost:" + port + "/health")
	log.Fatal(app.Listen(":" + port))
}