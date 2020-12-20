package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/anggri-microservice/users-service/internal/db"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// In this example we use the html template engine
	"github.com/gofiber/template/html"
)

//ResponseTemplate variable
var ResponseTemplate map[string]interface{}

const grpcServer string = "192.168.1.222:6003"

func main() {
	// Create a new engine by passing the template folder
	// and template extension using <engine>.New(dir, ext string)
	engine := html.New("./views", ".html")

	// We also support the http.FileSystem interface
	// See examples below to load templates from embedded files
	engine = html.NewFileSystem(http.Dir("./views"), ".html")

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false

	// Debug will print each template that is parsed, good for debugging
	engine.Debug(true) // Optional. Default: false

	// Layout defines the variable name that is used to yield templates within layouts
	engine.Layout("embed") // Optional. Default: "embed"

	// Delims sets the action delimiters to the specified strings
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters

	// After you created your engine, you can pass it to Fiber's Views Engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Cors headers
	app.Use(cors.New())

	// Load go env
	var err error
	if os.Getenv("SRV_DOT_ENV") == "true" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Connect db postgres
	go initServer(app)

	log.Fatal(app.Listen(":5000"))
}

func initServer(app *fiber.App) {
	var err error
	db.DBConn, err = db.PostgreSQL.ConnectSqlx()

	if err != nil {
		log.Println(err)
	}

	// To render a template, you can call the ctx.Render function
	// Render(tmpl string, values interface{}, layout ...string)
	app.Get("/users/:limit/:page", GetUsers)
	app.Get("/users/:id", GetUserByID)
	app.Post("/users", CreateUsers)
	app.Delete("/users/:id", DeleteUsers)
}

func getTotalUsers() (total int64, err error) {
	SQL := `select count(*) as total from users`
	err = db.DBConn.QueryRow(SQL).Scan(&total)

	return
}

//DeleteUsers tambah users
func DeleteUsers(c *fiber.Ctx) error {
	SQL := `delete from users where id=$1`
	_, err := db.DBConn.Exec(SQL, c.Params("id"))
	if err != nil {
		log.Println(err)
		messageError := err.Error()
		errTemplate := map[string]interface{}{
			"status":  fiber.StatusBadRequest,
			"message": messageError,
			"total":   0,
			"data":    nil,
		}

		return c.Status(fiber.StatusBadRequest).JSON(errTemplate)
	}

	successTemplate := map[string]interface{}{
		"status":  fiber.StatusOK,
		"message": "User berhasil didelete",
		"total":   1,
		"data":    nil,
	}
	return c.JSON(successTemplate)
}

//CreateUsers tambah users
func CreateUsers(c *fiber.Ctx) error {
	var request map[string]interface{}
	err := json.Unmarshal(c.Body(), &request)
	SQL := `INSERT INTO "public"."users" ( "name", "occupation", "nik", "gender", "status", "religion", "address") 
	VALUES ( $1, $2, $3, $4, $5, $6, $7 );`

	_, err = db.DBConn.Exec(SQL, request["name"], request["occupation"], request["nik"], request["gender"], request["status"], request["religion"], request["address"])
	if err != nil {
		log.Println(err)
		messageError := err.Error()
		errTemplate := map[string]interface{}{
			"status":  fiber.StatusBadRequest,
			"message": messageError,
			"total":   0,
			"data":    nil,
		}

		return c.Status(fiber.StatusBadRequest).JSON(errTemplate)
	}

	successTemplate := map[string]interface{}{
		"status":  fiber.StatusOK,
		"message": "Berhasil menambahkan user",
		"total":   1,
		"data":    nil,
	}
	return c.JSON(successTemplate)
}

//GetUsers get list users
func GetUsers(c *fiber.Ctx) error {
	SQL := `select ROW_NUMBER () OVER (ORDER BY id) as position, id, name, occupation, nik, gender, status, religion, address, created_at, updated_at from users limit $1 offset $2`
	type Users struct {
		Position   int64  `json:"position"`
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		Occupation string `json:"occupation"`
		Nik        string `json:"nik"`
		Gender     string `json:"gender"`
		Status     int    `json:"status"`
		Religion   string `json:"religion"`
		Address    string `json:"address"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	//set limit and offset list
	limit, err := strconv.Atoi(c.Params("limit"))
	page, err := strconv.Atoi(c.Params("page"))
	offset := (limit * page) - limit

	var users []Users
	err = db.DBConn.Select(&users, SQL, c.Params("limit"), offset)
	if err != nil {
		log.Println(err)
		messageError := err.Error()
		errTemplate := map[string]interface{}{
			"status":  fiber.StatusBadRequest,
			"message": messageError,
			"total":   0,
			"data":    nil,
		}

		return c.Status(fiber.StatusBadRequest).JSON(errTemplate)
	}

	total, _ := getTotalUsers()
	successTemplate := map[string]interface{}{
		"status":  fiber.StatusOK,
		"message": "Sukses menampilkan list user",
		"total":   total,
		"data":    users,
	}
	return c.JSON(successTemplate)
}

//GetUserByID get list users
func GetUserByID(c *fiber.Ctx) error {
	SQL := `select id, name, occupation, nik, gender, status, religion, address, created_at, updated_at from users where id=$1`
	type Users struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		Occupation string `json:"occupation"`
		Nik        string `json:"nik"`
		Gender     string `json:"gender"`
		Status     int    `json:"status"`
		Religion   string `json:"religion"`
		Address    string `json:"address"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	var users []Users
	err := db.DBConn.Select(&users, SQL, c.Params("id"))

	if err != nil {
		log.Println(err)
		messageError := err.Error()
		errTemplate := map[string]interface{}{
			"status":  fiber.StatusBadRequest,
			"message": messageError,
			"total":   0,
			"data":    nil,
		}

		return c.Status(fiber.StatusBadRequest).JSON(errTemplate)
	}

	var data interface{}
	if len(users) != 0 {
		data = users[0]
	}

	successTemplate := map[string]interface{}{
		"status":  fiber.StatusOK,
		"message": "Sukses menampilkan list user",
		"total":   1,
		"data":    data,
	}
	return c.JSON(successTemplate)
}
