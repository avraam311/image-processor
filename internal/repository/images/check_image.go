package images

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	statusInProcess = "in process"
)

func (r *Repository) CheckImage(ctx context.Context, id uint) error {
	query := `
		SELECT status
		FROM image
		WHERE id = $1;
	`

	var status string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrImageNotFound
		}

		return fmt.Errorf("repository/check_image.go - failed to check image - %w", err)
	}
	if status == statusInProcess {
		return ErrImageInProcess
	}

	return nil
}
