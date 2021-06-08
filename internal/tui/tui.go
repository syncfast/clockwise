// Package tui creates the Terminal User Interface, calculates cost over time,
// and collects user input in manual mode.
package tui

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// refreshInterval determines the refresh frequency of the various goroutines.
const refreshInterval time.Duration = time.Millisecond * 500

// initScreen initializes, configures, and returns a tcell screen.
func initScreen() (tcell.Screen, error) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	s, err := tcell.NewScreen()
	if err != nil {
		return s, err
	}

	if err = s.Init(); err != nil {
		return s, err
	}

	s.SetStyle(tcell.StyleDefault)
	s.Clear()

	return s, nil
}

func Start(manual bool, data *Data) {
	// var data Data

	s, err := initScreen()
	if err != nil {
		log.Fatal(err)
	}

	// Start cost calculation goroutine.
	go func() {
		calculateCost(data)
	}()

	// Start cost file generation subroutine.
	go func() {
		writeCostFile(data)
	}()

	quit := make(chan struct{})

	// Start manual user input goroutine.
	go func() {
		for {
			ev := s.PollEvent()

			// Handling exit and resize separate from input so that we toggle manual input.
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlC:
					close(quit)
					return
				}
			case *tcell.EventResize:
				s.Sync()
				draw(s, data, manual)
			}

			if manual {
				switch ev := ev.(type) {
				case *tcell.EventKey:
					switch ev.Key() {
					case tcell.KeyEscape, tcell.KeyCtrlC:
						close(quit)
						return
					case tcell.KeyUp, tcell.KeyRight:
						c := data.GetCount()
						c++
						data.SetCount(c)
					case tcell.KeyDown, tcell.KeyLeft:
						c := data.GetCount()
						if c > 0 {
							c--
							data.SetCount(c)
						}
					case tcell.KeyRune:
						r := ev.Rune()
						if unicode.IsDigit(r) {
							i := data.getInput()
							i += string(r)
							data.setInput(i)
						}
					case tcell.KeyEnter:
						i := data.getInput()
						if i == "" {
							continue
						}
						c, _ := strconv.Atoi(i)
						data.SetCount(c)
						data.setInput("")
					case tcell.KeyBackspace, tcell.KeyBackspace2:
						i := data.getInput()
						if i != "" {
							data.setInput(i[0 : len(i)-1])
						}
					}

					// Render TUI after processing input.
					draw(s, data, manual)
				}
			}
		}
	}()

	tick(s, data, manual, quit)
	s.Fini()

	// log.Info("Clockwise has been terminated.")
}

// data stores variables passed around between the various goRoutines.
type Data struct {
	mtx   sync.Mutex
	count int
	cost  float32
	input string
}

// Get count.
func (data *Data) GetCount() int {
	data.mtx.Lock()
	defer data.mtx.Unlock()
	return data.count
}

// Set count.
func (data *Data) SetCount(count int) {
	data.mtx.Lock()
	defer data.mtx.Unlock()
	data.count = count
}

// Get cost.
func (data *Data) getCost() float32 {
	data.mtx.Lock()
	defer data.mtx.Unlock()
	return data.cost
}

// Set cost.
func (data *Data) setCost(cost float32) {
	data.mtx.Lock()
	defer data.mtx.Unlock()
	data.cost = cost
}

// Get input.
func (data *Data) getInput() string {
	data.mtx.Lock()
	defer data.mtx.Unlock()
	return data.input
}

// Set input.
func (data *Data) setInput(input string) {
	data.mtx.Lock()
	defer data.mtx.Unlock()
	data.input = input
}

// tick configures the goroutine for the scheduled calculateCost update.
func tick(s tcell.Screen, data *Data, manual bool, quit <-chan struct{}) {
	t := time.NewTicker(refreshInterval)
	for {
		select {
		case <-quit:
			return
		case <-t.C:
			draw(s, data, manual)
		}
	}
}

// calculateCost calculates the total cost.
func calculateCost(data *Data) {
	for {
		count := data.GetCount()
		total := data.getCost()
		fullyLoadedCostMultiplier := float32(1.75)
		cps := float32(count) * fullyLoadedCostMultiplier * float32(150000) / 7488000
		total += cps
		data.setCost(total)
		time.Sleep(refreshInterval)
	}
}

// Scaffolding for draw functionality.
func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

// draw renders the TUI.
func draw(s tcell.Screen, data *Data, manual bool) {
	s.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorCornflowerBlue)
	emitStr(s, 0, 0, style, "Clockwise")

	costString := fmt.Sprintf("Total cost: $%.2f", data.getCost())
	emitStr(s, 0, 1, tcell.StyleDefault, costString)

	countString := fmt.Sprintf("Participant count: %s", strconv.Itoa((data.GetCount())))
	emitStr(s, 0, 2, tcell.StyleDefault, countString)

	if manual {
		inputString := fmt.Sprintf("Input: %s", data.getInput())
		emitStr(s, 0, 3, tcell.StyleDefault, inputString)
	}

	s.Show()
}

// writeCostFile outputs the cost that gets consumed by OBS.
func writeCostFile(data *Data) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	outputFolder := home + "/Documents/clockwise/"
	_ = os.Mkdir(outputFolder, os.ModePerm)
	outputFile := outputFolder + "clockwise.txt"

	for {
		costString := fmt.Sprintf("Total cost: $%.2f\n", data.getCost())

		if err := ioutil.WriteFile(outputFile, []byte(costString), 0600); err != nil {
			log.Fatal(err)
		}

		time.Sleep(refreshInterval)
	}
}
