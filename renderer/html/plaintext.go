// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package html

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
)

// PlaintextRenderer renders plaintext in a single "<pre>" element.
// Use "NewPlaintextRenderer" to instantiate.
type PlaintextRenderer struct {
	policy *bluemonday.Policy
}

func NewPlaintextRenderer(policy *bluemonday.Policy) PlaintextRenderer {
	p := policy
	if p == nil {
		p = bluemonday.UGCPolicy()
	}

	return PlaintextRenderer{
		policy: p,
	}
}

func (r *PlaintextRenderer) Render(code []byte, _ Transformer) ([]byte, error) {
	safe := r.policy.SanitizeBytes(code)

	return fmt.Appendf([]byte{}, "<pre>%s</pre>", safe), nil
}
