package api

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
)

type Api struct {
    e *echo.Echo
    port string
}

func NewApi(port string) *Api {
    e := echo.New()
    registerHealthCheckRoute(e)
    
    return &Api{
        e: e,
        port: port,
    }
}

func registerHealthCheckRoute(e *echo.Echo) {
    e.GET("/health", func(c echo.Context) error {
        return c.String(http.StatusOK, "OK")
    })
}

func (a Api) Run() {
	a.e.Logger.Fatal(a.e.Start(a.port))
}
