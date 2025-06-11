package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
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
		fmt.Println("Server error:", err)
		return
	}
	fmt.Println("💀 TODO Vuln Server listening on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	conn.Write([]byte("Welcome to VulnTODO!\n"))
	conn.Write([]byte("Enter your username: "))
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Printf("⚠️  User '%s' connected from %s\n", username, conn.RemoteAddr())

	loadTasks(username)
	conn.Write([]byte("Commands: add <task>, list, done <no>, exec <no>, exit\n\n"))

	for {
		conn.Write([]byte("> "))
		cmdline, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("User '%s' disconnected.\n", username)
			return
		}
		cmdline = strings.TrimSpace(cmdline)
		parts := strings.SplitN(cmdline, " ", 2)
		cmd := parts[0]

		mu.Lock()
		switch cmd {
		case "add":
			if len(parts) < 2 {
				conn.Write([]byte("Usage: add <task>\n"))
			} else {
				userTasks[username] = append(userTasks[username], Task{Description: parts[1]})
				conn.Write([]byte("Task added.\n"))
			}
		case "list":
			tasks := userTasks[username]
			for i, t := range tasks {
				status := "[ ]"
				if t.Done {
					status = "[x]"
				}
				conn.Write([]byte(fmt.Sprintf("%d. %s %s\n", i+1, status, t.Description)))
			}
		case "done":
			if len(parts) < 2 {
				conn.Write([]byte("Usage: done <no>\n"))
			} else {
				var idx int
				fmt.Sscanf(parts[1], "%d", &idx)
				if idx > 0 && idx <= len(userTasks[username]) {
					userTasks[username][idx-1].Done = true
					conn.Write([]byte("Marked done.\n"))
				} else {
					conn.Write([]byte("Invalid number.\n"))
				}
			}
		case "exec":
			if len(parts) < 2 {
				conn.Write([]byte("Usage: exec <no>\n"))
			} else {
				var idx int
				fmt.Sscanf(parts[1], "%d", &idx)
				if idx > 0 && idx <= len(userTasks[username]) {
					cmd := userTasks[username][idx-1].Description
					out, _ := exec.Command("sh", "-c", cmd).CombinedOutput()
					conn.Write([]byte(string(out)))
				} else {
					conn.Write([]byte("Invalid number.\n"))
				}
			}
		case "exit":
			conn.Write([]byte("Goodbye.\n"))
			mu.Unlock()
			saveTasks(username)
			return
		default:
			conn.Write([]byte("Unknown command.\n"))
		}
		mu.Unlock()

		saveTasks(username)
	}
}

func taskFileName(username string) string {
	return fmt.Sprintf("tasks_%s.json", username) // Vulnerable to traversal
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
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	_ = encoder.Encode(userTasks[username])
}
