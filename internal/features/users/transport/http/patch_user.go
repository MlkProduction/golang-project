package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MlkProduction/golang-project/internal/core/domain"
	core_logger "github.com/MlkProduction/golang-project/internal/core/logger"
	core_http_request "github.com/MlkProduction/golang-project/internal/core/transport/http/request"
	core_http_response "github.com/MlkProduction/golang-project/internal/core/transport/http/response"
	core_http_types "github.com/MlkProduction/golang-project/internal/core/transport/http/types"
	core_http_utils "github.com/MlkProduction/golang-project/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("Full Name can't be NULL")
		}
		fullNameLen := len([]rune(*r.FullName.Value))

		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("Full Name can't be > 3 or < 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("Phone number can't be < 10 or > 15 symbols")
			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("Phone number must starts with '+'")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	var request PatchUserRequest

	userID, err := core_http_utils.GetIntPathValue(r, "id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed get user id")
		return
	}

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed validate and decode http.request")
		return
	}
	userPatch := UserPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func UserPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
