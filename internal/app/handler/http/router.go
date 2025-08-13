package http

import (
	"case-management/internal/app/handler"

	"github.com/gin-gonic/gin"
)

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
			authRoutes.POST("/logout", handler.ValidateToken(), H.Auth.Logout)
			authRoutes.GET("/profile", handler.ValidateToken(), H.Auth.Profile)
		}

		userRoutes := apiV1.Group("/users")
		// userRoutes.Use(handler.ValidateToken())
		{
			userRoutes.POST("", H.User.CreateUser)
			userRoutes.GET("/", H.User.GetAllUsers)
			userRoutes.GET("/:id", H.User.GetUserByID)
			userRoutes.PUT("/:id", H.User.UpdateUserByID)
		}

		masterDataRoutes := apiV1.Group("/master-data")
		// masterDataRoutes.Use(handler.ValidateToken())
		{
			masterDataRoutes.GET("/lookups", H.MasterData.GetAllLookups)
		}

		permissionRoutes := apiV1.Group("/permissions")
		// permissionRoutes.Use(handler.ValidateToken())
		{
			permissionRoutes.GET("/", H.Permission.GetAllPermissions)
			permissionRoutes.PATCH("/update", H.Permission.UpdatePermission)
		}

		logRoutes := apiV1.Group("/logs")
		logRoutes.Use(handler.ValidateToken())
		{
			logRoutes.GET("/", H.Log.GetAllApiLogs)
		}

		caseManagementRoutes := apiV1.Group("/case-management")
		caseManagementRoutes.Use(handler.ValidateToken())
		{
			caseManagementRoutes.POST("/cases", H.Case.CreateCase)
			caseManagementRoutes.GET("/cases", H.Case.GetAllCases)
			caseManagementRoutes.GET("/cases/:id", H.Case.GetCaseByID)
		}

		customerRoutes := apiV1.Group("/customers")
		customerRoutes.Use(handler.ValidateToken())
		{
			customerRoutes.GET("/search/:id", H.Customer.SearchByCustomerId)
			customerRoutes.POST("/note", H.Customer.CreateCustomerNote)
			customerRoutes.GET("/:customerId/notes", H.Customer.GetAllCustomerNotes)
		}
	}
}
