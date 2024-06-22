package main

import (
	"flag"
	"fmt"
	"os"
	"unsafe"
)

/*
#include <stdio.h>
#include "custom_string.c"
#include "string_array.c"
#include "debug.c"
*/
import "C"

const (
	DEFAULT_PRINT_MEMORY_BUDGET = (1 << 15)
	MAX_FILE_SIZE               = 1 << 16
)

var PRINT_MEMORY_BUDGET = DEFAULT_PRINT_MEMORY_BUDGET

func SetFlags() {
	flag.BoolVar(&VERBOSE, "v", false, "")
	flag.BoolVar(&VERBOSE, "verbose", false, "Show more logs")
	flag.IntVar(&PRINT_MEMORY_BUDGET, "m", DEFAULT_PRINT_MEMORY_BUDGET, "")
	flag.IntVar(&PRINT_MEMORY_BUDGET, "memory", DEFAULT_PRINT_MEMORY_BUDGET, "Memory budget")
	flag.Parse()
	if !VERBOSE {
		PRINT_MEMORY_BUDGET = 0
	}
}

func GetFileSize(filepath string) (int64, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func main() {
	var err error
	VPrintln("start main in player.go")
	SetFlags()
	VPrintln("flags set")
	tail := flag.Args()
	if len(tail) != 2 {
		fmt.Println("Error: Expected 2 lua files in arguments")
		return
	}

	fileName1 := tail[0]
	fileName2 := tail[1]
	fileSize1, err := GetFileSize(fileName1)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	} else if fileSize1 > MAX_FILE_SIZE {
		fmt.Printf("max fiel size exceeded: %v/%v\n", fileSize1, MAX_FILE_SIZE)
		return
	}
	fileSize2, err := GetFileSize(fileName2)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	} else if fileSize2 > MAX_FILE_SIZE {
		fmt.Printf("max fiel size exceeded: %v/%v\n", fileSize2, MAX_FILE_SIZE)
		return
	}

	b1, err := os.ReadFile(fileName1)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	b2, err := os.ReadFile(fileName2)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	script1 := string(b1)
	script2 := string(b2)

	C.PRINT_MEMORY_BUDGET = C.size_t(PRINT_MEMORY_BUDGET)

	pl1 := NewState()
	PushFunction(pl1, unsafe.Pointer(C.printInLua), "print")
	err = InitRG(pl1, script1, fileName1)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	} else {
		VPrintln("[Successfully initialized]\n")
	}
	pl2 := NewState()
	PushFunction(pl2, unsafe.Pointer(C.printInLua), "print")
	err = InitRG(pl2, script2, fileName2)

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	} else {
		VPrintln("[Successfully initialized]\n")
	}

	game, err := PlayGame(pl1, pl2)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	VPrintf("game:\n%v\n", game)
	score1 := 0
	score2 := 0
	for _, botState := range game[len(game)-1] {
		if botState.bot.playerId == 1 {
			score1 += 1
		} else {
			score2 += 1
		}
	}
	fmt.Printf("%v - %v\n", score1, score2)

	VPrintln("close state")
	CloseState(pl1)
	CloseState(pl2)

	VPrintln("end main in go")
}
