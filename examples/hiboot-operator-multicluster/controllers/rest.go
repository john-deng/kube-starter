package controllers

import (
	goctx "context"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/kube-starter/pkg/kubeclient"
	"github.com/hidevopsio/kube-starter/pkg/oidc"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PodListResponse struct {
	model.BaseResponseInfo
	Data *corev1.PodList `json:"data"`
}

// Controller Rest Controller with path /
// RESTful Controller, derived from at.RestController. The context mapping of this controller is '/' by default
type Controller struct {
	// at.RestController or at.RestController must be embedded here
	at.RestController
	at.RequestMapping `value:"/api/v1/namespaces/{namespace}"`
}

func newController() *Controller {
	log.Infof("newController")
	return &Controller{}
}

func init() {
	app.Register(newController)
}

// ListPods list all pods for specific namespace
func (c *Controller) ListPods(_ struct {
	at.GetMapping `value:"/pods"`
	at.Operation  `id:"List Pods" description:"List Pods of giving namespace"`
	at.Consumes   `values:"application/json"`
	at.Produces   `values:"application/json"`
	Parameters    struct {
		at.Parameter `type:"string" name:"namespace" in:"path" description:"Path Variable（Namespace）" required:"true"`
	}
	Responses struct {
		StatusOK struct {
			at.Parameter `name:"Namespace" in:"body" description:"Get Pod List"`
			PodListResponse
		}
	}
}, namespace string, kubeClient *kubeclient.Client) (response *PodListResponse, err error) {
	response = new(PodListResponse)
	var podList corev1.PodList

	err = kubeClient.List(goctx.TODO(), &podList, client.InNamespace(namespace))
	if err == nil {
		response.Data = &podList
	}

	return
}

// ListPodsByUser list all pods by user
func (c *Controller) ListPodsByUser(_ struct {
	at.GetMapping `value:"/pods/user"`
	at.Operation  `id:"List Pods" description:"List Pods of giving namespace"`
	at.Consumes   `values:"application/json"`
	at.Produces   `values:"application/json"`
	Parameters    struct {
		at.Parameter `type:"string" name:"namespace" in:"path" description:"Path Variable（Namespace）" required:"true"`
	}
	Responses struct {
		StatusOK struct {
			at.Parameter `name:"Namespace" in:"body" description:"Get Pod List"`
			PodListResponse
		}
	}
}, namespace string, runtimeClient *kubeclient.RuntimeClient) (response *PodListResponse, err error) {
	response = new(PodListResponse)
	var podList corev1.PodList

	err = runtimeClient.List(goctx.TODO(), &podList, client.InNamespace(namespace))
	// TODO: error handling
	if err == nil {
		response.Data = &podList
	}

	return
}

type ServiceListResponse struct {
	model.BaseResponseInfo
	Data *corev1.ServiceList `json:"data"`
}

// ListServices list all services
func (c *Controller) ListServices(_ struct {
	at.GetMapping `value:"/services"`
	at.Operation  `id:"List Services" description:"List Services of giving namespace"`
	at.Consumes   `values:"application/json"`
	at.Produces   `values:"application/json"`
	Parameters    struct {
		at.Parameter `type:"string" name:"namespace" in:"path" description:"Path Variable（Namespace）" required:"true"`
	}
	Responses struct {
		StatusOK struct {
			at.Parameter `name:"Namespace" in:"body" description:"Get Service List"`
			ServiceListResponse
		}
	}
}, namespace string, runtimeClient *kubeclient.RuntimeClient) (response *ServiceListResponse, err error) {
	response = new(ServiceListResponse)
	var serviceList corev1.ServiceList

	err = runtimeClient.List(goctx.TODO(), &serviceList, client.InNamespace(namespace))
	if err == nil {
		response.Data = &serviceList
	}

	return
}

type DeploymentListResponse struct {
	model.BaseResponseInfo
	Data *appsv1.DeploymentList `json:"data"`
}

// ListDeployment list all deployments
func (c *Controller) ListDeployment(_ struct {
	at.GetMapping `value:"/deployments"`
	at.Operation  `id:"List Deployments" description:"List Deployments of giving namespace"`
	at.Consumes   `values:"application/json"`
	at.Produces   `values:"application/json"`
	Parameters    struct {
		at.Parameter `type:"string" name:"namespace" in:"path" description:"Path Variable（Namespace）" required:"true"`
	}
	Responses struct {
		StatusOK struct {
			at.Parameter `name:"Namespace" in:"body" description:"Get Deployment List"`
			PodListResponse
		}
	}
}, namespace string, token *oidc.Token, kubeClient *kubeclient.Client) (response *DeploymentListResponse, err error) {
	response = new(DeploymentListResponse)
	var deploymentList appsv1.DeploymentList

	err = kubeClient.List(goctx.TODO(), &deploymentList, client.InNamespace(namespace))
	if err == nil {
		user := "unknown"
		if token.Claims != nil {
			user = token.Claims.Subject
		}
		response.Message = user + " Got Deployment List"
		response.Data = &deploymentList
	}

	return
}
