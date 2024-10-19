package server

import (
	"fmt"
	"gateway_mrc/config"
	db "gateway_mrc/db/sqlc"
	"gateway_mrc/infrastructure/jwt_token"
	"gateway_mrc/interfaces/controllers"
	"gateway_mrc/interfaces/middleware"
	"gateway_mrc/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Server struct {
	config     config.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker jwt_token.Maker
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	tokenMaker, err := jwt_token.NewJWTMaker(
		config.TokenSymmetricKey,
		config.AccessTokenExpiresIn,
		config.RefreshTokenExpiresIn,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Укажите домен вашего клиента
		AllowCredentials: true,                              // Разрешить использование учетных данных (например, куки)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
	}
	router.Use(cors.New(corsConfig))

	// SESSION
	sessionMiddleware := middleware.SessionMiddleware(server.store)
	router.Use(sessionMiddleware)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Route %s not found", ctx.Request.URL)})
	})
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	//
	sessionUsecase := usecase.NewSessionUsecase(server.store)
	controller := controllers.NewProxyController(server.tokenMaker, sessionUsecase, server.config)
	server.setupAuthRoutes(router, controller)
	server.setupUsersRoutes(router, controller)
	server.router = router
}

func (server *Server) setupAuthRoutes(router *gin.Engine, controller controllers.ProxyController) {
	authMiddleware := middleware.AuthMiddleware(server.tokenMaker, server.store, server.config)
	public := router.Group("/api/v1/auth/public")
	{
		public.POST("/sign-up", func(ctx *gin.Context) {
			controller.ProxyCommonReq(ctx, server.config.HttpAuthAddress)
		})
		public.POST("/sign-in", func(ctx *gin.Context) {
			controller.ProxySignInReq(ctx, server.config.HttpAuthAddress)
		})
	}

	private := router.Group("/api/v1/auth/private")
	private.Use(authMiddleware)
	{
		private.GET("/logout", func(ctx *gin.Context) {
			controller.ProxyLogoutReq(ctx)
		})
		private.POST("/*path", func(ctx *gin.Context) {
			controller.ProxyCommonReq(ctx, server.config.HttpAuthAddress)
		})
		private.DELETE("/*path", func(ctx *gin.Context) {
			controller.ProxyCommonReq(ctx, server.config.HttpAuthAddress)
		})
		private.PUT("/*path", func(ctx *gin.Context) {
			controller.ProxyCommonReq(ctx, server.config.HttpAuthAddress)
		})
	}
}

func (server *Server) setupUsersRoutes(router *gin.Engine, controller controllers.ProxyController) {
	authMiddleware := middleware.AuthMiddleware(server.tokenMaker, server.store, server.config)
	private := router.Group("/api/v1/users/private")
	address := server.config.HttpAuthAddress
	private.Use(authMiddleware)
	{
		registerRoutes(private, controller, address)
	}
}

func registerRoutes(group *gin.RouterGroup, controller controllers.ProxyController, address string) {
	group.GET("/*path", func(ctx *gin.Context) {
		controller.ProxyCommonReq(ctx, address)
	})
	group.POST("/*path", func(ctx *gin.Context) {
		controller.ProxyCommonReq(ctx, address)
	})
	group.DELETE("/*path", func(ctx *gin.Context) {
		controller.ProxyCommonReq(ctx, address)
	})
	group.PUT("/*path", func(ctx *gin.Context) {
		controller.ProxyCommonReq(ctx, address)
	})
}

func (server *Server) Start(address string) error {
	log.Printf("Starting server on %s\n", address)
	return server.router.Run(address)
}

func RunGinServer(config config.Config, store db.Store) error {
	server, err := NewServer(config, store)
	if err != nil {
		return err
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		return err
	}
	return nil
}
