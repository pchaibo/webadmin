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
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Mmurl     string `json:"mmurl"`
	Mmtext    string `json:"mmtext"`
	Checkurl  string `json:"checkurl"`
	Checktext string `json:"checktext"`
	Status    int    `json:"status"` //状态  1:正常 2:停用
	Addtime   int    `json:"addtime"`
}
