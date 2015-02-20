package main

import (
    "bufio"
    "fmt"
    "os"
    "os/user"
    "strconv"
    "strings"
    "path/filepath"

    "gopkg.in/alecthomas/kingpin.v1"
)

type ProcNotes map[int]string

func dirExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func printNotes(notes ProcNotes) {
    if len(notes) == 0 {
        fmt.Println("No notes.")
        return
    }

    fmt.Printf("PID\tNOTE\n")
    for pid, note := range notes {
        exists, err := dirExists(filepath.Join("/proc", strconv.Itoa(pid)))
        status := "Unknown"
        if err == nil {
            if exists {
                status = "Running"
            } else {
                status = "Stopped"
            }
        }

        fmt.Printf("%d\t%s\t%s\n", pid, status, note)
    }
}

func searchNotes(notes ProcNotes, query string) ProcNotes {
    query = strings.ToLower(query)
    found := make(ProcNotes)

    for pid, note := range notes {
        if strings.Index(strings.ToLower(note), query) >= 0 {
            found[pid] = note
        }
    }

    return found
}

func saveNoteFile(notes ProcNotes, file *os.File) {
    file.Truncate(0)
    file.Seek(0, 0)
    for pid, note := range notes {
        _, err := file.WriteString(fmt.Sprintf("%d %s\n", pid, note))
        if err != nil {
            fmt.Println(err)
        }
    }
    file.Sync()
}

func makeNoteFile(path string) (*os.File, error) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        fh, err := os.Create(path)
        if err != nil {
            return nil, err
        }
        return fh, nil
    }

    return os.OpenFile(path, os.O_RDWR, os.FileMode(0666))
}

func readNoteFile(file *os.File) (ProcNotes, error) {
    notes := make(ProcNotes)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        entry := strings.SplitN(strings.TrimSpace(scanner.Text()), " ", 2)
        pid, err := strconv.Atoi(entry[0])
        if err != nil {
            fmt.Println("Unable to read note entry: ", entry, entry[0])
            continue
        }
        notes[pid] = entry[1]
    }
    return notes, nil
}

func main() {
    var (
        addAction = kingpin.Command("add", "Adds a note to a process.")
        addPid = addAction.Arg("pid", "Process ID.").Required().Int()
        note = addAction.Arg("note", "Note for process.").Required().String()
        delAction = kingpin.Command("del", "Deletes a note to a process.")
        delPid = delAction.Arg("pid", "Process ID.").Required().Int()

        _ = kingpin.Command("list", "Lists all process notes.")
        _ = kingpin.Command("clear", "Clears all process notes.")
        searchAction = kingpin.Command("search", "Searches for a process " +
                                       "matching the query.")
        query = searchAction.Arg("query",
                                 "The query to run.").Required().String()
    )

    usr, err := user.Current()
    if err != nil {
        fmt.Println("Unable to get current user.", err)
        return
    }

    noteFile, err := makeNoteFile(filepath.Join(usr.HomeDir, ".procnote"))
    defer noteFile.Close()

    if err != nil {
        fmt.Println("Unable to create or open note file:", err)
        return
    }

    notes, err := readNoteFile(noteFile)
    if err != nil {
        fmt.Println("Unable to read note file:", err)
        return
    }

    switch kingpin.Parse() {
        case "add":
            notes[*addPid] = *note
            saveNoteFile(notes, noteFile)
            break
        case "del":
            delete(notes, *delPid)
            saveNoteFile(notes, noteFile)
            break
        case "list":
            printNotes(notes)
            break
        case "search":
            printNotes(searchNotes(notes, *query))
            break
        case "clear":
            notes := make(ProcNotes)
            saveNoteFile(notes, noteFile)
            break
    }
}
