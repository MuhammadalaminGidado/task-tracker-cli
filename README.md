
# 📝 Go CLI Task Manager

A lightweight command-line task management application written in Go. Tasks are persisted to disk as JSON and can be managed via an interactive terminal menu.

---

## ✨ Features

- Add new tasks
- View all tasks
- View task by ID
- Filter tasks by status (Not Completed, Pending, Completed)
- Update task descriptions
- Mark tasks as complete
- Delete tasks
- Persistent storage using JSON file

---

## 🧱 Project Structure

├── main.go
├── task.go
└── taskList.json # auto-generated on first run

---

## ⚙️ Requirements

- Go 1.20+

---

## 🚀 Installation & Running

```bash
git clone <your-repo-url>
cd <your-repo>
go run .
```
---

📖 Usage

When you run the program:

Enter your input from these options:
1. Add Task
2. View Tasks
3. View Task by ID
4. View Tasks by Status
5. Complete Task
6. Update Task
7. Delete Task
8. Exit

Enter the corresponding number to execute an action.

___

🗂 Data Storage

Tasks are stored in:

taskList.json

Example:

{
  "tasks": [
    {
      "id": 1,
      "description": "Buy groceries",
      "status": "Pending",
      "createdAt": "2026-01-22T18:10:00Z",
      "updatedAt": "2026-01-22T18:10:00Z",
      "completed": false
    }
  ]
}

🧠 Status Values

| Status          | Meaning                      |
| --------------- | ---------------------------- |
| `Not Completed` | Created but not yet acted on |
| `Pending`       | Active task                  |
| `Completed`     | Finished                     |

Live URL: https://task-tracker-cli-brto.onrender.com

📄 License

MIT
