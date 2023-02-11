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

import "fmt"

// PopulationGroup defines the interface for working with groups of people.
type PopulationGroup interface {
	// FoodNeeded returns the number of FOOD units needed to sustain the population.
	FoodNeeded() float64
	// LifeSupportNeeded returns the number of LS units needed to sustain the population.
	LifeSupportNeeded() float64
	// NaturalBirthRate returns the percentage of natural births in the group.
	NaturalBirthRate() float64
	// NaturalDeathRate returns the percentage of natural deaths in the group.
	NaturalDeathRate() float64
	// Population returns total population of the unit.
	Population() int
	// Rebels returns the number of rebels in the population.
	Rebels() int
}

// naturalBirthRate calculates the birth rate for a population.
// The basic birth rate ranges from 0.25% to 10% of the population.
// The variation depends on the standard of living as well as the
// availability of "open" living space in the colony.
func naturalBirthRate(techLevel int, standardOfLiving, pctCapacity float64, isOnShip, isResortColony bool) float64 {
	if isOnShip { // births never happen on a ship
		return 0
	}
	// clamp the standard of living and percent capacity
	standardOfLiving = clamp(standardOfLiving, 0.01, 3.0)
	pctCapacity = clamp(pctCapacity, 0.01, 1.0)

	// the base rate is determined by tech level
	birthRate := float64(11-techLevel) * 0.1
	if birthRate < 0.0025 {
		birthRate = 0.0025
	} else if birthRate > 0.10 {
		birthRate = 0.10
	}

	// resort colonies increase the birth rate
	if isResortColony {
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

// naturalDeathRate calculates the basic death rate for a population.
// The rate is based on the tech level, standard of living, and
// availability of living space in the colony or ship.
func naturalDeathRate(techLevel int, standardOfLiving, pctCapacity float64) float64 {
	// clamp the standard of living and percent capacity
	standardOfLiving = clamp(standardOfLiving, 0.01, 3.0)
	pctCapacity = clamp(pctCapacity, 0.01, 1.0)

	// the base rate is determined by tech level
	var deathRate float64
	switch techLevel {
	case 0:
		deathRate = 1_500.0 / 100_000.0
	case 1:
		deathRate = 1_400.0 / 100_000.0
	case 2:
		deathRate = 1_300.0 / 100_000.0
	case 3:
		deathRate = 1_200.0 / 100_000.0
	case 4:
		deathRate = 1_100.0 / 100_000.0
	case 5:
		deathRate = 1_000.0 / 100_000.0
	case 6:
		deathRate = 900.0 / 100_000.0
	case 7:
		deathRate = 800.0 / 100_000.0
	case 8:
		deathRate = 700.0 / 100_000.0
	case 9:
		deathRate = 600.0 / 100_000.0
	case 10:
		deathRate = 500.0 / 100_000.0
	default:
		panic(fmt.Sprintf("assert(0 <= %d <= 10)", techLevel))
	}

	// standard of living influences it
	if standardOfLiving > 1.500 {
		deathRate *= 0.975
	} else if standardOfLiving > 1.250 {
		deathRate *= 0.950
	} else if standardOfLiving > 0.990 {
		// base rate
	} else if standardOfLiving > 0.875 {
		deathRate *= 1.025
	} else if standardOfLiving > 0.750 {
		deathRate *= 1.050
	} else if standardOfLiving > 0.625 {
		deathRate *= 1.075
	} else if standardOfLiving > 0.500 {
		deathRate *= 1.100
	} else if standardOfLiving > 0.375 {
		deathRate *= 1.125
	} else if standardOfLiving > 0.250 {
		deathRate *= 1.150
	} else if standardOfLiving > 0.125 {
		deathRate *= 1.175
	}

	// overcrowding increases it
	if pctCapacity > 2.000 {
		deathRate *= 3.000
	} else if pctCapacity > 1.500 {
		deathRate *= 2.000
	} else if pctCapacity > 0.990 {
		deathRate *= 1.500
	} else if pctCapacity > 0.975 {
		deathRate *= 1.250
	} else if pctCapacity > 0.950 {
		deathRate *= 1.100
	} else if pctCapacity > 0.925 {
		deathRate *= 1.025
	} else if pctCapacity > 0.900 {
		deathRate *= 1.010
	} else {
		// base rate
	}

	// death rate is never less than 0.25% or higher than 75%
	return clamp(deathRate, 0.00_2500, 0.75_0000)
}
