// wge - the wraith game engine
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package wge_test

import (
	"testing"

	"github.com/maloquacious/wge"
)

func TestCivilians(t *testing.T) {
	// verify that when tech is merged, we use the weighted
	// average (rounded down) for the new population
	for _, tc := range []struct {
		id          int
		pPop, pTech int
		qPop, qTech int
		expect      int
	}{
		{1, 100, 2, 100, 2, 2},
		{2, 100, 2, 100, 4, 3},
		{3, 100, 2, 100, 6, 4},
		{4, 300, 2, 100, 6, 3},
		{5, 100, 10, 1000, 1, 1},
	} {
		p := wge.NewCivilian(tc.pPop, tc.pTech)
		q := wge.NewCivilian(tc.qPop, tc.qTech)
		m := p.Merge(q)
		if tc.expect != m.TechLevel() {
			t.Errorf("merge: %d: expected tech-level %d, got %d\n", tc.id, tc.expect, m.TechLevel())
		}
	}
}
