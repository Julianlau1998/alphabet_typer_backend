package main

import (
	"alphabet_typer/records"
	"alphabet_typer/utility"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var dbClient *sql.DB

func startup() {

	dbClient = utility.NewDbClient()
	for utility.Migrate(dbClient) != nil {
		fmt.Printf("Verbindung Fehlgeschlagen, %s", utility.Migrate(dbClient))
		time.Sleep(20 * time.Second)
	}
}

func CORSMiddlewareWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		dynamicCORSConfig := middleware.CORSConfig{
			AllowOrigins: []string{"https://alphabet-typer.com", "https://alphabet-typer.netlify.app"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		}
		CORSMiddleware := middleware.CORSWithConfig(dynamicCORSConfig)
		CORSHandler := CORSMiddleware(next)
		return CORSHandler(ctx)
	}
}

func main() {
	startup()

	RecordRepository := records.NewRepository(dbClient)
	RecordService := records.NewService(RecordRepository)
	RecordDelivery := records.NewDelivery(RecordService)

	e := echo.New()

	e.Use(CORSMiddlewareWrapper)

	r := e.Group("/api")

	r.GET("/records", RecordDelivery.GetAll)
	r.POST("/records", RecordDelivery.Post)

	port := "8081"
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
