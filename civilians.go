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

// auxCivilian is a helper to convert to/from json.
// used to implement json.Marshaler and json.Unmarshaler interfaces.
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

// Code implements the Unit interface.
func (p Civilian) Code() string {
	return "CIV"
}

// FoodNeeded implements the PopulationGroup interface
func (p Civilian) FoodNeeded() float64 {
	return float64(p.qty.loyal+p.qty.rebel) * 0.01 * 0.0125
}

// IsOnClosedColony returns true if the population is on a closed colony.
func (p Civilian) IsOnClosedColony() bool {
	return false
}

// IsOnLifeSupport returns true if the population depends on life support for survival.
// This is true for all ships and closed colonies.
func (p Civilian) IsOnLifeSupport() bool {
	return p.IsOnShip() || p.IsOnClosedColony()
}

// IsOnOpenColony returns true if the population is on an open colony.
func (p Civilian) IsOnOpenColony() bool {
	return false
}

// IsOnShip returns true if the population is on a ship.
func (p Civilian) IsOnShip() bool {
	return false
}

// IsResortColony returns true if the population is in a resort colony
func (p Civilian) IsResortColony() bool {
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

// Mass implements the Unit interface.
func (p Civilian) Mass() float64 {
	const massPerUnit = 1.00 // per 100
	return p.Quantity() * massPerUnit
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
			deltaTech := q.techLevel - n.techLevel
			deltaRebels = q.qty.rebel * deltaTech / 100
		}
	}
	if deltaRebels < 1 {
		deltaRebels = 1
	}
	n.qty.loyal, n.qty.rebel = n.qty.loyal-deltaRebels, n.qty.rebel+deltaRebels

	return n
}

// NaturalBirthRate implements the PopulationGroup interface.
func (p Civilian) NaturalBirthRate(standardOfLiving, pctCapacity float64) float64 {
	return naturalBirthRate(p.techLevel, standardOfLiving, pctCapacity, p.IsOnShip(), p.IsResortColony())
}

// NaturalDeathRate implements the PopulationGroup interface.
func (p Civilian) NaturalDeathRate(standardOfLiving, pctCapacity float64) float64 {
	return naturalDeathRate(p.techLevel, standardOfLiving, pctCapacity)
}

// Population implements the PopulationGroup interface.
func (p Civilian) Population() int {
	return p.qty.loyal + p.qty.rebel
}

// Quantity implements the Unit interface.
func (p Civilian) Quantity() float64 {
	// there are 100 people per population unit
	return float64(p.Population()) * 0.01
}

// Rebels implements the PopulationGroup interface.
func (p Civilian) Rebels() int {
	return p.qty.rebel
}

// TechLevel implements the TechLevel interface.
func (p Civilian) TechLevel() int {
	return p.techLevel
}

// UnmarshalJSON implements the json.Unmarshaler interface
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

// Volume implements the Unit interface.
func (p Civilian) Volume() float64 {
	const volumePerUnit = 1.00 // per 100
	return p.Quantity() * volumePerUnit
}
