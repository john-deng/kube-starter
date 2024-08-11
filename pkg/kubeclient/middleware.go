package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// PostHandler is the middleware post handler
func (m *middleware) PostHandler(_ struct {
	at.MiddlewarePostHandler `value:"/" `
}, ctx context.Context) {
	responses := ctx.GetResponses()
	var baseResponseInfo model.BaseResponseInfo
	var statusCode int
	var response interface{}

	for _, resp := range responses {
		log.Debug(resp)
		if reflector.HasEmbeddedFieldType(resp, model.BaseResponseInfo{}) {
			response = resp
			respVal := reflector.GetFieldValue(resp, "BaseResponseInfo")
			if respVal.IsValid() {
				r := respVal.Interface()
				baseResponseInfo = r.(model.BaseResponseInfo)
			}
		}
		if resp != nil {
			switch resp.(type) {
			case error:
				log.Debug(resp)
				err := resp.(error)
				errStatusVal := reflector.GetFieldValue(err, "ErrStatus")
				if errStatusVal.IsValid() {
					esi := errStatusVal.Interface()
					errStatus := esi.(v1.Status)
					statusCode = int(errStatus.Code)
					log.Warn(errStatus)
				}
			}
		}
	}

	if statusCode != 0 {
		baseResponseInfo.SetCode(statusCode)
		err := reflector.SetFieldValue(response, "BaseResponseInfo", baseResponseInfo)
		log.Debugf("set BaseResponseInfo %v", err)
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}
