package operator

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/instantiate"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/kube-starter/pkg/kube"
	"github.com/hidevopsio/kube-starter/pkg/kubeclient"
	"github.com/hidevopsio/kube-starter/pkg/kubeconfig"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

type bootstrap struct {
	at.AfterInit
	at.EnableScheduling
	configurableFactory factory.ConfigurableFactory
	managerFactory      *instantiate.ScopedInstanceFactory[*Manager]
	clientFactory       *instantiate.ScopedInstanceFactory[*kubeclient.ClientCreation]
	managers            map[string]*Manager
	controllers         []*factory.MetaData
	instanceContainer   factory.InstanceContainer
}

func newBootstrap(
	configurableFactory factory.ConfigurableFactory,
	instFactory *instantiate.ScopedInstanceFactory[*Manager],
	clientFactory *instantiate.ScopedInstanceFactory[*kubeclient.ClientCreation],
	prop *kubeclient.Properties,
) (b *bootstrap, err error) {
	log.Infof("newBootstrap")
	b = &bootstrap{
		configurableFactory: configurableFactory,
		managerFactory:      instFactory,
		managers:            make(map[string]*Manager),
		clientFactory:       clientFactory,
	}
	// init all operator controllers
	// TODO: for loop to iterate all connected clusters
	// TODO: watch new connected cluster
	b.controllers = b.configurableFactory.GetInstances(kube.Controller{})
	var ctx = ctrl.SetupSignalHandler()
	for cm, cluster := range prop.Clusters {
		// init operator manager
		var manager *Manager
		var instanceContainer factory.InstanceContainer
		instanceContainer, err = b.managerFactory.GetInstanceContainer(&kubeconfig.ClusterConfig{
			ClusterInfo: cluster,
		})
		if err != nil {
			log.Errorf("unable to get instance container for manager: %v", err)
			os.Exit(1)
		}
		// new manager
		manager = b.managerFactory.GetInstanceFromContainer(instanceContainer)
		err = manager.AddHealthzCheck("healthz", healthz.Ping)
		if err != nil {
			log.Errorf("unable to set up health check: %v", err)
			os.Exit(1)
		}
		err = manager.AddReadyzCheck("readyz", healthz.Ping)
		if err != nil {
			log.Errorf("unable to set up ready check: %v", err)
			os.Exit(1)
		}

		err = b.registerControllers(instanceContainer)

		// start operator
		go func() {
			log.Info("starting operator manager")
			err = manager.Manager.Start(ctx)
			if err != nil {
				log.Error("problem running manager: %v", err)
				os.Exit(1)
			}
		}()

		b.managers[cm] = manager
	}

	return
}

func (b *bootstrap) registerControllers(instanceContainer factory.InstanceContainer) (err error) {
	if instanceContainer != nil {
		for _, controller := range b.controllers {
			err = b.configurableFactory.InjectDependency(instanceContainer, controller)
			if err != nil {
				log.Errorf("unable to inject dependency: %v", err)
				os.Exit(1)
			}
		}
	}
	return err
}

func init() {
	app.Register(
		newBootstrap,
		new(instantiate.ScopedInstanceFactory[*Manager]),
	)
}
