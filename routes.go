package main

import (
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/gofiber/fiber/v2"
	"github.com/srevinsaju/chibi/diag"
	"gorm.io/gorm"
	"math/rand"
	"net/url"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func CreateRoutes(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Post("/create", func(c *fiber.Ctx) error {
		var diags diag.Diagnostics
		var redirect RedirectComponent

		fun := c.FormValue("type") == "fun"
		short := c.FormValue("type") == "short"
		var id string
		if fun {
			for {
				id = petname.Generate(3, "-")
				db.First(&redirect, "id = ?", id)
				if redirect.Id == "" {
					break
				}
			}
		} else if short {
			for {
				// create a random string, 6 characters long
				id = randSeq(6)
				db.First(&redirect, "id = ?", id)
				if redirect.Id == "" {
					break
				}
			}
		} else {
			return diags.Append(diag.Error, "invalid type").JSON(c, nil)
		}

		redirect.Id = id
		redirect.RedirectURL = c.FormValue("url")
		if redirect.RedirectURL == "" {
			return diags.Append(diag.Error, "url is empty").JSONWithStatus(c, nil, fiber.StatusBadRequest)
		}
		u, err := url.Parse(redirect.RedirectURL)
		if err != nil {

			return diags.Append(diag.Error, "invalid url").JSONWithStatus(c, nil, fiber.StatusBadRequest)
		}
		redirect.RedirectURL = u.String()
		redirect.Alive = true
		db.Create(&redirect)
		diags = diags.Append(diag.Info, "created redirect")
		return diags.JSON(c, redirect)
	})

	app.Post("/create/:id", func(c *fiber.Ctx) error {
		var diags diag.Diagnostics
		var redirect RedirectComponent
		db.First(&redirect, "id = ?", c.Params("id"))
		if redirect.Id != "" {
			return diags.Append(diag.Error, "id already exists").JSON(c, nil)
		}
		redirect.Id = c.Params("id")
		redirect.RedirectURL = c.FormValue("url")
		if redirect.RedirectURL == "" {
			return diags.Append(diag.Error, "url is empty").JSONWithStatus(c, nil, fiber.StatusBadRequest)
		}
		u, err := url.Parse(redirect.RedirectURL)
		if err != nil {
			return diags.Append(diag.Error, "invalid url").JSONWithStatus(c, nil, fiber.StatusBadRequest)
		}
		redirect.RedirectURL = u.String()
		redirect.Alive = true
		db.Create(&redirect)
		diags = diags.Append(diag.Info, "created redirect")
		return diags.JSON(c, redirect)
	})

	app.Get("/url/:id", func(c *fiber.Ctx) error {
		var redirect RedirectComponent
		db.First(&redirect, "id = ?", c.Params("id"))
		var diags diag.Diagnostics
		if redirect.Id == "" {
			return diags.Append(diag.Error, "id not found").JSONWithStatus(c, nil, fiber.StatusNotFound)
		}
		return diags.JSONWithStatus(c, redirect, fiber.StatusOK)
	})

	return app
}
