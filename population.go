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
	// Births returns the number of new units added from natural births.
	Births() int
	// FoodNeeded returns the number of FOOD units needed to sustain the population.
	FoodNeeded() float64
	// LifeSupportNeeded returns the number of LS units needed to sustain the population.
	LifeSupportNeeded() float64
	// Population returns total population of the unit.
	Population() int
	// Rebels returns the number of rebels in the population.
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

// Births implements the Population interface.
// The basic birth rate ranges from 0.25% to 10% of the population.
// The variation depends on the standard of living as well as the
// availability of "open" living space in the colony.
func (p Civilian) Births(standardOfLiving, pctCapacity float64) int {
	if p.IsOnShip() { // births never happen on a ship
		return 0
	}

	// rate is determined by tech level, standard of living, and how crowded the colony is.
	birthRate := float64(11-p.techLevel) * standardOfLiving * (1.0 - pctCapacity)
	// birth rate is never less than 0.25% or higher than 10%
	if birthRate < 0.0025 {
		birthRate = 0.0025
	} else if birthRate > 0.1 {
		birthRate = 0.1
	}

	return int(float64(p.Population()) * birthRate)
}

// Code implements the Unit interface.
func (p Civilian) Code() string {
	return "UNK"
}

// FoodNeeded implements the Population interface
func (p Civilian) FoodNeeded() float64 {
	return float64(p.qty.loyal+p.qty.rebel) * 0.01 * 0.0125
}

// IsOnShip returns true if the population is on a ship.
func (p Civilian) IsOnShip() bool {
	return false
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
