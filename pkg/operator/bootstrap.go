package operator

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/kube-starter/pkg/kube"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

type bootstrap struct {
	at.AfterInit
	at.EnableScheduling
	configurableFactory factory.ConfigurableFactory
	managerFactory      *instantiate.ScopedInstanceFactory[*Manager]
	manager             *Manager
}

func newBootstrap(
	configurableFactory factory.ConfigurableFactory,
	instFactory *instantiate.ScopedInstanceFactory[*Manager],
) (b *bootstrap, err error) {
	log.Infof("newBootstrap")
	b = &bootstrap{
		configurableFactory: configurableFactory,
		managerFactory:      instFactory,
	}

	// init operator manager
	b.manager, err = b.managerFactory.GetInstance()
	err = b.manager.AddHealthzCheck("healthz", healthz.Ping)
	if err != nil {
		log.Errorf("unable to set up health check: %v", err)
		os.Exit(1)
	}
	err = b.manager.AddReadyzCheck("readyz", healthz.Ping)
	if err != nil {
		log.Errorf("unable to set up ready check: %v", err)
		os.Exit(1)
	}

	// init all operator controllers
	controllers := b.configurableFactory.GetInstances(kube.Controller{})
	instanceContainer := b.managerFactory.GetInstanceContainer()
	if instanceContainer != nil {
		for _, controller := range controllers {
			err = b.configurableFactory.InjectDependency(instanceContainer, controller)
			if err != nil {
				log.Errorf("unable to inject dependency: %v", err)
				os.Exit(1)
			}
		}
	}

	// start operator
	go func() {
		log.Info("starting operator manager")
		err = b.manager.Manager.Start(ctrl.SetupSignalHandler())
		if err != nil {
			log.Error("problem running manager: %v", err)
			os.Exit(1)
		}
	}()

	return
}

func init() {
	app.Register(
		newBootstrap,
		new(instantiate.ScopedInstanceFactory[*Manager]),
	)
}
