// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package indexer

import (
	"context"
	"fmt"

	"github.com/gardener/gardener/pkg/apis/seedmanagement"
	seedmanagementv1alpha1 "github.com/gardener/gardener/pkg/apis/seedmanagement/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AddManagedSeedShootName adds an index for seedmanagement.ManagedSeedShootName to the given indexer.
func AddManagedSeedShootName(ctx context.Context, indexer client.FieldIndexer) error {
	if err := indexer.IndexField(ctx, &seedmanagementv1alpha1.ManagedSeed{}, seedmanagement.ManagedSeedShootName, func(obj client.Object) []string {
		managedSeed, ok := obj.(*seedmanagementv1alpha1.ManagedSeed)
		if !ok {
			return []string{""}
		}
		if managedSeed.Spec.Shoot == nil {
			return []string{""}
		}
		return []string{managedSeed.Spec.Shoot.Name}
	}); err != nil {
		return fmt.Errorf("failed to add indexer for %s to ManagedSeed Informer: %w", seedmanagement.ManagedSeedShootName, err)
	}
	return nil
}