package helper

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func HandleDbErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return errors.New(pgErr.Message)
	}
	return err
}
