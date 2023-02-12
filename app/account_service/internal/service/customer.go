package service

import (
	"Server-for-Ecommerce/app/account_service/api"
	"Server-for-Ecommerce/app/account_service/internal/store"
	"Server-for-Ecommerce/app/account_service/util"
	"Server-for-Ecommerce/library/slice"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

func (s *Service) GetUsers(ctx context.Context, request *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	defaultResponse := func(code int32, message string) *api.GetUsersResponse {
		return &api.GetUsersResponse{
			Code:    code,
			Message: message,
			Data:    nil,
		}
	}
	logger := s.log.WithName("GetUsers").WithValues("request", request)
	logger.Info("Start process")
	defer logger.Info("End process")

	// Validate data
	accountIdStrings := strings.Split(request.Ids, ",")
	if len(accountIdStrings) == util.ZeroLength || accountIdStrings[0] == util.EmptyString {
		return nil, status.Errorf(codes.InvalidArgument, "Phải nhập thông tin Id user")
	}
	accountIds := make([]int64, len(accountIdStrings), len(accountIdStrings))
	for i, idString := range accountIdStrings {
		id, err := strconv.ParseInt(idString, util.Base10Int, util.Bit64Size)
		if err != nil {
			logger.Error(err, "Parse string to int64 fail")
			return nil, status.Errorf(codes.InvalidArgument, "Id phải có định dạng kiểu số nguyên")
		}
		accountIds[i] = id
	}

	customers, err := s.storeDb.GetUsers(ctx, store.GetUsersParams{CustomerIds: accountIds})
	if err != nil {
		logger.Error(err, "Call store GetUsers fail")
		return util.Response(err, defaultResponse)
	}
	return &api.GetUsersResponse{
		Code:    int32(codes.OK),
		Message: "OK",
		Data: &api.GetUsersResponseData{
			TotalItems: int32(len(customers.Users)),
			Items: slice.Map(customers.Users, func(user *store.GetUsersItem) *api.GetUsersResponseItem {
				return &api.GetUsersResponseItem{
					UserId:   user.CustomerId,
					Username: user.UserName,
					Name:     user.Name,
					Phone:    user.Phone,
					Address:  user.Address,
				}
			}),
		},
	}, nil
}
