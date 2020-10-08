package packect

type ReportV3Packet struct {
	Count int32          `json:"count"`
	Items []ReportV3Item `json:"items"`
}

type ReportV3Item struct {
	PackectSize int32     `json:"m_cbSize"`
	Timestamp   int64     `json:"m_time"`
	MarketId    int32     `json:"m_wMarket"`
	Label       []string  `json:"m_szLabel"`
	Name        []string  `json:"m_szName"`
	Open        float32   `json:"m_fOpen"`
	Close       float32   `json:"m_fLastClose"`
	High        float32   `json:"m_fHigh"`
	Low         float32   `json:"m_fLow"`
	Last        float32   `json:"m_fNewPrice"`
	Volume      float32   `json:"m_fVolume"`
	Amount      float32   `json:"m_fAmount"`
	BuyPrice    []float32 `json:"m_fBuyPrice"`
	BuyVolume   []float32 `json:"m_fBuyVolume"`
	SellPrice   []float32 `json:"m_fSellPrice"`
	SellVolume  []float32 `json:"m_fSellVolume"`
	BuyPrice4   float32   `json:"m_fBuyPrice4"`
	BuyVolume4  float32   `json:"m_fBuyVolume4"`
	SellPrice4  float32   `json:"m_fSellPrice4"`
	SellVolume4 float32   `json:"m_fSellVolume4"`
	BuyPrice5   float32   `json:"m_fBuyPrice5"`
	BuyVolume5  float32   `json:"m_fBuyVolume5"`
	SellPrice5  float32   `json:"m_fSellPrice5"`
	SellVolume5 float32   `json:"m_fSellVolume5"`
}

type FenbiPacket struct {
	Count  int32       `json:"count"`
	Header FenbiHeader `json:"header"`
	Items  []FenbiItem `json:"items"`
}

type FenbiHeader struct {
	Market int32    `json:"m_wMarket"`
	Label  []string `json:"m_szLabel"`
	Date   int64    `json:"m_lDate"`
	Open   float32  `json:"m_fOpen"`
	Close  float32  `json:"m_fLastClose"`
	Count  int32    `json:"m_nCount"`
}

type FenbiItem struct {
	Timestamp  int64     `json:"m_lTime"`
	Open       float32   `json:"m_fOpen"`
	Close      float32   `json:"m_fLastClose"`
	High       float32   `json:"m_fHigh"`
	Low        float32   `json:"m_fLow"`
	Last       float32   `json:"m_fNewPrice"`
	Volume     float32   `json:"m_fVolume"`
	Amount     float32   `json:"m_fAmount"`
	Stroke     int64     `json:"m_lStroke"`
	BuyPrice   []float32 `json:"m_fBuyPrice"`
	BuyVolume  []float32 `json:"m_fBuyVolume"`
	SellPrice  []float32 `json:"m_fSellPrice"`
	SellVolume []float32 `json:"m_fSellVolume"`
}

type MarketPacket struct {
	Count  int32        `json:"count"`
	Header MarketHeader `json:"header"`
	Items  []MarketItem `json:"items"`
}

type MarketHeader struct {
	Market      int32    `json:"m_wMarket"`
	Name        []string `json:"m_Name"`
	Property    int32    `json:"m_lProperty"`
	Date        int32    `json:"m_lDate"`
	PeriodCount int32    `json:"m_PeriodCount"`
	OpenTime    []int32  `json:"m_OpenTime"`
	CloseTime   []int32  `json:"m_CloseTime"`
	Count       int32    `json:"m_nCount"`
}

type MarketItem struct {
	Label    []string `json:"m_szLabel"`
	Name     []string `json:"m_szName"`
	Property int32    `json:"m_cProperty"`
}

type FinancePacket struct {
	Count int32         `json:"count"`
	Items []FinanceItem `json:"items"`
}

type FinanceItem struct {
	Market  int32    `json:"m_wMarket"` // 股票市场类型
	N1      int32    `json:"N1"`        // 保留字段
	Label   []string `json:"m_szLabel"` // 股票代码,以'\0'结尾,如 "600050"  10个字节 同通视规范定义
	BGRQ    int64    `json:"ZGB"`       // 财务数据的日期 如半年报 季报等 如 20090630 表示 2009年半年报
	ZGB     float32  `json:"ZGB"`       // 总股本
	GJG     float32  `json:"GJG"`       // 国家股
	FQFRG   float32  `json:"FQFRG"`     // 发起人法人股
	FRG     float32  `json:"FRG"`       // 法人股
	BGS     float32  `json:"BGS"`       // B股
	HGS     float32  `json:"HGS"`       // H股
	MQLT    float32  `json:"MQLT"`      // 目前流通
	ZGG     float32  `json:"ZGG"`       // 职工股
	A2ZPG   float32  `json:"A2ZPG"`     // A2转配股
	ZZC     float32  `json:"ZZC"`       // 总资产(千元)
	LDZC    float32  `json:"LDZC"`      // 流动资产
	GDZC    float32  `json:"GDZC"`      // 固定资产
	WXZC    float32  `json:"WXZC"`      // 无形资产
	CQTZ    float32  `json:"CQTZ"`      // 长期投资
	LDFZ    float32  `json:"LDFZ"`      // 流动负债
	CQFZ    float32  `json:"CQFZ"`      // 长期负债
	ZBGJJ   float32  `json:"ZBGJJ"`     // 资本公积金
	MGGJJ   float32  `json:"MGGJJ"`     // 每股公积金
	GDQY    float32  `json:"GDQY"`      // 股东权益
	ZYSR    float32  `json:"ZYSR"`      // 主营收入
	ZYLR    float32  `json:"ZYLR"`      // 主营利润
	QTLR    float32  `json:"QTLR"`      // 其他利润
	YYLR    float32  `json:"YYLR"`      // 营业利润
	TZSY    float32  `json:"TZSY"`      // 投资收益
	BTSR    float32  `json:"BTSR"`      // 补贴收入
	YYWSZ   float32  `json:"YYWSZ"`     // 营业外收支
	SNSYTZ  float32  `json:"SNSYTZ"`    // 上年损益调整
	LRZE    float32  `json:"LRZE"`      // 利润总额
	SHLR    float32  `json:"SHLR"`      // 税后利润
	JLR     float32  `json:"JLR"`       // 净利润
	WFPLR   float32  `json:"WFPLR"`     // 未分配利润
	MGWFP   float32  `json:"MGWFP"`     // 每股未分配
	MGSY    float32  `json:"MGSY"`      // 每股收益
	MGJZC   float32  `json:"MGJZC"`     // 每股净资产
	TZMGJZC float32  `json:"TZMGJZC"`   // 调整每股净资产
	GDQYB   float32  `json:"GDQYB"`     // 股东权益比
	JZCSYL  float32  `json:"JZCSYL"`    // 净资收益率
}

type FileHeaderEKE struct {
	Tag    int64    `json:"m_dwHeadTag"`
	Market int32    `json:"m_wMarket"`
	Label  []string `json:"m_szLabel"`
}

type FileHistoryPacket struct {
	Count   int32             `json:"count"`
	Headers []FileHeaderEKE   `json:"headers"`
	Items   []FileHistoryItem `json:"items"`
}

type FileHistoryItem struct {
	Timestamp int64   `json:"m_time"`
	Open      float32 `json:"m_fOpen"`
	Close     float32 `json:"m_fClose"`
	High      float32 `json:"m_fHigh"`
	Low       float32 `json:"m_fLow"`
	Volume    float32 `json:"m_fVolume"`
	Amount    float32 `json:"m_fAmount"`
	Advance   float32 `json:"m_wAdvance"`
	Decline   float32 `json:"m_wDecline"`
}

type FileMinutePacket struct {
	Count   int32            `json:"count"`
	Headers []FileHeaderEKE  `json:"headers"`
	Items   []FileMinuteItem `json:"items"`
}

type FileMinuteItem struct {
	Timestamp int64   `json:"m_time"`
	Price     float32 `json:"m_fPrice"`
	Volume    float32 `json:"m_fVolume"`
	Amount    float32 `json:"m_fAmount"`
}

type FileHISMinutePackect struct {
	Count   int32               `json:"count"`
	Headers []FileHeaderEKE     `json:"headers"`
	Items   []FileHISMinuteItem `json:"items"`
}

type FileHISMinuteItem struct {
	Timestamp  int64   `json:"m_time"`
	Open       float32 `json:"m_fOpen"`
	Close      float32 `json:"m_fClose"`
	High       float32 `json:"m_fHigh"`
	Low        float32 `json:"m_fLow"`
	Volume     float32 `json:"m_fVolume"`
	Amount     float32 `json:"m_fAmount"`
	ActivedBuy float32 `json:"m_fActiveBuyVol"`
}

type FilePowerPacket struct {
	Count   int32           `json:"count"`
	Headers []FileHeaderEKE `json:"headers"`
	Items   []FilePowerItem `json:"items"`
}

type FilePowerItem struct {
	Timestamp int64   `json:"m_time"`
	Give      float32 `json:"m_fGive"`
	Pei       float32 `json:"m_fPei"`
	PeiPrice  float32 `json:"m_fPeiPrice"`
	Profit    float32 `json:"m_fProfit"`
}

type FileHeaderEx struct {
	Attrib   int64    `json:"m_dwAttrib"`
	Len      int64    `json:"m_dwLen"`
	SerialNo int64    `json:"m_dwSerialNo"`
	FileName []string `json:"m_szFileName"`
}
