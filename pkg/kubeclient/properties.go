package kubeclient

import (
	"github.com/hidevopsio/hiboot/pkg/at"
)

// Properties the operator properties
type Properties struct {
	at.ConfigurationProperties `value:"kubeclient"`
	at.AutoWired

	// use DefaultInCluster as default
	DefaultInCluster *bool `json:"defaultInCluster"`
}
