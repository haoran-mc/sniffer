package proxy

import (
	"log"
	"strconv"

	"github.com/haoran-mc/sniffer/output/db"
	"github.com/valyala/fasthttp"
)

func handleReplayRequest(ctx *fasthttp.RequestCtx) {
	client := &fasthttp.HostClient{
		Addr: "127.0.0.1:9523",
	}
	req := &ctx.Request

	var (
		host   = string(req.Host())
		method = string(req.Header.Method())
		url    = string(req.URI().Path())
		ip     = "127.0.0.1" // TODO
	)

	req.SetHost("127.0.0.1:9523")
	req.URI().SetScheme("http")

	var resp fasthttp.Response
	if err := client.Do(req, &resp); err != nil {
		ctx.Error("Proxy error", fasthttp.StatusBadGateway)
		return
	}

	db.WriteTraffic(host, method, url, ip,
		strconv.Itoa(resp.StatusCode()),
		string(req.Body()),
		string(resp.Body()))

	ctx.SetStatusCode(resp.StatusCode())
	ctx.Response.SetBodyRaw(resp.Body())
	resp.Header.CopyTo(&ctx.Response.Header)
}

func StartServer() {
	if err := fasthttp.ListenAndServe("127.0.0.1:9522", handleReplayRequest); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
