package images

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/avraam311/image-processor/internal/api/handlers"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) DeleteImage(c *ginext.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		zlog.Logger.Warn().Err(err).Msg("id is not proper unsigned integer or empty parameter")
		handlers.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("non-empty and proper id required"))
		return
	}
	id := uint(idInt)

	if err := h.service.DeleteImage(c.Request.Context(), id); err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to delete image")
		handlers.Fail(c.Writer, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	handlers.OK(c.Writer, "image deleted")
}
