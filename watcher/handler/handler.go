package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	pb "github.com/emuta/haililive/protobuf/haililive"
	"github.com/emuta/haililive/watcher/packect"
)

func TimestampToString(unixTimestamp int64, layout string) string {
	t := time.Unix(unixTimestamp, 0)
	return t.Format(layout)
}

func TimeStringLong(unixTimestamp int64) string {
	return TimestampToString(unixTimestamp, "2006-01-02 15:04:05")
}

func TrimSpaceString(strSlice []string) string {
	s := strings.Join(strSlice, "")
	var r = strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "", "\x00", "")
	// strings.TrimSpace(s)
	return r.Replace(s)
}

func TimestampProto(unixTimestamp int64) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{Seconds: unixTimestamp, Nanos: 0}
}

type Handler struct {
	cli pb.HaiLiLiveServiceClient
}

func NewHandler(client pb.HaiLiLiveServiceClient) *Handler {
	return &Handler{cli: client}
}

func (h *Handler) OnMessage(ctx context.Context, delivery amqp.Delivery) error {

	switch delivery.Type {
	// just a test, pls ignore this case.
	case "relayer.connection.test":
		go h.ConnectionTest(ctx, &delivery)

	case "PACKECT_REPORT_STRUCTExV3":
		go h.ReportV3Packet(ctx, &delivery)

	case "PACKECT_FENBIDATA_STRUCT":
		go h.FenbiPacket(ctx, &delivery)

	case "PACKECT_MARKET_STRUCT":
		go h.MarketPacket(ctx, &delivery)

	case "PACKECT_FIN_LJF_STRUCTEx":
		go h.FinancePacket(ctx, &delivery)

	case "PACKECT_HISTORY_STRUCTEx":
		go h.FileHistoryPacket(ctx, &delivery)

	case "PACKECT_MINUTE_STRUCTEx":
		go h.FileMinutePacket(ctx, &delivery)

	case "PACKECT_HISMINUTE_STRUCTEx":
		go h.FileHISMinutePackect(ctx, &delivery)

	case "PACKECT_POWER_STRUCTEx":
		go h.FilePowerPacket(ctx, &delivery)

	default:
		// do nothing
		go h.UndefinedTypeMessage(ctx, &delivery)
	}

	return nil
}

func (h *Handler) UndefinedTypeMessage(ctx context.Context, msg *amqp.Delivery) {
	log.WithFields(log.Fields{
		"typ": "undefined",
	}).Infof("%s", msg.Body)
}

func (h *Handler) ConnectionTest(ctx context.Context, msg *amqp.Delivery) {
	log.WithFields(log.Fields{
		"typ":       "test",
		"Timestamp": msg.Timestamp,
		"AppId":     msg.AppId,
		"Type":      msg.Type,
		"MessageId": msg.MessageId,
	}).Debug("%s", msg.Body)
}

func (h *Handler) ReportV3Packet(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.ReportV3Packet
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to PACKECT_REPORT_STRUCTExV3")
		return
	}

	for i, item := range pkt.Items {
		log.WithFields(log.Fields{
			"i":      i,
			"ts":     TimeStringLong(item.Timestamp),
			"label":  TrimSpaceString(item.Label),
			"name":   TrimSpaceString(item.Name),
			"price":  item.Last,
			"open":   item.Open,
			"close":  item.Close,
			"high":   item.High,
			"low":    item.Low,
			"amount": item.Amount,
		}).Debug("reportV3")

		go func() {
			req := pb.CreateReportV3Request{
				Timestamp:  TimestampProto(item.Timestamp),
				MarketId:   item.MarketId,
				Code:       TrimSpaceString(item.Label),
				Name:       TrimSpaceString(item.Name),
				Open:       item.Open,
				Close:      item.Close,
				High:       item.High,
				Low:        item.Low,
				Last:       item.Last,
				Volume:     item.Volume,
				Amount:     item.Volume,
				BuyPrice:   item.BuyPrice,
				BuyVolume:  item.BuyVolume,
				SellPrice:  item.SellPrice,
				SellVolume: item.SellVolume,
			}

			req.BuyPrice = append(req.BuyPrice, item.BuyPrice4)
			req.BuyVolume = append(req.BuyVolume, item.BuyVolume4)
			req.SellPrice = append(req.SellPrice, item.SellPrice4)
			req.SellVolume = append(req.SellVolume, item.SellVolume4)

			req.BuyPrice = append(req.BuyPrice, item.BuyPrice5)
			req.BuyVolume = append(req.BuyVolume, item.BuyVolume5)
			req.SellPrice = append(req.SellPrice, item.SellPrice5)
			req.SellVolume = append(req.SellVolume, item.SellVolume5)

			_, _ = h.cli.CreateReportV3(ctx, &req)
		}()

	}

	log.WithFields(log.Fields{
		"count": pkt.Count,
	}).Info("ReportV3Packet")
}

func (h *Handler) FenbiPacket(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.FenbiPacket
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to FenbiPacket")
		return
	}

	log.WithFields(log.Fields{
		"date":   pkt.Header.Date,
		"close":  pkt.Header.Close,
		"open":   pkt.Header.Open,
		"count":  pkt.Count,
		"label":  TrimSpaceString(pkt.Header.Label),
		"market": pkt.Header.Market,
	}).Info("fenbi")

	for i, item := range pkt.Items {
		log.WithFields(log.Fields{
			"i":      i,
			"high":   item.High,
			"low":    item.Low,
			"open":   item.Open,
			"close":  item.Close,
			"ts":     TimeStringLong(item.Timestamp),
			"label":  TrimSpaceString(pkt.Header.Label),
			"market": pkt.Header.Market,
		}).Debug("fenbi.item")

		go func() {
			req := pb.CreateTickRequest{
				Timestamp:  TimestampProto(item.Timestamp),
				MarketId:   pkt.Header.Market,
				Code:       TrimSpaceString(pkt.Header.Label),
				Open:       item.Open,
				Close:      item.Close,
				High:       item.High,
				Low:        item.Low,
				Last:       item.Last,
				BuyPrice:   item.BuyPrice,
				BuyVolume:  item.BuyVolume,
				SellPrice:  item.SellPrice,
				SellVolume: item.SellVolume,
			}

			_, _ = h.cli.CreateTick(ctx, &req)
		}()
	}

	log.WithFields(log.Fields{
		"count": pkt.Count,
	}).Info("FenbiPacket")
}

func (h *Handler) MarketPacket(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.MarketPacket
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to MarketPacket")
		return
	}

	log.WithFields(log.Fields{
		"date":         pkt.Header.Date,
		"period_count": pkt.Header.PeriodCount,
		"count":        pkt.Count,
		"market":       TrimSpaceString(pkt.Header.Name),
	}).Info("market")

	for i, item := range pkt.Items {
		log.WithFields(log.Fields{
			"i":        i,
			"market":   pkt.Header.Market,
			"date":     pkt.Header.Date,
			"label":    TrimSpaceString(item.Label),
			"name":     TrimSpaceString(item.Name),
			"property": pkt.Items[i].Property,
		}).Debug("market.item")

		go func() {

			req := pb.CreateMarketRequest{
				Code:     TrimSpaceString(item.Label),
				Name:     TrimSpaceString(item.Name),
				MarketId: pkt.Header.Market,
				Date:     pkt.Header.Date,
			}

			_, _ = h.cli.CreateMarket(ctx, &req)
		}()
	}

	log.WithFields(log.Fields{
		"count": pkt.Count,
	}).Info("MarketPacket")
}

func (h *Handler) FinancePacket(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.FinancePacket
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to FinancePacket")
		return
	}

	for _, item := range pkt.Items {
		go func() {
			req := pb.CreateFinanceRequest{
				Market:  item.Market,
				Code:    TrimSpaceString(item.Label),
				BGRQ:    item.BGRQ,    // 财务数据的日期 如半年报 季报等 如 20090630 表示 2009年半年报
				ZGB:     item.ZGB,     // 总股本
				GJG:     item.GJG,     // 国家股
				FQFRG:   item.FQFRG,   // 发起人法人股
				FRG:     item.FRG,     // 法人股
				BGS:     item.BGS,     // B股
				HGS:     item.HGS,     // H股
				MQLT:    item.MQLT,    // 目前流通
				ZGG:     item.ZGG,     // 职工股
				A2ZPG:   item.A2ZPG,   // A2转配股
				ZZC:     item.ZZC,     // 总资产(千元)
				LDZC:    item.LDZC,    // 流动资产
				GDZC:    item.GDZC,    // 固定资产
				WXZC:    item.WXZC,    // 无形资产
				CQTZ:    item.CQTZ,    // 长期投资
				LDFZ:    item.LDFZ,    // 流动负债
				CQFZ:    item.CQFZ,    // 长期负债
				ZBGJJ:   item.ZBGJJ,   // 资本公积金
				MGGJJ:   item.MGGJJ,   // 每股公积金
				GDQY:    item.GDQY,    // 股东权益
				ZYSR:    item.ZYSR,    // 主营收入
				ZYLR:    item.ZYLR,    // 主营利润
				QTLR:    item.QTLR,    // 其他利润
				YYLR:    item.YYLR,    // 营业利润
				TZSY:    item.TZSY,    // 投资收益
				BTSR:    item.BTSR,    // 补贴收入
				YYWSZ:   item.YYWSZ,   // 营业外收支
				SNSYTZ:  item.SNSYTZ,  // 上年损益调整
				LRZE:    item.LRZE,    // 利润总额
				SHLR:    item.SHLR,    // 税后利润
				JLR:     item.JLR,     // 净利润
				WFPLR:   item.WFPLR,   // 未分配利润
				MGWFP:   item.MGWFP,   // 每股未分配
				MGSY:    item.MGSY,    // 每股收益
				MGJZC:   item.MGJZC,   // 每股净资产
				TZMGJZC: item.TZMGJZC, // 调整每股净资产
				GDQYB:   item.GDQYB,   // 股东权益比
				JZCSYL:  item.JZCSYL,  // 净资收益率
			}

			_, _ = h.cli.CreateFinance(ctx, &req)
		}()
	}

	log.WithFields(log.Fields{
		"count": pkt.Count,
	}).Info("FinancePacket")
}

type fileHistories struct {
	Label string
	Items []packect.FileHistoryItem
}

func (h *Handler) FileHistoryPacket(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.FileHistoryPacket
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to FileHistoryPacket")
		return
	}

	rows := map[string]fileHistories{}
	var curLabel string

	for i, item := range pkt.Items {
		if item.Timestamp == -1 {

			curLabel = TrimSpaceString(pkt.Headers[i].Label)
			rows[curLabel] = fileHistories{Label: curLabel}
		} else {
			items := rows[curLabel].Items
			items = append(items, packect.FileHistoryItem{
				Timestamp: item.Timestamp,
				Open:      item.Open,
				Close:     item.Close,
				High:      item.High,
				Low:       item.Low,
				Amount:    item.Amount,
				Volume:    item.Volume,
				Advance:   item.Advance,
				Decline:   item.Decline,
			})
			rows[curLabel] = fileHistories{Label: curLabel, Items: items}
		}
	}

	for _, row := range rows {
		for _, item := range row.Items {
			go func() {
				req := pb.CreateFileHistoryRequest{
					Code:      row.Label,
					Timestamp: TimestampProto(item.Timestamp),
					Open:    item.Open,
					Close:   item.Close,
					High:    item.High,
					Low:     item.Low,
					Volume:  item.Volume,
					Amount:  item.Amount,
					Advance: item.Advance,
					Decline: item.Decline,
				}

				_, _ = h.cli.CreateFileHistory(ctx, &req)
			}()
		}
	}

	log.WithFields(log.Fields{
		"count": pkt.Count,
		"rows":  len(rows),
	}).Info("FileHistoryPacket")
}

type fileMinutes struct {
	Label string
	Items []packect.FileMinuteItem
}

func (h *Handler) FileMinutePacket(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.FileMinutePacket

	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to FileMinutePacket")
		return
	}
	rows := map[string]fileMinutes{}
	var curLabel string

	for i, item := range pkt.Items {
		if item.Timestamp == -1 {

			curLabel = TrimSpaceString(pkt.Headers[i].Label)
			rows[curLabel] = fileMinutes{Label: curLabel}
		} else {
			items := rows[curLabel].Items
			items = append(items, packect.FileMinuteItem{
				Timestamp: item.Timestamp,
				Price:     item.Price,
				Amount:    item.Amount,
				Volume:    item.Volume,
			})
			rows[curLabel] = fileMinutes{Label: curLabel, Items: items}
		}
	}

	for _, row := range rows {
		for _, item := range row.Items {
			go func() {
				
				req := pb.CreateFileMinuteRequest{
					Code:      row.Label,
					Timestamp: TimestampProto(item.Timestamp),
					Price:  item.Price,
					Volume: item.Volume,
					Amount: item.Amount,
				}

				_, _ = h.cli.CreateFileMinute(ctx, &req)
			}()
		}
	}

	log.WithFields(log.Fields{
		"count": pkt.Count,
		"rows":  len(rows),
	}).Info("FileMinutePacket")
}

type fileHISMinutes struct {
	Label string
	Items []packect.FileHISMinuteItem
}

func (h *Handler) FileHISMinutePackect(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.FileHISMinutePackect
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to FileHISMinutePackect")
		return
	}

	rows := map[string]fileHISMinutes{}
	var curLabel string

	for i, item := range pkt.Items {
		if item.Timestamp == -1 {

			curLabel = TrimSpaceString(pkt.Headers[i].Label)
			rows[curLabel] = fileHISMinutes{Label: curLabel}

		} else {
			items := rows[curLabel].Items
			items = append(items, packect.FileHISMinuteItem{
				Timestamp:  item.Timestamp,
				Open:       item.Open,
				Close:      item.Close,
				High:       item.High,
				Low:        item.Low,
				Amount:     item.Amount,
				Volume:     item.Volume,
				ActivedBuy: item.ActivedBuy,
			})
			rows[curLabel] = fileHISMinutes{Label: curLabel, Items: items}
		}
	}

	for _, row := range rows {
		for _, item := range row.Items {
			go func() {
				
				req := pb.CreateFileHISMinuteRequest{
					Code:      row.Label,
					Timestamp: TimestampProto(item.Timestamp),
					Open:             item.Open,
					Close:            item.Close,
					High:             item.High,
					Low:              item.Low,
					Amount:           item.Amount,
					Volume:           item.Volume,
					ActivedBuyVolume: item.ActivedBuy,
				}

				_, _ = h.cli.CreateFileHISMinute(ctx, &req)
			}()
		}
	}

	log.WithFields(log.Fields{
		"count": pkt.Count,
		"rows":  len(rows),
	}).Info("FileHISMinutePackect")
}

type filePowers struct {
	Label string
	Items []packect.FilePowerItem
}

func (h *Handler) FilePowerPacket(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.FilePowerPacket
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to FilePowerPacket")
		return
	}

	rows := map[string]filePowers{}
	var curLabel string

	for i, item := range pkt.Items {
		if item.Timestamp == -1 {
			curLabel = TrimSpaceString(pkt.Headers[i].Label)
			rows[curLabel] = filePowers{Label: curLabel}
		} else {
			items := rows[curLabel].Items
			items = append(items, packect.FilePowerItem{
				Timestamp: item.Timestamp,
				Give:      item.Give,
				Pei:       item.Pei,
				PeiPrice:  item.PeiPrice,
				Profit:    item.Profit,
			})
			rows[curLabel] = filePowers{Label: curLabel, Items: items}
		}
	}

	for _, row := range rows {
		for _, item := range row.Items {
			go func() {
				
				req := pb.CreateFilePowerRequest{
					Code:      row.Label,
					Timestamp: TimestampProto(item.Timestamp),
					Give:        item.Give,
					Volume: item.Pei,
					Price:  item.PeiPrice,
					Profit:      item.Profit,
				}

				_, _ = h.cli.CreateFilePower(ctx, &req)
			}()
		}
	}

	log.WithFields(log.Fields{
		"count": pkt.Count,
		"rows":  len(rows),
	}).Info("FilePowerPacket")
}

func (h *Handler) FileBase(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.FileHeaderEx
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to FileHeaderEx")
		return
	}

	log.WithFields(log.Fields{
		"attrib":    pkt.Attrib,
		"len":       pkt.Len,
		"serial.no": pkt.SerialNo,
		"file.name": TrimSpaceString(pkt.FileName),
	}).Debug("file.base")

	req := pb.CreateFileBaseRequest{
		Attrib:   pkt.Attrib,
		Length:   pkt.Len,
		SerialNo: pkt.SerialNo,
		FileName: TrimSpaceString(pkt.FileName),
	}
	_, _ = h.cli.CreateFileBase(ctx, &req)
}

func (h *Handler) FileNews(ctx context.Context, msg *amqp.Delivery) {
	var pkt packect.FileHeaderEx
	if err := json.Unmarshal(msg.Body, &pkt); err != nil {
		log.WithError(err).Error("fail to unmarshal bytes to FileHeaderEx")
		return
	}

	log.WithFields(log.Fields{
		"attrib":    pkt.Attrib,
		"len":       pkt.Len,
		"serial.no": pkt.SerialNo,
		"file.name": TrimSpaceString(pkt.FileName),
	}).Debug("file.news")

	req := pb.CreateFileNewsRequest{
		Attrib:   pkt.Attrib,
		Length:   pkt.Len,
		SerialNo: pkt.SerialNo,
		FileName: TrimSpaceString(pkt.FileName),
	}
	_, _ = h.cli.CreateFileNews(ctx, &req)
}
