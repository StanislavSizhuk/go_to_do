package tasks_transport_http

import (
	"net/http"

	"github.com/StanislavSizhuk/go_to_do/internal/core/domain"
	core_logger "github.com/StanislavSizhuk/go_to_do/internal/core/logger"
	core_http_request "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/request"
	core_http_response "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"Buy groceries"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Milk, eggs, bread"`
	AuthorUserID int     `json:"author_user_id" validate:"required" example:"1"`
}

type CreateTaskResponse TaskDTOResponse

// CreateTask godoc
// @Summary      Create a new task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request  body      CreateTaskRequest  true  "Task payload"
// @Success      201      {object}  CreateTaskResponse
// @Failure      400      {object}  core_http_response.ErrorResponse
// @Failure      500      {object}  core_http_response.ErrorResponse
// @Router       /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}
	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)

		return
	}
	response := CreateTaskResponse( taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

