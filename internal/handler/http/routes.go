package http

import (
	"github.com/gin-gonic/gin"
	"github.com/penonpaper/event-api/internal/metrics"
	"github.com/penonpaper/event-api/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Routes struct {
	UserHandler       *UserHandler
	EventHandler      *EventHandler
	AttendeeHandler   *AttendeeHandler
	AuthHandler       *AuthHandler
	MetricsHandler    *metrics.Metrics
	AccessTokenSecret string
}

func NewRoutes(userHandler *UserHandler, eventHandler *EventHandler, attendeeHandler *AttendeeHandler, authHandler *AuthHandler, metrics *metrics.Metrics, accessTokenSecret string) *Routes {
	return &Routes{
		UserHandler:       userHandler,
		EventHandler:      eventHandler,
		AttendeeHandler:   attendeeHandler,
		AuthHandler:       authHandler,
		MetricsHandler:    metrics,
		AccessTokenSecret: accessTokenSecret,
	}
}

func (r *Routes) RegisteredRoutes(router *gin.Engine) {

	router.Use(middleware.PrometheusMiddleware(r.MetricsHandler))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := router.Group("/api/v1")
	{
		public.POST("/auth/signup", r.UserHandler.SignUp)
		public.POST("/auth/signin", r.AuthHandler.SignIn)
		public.GET("/events", r.EventHandler.GetEvents)
		public.GET("/events/tag", r.EventHandler.FilterByTags)
		//public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(r.AccessTokenSecret))
	{
		protected.POST("/event", middleware.OrganizerMiddleware("admin", "manager"),
			r.EventHandler.CreateEvent)
		protected.DELETE("/event/:id", middleware.OrganizerMiddleware("admin", "manager"),
			r.EventHandler.DeleteEvent)

		protected.POST("/auth/refresh", r.AuthHandler.Refresh)
		protected.POST("/event/:id/register", r.AttendeeHandler.Register)
		protected.POST("/event/:id/unregister", r.AttendeeHandler.Unregister)
		protected.POST("/auth/logout", r.AuthHandler.Logout)
	}

}
