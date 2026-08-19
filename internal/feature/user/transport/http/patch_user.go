package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/StanislavSizhuk/go_to_do/internal/core/domain"
	core_logger "github.com/StanislavSizhuk/go_to_do/internal/core/logger"
	core_http_request "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/request"
	core_http_response "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/response"
	core_http_types "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"John Doe"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+380501234567"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` cannot be null: ")
		}
		fullNameLength := len([]rune(*r.FullName.Value))
		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 characters: ")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 characters: ")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must start with a '+' sign: ")
			}
		}
	}

	return nil

}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
// @Summary      Partially update a user
// @Description  Updates one or more fields of a user. Only fields present in the request body are modified.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int               true  "User ID"
// @Param        request  body      PatchUserRequest  true  "Fields to update"
// @Success      200      {object}  PatchUserResponse
// @Failure      400      {object}  core_http_response.ErrorResponse
// @Failure      404      {object}  core_http_response.ErrorResponse
// @Failure      500      {object}  core_http_response.ErrorResponse
// @Router       /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user ID from path",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,

			"failed to decode and validate HTTP request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain() ,
		request.PhoneNumber.ToDomain(),
	)
}
