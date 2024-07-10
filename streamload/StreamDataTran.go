/*
 *@author  chengkenli
 *@project 2023
 *@package streamload
 *@file    StreamDatastream
 *@date    2024/7/10 14:59
 */

package streamload

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

// StreamDataTransInit 初始化Resty客户端
func (s *StreamParms) StreamDataTransInit() (*resty.Client, error) {
	//创建Resty客户端
	Client := resty.New()
	//进行重定向
	Client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		// 这里你可以根据需要添加自定义逻辑，比如保留headers等
		for key, values := range via[0].Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		// 如果想要完全信任所有重定向，只需返回nil
		return nil
	})).SetHeaders(map[string]string{
		"label":             s.Label,           /*label*/
		"format":            s.Format,          /*导入数据的格式。取值包括 CSV 和 JSON*/
		"timezone":          s.TimeZone,        /*默认为东八区 (Asia/Shanghai)*/
		"max_filter_ratio":  s.MaxFilterRatio,  /*指定导入作业的最大容错率 取值范围：0~1*/
		"strip_outer_array": s.StripOuterArray, /*裁剪最外层的数组结构*/
		"ignore_json_size":  s.IgnoreJsonSize,  /*是否检查 HTTP 请求中 JSON Body 的大小*/
		"db":                s.Schema,          /*库名*/
		"table":             s.Table,           /*表名*/
	}).SetBasicAuth(s.User, s.Pass)
	return Client, nil
}

// StreamDataTransToBegin 发送POST请求并处理响应（覆盖处理重定向-开始事务）
func (s *StreamParms) StreamDataTransToBegin(Client *resty.Client) ([]byte, error) {
	begin, err := Client.R().Post(fmt.Sprintf("http://%s:%s/api/transaction/begin", s.Host, s.Port))
	if err != nil {
		return begin.Body(), err
	}
	return begin.Body(), nil
}

// StreamDataTransToLoad 发送POST请求并处理响应（写入数据）
func (s *StreamParms) StreamDataTransToLoad(Client *resty.Client) ([]byte, error) {
	marshal, err := json.Marshal(s.SchemaDataMap)
	if err != nil {
		return marshal, err
	}
	load, err := Client.R().SetBody(marshal).Put(fmt.Sprintf("http://%s:%s/api/transaction/load", s.Host, s.Port))
	if err != nil {
		return load.Body(), err
	}
	return load.Body(), nil
}

// StreamDataTransToCommit 发送POST请求并处理响应（提交事务）
func (s *StreamParms) StreamDataTransToCommit(Client *resty.Client) ([]byte, error) {
	commit, err := Client.R().Post(fmt.Sprintf("http://%s:%s/api/transaction/commit", s.Host, s.Port))
	if err != nil {
		return commit.Body(), err
	}
	return commit.Body(), nil
}

// StreamDataTransToRollback 发送POST请求并处理响应（事务回滚） 只有当作业异常时，事务回滚才会生效，当事务提交正常，事务回滚无法使用。
func (s *StreamParms) StreamDataTransToRollback(Client *resty.Client) ([]byte, error) {
	rollback, err := Client.R().Post(fmt.Sprintf("http://%s:%s/api/transaction/rollback", s.Host, s.Port))
	if err != nil {
		return rollback.Body(), err
	}
	return rollback.Body(), nil
}
