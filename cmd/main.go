package main

import (
	_ "insightly/docs"
	"insightly/internal/ai"
	"insightly/internal/db"
	"insightly/internal/handlers"
	"insightly/internal/middlewares"
	"insightly/internal/repositories"
	"insightly/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Insightly API
// @version 1.0
// @description AI-powered CSV analyzer REST API
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате "Bearer {access_token}"
func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Println("[WARN] .env file not found, relying on system environment variables")
	}
	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	aiClient := ai.NewOpenAIClient(&client)

	database, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	secret := os.Getenv("JWT_SECRET")

	//Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	limiter := redis_rate.NewLimiter(rdb)

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authRoutes := r.Group("/auth")
	{

		authRoutes.POST("/register", userHandler.Register)
		authRoutes.POST("/login", userHandler.Login)
		authRoutes.POST("/refresh", refreshTokensHandler.Refresh)
	}

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware(secret))
	protected.Use(middlewares.RateLimiter(limiter))
	{
		protected.POST("/files/upload", filesHandler.CreateFiles)
		protected.GET("/files", filesHandler.GetFilesByUserId)
		protected.DELETE("/files/:id", filesHandler.DeleteFile)

		protected.POST("/analyze", queriesHandler.CreateQueries)
		protected.GET("/analyze/history", queriesHandler.GetQueriesByUserId)
	}

	r.Run(":8080")

}
