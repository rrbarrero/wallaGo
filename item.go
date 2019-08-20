package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type Item struct {
	Name         string
	ValidWords   []string
	InvalidWords []string
	Max          float64
	Min          float64
	ExceptionIDs []int64
}

func NewItem(itemLine string) *Item {
	fields := strings.Split(itemLine, "|")
	item := new(Item)
	item.Min = 0
	item.Max = 0
	item.Name = fields[0]
	item.ValidWords = strings.Split(fields[1], " ")
	item.InvalidWords = strings.Split(fields[2], " ")
	if max, err := strconv.ParseFloat(fields[3], 64); err == nil {
		item.Max = max
	}
	if min, err := strconv.ParseFloat(fields[4], 64); err == nil {
		item.Min = min
	}
	item.ExceptionIDs = get_exception_ids(fields[5])
	return item
}

func get_exception_ids(field string) []int64 {
	if field == "" {
		return []int64{0}
	}
	var ids []int64
	for _, id := range strings.Split(field, " ") {
		wId, err := strconv.ParseInt(id, 10, 64)
		check(err)
		ids = append(ids, wId)
	}
	return ids
}

// Hace la búsqueda en wallapop y nos devuelve
// un objeto tipo WallaItems que a su vez contiene
// Items que son los resultados de dicha búsqueda.
func search(url string) ([]WallaItem, error) {
	res, err := http.Get(url)
	check(err)
	defer res.Body.Close()
	var wallaItems []WallaItem
	body, err := ioutil.ReadAll(res.Body)
	result := gjson.Get(string(body), "items")
	result.ForEach(func(key, value gjson.Result) bool {
		var item WallaItem
		_item := gjson.Get(value.String(), "item")
		err = json.Unmarshal([]byte(_item.String()), &item)
		check(err)
		if config.DEBUG == true {
			log.Debug(item.URL)
		}
		wallaItems = append(wallaItems, item)
		return true // keep iterating
	})
	return wallaItems, err
}

// Retorna True si se encuentra en las excepciones
func check_if_exception_id(item Item, wallaItem WallaItem) bool {
	for _, eid := range item.ExceptionIDs {
		if eid == wallaItem.ItemId {
			return true
		}
	}
	return false
}

// Retorna true si se encuentra la palabra en el campo
func check_if_word(word string, field_text string) bool {
	formated_word := strings.ToLower(word)
	formated_field := strings.ToLower(field_text)
	if strings.Contains(formated_field, formated_word) {
		return true
	}
	return false
}

// Pasada una lista de palabras invalidas comprueba una por una si se
// encuentra en el campo. Si una palabra es encontrada sale de la función
// devolviendo true
func check_invalid_words(invalid_words []string, field_text string) bool {
	if len(invalid_words) > 0 && invalid_words[0] != "" {
		for _, word := range invalid_words {
			if check_if_word(word, field_text) == true {
				return true
			}
		}
	}
	return false
}

// Complementaria a la anterior
func check_valid_words(valid_words []string, field_text string) bool {
	any := true
	if len(valid_words) > 0 && valid_words[0] != "" {
		any = false
		for _, word := range valid_words {
			if check_if_word(word, field_text) == true {
				any = true
			}
		}
	}
	return any
}

// Compara si un item de wallapop cumple las condiciones
// propuestas en nuestra lista de items
func compare(item Item, wallaItem WallaItem) bool {
	if wallaItem.Reserved == true {
		return false
	}
	if wallaItem.SalePrice > item.Max {
		return false
	}
	if wallaItem.SalePrice < item.Min {
		return false
	}
	if wallaItem.Sold == true {
		return false
	}
	if wallaItem.SellerUser.Banned == true {
		return false
	}
	if config.SCORING_VALIDATION == true {
		if wallaItem.SellerUser.Validation.ScoringStarts < config.SCORING_VALIDATION_MIN_STARTS {
			return false
		}
	}
	if check_invalid_words(item.InvalidWords, wallaItem.Title) == true {
		if config.DEBUG == true {
			log.Debug("Compare false by invalid words, title")
			log.Debug(WALLA_ITEM_URL + wallaItem.URL)
		}
		return false
	}
	if check_invalid_words(item.InvalidWords, wallaItem.Description) == true {
		if config.DEBUG == true {
			log.Debug("Compare false by invalid words, description")
			log.Debug(WALLA_ITEM_URL + wallaItem.URL)
		}
		return false
	}
	valid_word_desc := check_valid_words(item.ValidWords, wallaItem.Description)
	valid_word_title := check_valid_words(item.ValidWords, wallaItem.Title)
	if valid_word_desc == false && valid_word_title == false {
		if config.DEBUG == true {
			log.Debug("Compare false by valid words")
			log.Debug(WALLA_ITEM_URL + wallaItem.URL)
		}
		return false
	}
	if check_if_exception_id(item, wallaItem) == true {
		if config.DEBUG == true {
			log.Debug("Compare false by check_if_exception_id")
		}
		return false
	}
	return true
}

func (item Item) CheckItem() []WallaItem {
	var result []WallaItem
	query := fmt.Sprintf(
		config.URL_TPLE,
		fmt.Sprintf("%f", item.Min),
		fmt.Sprintf("%f", item.Max),
		url.QueryEscape(item.Name),
	)
	if config.DEBUG == true {
		log.Debug(query)
	}
	wallaItems, err := search(query)
	check(err)
	for _, wallaItem := range wallaItems {
		if compare(item, wallaItem) == true {
			result = append(result, wallaItem)
		}
	}
	return result
}
