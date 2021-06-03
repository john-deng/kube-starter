package operator

import (
	"flag"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
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
type Manager struct {
	//at.ContextAware
	ctrl.Manager
}

func (c *configuration) Manager(scheme *runtime.Scheme) (mgr *Manager, err error) {
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	var m ctrl.Manager
	m, err = ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     c.Properties.MetricsBindAddress,
		Port:                   c.Properties.Port,
		HealthProbeBindAddress: c.Properties.HealthProbeBindAddress,
		LeaderElection:         c.Properties.LeaderElection,
		LeaderElectionID:       c.Properties.LeaderElectionID,
	})

	if err != nil {
		log.Error(err)
		return
	}

	mgr = &Manager{
		Manager: m,
	}
	return
}
