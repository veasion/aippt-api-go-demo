package main

import (
	"fmt"
	"os"
)

func main() {
	// PPT 流式生成
	// go run aippt_demo1.go api.go http_utils.go
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

	// 生成大纲
	fmt.Println("\n\n========== 正在生成大纲 ==========")
	outline, err := GenerateOutline(apiToken, subject, "", "")
	if err != nil {
		fmt.Println("异常:", err)
		return
	}

	// 生成大纲内容
	fmt.Println("\n\n========== 正在生成大纲内容 ==========")
	markdown, err := GenerateContent(apiToken, outline, "", "")
	if err != nil {
		fmt.Println("异常:", err)
		return
	}

	// 随机一个模板
	fmt.Println("\n\n========== 随机选择模板 ==========")
	templateId, err := RandomOneTemplateId(apiToken)
	if err != nil {
		fmt.Println("异常:", err)
		return
	}
	fmt.Println("templateId: " + templateId)

	// 生成PPT
	fmt.Println("\n\n========== 正在生成PPT ==========")
	pptInfo, err := GeneratePptx(apiToken, templateId, markdown, false)
	if err != nil {
		fmt.Println("异常:", err)
		return
	}

	pptId := pptInfo["id"].(string)
	fmt.Println("pptId: " + pptId)
	fmt.Println("ppt主题：" + pptInfo["subject"].(string))
	fmt.Println("ppt封面：" + pptInfo["coverUrl"].(string) + "?token=" + apiToken)

	// 下载PPT
	fmt.Println("\n\n========== 正在下载PPT ==========")
	url, err := DownloadPptx(apiToken, pptId)
	if err != nil {
		fmt.Println("异常:", err)
		return
	}
	fmt.Println("ppt链接：" + url)
	savePath, _ := os.Getwd()
	savePath += "/" + pptId + ".pptx"
	err = Download(url, savePath)
	if err != nil {
		fmt.Println("异常:", err)
		return
	}
	fmt.Println("ppt下载完成，保存路径：" + savePath)
}
