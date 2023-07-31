package ascii

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	Ping(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "PONG!", rr.Body.String())
}

func TestImageToAsciiSupportedFileTypes(t *testing.T) {
	type testcase struct {
		path     string
		filename string
	}
	cases := map[string]testcase{
		"jpg": {path: "files/cat2.jpg", filename: "cat2.jpg"},
		"png": {path: "files/cat2.png", filename: "cat2.png"},
	}

	for _, c := range cases {
		file, _ := os.Open(c.path)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("media", c.filename)
		io.Copy(part, file)
		writer.Close()
		req := httptest.NewRequest("POST", "/image-to-ascii", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		ImageToAscii(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)
	}
}

func TestImageToAsciiInvalidFile(t *testing.T) {
	file, _ := os.Open("files/dummy.pdf")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("media", "dummy.pdf")
	io.Copy(part, file)
	writer.Close()
	req := httptest.NewRequest("POST", "/image-to-ascii", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	ImageToAscii(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
