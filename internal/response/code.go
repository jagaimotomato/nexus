package response

// 定义业务状态码
const (
	Success = 0
	Error   = 500

	InvalidParams   = 400
	AuthFailed      = 401
	Forbidden       = 403 // 新增：无权限/黑名单
	NotFound        = 404
	TooManyRequests = 429 // 新增：限流
)

// MsgFlags 状态码对应的默认提示信息
var MsgFlags = map[int]string{
	Success:         "操作成功",
	Error:           "操作失败",
	InvalidParams:   "请求参数错误",
	AuthFailed:      "身份认证失败",
	Forbidden:       "无权限访问",
	NotFound:        "未找到相关资源",
	TooManyRequests: "请求过于频繁，请稍后再试",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Error]
}