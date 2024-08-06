package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/kube-starter/pkg/kubeconfig"
	"time"
)

// Properties the operator properties
type Properties struct {
	at.ConfigurationProperties `value:"kubeclient"`
	at.AutoWired

	// operator deployment namespace
	Namespace string `json:"namespace" default:"kube-system"`

	// use DefaultInCluster as default
	// Deprecated
	DefaultInCluster *bool `json:"defaultInCluster"`

	//OIDC Scope Impersonate
	OIDCScope string `json:"oidcScope"`

	QPS float32 `json:"qps"`

	// Maximum burst for throttle.
	// If it's zero, the created RESTClient will use DefaultBurst: 10.
	Burst int `json:"burst"`

	// The maximum length of time to wait before giving up on a server request. A value of zero means no timeout.
	Timeout time.Duration `json:"timeout"`

	// Use default cluster selector
	DefaultClusterSelector bool `json:"defaultClusterSelector"`

	// the default kube config in base64
	Clusters map[string]kubeconfig.ClusterInfo `json:"clusters"`
}
