package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/StanislavSizhuk/go_to_do/internal/core/logger"
	core_http_request "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/request"
	core_http_response "github.com/StanislavSizhuk/go_to_do/internal/core/transport/http/response"
)

// DeleteTask godoc
// @Summary      Delete a task by ID
// @Tags         tasks
// @Param        id   path  int  true  "Task ID"
// @Success      204  "No Content"
// @Failure      400  {object}  core_http_response.ErrorResponse
// @Failure      404  {object}  core_http_response.ErrorResponse
// @Failure      500  {object}  core_http_response.ErrorResponse
// @Router       /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
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
	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)

		return
	}

	responseHandler.NoContentResponse()

}
