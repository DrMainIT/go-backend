package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Urls struct {
	gorm.Model
	Name     string
	Approved bool
}

func Config(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		return "ok"
	}
	return os.Getenv(key)

}

func queryName(db *gorm.DB,approved bool) []string {
	var urls []Urls
	result := db.Find(&urls)
	if result.Error != nil {
		return nil
	}
	links := []string{}
	for _, url := range urls {
		if url.Approved == approved {
		links = append(links, url.Name)
		}
	}
	return links
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	db, err := gorm.Open(postgres.Open(Config("POSTGRES_URL")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Urls{})
	app.Get("/", func(c fiber.Ctx) error {
		// list all elements
		return c.Render("index", fiber.Map{})
	})
	app.Get("/add", func(c fiber.Ctx) error {
		// get values from the form
		url := c.FormValue("url")

		db.Create(&Urls{Name: url,Approved: false})
		//names := queryName(db)
		return c.Render("index", fiber.Map{})
	})

	app.Get("/login", func(c fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})
	app.Post("/login", func(c fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		if os.Getenv("USERNAME") == username && os.Getenv("PASSWORD") == password {
			return c.Redirect().To("/admin")
		}
		return c.Redirect().To("/login")
	})

	app.Get("/admin", func(c fiber.Ctx) error {
		names := queryName(db,false)

		return c.Render("admin", fiber.Map{
			"urls": names,
			"Title":       "Articles to buy",
		})
	})
	app.Post("/admin/accept", func(c fiber.Ctx) error{
		names := queryName(db,false)
		url := c.FormValue("link")
		fmt.Print(url)
		db.Model(&Urls{}).Where("Name = ?", url).Update("approved", true)
		return c.Render("admin", fiber.Map{
			"urls": names,
			"Title": "Articles to buy",
		})
	})
	
	app.Listen(":3000")
}
