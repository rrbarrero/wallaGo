package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	logging "github.com/op/go-logging"
	"github.com/tkanos/gonfig"
)

var config Configuration

var log = logging.MustGetLogger("wallago")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func check(e error) {
	if e != nil {
		log.Error(e)
		panic(e)
	}
}

func check_item(wg *sync.WaitGroup, line string) {
	defer wg.Done()
	if len(line) == 0 {
		return
	}
	if line[0] == '#' {
		return
	}
	item := NewItem(line)
	for i, e := range item.CheckItem() {
		fmt.Printf("%v.- %v - %v - %v - %v\n", i, e.SalePrice, e.Title, WALLA_ITEM_URL+e.URL, e.ItemSaleConditions.FixPrice)
	}
}

func main() {
	var wg sync.WaitGroup
	config = Configuration{}

	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		panic(err)
	}

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logging.WARNING, "")
	logging.SetBackend(backend1Leveled)
	log.Info("Started...")

	if config.DEBUG == true {
		file_log, err := os.Create("wallago.log")
		check(err)
		backend2 := logging.NewLogBackend(file_log, "", 0)
		backend2Formatter := logging.NewBackendFormatter(backend2, format)
		backend2Leveled := logging.AddModuleLevel(backend2Formatter)
		backend2Leveled.SetLevel(logging.DEBUG, "")
		logging.SetBackend(backend2Leveled)
	}

	file, err := os.Open("items.list")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go check_item(&wg, line)
	}

	err = scanner.Err()
	check(err)
	wg.Wait()

}
