package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tkanos/gonfig"
)

var config Configuration

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check_item(line string) {
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
	config = Configuration{}
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		panic(err)
	}
	file, err := os.Open("items.list")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		check_item(line)
	}

	err = scanner.Err()
	check(err)

}
