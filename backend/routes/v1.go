package routes

import (
	chatgpt "alle/client/chatGPT"
	"alle/handlers"

	"net/http"

	"github.com/gorilla/mux"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Router(chatGptClient chatgpt.ChatGPTClient) *mux.Router {
	router := mux.NewRouter()
	v1Router := router.PathPrefix("/api/v1").Subrouter()
	ChatRoutes(v1Router, chatGptClient)
	return router
}

func ChatRoutes(router *mux.Router, chatGptClient chatgpt.ChatGPTClient) {

	chatHandler := handlers.ChatHandler{
		ChatGPTClient: chatGptClient,
	}

	chatRouters := router.PathPrefix("/chat").Subrouter()
	chatRouters.HandleFunc("", chatHandler.AddChat).Methods("POST")
	chatRouters.HandleFunc("/all", chatHandler.AllChat).Methods("GET")

}
