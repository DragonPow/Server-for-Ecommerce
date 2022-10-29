package service

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/internal/store"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/slice"
	"google.golang.org/grpc/codes"
)

func (s *Service) CreateImportBill(ctx context.Context, request *api.CreateImportBillRequest) (*api.CreateImportBillResponse, error) {
	return &api.CreateImportBillResponse{
		Code:    int32(codes.OK),
		Message: request.String(),
		Data:    nil,
	}, nil
}

func (s Service) GetImportBill(ctx context.Context, request *api.GetImportBillRequest) (*api.GetImportBillResponse, error) {
	defaultResponse := func(code int32, message string) *api.GetImportBillResponse {
		return &api.GetImportBillResponse{
			Code:    code,
			Message: message,
			Data:    nil,
		}
	}
	logger := s.log.WithName("GetImportBill").WithValues("request", request)
	logger.Info("Start process")
	defer logger.Info("End process")

	importBill, err := s.storeDb.GetImportDataBill(ctx, store.GetImportDataBillParams{ImportId: request.ImportId})
	if err != nil {
		logger.Error(err, "Call GetImportDataBill fail")
		return util.Response(err, defaultResponse)
	}

	// TODO: Call another service to get LastActionByName and CreateByName

	// TODO: Call another service to get ProductName

	return &api.GetImportBillResponse{
		Code:    int32(codes.OK),
		Message: "OK",
		Data: &api.GetImportBillResponseData{
			Item: &api.GetImportBillItem{
				Id:               importBill.ID,
				Code:             importBill.Code,
				LastActionById:   importBill.WriteUid.Int64,
				LastActionByName: "",
				CreateById:       importBill.CreateUid.Int64,
				CreateByName:     "",
				ItemDetails: slice.Map(importBill.DetailItems, func(item *store.ImportBillDetail) *api.GetImportBillItemDetail {
					return &api.GetImportBillItemDetail{
						ProductId:   item.ProductID,
						ProductName: "",
						UomName:     "",
						Quantity:    item.Quantity,
					}
				}),
			},
		},
	}, nil
}
