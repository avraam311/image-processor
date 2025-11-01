package images

import (
	"context"
	"fmt"
)

func (r *Repository) ChangeImageStatus(ctx context.Context, id uint, status string) error {
	query := `
		UPDATE image
		SET status = $2
		WHERE id = $1;
	`

	res, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("repository/change_image_status.go - failed to change image status - %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrImageNotFound
	}

	return nil
}
