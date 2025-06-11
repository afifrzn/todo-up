package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Task struct {
    Description string
    Done        bool
}

var tasks []Task
var username string
var fileName string

func main() {
    // Input username di awal
    fmt.Print("Masukkan username: ")
    fmt.Scanln(&username)
    fileName = fmt.Sprintf("tasks_%s.json", username)

    loadTasks()

    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println("\n===== MENU TO-DO =====")
        fmt.Println("1. Tambah tugas")
        fmt.Println("2. Lihat semua tugas")
        fmt.Println("3. Tandai tugas selesai")
        fmt.Println("4. Keluar")
        fmt.Print("Pilih menu: ")

        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        switch input {
        case "1":
            fmt.Print("Masukkan deskripsi tugas: ")
            desc, _ := reader.ReadString('\n')
            addTask(strings.TrimSpace(desc))
        case "2":
            listTasks()
        case "3":
            listTasks()
            fmt.Print("Masukkan nomor tugas yang selesai: ")
            numStr, _ := reader.ReadString('\n')
            index, err := strconv.Atoi(strings.TrimSpace(numStr))
            if err != nil || index < 1 || index > len(tasks) {
                fmt.Println("Nomor tidak valid.")
                continue
            }
            markDone(index - 1)
        case "4":
            fmt.Println("Keluar. Sampai jumpa!")
            saveTasks()
            return
        default:
            fmt.Println("Pilihan tidak dikenal.")
        }

        saveTasks()
    }
}

func addTask(description string) {
    tasks = append(tasks, Task{Description: description})
    fmt.Println("Tugas ditambahkan:", description)
}

func listTasks() {
    if len(tasks) == 0 {
        fmt.Println("Belum ada tugas.")
        return
    }
    fmt.Println("\nDaftar Tugas:")
    for i, task := range tasks {
        status := " "
        if task.Done {
            status = "x"
        }
        fmt.Printf("[%s] %d: %s\n", status, i+1, task.Description)
    }
}

func markDone(index int) {
    tasks[index].Done = true
    fmt.Println("Tugas ditandai selesai:", tasks[index].Description)
}

func saveTasks() {
    file, err := os.Create(fileName)
    if err != nil {
        fmt.Println("Gagal menyimpan tugas:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(tasks)
    if err != nil {
        fmt.Println("Gagal encode data:", err)
    }
}

func loadTasks() {
    file, err := os.Open(fileName)
    if err != nil {
        if os.IsNotExist(err) {
            return
        }
        fmt.Println("Gagal membuka file:", err)
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&tasks)
    if err != nil {
        fmt.Println("Gagal membaca data tugas:", err)
    }
}
