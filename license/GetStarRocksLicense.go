package license

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type ParmsLicense struct {
	ManagerHost string
	ManagerPort string
	UserName    string
	PassWord    string
}
type InfoLicense struct {
	Cores        int
	Hosts        int
	ExpireAtUnix int64
	ExpireAtTime string
	License      string
}

func GetStarRocksLicense(p ParmsLicense) (*InfoLicense, error) {
	type license struct {
		Code int `json:"code"`
		List []struct {
			Cores    int   `json:"cores"`
			ExpireAt int64 `json:"expire_at"`
			Hosts    int   `json:"hosts"`
		} `json:"list"`
		Total int `json:"total"`
	}
	type key struct {
		Code int `json:"code"`
		Data struct {
			Cores int    `json:"cores"`
			D     string `json:"d"`
		} `json:"data"`
	}

	//创建Resty客户端
	Client := resty.New()
	//发送POST请求并处理响应
	respones1, err := Client.R().SetBody(map[string]string{"name": p.UserName, "password": p.PassWord}).Post(fmt.Sprintf("%s:%s/api/user/login", p.ManagerHost, p.ManagerPort))
	if err != nil {
		return &InfoLicense{}, err
	}
	if err != nil {
		return &InfoLicense{}, err
	}
	fmt.Println(string(respones1.Body()))
	/*------------------------------------------------*/
	respones2, err := Client.R().Get(fmt.Sprintf("%s:%s/api/license/list", p.ManagerHost, p.ManagerPort))
	if err != nil {
		return &InfoLicense{}, err
	}
	fmt.Println(string(respones2.Body()))
	/*------------------------------------------------*/
	respones3, err := Client.R().Get(fmt.Sprintf("%s:%s/api/license/collect-hosts-info", p.ManagerHost, p.ManagerPort))
	if err != nil {
		return &InfoLicense{}, err
	}
	fmt.Println(string(respones3.Body()))
	/*------------------------------------------------*/
	var k key
	err = json.Unmarshal(respones3.Body(), &k)
	if err != nil {
		return &InfoLicense{}, err
	}

	var l license
	err = json.Unmarshal(respones2.Body(), &l)
	if err != nil {
		return &InfoLicense{}, err
	}

	var info *InfoLicense
	for _, s := range l.List {
		info = &InfoLicense{
			Cores:        s.Cores,
			Hosts:        s.Hosts,
			ExpireAtUnix: s.ExpireAt,
			ExpireAtTime: unixToTime(strconv.FormatInt(s.ExpireAt, 10)).Format("2006-01-02 15:04:05"),
			License:      k.Data.D,
		}
	}
	return info, nil
}

func unixToTime(e string) (datatime time.Time) {
	data, _ := strconv.ParseInt(e, 10, 64)
	datatime = time.Unix(data/1000, 0)
	return
}
