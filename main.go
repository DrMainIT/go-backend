package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"os"
)

type Spesa struct {
	gorm.Model
	Name string
}
func Config(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		return "ok"
	}
	return os.Getenv(key)

}

func queryName(db *gorm.DB) []string {
	var products []Spesa 
	result := db.Find(&products)
	if result.Error != nil {
		return nil
	}
	names := []string{}
	for _,product := range products {
		names = append(names,product.Name)
	}
	return names
}

func main() {
	engine := html.New("./views",".html")
	app := fiber.New(fiber.Config{Views: engine})
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	db,err := gorm.Open(postgres.Open(Config("POSTGRES_URL")),&gorm.Config{})
	db.AutoMigrate(&Spesa{})
	app.Get("/", func(c fiber.Ctx) error {
		// list all elements
		names := queryName(db)

		return c.Render("index",fiber.Map{
			"ingredients": names,
			"Title": "Articles to buy",
		})
	})
	app.Get("/add", func(c fiber.Ctx) error {
		// get values from the form
		product := c.FormValue("product")
		db.Create(&Spesa{Name: product})
		names := queryName(db)
		return c.Render("index",fiber.Map{
			"ingredients": names,
			"Title": "Articles to buy",
		})
	})
	app.Listen(":3000")
}