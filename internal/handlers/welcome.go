package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Welcome(c *fiber.Ctx) error {
	err := c.Render("Welcome to Video Call!", nil, "layouts/main")
	if err != nil {
		log.Fatalf("Error in rendering welcome page: %v", err)
	}
	return nil
}
