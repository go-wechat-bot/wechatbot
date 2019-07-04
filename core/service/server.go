package service

import (
	"fmt"
)

func Start(){
	qrcode := GetWechatQRLogin()
	fmt.Println("二维码地址："+qrcode)
}