package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gizarash/task-manager/internal/store"
)

const storeFile = "store.json"

func main() {

	store, err := store.New()
	if err != nil {
		fmt.Println(err)
		return
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
				store.Add(value)
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
				if ok := store.MarkDone(id); ok {
					changed = true
					fmt.Printf("Пункт с id = \"%d\" помечен выполненным\n", id)
				} else {
					fmt.Printf("Пункт с id = \"%d\" не найден или уже был отмечен выполненным ранее\n", id)
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
				if ok := store.Delete(id); ok {
					changed = true
					fmt.Printf("Пункт с id = \"%d\" удален\n", id)
				} else {
					fmt.Printf("id \"%d\" не найден\n", id)
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
		err := store.Save()
		if err != nil {
			fmt.Printf("Ошибка при кодировании json в файл %s: %s\n", storeFile, err)
		}
	}

}
