// Copyright (c) 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package constants

const (
	// LabelKey is the key of a label used for the identification of CoreDNS pods.
	LabelKey = "k8s-app"
	// LabelValue is the value of a label used for the identification of CoreDNS pods (it's 'kube-dns' for legacy
	// reasons).
	LabelValue = "kube-dns"
	// PortServiceServer is the service port used for the DNS server.
	PortServiceServer = 53
	// PortServer is the target port used for the DNS server.
	PortServer = 8053
)
