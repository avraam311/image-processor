package images

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/avraam311/image-processor/internal/api/handlers"
	"github.com/avraam311/image-processor/internal/repository/images"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) GetProcessedImage(c *ginext.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("id is not proper unsigned integer or empty parameter")
		handlers.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("non-empty and proper id required"))
		return
	}
	id := uint(idInt)

	im, err := h.service.GetProcessedImage(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, images.ErrImageNotFound) {
			zlog.Logger.Warn().Err(err).Msg("image not found")
			handlers.Fail(c.Writer, http.StatusNotFound, fmt.Errorf("image not found"))
			return
		} else if errors.Is(err, images.ErrImageInProcess) {
			zlog.Logger.Warn().Err(err).Msg("image in process")
			handlers.Fail(c.Writer, http.StatusServiceUnavailable, fmt.Errorf("image in process"))
			return
		}

		zlog.Logger.Error().Err(err).Msg("failed to get image")
		handlers.Fail(c.Writer, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	handlers.OK(c.Writer, im)
}
