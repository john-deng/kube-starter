package operator

import (
	"fmt"
	"github.com/hidevopsio/kube-starter/pkg/kubeclient"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strconv"
	"time"

	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	_ "github.com/hidevopsio/kube-starter/pkg/kubeclient"
	"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	Profile = "operator"
)

type Manager struct {
	at.Scope `value:"prototype"`

	manager.Manager
}

type configuration struct {
	at.AutoConfiguration

	Properties *Properties
	portOffset int
}

func newConfiguration() *configuration {
	return &configuration{}
}

func init() {
	app.Register(newConfiguration)
}

// Manager is the controller runtime manager
func (c *configuration) Manager(scheme *runtime.Scheme, cfg *kubeclient.RestConfig) (mgr *Manager, err error) {
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

	var port string
	port, err = addOffsetToPort(c.Properties.MetricsBindAddress, c.portOffset)
	options.Metrics.BindAddress = port
	port, err = addOffsetToPort(c.Properties.HealthProbeBindAddress, c.portOffset)
	options.HealthProbeBindAddress = port
	options.LeaderElection = c.Properties.LeaderElection
	options.WebhookServer = webhook.NewServer(webhook.Options{
		Port: c.Properties.Port + c.portOffset, // Specify your desired port
	})

	c.portOffset = c.portOffset + 1

	log.Infof("started operator with qps: %v, burst: %v", cfg.QPS, cfg.Burst)
	mgr = new(Manager)
	mgr.Manager, err = ctrl.NewManager(cfg.Config, options)

	if err != nil {
		log.Errorf("ctrl.NewManager() returned error: %v", err)
	}
	return
}

// addOffsetToPort takes a port string (with leading colon) and an offset integer,
// and returns a new port string with the offset applied.
func addOffsetToPort(port string, offset int) (string, error) {
	// Strip the leading colon and convert the port number to an integer
	portNumber, err := strconv.Atoi(port[1:])
	if err != nil {
		return "", fmt.Errorf("error converting port: %v", err)
	}

	// Add the offset to the port number
	newPortNumber := portNumber + offset

	// Convert the new port number back to a string and add the leading colon
	newPort := ":" + strconv.Itoa(newPortNumber)

	return newPort, nil
}
