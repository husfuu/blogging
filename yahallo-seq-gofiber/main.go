package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SayYahalloReq struct {
	Hello string `json:"hello"`
	Name  string `json:"name"`
}

func main() {
	app := fiber.New()

	log := NewLogger()

	app.Get("/", func(c *fiber.Ctx) error {
		data := "Yahallo, World!"
		fmt.Println(c.Request().URI())
		log.Infof("Info | %s |%s", c.Request().URI().String(), data)
		return c.SendString(data)
	})

	app.Get("/yahallos", func(c *fiber.Ctx) error {

		data := []string{
			"yahallo 1",
			"yahallo 2",
			"yahallo 3",
		}

		return Response_Log(c, log, fiber.StatusOK, "Success to get all yahallos", data)
	})

	app.Post("/say-yahallo", func(c *fiber.Ctx) error {
		var payload SayYahalloReq

		err := c.BodyParser(&payload)
		if err != nil {
			return Response_Log(c, log, fiber.StatusBadRequest, "fail to parse request body", nil)
		}

		data := fmt.Sprintf("%s (yahallo) %s", payload.Hello, payload.Name)

		return Response_Log(c, log, fiber.StatusOK, "Success to say hello", data)
	})

	app.Listen(":3000")
}
