package main

import (
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"strconv"
	"test-server/myapp/diff-utils"
	"test-server/myapp/json-utils"
)

type Todo struct {
	Id     int    `json:"id"`
	Todo   string `json:"todo"`
	Status string `json:"status"`
}

type ResponseMessage struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
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
