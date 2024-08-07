package operator

import (
	"time"

	"github.com/hidevopsio/hiboot/pkg/at"
)

// Properties the operator properties
type Properties struct {
	at.ConfigurationProperties `value:"operator"`
	at.AutoWired

	Development bool `json:"development"`

	// SyncPeriod determines the minimum frequency at which watched resources are
	// reconciled. A lower period will correct entropy more quickly, but reduce
	// responsiveness to change if there are many watched resources. Change this
	// value only if you know what you are doing. Defaults to 10 hours if unset.
	// there will a 10 percent jitter between the SyncPeriod of all controllers
	// so that all controllers will not send list requests simultaneously.
	SyncPeriod *time.Duration `json:"syncPeriod"`

	// LeaderElection determines whether or not to use leader election when
	// starting the manager.
	LeaderElection bool `json:"leaderElection"`

	// LeaderElectionResourceLock determines which resource lock to use for leader election,
	// defaults to "configmapsleases". Change this value only if you know what you are doing.
	// Otherwise, users of your controller might end up with multiple running instances that
	// each acquired leadership through different resource locks during upgrades and thus
	// act on the same resources concurrently.
	// If you want to migrate to the "leases" resource lock, you might do so by migrating to the
	// respective multilock first ("configmapsleases" or "endpointsleases"), which will acquire a
	// leader lock on both resources. After all your users have migrated to the multilock, you can
	// go ahead and migrate to "leases". Please also keep in mind, that users might skip versions
	// of your controller.
	//
	// Note: before controller-runtime version v0.7, the resource lock was set to "configmaps".
	// Please keep this in mind, when planning a proper migration path for your controller.
	LeaderElectionResourceLock string `json:"leaderElectionResourceLock"`

	// LeaderElectionNamespace determines the namespace in which the leader
	// election resource will be created.
	LeaderElectionNamespace string `json:"leaderElectionNamespace"`

	// LeaderElectionID determines the name of the resource that leader election
	// will use for holding the leader lock.
	LeaderElectionID string `json:"leaderElectionId"`

	// LeaderElectionConfig can be specified to override the default configuration
	// that is used to build the leader election client.
	//LeaderElectionConfig *rest.Config

	// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
	// when the Manager ends. This requires the binary to immediately end when the
	// Manager is stopped, otherwise this setting is unsafe. Setting this significantly
	// speeds up voluntary leader transitions as the new leader doesn't have to wait
	// LeaseDuration time first.
	LeaderElectionReleaseOnCancel bool `json:"leaderElectionReleaseOnCancel"`

	// LeaseDuration is the duration that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack. Default is 15 seconds.
	LeaseDuration *time.Duration `json:"leaseDuration"`
	// RenewDeadline is the duration that the acting controlplane will retry
	// refreshing leadership before giving up. Default is 10 seconds.
	RenewDeadline *time.Duration `json:"renewDeadline"`
	// RetryPeriod is the duration the LeaderElector clients should wait
	// between tries of actions. Default is 2 seconds.
	RetryPeriod *time.Duration `json:"retryPeriod"`

	// Namespace if specified restricts the manager's cache to watch objects in
	// the desired namespace Defaults to all namespaces
	//
	// Note: If a namespace is specified, controllers can still Watch for a
	// cluster-scoped resource (e.g Node).  For namespaced resources the cache
	// will only hold objects from the desired namespace.
	Namespace string `json:"namespace"`

	// MetricsBindAddress is the TCP address that the controller should bind to
	// for serving prometheus metrics.
	// It can be set to "0" to disable the metrics serving.
	MetricsBindAddress string `json:"metricsBindAddress" default:":9000"`

	// HealthProbeBindAddress is the TCP address that the controller should bind to
	// for serving health probes
	HealthProbeBindAddress string `json:"healthProbeBindAddress" default:":9100"`

	// Readiness probe endpoint name, defaults to "readyz"
	ReadinessEndpointName string `json:"readinessEndpointName"`

	// Liveness probe endpoint name, defaults to "healthz"
	LivenessEndpointName string `json:"livenessEndpointName"`

	// Port is the port that the webhook server serves at.
	// It is used to set webhook.Server.Port.
	Port int `json:"port" default:"8000"`
	// Host is the hostname that the webhook server binds to.
	// It is used to set webhook.Server.Host.
	Host string `json:"host"`

	// CertDir is the directory that contains the server key and certificate.
	// if not set, webhook server would look up the server key and certificate in
	// {TempDir}/k8s-webhook-server/serving-certs. The server key and certificate
	// must be named tls.key and tls.crt, respectively.
	CertDir string `json:"certDir"`
	// Functions to all for a user to customize the values that will be injected.

	// DryRunClient specifies whether the client should be configured to enforce
	// dryRun mode.
	DryRunClient bool `json:"dryRunClient"`
}
