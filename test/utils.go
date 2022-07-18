package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func ServeTestHTTPRequest(method, url string, body io.Reader, router *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, "/api/v0"+url, body)
	fmt.Println(req.URL)
	router.ServeHTTP(w, req)

	return w
}
