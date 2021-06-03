/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//Here is a few steps to generate an operator
//operator-sdk init --domain icloudnative.net --repo github.com/hidevopsio/kube-starter/examples/demo-operator
//go:generate operator-sdk create api --group admin --version v1alpha1 --kind User --resource --namespaced=false --controller=false
//go:generate make generate manifests
//go:generate kubectl apply -f config/crd/bases

package main

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"github.com/hidevopsio/kube-starter/pkg/operator"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	adminv1alpha1 "github.com/hidevopsio/kube-starter/examples/demo-operator/apis/admin/v1alpha1"
	_ "github.com/hidevopsio/kube-starter/examples/demo-operator/controllers/admin"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	//+kubebuilder:scaffold:imports
)

func init() {
	var scheme = runtime.NewScheme()

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(adminv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme


	app.Register(scheme)
}

func main() {
	web.NewApplication().
		SetProperty(logging.Level, logging.LevelDebug).
		SetProperty(app.ProfilesInclude,
			actuator.Profile,
			logging.Profile,
			operator.Profile,
		).
		Run()
}
