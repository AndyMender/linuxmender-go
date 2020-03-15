package paths

import "path/filepath"

// ProjectPath is the main project's root path
var ProjectPath, _ = filepath.Abs(".")

// StaticPath points to the "static" sub-dir
var StaticPath = filepath.Join(ProjectPath, "static")

// ViewsPath points to the "views" sub-dir
var ViewsPath = filepath.Join(ProjectPath, "views")

// EntriesPath points to the JSON containing blog entry definitions
var EntriesPath = filepath.Join(StaticPath, "entries.json")

// DBPath points to the SQLite3 database containing static Web page definitions
var DBPath = filepath.Join(StaticPath, "db.sqlite3")
