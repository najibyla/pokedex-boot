package main

import (
	"fmt"
	"math/rand"
)

func calcDamage(atk, def int) int {
	dmg := atk - def/2 + rand.Intn(10) + 1
	if dmg < 1 {
		return 1
	}
	return dmg
}

func clampPos(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

func simulateBattle(attacker, defender PartyMember) (won bool, xpGained int) {
	atkHP := attacker.stat("hp") * 2
	defHP := defender.stat("hp") * 2

	fmt.Printf("%s (Lv.%d, HP:%d) vs %s (Lv.%d, HP:%d)\n",
		attacker.Name, attacker.Level, atkHP,
		defender.Name, defender.Level, defHP)

	for round := 1; round <= 30 && atkHP > 0 && defHP > 0; round++ {
		fmt.Printf("  --- Round %d ---\n", round)
		if attacker.stat("speed") >= defender.stat("speed") {
			dmg := calcDamage(attacker.stat("attack"), defender.stat("defense"))
			defHP -= dmg
			fmt.Printf("  %s hits %s for %d! (%s HP: %d)\n", attacker.Name, defender.Name, dmg, defender.Name, clampPos(defHP))
			if defHP <= 0 {
				break
			}
			dmg = calcDamage(defender.stat("attack"), attacker.stat("defense"))
			atkHP -= dmg
			fmt.Printf("  %s hits %s for %d! (%s HP: %d)\n", defender.Name, attacker.Name, dmg, attacker.Name, clampPos(atkHP))
		} else {
			dmg := calcDamage(defender.stat("attack"), attacker.stat("defense"))
			atkHP -= dmg
			fmt.Printf("  %s hits %s for %d! (%s HP: %d)\n", defender.Name, attacker.Name, dmg, attacker.Name, clampPos(atkHP))
			if atkHP <= 0 {
				break
			}
			dmg = calcDamage(attacker.stat("attack"), defender.stat("defense"))
			defHP -= dmg
			fmt.Printf("  %s hits %s for %d! (%s HP: %d)\n", attacker.Name, defender.Name, dmg, defender.Name, clampPos(defHP))
		}
	}

	won = atkHP > 0
	if won {
		xpGained = defender.BaseExperience / 2
	}
	return
}
