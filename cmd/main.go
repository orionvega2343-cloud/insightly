package main

import (
	"insightly/internal/ai"
	"insightly/internal/db"
	"insightly/internal/handlers"
	"insightly/internal/middlewares"
	"insightly/internal/repositories"
	"insightly/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("[FATAL] Failed to load .env: ", err)
	}

	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	aiClient := ai.NewOpenAIClient(&client)

	database, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	secret := os.Getenv("JWT_SECRET")

	// Repositories
	userRepo := repositories.NewUserRepo(database)
	refreshTokensRepo := repositories.NewRefreshTokensRepo(database)
	filesRepo := repositories.NewFilesRepo(database)
	queriesRepo := repositories.NewQueriesRepo(database)

	// Services
	userService := services.NewUserServiceImpl(userRepo, refreshTokensRepo, secret)
	refreshTokensService := services.NewRefreshTokensServiceImpl(refreshTokensRepo)
	filesService := services.NewFilesService(filesRepo)
	queriesService := services.NewQueriesService(filesRepo, queriesRepo, aiClient)

	// Handlers
	userHandler := handlers.NewUserHandler(userService, refreshTokensService, secret)
	refreshTokensHandler := handlers.NewRefreshTokensHandlerImpl(refreshTokensService, secret, userService)
	filesHandler := handlers.NewFilesHandler(filesService)
	queriesHandler := handlers.NewQueriesHandler(queriesService)

	// Router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", userHandler.Register)
		authRoutes.POST("/login", userHandler.Login)
		authRoutes.POST("/refresh", refreshTokensHandler.Refresh)
	}

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware(secret))
	{
		protected.POST("/files/upload", filesHandler.CreateFiles)
		protected.GET("/files", filesHandler.GetFilesByUserId)
		protected.DELETE("/files/:id", filesHandler.DeleteFile)

		protected.POST("/analyze", queriesHandler.CreateQueries)
		protected.GET("/analyze/history", queriesHandler.GetQueriesByUserId)
	}

	r.Run(":8080")

}
