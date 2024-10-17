package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	diff_utils "test-server/myapp/diff-utils"
	json_utils "test-server/myapp/json-utils"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Todo struct {
	Id     int    `json:"id"`
	Todo   string `json:"todo"`
	Status string `json:"status"`
}

type TodoModel struct {
	gorm.Model
	ID     uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Todo   string    `json:"todo" gorm:"type:varchar(110);"`
	Status string    `json:"status" gorm:"varchar(110);default:in_progress"`
}

type ResponseMessage struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

var db *sql.DB

func dbConnect() {
	var (
		host     = os.Getenv("DATABASE_HOST")
		port     = os.Getenv("DATABASE_PORT")
		user     = os.Getenv("DATABASE_USER")
		password = os.Getenv("DATABASE_PASSWORD")
		dbname   = os.Getenv("DATABASE_NAME")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Database connected!")

	err = db.AutoMigrate(&TodoModel{})
	if err != nil {
		panic(err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbConnect()

	var todoList []Todo

	e := echo.New()
	e.GET("/todos", func(c echo.Context) error {
		return c.JSON(http.StatusOK, todoList)
	})

	e.POST("/todos", func(c echo.Context) error {
		jsonBody := Todo{Id: rand.Intn(1000)}
		jsonBody, err := json_utils.Parse(jsonBody, c)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		todoList = append(todoList, jsonBody)
		return c.JSON(http.StatusOK, todoList)
	})

	e.DELETE("/todos/:id", func(c echo.Context) error {
		idParam, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Error on parsing")
		}

		currentTodo := diff_utils.FindInArr(todoList, "Id", idParam)

		todoList = append(todoList[:currentTodo.Index], todoList[currentTodo.Index+1:]...)

		return c.JSON(
			http.StatusOK,
			ResponseMessage{
				Message: "todo with id " + c.Param("id") + " deleted",
				Status:  "success",
			},
		)
	})

	e.PATCH("/todos/:id", func(c echo.Context) error {
		idParam, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Error on parsing")
		}

		var jsonBody Todo
		jsonBody, err = json_utils.Parse(jsonBody, c)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		var indexOfUpdatedTodo int

		for i := 0; i < len(todoList); i++ {
			if todoList[i].Id == idParam {
				indexOfUpdatedTodo = i
				break
			}
		}

		if jsonBody.Status != "" {
			todoList[indexOfUpdatedTodo].Status = jsonBody.Status
		}

		if jsonBody.Todo != "" {
			todoList[indexOfUpdatedTodo].Todo = jsonBody.Todo
		}

		return c.JSON(http.StatusOK, todoList)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
