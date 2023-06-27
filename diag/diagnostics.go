package diag

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type Severity int

const (
	Info Severity = iota
	Warning
	Error
)

type Diagnostic struct {
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
}

type Diagnostics []Diagnostic

func (d Diagnostics) HasErrors() bool {
	for i := range d {
		if d[i].Severity == Error {
			return true
		}
	}
	return false
}

func (d Diagnostics) Append(severity Severity, message string) Diagnostics {
	return append(d, Diagnostic{
		Severity: severity,
		Message:  message,
	})
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Diags   Diagnostics `json:"diagnostics"`
}

func (d Diagnostics) JSON(c *fiber.Ctx, data interface{}) error {
	return d.JSONWithStatus(c, data, fiber.StatusOK)
}

func (d Diagnostics) JSONWithStatus(c *fiber.Ctx, data interface{}, status int) error {

	var resp Response
	resp.Diags = d
	resp.Success = !d.HasErrors()
	resp.Data = data

	j, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set("Content-Type", "application/json")
	c.Status(status)
	return c.Send(j)
}
