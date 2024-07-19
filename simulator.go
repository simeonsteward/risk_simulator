package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run simulator.go [yourArmy#] [territory1Army#] ...")
		return
	}
	ourStartingArmies, err := strconv.Atoi(os.Args[1])
	if err != nil {
			fmt.Printf("Non Integer numArmies: %s", os.Args[1])
			os.Exit(1)
		}
	// -2 for our armies and the executable
	enemyStartingArmies := make([]int, len(os.Args)-2)
	for i, numArmies := range(os.Args[2:]){
		enemyStartingArmies[i], err = strconv.Atoi(numArmies)
		if err != nil {
			fmt.Printf("Non Integer numArmies: %s\n", numArmies)
			os.Exit(1)
		}
	}
	fmt.Printf("Your starting armies: %d\n", ourStartingArmies)
	fmt.Println("Territories you want to conquer have")
	for i, armies := range(enemyStartingArmies){
		fmt.Printf("Territory %d has %d armies\n", i, armies)
	}

	startSimulation(10000, ourStartingArmies, enemyStartingArmies)
}

type battleResult struct {
	win int //1 if win 0 if lose
	started int // 1 if started 0 if not started
	unitsStart int
	unitsEnd int
}

func startSimulation(n int, units int, enemyUnitsByTerritory []int) {
	totalBattleResults := make([]battleResult, len(enemyUnitsByTerritory))
	for range n {
		armies := units
		for i := range enemyUnitsByTerritory{
			result := performBattle(armies, enemyUnitsByTerritory[i])

			prev := totalBattleResults[i]
			totalBattleResults[i] = battleResult{
				win: prev.win + result.win,
				started: prev.started + result.started,
				unitsStart: prev.unitsStart + result.unitsStart,
				unitsEnd: prev.unitsEnd + result.unitsEnd,
			}
			// subtract one to account for leaving an army in the territory
			armies = result.unitsEnd - 1
		}
	}

	for i, total := range totalBattleResults{
		fmt.Printf("-- Battle %d results --\n", i)
		startedRate := float32(total.started)/float32(n)
		fmt.Printf("Started Rate (%%): %.2f\n", startedRate)
		winRate := float32(total.win)/float32(n)
		fmt.Printf("Win Rate (%%): %.2f\n", winRate)
		avgUnitsStart := float32(total.unitsStart)/float32(n)
		fmt.Printf("avg units start: %.2f\n", avgUnitsStart)
		avgUnitsEnd := float32(total.unitsEnd)/float32(n)
		fmt.Printf("avg units end: %.2f\n", avgUnitsEnd)
	}
	
}

func min(a int, b int) int{
	if a < b {
		return a
	}
	return b
}

func performBattle(attackerArmy int, defenderArmy int) battleResult{
	result := battleResult{
		unitsStart: attackerArmy,
	}
	if attackerArmy <= 1{
		result.started = 0
		result.win = 0
		result.unitsEnd = attackerArmy
		return result
	}
	result.started = 1
	for attackerArmy > 1 && defenderArmy > 0 {
		attackersInPlay := min(attackerArmy - 1, 3)
		defendersInPlay := min(defenderArmy, 2)
		attackerRolls := make([]int, 0, attackersInPlay)
		defenderRolls := make([]int, 0, defendersInPlay)
		for range attackersInPlay {
			attackerRolls = append(attackerRolls, rand.Intn(6) + 1)
		}
		for range defendersInPlay {
			defenderRolls = append(defenderRolls, rand.Intn(6) + 1)
		}

		attackerCasualties, defenderCasualties := calcCasualties(attackerRolls, defenderRolls)
		attackerArmy += attackerCasualties
		defenderArmy += defenderCasualties

	}
	result.unitsEnd = attackerArmy
	if attackerArmy > 1 {
		result.win = 1
	} else{
		result.win = 0
	} 
	return result
	}

func calcCasualties(attackerRolls []int, defenderRolls []int) (int, int) {
	defenderCasualties := 0
	attackerCasualties := 0
	sort.Sort(sort.Reverse(sort.IntSlice(attackerRolls)))
	sort.Sort(sort.Reverse(sort.IntSlice(defenderRolls)))
	for i := range min(len(attackerRolls), len(defenderRolls)) {
		if attackerRolls[i] > defenderRolls[i] {
			defenderCasualties--
		} else {
			attackerCasualties--
		}
	}
	return attackerCasualties, defenderCasualties
}