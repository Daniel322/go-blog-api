package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"strconv"
)

type User struct {
	email    string
	password string
}

type Todo struct {
	Id     int    `json:"id"`
	Todo   string `json:"todo"`
	Status string `json:"status"`
}

type ResponseMessage struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

var mockUser = User{email: "test@test.com", password: "qwerty"}

func jsonParse[T comparable](jsonBody T, c echo.Context) error {
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "Error on JSON parsing")
	}
}

func main() {
	var todoList []Todo
	e := echo.New()
	e.GET("/todos", func(c echo.Context) error {
		return c.JSON(http.StatusOK, json.NewEncoder(c.Response()).Encode(todoList))
	})

	e.POST("/todos", func(c echo.Context) error {
		jsonBody := Todo{Id: rand.Intn(1000)}
		err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "Error on JSON parsing")
		}
		todoList = append(todoList, jsonBody)
		fmt.Println(todoList)
		return c.JSON(http.StatusOK, json.NewEncoder(c.Response()).Encode(todoList))
	})

	e.DELETE("/todos/:id", func(c echo.Context) error {
		idParam, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Error on parsing")
		}

		var indexOfDeletedTodo int

		for i := 0; i < len(todoList); i++ {
			if todoList[i].Id == idParam {
				indexOfDeletedTodo = i
				break
			}
		}

		todoList = append(todoList[:indexOfDeletedTodo], todoList[indexOfDeletedTodo+1:]...)

		return c.JSON(
			http.StatusOK,
			json.NewEncoder(c.Response()).Encode(
				ResponseMessage{
					Message: "todo with id " + c.Param("id") + " deleted",
					Status:  "success",
				},
			),
		)
	})

	e.PATCH("/todos/:id", func(c echo.Context) error {
		idParam, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Error on parsing")
		}

		var jsonBody Todo
		err = json.NewDecoder(c.Request().Body).Decode(&jsonBody)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusBadRequest, "Error on JSON parsing")
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

		return c.JSON(http.StatusOK, json.NewEncoder(c.Response()).Encode(todoList))
	})
	e.Logger.Fatal(e.Start(":1323"))
}
