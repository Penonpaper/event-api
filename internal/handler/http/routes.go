package http

import (
	"github.com/gin-gonic/gin"
	"github.com/penonpaper/event-api/internal/middleware"
)

type Routes struct {
	UserHandler  *UserHandler
	EventHandler *EventHandler
	Jwtsecret    string
}

func NewRoutes(userHandler *UserHandler, eventHandler *EventHandler, jwttoken string) *Routes {
	return &Routes{
		UserHandler:  userHandler,
		EventHandler: eventHandler,
		Jwtsecret:    jwttoken,
	}
}

func (r *Routes) RegisteredRoutes(router *gin.Engine) {

	public := router.Group("/api/v1")
	{
		public.POST("/auth/signup", r.UserHandler.SignUp)
		public.POST("/auth/signin", r.UserHandler.SignIn)
		public.GET("/events", r.EventHandler.GetEvents)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(r.Jwtsecret))
	{
		protected.POST("/events", middleware.OrganizerMiddleware("admin", "manager"),
			r.EventHandler.CreateEvent)
		protected.DELETE("/event/:id", middleware.OrganizerMiddleware("admin", "manager"),
			r.EventHandler.DeleteEvent)
	}

}
