package main

import (
	"strconv"
	"strings"
	"time"
	"flag"
	"fmt"
	"math/rand"
)

var (
	lane      [][]string
	numLaps   int
	numRacers int
	winners   []int

	//Rune Idea from: // https://stackoverflow.com/questions/16682797/how-to-convert-a-rune-to-unicode-style-string-like-u554a-in-golang
	r1 = rune('🚗')
	r2 = rune('🚓')
	r3 = rune('🚕')
	r4 = rune('🛺')
	r5 = rune('🚌')
	r6 = rune('🛵')
	r7 = rune('🚒')
	r8 = rune('🚚')
	r9 = rune('🛸')
	r10 = rune('🚜')

	racersNames = [11]string{"-", "Automovilista " + string(r1),"Policia " + string(r2), "Taxista " + string(r3), "Mototaxista " + string(r4), "Microbusero " + string(r5), 
							"Motociclista " + string(r6), "Bombero " + string(r7), "Trailero " + string(r8),"Marcianito " + string(r9), "Agronomo " + string(r10)}
	racersRunes = [11]string{"-", string(r1), string(r2), string(r3), string(r4), string(r5), string(r6), string(r7), string(r8), string(r9), string(r10)}
	racers = make(map[int]chan bool) 
	//Through this channel the racers will ask to move
	moveRequest = make(chan Data)
	//Channel used to clear the previusly position of a racer
	clearRequest = make(chan Data, 60)
	//Channel used to communicate racers data throughout the race
	updateRequest  = make(chan PrintData, 60)  	
)

//How long the race track is
const circuitDistance = 170

type PrintData struct {
	name       int
	lane       int
	position   int
	lap        int
	speed      float64
	lapTime    string
}

type Data struct {
	name       int
	lane       int
	position   int
	currentLap int
}

func printLanes() {
	for i := range lane {
		fmt.Print("|")
		tmp := ""
		for j := range lane[i] {
			tmp += lane[i][j]
		}
		fmt.Print(tmp)
		fmt.Print("|")
		fmt.Println("")
		tmp = "|" + strings.Repeat("═", circuitDistance) + "|"
		fmt.Println(tmp)
	}
}

// Table format idea from: https://github.com/jedib0t/go-pretty
func spacedPrint(string2print string, numSpaces int) string {
	if len(string2print) < numSpaces {
		string2print += strings.Repeat(" ", (numSpaces - len(string2print)))
	}
	return string2print
}

func printRacersData(killPrint chan struct{}) {
	tmp := ""
	start := time.Now()
	updateList := make([]PrintData, numRacers)
	const numSpaces = 23
	const numSpacesNames = 25
	info := [6]string{"", "Lane: ", "Position: km ", "Lap: ", "Speed: ", "Lap Time: "}
	for {
		for i := 0; i < 15; i++ {
			select {
			case x := <-updateRequest:
				updateList[x.name-1] = x
			case <-killPrint:
				return
			}
		}

		//Clears the terminal
		fmt.Print("\033[H\033[2J")
		fmt.Println("")
		printLanes()
		
		fmt.Println("")
		tmp = strings.Repeat("━━━━━━━━━━━━━━━━━━━━━━━━╋", numRacers)
		fmt.Println(tmp)
		params := make([]string, 6)
		for j := 0; j < numRacers; j++ {
			auxStr := info[0] + racersNames[updateList[j].name]
			//Names
			params[0] += " " + spacedPrint(auxStr, numSpacesNames) + "╋"                                   
			//Lanes
			params[1] += " " + spacedPrint(info[1]+strconv.Itoa(updateList[j].lane), numSpaces) + "╋"       
			//Positions
			params[2] += " " + spacedPrint(info[2]+strconv.Itoa(updateList[j].position), numSpaces) + "╋"   
			//Current Lap
			if updateList[j].lap == numLaps+1 {
				auxStr = info[3] + "Finish"
			} else {
				auxStr = info[3] + strconv.Itoa(updateList[j].lap) + "/" + strconv.Itoa(numLaps)
			}
			params[3] += " " + spacedPrint(auxStr, numSpaces) + "╋"  
			//Current Speed
			speedString := fmt.Sprintf("%.2f" + " km/h", updateList[j].speed)
			params[4] += " " + spacedPrint(info[4]+speedString, numSpaces) + "╋"   
			//Lap Time                       
			params[5] += " " + spacedPrint(info[5]+updateList[j].lapTime, numSpaces) + "╋"      
		}
		for _, v := range params {
			fmt.Println(v)
		}

		tmp = strings.Repeat("━━━━━━━━━━━━━━━━━━━━━━━━╋", numRacers)
		fmt.Println(tmp)
		fmt.Println("Total Time: ", time.Now().Sub(start).String())
	}

}

func racerBehaviour(initLocation Data, maxSpeed float64, acceleration float64, moveChannelRequest chan Data, response chan bool) {
	startLap := time.Now()
	endLap := time.Now().Sub(startLap)
	name := initLocation.name
	currentLocation := initLocation
	currentSpeed := acceleration
	slowValue := -80.0
	sleep := 400.0
	
	nextLocation := Data{0, 0, 0, 0}
	nextAcceleration := 0.0
	lap := initLocation.currentLap

	//If there are unfinished laps
	for lap < numLaps+1 {
		//Sleep time for each racer, allows to synchronize the behaviour
		time.Sleep(time.Duration(sleep-currentSpeed) * time.Millisecond)
		for {
			frontCar := false
			//Check if there are cars 6 positions after him.
			for i := (currentLocation.position + 1) % circuitDistance; i != (currentLocation.position + 6) % circuitDistance; i = (i + 1) % circuitDistance {
				//If so, check that it is in the same lane
				if lane[currentLocation.lane][i] != " " { 
					frontCar = true

					//We check if the car can switch to an adjacent lane

					//Top Lane
					if currentLocation.lane == 0 {
						//If the right lane is clear, then switch
						if lane[currentLocation.lane+1][(currentLocation.position+1)%circuitDistance] == " " {
							nextLocation = Data{name, currentLocation.lane + 1, currentLocation.position + 1, lap}
							nextAcceleration = acceleration
						
						//If not, slow down
						} else {
							nextLocation = Data{name, currentLocation.lane, currentLocation.position + 1, lap}
							nextAcceleration = slowValue
						}
					//Bottom Lane	
					} else if currentLocation.lane == 4 {
						//If the left lane is clear, then switch
						if lane[currentLocation.lane-1][(currentLocation.position+1)%circuitDistance] == " " {
							nextLocation = Data{name, currentLocation.lane - 1, currentLocation.position + 1, lap}
							nextAcceleration = acceleration

						//If not, slow down
						} else {
							nextLocation = Data{name, currentLocation.lane, currentLocation.position + 1, lap}
							nextAcceleration = slowValue
						}
					//Middle Lanes	
					} else {
						//If the right lane is clear, then switch
						if lane[currentLocation.lane+1][(currentLocation.position+1)%circuitDistance] == " " {
							nextLocation = Data{name, currentLocation.lane + 1, currentLocation.position + 1, lap}
							nextAcceleration = acceleration
						//If the left lane is clear, then switch
						} else if lane[currentLocation.lane-1][(currentLocation.position+1)%circuitDistance] == " " {
							nextLocation = Data{name, currentLocation.lane - 1, currentLocation.position + 1, lap}
							nextAcceleration = acceleration
						//If not, slow down			
						} else {
							nextLocation = Data{name, currentLocation.lane, currentLocation.position + 1, lap}
							nextAcceleration = slowValue
						}
					}
				}
				//If we already find a car in the path, then stop searching
				if frontCar {
					break
				}
			}

			//If the lane is clear, then just continue
			if frontCar == false { 
				nextLocation = Data{name, currentLocation.lane, currentLocation.position + 1, lap}
				nextAcceleration = acceleration
			}
			
			//If the racer finish the lap then we relocate the racer at the start of the lane
			if nextLocation.position >= circuitDistance {
				nextLocation.position = 0
			}
			//Request to move to the next position
			moveChannelRequest <- nextLocation

			//If the path is clear then we clear the previus position
			if <-response == true {
				clearRequest <- currentLocation
				break
			}
		}
		//We adjust the data at the start of any lap
		if nextLocation.position == 0 {
			endLap = time.Now().Sub(startLap)
			startLap = time.Now()
			lap++
		}
		currentLocation = nextLocation
		if nextAcceleration > 0 {
			if currentSpeed < maxSpeed {
				currentSpeed += nextAcceleration
			} else {
				currentSpeed = maxSpeed
			}

		} else {
			if currentSpeed+nextAcceleration < 0 {
				currentSpeed = 0
			} else {
				currentSpeed += nextAcceleration
			}
		}
		//Request to update the data of the racer
		updateRequest <- PrintData{name, currentLocation.lane, currentLocation.position, lap, currentSpeed, endLap.String()}
	}
}

func main() {
	winners = []int{1, 2, 3}
	winners = winners[:0]
	lane = make([][]string, 5)
	racers := make(map[int]chan bool)
	startPos := [10]int{0, 0, 0, 0, 0, 3, 3, 3, 3, 3}

	nr := flag.Int("racers", 5, "N° of Racers")
	nl := flag.Int("laps", 1, "N° of laps")
	flag.Parse()
	numRacers = *nr
	numLaps = *nl
	//validate inputs
	if numRacers > 10 || numRacers < 2 {
		fmt.Println("Racers value must be between 2 and 10")
		return
	}

	if numLaps < 1 {
		fmt.Println("Laps value must be greater than 1")
		return
	}

	//Initialize empty lanes array
	for i := range lane {
		lane[i] = make([]string, circuitDistance)
	}
	for i := range lane {
		for j := range lane[i] {
			lane[i][j] = " "
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	//Initialize racers
	for i := 1; i < numRacers+1; i++ {
		//This channel will allow the communication from main to racerBehaviour 
		auxResponseChan := make(chan bool)
		racers[i] = auxResponseChan
		//Random speed and acceleration
		tmpMaxSpeed := float64(r.Intn(220-125) + 125)
		tmpAcceleration := float64(r.Intn(45-15) + 15)
		go racerBehaviour(Data{i, i % 5, startPos[i-1], 1}, tmpMaxSpeed, tmpAcceleration, moveRequest, auxResponseChan)
	}
	printChan := make(chan struct{})
	go printRacersData(printChan)
	for {
		select {
		//Received requests to check if the racer can move to the specified position
		case recievedRequest := <-moveRequest:
			//If its empty then we move the racer to that position
			if lane[recievedRequest.lane][recievedRequest.position] == " " {
				lane[recievedRequest.lane][recievedRequest.position] = racersRunes[recievedRequest.name]
				racers[recievedRequest.name] <- true
				//We save the names of the racers who finish the race 
				if recievedRequest.currentLap == numLaps && recievedRequest.position == 0 { 
					//When a racer finish all the laps we delete it from the circuit to avoid incorrect behaviours 
					lane[recievedRequest.lane][recievedRequest.position] = " "
					winners = append(winners, recievedRequest.name)
					//When everyone finish the race then we show the results
					if len(winners) == numRacers {
						//Stop printing
						printChan <- struct{}{}

						//Clears the terminal
						fmt.Print("\033[H\033[2J")

						fmt.Println("        Race is over!")
						fmt.Println("╋━━━━━━━━━━Results━━━━━━━━━━╋")
						
						for i := 0; i < numRacers; i++ {
							fmt.Println("╋ " + strconv.Itoa(i+1) + ".-" + racersNames[winners[i]])
							fmt.Println("╋━━━━━━━━━━━━━━━━━━━━━━━━━━━╋")
						}
						return
					}
				}
			} else {
				racers[recievedRequest.name] <- false
			}
		//Request to clear the space previusly used by a racer
		case recievedRequest := <-clearRequest:
			if lane[recievedRequest.lane][recievedRequest.position] == racersRunes[recievedRequest.name] {
				lane[recievedRequest.lane][recievedRequest.position] = " "
			} 
		}

	}

}