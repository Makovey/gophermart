package config

import "flag"

const (
	flagRunAddress     = "a"
	flagDatabaseURI    = "d"
	flagAccrualAddress = "r"
)

type flagsValue struct {
	runAddress     string
	databaseURI    string
	accrualAddress string
}

func (v *flagsValue) parseFlagsIfNeeded() {
	if flag.Lookup(flagRunAddress) == nil {
		flag.StringVar(
			&v.runAddress,
			flagRunAddress,
			"",
			"the address and port to listen on for HTTP requests. In format: \"address:port\"",
		)
	} else {
		v.runAddress = flag.Lookup(flagRunAddress).Value.(flag.Getter).Get().(string)
	}

	if flag.Lookup(flagDatabaseURI) == nil {
		flag.StringVar(
			&v.databaseURI,
			flagDatabaseURI,
			"",
			"database DSN. In format: \"postgres://username:password@host:port/databaseName\"",
		)
	} else {
		v.databaseURI = flag.Lookup(flagDatabaseURI).Value.(flag.Getter).Get().(string)
	}

	if flag.Lookup(flagAccrualAddress) == nil {
		flag.StringVar(&v.accrualAddress, flagAccrualAddress, "", "address for accrual api. In format: \"address:port\"")
	} else {
		v.accrualAddress = flag.Lookup(flagAccrualAddress).Value.(flag.Getter).Get().(string)
	}

	flag.Parse()
}

func newFlagsValue() flagsValue {
	var f flagsValue
	f.parseFlagsIfNeeded()

	return f
}
