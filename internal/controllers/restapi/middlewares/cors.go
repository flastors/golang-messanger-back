package middlewares

import (
	"net/http"
)

// CORS — простой middleware для разработки.
// Устанавливает заголовки CORS и обрабатывает preflight запросы (OPTIONS).
// В production рекомендуется ограничить Access-Control-Allow-Origin конкретными доменами.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		// Если это preflight запрос — ответим 200 и не пропустим дальше
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
