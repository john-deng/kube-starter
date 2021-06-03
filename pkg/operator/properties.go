package operator

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Properties the operator properties
type Properties struct {
	at.ConfigurationProperties `value:"operator"`
	at.AutoWired

	ctrl.Options
}
