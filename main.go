package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)


func main(){
	if len(os.Args) == 1 {
		fmt.Println("Введите ссылку для получения данных в строке вызова.")
		return
	}
	incomingArgs := os.Args[1]
	if incomingArgs == "" {
		log.Fatal("Введите адресную строку для запроса")
	}
	if !strings.Contains(incomingArgs, "http"){
		incomingArgs = "https://" + incomingArgs
	}
	text, err := Response(incomingArgs)
	if err != nil {
		log.Printf("Ошибка запроса %s", err)
	}
	_ = GetAllUrls(text)
}

func Response(url string) (string, error){
	var text []byte
	resp, err := http.Get(url)

	if err != nil {
		log.Println("Ошибка получения ответа")
		return "", err
	}
	if resp.StatusCode != http.StatusOK{
		return "", fmt.Errorf("Ошибка %s запроса к %s", resp.Status, url)
	}

	text, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error by %s reading %v", url, resp.Body)
		return "", err
	}
	resp.Body.Close()

	return string(text), err
}

func GetAllUrls(s string) []string {
	var urls, set []string
	var setMap = make(map[string]int)


	//Выбор всех совпадений по шаблону в тексте -> []string
	pattern := `("http.+?")`
	if ok, _ := regexp.Match(pattern, []byte(s)); ok {
		re := regexp.MustCompile(pattern)
		urls = re.FindAllString(s, -1)
	}

	//Выбор уникальных значений -> map[string]int
	for _, str := range urls {
		setMap[str]++
	}
	for key := range setMap{
		set = append(set, key)
	}
	//Сортировка
	sort.Strings(set)
	for _, url := range set {
		fmt.Printf("%d совпадений: %s\n", setMap[url], url)
	}
	return set
}