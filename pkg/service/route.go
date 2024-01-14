package service

import "net/http"

type Route interface {
	GetPattern() string
	GetMethod() string
	GetHandler() http.Handler
	GetOverrideAuthMiddlewareFunc() AuthMiddlewareFunc
}

type ForgeRoute struct {
	Pattern                    string
	Method                     string
	Handler                    http.Handler
	OverrideAuthMiddlewareFunc AuthMiddlewareFunc
}

func (route ForgeRoute) GetPattern() string {
	return route.Pattern
}

func (route ForgeRoute) GetMethod() string {
	return route.Method
}

func (route ForgeRoute) GetHandler() http.Handler {
	return route.Handler
}

func (route ForgeRoute) GetOverrideAuthMiddlewareFunc() AuthMiddlewareFunc {
	return route.OverrideAuthMiddlewareFunc
}
