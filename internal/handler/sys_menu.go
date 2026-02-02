package handler

import (
	"nexus/internal/data"
	"nexus/internal/response"
	"nexus/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct{}

func (h *MenuHandler) GetList(c *gin.Context) {
	menus, err := service.GetMenuTree()
	if err != nil {
		response.FailWithMessage(c, "获取菜单失败")
		return
	}
	response.OKWithData(c, menus)
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req data.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.InvalidParams)
		return
	}
	err := service.CreateMenu(&req)
	if err != nil {
		response.FailWithMessage(c, "创建失败：" + err.Error())
		return
	}
	response.OK(c)
}

func (h *MenuHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var req data.Menu
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.InvalidParams)
		return
	}
	err := service.UpdateMenu(uint(id), &req)
	if err != nil {
		response.FailWithMessage(c, "更新失败：" + err.Error())
		return
	}
	response.OK(c)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	err := service.DeleteMenu(uint(id))
	if err != nil {
		response.FailWithMessage(c, "删除失败：" + err.Error())
		return
	}
	response.OK(c)
}