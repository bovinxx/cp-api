package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/bovinxx/code-processor/api/pkg/metrics"
	"github.com/bovinxx/code-processor/api/pkg/services"
)

type Handler struct {
	Services *services.Services
	Request  *Request
	Response *Response
}

var (
	metric = metrics.NewMetrics()
)

func NewHandler(service *services.Services) *Handler {
	return &Handler{
		Services: service,
		Request:  &Request{},
		Response: &Response{},
	}
}

// @Summary Create task
// @Description Create task and return taskID
// @Success 201 {string} string "taskID"
// @Failure 400 {string} string "invalid json"
// @Failure 401 {string} string "unauthorized user"
// @Router /task [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	defer metric.Record(time.Now(), "create")
	if err := h.userIdentity(r); err != nil {
		h.Response.ErrorMessage(w, http.StatusUnauthorized, err)
		return
	}

	task, err := h.Request.ParseTask(r)
	if err != nil {
		h.Response.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.Services.Task.CreateTask(*task)
	if err != nil {
		h.Response.ErrorMessage(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}
	h.Response.JsonResponse(w, http.StatusCreated, map[string]interface{}{"id": id})
}

// @Summary Get status task
// @Description Return status task by id
// @Param id path string true "taskID"
// @Success 200 {string} string "status"
// @Failure 400 {string} string "there's no such task"
// @Failure 401 {string} string "unauthorized user"
// @Router /status/{id} [get]
func (h *Handler) GetStatusTask(w http.ResponseWriter, r *http.Request) {
	defer metric.Record(time.Now(), "status")
	if err := h.userIdentity(r); err != nil {
		h.Response.ErrorMessage(w, http.StatusUnauthorized, err)
		return
	}

	id := h.Request.ParseTaskID(r)
	if id == "" {
		h.Response.ErrorMessage(w, http.StatusBadRequest, errors.New("empty id"))
		return
	}
	status, err := h.Services.Task.GetStatus(id)
	if err != nil {
		h.Response.ErrorMessage(w, http.StatusNotFound, err)
		return
	}
	h.Response.JsonResponse(w, http.StatusOK, map[string]interface{}{"status": status})
}

// @Summary Get result task
// @Description Return result task by id
// @Param id path string true "taskID"
// @Success 200 {string} string "Result"
// @Failure 400 {string} string "there's not such task"
// @Failure 401 {string} string "unauthorized user"
// @Router /result/{id} [get]
func (h *Handler) GetResultTask(w http.ResponseWriter, r *http.Request) {
	defer metric.Record(time.Now(), "result")

	if err := h.userIdentity(r); err != nil {
		h.Response.ErrorMessage(w, http.StatusUnauthorized, err)
		return
	}

	id := h.Request.ParseTaskID(r)
	if id == "" {
		h.Response.ErrorMessage(w, http.StatusBadRequest, errors.New("empty id"))
		return
	}
	result, err := h.Services.Task.GetResult(id)
	if err != nil {
		h.Response.ErrorMessage(w, http.StatusNotFound, err)
		return
	}
	h.Response.JsonResponse(w, http.StatusOK, map[string]interface{}{"result": result})
}

// @Summary Register user
// @Description Register user and return usedID
// @Accept  json
// @Produce  json
// @param user body model.User true "Register user"
// @Success 201 {string} string "usedID"
// @Failure 400 {string} string "bad json"
// @Router /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer metric.Record(time.Now(), "register")

	user, err := h.Request.ParseUser(r)
	if err != nil {
		h.Response.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	id, err := h.Services.Auth.CreateUser(user.Username, user.Password)
	if err != nil {
		h.Response.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	h.Response.JsonResponse(w, http.StatusCreated, map[string]string{"id": id})
}

// @Summary Login user
// @Description Login user and return sessionID
// @Accept  json
// @Produce  json
// @param user body model.User true "Login user"
// @Success 201 {string} string "sessionID"
// @Failure 400 {string} string "bad json"
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer metric.Record(time.Now(), "login")

	user, err := h.Request.ParseUser(r)
	if err != nil {
		h.Response.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	sessionID, err := h.Services.Auth.LoginUser(user.Username, user.Password)
	if err != nil {
		h.Response.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	h.Response.JsonResponse(w, http.StatusOK, map[string]string{"token": sessionID})
}

func (h *Handler) userIdentity(r *http.Request) error {
	authorizationParts := strings.Split(h.Request.ParseAuthToken(r), " ")
	if len(authorizationParts) < 2 || authorizationParts[0] != "Bearer" {
		return errors.New("invalid token")
	}
	token := authorizationParts[1]
	if token == "" {
		return errors.New("token is empty")
	}

	ok, err := h.Services.Auth.Authorization(token)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("incorrect token")
	}
	return nil
}
