package users_transport_http

import (
	"net/http"

	core_logger "github.com/MlkProduction/golang-project/internal/core/logger"
	core_http_response "github.com/MlkProduction/golang-project/internal/core/transport/http/response"
	core_http_utils "github.com/MlkProduction/golang-project/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "Don't delete user id path value")
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(err, "don't can delete user from this id")
		return
	}

	responseHandler.NoContentResponse()

}
