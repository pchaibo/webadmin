package model

type Admin struct {
	Id           int    `gorm:"primaryKey" json:"id" `
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Status       int    `json:"status"`
	Registertime string `json:"register_time"`
}

type AuthGroup struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Rules  string `json:"rules"`
}

type AuthRule struct {
	Id        int    `json:"id"`
	Pid       int    `json:"pid"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	Type      int    `json:"type"`
	Status    int    `json:"status"`
	Condition string `json:"condition"`
}

type ShellGroup struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Mmurl   string `json:"mmurl" gorm:"column:url"`
	Addtime int    `json:"addtime"`
}

type Shell struct {
	Id      int         `json:"id"`
	Host    string      `json:"host"`
	Scheme  string      `json:"scheme"`
	GroupId int         `json:"group_id"`
	Status  int         `json:"status"`
	Num     int         `json:"num"`
	Sitenum int         `json:"sitenum"`
	Maxurl  string      `json:"maxurl"`
	Minurl  string      `json:"minurl"`
	Dir     int         `json:"dir"`
	Lock    int         `json:"lock"`
	Remark  string      `json:"remark"`
	Addtime int         `json:"addtime"`
	Uptime  int         `json:"uptime"`
	Group   *ShellGroup `gorm:"foreignKey:GroupId" json:"group,omitempty"`
}

type ShellMax struct {
	Id      int    `json:"id"`
	ShellId int    `json:"shell_id"`
	Url     string `json:"url"`
	Addtime int    `json:"addtime"`
	Status  int    `json:"status"`
}

// func (ShellMax) TableName() string {// 	return "tu_shellmax"// }

type ShellMin struct {
	Id      int    `json:"id"`
	ShellId int    `json:"shell_id"`
	Url     string `json:"url"`
	Addtime int    `json:"addtime"`
	Status  int    `json:"status"`
}

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	Bnaccess  string `json:"bnaccess"`
	Bnasecret string `json:"bnasecret"`
	Status    int    `json:"status"`
}

// heyue:合约
type Coin struct {
	Id             int64   `gorm:"primaryKey" json:"id" `
	Name           string  `json:"name"`           //名称
	Symbol         string  `json:"symbol"`         //合约代码
	Close          float64 `json:"close"`          //最新价
	Priceprecision int     `json:"priceprecision"` //价格精度
	Open           float64 `json:"open"`
	Low            float64 `json:"low"`
	High           float64 `json:"high"`
	Status         uint    `json:"status"` //状态  1:正常 0:停用
	AddTime        int64   `json:"addtime"`
	UpdateTime     int64   `json:"updatetime"`
}

func (Coin) TableName() string {
	return "tu_heyue_coin"
}

type Heyue struct {
	Id              uint    `json:"id"`
	UserId          uint    `json:"userid"`
	UserName        string  `json:"username"`
	Symbol          string  `json:"symbol"`
	Side            int32   `json:"side"`
	Num             int32   `json:"num"`
	Is_num          int32   `json:"is_num"`
	Status          int32   `json:"status"`
	Sellprice       float64 `json:"sellprice"`                        //收益百分比
	Oneprice        float64 `from:"oneprice" json:"oneprice"`         //首仓usdt
	Repeatprice     float64 `from:"Repeatprice" json:"repeatprice"`   //补仓usdt
	Rangetype       int     `json:"rangetype"`                        //网格类型 1:差价usdt 2:网格保证金
	Rangeprice      float64 `from:"rangeprice" json:"rangeprice"`     //网格差价usdt
	Rangepercent    int     `json:"rangepercent"`                     //网格保证金百分比
	Rangeclosingpct int     `json:"rangeclosingpct"`                  //网格平仓百分比
	Rangeclosing    int     `json:"rangeclosing"`                     //网格平仓 1:不平 2：平仓
	Closingprice    float64 `from:"closingprice" json:"closingprice"` //强平价格
	Risk            int     `json:"risk"`                             //风控
	RiskTime        int     `json:"risktime"`                         //风控时间
	Newprice        float64 `from:"newprice" json:"newprice"`         //最新价格
	NewTime         int64   `json:"newtime"`                          //价格更新时间
	AddTime         int64   `json:"addtime"`
	UpdateTime      int64   `json:"updatetime"`
}

type Heyuesorder struct {
	Id         uint    `json:"id"`
	Ordertype  int32   `json:"ordertype"` //交易类型:1=开仓,2=平仓
	UserId     uint    `json:"userid"`
	Username   string  `json:"username"`
	Symbol     string  `json:"symbol"`
	Side       int32   `json:"side"` //状态:1=开多,2=开空
	Price      float64 `json:"price"`
	Total      float64 `json:"total"`    //总金额
	Quantity   float64 `json:"quantity"` //数量
	Num        int32   `json:"num"`      //第几次
	Orderid    int64   `json:"orderid"`
	Log        string  `json:"log"`
	Status     int32   `json:"status"`
	Usdt       float64 `json:"usdt"` //收效
	AddTime    int64   `json:"addtime"`
	UpdateTime int64   `json:"updatetime"`
}
