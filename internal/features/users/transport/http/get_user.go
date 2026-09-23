package users_transport_http

import (
	"net/http"

	core_logger "github.com/MlkProduction/golang-project/internal/core/logger"
	core_http_response "github.com/MlkProduction/golang-project/internal/core/transport/http/response"
	core_http_utils "github.com/MlkProduction/golang-project/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "Don't get user id path value")
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)

	if err != nil {
		responseHandler.ErrorResponse(err, "don't can get user from this id")
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
