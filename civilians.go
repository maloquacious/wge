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

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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

// auxCivilian is a helper to convert to/from json
type auxCivilian struct {
	LoyalCitizens int `json:"loyal-citizens"`
	RebelCitizens int `json:"rebel-citizens"`
	TechLevel     int `json:"tech-level"`
}

func NewCivilian(pop, techLevel int) Civilian {
	var p Civilian
	p.qty.loyal = pop
	p.techLevel = techLevel
	return p
}

// BirthRate implements the PopulationGroup interface.
// The basic birth rate ranges from 0.25% to 10% of the population.
// The variation depends on the standard of living as well as the
// availability of "open" living space in the colony.
func (p Civilian) BirthRate(standardOfLiving, pctCapacity float64) float64 {
	if p.IsOnShip() { // births never happen on a ship
		return 0
	}

	// clamp the standard of living and percent capacity
	standardOfLiving = clamp(standardOfLiving, 0.01, 3.0)
	pctCapacity = clamp(pctCapacity, 0.01, 1.0)

	// the base rate is determined by tech level
	birthRate := float64(11-p.techLevel) * 0.1
	if birthRate < 0.0025 {
		birthRate = 0.0025
	} else if birthRate > 0.10 {
		birthRate = 0.10
	}

	// resort colonies increase the birth rate
	if p.IsResortColony() {
		birthRate *= 2
	}

	// standard of living influences it
	if standardOfLiving < 0.25 {
		birthRate *= 1.5
	} else if standardOfLiving < 0.80 {
		birthRate *= 1.25
	} else if standardOfLiving < 1.20 {
		// 80% to 120% is the standard range
	} else if standardOfLiving > 1.20 {
		birthRate *= 0.75
	} else if standardOfLiving > 1.75 {
		birthRate *= 0.5
	}

	// overcrowding reduces the birth rate
	if pctCapacity < 0.25 {
		birthRate *= 1.25
	} else if pctCapacity < 0.40 {
		birthRate *= 1.10
	} else if pctCapacity < 0.65 {
		// 40% to 65% is the standard range
	} else if pctCapacity < 0.70 {
		birthRate *= 0.90
	} else if pctCapacity < 0.80 {
		birthRate *= 0.60
	} else if pctCapacity < 0.90 {
		birthRate *= 0.25
	} else if pctCapacity < 0.95 {
		birthRate *= 0.1
	} else {
		birthRate *= 0.05
	}

	// birth rate is never less than 0.25% or higher than 10%
	return clamp(birthRate, 0.0025, 0.10)
}

// Code implements the Unit interface.
func (p Civilian) Code() string {
	return "UNK"
}

// FoodNeeded implements the PopulationGroup interface
func (p Civilian) FoodNeeded() float64 {
	return float64(p.qty.loyal+p.qty.rebel) * 0.01 * 0.0125
}

// IsResortColony returns true if the population is in a resort colony
func (p Civilian) IsResortColony() bool {
	return false
}

// IsOnShip returns true if the population is on a ship.
func (p Civilian) IsOnShip() bool {
	return false
}

// LifeSupportNeeded implements the PopulationGroup interface
func (p Civilian) LifeSupportNeeded() float64 {
	return float64(p.qty.loyal+p.qty.rebel) * 0.01 * 0.5
}

// MarshalJSON implements the json.Marshaler interface
func (p Civilian) MarshalJSON() ([]byte, error) {
	var aux auxCivilian
	aux.LoyalCitizens = p.qty.loyal
	aux.RebelCitizens = p.qty.rebel
	aux.TechLevel = p.techLevel
	return json.Marshal(&aux)
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

// Population implements the PopulationGroup interface.
func (p Civilian) Population() int {
	return p.qty.loyal + p.qty.rebel
}

// Rebels implements the PopulationGroup interface.
func (p Civilian) Rebels() int {
	return p.qty.rebel
}

// TechLevel implements the TechLevel interface.
func (p Civilian) TechLevel() int {
	return p.techLevel
}

// UnmarshalJSON implements the json.Marshaler interface
func (p *Civilian) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	var aux auxCivilian
	if err := dec.Decode(&aux); err != nil {
		return fmt.Errorf("decode civilian: %w", err)
	}

	p.qty.loyal = aux.LoyalCitizens
	p.qty.rebel = aux.RebelCitizens
	p.techLevel = aux.TechLevel

	return nil
}
