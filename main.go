package main

import (
	"context"
	"fmt"
	"github.com/raythx98/gohelpme/builder/httpclient"
	"github.com/raythx98/gohelpme/builder/httprequest"
	"github.com/raythx98/gohelpme/middleware"
	"github.com/raythx98/gohelpme/tool/logger"
	"io"
	"log"
	"net/http"
)

func main() {
	// converting our handler function to handler
	// type to make use of our middleware
	mux := http.NewServeMux()

	l := logger.NewDefault()
	opt := middleware.LogOptions{
		RequestRedact:  []string{"body.variables.password"},
		ResponseRedact: []string{"body.data.login.token"},
	}

	finalHandler := http.HandlerFunc(handler)
	mux.Handle("/", middleware.JsonResponse(middleware.AddRequestId(middleware.Log(l, opt)(finalHandler))))
	mux.Handle("/test", middleware.Log(l)(finalHandler))

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}

func handler(w http.ResponseWriter, r *http.Request) {

	requestBody, _ := io.ReadAll(r.Body)
	fmt.Println("test receive request body", string(requestBody))

	callHttp(r.Context(), &httprequest.Builder{})

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("{\"response\": \"body\"}"))
}

func callHttp(ctx context.Context, builder *httprequest.Builder) {
	httpClient := httpclient.New(logger.NewDefault())

	//req, err := builder.New(ctx, httprequest.Get, "https://echo.free.beeceptor.com").WithBody("").Build()
	req, err := builder.New(ctx, httprequest.Get, "https://echo.free.beeceptor.com").Build()
	if err == nil {
		resp, _ := httpClient.Do(req)
		if resp != nil {
			respBody, _ := io.ReadAll(resp.Body)
			fmt.Println("test receive response body", string(respBody))
		}
	}
}
