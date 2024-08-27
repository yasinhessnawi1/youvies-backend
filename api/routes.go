package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// User Endpoints
	r.POST("/youvies/v1/api/register", RegisterUser)
	r.POST("/youvies/v1/api/login", LoginUser)
	r.POST("/youvies/v1/api/logout", AuthMiddleware("user"), LogoutUser)
	r.PUT("/youvies/v1/api/user", AuthMiddleware("user"), EditUser)

}
