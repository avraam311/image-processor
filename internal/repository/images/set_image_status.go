package images

import (
	"context"
	"fmt"
)

func (r *Repository) SetImageStatus(ctx context.Context, status string) (uint, error) {
	query := `
		INSERT INTO image (status)
		VALUES ($1)
		RETURNING id;
	`

	var id uint
	err := r.db.QueryRowContext(ctx, query, status).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repository/set_image_status.go - failed to scan id")
	}

	return id, nil
}
