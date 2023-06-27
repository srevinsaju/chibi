package main

import (
	"embed"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v2"
	"net/http"
	"os"
)

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Embed a directory
//
//go:embed views/*.html
var embedViews embed.FS

type RedirectComponent struct {
	gorm.Model
	RedirectURL string `json:"redirect_url"`
	Id          string `gorm:"unique,primaryKey" json:"id"`
	Alive       bool   `json:"alive"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		frontendUrl = "localhost:8080"
	}

	// Migrate the schema
	err = db.AutoMigrate(&RedirectComponent{})
	if err != nil {
		panic(err)
	}

	engine := django.NewFileSystem(http.FS(embedViews), ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).SendFile(fmt.Sprintf("views/%d.html", code))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	})

	app.Mount("/api/v1", CreateRoutes(db))
	app.Get("/:id", func(c *fiber.Ctx) error {
		var redirect RedirectComponent
		db.First(&redirect, "id = ?", c.Params("id"))
		if redirect.Id == "" {
			c.Status(fiber.StatusNotFound)
			return c.Render("views/404", fiber.Map{})
		}
		return c.Redirect(redirect.RedirectURL, fiber.StatusTemporaryRedirect)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/index", fiber.Map{
			"frontend_url": frontendUrl,
		})
	})

	err = app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
