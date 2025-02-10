package main

import (
	"alphabet_typer/records"
	"alphabet_typer/utility"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/joho/godotenv"
	"net/http"
	"time"
	"os"

	"github.com/labstack/echo/v4/middleware"
)

var dbClient *sql.DB

func startup() {
	err := godotenv.Load()
	if err != nil {
        fmt.Sprintf("%v", os.Getenv("dbName"))
    }

	dbClient = utility.NewDbClient()

	// ToDo !!Migration currently not working!!

	// for utility.Migrate(dbClient) != nil {
	// 	fmt.Printf("Verbindung Fehlgeschlagen, %s", utility.Migrate(dbClient))
	// 	time.Sleep(20 * time.Second)
	// }
}

func CORSMiddlewareWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		dynamicCORSConfig := middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
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

	port := "8080"
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
