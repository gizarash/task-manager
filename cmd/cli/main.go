package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gizarash/task-manager/internal/model"
)

const storeFile = "store.json"

func main() {
	
	file, err := os.OpenFile(storeFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Ошибка при открытии/создании файла %s: %s\n", storeFile, err)
		return
	}
	defer file.Close()
	
	stat, err := file.Stat()
	if err != nil {
		fmt.Printf("Ошибка при получении информации о файле %s: %s\n", storeFile, err)
		return
	}

	var store model.Store
	if stat.Size() == 0 {
		store.CurrentId = 1
		store.Todos = []model.Todo{}
	} else {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&store)
		if err != nil {
			fmt.Printf("Ошибка при декодировании файла %s: %s\n", storeFile, err)
			return
		}
	}
	
	changed := false
	if len(os.Args) > 1 {
		command := os.Args[1]
		value := ""
		if len(os.Args) > 2 {
			value = strings.Join(os.Args[2:], " ")
		}
		switch command {
		case "add":
			if value == "" {
				fmt.Println("Добавляемое значение не может быть пустым")
			} else {
				store.Todos = append(store.Todos, model.Todo{Id: store.CurrentId, Title: value, Done: false})
				store.CurrentId++
				changed = true
				fmt.Printf("Значение \"%s\" добавлено в список\n", value)
			}
		case "list":
			if len(store.Todos) > 0 {
				for _, t := range store.Todos {
					doneFlag := " "
					if t.Done {
						doneFlag = "x"
					}
					fmt.Printf("[%d] [%s] %s\n", t.Id, doneFlag, t.Title)
				}
			} else {
				fmt.Println("Список дел пуст")
			}
		case "done":
			if value == "" {
				fmt.Println("id не может быть пустым")
			} else {
				id, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("Необходимо передать id, вы передали %s\n", value)
					return
				}
				notFound := true
				for i, t := range store.Todos {
					if t.Id == id {
						store.Todos[i].Done = true
						fmt.Printf("Пункт \"%s\" помечен выполненным\n", t.Title)
						notFound = false
						changed = true
						break
					}
				}
				if notFound {
					fmt.Printf("id \"%s\" не найден\n", value)
				}
			}
		case "delete":
			if value == "" {
				fmt.Println("id не может быть пустым")
			} else {
				id, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("Необходимо передать id, вы передали %s\n", value)
					return
				}
				notFound := true
				for i, t := range store.Todos {
					if t.Id == id {
						store.Todos = slices.Delete(store.Todos, i, i + 1)
						notFound = false
						changed = true
						fmt.Printf("Пункт \"%s\" удален\n", t.Title)
						break
					}
				}
				if notFound {
					fmt.Printf("id \"%s\" не найден\n", value)
				}
			}
		default:
			fmt.Println("Неизвестная команда")
		}
	} else {
		fmt.Println("Примеры аргументов программы:")
		fmt.Println("Добавить в список новый пункт - add \"Купить молоко\"")
		fmt.Println("Просмотр списка - list")
		fmt.Println("Пометить пункт выполненным - done 1")
		fmt.Println("Удалить пункт - delete 1")
	}

	if changed {
		encoder := json.NewEncoder(file)
		file.Truncate(0)
		file.Seek(0, 0)
		err = encoder.Encode(store)
		if err != nil {
			fmt.Printf("Ошибка при кодировании json в файл %s: %s\n", storeFile, err)
		}
	}

}
