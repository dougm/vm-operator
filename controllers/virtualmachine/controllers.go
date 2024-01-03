// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package virtualmachine

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/vmware-tanzu/vm-operator/controllers/virtualmachine/v1alpha1"
	"github.com/vmware-tanzu/vm-operator/controllers/virtualmachine/v1alpha2"
	pkgconfig "github.com/vmware-tanzu/vm-operator/pkg/config"
	"github.com/vmware-tanzu/vm-operator/pkg/context"
)

// AddToManager adds the controller to the provided manager.
func AddToManager(ctx *context.ControllerManagerContext, mgr manager.Manager) error {
	if pkgconfig.FromContext(ctx).Features.VMOpV1Alpha2 {
		return v1alpha2.AddToManager(ctx, mgr)
	}
	return v1alpha1.AddToManager(ctx, mgr)
}
