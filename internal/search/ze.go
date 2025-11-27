package search

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/itihey/tikuAdapter/pkg/errors"
	"github.com/itihey/tikuAdapter/pkg/model"
	"github.com/itihey/tikuAdapter/pkg/util"
	"time"
)

// ZE API 的响应结构体
type zeResponseData struct {
	Question string `json:"question"` // 原始题目
	Answer   string `json:"answer"`   // 答案文本
	IsAI     bool   `json:"is_ai"`
}

type zeResponse struct {
	Code int            `json:"code"`
	Msg  string         `json:"message"` // 对应 ZE API 的 message 字段
	Data zeResponseData `json:"data"`    // 对应 ZE API 的 data 字段
}

// ZEClient ZE题库客户端
type ZEClient struct {
	Enable bool
}

func (in *ZEClient) getHTTPClient() *resty.Client {
	return resty.New().SetTimeout(3 * time.Second)
}

// SearchAnswer 搜索答案
func (in *ZEClient) SearchAnswer(req model.SearchRequest) (answer [][]string, err error) {
	answer = make([][]string, 0)
	if !in.Enable {
		return answer, nil
	}

	// 1. 设置 ZE API 地址
	url := "http://localhost:3000/query"
	client := in.getHTTPClient()

	// 2. 构造请求体：添加 req.Options 字段
	resp, err := client.R().
		SetBody(map[string]interface{}{
			"title":   req.Question,
			"options": req.Options, // **已加入 options 字段**，直接使用请求模型中的 Options 数组
			"type":    util.GetTypeInt(req.Type),
		}).
		Post(url)

	if err != nil {
		return answer, errors.ErrRequest
	}

	var res zeResponse
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return answer, errors.ErrParserJSON
	}

	// 3. 解析响应
	if res.Code == 0 { // 假设 code: 0 为错误或未找到
		return answer, errors.New(res.Msg)
	}

	// 4. 提取和格式化答案
	if res.Data.Answer != "" {
		// ZE 题库返回的是文本答案，直接作为结果
		answer = append(answer, []string{res.Data.Answer})
	}

	return answer, nil
}
