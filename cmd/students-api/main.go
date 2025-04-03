package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NiteshSGupta/students-api/internal/config"
	"github.com/NiteshSGupta/students-api/internal/http/handlers/student"
	"github.com/NiteshSGupta/students-api/internal/storage/sqlite"
)

func main() {
	// fmt.Println("My first project in go")

	//1. load config
	cfg := config.MustLoad()

	//2. database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//3. setup router
	//net http is default package , we don't have to install
	// The net/http package in Go is used whenever you want to work with HTTP—things like building web servers, making web requests, or handling web stuff in your program. It’s like a toolbox for talking over the internet using the HTTP protocol (the language web browsers and servers use). Let’s break it down simply, like figuring out when to use a hammer in your toolbox!
	//net http is used for methods like get,post , query parameter , path perameter, this package used for both rounting and making server

	//NewServeMux return router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to students api"))
	})

	router.HandleFunc("POST /api/students", student.New(storage))

	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	//4. setup server
	server :=
		http.Server{
			Addr:    cfg.Addr,
			Handler: router,
		}

	// fmt.Printf("Server start %s", cfg.HttpServer.Addr)
	slog.Info("Server start", slog.String("address", cfg.Addr))

	// created channel
	// Channels = A way for those helpers to pass stuff (like info or results) to each other neatly.
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Goroutines = Helpers doing tasks at the same time, so your program doesn’t have to wait around.
	go func() {
		// fmt.Println("Server start")
		//here we used go routines for continous running with diffrent request and response, multithreading
		// fmt.Printf("Server start %s", cfg.HttpServer.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	//so after the above code
	slog.Info("shuting down the server")

	//this will notify the us after 5 seconds to shutdown the server if any process is running more than 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer used for automatic call
	// In Go, defer is a powerful and unique feature that schedules a function call to be executed after the surrounding function completes, regardless of how it completes (normal return, panic, etc.).
	defer cancel()

	//now server get 5 seconds notify
	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	// if err := server.Shutdown(ctx); err != nil {
	// 	slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	// }

	slog.Info("server shutdown succesfully")

}
