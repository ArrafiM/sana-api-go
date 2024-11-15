package routes

import (
	"github.com/gin-gonic/gin"
)

func IndexRoutes(r *gin.Engine) {
	AlbumRoutes(r)
	RoleRoute(r)
	UserRoute(r)
	AuthRoute(r)
	LocationRoutes(r)
	MerchantRoute(r)
	MerchandiseRoute(r)
	//websocketchat
	WebscoketRoute(r)
	//--
	ChatRoutes(r)
	DeviceTokenRoutes(r)

}
