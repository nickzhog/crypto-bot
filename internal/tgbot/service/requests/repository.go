package requests

import "context"

type Repository interface {
	Create(ctx context.Context, usrID int64, cryptoID string) error
}
