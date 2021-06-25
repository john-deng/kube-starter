package operator

import (
	"time"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	Profile = "operator"
)

type configuration struct {
	at.AutoConfiguration

	Properties *Properties
}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

// Manager is the controller runtime manager
func (c *configuration) Manager(scheme *runtime.Scheme) (mgr manager.Manager, err error) {
	opts := zap.Options{
		Development: c.Properties.Development,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	var options ctrl.Options
	_ = copier.CopyWithOption(&options, &c.Properties, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	options.Scheme = scheme
	if c.Properties.LeaseDuration != nil {
		second := *c.Properties.LeaseDuration * time.Second
		options.LeaseDuration = &second
	}
	if c.Properties.RenewDeadline != nil {
		second := *c.Properties.RenewDeadline * time.Second
		options.RenewDeadline = &second
	}
	if c.Properties.RetryPeriod != nil {
		second := *c.Properties.RetryPeriod * time.Second
		options.RetryPeriod = &second
	}
	if c.Properties.SyncPeriod != nil {
		second := *c.Properties.SyncPeriod * time.Second
		options.SyncPeriod = &second
	}

	options.LeaderElectionConfig, err = Kubeconfig()

	mgr, err = ctrl.NewManager(ctrl.GetConfigOrDie(), options)

	if err != nil {
		log.Error(err)
	}
	return
}
