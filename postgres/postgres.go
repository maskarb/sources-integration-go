package postgres

import (
	"net"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/pgext"
	"github.com/maskarb/sources-integration-go/xconfig"
)

var (
	pgMainOnce sync.Once
	pgMain     *pg.DB
)

func PGMain() *pg.DB {
	pgMainOnce.Do(func() {
		pgMain = NewPostgres(xconfig.Config.PGMain, hasPgbouncer())
	})
	return pgMain
}

var (
	pgMainTxOnce sync.Once
	pgMainTx     *pg.DB
)

func PGMainTx() *pg.DB {
	pgMainTxOnce.Do(func() {
		pgMainTx = NewPostgres(xconfig.Config.PGMain, false)
	})
	return pgMainTx
}

func hasPgbouncer() bool {
	switch xconfig.Config.Env {
	case "test", "dev":
		return false
	default:
		return true
	}
}

func NewPostgres(cfg *xconfig.Postgres, usePool bool) *pg.DB {
	addr := cfg.Addr
	if usePool {
		addr = replacePort(addr, cfg.ConnectionPoolPort)
	}

	opt := cfg.Options()
	opt.Addr = addr

	db := pg.Connect(opt)
	defer db.Close()

	db.AddQueryHook(pgext.OpenTelemetryHook{})

	return db
}

func replacePort(s, newPort string) string {
	if newPort == "" {
		return s
	}
	host, _, err := net.SplitHostPort(s)
	if err != nil {
		host = s
	}
	return net.JoinHostPort(host, newPort)
}
