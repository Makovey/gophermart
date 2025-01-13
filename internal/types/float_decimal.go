package types

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type FloatDecimal decimal.Decimal

func (f FloatDecimal) MarshalJSON() ([]byte, error) {
	dec := decimal.Decimal(f)
	floatVal, _ := dec.Float64()
	return json.Marshal(floatVal)
}

func (f *FloatDecimal) UnmarshalJSON(data []byte) error {
	if f == nil {
		return nil
	}

	var floatVal float64
	if err := json.Unmarshal(data, &floatVal); err == nil {
		*f = FloatDecimal(decimal.NewFromFloat(floatVal))
		return nil
	}

	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return fmt.Errorf("can' unmarshal to string: %v", err)
	}

	dec, err := decimal.NewFromString(strVal)
	if err != nil {
		return fmt.Errorf("can't create decimal from string %s", err)
	}

	*f = FloatDecimal(dec)
	return nil
}
