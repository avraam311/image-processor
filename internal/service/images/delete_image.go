package images

import (
	"context"
	"fmt"
)

func (s *Service) DeleteImage(ctx context.Context, id uint) error {
	err := s.repo.DeleteImage(ctx, id)
	if err != nil {
		return fmt.Errorf("service/images - %w", err)
	}

	return nil
}
