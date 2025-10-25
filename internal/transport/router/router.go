package router

import (
	"hackaton/internal/config"
	"hackaton/internal/handler"
	custommiddleware "hackaton/internal/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer(jwtSecret []byte, h *handler.Handler, config *config.Config) *http.Server {
	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
	}))

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/login", h.Login)
			v1.POST("/register", h.CreateUser)

			protected := v1.Group("")
			protected.Use(custommiddleware.AuthRequired(jwtSecret))
			{

				users := protected.Group("/users")
				users.GET("", h.GetUsers)
				users.GET("/:id", h.GetUserByID)
				users.PUT("/:id", h.UpdateUser)
				users.DELETE("/:id", h.DeleteUser)
			}
		}
	}

	return &http.Server{
		Addr:         config.HTTPServer.Port,
		Handler:      router,
		WriteTimeout: config.HTTPServer.WriteTimeout,
		ReadTimeout:  config.HTTPServer.ReadTimeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}
}
