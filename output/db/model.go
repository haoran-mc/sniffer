package db

import "time"

type Traffic struct {
	Ts      time.Time `gorm:"type:DateTime('Asia/Shanghai'); comment:时间(秒)"`
	Host    string    `gorm:"type:LowCardinality(String)"`
	Method  string    `gorm:"type:LowCardinality(String)"`
	Url     string    `gorm:"type:String"`
	Ip      string    `gorm:"type:String"`
	Status  string    `gorm:"type:LowCardinality(FixedString(3)); comment:响应状态码"`
	ReqBody string    `gorm:"type:String; comment:请求体"`
	ResBody string    `gorm:"type:String; comment:响应体"`
}
