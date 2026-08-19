package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/StanislavSizhuk/go_to_do/internal/core/logger"
	core_http_request "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/request"
	core_http_response "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask godoc
// @Summary      Get a task by ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  GetTaskResponse
// @Failure      400  {object}  core_http_response.ErrorResponse
// @Failure      404  {object}  core_http_response.ErrorResponse
// @Failure      500  {object}  core_http_response.ErrorResponse
// @Router       /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)

		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}
