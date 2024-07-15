package main

import (
	"fmt"
	"os"
)

func main() {
	// PPT 直接生成
	// go run aippt_demo2.go api.go http_utils.go
	// 官网 https://docmee.cn
	// 开放平台 https://docmee.cn/open-platform/api

	// 填写你的API-KEY
	apiKey := "YOUR API KEY"

	// 第三方用户ID（数据隔离）
	uid := "test"
	subject := "AI未来的发展"

	// 创建 api token (有效期2小时，建议缓存到redis，同一个 uid 创建时之前的 token 会在10秒内失效)
	apiToken, err := CreateApiToken(apiKey, uid, -1)
	if err != nil {
		fmt.Println("异常:", err)
		return
	}
	fmt.Println("api token: " + apiToken)

	// 生成PPT
	fmt.Println("\n\n========== 正在生成PPT ==========")
	pptInfo, err := DirectGeneratePptx(apiToken, true, "", subject, "", "", false)
	if err != nil {
		fmt.Println("异常:", err)
		return
	}

	pptId := pptInfo["id"].(string)
	fileUrl := pptInfo["fileUrl"].(string)
	fmt.Println("pptId: " + pptId)
	fmt.Println("ppt主题：" + pptInfo["subject"].(string))
	fmt.Println("ppt封面：" + pptInfo["coverUrl"].(string) + "?token=" + apiToken)
	fmt.Println("ppt链接：" + fileUrl)

	// 下载PPT
	savePath, _ := os.Getwd()
	savePath += "/" + pptId + ".pptx"
	err = Download(fileUrl, savePath)
	if err != nil {
		fmt.Println("异常:", err)
		return
	}
	fmt.Println("ppt下载完成，保存路径：" + savePath)
}
