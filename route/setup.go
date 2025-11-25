// route/setup.go
package route

import (
	"uas-pelaporan-prestasi-backend/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Delete("/achievements/:id", service.DeleteAchievement)
}