package kubeclient

import (
	"errors"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
)

var (
	ErrNilKubeClient = errors.New("kube client is nil, please check if API Server is available")
)

type middleware struct {
	at.Middleware
}

func newMiddleware() *middleware {
	return &middleware{}
}

func init() {
	app.Register(newMiddleware)
}

// CheckKubeClient is the middleware handler,it supports dependency injection, method annotation
// middleware handler can be annotated to specific purpose or general purpose
func (m *middleware) CheckKubeClient(_ struct {
	at.MiddlewareHandler `value:"/" `
},
	//cli *Client,
	ctx context.Context) (err error) {

	//if cli.Client == nil {
	//	err = fmt.Errorf("[middleware] kube client is nil")
	//	log.Warn(err)
	//	ctx.StatusCode(500)
	//	ctx.ResponseBody(err.Error(), nil)
	//	return
	//}
	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}
