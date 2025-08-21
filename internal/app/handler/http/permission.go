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
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid limit parameter"))
		return
	}

	page, err := getPage(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid page parameter"))
		return
	}

	permissionName := ctx.Query("keyword")
	var sectionID, departmentID *uuid.UUID

	if sid := ctx.Query("sectionId"); sid != "" {
		if parsed, err := uuid.Parse(sid); err == nil {
			sectionID = &parsed
		}
	}

	if did := ctx.Query("departmentId"); did != "" {
		if parsed, err := uuid.Parse(did); err == nil {
			departmentID = &parsed
		}
	}

	permissions, total, permRoleCount, err := h.UseCase.GetAllPermissions(ctx, page, limit, permissionName, sectionID, departmentID)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	if permissions == nil {
		permissions = []model.PermissionWithRolesResponse{}
	}

	lib.HandlePaginatedResponse(ctx, page, limit, total, permissions, lib.OptionField{PermRoleCount: &permRoleCount})
}

func (h *PermissionHandler) UpdatePermission(ctx *gin.Context) {
	var deptUUID, secUUID uuid.UUID
	if deptId := ctx.Query("departmentId"); deptId != "" {
		if parsed, err := uuid.Parse(deptId); err != nil {
			lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid department ID, err :"+err.Error()))
			return
		} else {
			deptUUID = parsed
		}
	}

	if secId := ctx.Query("sectionId"); secId != "" {
		if parsed, err := uuid.Parse(secId); err != nil {
			lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid section ID, err :"+err.Error()))
			return
		} else {
			secUUID = parsed
		}
	}

	var reqs []model.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&reqs); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	if err := h.UseCase.UpdatePermission(ctx, reqs, deptUUID, secUUID); err != nil {
		lib.HandleError(ctx, lib.CannotUpdate.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, gin.H{"message": "Permissions updated successfully"})
}
