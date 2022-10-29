package store

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/slice"
)

func (s *Store) GetImportDataBill(ctx context.Context, params GetImportDataBillParams) (*GetImportDataBillResponse, error) {
	// Get Import
	imports, err := s.getImportBills(ctx, []int64{params.ImportId})
	if err != nil {
		s.log.Error(err, "Call getImportBills fail")
		return nil, err
	}
	if len(imports) == util.ZeroLength {
		return nil, util.ErrorGrpcf(util.NotFoundModel, "Không tìm thấy import với id = %v", params.ImportId)
	}
	// Get Detail Imports
	importDetails, err := s.getImportBillDetails(ctx, getImportBillDetailsParams{
		ImportID: []int64{params.ImportId},
	})
	if err != nil {
		s.log.Error(err, "Call getImportBillDetails fail")
		return nil, err
	}
	return &GetImportDataBillResponse{
		ImportBill: &imports[0],
		DetailItems: slice.Map(importDetails, func(detail ImportBillDetail) *ImportBillDetail {
			return &detail
		}),
	}, nil
}
