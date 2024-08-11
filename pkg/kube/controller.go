package kube

import "github.com/hidevopsio/hiboot/pkg/at"

// Controller annotation in kube-starter is used to identify the controller
//
//	type MyController struct {
//	  kube.Controller
//	  ...
//	}
//

type Controller struct {
	at.Annotation `json:"-"`

	at.BaseAnnotation
}
