package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	core_logger "github.com/MlkProduction/golang-project/internal/core/logger"
	core_http_response "github.com/MlkProduction/golang-project/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RequestIDHeader = "X-Request-ID"
)

/*
Получается как надстройка над Middleware
1. RequestID() возвращает Middleware, Middleware является надстройкой над конкретным http.Handler
2. Надстройка заключается в том, что добавляем обработку запроса/ответа вокруг следующего Handler.
3. requestID получает ID хедера(заголовка) (Берём из входящего HTTP-запроса значение заголовка)
4. Если ответ пустая строка, то uuid.NewString() - создаёт новый уникальный UUID
5. Если всё ок, то видимо получаем "X-Request-ID" : requestID. r.Header.Set(RequestIDHeader, requestID) положи в заголовки запроса X-Request-ID со значением requestID.
6. То же самое, но уже для ответа
7. next.ServeHTTP(w, r) - Мы закончили свою работу и передаём запрос следующему обработчику:
*/

func RequestID() Middleware {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(RequestIDHeader, requestID)
			w.Header().Set(RequestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

/*
Получается как надстройка над Middleware, но с логами
1. Берём наш логгер и возвращаем Middleware
2. добавляем обработку запроса/ответа вокруг следующего Handler
3. requestID получает ID хедера(заголовка) (Берём из входящего HTTP-запроса значение заголовка)
4. requestID := r.Header.Get(RequestIDHeader) - получаем значение заголовка
5. l = видимо соединияем в лог ID заголовка и url
6. Берём context текущего запроса и кладём туда наш логгер под ключом "log".
7. r.WithContext(ctx) создаёт новый request, такой же как старый, но с нашим новым ctx. передаёт этот обновлённый request дальше.
*/

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			//	zap.String("http.Method", r.Method),
			)

			ctx := context.WithValue(r.Context(), "log", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

/*
Получается как надстройка над Middleware, но с логам
1. перехватить panic, чтобы она не уронила весь HTTP-сервер.
*/

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle http request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

/*

 */

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http.Method", r.Method),
				zap.Time("time", before.UTC()),
			)
			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCodeOrPanic()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}
