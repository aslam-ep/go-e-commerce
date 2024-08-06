package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/aslam-ep/go-e-commerce/docs/swagger"
	"github.com/aslam-ep/go-e-commerce/internal/address"
	"github.com/aslam-ep/go-e-commerce/internal/auth"
	"github.com/aslam-ep/go-e-commerce/internal/user"
	"github.com/aslam-ep/go-e-commerce/router/middleware"
	"github.com/aslam-ep/go-e-commerce/utils"
)

func SetupRoutes(r chi.Router, userHandler *user.UserHandler, authHandler *auth.AuthHandler, addressHandler *address.AddressHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			utils.WriteResponse(w, http.StatusAccepted, &struct {
				Message string `json:"message"`
			}{
				Message: "API up and running",
			})
		})

		// Registering the swagger UI handler
		r.Get("/swagger/*", httpSwagger.WrapHandler)

		// Auth Router group
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/refresh-token", authHandler.RefreshToken)
		})

		// User Router group
		r.Route("/users", func(r chi.Router) {
			r.Post("/create", userHandler.CreateUser)

			r.With(middleware.AuthMiddleware, middleware.ProfileMiddleware).Route("/{id}", func(r chi.Router) {
				r.Get("/", userHandler.GetUser)
				r.Put("/update", userHandler.UpdateUser)
				r.Put("/reset-password", userHandler.ResetPassword)
				r.Delete("/delete", userHandler.DeleteUser)

				// Address Router group
				r.Route("/addresses", func(r chi.Router) {
					r.Get("/", addressHandler.GetAllAddress)
					r.Post("/create", addressHandler.CreateAddress)
					r.Route("/{address_id}", func(r chi.Router) {
						r.Get("/", addressHandler.GetAddressByID)
						r.Put("/update", addressHandler.UpdateAddress)
						r.Put("/set-default", addressHandler.SetDefaultAddress)
						r.Delete("/delete", addressHandler.DeleteAddress)
					})
				})
			})
		})
	})
}
