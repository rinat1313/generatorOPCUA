package controller

import (
	"generatorOPCUA/internal/server/httpServer"
	"log"
	"net/http"
)

type Controller struct {
}

func (c *Controller) StartWork() {
	handler := httpServer.NewHttpHandler()
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Println(err)
	}
}
