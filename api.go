package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const BaseUrl = "https://chatmee.cn"

func CreateApiToken(apiKey string, uid string, limit int) (string, error) {
	url := BaseUrl + "/api/user/createApiToken"
	reqData := map[string]interface{}{
		"uid": uid,
	}
	if limit > 0 {
		reqData["limit"] = limit
	}
	headers := map[string]string{
		"Api-Key": apiKey,
	}
	body, _ := json.Marshal(reqData)
	resp, err := PostJson(url, headers, string(body))
	if err != nil {
		return "", err
	}
	if resp.status != 200 {
		return "", errors.New(fmt.Sprintf("创建apiToken失败，httpStatus=%d", resp.status))
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(resp.text), &jsonData)
	if err != nil {
		return "", err
	}
	if code, ok := jsonData["code"]; ok && code.(float64) != 0 {
		return "", errors.New(fmt.Sprintf("创建apiToken异常: %s", jsonData["message"]))
	}
	return jsonData["data"].(map[string]interface{})["token"].(string), nil
}

func GenerateOutline(apiToken string, subject string, prompt string, dataUrl string) (string, error) {
	url := BaseUrl + "/api/ppt/generateOutline"
	reqData := map[string]interface{}{
		"subject": subject,
	}
	if prompt != "" {
		reqData["prompt"] = prompt
	}
	if dataUrl != "" {
		reqData["dataUrl"] = dataUrl
	}
	headers := map[string]string{
		"token": apiToken,
	}
	body, _ := json.Marshal(reqData)
	var sb []string
	consumer := func(data string) error {
		var jsonData map[string]interface{}
		err := json.Unmarshal([]byte(data), &jsonData)
		if err != nil {
			return err
		}
		if status, ok := jsonData["status"]; ok && status.(float64) == -1 {
			return errors.New(fmt.Sprintf("请求异常: %s", jsonData["error"]))
		}
		if text, ok := jsonData["text"].(string); ok {
			sb = append(sb, text)
			fmt.Print(text)
		}
		return nil
	}
	resp, err := PostSse(url, headers, string(body), consumer)
	if err != nil {
		return "", err
	}
	if resp.status != 200 {
		err = errors.New(fmt.Sprintf("生成大纲失败，httpStatus=%d", resp.status))
		return "", err
	}
	if strings.Contains(resp.contentType, "application/json") {
		var jsonData map[string]interface{}
		err := json.Unmarshal([]byte(resp.text), &jsonData)
		if err != nil {
			return "", err
		}
		err = errors.New(fmt.Sprintf("生成大纲异常: %s", jsonData["message"]))
		return "", err
	}
	return strings.Join(sb, ""), nil
}

func GenerateContent(apiToken string, outlineMarkdown string, prompt string, dataUrl string) (string, error) {
	url := BaseUrl + "/api/ppt/generateContent"
	reqData := map[string]interface{}{
		"outlineMarkdown": outlineMarkdown,
	}
	if prompt != "" {
		reqData["prompt"] = prompt
	}
	if dataUrl != "" {
		reqData["dataUrl"] = dataUrl
	}
	headers := map[string]string{
		"token": apiToken,
	}
	body, _ := json.Marshal(reqData)
	var sb []string
	consumer := func(data string) error {
		var jsonData map[string]interface{}
		err := json.Unmarshal([]byte(data), &jsonData)
		if err != nil {
			return err
		}
		if status, ok := jsonData["status"]; ok && status.(float64) == -1 {
			return errors.New(fmt.Sprintf("请求异常: %s", jsonData["error"]))
		}
		if text, ok := jsonData["text"].(string); ok {
			sb = append(sb, text)
			fmt.Print(text)
		}
		return nil
	}
	resp, err := PostSse(url, headers, string(body), consumer)
	if err != nil {
		return "", err
	}
	if resp.status != 200 {
		err = errors.New(fmt.Sprintf("生成大纲内容失败，httpStatus=%d", resp.status))
		return "", err
	}
	if strings.Contains(resp.contentType, "application/json") {
		var jsonData map[string]interface{}
		err := json.Unmarshal([]byte(resp.text), &jsonData)
		if err != nil {
			return "", err
		}
		err = errors.New(fmt.Sprintf("生成大纲内容异常: %s", jsonData["message"]))
		return "", err
	}
	return strings.Join(sb, ""), nil
}

func RandomOneTemplateId(apiToken string) (string, error) {
	url := BaseUrl + "/api/ppt/randomTemplates"
	reqData := map[string]interface{}{
		"size": 1,
		"filters": map[string]interface{}{
			"type": 1,
		},
	}
	headers := map[string]string{
		"token": apiToken,
	}
	body, _ := json.Marshal(reqData)
	resp, err := PostJson(url, headers, string(body))
	if err != nil {
		return "", err
	}
	if resp.status != 200 {
		return "", errors.New(fmt.Sprintf("获取模板失败，httpStatus=%d", resp.status))
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(resp.text), &jsonData)
	if err != nil {
		return "", err
	}
	if code, ok := jsonData["code"]; ok && code.(float64) != 0 {
		return "", errors.New(fmt.Sprintf("获取模板异常: %s", jsonData["message"]))
	}
	return jsonData["data"].([]interface{})[0].(map[string]interface{})["id"].(string), nil
}

func GeneratePptx(apiToken string, templateId string, markdown string, pptxProperty bool) (map[string]interface{}, error) {
	url := BaseUrl + "/api/ppt/generatePptx"
	reqData := map[string]interface{}{
		"outlineContentMarkdown": markdown,
		"pptxProperty":           pptxProperty,
	}
	if templateId != "" {
		reqData["templateId"] = templateId
	}
	headers := map[string]string{
		"token": apiToken,
	}
	body, _ := json.Marshal(reqData)
	resp, err := PostJson(url, headers, string(body))
	if err != nil {
		return nil, err
	}
	if resp.status != 200 {
		return nil, errors.New(fmt.Sprintf("生成PPT失败，httpStatus=%d", resp.status))
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(resp.text), &jsonData)
	if err != nil {
		return nil, err
	}
	if code, ok := jsonData["code"]; ok && code.(float64) != 0 {
		return nil, errors.New(fmt.Sprintf("生成PPT异常: %s", jsonData["message"]))
	}
	return jsonData["data"].(map[string]interface{})["pptInfo"].(map[string]interface{}), nil
}

func DownloadPptx(apiToken string, id string) (string, error) {
	url := BaseUrl + "/api/ppt/downloadPptx"
	reqData := map[string]interface{}{
		"id": id,
	}
	headers := map[string]string{
		"token": apiToken,
	}
	body, _ := json.Marshal(reqData)
	resp, err := PostJson(url, headers, string(body))
	if err != nil {
		return "", err
	}
	if resp.status != 200 {
		return "", errors.New(fmt.Sprintf("下载PPT失败，httpStatus=%d", resp.status))
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(resp.text), &jsonData)
	if err != nil {
		return "", err
	}
	if code, ok := jsonData["code"]; ok && code.(float64) != 0 {
		return "", errors.New(fmt.Sprintf("下载PPT异常: %s", jsonData["message"]))
	}
	return jsonData["data"].(map[string]interface{})["fileUrl"].(string), nil
}

func DirectGeneratePptx(apiToken string, stream bool, templateId string, subject string, prompt string, dataUrl string, pptxProperty bool) (map[string]interface{}, error) {
	url := BaseUrl + "/api/ppt/directGeneratePptx"
	reqData := map[string]interface{}{
		"stream":       stream,
		"pptxProperty": pptxProperty,
	}
	if templateId != "" {
		reqData["templateId"] = templateId
	}
	if subject != "" {
		reqData["subject"] = subject
	}
	if prompt != "" {
		reqData["prompt"] = prompt
	}
	if dataUrl != "" {
		reqData["dataUrl"] = dataUrl
	}
	headers := map[string]string{
		"token": apiToken,
	}
	body, _ := json.Marshal(reqData)
	if stream {
		var pptInfo []interface{}
		consumer := func(data string) error {
			var jsonData map[string]interface{}
			err := json.Unmarshal([]byte(data), &jsonData)
			if err != nil {
				return err
			}
			status, ok := jsonData["status"]
			if ok && status.(float64) == -1 {
				return errors.New(fmt.Sprintf("请求异常: %s", jsonData["error"]))
			}
			if ok && status.(float64) == 4 {
				result := jsonData["result"].(map[string]interface{})
				pptInfo = append(pptInfo, result)
			}
			if text, ok := jsonData["text"].(string); ok {
				fmt.Print(text)
			}
			return nil
		}
		resp, err := PostSse(url, headers, string(body), consumer)
		if err != nil {
			return nil, err
		}
		if resp.status != 200 {
			err = errors.New(fmt.Sprintf("生成PPT失败，httpStatus=%d", resp.status))
			return nil, err
		}
		if strings.Contains(resp.contentType, "application/json") {
			var jsonData map[string]interface{}
			err := json.Unmarshal([]byte(resp.text), &jsonData)
			if err != nil {
				return nil, err
			}
			err = errors.New(fmt.Sprintf("生成PPT异常: %s", jsonData["message"]))
			return nil, err
		}
		return pptInfo[0].(map[string]interface{}), nil
	} else {
		resp, err := PostJson(url, headers, string(body))
		if err != nil {
			return nil, err
		}
		if resp.status != 200 {
			return nil, errors.New(fmt.Sprintf("生成PPT失败，httpStatus=%d", resp.status))
		}
		var jsonData map[string]interface{}
		err = json.Unmarshal([]byte(resp.text), &jsonData)
		if err != nil {
			return nil, err
		}
		if code, ok := jsonData["code"]; ok && code.(float64) != 0 {
			return nil, errors.New(fmt.Sprintf("生成PPT异常: %s", jsonData["message"]))
		}
		pptInfo := jsonData["data"].(map[string]interface{})["pptInfo"].(map[string]interface{})
		return pptInfo, nil
	}
}
