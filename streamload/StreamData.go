/*
 *@author  chengkenli
 *@project 2023
 *@package streamload
 *@file    schemastream
 *@date    2024/7/10 13:32
 */

package streamload

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func (s *StreamParms) StreamDataToByte() ([]byte, error) {
	marshal, err := json.Marshal(s.SchemaDataMap)
	if err != nil {
		return marshal, err
	}

	//创建Resty客户端
	Client := resty.New()
	//发送POST请求并处理响应
	response, err := Client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		// 这里你可以根据需要添加自定义逻辑，比如保留headers等
		for key, values := range via[0].Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		// 如果想要完全信任所有重定向，只需返回nil
		return nil
	})).R().
		SetHeaders(map[string]string{
			"label":             s.Label,           /*label*/
			"Expect":            s.Expect,          /*在服务器拒绝导入作业请求的情况下，避免不必要的数据传输，减少不必要的资源开销。*/
			"format":            s.Format,          /*导入数据的格式。取值包括 CSV 和 JSON*/
			"timezone":          s.TimeZone,        /*默认为东八区 (Asia/Shanghai)*/
			"max_filter_ratio":  s.MaxFilterRatio,  /*指定导入作业的最大容错率 取值范围：0~1*/
			"strip_outer_array": s.StripOuterArray, /*裁剪最外层的数组结构*/
			"ignore_json_size":  s.IgnoreJsonSize,  /*是否检查 HTTP 请求中 JSON Body 的大小*/
		}).SetBasicAuth(s.User, s.Pass).
		SetBody(marshal).
		Put(fmt.Sprintf("http://%s:%s/api/%s/%s/_stream_load", s.Host, s.Port, s.Schema, s.Table))
	if err != nil {
		return response.Body(), err
	}
	return response.Body(), nil
}
