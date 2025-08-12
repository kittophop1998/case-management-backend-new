package http

import (
	"case-management/pkg/monitoring"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	prom *monitoring.Prometheus,
) {

	router.GET("/api/healthz", HealthCheck)
	router.GET("/api/metrics", prom.Handler())

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

		dashboardRoutes := apiV1.Group("/dashboard")
		{
			dashboardRoutes.GET("/custinfo/:id", H.Dashboard.GetCustInfo)
			dashboardRoutes.GET("/custprofile/:id", H.Dashboard.GetCustProfile)
			dashboardRoutes.GET("/custsegment/:id", H.Dashboard.GetCustSegment)
			dashboardRoutes.GET("/custsuggestion/:id", H.Dashboard.GetCustSuggestion)
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
