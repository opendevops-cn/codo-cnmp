package xerr

var message map[ErrCode]string

func init() {
	message = make(map[ErrCode]string)
	message[OK] = "SUCCESS"
	message[ServerCommonError] = "操作失败，请稍后再试"
	message[RequestParamError] = "参数错误"
	message[TokenExpireError] = "token失效，请重新登录"
	message[ErrTokenInvalid] = "token无效"
	message[TokenGenerateError] = "生成token失败"
	message[DbError] = "数据库繁忙,请稍后再试"
	message[DbUpdateAffectedZeroError] = "更新数据影响行数为0"
	message[ErrNotAllowed] = "访问未授权"
	message[ErrResourceNotFound] = "资源对象不存在"
	message[ErrItemExist] = "对象已存在"
	message[ErrItemCantEditOrDel] = "对象不可编辑"
	message[ErrUnAuthorization] = "没有权限, 请登录之后操作"
	message[ErrHttpRequestError] = "http请求返回异常"

	message[ErrPasswordOrEmail] = "用户名或密码错误"
	message[ErrPassword] = "密码错误"
	message[ErrUserForbidden] = "用户已经被禁用"
	message[ErrInnerUserNotChanPwd] = "内部用户不支持修改密码"
	message[ErrUserNotFound] = "用户不存在"
}

func MapErrMsg(errcode ErrCode) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}
