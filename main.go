package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	ADD_TASK               = "Add Task"
	VIEW_TASKS             = "View Tasks"
	VIEW_TASK_BY_ID        = "View Task by ID"
	VIEW_BY_STATUS         = "View Tasks by Status"
	COMPLETE_TASK          = "Complete Task"
	UPDATE_TASK            = "Update Task"
	DELETE_TASK            = "Delete Task"
	EXIT                   = "Exit"
	NOT_COMPLETED   Status = "Not Completed"
	PENDING         Status = "Pending"
	COMPLETED       Status = "Completed"
)

type Command struct {
	Label string
	Run   func(*TaskList)
}

var commands = []Command{
	{Label: ADD_TASK, Run: addTask},
	{Label: VIEW_TASKS, Run: viewTasks},
	{Label: VIEW_TASK_BY_ID, Run: viewTaskByID},
	{Label: VIEW_BY_STATUS, Run: viewTasksByStatus},
	{Label: COMPLETE_TASK, Run: markAsComplete},
	{Label: UPDATE_TASK, Run: updateTaskDescription},
	{Label: DELETE_TASK, Run: deleteTask},
	{Label: EXIT, Run: nil},
}

func printCommands() {
	fmt.Println("Enter your input from these options:")
	for i, cmd := range commands {
		fmt.Printf("%d. %s\n", i+1, cmd.Label)
	}
}
func terminalInput() int {
	printCommands()
	var input int
	fmt.Scanln(&input)
	return input
}

func execute(option int, tasks *TaskList) bool {
	if option < 1 || option > len(commands) {
		fmt.Println("Invalid option. Please try again.")
		return true
	}
	cmd := commands[option-1]
	if cmd.Run == nil {
		fmt.Println("Exiting the application. Goodbye!")
		return false
	}
	cmd.Run(tasks)
	return true
}

// task.go
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Completed   bool      `json:"completed"`
}

var fileName = "taskList.json"

type TaskList struct {
	Tasks []Task `json:"tasks"`
}
type Status string

var reader = bufio.NewReader(os.Stdin)

func nextID(tasks *TaskList) int {
	if len(tasks.Tasks) == 0 {
		return 1
	}
	return tasks.Tasks[len(tasks.Tasks)-1].ID + 1
}

func idExists(tasks *TaskList, id int) bool {
	for _, task := range tasks.Tasks {
		if task.ID == id {
			return true
		}
	}
	fmt.Println("Task ID not found.")
	return false
}

func addTask(tasks *TaskList) {
	fmt.Print("Enter task description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	if description == "" {
		fmt.Println("Task description cannot be empty.")
		return
	}
	id := nextID(tasks)

	task := Task{
		ID:          id,
		Description: description,
		Status:      PENDING,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Completed:   false,
	}

	tasks.Tasks = append(tasks.Tasks, task)
	writeToFile(tasks)
	fmt.Println("New task added successfully!")
}

func viewTasks(tasks *TaskList) {
	if len(tasks.Tasks) == 0 {
		fmt.Println("No tasks available.")
		return
	}

	for _, task := range tasks.Tasks {
		fmt.Printf("[#%d] %s | %s | Completed: %v | Created: %s\n", task.ID, task.Description, task.Status, task.Completed, task.CreatedAt.Format(time.RFC3339))
	}
}

func viewTaskByID(tasks *TaskList) {
	var id int
	fmt.Print("Enter task ID to view: ")
	fmt.Scanln(&id)

	if !idExists(tasks, id) {
		return
	}
	for _, task := range tasks.Tasks {
		if task.ID == id {
			fmt.Printf("[#%d] %s | %s | Completed: %v | Created: %s\n", task.ID, task.Description, task.Status, task.Completed, task.CreatedAt.Format(time.RFC3339))
			return
		}
	}
}

func viewTasksByStatus(tasks *TaskList) {
	var statusInput string
	fmt.Print("Enter task status to view (N (Not Completed), P (Pending), C (Completed)): ")
	statusInput, _ = reader.ReadString('\n')
	statusInput = strings.TrimSpace(statusInput)
	statusInput = strings.ToUpper(statusInput)

	if statusInput != "N" && statusInput != "P" && statusInput != "C" {
		fmt.Println("Invalid status input. Please enter N, P, or C.")
		return
	}

	var found bool
	for _, task := range tasks.Tasks {
		if string(task.Status[0]) == statusInput {
			fmt.Printf("[#%d] %s | %s | Completed: %v | Created: %s\n", task.ID, task.Description, task.Status, task.Completed, task.CreatedAt.Format(time.RFC3339))
			found = true
		}
	}
	if !found {
		fmt.Printf("No tasks found with status: %s\n", statusInput)
	}
}

func deleteTask(tasks *TaskList) {
	var id int
	fmt.Print("Enter task ID to delete: ")
	fmt.Scanln(&id)

	if !idExists(tasks, id) {
		return
	}

	for i, task := range tasks.Tasks {
		if task.ID == id {
			tasks.Tasks = append(tasks.Tasks[:i], tasks.Tasks[i+1:]...)
			writeToFile(tasks)
			fmt.Printf("Task with ID %d deleted.\n", id)
			return
		}
	}
}

func updateTaskDescription(tasks *TaskList) {
	var id int
	fmt.Print("Enter task ID to update: ")
	fmt.Scanln(&id)

	if !idExists(tasks, id) {
		return
	}

	newDescription, _ := reader.ReadString('\n')
	newDescription = strings.TrimSpace(newDescription)

	for i, task := range tasks.Tasks {
		if task.ID == id {
			tasks.Tasks[i].Description = newDescription
			tasks.Tasks[i].UpdatedAt = time.Now()
			break
		}
	}
	fmt.Printf("Task with ID %d updated.\n", id)
	viewTasks(tasks)
}

func markAsComplete(tasks *TaskList) {
	var id int
	fmt.Print("Enter task ID to complete: ")
	fmt.Scanln(&id)

	if !idExists(tasks, id) {
		return
	}

	for i, task := range tasks.Tasks {
		if task.ID == id {
			if task.Completed {
				fmt.Printf("Task with ID %d is already completed.\n", id)
				return
			}
			tasks.Tasks[i].Completed = true
			tasks.Tasks[i].Status = COMPLETED
			break
		}
	}
	fmt.Printf("\nTask with ID %d marked as complete.\n", id)
	viewTasks(tasks)
}

func loadTasks() (*TaskList, error) {

	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return &TaskList{Tasks: []Task{}}, nil
		}
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}

	var taskList TaskList
	err = json.Unmarshal(file, &taskList)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %v", err)
	}

	return &taskList, nil
}

func writeToFile(taskList *TaskList) {
	data, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling tasks: %v\n", err)
		return
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func main() {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	for {
		userOption := terminalInput()
		if !execute(userOption, tasks) {
			break
		}
	}

}
