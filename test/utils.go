package test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func ServeTestHTTPRequest(method, url string, body io.Reader, router *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	router.ServeHTTP(w, req)

	return w
}
