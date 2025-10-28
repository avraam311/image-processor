package images

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avraam311/image-processor/internal/api/handlers"
	"github.com/avraam311/image-processor/internal/models"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) UploadImage(c *ginext.Context) {
	var im models.Image

	if err := json.NewDecoder(c.Request.Body).Decode(&im); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to decode request body")
		handlers.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid request body: %s", err.Error()))
		return
	}

	if err := h.validator.Struct(im); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to validate request body")
		handlers.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("validation error: %s", err.Error()))
		return
	}

	id, err := h.service.UploadImage(c.Request.Context(), &im)
	if err != nil {
		zlog.Logger.Error().Err(err).Interface("image", im).Msg("failed to upload image")
		return
	}

	handlers.Created(c.Writer, id)
}
