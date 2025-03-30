package main

import (
	"fmt"
	"os"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"  
    "github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/encryptcookie"
)

type Urls struct {
	gorm.Model
	Name     string
	Approved bool
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
	// Cookie 
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: "secret-thirty-2-character-string",
	}))
	// Initialize default config
	app.Use(cors.New())
	
	app.Use(logger.New())
	
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://gofiber.io", "https://gofiber.net", os.Getenv("FRONTEND_URL")},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Urls{})
	app.Get("/", func(c fiber.Ctx) error {
		// list all elements
		urls := queryName(db,true)
		return c.JSON(fiber.Map{
			"urls":  urls,
		})
	})
	app.Post("/add", func(c fiber.Ctx) error {
		// get values from json
		url := c.FormValue("url")

		db.Create(&Urls{Name: url,Approved: false})
		return c.Redirect().To(os.Getenv("FRONTEND_URL")+"/thankyou")
	
	})

	app.Get("/admin", func(c fiber.Ctx) error {
		names := queryName(db,false)
		return c.JSON(fiber.Map{
			"urls": names,
		})
	})
	app.Post("/admin/accept", func(c fiber.Ctx) error{
		url := c.FormValue("url")
		password := c.FormValue("password")

		if password == os.Getenv("PASSWORD") {
			db.Model(&Urls{}).Where("Name = ?", url).Update("approved", true)
		}
		return c.Redirect().To(os.Getenv("FRONTEND_URL")+ "/admin")
	})
	
	app.Listen(":8000")
}


