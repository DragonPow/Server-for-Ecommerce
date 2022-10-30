package store

import "time"

type GetUsersParams struct {
	CustomerIds []int64
}

type GetUsersResponse struct {
	Users []*GetUsersItem
}

type GetUsersItem struct {
	AccountId  int64
	CustomerId int64
	Name       string
	UserName   string
	Phone      string
	Address    string
	CreateDate time.Time
	WriteDate  time.Time
}
