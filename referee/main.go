package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ClementTariel/rg-lua/rgcore"
)

const (
	DEFAULT_PRINT_MEMORY_BUDGET = (1 << 15)
	MAX_FILE_SIZE               = 1 << 16
)

var PRINT_MEMORY_BUDGET = DEFAULT_PRINT_MEMORY_BUDGET

func SetFlags() {
	flag.BoolVar(&rgcore.VERBOSE, "v", false, "")
	flag.BoolVar(&rgcore.VERBOSE, "verbose", false, "Show more logs")
	flag.IntVar(&PRINT_MEMORY_BUDGET, "m", DEFAULT_PRINT_MEMORY_BUDGET, "")
	flag.IntVar(&PRINT_MEMORY_BUDGET, "memory", DEFAULT_PRINT_MEMORY_BUDGET, "Memory budget")
	flag.Parse()
	if !rgcore.VERBOSE {
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

func callPost(url string, postBody []byte) (*http.Response, error) {
	responseBody := bytes.NewBuffer(postBody)
	rgcore.VPrintln("call POST")
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	return resp, err
}

func callGet(url string) (*http.Response, error) {
	rgcore.VPrintln("call GET")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	return resp, err
}

// TODO: run game on api call not on start
func main() {
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	// e.Logger.Fatal(e.Start(":3333"))

	var err error
	rgcore.VPrintln("start main in referee.go")
	SetFlags()
	rgcore.VPrintln("flags set")
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
		fmt.Printf("max file size exceeded: %v/%v\n", fileSize1, MAX_FILE_SIZE)
		return
	}
	fileSize2, err := GetFileSize(fileName2)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	} else if fileSize2 > MAX_FILE_SIZE {
		fmt.Printf("max file size exceeded: %v/%v\n", fileSize2, MAX_FILE_SIZE)
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

	rgcore.SetPrintMemoryBudget(PRINT_MEMORY_BUDGET)

	rgcore.VPrintln("files read")
	game, err := PlayGame(fileName1, script1, fileName2, script2)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	rgcore.VPrintf("game:\n%v\n", game)
	score1 := 0
	score2 := 0
	for _, botState := range game[len(game)-1] {
		if botState.Bot.PlayerId == 1 {
			score1 += 1
		} else {
			score2 += 1
		}
	}
	fmt.Printf("%v - %v\n", score1, score2)

	_, err = callGet(fmt.Sprintf("http://localhost:%d/kill", rgcore.PORT_PLAYER_1))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	_, err = callGet(fmt.Sprintf("http://localhost:%d/kill", rgcore.PORT_PLAYER_2))
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	rgcore.VPrintln("end main in go")
}
