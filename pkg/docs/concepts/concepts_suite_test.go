// Copyright (c) 2024 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package concepts

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConcepts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docs Concepts Test Suite")
}
