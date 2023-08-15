package common

// 用于service层与controller层之间的错误传递
type LError struct {
	HttpCode int32  // 状态码
	Msg      string // 返回状态描述
	Err      error  // 错误
}
