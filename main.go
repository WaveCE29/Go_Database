package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	database = "mydatabase"
	username = "admin"
	password = "1244"
)

var db *sql.DB

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)

	dbc, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	db = dbc

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")

	app := fiber.New()

	app.Get("/", hello)
	app.Get("/product/:id", getProductHandler)
	app.Get("/product/", getProductsHandler)
	app.Post("/product", createProductHandler)
	app.Put("/product/:id", updateProductHandler)
	app.Delete("/product/:id", deleteProductHandler)

	app.Listen(":8080")
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Welcome to the Product API!")
}

func getProductHandler(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid product ID")
	}

	product, err := getProduct(productId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).SendString("Product not found")
		}
	}
	return c.JSON(product)
}

func getProductsHandler(c *fiber.Ctx) error {

	product, err := getProducts()
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).SendString("Product not found")
		}
	}
	return c.JSON(product)
}

func createProductHandler(c *fiber.Ctx) error {
	p := new(Product)

	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := createProduct(p)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(p)
}

func updateProductHandler(c *fiber.Ctx) error {

	productId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid product ID")
	}
	p := new(Product)

	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	product, err := updateProduct(productId, p)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).SendString("Product not found")
		}
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(product)
}

func deleteProductHandler(c *fiber.Ctx) error {

	productId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("Invalid product ID")
	}

	err = deleteProduct(productId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).SendString("Product not found")
		}
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
