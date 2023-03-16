package cryptocurrency

import "context"

type Repository interface {
	UpsertMany(ctx context.Context, currencies []Currency) error
	FindCurrencies(ctx context.Context, index int) ([]Currency, error)
	FindOne(ctx context.Context, name string) (Currency, error)
}
