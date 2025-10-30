package images

import (
	"context"
	"fmt"
)

const (
	rowsAffected = 0
)

func (r *Repository) ChangeImageStatus(ctx context.Context, id uint) error {
	query := `
		DELETE
		FROM image
		WHERE id = $1;
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository/change_image_status.go - failed to delete image - %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == rowsAffected {
		return ErrImageNotFound
	}

	return nil
}
