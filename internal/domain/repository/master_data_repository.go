package repository

import "github.com/gin-gonic/gin"

type MasterDataRepository interface {
	FindAll(ctx *gin.Context) (map[string]interface{}, error)
}
