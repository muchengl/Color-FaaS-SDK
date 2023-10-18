package cfaas

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"log"
	"reflect"
)

const (
	codeFunctionPanic    = 512
	codeInitializerPanic = 513
	functionSuccess      = 200
	functionFail         = 500
	unknownError         = 590
)

type funcRuntime struct {
	fun    interface{}
	server server.Hertz
}

func Run(function interface{}) {
	// init env
	runt := funcRuntime{
		fun: function,
	}
	runt.launchServer()
}

func (r *funcRuntime) launchServer() {
	h := server.Default()
	r.server = *h

	h.GET("/heartbeat", r.heartbeat)
	h.GET("/invoke", r.invokeGet)
	h.POST("/invoke", r.invokePost)
	h.GET("/quit", r.exit)

	h.Spin()
}

func (r *funcRuntime) heartbeat(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusOK, "ok")
}

// invokeGet get interface, debug only
func (r *funcRuntime) invokeGet(ctx context.Context, c *app.RequestContext) {
	// get invoke context
	umsg := c.Query("umsg")
	log.Default().Printf("user msg: %s", umsg)

	// prepare for running
	funcContext := FaaSContext{}
	msg := FuncRequest(umsg)
	args := []reflect.Value{reflect.ValueOf(funcContext), reflect.ValueOf(msg)}

	r.invoke(ctx, c, args)
}

func (r *funcRuntime) invokePost(ctx context.Context, c *app.RequestContext) {
	req := funcInvokeRequest{}
	if err := c.Bind(&req); err != nil {
		c.String(consts.StatusOK, string(codeInitializerPanic))
	}

	reqByte, _ := json.Marshal(req)
	log.Default().Printf("user msg: %s", string(reqByte))

	funcContext := FaaSContext{}
	msg := FuncRequest(req.Msg)
	args := []reflect.Value{reflect.ValueOf(funcContext), reflect.ValueOf(msg)}

	r.invoke(ctx, c, args)
}

// todo for invoke, if user's code get a panic. sdk need to handle this, sdk can't panic.
func (r *funcRuntime) invoke(ctx context.Context, c *app.RequestContext, args []reflect.Value) {
	// run function and get res
	function := reflect.ValueOf(r.fun)
	result := function.Call(args)

	// return
	if result[1].Interface() != nil {
		c.String(consts.StatusOK, string(functionFail))
	}

	response := funcInvokeResponse{
		Status: functionSuccess,
		Msg:    string(result[0].Interface().(FuncResponse)),
	}
	responseByte, _ := json.Marshal(response)
	c.String(consts.StatusOK, string(responseByte))
}

func (r *funcRuntime) exit(ctx context.Context, c *app.RequestContext) {
	r.server.Shutdown(ctx)
}
