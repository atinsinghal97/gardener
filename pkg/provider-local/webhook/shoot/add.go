// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shoot

import (
	"os"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/extensions/pkg/webhook/shoot"
)

var (
	logger = log.Log.WithName("local-shoot-webhook")

	// DefaultAddOptions are the default AddOptions for AddToManager.
	DefaultAddOptions = AddOptions{}
)

// AddOptions are options to apply when adding the local shoot webhook to the manager.
type AddOptions struct{}

// AddToManagerWithOptions creates a webhook with the given options and adds it to the manager.
func AddToManagerWithOptions(mgr manager.Manager, _ AddOptions) (*extensionswebhook.Webhook, error) {
	logger.Info("Adding webhook to manager")

	failurePolicy := admissionregistrationv1.Fail

	args := shoot.Args{
		Types:         []extensionswebhook.Type{{Obj: &corev1.ConfigMap{}}},
		Mutator:       NewMutator(),
		FailurePolicy: &failurePolicy,
	}

	// With https://github.com/gardener/gardener/pull/7578, provider-local temporarily mutates the CoreDNS configuration
	// since the cluster network is not working properly in the shoots. However, this breaks the development setup
	// (https://github.com/gardener/gardener/blob/master/docs/development/getting_started_locally.md) because it is not
	// possible to determine the shoot namespace when running out of the seed cluster. Hence, `MutatorWithShootClient`
	// cannot be used.
	// Since all this is just a temporary workaround and was motivated by fixing flaky e2e tests, we can also work
	// without it for development scenarios for the time being.
	// So when the `PROVIDER_LOCAL_DISABLE_COREDNS_MUTATION` environment variable is not set then we enable the webhook
	// (default scenario), otherwise we continue with above mutator (only for the local development scenario from
	// above).
	if os.Getenv("PROVIDER_LOCAL_DISABLE_COREDNS_MUTATION") == "" {
		args.Mutator = nil
		args.MutatorWithShootClient = NewMutatorWithShootClient(NewMutator())
		args.Types = append(args.Types,
			extensionswebhook.Type{Obj: &corev1.Node{}},
			extensionswebhook.Type{Obj: &corev1.Service{}},
			extensionswebhook.Type{Obj: &appsv1.Deployment{}},
		)
	}

	return shoot.New(mgr, args)
}

// AddToManager creates a webhook with the default options and adds it to the manager.
func AddToManager(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	return AddToManagerWithOptions(mgr, DefaultAddOptions)
}
