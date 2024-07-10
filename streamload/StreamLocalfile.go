/*
 *@author  chengkenli
 *@project 2023
 *@package streamload
 *@file    localfile
 *@date    2024/7/10 13:29
 */

package streamload

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
)

func (s *StreamParms) StreamLocalfile() {
	path := "C:\\Users\\chengken\\Desktop\\2024\\demodata"
	readFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("文件读出成功!")

	fmt.Println("提交Stream Load!")
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
			"label":            s.Label,           /*label*/
			"Expect":           s.Expect,          /*在服务器拒绝导入作业请求的情况下，避免不必要的数据传输，减少不必要的资源开销。*/
			"format":           s.Format,          /*导入数据的格式。取值包括 CSV 和 JSON*/
			"column_separator": s.ColumnSeparator, /*列分隔符*/
			"skip_header":      s.SkipHeader,      /*指定跳过 CSV 文件最开头的几行数据*/
			"timezone":         s.TimeZone,        /*默认为东八区 (Asia/Shanghai)*/
		}).SetBasicAuth(s.User, s.Pass).
		SetBody(readFile).
		Put(fmt.Sprintf("http://%s:%s/api/%s/%s/_stream_load", s.Host, s.Port, s.Schema, s.Table))
	if err != nil {
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(response.Body()))

}
