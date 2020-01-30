package paths

import "path/filepath"

// ProjectPath is the main project's root path
var ProjectPath, _ = filepath.Abs(".")

// StaticPath points to the "static" sub-dir
var StaticPath = filepath.Join(ProjectPath, "static")

// EntriesPath points to the JSON containing blog entry definitions
var EntriesPath = filepath.Join(StaticPath, "entries.json")
