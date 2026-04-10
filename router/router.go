package router

import (
	"net/http"
	"sdk/handler"
	"sdk/repository"
	"sdk/service"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	repo := repository.NewRepository()
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
	
	r.GET("/taxbandits/auth", h.GetAccessToken)
	
	r.POST("/businesses", h.CreateBusiness)
	r.PUT("/businesses", h.UpdateBusiness)
	r.GET("/businesses", h.ListBusinesses)
	r.DELETE("/businesses", h.DeleteBusiness)
	r.GET("/businesses/details", h.GetBusiness)

	r.GET("/businesses/request-url", h.RequestByURL)
	r.GET("/whcertificate/list", h.ListWhCertificate)
	r.POST("/whcertificate/request-email", h.RequestByEmail)

	r.GET("/recipient/list", h.ListRecipient)

	r.GET("/pdf", h.GetPDF)

	return r
}
