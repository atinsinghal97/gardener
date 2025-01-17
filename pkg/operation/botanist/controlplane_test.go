// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package botanist

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/operation"
	mockcontrolplane "github.com/gardener/gardener/pkg/operation/botanist/component/extensions/controlplane/mock"
	mockdnsrecord "github.com/gardener/gardener/pkg/operation/botanist/component/extensions/dnsrecord/mock"
	mockinfrastructure "github.com/gardener/gardener/pkg/operation/botanist/component/extensions/infrastructure/mock"
	"github.com/gardener/gardener/pkg/operation/shoot"
)

var _ = Describe("controlplane", func() {
	var (
		ctrl *gomock.Controller

		infrastructure       *mockinfrastructure.MockInterface
		controlPlane         *mockcontrolplane.MockInterface
		controlPlaneExposure *mockcontrolplane.MockInterface
		externalDNSRecord    *mockdnsrecord.MockInterface
		internalDNSRecord    *mockdnsrecord.MockInterface
		botanist             *Botanist

		ctx     = context.TODO()
		fakeErr = fmt.Errorf("fake err")
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		infrastructure = mockinfrastructure.NewMockInterface(ctrl)
		controlPlane = mockcontrolplane.NewMockInterface(ctrl)
		controlPlaneExposure = mockcontrolplane.NewMockInterface(ctrl)
		externalDNSRecord = mockdnsrecord.NewMockInterface(ctrl)
		internalDNSRecord = mockdnsrecord.NewMockInterface(ctrl)

		botanist = &Botanist{
			Operation: &operation.Operation{
				Shoot: &shoot.Shoot{
					Components: &shoot.Components{
						Extensions: &shoot.Extensions{
							ControlPlane:         controlPlane,
							ControlPlaneExposure: controlPlaneExposure,
							ExternalDNSRecord:    externalDNSRecord,
							InternalDNSRecord:    internalDNSRecord,
							Infrastructure:       infrastructure,
						},
					},
				},
			},
		}
		botanist.Shoot.SetInfo(&gardencorev1beta1.Shoot{})
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#DeployControlPlane", func() {
		var infrastructureStatus = &runtime.RawExtension{Raw: []byte("infra-status")}

		BeforeEach(func() {
			infrastructure.EXPECT().ProviderStatus().Return(infrastructureStatus)
			controlPlane.EXPECT().SetInfrastructureProviderStatus(infrastructureStatus)
		})

		Context("deploy", func() {
			It("should deploy successfully", func() {
				controlPlane.EXPECT().Deploy(ctx)
				Expect(botanist.DeployControlPlane(ctx)).To(Succeed())
			})

			It("should return the error during deployment", func() {
				controlPlane.EXPECT().Deploy(ctx).Return(fakeErr)
				Expect(botanist.DeployControlPlane(ctx)).To(MatchError(fakeErr))
			})
		})

		Context("restore", func() {
			var shootState = &gardencorev1beta1.ShootState{}

			BeforeEach(func() {
				botanist.SetShootState(shootState)
				botanist.Shoot.SetInfo(&gardencorev1beta1.Shoot{
					Status: gardencorev1beta1.ShootStatus{
						LastOperation: &gardencorev1beta1.LastOperation{
							Type: gardencorev1beta1.LastOperationTypeRestore,
						},
					},
				})
			})

			It("should restore successfully", func() {
				controlPlane.EXPECT().Restore(ctx, shootState)
				Expect(botanist.DeployControlPlane(ctx)).To(Succeed())
			})

			It("should return the error during restoration", func() {
				controlPlane.EXPECT().Restore(ctx, shootState).Return(fakeErr)
				Expect(botanist.DeployControlPlane(ctx)).To(MatchError(fakeErr))
			})
		})
	})

	Describe("#DeployControlPlaneExposure()", func() {
		Context("deploy", func() {
			It("should deploy successfully", func() {
				controlPlaneExposure.EXPECT().Deploy(ctx)
				Expect(botanist.DeployControlPlaneExposure(ctx)).To(Succeed())
			})

			It("should return the error during deployment", func() {
				controlPlaneExposure.EXPECT().Deploy(ctx).Return(fakeErr)
				Expect(botanist.DeployControlPlaneExposure(ctx)).To(MatchError(fakeErr))
			})
		})

		Context("restore", func() {
			var shootState = &gardencorev1beta1.ShootState{}

			BeforeEach(func() {
				botanist.SetShootState(shootState)
				botanist.Shoot.SetInfo(&gardencorev1beta1.Shoot{
					Status: gardencorev1beta1.ShootStatus{
						LastOperation: &gardencorev1beta1.LastOperation{
							Type: gardencorev1beta1.LastOperationTypeRestore,
						},
					},
				})
			})

			It("should restore successfully", func() {
				controlPlaneExposure.EXPECT().Restore(ctx, shootState)
				Expect(botanist.DeployControlPlaneExposure(ctx)).To(Succeed())
			})

			It("should return the error during restoration", func() {
				controlPlaneExposure.EXPECT().Restore(ctx, shootState).Return(fakeErr)
				Expect(botanist.DeployControlPlaneExposure(ctx)).To(MatchError(fakeErr))
			})
		})
	})
})
