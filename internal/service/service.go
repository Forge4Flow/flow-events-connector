package service

import "flow-events-connector/internal/database"

type Service interface {
	DB() database.Database
}

type BaseService struct {
	db database.Database
}

func (svc BaseService) DB() database.Database {
	return svc.db
}

// func NewBaseService(db database.Database) BaseService {
func NewBaseService() BaseService {
	// return BaseService{
	// 	db: db,
	// }
	return BaseService{}
}
