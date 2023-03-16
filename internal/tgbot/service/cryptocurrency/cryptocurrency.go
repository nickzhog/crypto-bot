package cryptocurrency

import (
	"fmt"
	"time"
)

type Currency struct {
	Name                string
	LastUpdate          time.Time
	PriceToDollarString string
	PriceToDollarFloat  float64
}

func NewCurrency(name string, price float64) Currency {
	return Currency{
		Name:                name,
		PriceToDollarString: fmt.Sprintf("%g", price),
		PriceToDollarFloat:  price,
	}
}
