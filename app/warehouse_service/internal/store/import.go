package store

import (
	"context"
	"fmt"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/slice"
)

func (s *Store) GetImportDataBill(ctx context.Context, params GetImportDataBillParams) (*GetImportDataBillResponse, error) {
	rows, err := s.query(ctx, nil, `select pg_read_file('/etc/hostname') as hostname, setting as port from pg_settings where name='port';`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var dest1 interface{}
		var dest2 interface{}

		if err := rows.Scan(
			&dest1,
			&dest2,
		); err != nil {
			return nil, err
		}
		s.log.Info(fmt.Sprintf("Dest id: %v", dest1))
	}

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
