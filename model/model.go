package model

import (
	"time"
)

type ReportV3 struct {
	Id        int64
	MarketId  int32
	Timestamp time.Time
	Code      string
	Name      string
	Last      float32
	Open      float32
	Close     float32
	High      float32
	Low       float32
	Amount    float32
	Volume    float32

	BuyPrice   []float32
	BuyVolume  []float32
	SellPrice  []float32
	SellVolume []float32
}

func (ReportV3) TableName() string {
	return "reportv3"
}

type Market struct {
	MarketId int32
	Date     int32
	Code     string
	Name     string
}

func (Market) TableName() string {
	return "market"
}

type Fenbi struct {
	Id        int64
	MarketId  int32
	Timestamp time.Time
	Code      string
	Open      float32
	Close     float32
	High      float32
	Low       float32
	Last      float32

	BuyPrice   []float32
	BuyVolume  []float32
	SellPrice  []float32
	SellVolume []float32
}

func (Fenbi) TableName() string {
	return "fenbi"
}

type Finance struct {
	Id      int64   `json:"id"`
	Market  int32   `json:"market"`  // 股票市场类型
	Code    string  `json:"code"`    // 股票代码,以'\0'结尾,如 "600050"  10个字节 同通视规范定义
	BGRQ    int64   `json:"ZGB"`     // 财务数据的日期 如半年报 季报等 如 20090630 表示 2009年半年报
	ZGB     float32 `json:"ZGB"`     // 总股本
	GJG     float32 `json:"GJG"`     // 国家股
	FQFRG   float32 `json:"FQFRG"`   // 发起人法人股
	FRG     float32 `json:"FRG"`     // 法人股
	BGS     float32 `json:"BGS"`     // B股
	HGS     float32 `json:"HGS"`     // H股
	MQLT    float32 `json:"MQLT"`    // 目前流通
	ZGG     float32 `json:"ZGG"`     // 职工股
	A2ZPG   float32 `json:"A2ZPG"`   // A2转配股
	ZZC     float32 `json:"ZZC"`     // 总资产(千元)
	LDZC    float32 `json:"LDZC"`    // 流动资产
	GDZC    float32 `json:"GDZC"`    // 固定资产
	WXZC    float32 `json:"WXZC"`    // 无形资产
	CQTZ    float32 `json:"CQTZ"`    // 长期投资
	LDFZ    float32 `json:"LDFZ"`    // 流动负债
	CQFZ    float32 `json:"CQFZ"`    // 长期负债
	ZBGJJ   float32 `json:"ZBGJJ"`   // 资本公积金
	MGGJJ   float32 `json:"MGGJJ"`   // 每股公积金
	GDQY    float32 `json:"GDQY"`    // 股东权益
	ZYSR    float32 `json:"ZYSR"`    // 主营收入
	ZYLR    float32 `json:"ZYLR"`    // 主营利润
	QTLR    float32 `json:"QTLR"`    // 其他利润
	YYLR    float32 `json:"YYLR"`    // 营业利润
	TZSY    float32 `json:"TZSY"`    // 投资收益
	BTSR    float32 `json:"BTSR"`    // 补贴收入
	YYWSZ   float32 `json:"YYWSZ"`   // 营业外收支
	SNSYTZ  float32 `json:"SNSYTZ"`  // 上年损益调整
	LRZE    float32 `json:"LRZE"`    // 利润总额
	SHLR    float32 `json:"SHLR"`    // 税后利润
	JLR     float32 `json:"JLR"`     // 净利润
	WFPLR   float32 `json:"WFPLR"`   // 未分配利润
	MGWFP   float32 `json:"MGWFP"`   // 每股未分配
	MGSY    float32 `json:"MGSY"`    // 每股收益
	MGJZC   float32 `json:"MGJZC"`   // 每股净资产
	TZMGJZC float32 `json:"TZMGJZC"` // 调整每股净资产
	GDQYB   float32 `json:"GDQYB"`   // 股东权益比
	JZCSYL  float32 `json:"JZCSYL"`  // 净资收益率
}

func (Finance) TableName() string {
	return "finance"
}

type FileHistory struct {
	Id        int64
	Code      string
	Timestamp time.Time
	Open      float32
	Close     float32
	High      float32
	Low       float32
	Volume    float32
	Amount    float32
	Advance   float32
	Decline   float32
}

func (FileHistory) TableName() string {
	return "file_history"
}

type FileMinute struct {
	Id        int64
	Code      string
	Timestamp time.Time
	Price     float32
	Volume    float32
	Amount    float32
}

func (FileMinute) TableName() string {
	return "file_minute"
}

type FileHISMinute struct {
	Id               int64
	Code             string
	Timestamp        time.Time
	Open             float32
	Close            float32
	High             float32
	Low              float32
	Amount           float32
	Volume           float32
	ActivedBuyVolume float32
}

func (FileHISMinute) TableName() string {
	return "file_hisminute"
}

type FilePower struct {
	Id          int64
	Code        string
	Timestamp   time.Time
	Give        float32
	Price  float32
	Volume float32
	Profit      float32
}

func (FilePower) TableName() string {
	return "file_power"
}

type FileBase struct {
	Attrib       int64
	Length       int64
	SerialNumber int64
	FileName     string
}

func (FileBase) TableName() string {
	return "file_base"
}

type FileNews struct {
	Attrib       int64
	Length       int64
	SerialNumber int64
	FileName     string
}

func (FileNews) TableName() string {
	return "file_news"
}
