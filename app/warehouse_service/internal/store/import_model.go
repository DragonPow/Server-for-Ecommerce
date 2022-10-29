package store

type GetImportDataBillParams struct {
	ImportId int64
}

type GetImportDataBillResponse struct {
	*ImportBill
	DetailItems []*ImportBillDetail
}
