package kubeclient

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"unsafe"
)

type middleware struct {
	at.Middleware
	at.UseMiddleware
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
}, ctx context.Context) (response model.ResponseInfo, err error) {
	responses := ctx.GetResponses()
	//var baseResponseInfo model.BaseResponseInfo
	var status v1.Status

	for _, resp := range responses {
		log.Debug(resp)
		if reflector.HasEmbeddedFieldType(resp, model.BaseResponseInfo{}) {
			response = resp.(model.ResponseInfo)
		}
		if resp != nil {
			switch resp.(type) {
			case error:
				log.Debug(resp)
				err = resp.(error)
				errStatusVal := reflector.GetFieldValue(err, "ErrStatus")
				if errStatusVal.IsValid() {
					esi := errStatusVal.Interface()
					status = esi.(v1.Status)
				} else {
					status, _ = m.getErrorStatus(err)
				}
				log.Warn(status)
			}
		}
	}

	if err != nil {
		var msg string
		if status.Code != 0 && status.Code != 200 {
			msg = fmt.Sprintf("%v, %v, %v", status.Status, status.Message, status.Reason)
		}
		if response == nil {
			response = new(model.BaseResponseInfo)
		}
		response.SetCode(int(status.Code))
		response.SetMessage(msg)
		log.Debugf("set response: %v, err: %v", response, err)
		return
	}

	// call ctx.Next() if you want to continue, otherwise do not call it
	ctx.Next()
	return
}

func (m *middleware) getErrorStatus(wrapError interface{}) (status v1.Status, ok bool) {
	// Use reflection to access the private 'err' field
	actualValue := reflect.ValueOf(wrapError).Elem()
	errField := actualValue.FieldByName("err")

	// Check if the 'err' field is valid and can be addressed
	if errField.IsValid() {
		// Since it's a private field, we need to use .Interface() on the field directly
		errValue := reflect.NewAt(errField.Type(), unsafe.Pointer(errField.UnsafeAddr())).Elem()

		// Convert the reflect.Value to an interface{} and then type assert it
		errMap, converted := errValue.Interface().(*apiutil.ErrResourceDiscoveryFailed)
		if !converted {
			return
		}

		// Iterate through the map and use reflection to get the ErrStatus
		for _, err := range *errMap {
			// Get the reflect value of the error
			value := reflect.ValueOf(err).Elem()

			// Access the ErrStatus field
			errStatusField := value.FieldByName("ErrStatus")
			if errStatusField.IsValid() {
				// Convert the reflect value to the appropriate type
				status, ok = errStatusField.Interface().(v1.Status)
			}
		}
	}
	return
}
