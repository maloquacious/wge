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

package wge

// PopulationGroup defines the interface for working with groups of people.
type PopulationGroup interface {
	// FoodNeeded returns the number of FOOD units needed to sustain the population.
	FoodNeeded() float64
	// LifeSupportNeeded returns the number of LS units needed to sustain the population.
	LifeSupportNeeded() float64
	// NaturalBirthRate returns the percentage of natural births in the group.
	NaturalBirthRate() float64
	// Population returns total population of the unit.
	Population() int
	// Rebels returns the number of rebels in the population.
	Rebels() int
}
