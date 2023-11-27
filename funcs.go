package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/robertkrimen/otto"
)

func getAverage(ao uint8, solves []string) (result string) {
	var number float64 = 0
	switch ao {
	case 5:
		if len(solves) >= 5 {
			for i := 0; i < 5; i++ {
				a, _ := strconv.ParseFloat(solves[i], 64)
				number += a
				if i == 4 {
					result = "Ao5 " + strconv.FormatFloat(number/5, 'f', 2, 64)
					break
				}
			}
		} else {
			result = "Ao5  -"
		}

	case 12:
		if len(solves) >= 12 {
			for i := 0; i < 12; i++ {
				a, _ := strconv.ParseFloat(solves[i], 64)
				number += a
				if i == 11 {
					result = "Ao12 " + strconv.FormatFloat(number/12, 'f', 2, 64)
					break
				}
			}
		} else {
			result = "Ao12  -"
		}

	case 50:
		if len(solves) >= 50 {
			for i := 0; i < 50; i++ {
				a, _ := strconv.ParseFloat(solves[i], 64)
				number += a
				if i == 49 {
					result = "Ao50" + strconv.FormatFloat(number/50, 'f', 2, 64)
					break
				}
			}
		} else {
			result = "Ao50  -"
		}
	case 100:
		if len(solves) >= 100 {
			for i := 0; i < 100; i++ {
				a, _ := strconv.ParseFloat(solves[i], 64)
				number += a
				if i == 99 {
					result = "Ao100" + strconv.FormatFloat(number/100, 'f', 2, 64)
					break
				}
			}
		} else {
			result = "Ao100  -"
		}
	}
	return result
}

func startTimer(f bool) {
	var a uint8 = 0
	switch f {
	case true:
		go func() {
			timerRunning = true
			fmt.Printf("Timer started... ")
			for range time.Tick(10 * time.Millisecond) {
				select {
				case <-ch:
					return
				default:
					a++
					seconds += 0.01
					t := strconv.FormatFloat(seconds, 'f', 2, 64)
					timer.SetText(t)
				}
			}
		}()

	default:
		ch <- true
		timerRunning = false
		solve := seconds
		seconds = 0
		t := strconv.FormatFloat(solve, 'f', 2, 64)
		timer.SetText(t)
		data, _ := os.ReadFile(timesPath)
		fmtdata := string(data)
		if fmtdata == "" {
			fmt.Fprintf(timesFile, "%s", t)
			fmt.Printf("Saved %s!\n", t)
		} else {
			fmt.Fprintf(timesFile, "\n%s", t)
			fmt.Printf("Saved %s!\n", t)
		}
	}
}
func manageFile() []string {
	data, _ := os.ReadFile(timesPath)
	fmtdata := string(data)
	content := strings.Split(fmtdata, "\n")
	slices.Reverse(content)
	return content
}

func getScramble() string {
	vm := otto.New()
	vm.Run(`Array.prototype.choose = function() {
		var index = Math.floor(Math.random() * this.length);
		return this[index];
};

function scramble(length) {
		// generate map of planes associated with each side
		var planes = {x: ['L', 'R'], y: ['U', 'D'], z: ['F', 'B']};
		var planeMap = {};
		for (var plane in planes) {
			var sides = planes[plane];
			for (var i = 0; i < sides.length; i++) {
				var side = sides[i];
				planeMap[side] = plane;
			}
		}

		var sides = ['F', 'B', 'R', 'L', 'U', 'D'];
		var modifiers = ['2', '\'', ''];

		// create buffer of moved sides
		// once plane crossed, move on those sides is then permittalbe again

		var weakBuffer = [], moves = [];
		for (var i = 0; i < length; i++) {
			var mod = modifiers.choose(), side;
			if (weakBuffer.length == 0) {
				side = sides.choose();
			}
			else if (weakBuffer.length == 1) {
				var badSide = weakBuffer[0],
				newSides = sides.slice(),
				badIndex = newSides.indexOf(badSide);
				newSides.splice(badIndex, 1);

				side = newSides.choose();

				if (planeMap[side] != planeMap[badSide]) {
					weakBuffer = []; // planes have been crossed
				}
			}
			else {
				var newSides = sides.slice();
				for (var a = 0; a < weakBuffer.length; a++) {
					var badSide = weakBuffer[a],
					badIndex = newSides.indexOf(badSide);
					newSides.splice(badIndex, 1);
				}

				side = newSides.choose();

				weakBuffer = [];
			}
			moves.push(side + mod);
			weakBuffer.push(side);
		}
		return moves.join(' ');
	}
var scramble = scramble(20);
		`)
	value, _ := vm.Get("scramble")
	text, _ := value.ToString()
	return text
}

func gen_avg_tab(num uint8) fyne.CanvasObject {
	switch num {
	case 5:
		widget := widget.NewLabel(getAverage(num, manageFile()))
		widget.TextStyle.Bold = true
		widget.Alignment = fyne.TextAlignTrailing
		return widget
	case 12:
		widget := widget.NewLabel(getAverage(num, manageFile()))
		widget.TextStyle.Bold = true
		widget.Alignment = fyne.TextAlignTrailing
		return widget
	case 50:
		widget := widget.NewLabel(getAverage(num, manageFile()))
		widget.TextStyle.Bold = true
		widget.Alignment = fyne.TextAlignTrailing
		return widget
	case 100:
		widget := widget.NewLabel(getAverage(num, manageFile()))
		widget.TextStyle.Bold = true
		widget.Alignment = fyne.TextAlignTrailing
		return widget
	}
	var widget fyne.CanvasObject
	return widget
}
