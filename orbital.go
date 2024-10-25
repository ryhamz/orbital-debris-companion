package main

import (
	"fmt"
	"math/rand"
	"time"
)

// CrewMember represents a single player in the game
type CrewMember struct {
	Name   string
	Action string
	Stress int
	Rested bool
}

// Hex represents a hex on the map
type Hex struct {
	ID           int
	IsExplored   bool
	IsDangerous  bool
	HasArtifact  bool
	HiddenItems  bool
	IsSafeToRest bool
}

// Watch represents a 4-hour cycle in the game
type Watch struct {
	Number int
	Event  int
}

// HexMap represents the map with multiple hexes
type HexMap struct {
	Hexes      []Hex
	CurrentHex int
}

// Warden represents the game master handling the procedure
type Warden struct {
	EventDie int
	Time     int // Tracks how many watches have passed (6 watches = 1 day)
}

// Actions for the crew members
const (
	Trek   = "Trek"
	Guard  = "Guard"
	Encamp = "Encamp"
	Rest   = "Rest"
	Explore = "Explore"
	Collect = "Collect"
)

// InitializeCrew initializes a group of crew members
func InitializeCrew(names []string) []CrewMember {
	var crew []CrewMember
	for _, name := range names {
		crew = append(crew, CrewMember{Name: name, Stress: 0, Rested: true})
	}
	return crew
}

// RollEventDie simulates rolling a d10 to determine an event
func (w *Warden) RollEventDie() int {
	w.EventDie = rand.Intn(10) + 1
	return w.EventDie
}

// Trek action allows the crew to move to an adjacent hex
func TrekAction(hexMap *HexMap, warden *Warden) {
	fmt.Println("Crew treks to an adjacent hex.")
	// Move to the next hex on the map (wraps around if at the end)
	hexMap.CurrentHex = (hexMap.CurrentHex + 1) % len(hexMap.Hexes)
	warden.RollEventDie()
	fmt.Printf("Event die rolled: %d\n", warden.EventDie)
	// Apply any effects based on the current hex
	if hexMap.Hexes[hexMap.CurrentHex].IsDangerous {
		fmt.Println("The new hex is dangerous! Rolling for events...")
		// Roll 2 event dice if the crew is in danger
		warden.RollEventDie()
	}
}

// Guard action gives advanced warning for any dangers
func GuardAction() {
	fmt.Println("Crew member is guarding, they will have advanced warning of any danger.")
}

// Encamp action allows the crew to rest safely
func EncampAction(crew []CrewMember, hexMap *HexMap) {
	if hexMap.Hexes[hexMap.CurrentHex].IsSafeToRest {
		fmt.Println("The crew finds a safe place to rest.")
		for i := range crew {
			crew[i].Rested = true
			crew[i].Stress = 0
		}
	} else {
		fmt.Println("No safe place to rest. The crew cannot reduce stress.")
	}
}

// Rest action attempts to reduce stress and recover
func RestAction(crewMember *CrewMember) {
	fmt.Printf("%s is resting.\n", crewMember.Name)
	if crewMember.Rested {
		crewMember.Stress = 0
		fmt.Println("Stress reduced!")
	} else {
		fmt.Println("No safe rest, stress remains.")
	}
}

// Explore action discovers major hex features or hidden items
func ExploreAction(hexMap *HexMap) {
	currentHex := &hexMap.Hexes[hexMap.CurrentHex]
	fmt.Println("The crew explores the hex.")
	if !currentHex.IsExplored {
		fmt.Println("They discover major features of the hex.")
		currentHex.IsExplored = true
		if currentHex.HiddenItems {
			fmt.Println("Hidden items are found!")
		}
	} else {
		fmt.Println("No new discoveries.")
	}
}

// Collect action allows the crew to gather resources
func CollectAction(crewMember *CrewMember, hexMap *HexMap) {
	currentHex := hexMap.Hexes[hexMap.CurrentHex]
	fmt.Printf("%s is collecting resources.\n", crewMember.Name)
	roll := rand.Intn(20) + 1 // Simulate skill check
	if roll == 20 {
		fmt.Println("Critical success! An artifact is discovered.")
		currentHex.HasArtifact = true
	} else if roll >= 15 {
		fmt.Println("Success! Useful goods are found.")
	} else if roll >= 10 {
		fmt.Println("Partial success. Some scrap is found.")
	} else {
		fmt.Println("Failure. Nothing of value is found.")
		if roll == 1 {
			fmt.Println("Critical failure! The crew member is lost until another crew member explores.")
		}
	}
}

// TimePasses updates the time and checks for exhaustion
func TimePasses(warden *Warden, crew []CrewMember) {
	warden.Time++
	fmt.Printf("Time passes. Watch %d is over.\n", warden.Time)
	if warden.Time%6 == 0 {
		fmt.Println("End of the day. Checking for exhaustion penalties...")
		for i := range crew {
			if !crew[i].Rested {
				fmt.Printf("%s suffers exhaustion penalties!\n", crew[i].Name)
				crew[i].Stress += 10
			}
			crew[i].Rested = false // Reset rest status at the end of the day
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize the crew
	crew := InitializeCrew([]string{"Alice", "Bob", "Charlie"})

	// Initialize the map (10 hexes for this example)
	hexMap := HexMap{
		Hexes: []Hex{
			{ID: 1, IsDangerous: false, IsSafeToRest: true, HiddenItems: true},
			{ID: 2, IsDangerous: true, IsSafeToRest: false},
			{ID: 3, IsDangerous: false, IsSafeToRest: true},
			// Add more hexes as needed
		},
		CurrentHex: 0,
	}

	// Initialize the Warden
	warden := Warden{}

	// Simulate a few watches
	for watch := 1; watch <= 4; watch++ {
		fmt.Printf("\n--- Watch %d ---\n", watch)
		
		// Crew declares actions
		crew[0].Action = Trek   // Example: First crew member treks
		crew[1].Action = Guard  // Second guards
		crew[2].Action = Rest   // Third rests

		// Resolve each crew member's action
		for i := range crew {
			switch crew[i].Action {
			case Trek:
				TrekAction(&hexMap, &warden)
			case Guard:
				GuardAction()
			case Rest:
				RestAction(&crew[i])
			case Explore:
				ExploreAction(&hexMap)
			case Collect:
				CollectAction(&crew[i], &hexMap)
			}
		}

		// Time passes, check for exhaustion
		TimePasses(&warden, crew)
	}
}
