package core

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// This assumes we are using a web action with raw=true,
type MainRequestArgs struct {
	Path                  string            `json:"__ow_path"`            // The url path for the caller
	HTTPMethod            string            `json:"__ow_method"`          // The http method of the function
	Headers               map[string]string `json:"__ow_headers"`         // the request headers.
	QueryStringParameters string            `json:"__ow_query"`           // the query parameters from the request as an unparsed string separated by "&"
	Body                  string            `json:"__ow_body"`            // the request body entity, as a base64 encoded string when content is binary or JSON object/array, or plain string otherwise.
	IsBase64Encoded       bool              `json:"__ow_isBase64Encoded"` // set to true when the body content is binary
}

type MainRequest struct{}

// Receives the web function arguments, converts it to an http
func (r *MainRequest) MainArgsToHTTPRequest(req MainRequestArgs) (*http.Request, error) {
	return toRequest(r, req)
}

func toRequest(r *MainRequest, req MainRequestArgs) (*http.Request, error) {
	decodedBody := []byte(req.Body)

	if req.IsBase64Encoded {
		base64Body, err := base64.StdEncoding.DecodeString(req.Body)
		if err != nil {
			return nil, err
		}
		decodedBody = base64Body
	}

	path := req.Path
	if path == "" {
		path = "/"
	}
	path = path + "?" + req.QueryStringParameters

	method := req.HTTPMethod
	if method == "" {
		method = "POST"
	}

	httpRequest, err := http.NewRequest(
		strings.ToUpper(req.HTTPMethod),
		path,
		bytes.NewReader(decodedBody),
	)

	if err != nil {
		fmt.Printf("Could not convert request %s:%s to http.Request\n", req.HTTPMethod, req.Path)
		log.Println(err)
		return nil, err
	}

	for h := range req.Headers {
		httpRequest.Header.Add(h, req.Headers[h])
	}

	httpRequest.RequestURI = httpRequest.URL.RequestURI()

	return httpRequest, nil
}
