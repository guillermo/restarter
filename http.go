package main

import (
	"fmt"
	"github.com/cespare/hutil/apachelog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Http struct {
	From int
	To   int
}

func (h *Http) ListenAndServe() error {

	urlString := fmt.Sprintf("http://localhost:%d/", h.To)
	target, err := url.Parse(urlString)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	addr := fmt.Sprintf(":%d", h.From)

	var handler http.Handler
	handler = proxy
	handler = apachelog.NewDefaultHandler(handler)

	err = http.ListenAndServe(addr, handler)
	if err != nil {
		panic(err)
	}

	return nil
}
