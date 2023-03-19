package server

import (
	"github.com/qiniu/qmgo"
)

type Config struct {
	Database *qmgo.Database
	GRPCPort string
}
