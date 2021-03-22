package kubeclient

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hidevopsio/hiboot/pkg/app/web/context"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/kube-starter/pkg/jwt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Subject = "subject"
)

// DefaultKubeconfig
func DefaultKubeconfig() string {
	fname := os.Getenv("KUBECONFIG")
	if fname != "" {
		return fname
	}
	home, err := os.UserHomeDir()
	if err != nil {
		klog.Warningf("failed to get home directory: %v", err)
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

// Kubeconfig
func Kubeconfig() (cfg *rest.Config, err error) {
	cfg, err = rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", DefaultKubeconfig())
		if err != nil {
			klog.Warningf("Error building kubeconfig: %s", err.Error())
		}
	}
	return
}

// KubeClient
func KubeClient(scheme *runtime.Scheme, cfg *RestConfig) (k8sClient client.Client, err error)  {
	if cfg.Config == nil {
		return
	}
	k8sClient, err = client.New(cfg.Config, client.Options{Scheme: scheme})
	return
}

// RuntimeKubeClient
func RuntimeKubeClient(ctx context.Context, scheme *runtime.Scheme, cfg *RestConfig) (runtimeClient *RuntimeClient, err error)  {
	if cfg.Config == nil {
		return
	}
	runtimeClient = new(RuntimeClient)
	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Replace(bearerToken, "Bearer ", "", -1)
	var claims *jwt.Claims
	claims, err = jwt.DecodeWithoutVerify(token)
	if err == nil {
		cfg.Impersonate.UserName = claims.Issuer + "#" + claims.Subject
	} else {
		// unauthorized user
		ctx.StatusCode(http.StatusUnauthorized)
		scheme = runtime.NewScheme()
	}
	runtimeClient.Claims = claims
	runtimeClient.Context = ctx
	runtimeClient.Client, err = client.New(cfg.Config, client.Options{Scheme: scheme})
	if err != nil {
		log.Error(err)
	}
	return
}

