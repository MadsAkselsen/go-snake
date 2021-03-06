package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

// Extra challenges:
// 1. Make the game dynamically adjust when the screen is resized
// 2. Prettify the game frame

const SnakeSymbol = 0x2588
const AppleSymbol = 0x25CF
const GameFrameWidth = 30
const GameFrameHeight = 15
const GameFrameSymbol = '║'

type Point struct {
	row, col int
}

type Snake struct {
	parts []*Point
	velRow, velCol int
	symbol rune
}

type Apple struct {
	row, col int
	symbol rune
}

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var snake *GameObject
var isGamePaused bool
var debugLog string

func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	for {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()

		time.Sleep(75 * time.Millisecond)
	}

	screen.Fini()
}

func UpdateState() {
	if isGamePaused {
		return
	}

	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}
}

func DrawState() {
	if isGamePaused {
		return
	}

	screen.Clear()
	PrintString(0, 0, debugLog)
	PrintGameFrame()
	for _, obj := range gameObjects {
		PrintFilledRectInGameFrame(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}

	screen.Show()
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func InitGameState() {
	snake := &GameObject{
		row: 
	}
	gameObjects = []*GameObject{}
}

func HandleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	}
}

func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}

	return key
}

func PrintGameFrame() {
	frameStartRow, frameStartCol := GetFrameTopLeft()
	PrintUnfilledRect(
		frameStartRow, frameStartCol, GameFrameWidth, GameFrameHeight, GameFrameSymbol)
}

func GetFrameTopLeft() (int, int) {
	screenWidth, screenHeight := screen.Size()
	return screenHeight/2 - GameFrameHeight/2, screenWidth/2 - GameFrameWidth/2
}

func PrintStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		PrintFilledRect(row, col, 1, 1, c)
		col += 1
	}
}

func PrintFilledRectInGameFrame(row, col, width, height int, ch rune) {
	frameRow, frameCol := GetFrameTopLeft()
	PrintFilledRect(frameRow+row, frameCol+col, width, height, ch)
}

func PrintFilledRect(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func PrintUnfilledRect(row, col, width, height int, ch rune) {
	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row, ch, nil, tcell.StyleDefault)
	}

	for r := 0; r < height; r++ {
		screen.SetContent(col, row+r, ch, nil, tcell.StyleDefault)
		screen.SetContent(col+width-1, row+r, ch, nil, tcell.StyleDefault)
	}

	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row+height-1, ch, nil, tcell.StyleDefault)
	}
}