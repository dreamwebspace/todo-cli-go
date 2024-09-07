package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	Description string `json:"description"`
	IsCompleted bool   `json:"isCompleted"`
}

type TodoApp struct {
	tasks    []Task
	fileName string
}

func NewTodoApp() *TodoApp {
	app := &TodoApp{
		fileName: "tasks.json",
	}
	app.loadTasks()
	return app
}

func (app *TodoApp) loadTasks() {
	data, err := ioutil.ReadFile(app.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	err = json.Unmarshal(data, &app.tasks)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
	}
}

func (app *TodoApp) saveTasks() {
	data, err := json.Marshal(app.tasks)
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	err = ioutil.WriteFile(app.fileName, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
}

func (app *TodoApp) addTask(description string) {
	app.tasks = append(app.tasks, Task{Description: description, IsCompleted: false})
	app.saveTasks()
	app.listTasks()
}

func (app *TodoApp) listTasks() {
	if len(app.tasks) == 0 {
		fmt.Println("No tasks.")
	} else {
		fmt.Println()
		for i, task := range app.tasks {
			status := "[ ]"
			if task.IsCompleted {
				status = "[X]"
			}
			fmt.Printf("%d. %s %s\n", i+1, status, task.Description)
		}
		fmt.Println()
	}
}

func (app *TodoApp) toggleTaskCompletion(index int) {
	if index >= 0 && index < len(app.tasks) {
		app.tasks[index].IsCompleted = !app.tasks[index].IsCompleted
		app.saveTasks()
		app.listTasks()
	} else {
		fmt.Println("Invalid task number.")
	}
}

func (app *TodoApp) removeTask(index int) {
	if index >= 0 && index < len(app.tasks) {
		app.tasks = append(app.tasks[:index], app.tasks[index+1:]...)
		app.saveTasks()
		app.listTasks()
	} else {
		fmt.Println("Invalid task number.")
	}
}

func (app *TodoApp) moveTaskUp(index int) {
	if index > 0 && index < len(app.tasks) {
		app.tasks[index], app.tasks[index-1] = app.tasks[index-1], app.tasks[index]
		app.saveTasks()
		app.listTasks()
	} else {
		fmt.Println("Cannot move task up.")
	}
}

func (app *TodoApp) moveTaskDown(index int) {
	if index >= 0 && index < len(app.tasks)-1 {
		app.tasks[index], app.tasks[index+1] = app.tasks[index+1], app.tasks[index]
		app.saveTasks()
		app.listTasks()
	} else {
		fmt.Println("Cannot move task down.")
	}
}

func (app *TodoApp) renameTask(index int, newDescription string) {
	if index >= 0 && index < len(app.tasks) {
		oldDescription := app.tasks[index].Description
		app.tasks[index].Description = newDescription
		app.saveTasks()
		fmt.Printf("  From: %s\n", oldDescription)
		fmt.Printf("  To:   %s\n", newDescription)
		app.listTasks()
	} else {
		fmt.Println("Invalid task number.")
	}
}

func (app *TodoApp) processCommand(command string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	action := strings.ToLower(parts[0])
	switch action {
	case "a":
		if len(parts) > 1 {
			app.addTask(strings.Join(parts[1:], " "))
		} else {
			fmt.Println("Usage: a <task description>")
		}
	case "t":
		app.listTasks()
	case "x":
		if len(parts) > 1 {
			if taskNumber, err := strconv.Atoi(parts[1]); err == nil {
				app.toggleTaskCompletion(taskNumber - 1)
			} else {
				fmt.Println("Invalid task number.")
			}
		} else {
			fmt.Println("Usage: x <task number>")
		}
	case "d":
		if len(parts) > 1 {
			if taskNumber, err := strconv.Atoi(parts[1]); err == nil {
				app.removeTask(taskNumber - 1)
			} else {
				fmt.Println("Invalid task number.")
			}
		} else {
			fmt.Println("Usage: d <task number>")
		}
	case "h":
		if len(parts) > 1 {
			if taskNumber, err := strconv.Atoi(parts[1]); err == nil {
				app.moveTaskUp(taskNumber - 1)
			} else {
				fmt.Println("Invalid task number.")
			}
		} else {
			fmt.Println("Usage: h <task number>")
		}
	case "l":
		if len(parts) > 1 {
			if taskNumber, err := strconv.Atoi(parts[1]); err == nil {
				app.moveTaskDown(taskNumber - 1)
			} else {
				fmt.Println("Invalid task number.")
			}
		} else {
			fmt.Println("Usage: l <task number>")
		}
	case "r":
		if len(parts) > 2 {
			if taskNumber, err := strconv.Atoi(parts[1]); err == nil {
				app.renameTask(taskNumber-1, strings.Join(parts[2:], " "))
			} else {
				fmt.Println("Invalid task number.")
			}
		} else {
			fmt.Println("Usage: r <task number> <new task description>")
		}
	case "q":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "?":
		app.printHelp()
	default:
		fmt.Println("Unknown command. Type \"?\" for help.")
	}
}

func (app *TodoApp) printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  a <task description> - Add a new task")
	fmt.Println("  t - List all tasks")
	fmt.Println("  x <task number> - Mark task as complete/incomplete")
	fmt.Println("  d <task number> - Remove task")
	fmt.Println("  h <task number> - Move task higher")
	fmt.Println("  l <task number> - Move task lower")
	fmt.Println("  r <task number> <new description> - Rename task")
	fmt.Println("  ? - Show this help message")
	fmt.Println("  q - Quit the application")
}

func (app *TodoApp) run() {
	app.listTasks() // Display tasks at the start

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			app.processCommand(scanner.Text())
		} else {
			break
		}
	}
}

func main() {
	app := NewTodoApp()
	app.run()
}
