# Code Conventions

These arent strict rules, just a description of how the code is generally written so you know what to expect when reading through it or contributing.

## Receiver name

Always `this`:
```go
func (this *SClient) Active() bool { return this.Connection != nil }
```

## Type prefixes

Theres a loose system of prefixes that tell you what kind of type youre looking at:

| Prefix | Kind | Examples |
|:-------|:-----|:---------|
| `S` | Struct | `SClient`, `SSocket`, `SRequest` |
| `M` | Map | `MFileState`, `MLogLevel` |
| `T` | Defined type / alias | `TLogLevel`, `TMethod` |
| `E` | Enum / bitmask | `EFileEvent`, `ELogLevel` |
| `Fn` | Function type | `FnValidator`, `FnExecution` |

That said, not every type follows this. Some packages just dont use prefixes at all.

## Pointer vs value receivers

Pointer receiver if the method mutates anything or if the type holds a mutex or a channel. Value receiver otherwise:
```go
func (this *SSocket) send(method TMethod, params any) *SMessage  // pointer: mutates
func (this SClog) Clone(override SClog) SClog                    // value: returns a copy
func (this argDef) String() string                               // value: read-only
```

## Variable declarations

`var` is the go-to:
```go
var str strings.Builder
var info, err = os.Stat(path)
```
`:=` only really shows up in `for range` loops and `if`-scoped declarations:
```go
for i, arg := range args { }
if idx := strings.LastIndex(partial, "/"); idx >= 0 { }
```

## Error handling

Check the error right after the call, no `fmt.Errorf` wrapping, no `defer`-error patterns. Just the straightforward `if err != nil` and log it.

## Logging

Every package that logs makes its own logger:

```go
var clog = clogger.Default.Clone(clogger.SClog{Name: "packagename"})
```

Each level has a regular and a formatted variant, so you get both `Error(msg)` and `Errorf(format, args...)` and the same goes for the rest.

## String building

`strings.Builder` for anything more than a one-liner:
```go
var builder strings.Builder
builder.WriteString(alias)
builder.WriteRune('\n')
return builder.String()
```
`str +=` should be avoided.

## Imports

Stdlib first, blank line, then everything else. Each group sorted:
```go
import (
    "os"
    "time"

    "github.com/CTNOriginals/BitburnerGoFilesync/clogger"
    "github.com/gorilla/websocket"
)
```

## Package structure

The entry file of a package is named after it - `config/config.go`, `watcher/watcher.go`, that kind of thing.

Sub-packages that implement CLI commands register themselves through `init()` by calling `commands.List.Push(&def)`. The `init.go` file at the project root imports them with a blank identifier to trigger all of that.

## Annotations

- `// NOTE: ` - Workarounds and non-obvious choices.
- `// TODO: ` - Incomplete features.
- `// BUG: ` - Known issues.

## Other bits

- **Mutex**: `sync.Mutex` with `Lock()`/`defer Unlock()` on types that get accessed from different goroutines.
- **Generics**: A `call[TResult any]()` helper in `websocket/socket.go` cuts down on boilerplate for the typed socket methods:
```go
func call[TResult any](this *SSocket, method TMethod, params any, callback func(*TResult)) {
    var message = this.send(method, params)
    var response = AwaitResponse[TResult](message)

    if callback != nil {
        callback(response)
    }
}
```
