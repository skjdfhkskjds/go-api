package main

import (
	"log"

	"github.com/skjdfhkskjds/go-api/internal/engine"
	"github.com/skjdfhkskjds/go-api/internal/types"
)

func main() {
	// Create engine with default config
	e := engine.New(nil)

	// Register routes
	e.GET("/", homeHandler)
	e.GET("/health", healthHandler)
	e.GET("/users/:id", getUserHandler)
	e.POST("/users", createUserHandler)

	// Create a route group
	api := e.Group("/api/v1")
	api.GET("/status", statusHandler)

	// Start server
	log.Println("Starting basic example server...")
	if err := e.Run(":8080"); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func homeHandler(ctx *types.Context) {
	ctx.JSON(200, map[string]string{
		"message": "Welcome to the Go API framework!",
		"version": "1.0.0",
	})
}

func healthHandler(ctx *types.Context) {
	ctx.JSON(200, map[string]string{
		"status": "healthy",
	})
}

func getUserHandler(ctx *types.Context) {
	userID := ctx.GetParam("id")
	ctx.JSON(200, map[string]string{
		"user_id": userID,
		"name":    "User " + userID,
	})
}

func createUserHandler(ctx *types.Context) {
	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.Error(400, err)
		return
	}

	ctx.Success(map[string]any{
		"id":    "123",
		"name":  user.Name,
		"email": user.Email,
	})
}

func statusHandler(ctx *types.Context) {
	ctx.JSON(200, map[string]string{
		"api_version": "v1",
		"status":      "running",
	})
}
