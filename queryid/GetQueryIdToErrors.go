package queryid

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type QueryError struct {
	ManagerWeb  string
	ManagerPort int
	User        string
	Pass        string
	QueryId     string
}

func (q *QueryError) GetQueryIdToError() (string, error) {
	/*-------------------------登录manager web-------------------------------*/
	//创建Resty客户端
	Client := resty.New()
	//发送POST请求并处理响应
	response, err := Client.R().SetBody(map[string]string{
		"name":     q.User,
		"password": q.Pass,
	}).Post(fmt.Sprintf("http://%s:%d/api/user/login", q.ManagerWeb, q.ManagerPort))
	if err != nil {
		return string(response.Body()), err
	}
	/*-------------------------登录manager web end----------------------------*/
	/*请求license*/
	response, err = Client.R().Get(fmt.Sprintf("http://%s:%d/api/query/detail/%s", q.ManagerWeb, q.ManagerPort, q.QueryId))
	if err != nil {
		return string(response.Body()), err
	}
	type mm struct {
		Code int `json:"code"`
		Data struct {
			QueryDetail struct {
				QueryID      string   `json:"queryId"`
				User         string   `json:"user"`
				Status       string   `json:"status"`
				ErrorMessage string   `json:"errorMessage"`
				StartTime    int      `json:"startTime"`
				EndTime      int      `json:"endTime"`
				TimeUsed     int      `json:"timeUsed"`
				SQL          string   `json:"sql"`
				ExplainText  []string `json:"explainText"`
				ProfileText  []string `json:"profileText"`
			} `json:"queryDetail"`
		} `json:"data"`
	}
	var m mm
	err = json.Unmarshal(response.Body(), &m)
	if err != nil {
		return string(response.Body()), err
	}
	if m.Data.QueryDetail.ErrorMessage == "" {
		return string(response.Body()), errors.New("error info result is nil")
	}
	return m.Data.QueryDetail.ErrorMessage, nil
}
