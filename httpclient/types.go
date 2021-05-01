package httpclient

import "fmt"

// HealthForm form data
type HealthForm struct {
	Form Entity `json:"entity"`
}

// Entity detail form data
type Entity struct {
	// data auto filled by filler
	ReportStatus      string `fill:"bgdzsq" json:"bgdzsq"`               // 报告当地社区状态
	HealthStatus      string `fill:"jrstzk" json:"jrstzk"`               // 医学状态
	CloseContact3     string `fill:"sfjcysqzrq" json:"sfjcysqzrq"`       // 感染新冠人群接触
	CloseContact2     string `fill:"sfyqgzdyqryjc" json:"sfyqgzdyqryjc"` // 中高风险地区人员密切接触
	CloseContact4     string `fill:"sfjchwry" json:"sfjchwry"`           // 海归人员接触
	HometownCity      string `fill:"jgshi" json:"jgshi"`                 // 籍贯市
	CollegeName       string `fill:"sqbmmc" json:"sqbmmc"`               // 学院名称
	Area              string `fill:"jrszd" json:"jrszd"`                 // 当前所在地
	StudentNumber     string `fill:"gh" json:"gh"`                       // 学号
	PhoneNumber       string `fill:"lxdh" json:"lxdh"`                   // 联系电话
	StudentName       string `fill:"sqrmc" json:"sqrmc"`                 // 申请人姓名
	StudentID         string `fill:"sqrid" json:"sqrid"`                 // 申请人ID
	ClassID           string `fill:"sqbmid" json:"sqbmid"`               // 班级ID
	Gender            string `fill:"xb" json:"xb"`                       // 性别
	Fever             string `fill:"sffr" json:"sffr"`                   // 发热
	TimeSeries        string `fill:"glqsrq" json:"glqsrq"`               // 时间区间
	CloseContact1     string `fill:"sfyyqryjc" json:"sfyyqryjc"`         // 中高风险地区途径人员密切接触
	IdentityNumber    string `fill:"sfzh" json:"sfzh"`                   // 身份证号码
	HighRiskArea      string `fill:"jgzgfxdq" json:"jgzgfxdq"`           // 经过的中高风险地区
	CloseContact5     string `fill:"jrsfjgzgfxdq" json:"jrsfjgzgfxdq"`   // 途径高风险地区
	Isolation         string `fill:"jzzt" json:"jzzt"`                   // 隔离状态
	HometownProvince  string `fill:"jgshen" json:"jgshen"`               // 籍贯省
	DetailArea        string `fill:"jrjzdxxdz" json:"jrjzdxxdz"`         // 详细居住地址
	LeaveArea         string `fill:"sflz" json:"sflz"`                   // 离开镇江/张家港
	Date              string `fill:"tbrq" json:"tbrq"`                   // 填表日期
	ReturnTransportID string `fill:"fhzjbc" json:"fhzjbc"`               // 返镇交通工具班次
	ReturnTransport   string `fill:"fhzjgj" json:"fhzjgj"`               // 返镇交通工具
	LeaveTransportID  string `fill:"lzbc" json:"lzbc"`                   // 离镇交通工具班次
	LeaveTransport    string `fill:"lzjtgj" json:"lzjtgj"`               // 离镇交通工具
	LeaveTime         string `fill:"lzsj" json:"lzsj"`                   // 离镇时间
	Rysf              string `fill:"rysf" json:"rysf"`                   // unknown

	// auto filled, todo: judge the status?
	ReturnArea        string `fill:"sffz" json:"sffz"`         // 返回镇江/张家港
	ReturnTime        string `fill:"fhzjsj" json:"fhzjsj"`     // 返镇时间
	ReturnThroughArea string `fill:"fztztkdd" json:"fztztkdd"` // 返镇途径地点

	// data can't be auto filled
	Time                 string `json:"tjsj"` // 填报时间 2021-05-01 10:32
	MorningTemperature   string `json:"tw"`   // 今晨体温
	LastNightTemperature string `json:"zwtw"` // 昨晚体温
	Remark               string `json:"bz"`   // 备注, default to ""
	Ext                  string `json:"_ext"` // not use， default to "{}"
}

type CookieNotFoundErr struct {
	cookie string
}

// ElementNotFoundErr error interface for element
type ElementNotFoundErr struct {
	element string
}

func (t CookieNotFoundErr) Error() string {
	return fmt.Sprintf("http: can't find cookie: %s", t.cookie)
}

func (t ElementNotFoundErr) Error() string {
	return fmt.Sprintf("http: can't found element: %s", t.element)
}
