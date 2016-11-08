package main

import(
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
	"expvar"
)

func main(){

	//normalServer()
	//fsServer();
}

func normalServer(){
	fmt.Print("normal server")
	if error:=fasthttp.ListenAndServe(":65530",func(ctx *fasthttp.RequestCtx){
		ctx.SetContentType("text/plain; charset=utf8")

		// Set arbitrary headers
		ctx.Response.Header.Set("X-My-Header", "my-header-value")

		// Set cookies
		var c fasthttp.Cookie
		c.SetKey("cookie-name")
		c.SetValue("cookie-value")
		ctx.Response.Header.SetCookie(&c)
	});error!=nil{
		fmt.Println(error.Error())
	}
}

func fsServer(){
	fmt.Print("fileSystem server test")
	fs := &fasthttp.FS{
		Root:               "/Users/bdiao",
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		Compress:           false,
		AcceptByteRange:    false,
	}
	fsHandler:=fs.NewRequestHandler();
	if error := fasthttp.ListenAndServe(":65531", func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/states":
			expvarhandler.ExpvarHandler(ctx)
		default:
			fsHandler(ctx)
			updateFSCounters(ctx)
		}
	}); error != nil {
		fmt.Println(error.Error())
	}

}

func updateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)
