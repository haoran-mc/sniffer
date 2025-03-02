package reassembler

import (
	"log"
	"strconv"

	"github.com/haoran-mc/sniffer/output/db"
	"github.com/valyala/fasthttp"
)

func proxyHandler(ctx *fasthttp.RequestCtx) {
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

	// 代理请求并接收响应
	var resp fasthttp.Response
	if err := client.Do(req, &resp); err != nil {
		// log.Println("Reverse proxy error:", err)
		ctx.Error("Proxy error", fasthttp.StatusBadGateway)
		return
	}

	db.WriteTraffic(host, method, url, ip,
		strconv.Itoa(resp.StatusCode()),
		string(req.Body()),
		string(resp.Body()))

	// 将响应复制回客户端
	ctx.SetStatusCode(resp.StatusCode())
	ctx.Response.SetBodyRaw(resp.Body())
	resp.Header.CopyTo(&ctx.Response.Header)
}

func StreamParserServerStart() {
	if err := fasthttp.ListenAndServe("127.0.0.1:9522", proxyHandler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
