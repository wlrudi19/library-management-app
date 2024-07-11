package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/wlrudi19/library-management-app/go-user-services/app/api"
	"github.com/wlrudi19/library-management-app/go-user-services/app/repository"
	"github.com/wlrudi19/library-management-app/go-user-services/app/service"
	"github.com/wlrudi19/library-management-app/go-user-services/config"
)

func main() {
	loadConfig := config.LoanConfig()
	connDB, connRedis, err := config.ConnectConfig(loadConfig.Database, loadConfig.Redis)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	defer connDB.Close()
	defer connRedis.Close()

	//init
	userRepository := repository.NewUserRepository(connDB, connRedis)
	userLogic := service.NewUserLogic(userRepository)
	userHandler := api.NewUserHandler(userLogic)
	userRouter := api.NewUserRouter(userHandler)

	//group route
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Mount("/users", userRouter)
	})

	server := http.Server{
		Addr:    "localhost:3012",
		Handler: r,
	}

	fmt.Println("starting server on port 3012...")

	err = server.ListenAndServe()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
