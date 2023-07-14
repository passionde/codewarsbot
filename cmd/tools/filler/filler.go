package filler

import (
	"bufio"
	"fmt"
	"github.com/Yarik-xxx/CodeWarsRestApi/cmd/tools/helper"
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/app/store"
	"github.com/Yarik-xxx/CodeWarsRestApi/pkg/codewars"
	"log"
	"os"
	"time"
)

func FillingDatabase() {
	file, err := os.Open("cmd/tools/filler/names.txt")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer file.Close()

	// Получение имен
	fileScanner := bufio.NewScanner(file)
	names := make(map[string]interface{})

	for fileScanner.Scan() {
		userName := fileScanner.Text()
		names[userName] = nil
	}

	// Получение ID кат
	ids := make(map[string]interface{})
	for name, _ := range names {
		totalPage := 1
		for page := 0; page <= totalPage; page++ {
			completed, err := codewars.GetCompletedChallenges(name, page)
			for err != nil && err.Error() == "429" {
				time.Sleep(5 * time.Second)
				completed, err = codewars.GetCompletedChallenges(name, page)
			}
			if err != nil {
				log.Println(err)
				continue
			}

			totalPage = completed.TotalPages

			for _, info := range completed.Data {
				ids[info.Id] = nil
			}
		}
		log.Println("Итого:", len(ids))
	}

	// Запись ID в файл (на всякий случай)
	file, err = os.Create("cmd/tools/filler/ids.txt")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()

	for id, _ := range ids {
		if _, err := file.WriteString(id + "\n"); err != nil {
			log.Println(err)
		}
	}

	// Заполнение БД информацией
	LazyFilling()

}

func LazyFilling() {
	file, err := os.Open("cmd/tools/filler/ids.txt")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	defer file.Close()

	// Получение имен
	fileScanner := bufio.NewScanner(file)

	st := store.New(helper.InitConfig())
	err = st.Open()
	if err != nil {
		log.Fatal(err)
	}

	type TmpStruct struct {
		inCache bool
		err     error
	}

	go func() {
		for {
			<-time.After(60 * time.Second)
			fmt.Println("Всего собрано информации о", st.Challenge().Count(), "катах")
		}

	}()

	for fileScanner.Scan() {
		id := fileScanner.Text()

		_, _, err := st.Challenge().Get(id)

		for err != nil && err.Error() == "429" {
			time.Sleep(5 * time.Second)
			_, _, err = st.Challenge().Get(id)
		}

		if err != nil {
			log.Println(err)
			continue
		}
	}
}
