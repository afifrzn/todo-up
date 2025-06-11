package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net"
    "os"
    "strconv"
    "strings"
)

type Task struct {
    Description string
    Done        bool
}

func main() {
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    fmt.Println("Listening on port 8080...")

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
    conn.Write([]byte("Welcome to TODO Netcat!\nEnter your username: "))

    scanner := bufio.NewScanner(conn)
    if !scanner.Scan() {
        return
    }
    username := strings.TrimSpace(scanner.Text())
    filename := fmt.Sprintf("tasks_%s.json", username)
    tasks := loadTasks(filename)

    conn.Write([]byte(fmt.Sprintf("Hi %s! You can now use commands: add, list, done\n", username)))

    for scanner.Scan() {
        input := scanner.Text()
        args := strings.Fields(input)
        if len(args) == 0 {
            continue
        }

        cmd := args[0]
        switch cmd {
        case "add":
            desc := strings.Join(args[1:], " ")
            tasks = append(tasks, Task{Description: desc})
            conn.Write([]byte("Task added.\n"))
        case "list":
            if len(tasks) == 0 {
                conn.Write([]byte("No tasks found.\n"))
                continue
            }
            for i, task := range tasks {
                status := " "
                if task.Done {
                    status = "x"
                }
                conn.Write([]byte(fmt.Sprintf("[%s] %d: %s\n", status, i+1, task.Description)))
            }
        case "done":
            if len(args) < 2 {
                conn.Write([]byte("Usage: done <task number>\n"))
                continue
            }
            idx, err := strconv.Atoi(args[1])
            if err != nil || idx < 1 || idx > len(tasks) {
                conn.Write([]byte("Invalid task number.\n"))
                continue
            }
            tasks[idx-1].Done = true
            conn.Write([]byte("Marked as done.\n"))
        default:
            conn.Write([]byte("Unknown command.\n"))
        }
        saveTasks(filename, tasks)
    }
}

func loadTasks(filename string) []Task {
    file, err := os.Open(filename)
    if err != nil {
        return []Task{}
    }
    defer file.Close()

    var tasks []Task
    json.NewDecoder(file).Decode(&tasks)
    return tasks
}

func saveTasks(filename string, tasks []Task) {
    file, err := os.Create(filename)
    if err != nil {
        fmt.Println("Error saving tasks:", err)
        return
    }
    defer file.Close()
    json.NewEncoder(file).Encode(tasks)
}
