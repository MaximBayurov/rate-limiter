package app

import (
	"errors"

	"github.com/MaximBayurov/rate-limiter/internal/configuration"
	"github.com/MaximBayurov/rate-limiter/internal/iplists"
	"github.com/MaximBayurov/rate-limiter/internal/logger"
	"github.com/jmoiron/sqlx"
)

type App interface {
	AddIP(ip string, listType string, overwrite bool) error
	DeleteIP(ip string, listType string) error
}

type BruteforceApp struct {
	logger  logger.Logger
	configs configuration.AppConf
	ipList  iplists.IPList
}

func New(logger logger.Logger, db *sqlx.DB, configs configuration.AppConf) App {
	ipList := iplists.NewIPList(db, logger)
	return BruteforceApp{
		logger:  logger,
		configs: configs,
		ipList:  ipList,
	}
}

func (a BruteforceApp) AddIP(ip, listType string, overwrite bool) (err error) {
	var ipListType iplists.ListType
	var ok bool
	if ipListType, ok = iplists.ParseType(listType); !ok {
		return ErrUnsupportedIPListType
	}

	err = a.ipList.Add(ip, ipListType)
	if errors.Is(err, iplists.ErrNotInserted) && overwrite {
		err = a.ipList.Update(ip, ipListType)
	}

	return err
}

func (a BruteforceApp) DeleteIP(ip, listType string) (err error) {
	var ipListType iplists.ListType
	var ok bool
	if ipListType, ok = iplists.ParseType(listType); !ok {
		return ErrUnsupportedIPListType
	}

	err = a.ipList.Delete(ip, ipListType)
	return err
}
