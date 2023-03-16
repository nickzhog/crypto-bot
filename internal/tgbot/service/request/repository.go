package request

import "context"

type Repository interface {
	Create(ctx context.Context, usrID int64, cryptoID, price string) error
	FindForUser(ctx context.Context, usrID int64) ([]Request, error)
}
