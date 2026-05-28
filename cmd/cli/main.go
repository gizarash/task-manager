package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gizarash/task-manager/internal/model"
)

func main() {
	
	currentId := 1
	var todos []model.Todo
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
				todos = append(todos, model.Todo{Id: currentId, Title: value, Done: false})
				currentId++
				fmt.Printf("Значение \"%s\" добавлено в список\n", value)
			}
		case "list":
			if len(todos) > 0 {
				for _, t := range todos {
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
				for i, t := range todos {
					if t.Id == id {
						todos[i].Done = true
						fmt.Printf("Пункт \"%s\" помечен выполненным\n", t.Title)
						notFound = false
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
				for i, t := range todos {
					if t.Id == id {
						todos = slices.Delete(todos, i, i + 1)
						notFound = false
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
}
