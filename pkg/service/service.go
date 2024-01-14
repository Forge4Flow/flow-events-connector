package service

import "github.com/forge4flow/forge4flow-core/pkg/database"

const (
	EventService       = "EventService"
	WarrantService     = "WarrantService"
	UserService        = "UserService"
	TenantService      = "TenantService"
	RoleService        = "RoleService"
	PricingTierService = "PricingTierService"
	PermissionService  = "PermissionService"
	ObjectTypeService  = "ObjectTypeService"
	ObjectSerice       = "ObjectSerice"
	FeatureService     = "FeatureService"
	CheckService       = "CheckService"
	SessionService     = "SessionService"
	NonceService       = "NonceService"
	FlowService        = "FlowService"
	ApiService         = "ApiService"
)

type Env interface {
	DB() database.Database
	EventDB() database.Database
}

type Service interface {
	Routes() ([]Route, error)
	Env() Env
	ID() string
}

type BaseService struct {
	env Env
}

func (svc BaseService) Env() Env {
	return svc.env
}

func NewBaseService(env Env) BaseService {
	return BaseService{
		env: env,
	}
}
