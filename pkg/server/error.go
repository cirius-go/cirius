package server

import "github.com/labstack/echo/v4"

// DefineCustomError define new custom error at package level.
func DefineCustomError(code int, msg ...interface{}) *echo.HTTPError {
	e := echo.NewHTTPError(code, msg...)

	return e
}
