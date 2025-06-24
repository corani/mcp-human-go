package web

import (
	"fmt"
	"time"

	"github.com/corani/mcp-human-go/internal/config"
	"github.com/corani/mcp-human-go/internal/memory"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type API struct {
	conf   *config.Config
	memory *memory.MemoryDB
	app    *fiber.App
}

func NewAPI(conf *config.Config, mem *memory.MemoryDB) *API {
	api := &API{
		conf:   conf,
		memory: mem,
		app: fiber.New(fiber.Config{
			Views:             html.New("./internal/web/views", ".html"),
			EnablePrintRoutes: true,
			StreamRequestBody: true,
			DisableKeepalive:  true,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      5 * time.Second,
		}),
	}

	api.app.Static("/", "./internal/web/static").Name("static")

	api.app.Get("/ui", api.HandleUI).Name("ui_index")
	api.app.Get("/ui/:name", api.HandleUI).Name("ui_name")
	api.app.Get("/api/memory", api.HandleList).Name("api_memory_list")
	api.app.Get("/api/memory/:id", api.HandleGet).Name("api_memory_get")
	api.app.Post("/api/memory/:id/answer", api.HandleAnswer).Name("api_memory_answer")

	return api
}

func (a *API) Start() error {
	addr := a.conf.WebPort

	return a.app.Listen(fmt.Sprintf(":%d", addr))
}

func (a *API) Shutdown() error {
	return a.app.Shutdown()
}

func (a *API) HandleUI(c *fiber.Ctx) error {
	name := c.Params("name", "index")

	return c.Render(name, fiber.Map{})
}

func (a *API) HandleList(c *fiber.Ctx) error {
	items := a.memory.ListQuestions()

	return c.JSON(items)
}

func (a *API) HandleGet(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "id is required")
	}

	item, ok := a.memory.Get(id)
	if !ok {
		return fiber.NewError(fiber.StatusNotFound)
	}

	return c.JSON(item)
}

func (a *API) HandleAnswer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "id is required")
	}

	answer := string(c.Body())
	if answer == "" {
		return fiber.NewError(fiber.StatusBadRequest, "answer is required in body")
	}

	err := a.memory.UpdateAnswer(id, answer)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
