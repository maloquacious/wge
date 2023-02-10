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
	m := wge.NewCivilian(100, 2).Merge(wge.NewCivilian(100, 2))
	if m.TechLevel() != 2 {
		t.Errorf("merge: expected tech-level %d, got %d\n", 2, m.TechLevel())
	}
	m = wge.NewCivilian(100, 2).Merge(wge.NewCivilian(100, 4))
	if m.TechLevel() != 3 {
		t.Errorf("merge: expected tech-level %d, got %d\n", 3, m.TechLevel())
	}
	m = wge.NewCivilian(100, 2).Merge(wge.NewCivilian(100, 6))
	if m.TechLevel() != 4 {
		t.Errorf("merge: expected tech-level %d, got %d\n", 4, m.TechLevel())
	}
	m = wge.NewCivilian(300, 2).Merge(wge.NewCivilian(100, 6))
	if m.TechLevel() != 3 {
		t.Errorf("merge: expected tech-level %d, got %d\n", 3, m.TechLevel())
	}
}
