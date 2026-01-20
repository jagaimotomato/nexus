package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// Result 基础方法 (保持不变)
func Result(c *gin.Context, httpStatus int, code int, msg string, data interface{}) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// ==========================================
// 成功响应 (利用 GetMsg(Success))
// ==========================================

// OK 普通成功
func OK(c *gin.Context) {
    // 自动获取 "success" 或 "操作成功"
	Result(c, http.StatusOK, Success, GetMsg(Success), nil)
}

// OKWithData 带数据的成功
func OKWithData(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, Success, GetMsg(Success), data)
}

// OKWithList 带分页的成功
func OKWithList(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	OKWithData(c, PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// OKWithMessage 成功，但覆盖默认的 msg
func OKWithMessage(c *gin.Context, msg string) {
	Result(c, http.StatusOK, Success, msg, nil)
}

// ==========================================
// 失败响应 (利用 GetMsg(code))
// ==========================================

// Fail 通用失败 (使用默认的 500 错误文案)
func Fail(c *gin.Context) {
    // 自动获取 "操作失败"
	Result(c, http.StatusOK, Error, GetMsg(Error), nil)
}

// FailWithMessage 失败，自定义消息 (依然是用 Error 状态码)
func FailWithMessage(c *gin.Context, msg string) {
	Result(c, http.StatusOK, Error, msg, nil)
}

// FailWithCode 【核心修改】根据错误码自动获取文案
// 场景：token 过期，你只需要传 response.AuthFailed，文案自动变成 "身份认证失败"
func FailWithCode(c *gin.Context, code int) {
	Result(c, http.StatusOK, code, GetMsg(code), nil)
}

// FailWithDetailed 既要指定错误码，又要自定义文案
func FailWithDetailed(c *gin.Context, code int, msg string) {
	Result(c, http.StatusOK, code, msg, nil)
}