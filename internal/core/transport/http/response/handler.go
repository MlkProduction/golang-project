package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_error "github.com/MlkProduction/golang-project/internal/core/errors"
	core_logger "github.com/MlkProduction/golang-project/internal/core/logger"
	"go.uber.org/zap"
)

/*
HTTP обработчиков ответов
Структура
1.  Берём ссылку на структуру Logger из core_logger (работаем с оригиналом)
2. Интерфейс, через который этот объект сможет записывать HTTP-ответ клиенту.
*/
type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

/*
Это видимо конструктор для структуры
1.  получает готовый логгер → создаёт HTTPResponseHandler → кладёт логгер внутрь → возвращает указатель на новый обработчик.
*/
func NewHTTPResponseHandler(log *core_logger.Logger, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.rw.WriteHeader(statusCode)
	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)

}

/*

Видимо если паника, то как то её обрабатываем
1.  статус код = ошибка 500
2. присваем p ошибку непонятной паники (непонятной в связи с чем)
3. в лог какой-то месседж и ошибку из zap
4. в ответ на запрос записываем в хедер статус код (500)
5. присваем ответу мапу с ключом - меседж, значение - ошибка
6. если ошибка (не знаю что делает енкодер и енкод)в, теории наверно он кодирует ответ как-то,
json.NewEncoder(h.rw) создаёт JSON-кодировщик, который умеет писать JSON прямо в h.rw.
Encode превращает Go-значение в JSON и записывает его в HTTP-ответ.
7. И последний if err := ... проверяет уже ошибку при записи этого ответа.
Если даже ответ с ошибкой не удалось отправить — тогда это дополнительно пишется в лог.

*/

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {

	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_error.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_error.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_error.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)

}

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {

	h.rw.WriteHeader(statusCode)

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponse(response, statusCode)

}
