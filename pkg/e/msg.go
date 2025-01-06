package e

var MsgFlags = map[int]string{
	SUCCESS:               "ok",
	UpdatePasswordSuccess: "修改密码成功",
	NotExistInentifier:    "该第三方账号未绑定",
	ERROR:                 "fail",
	InvalidParams:         "请求参数错误",
	ServerBusy:            "服务繁忙",

	ErrorExistNick:          "已存在该昵称",
	ErrorExistEmail:         "邮箱已注册",
	ErrorExistTelNum:        "手机号已注册",
	ErrorExistUser:          "已存在该用户名",
	ErrorNotExistUser:       "该用户不存在",
	ErrorNotCompare:         "密码错误",
	ErrorCheckCode:          "验证码错误",
	ErrorNotComparePassword: "两次密码输入不一致",
	ErrorComparePassword:    "新密码不能和旧密码相同",
	ErrorFailEncryption:     "加密失败",
	ErrorNotExistProduct:    "该商品不存在",
	ErrorNotExistAddress:    "该收获地址不存在",
	ErrorExistFavorite:      "已收藏该商品",
	ErrorUserNotFound:       "用户不存在",

	ErrorBossCheckTokenFail:        "商家的Token鉴权失败",
	ErrorBossCheckTokenTimeout:     "商家Token已超时",
	ErrorBossToken:                 "商家的Token生成失败",
	ErrorBoss:                      "商家Token错误",
	ErrorBossInsufficientAuthority: "商家权限不足",
	ErrorBossProduct:               "商家读文件错误",

	ErrorProductExistCart: "商品已经在购物车了，数量+1",
	ErrorProductMoreCart:  "超过最大上限",

	ErrorAuthCheckTokenFail:        "Token鉴权失败",
	ErrorAuthCheckTokenTimeout:     "Token已超时",
	ErrorAuthToken:                 "Token生成失败",
	ErrorAuth:                      "Token错误",
	ErrorAuthInsufficientAuthority: "权限不足",
	ErrorReadFile:                  "读文件失败",
	ErrorSendEmail:                 "发送邮件失败",
	ErrorCallApi:                   "调用接口失败",
	ErrorUnmarshalJson:             "解码JSON失败",

	ErrorUploadFile:    "上传失败",
	ErrorAdminFindUser: "管理员查询用户失败",

	ErrorDatabase: "数据库操作出错,请重试",

	ErrorOss: "OSS配置错误",

	WebsocketSuccessMessage: "解析content内容信息",
	WebsocketSuccess:        "发送信息，请求历史纪录操作成功",
	WebsocketEnd:            "请求历史纪录，但没有更多记录了",
	WebsocketOnlineReply:    "针对回复信息在线应答成功",
	WebsocketOfflineReply:   "针对回复信息离线回答成功",
	WebsocketLimit:          "请求收到限制",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
