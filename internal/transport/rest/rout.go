package rest

import (
	"github.com/gin-gonic/gin"
)

type handlers interface {
	AddUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetUsersByNationality(c *gin.Context)
	GetUsersByAge(c *gin.Context)
	GetUsersByGender(c *gin.Context)
	GetUsersByName(c *gin.Context)
}

type Router struct {
	router   *gin.Engine
	handlers handlers
}

func NewRouter(handlers handlers) *Router {
	router := gin.Default()

	return &Router{router: router, handlers: handlers}
}

func (r *Router) Run(addr string) error {
	r.router.POST("/add", r.handlers.AddUser)
	r.router.DELETE("/delete/:id", r.handlers.DeleteUser)
	r.router.PUT("/update/:id", r.handlers.UpdateUser)
	r.router.GET("/get/id/:id", r.handlers.GetUserByID)
	r.router.GET("/get/nationality/:nationality", r.handlers.GetUsersByNationality)
	r.router.GET("/get/age/:age", r.handlers.GetUsersByAge)
	r.router.GET("/get/gender/:gender", r.handlers.GetUsersByGender)
	r.router.GET("/get/name/:name", r.handlers.GetUsersByName)
	r.router.GET("/get/all", r.handlers.GetAllUsers)

	r.router.Use(gin.Logger())
	r.router.Use(gin.Recovery())
	r.router.Use(gin.ErrorLogger())

	return r.router.Run(addr)
}
