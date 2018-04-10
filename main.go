package main

import (
	"encoding/base64"
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
	"strings"
	"time"
)

func reverseProxyHandler(ctx *fasthttp.RequestCtx) {

	if defaultRouter(ctx) {
		return
	}
	path := strings.Split(string(ctx.Path()), "/")
	uri, err := decodePath(path)
	if err != nil {
		ctx.Logger().Printf("error when decode path: %s", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	ctx.Request.SetRequestURI(uri.Path)
	req := &ctx.Request
	resp := &ctx.Response
	prepareRequest(req, uri)
	proxyClient := &fasthttp.HostClient{
		Addr:                uri.Host,
		ReadTimeout:         60 * time.Second,
		WriteTimeout:        60 * time.Second,
		MaxResponseBodySize: 5242800, // 5mb
	}
	if err := proxyClient.Do(req, resp); err != nil {
		ctx.Logger().Printf("error when proxy the request: %s", err)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
	postprocessResponse(resp)
}

func decodePath(path []string) (*url.URL, error) {
	host, err := base64.RawURLEncoding.DecodeString(path[1])
	if err != nil {
		return nil, err
	}
	path[1] = string(host)
	s := len(path)
	uri, err := url.Parse(strings.Join(path[1:s], "/"))
	if err != nil {
		return nil, err
	}
	if uri.Scheme != "http" && uri.Scheme != "https" {
		//noinspection GoErrorStringFormat
		return nil, errors.New("Scheme must be http or https")
	}
	return uri, err
}

func defaultRouter(ctx *fasthttp.RequestCtx) bool {
	if string(ctx.Path()) == "/healthz" || string(ctx.Path()) == "/" {
		ctx.SetContentType("application/json; charset=utf8")
		ctx.SetBodyString("{\"alive\": true}")
		ctx.SetStatusCode(fasthttp.StatusOK)
		return true
	} else if string(ctx.Path()) == "/favicon.ico" {
		ctx.SetContentType("text/plain; charset=utf8")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return true
	}

	return false
}

func prepareRequest(req *fasthttp.Request, uri *url.URL) {
	// do not proxy "Connection" header.
	req.Header.Del("Connection")
	req.Header.Set("Host", uri.Host)
	req.Header.Set("Referer", uri.String())
}

func postprocessResponse(resp *fasthttp.Response) {
	// do not proxy "Connection" header
	resp.Header.Del("Connection")
}

func main() {
	if err := fasthttp.ListenAndServe(":8001", reverseProxyHandler); err != nil {
		log.Fatalf("error in fasthttp server: %s", err)
	}
	log.Printf("See stats at http://127.0.0.1:8001/healthz")
}
