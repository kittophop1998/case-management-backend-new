package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/handler"
	"encoding/json"
	"net/http"
	"time"

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
			userRoutes.GET("", H.User.GetAllUsers)
			userRoutes.GET("/:id", H.User.GetUserByID)
			userRoutes.PUT("/:id", H.User.UpdateUserByID)
		}

		masterDataRoutes := apiV1.Group("/master-data")
		{
			masterDataRoutes.GET("/lookups", H.MasterData.GetAllLookups)
		}

		permissionRoutes := apiV1.Group("/permissions")
		// permissionRoutes.Use(handler.ValidateToken())
		{
			permissionRoutes.GET("", H.Permission.GetAllPermissions)
			permissionRoutes.PATCH("/update", H.Permission.UpdatePermission)
		}

		logRoutes := apiV1.Group("/logs")
		logRoutes.Use(handler.ValidateToken())
		{
			logRoutes.GET("", H.Log.GetAllApiLogs)
		}

		caseManagementRoutes := apiV1.Group("/cases")
		caseManagementRoutes.Use(handler.ValidateToken())
		{
			caseManagementRoutes.POST("", H.Case.CreateCase)
			caseManagementRoutes.GET("", H.Case.GetAllCases)
			caseManagementRoutes.GET("/:id", H.Case.GetCaseByID)
			caseManagementRoutes.GET("/disposition", H.Case.GetAllDisposition)
		}

		customerRoutes := apiV1.Group("/customers")
		customerRoutes.Use(handler.ValidateToken())
		{
			customerRoutes.POST("/note", H.Customer.CreateCustomerNote)
			customerRoutes.GET("/:customerId/notes", H.Customer.GetAllCustomerNotes)
			customerRoutes.GET("/note-types", H.Customer.GetNoteTypes)
		}

		dashboardRoutes := apiV1.Group("/dashboard")
		// userRoutes.Use(handler.ValidateToken())
		{
			dashboardRoutes.GET("/custinfo/:id", H.Dashboard.GetCustInfo)
			dashboardRoutes.GET("/custprofile/:id", H.Dashboard.GetCustProfile)
			dashboardRoutes.GET("/custsegment/:id", H.Dashboard.GetCustSegment)
			dashboardRoutes.GET("/custsuggestion/:id", H.Dashboard.GetCustSuggestion)
		}

		// TODO: delete
		apiV1.GET("/mock/*any", func(c *gin.Context) {
			if c.Query("isError") != "" && c.Query("isError") != "false" {
				time.Sleep(2 * time.Second)
				lib.HandleError(c, lib.InternalServer)
				return
			}
			datamock := c.Query("datamock")
			var parsedData map[string]interface{}
			if err := json.Unmarshal([]byte(datamock), &parsedData); err != nil {
				lib.HandleError(c, lib.InternalServer)
				return
			}
			time.Sleep(2 * time.Second)
			c.JSON(http.StatusOK, gin.H{
				"data": parsedData,
			})
		})
	}
}
