package server

import "context"

// Server 对 Server 的生命周期进行管理
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}
