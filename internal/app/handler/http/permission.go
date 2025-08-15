package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionHandler struct {
	UseCase usecase.PermissionUseCase
}

func (h *PermissionHandler) GetAllPermissions(ctx *gin.Context) {
	limit, err := getLimit(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	page, err := getPage(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	permissionName := ctx.Query("name")
	var sectionID, departmentID *uuid.UUID

	if sid := ctx.Query("section_id"); sid != "" {
		if parsed, err := uuid.Parse(sid); err == nil {
			sectionID = &parsed
		}
	}

	if did := ctx.Query("department_id"); did != "" {
		if parsed, err := uuid.Parse(did); err == nil {
			departmentID = &parsed
		}
	}

	permissions, total, err := h.UseCase.GetAllPermissions(ctx, page, limit, permissionName, sectionID, departmentID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if permissions == nil {
		permissions = []model.PermissionWithRolesResponse{}
	}

	lib.HandlePaginatedResponse(ctx, page, limit, total, permissions)
}

func (h *PermissionHandler) UpdatePermission(ctx *gin.Context) {
	departmentId := ctx.Query("departmentId")

	var deptUUID, secUUID uuid.UUID
	if deptId := ctx.Query("departmentId"); deptId != "" {
		if parsed, err := uuid.Parse(departmentId); err == nil {
			deptUUID = parsed
		}
	}

	if secId := ctx.Query("SectionId"); secId != "" {
		if parsed, err := uuid.Parse(secId); err == nil {
			secUUID = parsed
		}
	}

	var reqs []model.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UseCase.UpdatePermission(ctx, reqs, deptUUID, secUUID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Permission updated successfully"})
}
