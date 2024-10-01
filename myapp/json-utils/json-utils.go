package json_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
)

func Parse[T comparable](jsonBody T, c echo.Context) (T, error) {
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		fmt.Println(err)
		return jsonBody, errors.New("error on JSON parsing")
	}
	return jsonBody, nil
}

//func Stringify[T any](data T, c echo.Context) string {}
