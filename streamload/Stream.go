/*
 *@author  chengkenli
 *@project starrocks
 *@package streamload
 *@file    Stream
 *@date    2024/7/10 22:39
 */

package streamload

// StreamParms Stream Load主要结果体
type StreamParms struct {
	Host            string                   /*主机：FE IP*/
	Port            string                   /*端口：8030*/
	User            string                   /*账号*/
	Pass            string                   /*密码*/
	Schema          string                   /*库名*/
	Table           string                   /*表名*/
	Label           string                   /*唯一标识*/
	Expect          string                   /*100-continue*/
	File            string                   /*CSV - 文件绝对路径*/
	Format          string                   /*CSV - 格式化文件，csv或json*/
	ColumnSeparator string                   /*CSV - 列分隔符*/
	SkipHeader      string                   /*CSV - 跳过文件第几行*/
	TimeZone        string                   /*时区*/
	MaxFilterRatio  string                   /*JSON - 指定导入作业的最大容错率 取值范围：0~1*/
	StripOuterArray string                   /*JSON - 裁剪最外层的数组结构*/
	IgnoreJsonSize  string                   /*JSON - 是否检查 HTTP 请求中 JSON Body 的大小*/
	SchemaDataMap   []map[string]interface{} /*[]byte 数据流*/
}
