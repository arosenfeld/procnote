package main

import (
    "fmt"
    "os"
    "os/user"
    "path/filepath"

    "gopkg.in/alecthomas/kingpin.v1"
)

func add(pid int, note string) {
}

func save(notes map[int]string, file *os.File) {
}

func makeNoteFile(path string) (*os.File, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        _, err := os.Create(path)
        if err != nil {
            return nil, err
        }
    }

    return os.Open(path)
}

func readNotes(file *os.File) (map[int]string, error) {
    notes := make(map[int]string)
    return notes, nil
}

func main() {
    var (
        action = kingpin.Command("add", "Adds a note to a process.")
        pid = action.Arg("pid", "Process ID.").Required().Int()
        note = action.Arg("note", "Note for process.").Required().String()
    )

    usr, err := user.Current()
    if err != nil {
        fmt.Println("Unable to get current user.", err)
        return
    }

    noteFile, err := makeNoteFile(filepath.Join(usr.HomeDir, ".procnote"))
    if err != nil {
        fmt.Println("Unable to create or open note file:", err)
        return
    }

    notes, err := readNotes(noteFile)
    if err != nil {
        fmt.Println("Unable to read note file:", err)
        return
    }

    switch kingpin.Parse() {
        case "add":
            add(*pid, *note)
            break
    }

    save(notes, noteFile)
}
