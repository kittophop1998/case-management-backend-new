package http

import "github.com/gin-gonic/gin"

func SetupRoutes(
	router *gin.Engine,
) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	apiV1 := router.Group("/api/v1")
	{
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/login", H.Auth.Login)
			authRoutes.POST("/logout", H.Auth.Logout)
			authRoutes.POST("/profile", H.Auth.Profile)
		}

		userRoutes := apiV1.Group("/users")
		{
			userRoutes.GET("/", H.User.GetAllUsers)
			userRoutes.GET("/:id", H.User.GetUserByID)
			userRoutes.PUT("/:id", H.User.UpdateUserByID)
		}

		masterDataRoutes := apiV1.Group("/master-data")
		{
			masterDataRoutes.GET("/lookups", H.MasterData.GetAllLookups)
		}

		permissionRoutes := apiV1.Group("/permissions")
		{
			permissionRoutes.GET("/", H.Permission.GetAllPermissions)
			permissionRoutes.PATCH("/update", H.Permission.UpdatePermission)
		}

		logRoutes := apiV1.Group("/logs")
		{
			logRoutes.GET("/", H.Log.GetAllApiLogs)
		}

		// caseManagementRoutes := apiV1.Group("/case-management")
		// {

		// }
	}

}
