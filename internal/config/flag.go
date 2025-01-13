package config

import "flag"

const (
	flagRunAddress           = "a"
	flagDatabaseURI          = "d"
	flagAccrualAddress       = "r"
	flagAccrualLocation      = "k"
	flagTickerTimer          = "t"
	flagAccrualClientTimeout = "m"
)

type flagsValue struct {
	runAddress           string
	databaseURI          string
	accrualAddress       string
	accrualFileLocation  string
	tickerTimer          string
	accrualClientTimeout string
}

func registerFlag(name, defaultValue, usage string, target *string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(target, name, defaultValue, usage)
	} else {
		*target = flag.Lookup(name).Value.(flag.Getter).Get().(string)
	}
}

func (v *flagsValue) parseFlagsIfNeeded() {
	registerFlag(
		flagRunAddress,
		"",
		"the address and port to listen on for HTTP requests. In format: \"address:port\"",
		&v.runAddress,
	)
	registerFlag(
		flagDatabaseURI,
		"",
		"database DSN. In format: \"postgres://username:password@host:port/databaseName\"",
		&v.databaseURI,
	)
	registerFlag(
		flagAccrualAddress,
		"",
		"address for accrual api. In format: \"address:port\"",
		&v.accrualAddress,
	)
	registerFlag(
		flagAccrualLocation,
		"",
		"accrual file location. In format: \"./cmd/accrual/bin\"",
		&v.accrualFileLocation,
	)
	registerFlag(
		flagTickerTimer,
		"",
		"ticker for accrual worker, needed for fetching new statuses. In format: 10s",
		&v.tickerTimer,
	)
	registerFlag(
		flagAccrualClientTimeout,
		"",
		"timeout for client, making requests accrual system. In format: 10s",
		&v.tickerTimer,
	)

	flag.Parse()
}

func newFlagsValue() flagsValue {
	var f flagsValue
	f.parseFlagsIfNeeded()
	return f
}
