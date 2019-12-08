package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var MediaFiles = []string{".jpg", ".png", ".ico", ".js", ".gif", ".webp", ".css",}

type UrlsList struct {
	List    []Url `json:"list"`
	Channel chan Url
}

type Url struct {
	Link string `json:"link"`
}

func (u *Url) Response() (string, error) {
	var text []byte
	resp, err := http.Get(u.Link)

	if err != nil {
		log.Println("Ошибка получения ответа")
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ошибка %s запроса к %s\n", resp.Status, u.Link)
	}

	text, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Error by %s reading %v\n", u.Link, resp.Body)
		return "", err
	}
	if err := resp.Body.Close(); err != nil {
		log.Println(err)
	}

	return string(text), err
}

func (l *UrlsList) GetUrls(text string, wg *sync.WaitGroup) {
	var setMap = make(map[Url]int)
	var urls []Url
	var key Url

	//Выбор всех совпадений по шаблону в тексте -> []string
	pattern := `"(http.+?)"`
	if ok, _ := regexp.Match(pattern, []byte(text)); ok {
		re := regexp.MustCompile(pattern)
		for _, url := range re.FindAllString(text, -1) {
			url = strings.ReplaceAll(url,`"`, ``)
			urls = append(urls, Url{url})
		}
	}
	//Функция проверки наличия медиа-файлов в суфиксе ссылки
	checkFunc := func(u Url) bool {
		for _, ext := range MediaFiles {
			if strings.Contains(filepath.Ext(u.Link), ext) {
				return true
			}
		}
		return false
	}

	//Выбор уникальных значений -> map[Url]int
	for _, url := range urls {
		if !checkFunc(url) {
			setMap[url]++
		}
	}

	for key = range setMap {
		l.List = append(l.List, key)
	}

	go SaveStr(l, wg)
}

func main() {
	var wg sync.WaitGroup
	var urlsList UrlsList
	var url Url

	if len(os.Args) == 1 {
		fmt.Println("Введите ссылку для получения данных в строке вызова.")
		return
	}
	url.Link = os.Args[1]
	if url.Link == "" {
		log.Fatal("Введите адресную строку для запроса")
	}
	if !strings.Contains(url.Link, "http") {
		url.Link = "https://" + url.Link
	}
	text, err := url.Response()
	if err != nil {
		log.Printf("Ошибка запроса %s\n", err)
	}

	go urlsList.GetUrls(text, &wg)
	wg.Add(1)
	wg.Wait()

}

func SaveStr(l *UrlsList, wg *sync.WaitGroup) {
	var data []Url
	fileName := fmt.Sprintf("%v.json", time.Now())
	data = append(data, l.List...)

	jsonList, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		log.Fatal("Ошибка сериализации:", err)
	}
	err = ioutil.WriteFile(fileName, jsonList, 777)
	if err != nil {
		log.Fatal("Ошибка записи файла:", err)
	}

	fmt.Printf("Файл %q сохранен.\nНайдено ссылок: %d\n", fileName, len(data))
	wg.Done()
}
