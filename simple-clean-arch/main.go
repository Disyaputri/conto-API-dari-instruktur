package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"simple-http/handler"
	"simple-http/repository"
	"simple-http/usecase"
)

func main() {
	userRepo := repository.NewUserRepository()
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUC)

	r := gin.Default()
	r.GET("/api/v1/users", userHandler.GetUsers)
	r.GET("/api/v1/users/:id", userHandler.GetUserByID)
	r.POST("/api/v1/users", userHandler.CreateUser)
	r.PUT("/api/v1/users/:id", userHandler.UpdateUser)
	r.DELETE("/api/v1/users/:id", userHandler.DeleteUser)

	log.Println("RUNNING on :3000")
	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
