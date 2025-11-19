// config/app.go
package config

import "github.com/gofiber/fiber/v2"

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Sistem Pelaporan Prestasi Mahasiswa Unair - UAS 2025",
	})
	return app
}