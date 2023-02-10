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

type Population interface {
	// the number of FOOD units needed to sustain the population.
	FoodNeeded() float64
	// the number of LS units needed to sustain the population
	LifeSupportNeeded() float64
	// Population is the total population of the unit.
	Population() int
	// Rebels is the number of rebels in the populations.
	Rebels() int
}

// Civilian is a population unit composed of the bourgeoisie, retirees,
// stay-at-home parents, and the unemployed.
// The state can order civilians to relocate to other planets or systems.
type Civilian struct {
	qty struct {
		loyal int
		rebel int
	}
	techLevel int
}

func NewCivilian(pop, techLevel int) Civilian {
	var p Civilian
	p.qty.loyal = pop
	p.techLevel = techLevel
	return p
}

// Code implements the Unit interface.
func (p Civilian) Code() string {
	return "UNK"
}

// FoodNeeded implements the Population interface
func (p Civilian) FoodNeeded() float64 {
	return float64(p.qty.loyal+p.qty.rebel) * 0.01 * 0.0125
}

// LifeSupportNeeded implements the Population interface
func (p Civilian) LifeSupportNeeded() float64 {
	return float64(p.qty.loyal+p.qty.rebel) * 0.01 * 0.5
}

// Merge combines two population units.
// Rebel population and tech levels are calculated as the weighted average of the units.
func (p Civilian) Merge(q Civilian) Civilian {
	if p.Population() == 0 {
		return q
	} else if q.Population() == 0 {
		return p
	}

	var n Civilian
	n.qty.loyal, n.qty.rebel = p.qty.loyal+q.qty.loyal, p.qty.rebel+q.qty.rebel
	deltaRebels := 0 // merging units always increases discontent
	if p.techLevel == q.techLevel {
		n.techLevel = p.techLevel
	} else {
		pTech, qTech := p.Population()*p.techLevel, q.Population()*q.techLevel
		n.techLevel = (pTech + qTech) / (p.Population() + q.Population())
		// the group losing tech levels gets especially cranky
		if n.techLevel < p.techLevel {
			deltaTech := p.techLevel - n.techLevel
			deltaRebels = p.qty.rebel * deltaTech / 100
		} else if n.techLevel < q.techLevel {
			deltaTech := p.techLevel - n.techLevel
			deltaRebels = p.qty.rebel * deltaTech / 100
		}
	}
	if deltaRebels < 1 {
		deltaRebels = 1
	}
	n.qty.loyal, n.qty.rebel = n.qty.loyal-deltaRebels, n.qty.rebel+deltaRebels

	return n
}

// Population implements the Population interface.
func (p Civilian) Population() int {
	return p.qty.loyal + p.qty.rebel
}

// Rebels implements the Population interface.
func (p Civilian) Rebels() int {
	return p.qty.rebel
}

// TechLevel implements the TechLevel interface.
func (p Civilian) TechLevel() int {
	return p.techLevel
}
