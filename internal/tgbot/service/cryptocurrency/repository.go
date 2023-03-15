package cryptocurrency

import "context"

type Repository interface {
	UpsertMany(ctx context.Context, currencies []Currency) error
}
