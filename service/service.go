package service

import (
	"context"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/emuta/haililive/model"
	pb "github.com/emuta/haililive/protobuf/haililive"
	"github.com/emuta/haililive/uniqueid"
)

type HaiLiLiveServiceServerImpl struct {
	db *gorm.DB
}

func NewHaiLiLiveServiceServerImpl(db *gorm.DB) *HaiLiLiveServiceServerImpl {
	return &HaiLiLiveServiceServerImpl{db: db}
}

func (s *HaiLiLiveServiceServerImpl) CreateReportV3(ctx context.Context, req *pb.CreateReportV3Request) (*pb.Response, error) {
	report := model.ReportV3{
		MarketId: req.MarketId,
		Code:     req.Code,
		Name:     req.Name,
		Last:     req.Last,
		Open:     req.Open,
		Close:    req.Close,
		High:     req.High,
		Low:      req.Low,
		Amount:   req.Amount,
		Volume:   req.Volume,
	}

	for _, v := range req.BuyPrice {
		report.BuyPrice = append(report.BuyPrice, v)
	}

	for _, v := range req.BuyVolume {
		report.BuyVolume = append(report.BuyVolume, v)
	}

	for _, v := range req.SellPrice {
		report.SellPrice = append(report.SellPrice, v)
	}

	for _, v := range req.SellVolume {
		report.SellVolume = append(report.SellVolume, v)
	}

	if ts, err := ptypes.Timestamp(req.Timestamp); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	} else {
		report.Timestamp = ts
	}

	report.Id = uniqueid.New(report.Timestamp, report.Code)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.ReportV3{}, "id = ?", report.Id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return tx.Create(&report).Error
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Response{Success: true}, nil
}

func (s *HaiLiLiveServiceServerImpl) CreateTick(ctx context.Context, req *pb.CreateTickRequest) (*pb.Response, error) {
	fenbi := model.Fenbi{
		MarketId: req.MarketId,
		Code:     req.Code,
		Last:     req.Last,
		Open:     req.Open,
		Close:    req.Close,
		High:     req.High,
		Low:      req.Low,
	}

	for _, v := range req.BuyPrice {
		fenbi.BuyPrice = append(fenbi.BuyPrice, v)
	}

	for _, v := range req.BuyVolume {
		fenbi.BuyVolume = append(fenbi.BuyVolume, v)
	}

	for _, v := range req.SellPrice {
		fenbi.SellPrice = append(fenbi.SellPrice, v)
	}

	for _, v := range req.SellVolume {
		fenbi.SellVolume = append(fenbi.SellVolume, v)
	}

	if ts, err := ptypes.Timestamp(req.Timestamp); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	} else {
		fenbi.Timestamp = ts
	}

	fenbi.Id = uniqueid.New(fenbi.Timestamp, fenbi.Code)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.Fenbi{}, "id = ?", fenbi.Id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return tx.Create(&fenbi).Error
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Response{Success: true}, nil
}

func (s *HaiLiLiveServiceServerImpl) CreateMarket(ctx context.Context, req *pb.CreateMarketRequest) (*pb.Response, error) {
	market := model.Market{
		MarketId: req.MarketId,
		Date:     req.Date,
		Code:     req.Code,
		Name:     req.Name,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.Market{}, "code = ?", market.Code).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return tx.Create(&market).Error
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Response{Success: true}, nil
}

func (s *HaiLiLiveServiceServerImpl) CreateFinance(ctx context.Context, req *pb.CreateFinanceRequest) (*pb.Response, error) {
	finance := model.Finance{
		Market:  req.Market,
		Code:    req.Code,
		BGRQ:    req.BGRQ,
		ZGB:     req.ZGB,
		GJG:     req.GJG,
		FQFRG:   req.FQFRG,
		FRG:     req.FRG,
		BGS:     req.BGS,
		HGS:     req.HGS,
		MQLT:    req.MQLT,
		ZGG:     req.ZGG,
		A2ZPG:   req.A2ZPG,
		ZZC:     req.ZZC,
		LDZC:    req.LDZC,
		GDZC:    req.GDZC,
		WXZC:    req.WXZC,
		CQTZ:    req.CQTZ,
		LDFZ:    req.LDFZ,
		CQFZ:    req.CQFZ,
		ZBGJJ:   req.ZBGJJ,
		MGGJJ:   req.MGGJJ,
		GDQY:    req.GDQY,
		ZYSR:    req.ZYSR,
		ZYLR:    req.ZYLR,
		QTLR:    req.QTLR,
		YYLR:    req.YYLR,
		TZSY:    req.TZSY,
		BTSR:    req.BTSR,
		YYWSZ:   req.YYWSZ,
		SNSYTZ:  req.SNSYTZ,
		LRZE:    req.LRZE,
		SHLR:    req.SHLR,
		JLR:     req.JLR,
		WFPLR:   req.WFPLR,
		MGWFP:   req.MGWFP,
		MGSY:    req.MGSY,
		MGJZC:   req.MGJZC,
		TZMGJZC: req.TZMGJZC,
		GDQYB:   req.GDQYB,
		JZCSYL:  req.JZCSYL,
	}

	date, err := time.Parse(strconv.FormatInt(finance.BGRQ, 10), "20060102")
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	finance.Id = uniqueid.New(date, finance.Code)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.Finance{}, "id = ?", finance.Id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return tx.Create(&finance).Error
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Response{Success: true}, nil
}

func (s *HaiLiLiveServiceServerImpl) CreateFileHistory(ctx context.Context, req *pb.CreateFileHistoryRequest) (*pb.Response, error) {
	m := model.FileHistory{
		Code:    req.Code,
		Open:    req.Open,
		Close:   req.Close,
		High:    req.High,
		Low:     req.Low,
		Volume:  req.Volume,
		Amount:  req.Amount,
		Advance: req.Advance,
		Decline: req.Decline,
	}

	if ts, err := ptypes.Timestamp(req.Timestamp); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	} else {
		m.Timestamp = ts
	}

	m.Id = uniqueid.New(m.Timestamp, m.Code)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.FileHistory{}, "id = ?", m.Id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return tx.Create(&m).Error
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Response{Success: true}, nil
}

func (s *HaiLiLiveServiceServerImpl) CreateFileMinute(ctx context.Context, req *pb.CreateFileMinuteRequest) (*pb.Response, error) {
	m := model.FileMinute{
		Code:   req.Code,
		Price:  req.Price,
		Volume: req.Volume,
		Amount: req.Amount,
	}

	if ts, err := ptypes.Timestamp(req.Timestamp); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	} else {
		m.Timestamp = ts
	}

	m.Id = uniqueid.New(m.Timestamp, m.Code)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.FileMinute{}, "id = ?", m.Id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return tx.Create(&m).Error
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Response{Success: true}, nil
}

func (s *HaiLiLiveServiceServerImpl) CreateFileHISMinute(ctx context.Context, req *pb.CreateFileHISMinuteRequest) (*pb.Response, error) {
	m := model.FileHISMinute{
		Code:             req.Code,
		Open:             req.Open,
		Close:            req.Close,
		High:             req.High,
		Low:              req.Low,
		Amount:           req.Amount,
		Volume:           req.Volume,
		ActivedBuyVolume: req.ActivedBuyVolume,
	}

	if ts, err := ptypes.Timestamp(req.Timestamp); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	} else {
		m.Timestamp = ts
	}

	m.Id = uniqueid.New(m.Timestamp, m.Code)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.FileHISMinute{}, "id = ?", m.Id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return tx.Create(&m).Error
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Response{Success: true}, nil
}

func (s *HaiLiLiveServiceServerImpl) CreateFilePower(ctx context.Context, req *pb.CreateFilePowerRequest) (*pb.Response, error) {
	m := model.FilePower{
		Code:   req.Code,
		Give:   req.Give,
		Price:  req.Price,
		Volume: req.Volume,
		Profit: req.Profit,
	}

	if ts, err := ptypes.Timestamp(req.Timestamp); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	} else {
		m.Timestamp = ts
	}

	m.Id = uniqueid.New(m.Timestamp, m.Code)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model.FilePower{}, "id = ?", m.Id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return tx.Create(&m).Error
			} else {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.Response{Success: true}, nil
}

func (s *HaiLiLiveServiceServerImpl) CreateFileBase(ctx context.Context, req *pb.CreateFileBaseRequest) (*pb.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFileBase not implemented")
}

func (s *HaiLiLiveServiceServerImpl) CreateFileNews(ctx context.Context, req *pb.CreateFileNewsRequest) (*pb.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFileNews not implemented")
}
