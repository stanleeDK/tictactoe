# Triple T Tantilizer

An unbeatable tic-tac-toe opponent, written in Go with no external dependencies.

Every time you move, the server enumerates **the entire remaining game** into an
explicit move tree — roughly 60,000 board states on your first move — scores it with
minimax, and plays the best reply. You can then click the node counter to render that
whole tree in your browser and see exactly why the computer chose what it chose.

```
    x | - | -          Your move  →  server builds 59,705 future board states
   ---+---+---                    →  minimax scores every one of them
    - | o | -                     →  plays the highest-scoring reply
   ---+---+---
    - | - | -
```

---

## Contents

- [Features](#features)
- [Quick start](#quick-start)
- [How it works](#how-it-works)
- [Project layout](#project-layout)
- [HTTP API](#http-api)
- [Configuration](#configuration)
- [Visualizing the move tree](#visualizing-the-move-tree)
- [Known limitations](#known-limitations)

---

## Features

- **Unbeatable play** — full-depth minimax over the complete game tree. No heuristics,
  no opening book, no randomness. The best you can do is draw.
- **Depth-aware scoring** — among equally winning lines the engine prefers the *fastest*
  win and the *slowest* loss, so it closes games out instead of dawdling.
- **Live progress streaming** — tree construction pushes its node count to the browser
  over Server-Sent Events, so you watch the search space grow in real time.
- **Inspectable search** — the scored tree can be rendered in the browser as an
  interactive chart, one 3×3 board per node, annotated with tree depth, node index,
  static board score, and final minimax score.
- **Multi-player sessions** — each browser gets its own mutex-guarded session with its
  own board and its own prediction tree, so concurrent games don't interfere.
- **Zero dependencies** — standard library only. `go.mod` requires nothing.

## Quick start

Requires Go 1.18 or newer.

```bash
git clone https://github.com/stanleeDK/tictactoe.git
cd tictactoe

# development mode — binds localhost:8888
GO_ENV=development go run .
```

Then open <http://localhost:8888>. You are always `x` and you always move first.

To build a binary instead:

```bash
go build -o ttt .
GO_ENV=development ./ttt
```

> **Note:** the server resolves `staticfiles/` relative to the working directory, so run
> the binary from the repository root.

## How it works

The interesting part is that the move tree is a **real, materialized data structure** —
not an implicit recursion — built breadth-first with an explicit queue.

### 1. Build the tree — `tree.go`, `BuildMoveTreeBreadthFirst()`

Starting from the current board, the root is enqueued and then:

1. Dequeue a node and look at whose turn it is (`MostRecentPlayer` alternates by level).
2. For every empty cell, produce one child board with that move applied.
3. Evaluate each child (`ReturnEvaluationOfGameBoard`) and give it a static score.
4. Enqueue only the children whose games are still in progress; terminal boards stay leaves.

Two copies of the parent board are kept per iteration, and the distinction matters:
`tempBoardIterator` is a scratch board that accumulates moves purely to drive the loop
over the remaining empty cells, while `originalBoardState` is the pristine parent that
each child is derived from. Without the second copy, every child would inherit the
previous sibling's move.

Each node created pushes the running node count into a buffered channel with a
`select`/`default`, so a slow or absent listener never blocks the search.

### 2. Score the leaves — `tree.go`, `findNextMove()`

| Outcome | Static score |
| --- | --- |
| Computer (`o`) wins | `10 - depth` |
| Human (`x`) wins | `-10 + depth` |
| Draw | `0` |
| Game not finished | `0` |

Depth is folded into the score so shallow wins beat deep ones, and a loss three moves
away is treated as less bad than a loss on the next move.

### 3. Minimax — `tree.go`, `Maximizer()` / `Minimizer()`

`Maximizer` and `Minimizer` recurse into each other down alternating levels. A leaf
returns its static score; an internal node returns the max (computer's turn) or min
(human's turn) of its children, and caches it on the node as `MiniMaxScore`.

### 4. Pick the move — `tree.go`, `SearchMiniMaxedTreeForNextComputerMove()`

With the tree scored, the computer's move is simply the root child with the highest
`MiniMaxScore`. That child's board becomes the new game board and is rendered to HTML
for the client.

### Request flow

```
  browser                      server
     │
     │  POST /startgame/                      GameSessions.CreateNewSession()
     ├──────────────────────────────────────► new board + empty prediction tree
     │  ◄── {"sessionid": "…"}
     │
     │  POST /processHumanMove/               apply move → evaluate
     ├──────────────────────────────────────► if game continues:
     │                                          BuildMoveTreeBreadthFirst()  (goroutine)
     │  GET /getprogress/   (EventSource)       Minimax()
     ├──────────────────────────────────────►   SearchMiniMaxedTreeForNextComputerMove()
     │  ◄── data: {"payloadtype":"41273", …}    ↕ node counts over SSE
     │  ◄── data: {"progress":"complete"}
     │
     │  GET /getCurrentBoardState/           prerendered <table> of the board
     ├──────────────────────────────────────►
     │  ◄── <table>…</table>
     │
     │  GET /showgametree/                   whole scored tree as Treant JSON
     ├──────────────────────────────────────►
     │  ◄── {"chart":…,"nodeStructure":…}
```

The server renders each board to an HTML `<table>` and the client scrapes the cell
values out of it — an unusual choice, but it means one rendering path serves both the
playable board and the debug tree.

## Project layout

| File | Responsibility |
| --- | --- |
| `main.go` | Entry point; route registration and server startup |
| `views.go` | Renders `staticfiles/index.html` for `GET /` |
| `gameplayfacilitationcontroller.go` | All HTTP handlers; orchestrates move → build tree → minimax → reply |
| `gamesession.go` | Per-browser `GameSession` and the mutex-guarded `GameSessions` map |
| `boardinstance.go` | The 3×3 board: win/draw detection, scoring fields, HTML rendering |
| `tree.go` | `Node`/`Tree`, breadth-first tree construction, minimax, Treant export |
| `queue/queue.go` | Minimal FIFO queue backing the breadth-first build |
| `treantchart/treantchart.go` | Structs that marshal to the JSON shape Treant.js expects |
| `utilityfuncs.go` | Console board printing and `isBoardFull` helper |
| `staticfiles/index.html` | The entire front end: board UI, SSE client, tree rendering |
| `staticfiles/playground.html` | Standalone CSS confetti experiment, not wired into the game |

### Board representation

`BoardInstance.Board` is a `[3][3]byte` using `'x'`, `'o'`, and `'-'` for empty. Game
state is the `GameState` enum: `Draw`, `ComputerIsOAndWon`, `HumanisXandWon`,
`GameNotFinished`. Wins are detected by three independent scans — horizontal, vertical,
and both diagonals.

## HTTP API

All handlers are registered with a trailing slash; the browser calls them without one
and relies on `http.ServeMux`'s redirect.

| Method | Path | Parameters | Returns |
| --- | --- | --- | --- |
| `GET` | `/` | — | The game page |
| `POST` | `/startgame/` | body ignored | `{"sessionid":"<32 hex chars>"}` |
| `POST` | `/processHumanMove/` | JSON `{cellIndex,value,sessionid}` | Empty; board is updated server-side |
| `GET` | `/getCurrentBoardState/` | `?sessionid=` | HTML `<table>` fragment of the current board |
| `GET` | `/getprogress/` | `?sessionid=` | `text/event-stream` of node counts |
| `GET` | `/restartGame/` | `?sessionid=` | Empty; resets board and drops the tree |
| `GET` | `/showgametree/` | `?sessionid=` | Treant-format JSON of the scored tree |
| `GET` | `/staticfiles/…` | — | Static assets |

`cellIndex` is a two-character row/column string: `"00"` through `"22"`.

Every `sessionid`-bearing endpoint looks the session up under the sessions mutex and
returns `404 Not Found` for an unknown id. `/getprogress/` reports the same condition on
the stream instead, since the client is an `EventSource`:

```
data: {"PayLoadType":"thisisprogressdata","data":"session not found"}
```

Progress events look like this, terminating when the channel closes:

```
data: {"payloadtype":"41273","progress":"pending"}
data: {"payloadtype":"complete","progress":"complete"}
```

## Configuration

One environment variable:

| `GO_ENV` | Listen address |
| --- | --- |
| `development` | `localhost:8888` |
| anything else (including unset) | `0.0.0.0:80`, with 60s write and 120s idle timeouts |

Binding port 80 needs elevated privileges, so use `GO_ENV=development` locally.

## Visualizing the move tree

After the computer replies, the node counter under the board ("Moves being considered")
is clickable. Clicking it fetches the scored tree and renders it with
[Treant.js](https://fperucic.github.io/treant-js/), loaded from a CDN along with Raphaël.

Each node shows the next player to move, whether that level is maximizing or minimizing,
its node index (`i`) and level (`l`), its static board score (`sc`), its minimax score
(`mmS`), and the board itself — colored green when the computer wins, red when the human
wins, grey for a draw.

Be aware that this is a genuinely large payload: on the first move the tree has ~59,700
nodes and the JSON runs to roughly 55 MB, because every node carries its own prerendered
HTML. It is best used a few moves into a game, once the tree has shrunk.

## Known limitations

Honest notes on the current state of the code:

- **No alpha-beta pruning.** The full tree is built and scored from scratch on every
  move rather than pruned or reused, which is exactly why the first reply is the slow
  one. It is also the point: the tree exists so it can be looked at.
- **Sessions are never reaped.** `GameSessionsRunning` grows for the lifetime of the
  process; abandoned games are kept in memory along with their trees.
- **Fixed roles.** The human is always `x` and moves first; the computer is always `o`.
  The plumbing for choosing a side is still visible in the code but unused.
- **All state is in memory.** Restarting the server drops every game in progress.
- **A compiled `ttt` binary is committed** to the repository.
