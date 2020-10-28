package main

import (
	"log"

	"github.com/codemicro/contactFormProcessor/internal/endpoints"
	"github.com/codemicro/contactFormProcessor/internal/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, _ error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(endpoints.StatusResponse{"internal server error"})
		},
	})

	app.Use(recover.New())

	app.Get("/captcha", endpoints.EndpointGenerateCaptcha)
	app.Post("/send", limiter.New(limiter.Config{Key: helpers.GetIP}), endpoints.EndpointProcessEmail)

	log.Panic(app.Listen(":80"))
}
