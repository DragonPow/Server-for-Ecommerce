package store

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app/account_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/slice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) GetUsers(ctx context.Context, params GetUsersParams) (*GetUsersResponse, error) {
	// Get Customers
	customers, err := s.getCustomerByIds(ctx, params.CustomerIds)
	if err != nil {
		s.log.Error(err, "Call getCustomerByIds fail")
		return nil, status.Error(codes.Internal, "Tìm danh sách khách hàng lỗi")
	}

	// Get Accounts
	if len(customers) == util.ZeroLength {
		return &GetUsersResponse{Users: []*GetUsersItem{}}, nil
	}
	accounts, err := s.getAccountByIds(ctx, slice.Map(customers, func(customer CustomerInfo) int64 { return customer.AccountID }))
	if err != nil {
		s.log.Error(err, "Call getAccountByIds fail")
		return nil, status.Error(codes.Internal, "Tìm danh sách tài khoản lỗi")
	}

	// Map account to Customer
	mapAccountId := slice.KeyBy(accounts, func(account getAccountByIdsRow) (int64, *getAccountByIdsRow) {
		return account.ID, &account
	})
	response := &GetUsersResponse{}
	for _, customer := range customers {
		account, ok := mapAccountId[customer.AccountID]
		if !ok {
			errRes := status.Errorf(codes.NotFound, "Không tìm thấy thông tin tài khoản của khách hàng %v", customer.Name)
			s.log.Error(errRes, "Not found account", "accounts", accounts, "customers", customers)
			return nil, errRes
		}
		response.Users = append(response.Users, &GetUsersItem{
			AccountId:  account.ID,
			CustomerId: customer.ID,
			Name:       customer.Name,
			UserName:   account.Username,
			Phone:      customer.Phone.String,
			Address:    customer.Address.String,
			CreateDate: customer.CreateDate,
			WriteDate:  customer.WriteDate,
		})
	}

	return response, nil
}
