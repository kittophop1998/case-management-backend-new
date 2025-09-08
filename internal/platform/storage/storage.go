package storage

import (
	"case-management/infrastructure/config"

	"github.com/minio/minio-go"
)

var Storage *minio.Client

func InitStorage(cfg config.IsilonConfig) {
	var err error

	Storage, err = minio.New(cfg.BaseURL, cfg.AccessKey, cfg.SecretKey, true)
	if err != nil {
		panic("cannot connect storage (minio) >> " + err.Error())
	}
}
