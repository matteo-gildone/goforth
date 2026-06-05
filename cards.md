Card 1: Bug -- Return scanner.Err() from Game.Run
File: game.go:33
Problem: If the reader fails mid-stream (I/O error, buffer overflow), Game.Run returns nil. The scanner error is silently swallowed.
What to do:
- Change return nil to return scanner.Err()
  Definition of done:
- Game.Run returns the scanner error when the reader fails
- A test passes a failing reader (e.g. iotest.ErrReader) and asserts the error is returned
- All existing tests still pass
---
Card 2: Cleanup -- Remove unnecessary loop variable copy
File: handlers.go:177
Problem: d := dir was needed before Go 1.22 to avoid closure-over-loop-variable bugs. Since Go 1.22 (you're on 1.25), loop variables are per-iteration. The copy is dead code that signals unfamiliarity with the change.
What to do:
- Remove d := dir
- Capture dir directly in the closure
  Definition of done:
- go test passes
- go vet clean
---
Card 3: Cleanup -- Replace Exit struct with plain string
File: room.go:4-6
Problem: Exit wraps a single string field (RoomID). It adds indirection without buying an extension point that's currently used. Every access becomes exit.RoomID instead of just exit.
What to do:
- Change exits map[Direction]Exit to exits map[Direction]string
- Update North, South, East, West, Up, Down to store to.ID directly
- Update Exit() method to return from the map directly
- Remove the Exit struct
  Definition of done:
- go test passes
- Room direction methods and Exit()/ExitDirections() behave identically
- No Exit struct remains
---
Card 4: Refactor -- Extract currentRoom helper to replace repeated room-existence checks
Files: handlers.go
Problem: Four handlers (LookHandler, GoHandler, TakeHandler, DropHandler) repeat the same 3-line room-existence check. The check is intentional (defensive against corrupt player state), but the repetition is boilerplate.
What to build:
An unexported helper:
func currentRoom(g *Game) (*Room, error) {
room, ok := g.World.RoomByID(g.Player.CurrentRoom())
if !ok {
return nil, &RoomNotFoundErr{ID: g.Player.CurrentRoom()}
}
return room, nil
}
Replace all inline room-existence checks in handlers with calls to currentRoom(g).
GoHandler calls it twice (before move to validate exit, after move to trigger OnEnter) -- both stay.
Definition of done:
- No inline RoomByID + if !ok pattern remains in handlers
- All handlers that need the current room use the shared helper
- go test passes
- Behaviour identical to before
---
Card 5: Feature -- Implement GameBuilder
Problem: Creating a game requires 12 steps with implicit ordering constraints. The consumer must manually wire World, Player, CommandRegistry, and Game together every time.
What to build:
New file: builder.go
Build() *GameBuilder
(*GameBuilder) AddRooms(rooms ...*Room) *GameBuilder
(*GameBuilder) AddObjects(objects ...*Object) *GameBuilder
(*GameBuilder) PlaceObject(room *Room, objects ...*Object) *GameBuilder
(*GameBuilder) StartRoom(room *Room) *GameBuilder
(*GameBuilder) Output(w io.Writer) *GameBuilder
(*GameBuilder) Done() (*Game, error)
Design decisions (settled):
- PlaceObject takes pointers, not strings -- matches the room-connection API, eliminates string typos
- StartRoom takes *Room pointer
- AddObjects stays separate from PlaceObject -- objects can exist unplaced (future loot drops)
- Done() registers default handlers automatically
- Room creation and connection stays as-is (method chaining on *Room)
- Existing low-level API (NewWorld, NewPlayer, etc.) remains available
  Validation in Done():
- Start room is set and exists in the added rooms
- All objects in PlaceObject calls exist in the added objects
- All rooms in PlaceObject calls exist in the added rooms
- Output writer is set
- At least one room exists
  Definition of done:
- example/main.go rewritten to use builder
- Builder tests cover: happy path, missing start room, invalid placement, missing output
- go test passes
- Existing low-level API still works (not removed)
---
Card 6: Cleanup -- Remove scratch code from handlers.go
File: handlers.go:201-208
Problem: Builder pseudocode pasted at end of file. Not valid Go, breaks compilation if uncommented.
What to do: Delete lines 201-208.
Definition of done: go build ./... passes.
---
Card 7: Refactor -- Move player inventory ownership out of World
Problem: World manages player inventory via the magic string "player" in objectLocations. This means:
- World is responsible for two domains (spatial placement + player possession)
- A room named "player" would collide with the sentinel
- PlayerInventory() scans the entire location map every call
- Object ownership is a player concern living in the wrong layer
  What to build:
- Player gets an inventory (e.g. map[string]bool or []string)
- New methods on Player: AddToInventory(objectID string), RemoveFromInventory(objectID string), HasObject(objectID string) bool, Inventory() []string
- Remove from World: MoveObjectToPlayer, PlayerHasObject, PlayerInventory
- Handlers update to use Player for inventory operations and World only for room-level placement
- TakeHandler: removes object from world locations, adds to player inventory
- DropHandler: removes from player inventory, places in world at current room
  Definition of done:
- No magic "player" string in the codebase
- World has no knowledge of player possession
- Player owns inventory state
- All existing tests pass (rewritten to use new API)
- A room with ID "player" does not break anything
---
Card 8: Refactor -- Consolidate test setup helpers
Problem: Three different helpers (setupGame, setupAliasWorld in game_test.go and newTestGame in handlers_test.go) do variations of the same thing. setupWorld in world_test.go is a fourth.
Depends on: Card 5 (builder) -- once the builder exists, test setup becomes trivial.
What to do:
- After builder is implemented, rewrite test setup to use Build().AddRooms(...).Done()
- Or: unify into a single newTestGame helper with clear parameters
- Remove setupGame, setupAliasWorld, setupWorld duplication
  Definition of done:
- One canonical way to set up a test game
- All tests pass
- No duplicate setup logic across test files
  Card 9: Feature -- Add error accumulation to Room direction methods (pipeline pattern)
  File: room.go
  Problem: Room direction methods silently overwrite exits. If a consumer makes contradictory connections (e.g. entrance.North(dining) followed by entrance.North(kitchen)), or a bidirectional conflict occurs, the graph becomes inconsistent with no feedback.
  What to build:
- Add err error field to Room struct (unexported)
- Each direction method (North, South, East, West, Up, Down) checks if that direction is already assigned on the receiver or if the reverse direction is already assigned on the target
- If conflict detected: store the error, return r without mutating (no-op)
- Once errored, subsequent direction calls on that room are no-ops (same as bufio.Scanner)
- Add Err() error method on *Room that returns the stored error
  func (r *Room) Err() error
  Integration with builder (Card 5):
  Done() iterates all added rooms and calls Err() on each. If any room has accumulated an error, Done() returns it. The consumer never needs to check Err() manually if they use the builder.
  Consumers not using the builder check manually:
  entrance.North(dining).East(sport)
  if err := entrance.Err(); err != nil {
  log.Fatal(err)
  }
  Error message format: Something like "room %q: exit %s already assigned"
  Definition of done:
- Direction methods are no-ops after first conflict
- Err() returns nil when no conflict exists
- Err() returns a descriptive error on conflict
- Bidirectional conflicts are caught (target's reverse direction already set)
- Builder's Done() checks all rooms for errors
- Tests cover: no conflict, direct conflict (same direction twice), bidirectional conflict, chaining after error is a no-op
- Fluent API (entrance.North(dining).East(sport)) still works unchanged
---
Updated ordering:
1. Card 6 -- remove scratch code from handlers.go
2. Card 1 -- return scanner.Err()
3. Card 2 -- remove loop variable copy
4. Card 3 -- replace Exit struct with string
5. Card 4 -- extract currentRoom helper
6. Card 9 -- room error accumulation (pipeline pattern)
7. Card 5 -- GameBuilder (depends on Card 9 for Done() room validation)
8. Card 7 -- move inventory to Player
9. Card 8 -- consolidate test helpers (last, benefits from everything else)
---