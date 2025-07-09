package xerr

type ErrCode uint32

// OK 成功返回
const OK = ErrCode(200)

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码

const (
	ServerCommonError         = ErrCode(100001)
	RequestParamError         = ErrCode(100002)
	TokenExpireError          = ErrCode(100003)
	TokenGenerateError        = ErrCode(100004)
	DbError                   = ErrCode(100005)
	DbUpdateAffectedZeroError = ErrCode(100006)
	ErrNotAllowed             = ErrCode(100007)
	ErrTokenInvalid           = ErrCode(100008)
	ErrItemExist              = ErrCode(100009)
	ErrItemCantEditOrDel      = ErrCode(100010)
	ErrUnAuthorization        = ErrCode(100011)
	ErrPermissionDenied       = ErrCode(100403)
	ErrResourceNotFound       = ErrCode(100404)
	ErrHttpRequestError       = ErrCode(100405)
	ErrResourceGetError       = ErrCode(100406)
	ErrWebsocketNotCONNECT    = ErrCode(100407)
	DryRunResourceError       = ErrCode(100408)
)

// base模块-101

const (
	ErrPasswordOrEmail     = ErrCode(101001)
	ErrUserForbidden       = ErrCode(101002)
	ErrInnerUserNotChanPwd = ErrCode(101003)
	ErrUserNotFound        = ErrCode(101004)
	ErrPassword            = ErrCode(101005)
)
