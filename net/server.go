package net

import (
	"fmt"

	"github.com/fc-testzone/apiserver/content"

	v1 "github.com/fc-testzone/apiserver/api/v1"
	"github.com/fc-testzone/apiserver/auth"
	"github.com/fc-testzone/apiserver/utils"
	jsoniter "github.com/json-iterator/go"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type WebServer struct {
	log  *utils.Log
	cfg  *utils.Configs
	auth *auth.Authorizer
	con  *content.Content
}

func NewWebServer(l *utils.Log, c *utils.Configs, a *auth.Authorizer, cn *content.Content) *WebServer {
	return &WebServer{
		log:  l,
		cfg:  c,
		auth: a,
		con:  cn,
	}
}

func (w *WebServer) IndexHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("application/json")

	ctx.WriteString("{\"server\": \"index\"}")

	w.log.Info("WEB", "Index handler was called")
}

func (w *WebServer) RegisterHandler(ctx *fasthttp.RequestCtx) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ctx.Response.Header.SetContentType("application/json")

	w.log.Info("WEB", "Register handler was called")

	var regData v1.RegisterData
	var err = json.Unmarshal(ctx.PostBody(), &regData)
	if err != nil {
		ctx.Response.Header.SetStatusCode(403)
		w.log.Error("WEB", "Fail to decode register request", err.Error())
		return
	}

	w.log.Info("WEB", "Registering new user \""+regData.Login+"\"")

	err = w.auth.CreateUser(regData.Login, regData.Passwd)
	if err != nil {
		ctx.Response.Header.SetStatusCode(403)
		w.log.Error("WEB", "Fail to register user", err.Error())
	}
}

func (w *WebServer) LoginHandler(ctx *fasthttp.RequestCtx) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ctx.Response.Header.SetContentType("application/json")

	w.log.Info("WEB", "Login handler was called")

	var regData v1.RegisterData
	var err = json.Unmarshal(ctx.PostBody(), &regData)
	if err != nil {
		ctx.Response.Header.SetStatusCode(403)
		w.log.Error("WEB", "Fail to decode login request", err.Error())
		return
	}

	w.log.Info("WEB", "Login user \""+regData.Login+"\"")

	var token, err2 = w.auth.CreateToken(regData.Login, regData.Passwd)
	if err2 != nil {
		ctx.Response.Header.SetStatusCode(401)
		w.log.Error("WEB", "Fail to login user", err2.Error())
		return
	}

	var c fasthttp.Cookie
	c.SetHTTPOnly(true)
	c.SetKey("token")
	c.SetValue(token)
	c.SetDomain(w.cfg.Settings().Domain)
	ctx.Response.Header.SetCookie(&c)
}

func (w *WebServer) ContentHandler(ctx *fasthttp.RequestCtx) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	ctx.Response.Header.SetContentType("application/json")

	w.log.Info("WEB", "Content handler was called")

	var token = string(ctx.Request.Header.Cookie("token"))
	var err = w.auth.CheckToken(token)
	if err != nil {
		ctx.Response.Header.SetStatusCode(401)
		w.log.Error("WEB", "Fail to login user", err.Error())
		return
	}

	var posts []content.Post
	var err2 = w.con.Posts(&posts)
	if err2 != nil {
		ctx.Response.Header.SetStatusCode(503)
		w.log.Error("WEB", "Posts not found", err2.Error())
		return
	}

	var out, _ = json.Marshal(posts)
	ctx.Write(out)
}

func (w *WebServer) Start(ip string, port int) error {
	var router = fasthttprouter.New()

	var makeReq = func(api string, req string) string {
		return "/api/" + api + req
	}

	router.GET(v1.IndexRequest, w.IndexHandler)
	router.GET(makeReq(v1.ApiVersion, v1.ContentRequest), w.ContentHandler)
	router.POST(makeReq(v1.ApiVersion, v1.RegisterRequest), w.RegisterHandler)
	router.POST(makeReq(v1.ApiVersion, v1.LoginRequest), w.LoginHandler)

	return fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), router.Handler)
}
