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
	"starrocks/conn"
	"time"
)

var c conn.Connect

type StreamParms struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Schema          string
	Table           string
	Label           string
	Expect          string
	Format          string /*file*/
	ColumnSeparator string /*file*/
	SkipHeader      string /*file*/
	TimeZone        string
	MaxFilterRatio  string /*json*/
	StripOuterArray string /*json*/
	IgnoreJsonSize  string /*json*/
}

func init() {
	c = conn.Connect{
		Host: "xxxx",
		Port: 0,
		User: "xxxx",
		Pass: "xxxx",
		Base: "xxxx",
	}
}

func (s *StreamParms) StreamData() {
	db, err := c.StarRocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("数据库读取")
	var m []map[string]interface{}
	r := db.Raw("select '2024-07-10' as ts,'sr-app' as app,`queryId`,'','','',`timestamp`,`queryType`,`clientIp`,`user`,`authorizedUser`,`resourceGroup`,`catalog`,`db`,`state`,`errorCode`,`queryTime`,`scanBytes`,`scanRows`,`returnRows`,`cpuCostNs`,`memCostBytes`,`stmtId`,`isQuery`,`feIp`,`stmt`,`digest`,`planCpuCosts`,`planMemCosts` from srapp.audit.starrocks_audit_log where timestamp > (NOW() - INTERVAL 1 MINUTE) and queryTime > 10000").Scan(&m)
	if r.Error != nil {
		fmt.Println(r.Error.Error())
		return
	}
	marshal, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("提交Stream Load")
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
			"label":             time.Now().Format("2006_01_02_15_04_05"), /*label*/
			"Expect":            s.Expect,                                 /*在服务器拒绝导入作业请求的情况下，避免不必要的数据传输，减少不必要的资源开销。*/
			"format":            s.Format,                                 /*导入数据的格式。取值包括 CSV 和 JSON*/
			"timezone":          s.TimeZone,                               /*默认为东八区 (Asia/Shanghai)*/
			"max_filter_ratio":  s.MaxFilterRatio,                         /*指定导入作业的最大容错率 取值范围：0~1*/
			"strip_outer_array": s.StripOuterArray,                        /*裁剪最外层的数组结构*/
			"ignore_json_size":  s.IgnoreJsonSize,                         /*是否检查 HTTP 请求中 JSON Body 的大小*/
		}).SetBasicAuth(s.User, s.Pass).
		SetBody(marshal).
		Put(fmt.Sprintf("http://%s:%s/api/%s/%s/_stream_load", s.Host, s.Port, s.Schema, s.Table))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(response.Body()))

}
