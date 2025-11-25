package service

import (
	"uas-pelaporan-prestasi-backend/app/repository"

	"github.com/gofiber/fiber/v2"
)

func DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.SoftDeleteAchievementMongo(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	err = repository.UpdateAchievementReferencePostgres(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal update reference: " + err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Prestasi draft dihapus (soft delete dengan status deleted)"})
}