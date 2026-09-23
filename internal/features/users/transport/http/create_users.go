package users_transport_http

import (
	"net/http"

	"github.com/MlkProduction/golang-project/internal/core/domain"
	core_logger "github.com/MlkProduction/golang-project/internal/core/logger"
	core_http_request "github.com/MlkProduction/golang-project/internal/core/transport/http/request"
	core_http_response "github.com/MlkProduction/golang-project/internal/core/transport/http/response"
)

/*
CreateUser — это переводчик между HTTP-миром и твоим Go-приложением.

Он получает HTTP → достаёт JSON → превращает его в CreateUserRequest → передаёт данные сервису → получает результат → превращает результат в HTTP/JSON.

А конкретно что такое FromContext, почему Decode(&request), зачем здесь rw, и как должен участвовать usersService — можем разобрать по одному.

*/

//################################

/*
Структура для создания юзера, тоесть принимаем данные
1. Фуллнейм - строка - А тег json:"full_name" говорит encoding/json, что JSON-поле "full_name" соответствует Go-полю FullName.
2. Телефон - * (необязательно) - А тег json:"phone_number" говорит encoding/json, что JSON-поле "phone_number" соответствует Go-полю PhoneNumber.
*/
type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

/*
Структура - что будет возвращать при получении данных
1. ID
2. Версию
3. Полное имя
4. Телефон
*/
type CreateUserResponse UserDTOResponse

/*
HTTP handler, который принимает запрос.
1. ctx получает контект из запроса
2. Мы достаём из контекста логгер.
3. Записываем лог с левелом Debug
4. Мы берём JSON из HTTP-запроса и превращаем его в Go-структуру.
5. Если декодировать не получилось, мы выводим в консоль ошибку, но всёравно записываем в header статус код 200
6. Есл получилось, возвращаем ответ - 200
*/
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("Invoke CreateUser handler")
	var request CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate http request")

		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)

}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUnitialized(dto.FullName, dto.PhoneNumber)

}
