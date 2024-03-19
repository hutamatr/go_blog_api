package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hutamatr/go-blog-api/exception"
	"github.com/hutamatr/go-blog-api/helpers"
	"github.com/hutamatr/go-blog-api/model/web"
	servicesU "github.com/hutamatr/go-blog-api/services/user"
	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	service servicesU.UserService
}

var accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
var refreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")

func NewUserController(userService servicesU.UserService) UserController {
	return &UserControllerImpl{
		service: userService,
	}
}

func (controller *UserControllerImpl) CreateUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var newUserRequest web.UserCreateRequest

	helpers.DecodeJSONFromRequest(request, &newUserRequest)

	newUser, accessToken, refreshToken := controller.service.SignUp(request.Context(), newUserRequest)

	cookie := http.Cookie{}
	cookie.Name = "rt"
	cookie.Value = refreshToken
	cookie.MaxAge = 7 * 24 * 60 * 60 * 1000
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	http.SetCookie(writer, &cookie)

	userResponse := web.ResponseJSON{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data: map[string]interface{}{
			"access_token": accessToken,
			"user":         newUser,
		},
	}

	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) SignInUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var signInRequest web.UserLoginRequest

	helpers.DecodeJSONFromRequest(request, &signInRequest)

	signInUser, accessToken, refreshToken := controller.service.SignIn(request.Context(), signInRequest)

	cookie := http.Cookie{}
	cookie.Name = "rt"
	cookie.Value = refreshToken
	cookie.MaxAge = 7 * 24 * 60 * 60 * 1000
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	http.SetCookie(writer, &cookie)

	userResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"access_token": accessToken,
			"user":         signInUser,
		},
	}

	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) SignOutUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	cookie := http.Cookie{}
	cookie.Name = "rt"
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	http.SetCookie(writer, &cookie)

	userResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
	}

	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) FindAllUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	users := controller.service.FindAll(request.Context())

	userResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   users,
	}

	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) FindByIdUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	user := controller.service.FindById(request.Context(), userId)

	userResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   user,
	}

	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) UpdateUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	var updateUserRequest web.UserUpdateRequest

	updateUserRequest.Id = userId

	helpers.DecodeJSONFromRequest(request, &updateUserRequest)

	updatedUser := controller.service.Update(request.Context(), updateUserRequest)

	userResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "UPDATED",
		Data:   updatedUser,
	}

	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) DeleteUser(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id := params.ByName("userId")
	userId, err := strconv.Atoi(id)
	helpers.PanicError(err)

	controller.service.Delete(request.Context(), userId)

	userResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "DELETED",
	}

	helpers.EncodeJSONFromResponse(writer, userResponse)
}

func (controller *UserControllerImpl) GetRefreshToken(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	refreshToken, err := request.Cookie("rt")
	if err != nil {
		panic(exception.NewUnauthorizedError(err.Error()))
	}

	claims, err := helpers.VerifyToken(refreshToken.Value, refreshTokenSecret)

	if err != nil {
		panic(exception.NewUnauthorizedError(err.Error()))
	}

	username := claims["sub"].(string)

	newAccessToken, err := helpers.GenerateToken(username, 5*time.Minute, accessTokenSecret)

	helpers.PanicError(err)

	userResponse := web.ResponseJSON{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"access_token": newAccessToken,
		},
	}

	helpers.EncodeJSONFromResponse(writer, userResponse)
}