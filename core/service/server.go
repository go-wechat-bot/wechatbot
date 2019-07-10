package service

import (
	"time"
)
// 当前系统状态
const (
	STATUS_INIT = 0;
	STATUS_WAITING_FOR_LOGIN = 1;
	STATUS_LOGIN = 2;
)
// 记录微信登录凭证
type wechat_info_type struct{
	qrcode string
}
var wechat_info wechat_info_type

// 开始
func Start(){
	now_status := STATUS_INIT
	// 获取登录二维码
	wechat_info.qrcode = GetWechatQRLogin()
	// 开始轮训
	for range time.Tick(time.Second) {
		switch now_status{
		case STATUS_INIT:
			UpdateWechatQrcodeStatus()
		}
	}
}

