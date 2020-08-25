package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

// 场景描述：我有很多行日志数据单行的格式是这样的
// HOST;000012000629948340196501;ipv4;3; ips: user_id=2;user_name=172.21.1.102;policy_id=1;src_mac=52:54:00:62:7f:4a;dst_mac=58:69:6c:7b:fa:e7;src_ip=172.21.1.102;dst_ip=172.22.2.3;src_port=48612;dst_port=80;app_name=网页浏览(HTTP);protocol=TCP;app_protocol=HTTP;event_id=1310909;event_name=Microsoft_IIS_5.1_Frontpage扩展路径信息漏洞;event_type=安全漏洞;level=info;ctime=2019-12-26 11:17:17;action=pass
// 其中ips: 之前的都是不规范的字段
// 我需要把他解析成结构化的数据,这样的
type IpsItem struct {
	UserId      int    `json:"user_id"`
	UserName    string `json:"user_name"`
	SrcIp       string `json:"src_ip"`
	DstIp       string `json:"dst_ip"`
	SrcPort     int    `json:"src_port"`
	DstPort     int    `json:"dst_port"`
	AppName     string `json:"app_name"`
	Protocol    string `json:"protocol"`
	AppProtocol string `json:"app_protocol"`
	EventId     int    `json:"event_id"`
	EventName   string `json:"event_name"`
	EventType   string `json:"event_type"`
	Level       string `json:"level"`
	Ctime       string `json:"ctime"`
	Action      string `json:"action"`
}

var testRawString = "HOST;000012000629948340196501;ipv4;3; ips: user_id=2;user_name=172.21.1.102;policy_id=1;src_mac=52:54:00:62:7f:4a;dst_mac=58:69:6c:7b:fa:e7;src_ip=172.21.1.102;dst_ip=172.22.2.3;src_port=48612;dst_port=80;app_name=网页浏览(HTTP);protocol=TCP;app_protocol=HTTP;event_id=1311495;event_name=HTTP_Nikto_WEB漏洞扫描;event_type=安全扫描;level=warning;ctime=2019-12-26 11:17:17;action=pass"

// 解析原始日志结构
func NewIpsItem(raw string) *IpsItem {
	//清除非法的字符
	raw = strings.ReplaceAll(raw, ":", ";")

	ins := IpsItem{}
	t := reflect.TypeOf(ins)
	//遍历结构体属性
	for i := 0; i < t.NumField(); i++ {
		//获取属性structField
		sf := t.Field(i)
		//属性名称
		fieldName := sf.Name
		//tag json的值
		tagName := sf.Tag.Get("json")

		//获取字段值
		// FieldByName returns the struct field with the given name.
		// It returns the zero Value if no field was found.
		// It panics if v's Kind is not struct.
		fieldValue := reflect.ValueOf(&ins).Elem().FieldByName(fieldName)

		//属性的值 type
		switch sf.Type.Name() {
		case "int":
			var someInt int64
			scanValueFromString(raw, tagName, tagName+"=%d", &someInt)
			//给属性赋值
			fieldValue.SetInt(someInt)
			//todo:: 支持更多类型
		default:
			var someString string
			scanValueFromString(raw, tagName, tagName+"=%s", &someString)
			////给属性赋值
			fieldValue.SetString(someString)
		}

	}
	return &ins
}

//scanValueFromString 字符串 字段的值
func scanValueFromString(raw string, tagJsonValue, format string, someV interface{}) {
	for _, ss := range strings.Split(raw, ";") {
		ele := strings.TrimSpace(ss)
		if strings.HasPrefix(ele, tagJsonValue) {
			fmt.Sscanf(ele, format, someV)
			//n, err := fmt.Sscanf(ele, format, someV)
			//fmt.Println(n, err)
			return
		}
	}
}

func TestParse() {
	ii := NewIpsItem(testRawString)
	fmt.Printf("%+v\n", ii)
}
