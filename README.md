# go-forth

# Card 1 — World and Room types

Define the foundational domain types and the World as the single source of truth
for all rooms and objects.

## Types to define

- `Direction` — string type with constants: `North`, `South`, `East`, `West`, `Up`, `Down`
- `Room` — fields: `ID string`, `Description string`, `Exits map[Direction]string`
  (exits map a direction to a room ID, not a pointer)
- `Object` — fields: `ID string`, `Name string`
- `World` — unexported fields: `rooms map[string]*Room`, `objects map[string]*Object`

## Constructors to implement

- `NewRoom(id, description string) *Room`
- `NewObject(id, name string) *Object`
- `NewWorld() *World`

## Methods to implement on `World`

- `AddRoom(r *Room)`
- `AddObject(o *Object)`
- `ConnectRooms(fromID string, dir Direction, toID string)`
- `RoomByID(id string) (*Room, bool)`
- `ObjectByID(id string) (*Object, bool)`

## Doc comments

Every exported type, constructor, and method needs a doc comment.
They will appear in `go doc github.com/matteo-gildone/goforth`.
Write them as complete sentences starting with the symbol name:

    // Room represents a single location in the game world.
    // Room is connected to other rooms via directional exits.

## Tests

Table-driven, covering happy path and not-found cases for lookups,
and connection via `ConnectRooms`.

## Design notes

Exits store room IDs as strings, not pointers. World is the only place that
resolves a room reference. This keeps the graph serialisable and avoids circular
pointer issues. Unexported fields require constructors — consumers of
`github.com/matteo-gildone/goforth` cannot use struct literals directly.

---

# Card 2 — Object location tracking

Add location tracking to World. Objects are not owned by rooms — their location
is a separate concern tracked centrally.

## Add to `World`

- `objectLocations map[string]string` — unexported. Maps object ID to room ID,
  or the sentinel string `"player"` when in the player's inventory.

## Methods to implement on `World`

- `PlaceObject(objectID, roomID string) error`
- `MoveObjectToPlayer(objectID string) error`
- `MoveObjectToRoom(objectID, roomID string) error`
- `PlayerHasObject(objectID string) bool`
- `ObjectsInRoom(roomID string) []*Object`

## Doc comments

Every exported method needs a doc comment. Example:

    // PlaceObject sets the initial location of an object within the world.
    // It returns an error if the object or room ID is not recognised.

## Tests

Table-driven, covering initial placement, take (room → player), drop
(player → room), unknown object ID, unknown room ID.

## Design notes

This model makes take/drop a single map update rather than splicing slices on
two structs. The Object type stays immutable.

---

# Card 3 — Player

Define the Player type. Kept deliberately thin — Player has no knowledge of World.

## Type to define

- `Player` — unexported fields: `currentRoomID string`

## Constructor to implement

- `NewPlayer(startingRoomID string) *Player`

## Methods to implement

- `MoveTo(roomID string)` — updates currentRoomID
- `CurrentRoomID() string` — accessor for the unexported field

## Doc comments

Every exported symbol needs a doc comment. Example:

    // Player represents the person playing the game.
    // Player has no knowledge of the world — movement validation
    // is the responsibility of the game loop.

## Tests

Minimal for the type itself. Interesting behaviour comes in integration (Card 7).

## Design notes

Player does not resolve rooms, validate exits, or touch World. All of that is
the game loop's responsibility, informed by World. Unexported field means
currentRoomID is set only via MoveTo and read via CurrentRoomID.

---

# Card 4 — Parser

Implement a parser completely decoupled from command handlers and world state.

## Types to define

- `Command` — exported fields: `Name string`, `Args []string`

## Functions to implement

- `Parse(line string) Command` — trim whitespace, lowercase, split on whitespace.
  First token becomes `Name`, remainder becomes `Args`. Empty or whitespace-only
  input returns zero `Command` (empty Name, nil Args).

## Doc comments

    // Command represents a parsed player instruction.
    // Name is the command verb. Args are the remaining tokens.

    // Parse tokenises a raw input line into a Command.
    // Input is trimmed, lowercased, and split on whitespace.
    // An empty or whitespace-only line returns a zero Command.

## Tests

Table-driven, covering empty string, whitespace-only, single word, multi-word,
extra internal spaces, all-caps input, mixed case.

## Design notes

The parser knows nothing about which commands exist. It only tokenises.
Validation and dispatch happen elsewhere.

---

# Card 5 — CommandRegistry

Implement the extensibility mechanism. Adding a new command to a game built on
`github.com/matteo-gildone/goforth` requires only a registration call —
nothing else changes.

## Types to define

- `HandlerFunc` — `func(args []string, g *Game) error`
- `CommandRegistry` — unexported field: `handlers map[string]HandlerFunc`

## Constructor to implement

- `NewCommandRegistry() *CommandRegistry`

## Methods to implement on `CommandRegistry`

- `Register(name string, h HandlerFunc)`
- `Dispatch(cmd Command, g *Game) error` — looks up handler by `cmd.Name`,
  calls it with `cmd.Args`. If name not found, prints
  `"I don't know how to do that."` and returns `nil`.

## Sentinel to define

- `ErrQuit = errors.New("quit")` — returned by the quit handler to signal
  clean loop exit. Mirrors the `io.EOF` pattern — a known terminal condition,
  not a bug.

## Doc comments

    // HandlerFunc is the type for all command handlers.
    // args are the tokens following the command name.
    // Returning ErrQuit signals the game loop to exit cleanly.

    // CommandRegistry maps command names to their handlers.
    // Consumers register their own commands alongside the built-in ones.

    // ErrQuit is returned by a handler to signal a clean game exit.
    // The game loop treats ErrQuit as a normal terminal condition, not an error.

## Tests

Table-driven, covering successful dispatch, unknown command returns nil,
`ErrQuit` propagates correctly.

## Design notes

Unknown command is not an error — it is normal gameplay feedback. Only `ErrQuit`
and unexpected errors should stop the game loop.

---

# Card 6 — Built-in handlers

Implement all built-in command handlers and register direction aliases.
These are the default verbs any game built on `github.com/matteo-gildone/goforth`
gets out of the box.

## Handlers to implement

- `LookHandler` — prints current room description and available exits
- `GoHandler` — resolves direction from `args[0]` if present, otherwise from
  `cmd.Name` itself. Moves player if exit exists, prints feedback if not.
- `TakeHandler` — moves named object from current room to player inventory
- `DropHandler` — moves named object from player inventory to current room
- `InventoryHandler` — lists objects currently held by player
- `QuitHandler` — returns `ErrQuit`

## Direction aliasing

Register direction words and short forms as wrappers around `GoHandler` so that
`north`, `go north`, and `n` all work. Iterate an alias map at registration
time, registering each as a closure:
```
"go"    → GoHandler
"north" → func(args, g) { GoHandler([]string{"north"}, g) }
"n"     → func(args, g) { GoHandler([]string{"north"}, g) }
// repeat for all directions
```

## Doc comments

Each exported handler needs a doc comment. Example:

    // LookHandler prints the description and exits of the player's current room.
    // Register it as "look" in a CommandRegistry.

## Tests

Per handler, covering happy path, missing args, invalid state (no exit that way,
object not in room, object not in inventory).

## Design notes

Handlers are exported so consumers can compose them — a consumer might want to
wrap LookHandler to add extra output, or call GoHandler programmatically during
setup.

---

# Card 7 — Game and game loop

Assemble the engine. The game loop is separate from the parser and separate
from the handlers.

## Type to define

- `Game` — fields: `World *World`, `Player *Player`, `Registry *CommandRegistry`

## Constructor to implement

- `NewGame(w *World, p *Player, r *CommandRegistry) *Game`

## Methods to implement

- `Game.Run(r io.Reader)` — reads lines from `r` in a loop, calls `Parse`,
  calls `Registry.Dispatch`. Stops on `ErrQuit` or an unexpected error.
  Using `io.Reader` rather than `os.Stdin` directly makes this testable.

## Loop shape
```
scanner loop:
  read line
  parse → Command
  dispatch → error
  if err == ErrQuit → break
  if err != nil → handle unexpected error
```

## Doc comments

    // Game holds the complete state of a running game session.
    // World, Player, and Registry are intentionally exported so consumers
    // can inspect and extend state between commands if needed.

    // Run starts the game loop, reading commands from r until the player quits
    // or an unexpected error occurs. Pass os.Stdin for interactive play.

## Tests

Integration test: construct a small world using the public API, feed a scripted
sequence of commands via `strings.NewReader`, assert final player position and
inventory state.

---

# Card 8 — Example consumer

Prove the public API is usable and validate the under-50-lines constraint.
This is the definition of done for the engine.

## What to build

`example/main.go` — a standalone program that imports
`github.com/matteo-gildone/goforth` and wires a playable world:

- At least four rooms connected in multiple directions
- At least two objects placed in different rooms
- All built-in handlers registered
- `game.Run(os.Stdin)` called to start

## Definition of done

- No engine logic in the example — only calls to the public API
- A developer unfamiliar with the codebase can read `example/main.go`
  and understand the full shape of a game world in under a minute
- Total line count is under 50
- `go doc github.com/matteo-gildone/goforth` produces clean, readable output
  for every exported symbol