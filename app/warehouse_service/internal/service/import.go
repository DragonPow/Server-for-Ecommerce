package service

import (
	"Server-for-Ecommerce/app/warehouse_service/api"
	"Server-for-Ecommerce/app/warehouse_service/internal/store"
	"Server-for-Ecommerce/app/warehouse_service/util"
	"Server-for-Ecommerce/library/slice"
	"context"
	"google.golang.org/grpc/codes"
	"strconv"
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

	var lastActionByName string
	var createByName string
	var uids []string
	if importBill.CreateUid.Valid {
		uids = append(uids, strconv.FormatInt(importBill.CreateUid.Int64, util.Base10Int))
	}
	if importBill.WriteUid.Valid {
		uids = append(uids, strconv.FormatInt(importBill.WriteUid.Int64, util.Base10Int))
	}
	// Get LastActionByName and CreateByName
	//if len(uids) > util.ZeroLength {
	//	userRes, err := s.accountClient.GetUsers(ctx, &accountApi.GetUsersRequest{Ids: strings.Join(slice.Uniq(uids), ",")})
	//	if err != nil {
	//		logger.Error(err, "Call GetUsers fail")
	//		return util.Response(err, defaultResponse)
	//	}
	//	if userRes.Data != nil && len(userRes.Data.Items) > util.ZeroLength {
	//		for _, user := range userRes.Data.Items {
	//			// Set create by
	//			if user.UserId == importBill.CreateUid.Int64 {
	//				createByName = user.Name
	//			}
	//			// Set write by
	//			if user.UserId == importBill.WriteUid.Int64 {
	//				lastActionByName = user.Name
	//			}
	//		}
	//	}
	//}

	// TODO: Call another service to get ProductName

	return &api.GetImportBillResponse{
		Code:    int32(codes.OK),
		Message: "OK",
		Data: &api.GetImportBillResponseData{
			Item: &api.GetImportBillItem{
				Id:               importBill.ID,
				Code:             importBill.Code,
				LastActionById:   importBill.WriteUid.Int64,
				LastActionByName: lastActionByName,
				CreateById:       importBill.CreateUid.Int64,
				CreateByName:     createByName,
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
