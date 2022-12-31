package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/THAI-DEV/assessment/database"
	"github.com/THAI-DEV/assessment/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var port string

// var dbUrl string

func init() {
	port = os.Getenv("PORT")
	// dbUrl = os.Getenv("DATABASE_URL")
}

func main() {
	// fmt.Println("Please use server.go for main file")
	// fmt.Println("start at port:", os.Getenv("PORT"))
	database.CreateTable()
	initGin()
}

func initGin() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", handler.Root)
	r.POST("/expenses", handler.Create)
	r.GET("/expenses/:id", handler.Read)
	r.PUT("/expenses/:id", handler.Update)

	srv := &http.Server{
		Addr:           port,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println("--- start at port:", port, "---")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit
	log.Println("... Shutdowning Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(" ... Server Shutdowning ...:", err)
	}

	<-ctx.Done()
	log.Println("--- Server exiting ---")
}
