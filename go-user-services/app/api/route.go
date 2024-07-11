package api

import (
	"github.com/go-chi/chi"
	"github.com/wlrudi19/library-management-app/go-user-services/infrastructure/middlewares"
)

func NewUserRouter(userHandler UserHandler) *chi.Mux {
	intAuth := middlewares.InternalAuth
	r := chi.NewRouter()
	r.Use(intAuth)

	r.Post("/login", userHandler.LoginUser)

	return r
}
