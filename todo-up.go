package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Task struct {
	Description string
	Done        bool
}

var (
	userTasks = make(map[string][]Task)
	mu        sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("TODO Netcat server running on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Simple TODO App Using GoLang\n"))
conn.Write([]byte("Masukkan Username: "))
reader := bufio.NewReader(conn)
username, err := reader.ReadString('\n')
if err != nil {
    return
}
username = strings.TrimSpace(username)
fmt.Printf("User '%s' connected from %s\n", username, conn.RemoteAddr())

conn.Write([]byte("Halo " + username + "!\n"))
conn.Write([]byte("Gunakan perintah: add <task>, list, done <no>, exit\n\n"))
	loadTasks(username)

	for {
		conn.Write([]byte("> ")) // prompt
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Connection closed for %s\n", username)
			return
		}
		text = strings.TrimSpace(text)

		if text == "exit" || text == "quit" || text == "keluar" {
			conn.Write([]byte("Babai\n"))
			fmt.Printf("User '%s' disconnected.\n", username)
			return
		}

		if text == "" {
			continue
		}

		parts := strings.SplitN(text, " ", 2)
		cmd := parts[0]

		mu.Lock()
		switch cmd {
		case "add":
			if len(parts) < 2 {
				conn.Write([]byte("Usage: add <task>\n"))
			} else {
				task := Task{Description: parts[1]}
				userTasks[username] = append(userTasks[username], task)
				conn.Write([]byte("Task added.\n"))
			}
		case "list":
			tasks := userTasks[username]
			if len(tasks) == 0 {
				conn.Write([]byte("No tasks found.\n"))
			} else {
				for i, task := range tasks {
					status := "[ ]"
					if task.Done {
						status = "[x]"
					}
					conn.Write([]byte(fmt.Sprintf("%d. %s %s\n", i+1, status, task.Description)))
				}
			}
		case "done":
			if len(parts) < 2 {
				conn.Write([]byte("Usage: done <task number>\n"))
			} else {
				index := -1
				fmt.Sscanf(parts[1], "%d", &index)
				if index >= 1 && index <= len(userTasks[username]) {
					userTasks[username][index-1].Done = true
					conn.Write([]byte("Task marked as done.\n"))
				} else {
					conn.Write([]byte("Invalid task number.\n"))
				}
			}
		default:
			conn.Write([]byte("Unknown command. Use: add, list, done, exit\n"))
		}
		mu.Unlock()

		saveTasks(username)
	}
}

func taskFileName(username string) string {
	return fmt.Sprintf("tasks_%s.json", username)
}

func loadTasks(username string) {
	file, err := os.Open(taskFileName(username))
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var tasks []Task
	if err := decoder.Decode(&tasks); err == nil {
		userTasks[username] = tasks
	}
}

func saveTasks(username string) {
	file, err := os.Create(taskFileName(username))
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	_ = encoder.Encode(userTasks[username])
}
