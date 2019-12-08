package save_package

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func SaveStr(c <-chan string, wg *sync.WaitGroup) {
	fileName := fmt.Sprintf("%v.txt", time.Now())
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Ошибка создания файла (%s): %s.\n", fileName, err)
	}
	closeFile := func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}

	defer closeFile()
	for {
		if <-c == "stop" {
			fmt.Println("Channel was stop.")
			break
		} else {
			if _, err := file.WriteString(<-c + "\n"); err != nil {
				log.Printf("Ошибка записи в файл")
			}
		}
	}

	fmt.Printf("Файл %q сохранен.\n", fileName)
	wg.Done()
}
