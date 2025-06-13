package filter

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin == "http://ec2-18-224-124-8.us-east-2.compute.amazonaws.com:3002" || origin == "http://localhost:3002" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Cookie, Content-Length, Origin, cache-control, X-Requested-With, Content-Type, Accept, Authorization, Token, PanGu-Token, Timestamp, UserId, x-request-id, apifoxtoken")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Authorization, New-Expires-At, PanGu-Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	}
}
