package handler

import (
	"nexus/internal/response"
	"nexus/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	svc *service.MenuService
}

func NewMenuHandler(s *service.MenuService) *MenuHandler {
	return &MenuHandler{svc: s}
}

func (h *MenuHandler) RegisterPublic(r *gin.RouterGroup) {}

func (h *MenuHandler) RegisterPrivate(r *gin.RouterGroup) {
	menu := r.Group("/menus")
	{
		menu.GET("", h.GetList)
		menu.POST("", h.Create)
		menu.PUT("/:id", h.Update)
		menu.DELETE("/:id", h.Delete)
	}
}

func (h *MenuHandler) GetList(c *gin.Context) {
	menus, err := h.svc.GetMenuTree()
	if err != nil {
		response.FailWithMessage(c, "获取菜单失败")
		return
	}
	response.OKWithData(c, menus)
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req service.MenuInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.InvalidParams)
		return
	}
	err := h.svc.CreateMenu(&req)
	if err != nil {
		response.FailWithMessage(c, "创建失败："+err.Error())
		return
	}
	response.OK(c)
}

func (h *MenuHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var req service.MenuInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.InvalidParams)
		return
	}
	err := h.svc.UpdateMenu(uint(id), &req)
	if err != nil {
		response.FailWithMessage(c, "更新失败："+err.Error())
		return
	}
	response.OK(c)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	err := h.svc.DeleteMenu(uint(id))
	if err != nil {
		response.FailWithMessage(c, "删除失败："+err.Error())
		return
	}
	response.OK(c)
}
