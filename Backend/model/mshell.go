package model

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
