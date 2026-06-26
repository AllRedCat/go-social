package main

import (
	"log"
	"os"

	"github.com/AllRedCat/go-social/internal/auth"
	"github.com/AllRedCat/go-social/internal/database"
	"github.com/AllRedCat/go-social/internal/posts"
	"github.com/gin-gonic/gin"
)

func main() {
	// Ensure database directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("Error creating database directory: ", err)
	}

	// Init db connection
	db, err := database.DbConnect("data/social.db")
	if err != nil {
		log.Fatal("Error on init database: ", err)
	}
	defer db.Close()

	// Create initial tables if does not exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		avatar_url TEXT DEFAULT '',
		created_at DATETIME NOT NULL,
		deleted_at DATETIME
	);
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INT NOT NULL,
		title TEXT NOT NULL,
		content TEXT DEFAULT '',
		image_url TEXT DEFAULT '',
		created_at DATETIME NOT NULL,
		updated_at DATETIME
	);`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatal("Fail on create table on database: ", err)
	}

	// Dependency injection
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	postsRepo := posts.NewRepository(db)
	postsService := posts.NewService(postsRepo)
	postsHandler := posts.NewHandler(postsService)

	// Router init
	router := gin.Default()

	// Static file server
	router.Static("/uploads", "./uploads")

	// Register routes
	authHandler.RegisterRoutes(router)
	postsHandler.RegisterRoutes(router)

	log.Println("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error on init server: ", err)
	}
}
