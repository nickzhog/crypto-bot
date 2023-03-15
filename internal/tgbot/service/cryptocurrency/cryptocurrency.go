package cryptocurrency

import (
	"fmt"
)

type Currency struct {
	Name                string
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
