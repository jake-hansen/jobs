// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package utils_test

import (
	"testing"

	"github.com/jake-hansen/jobs/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewAtomicBool(t *testing.T) {
	t.Run("success-init-false", func(t *testing.T) {
		atomicBool := utils.NewAtomicBool(false)

		assert.Equal(t, atomicBool.SafeGet(), false)
	})
	t.Run("success-init-true", func(t *testing.T) {
		atomicBool := utils.NewAtomicBool(true)

		assert.Equal(t, true, atomicBool.SafeGet())
	})
}

func TestAtomicBool_SafeSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		atomicBool := utils.NewAtomicBool(false)
		go atomicBool.SafeSet(true)

		assert.Equal(t, false, atomicBool.SafeGet())
	})
}

func TestAtomicBool_SafeGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		atomicBool := utils.NewAtomicBool(false)

		go atomicBool.SafeGet()
		go atomicBool.SafeGet()

		assert.Equal(t, false, atomicBool.SafeGet())
	})
}
