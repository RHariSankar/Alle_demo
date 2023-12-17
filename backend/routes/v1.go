package routes

import (
	chatgpt "alle/client/chatGPT"
	"alle/controllers"
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
	ImageRoutes(v1Router)
	return router
}

func ChatRoutes(router *mux.Router, chatGptClient chatgpt.ChatGPTClient) {

	chatHandler := handlers.ChatHandler{
		ChatGPTClient:  chatGptClient,
		ChatController: controllers.GetChatControllerInstance(),
	}

	chatRouter := router.PathPrefix("/chat").Subrouter()
	chatRouter.HandleFunc("", chatHandler.AddChat).Methods("POST")
	chatRouter.HandleFunc("/all", chatHandler.AllChat).Methods("GET")

}

func ImageRoutes(router *mux.Router) {

	imageHandler := handlers.ImageHandler{
		ImageController: controllers.GetImageControllerInstance(),
		ChatController:  controllers.GetChatControllerInstance(),
	}

	imageRoute := router.PathPrefix("/image").Subrouter()
	imageRoute.HandleFunc("", imageHandler.AddImage).Methods("POST")
	imageRoute.HandleFunc("/{id}", imageHandler.GetImage).Methods("GET")

}
