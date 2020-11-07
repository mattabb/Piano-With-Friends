package handlers


import (
	"net/http"
)

// RenderHome renders the home page
func RenderHome(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "UI/public/index.html")
}

