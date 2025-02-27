package main

import (
	"fmt"
	"html/template"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type FormData struct {
	Values map[string]string
	Errors []string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: nil,
	}
}

type Page struct {
	Form FormData
}

func newPage() Page {
	return Page{
		Form: newFormData(),
	}
}

type Order struct {
	Tea    int    `json:"tea"`
	Coffee int    `json:"coffee"`
	Cake   int    `json:"cake"`
	Soda   int    `json:"soda"`
	Name   string `json:"name"`
	Errors []string
}

func parseOrder(tea string, coffee string, cake string, soda string, name string) Order {

	order := Order{
		Tea:    parseInt(tea),
		Coffee: parseInt(coffee),
		Cake:   parseInt(cake),
		Soda:   parseInt(soda),
		Name:   name,
	}

	if order.Tea == 0 && order.Coffee == 0 && order.Cake == 0 && order.Soda == 0 {
		order.Errors = append(order.Errors, "Please order something")
	}

	if len(order.Name) < 3 {
		order.Errors = append(order.Errors, "Please provide 3 letters name")
	}

	return order
}

func parseInt(value string) int {
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = newTemplate()

	e.Static("/images", "images")
	e.Static("/css", "css")
	page := newPage()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/order", func(c echo.Context) error {

		order := parseOrder(
			c.FormValue("tea"),
			c.FormValue("coffee"),
			c.FormValue("cake"),
			c.FormValue("soda"),
			c.FormValue("name"),
		)
		e.Logger.Debugf("Received order: %+v", order)

		if order.Errors != nil {
			form := newFormData()
			form.Values["tea"] = strconv.Itoa(order.Tea)
			form.Values["caffee"] = strconv.Itoa(order.Coffee)
			form.Values["cake"] = strconv.Itoa(order.Cake)
			form.Values["soda"] = strconv.Itoa(order.Soda)
			form.Values["name"] = order.Name
			form.Errors = order.Errors
			return c.Render(422, "form", form)
		}

		c.Render(200, "form", newFormData())

		orderId := order.Name + "-" + time.Now().UTC().Format("2006-01-02 15:04:05.350")
		return c.Render(200, "confirmation", fmt.Sprintf("Your order id is: %v", orderId))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
