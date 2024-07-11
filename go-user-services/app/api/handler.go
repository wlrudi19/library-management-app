package api

import (
	"encoding/json"
	"net/http"

	"github.com/wlrudi19/library-management-app/go-user-services/app/model"
	"github.com/wlrudi19/library-management-app/go-user-services/app/service"
	"github.com/wlrudi19/library-management-app/go-user-services/utils/response"
)

type UserHandler interface {
	LoginUser(writer http.ResponseWriter, req *http.Request)
}

type userhandler struct {
	UserLogic service.UserLogic
}

func NewUserHandler(userLogic service.UserLogic) UserHandler {
	return &userhandler{
		UserLogic: userLogic,
	}
}

func (h *userhandler) LoginUser(writer http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var jsonReq model.UserRequest

	err := json.NewDecoder(req.Body).Decode(&jsonReq)
	if err != nil {
		resp := response.CustomBuilder(http.StatusBadRequest, false, jsonReq, "Request invalid. Format JSON not match")
		resp.Send(writer)
		return
	}

	var loginToken = model.LoginResponse{}
	loginToken, err = h.UserLogic.LoginUser(ctx, jsonReq.Email, jsonReq.Password)
	if err != nil {
		resp := response.CustomBuilder(http.StatusUnauthorized, false, jsonReq, "You are not authorized to access this resource")
		resp.Send(writer)
		return
	}

	resp := response.CustomBuilder(http.StatusOK, true, loginToken, "You are success login")
	resp.Send(writer)
	return
}
