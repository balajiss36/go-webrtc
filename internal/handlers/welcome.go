package handlers

import "github.com/gofiber/fiber/v2"

func Welcome(c *fiber.Ctx) error {
	c.Render("Welcome to Video Call!", nil, "layouts/main")
	return nil
}
