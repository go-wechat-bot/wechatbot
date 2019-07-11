package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const(
	JSLOGIN		= "https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=1476606163580";
	LOGINSTATUS	= "https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid="
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
	result := regexp.MustCompile(`([0-9]+).*"([a-zA-Z0-9=\-\_]+)"`).FindAllStringSubmatch(data,-1) //匹配
	if(len(result) == 0 || result[0][1] != "200"){
		fmt.Println(data)
		fmt.Println(result)
		panic("获取扫码失败")
	}
	fmt.Println("https://login.weixin.qq.com/qrcode/"+result[0][2])
	return result[0][2];
}

// 查询微信二维码状态
func UpdateWechatQrcodeStatus(){
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	url := (LOGINSTATUS+wechat_info.qrcode+"&tip=0&r=862560455&_="+fmt.Sprintf("%d",time.Now().UnixNano() / 1000000))
	rsp, _ := client.Get(url)
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("myHttpGet error is ", err)
		return
	}
	checkWechatLoginCode(string(body))
}

// 检测扫码状态
func checkWechatLoginCode(content string){
	split := strings.Split(content, ";")
	if code := strings.Split(split[0],"=") ; code[1] == "201"{
		// 保存头像
		writeWechatAvatar(split[1]+";"+split[2])
	}else if(code[1] == "200"){
		// 登录成功
		fmt.Println(split)
	}
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func writeWechatAvatar(base64 string){
	data := strings.Replace(strings.Replace(base64, "window.userAvatar = '", "", 1),"'","",1)

	fmt.Println(WriteFile("runtime/tmp", data))
}

func markDir(_dir string){
	exist, err := PathExists(_dir) ;
	if  err != nil{
		fmt.Printf("get dir error![%v]\n", err)
		return
	}
	os.Remove("runtime/tmp/avatar.jpg")
	if exist {
		fmt.Printf("has dir![%v]\n", _dir)
	} else {
		fmt.Printf("no dir![%v]\n", _dir)
		// 创建文件夹
		err := os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
}

//写入文件,保存
func WriteFile(path string, base64_image_content string) bool {
	b, _ := regexp.MatchString(`^data:\s*img\/(\w+);base64,`, base64_image_content)
	if !b {
		return false
	}
	re, _ := regexp.Compile(`^data:\s*img\/(\w+);base64,`)
	base64Str := re.ReplaceAllString(base64_image_content, "")
	if  ok, _ := PathExists(path); !ok {
		os.Mkdir(path, 0666)
	}
	file := path  + "/avatar.jpg"
	byte, _ := base64.StdEncoding.DecodeString(base64Str)

	fmt.Println(path  + "/avatar.jpg")
	err := ioutil.WriteFile(file, byte, 0666)
	if err != nil {
		log.Println(err)
	}

	return false
}
