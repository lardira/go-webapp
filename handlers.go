package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lardira/go-webapp/model"
	"github.com/lardira/go-webapp/utils"
	http_errors "github.com/lardira/go-webapp/utils/errors"
)

type UserHandler struct {
}

type VariantHandler struct {
}

type AuthHandler struct {
}

type TestHandler struct {
}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)

	id, err := strconv.Atoi(head)
	if len(head) > 0 && err != nil {
		http_errors.BadRequest(w, err)
		return
	}

	switch r.Method {

	case http.MethodPost:
		if id != 0 {
			http_errors.BadRequest(w, errors.New("id provided"))
			return
		}

		var request model.UserRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = request.Validate()
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		user, err := model.CreateUser(GlobalConnectionPool, request.Login, request.Password)
		if err != nil {
			http_errors.Error(w, errors.New("could not create user"))
			return
		}

		w.WriteHeader(http.StatusCreated)

		response, err := json.Marshal(user)
		if err != nil {
			http_errors.Error(w, errors.New("error occured"))
			return
		}

		fmt.Fprint(w, string(response))

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func returAvailableTasks(w http.ResponseWriter, id int) {
	tasks, err := model.GetAllTasksByVariantId(GlobalConnectionPool, int64(id))
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}

func returnAllVariants(w http.ResponseWriter) {
	variants, err := model.GetAllVariants(GlobalConnectionPool)
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	response, err := json.Marshal(variants)
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}

func returnTaskById(w http.ResponseWriter, taskId int, varId int) {
	task, err := model.GetTask(GlobalConnectionPool, int64(taskId), int64(varId))
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}

func (uh *VariantHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)

	varId, err := strconv.Atoi(head)
	if len(head) > 0 && err != nil {
		http_errors.BadRequest(w, err)
		return
	}

	variantIdPresented := len(head) > 0

	switch r.Method {

	case http.MethodGet:
		if variantIdPresented {
			head, r.URL.Path = utils.ShiftPath(r.URL.Path)

			//return variant's task data if id of task
			taskIdPresented := len(head) > 0
			if taskIdPresented {
				taskId, err := strconv.Atoi(head)
				if len(head) > 0 && err != nil {
					http_errors.BadRequest(w, err)
					return
				}
				returnTaskById(w, taskId, varId)

			} else {
				returAvailableTasks(w, varId)
			}

		} else {
			returnAllVariants(w)
		}

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)

	id, err := strconv.Atoi(head)
	if len(head) > 0 && err != nil {
		http_errors.BadRequest(w, err)
		return
	}

	switch r.Method {

	case http.MethodPost:
		type AuthResponse struct {
			Key string `json:"key"`
		}

		if id != 0 {
			http_errors.BadRequest(w, errors.New("id provided"))
			return
		}

		var request model.UserRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = request.Validate()
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = model.Authorize(
			GlobalConnectionPool,
			request.Login,
			request.Password,
		)

		if err != nil {
			http_errors.NotFound(w, errors.New("no user with such credentials"))
			return
		}

		response, err := json.Marshal(
			AuthResponse{
				Key: fmt.Sprintf("Basic %s:%s", request.Login, request.Password),
			},
		)

		if err != nil {
			http_errors.Error(w, errors.New("could not respond"))
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))

	case http.MethodPut:

		var request model.UserRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = request.Validate()
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = model.LogOutUser(
			GlobalConnectionPool,
			request.Login,
			request.Password,
		)
		if err != nil {
			http_errors.NotFound(w, errors.New("no user with such credentials"))
			return
		}

		w.WriteHeader(http.StatusOK)

		response, _ := json.Marshal(DefaultResponse{
			Message: "ok",
		})
		fmt.Fprint(w, string(response))

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func (th *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)

	varId, err := strconv.Atoi(head)
	if len(head) > 0 && err != nil {
		http_errors.BadRequest(w, err)
		return
	}

	idPresented := len(head) > 0

	switch r.Method {

	case http.MethodPost:
		if idPresented {
			login, password, ok := utils.ParseAuthHeader(w, r, AUTH_TYPE, AUTH_HEADER)
			if !ok {
				http_errors.BadRequest(w, err)
				return
			}

			user, err := model.GetUserByLoginAndPasword(GlobalConnectionPool, login, password)
			if err != nil {
				http_errors.BadRequest(w, err)
				return
			}

			test, err := model.CreateTest(GlobalConnectionPool, user.Id, int64(varId))
			if err != nil {
				http_errors.Error(w, err)
				return
			}

			response, err := json.Marshal(model.TestResponse{Id: test.Id})
			if err != nil {
				http_errors.Error(w, err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, string(response))
		}

	case http.MethodPut:
		if idPresented {
			var request model.TestAnswerRequest

			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				http_errors.BadRequest(w, err)
				return
			}

			err = model.AddTestAnswer(
				GlobalConnectionPool,
				request.TestId,
				request.Answer,
			)

			if err != nil {
				http_errors.Error(w, err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			response, _ := json.Marshal(DefaultResponse{Message: "ok"})
			fmt.Fprint(w, string(response))
		}

	case http.MethodGet:
		if idPresented {
			head, r.URL.Path = utils.ShiftPath(r.URL.Path)

			testId, err := strconv.Atoi(head)
			if len(head) > 0 && err != nil {
				http_errors.BadRequest(w, err)
				return
			}

			testResult, err := model.GetTestResult(
				GlobalConnectionPool,
				int64(testId),
				int64(varId),
			)

			if err != nil {
				http_errors.Error(w, err)
				return
			}

			response, err := json.Marshal(testResult)
			if err != nil {
				http_errors.Error(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(response))
		}

	default:
		http_errors.MethodNotAllowed(w)
	}
}
