package app

import (
	"errors"
	"time"

	"github.com/MaximBayurov/rate-limiter/internal/configuration"
	"github.com/MaximBayurov/rate-limiter/internal/iplists"
	"github.com/MaximBayurov/rate-limiter/internal/logger"
	"github.com/jmoiron/sqlx"
)

type App interface {
	AddIP(ip string, listType string, overwrite bool) error
	DeleteIP(ip string, listType string) error
	TryLogin(login, password, ip string) error
	ClearLoginAttempts(login, ip string) error
}

type BruteforceApp struct {
	logger    logger.Logger
	configs   configuration.AppConf
	ipList    iplists.IPList
	logins    bucketsPull
	passwords bucketsPull
	ips       bucketsPull
}

func New(logger logger.Logger, db *sqlx.DB, configs configuration.AppConf) App {
	ipList := iplists.NewIPList(db, logger)

	application := &BruteforceApp{
		logger:    logger,
		configs:   configs,
		ipList:    ipList,
		logins:    newBucketsPull(configs.LoginAttempts),
		passwords: newBucketsPull(configs.PasswordAttempts),
		ips:       newBucketsPull(configs.IPAttempts),
	}
	go application.CollectGarbage()
	return application
}

func (a *BruteforceApp) AddIP(ip, listType string, overwrite bool) (err error) {
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

func (a *BruteforceApp) DeleteIP(ip, listType string) (err error) {
	var ipListType iplists.ListType
	var ok bool
	if ipListType, ok = iplists.ParseType(listType); !ok {
		return ErrUnsupportedIPListType
	}

	err = a.ipList.Delete(ip, ipListType)
	return err
}

func (a *BruteforceApp) TryLogin(login, password, ip string) error {
	var err error

	loginOk := a.logins.Allow(login)
	if !loginOk {
		err = ErrLoginDisallow
	}
	passwordOk := a.passwords.Allow(password)
	if !passwordOk {
		err = ErrPasswordDisallow
	}
	ipOk := a.ips.Allow(ip)
	if !ipOk {
		err = ErrIPDisallow
	}

	listType, iperr := a.ipList.In(ip)
	if iperr == nil {
		switch listType {
		case iplists.Black:
			err = ErrIPInBlackList
		case iplists.White:
			err = nil
		}
	}

	return err
}

func (a *BruteforceApp) ClearLoginAttempts(login, ip string) error {
	a.logins.DeleteBucket(login)
	a.ips.DeleteBucket(ip)

	return nil
}

// CollectGarbage функция для систематической очистки мусорных bucket.
func (a *BruteforceApp) CollectGarbage() {
	t := time.NewTicker(time.Minute)
	for {
		<-t.C
		a.logins.ClearEmptyBuckets()
		a.passwords.ClearEmptyBuckets()
		a.ips.ClearEmptyBuckets()
	}
}
