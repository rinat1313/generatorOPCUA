package httpServer

import (
	"log"
	"net/http"
)

type httpFunc func(http.ResponseWriter, *http.Request)

type HandlerCommand struct {
	Name     string
	Method   string
	HttpFunc httpFunc
}

func checkValidMethod(method string) bool {
	switch method {
	case "GET":
		{
			return true
		}
	case "POST":
		{
			return true
		}
	case "PUT":
		{
			return true
		}
	case "DELETE":
		{
			return true
		}
	case "OPTIONS":
		{
			return true
		}
	default:
		{
			return false
		}
	}
}

func CreateHandlerCommand(path string, method string, httpFunc httpFunc) HandlerCommand {
	handler := HandlerCommand{}
	handler.Name = path
	if !checkValidMethod(method) {
		log.Fatalf("Invalid method: %s", method)
	}
	handler.Method = method
	handler.HttpFunc = httpFunc
	return handler
}
