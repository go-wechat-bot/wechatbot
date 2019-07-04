package service

import (
	"net/http"
	"time"
	"bytes"
	"io"
	"regexp"
)

const(
	JSLOGIN		= "https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=1476606163580"
)

func Get(url string) (response string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, error := client.Get(url)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	response = result.String()
	return
}


/**
	获取二维码认证code
 */
func GetWechatQRLogin() string{
	data := Get(JSLOGIN)
	result := regexp.MustCompile(`([0-9]+).*"([a-zA-Z0-9=]+)"`).FindAllStringSubmatch(data,-1) //匹配
	if(len(result) == 0 || result[0][1] != "200"){
		panic("获取扫码失败")
	}
	return "https://login.weixin.qq.com/qrcode/"+result[0][2];
}