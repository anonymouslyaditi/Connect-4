# API Documentation - 4 in a Row Game Server

## HTTP Endpoints

### GET /
Serves the static frontend (HTML/CSS/JS).

**Response:** HTML page

---

### GET /leaderboard
Returns the current leaderboard with player rankings.

**Response:**
```json
{
  "alice": 5,
  "bob": 3,
  "charlie": 2
}
```

**Status Codes:**
- `200 OK` - Success

---

## WebSocket Endpoint

### WS /ws
WebSocket connection for real-time game communication.

**Connection URL:** `ws://localhost:8080/ws` (or `wss://` for secure)

---

## WebSocket Messages

### Client → Server

#### 1. Join Game
Join matchmaking or reconnect to existing game.

**Message:**
```json
{
  "type": "join",
  "username": "alice",
  "gameId": "g_1729765800000000000"  // Optional, for reconnection
}
```

**Fields:**
- `type` (string, required): Must be "join"
- `username` (string, required): Player's username
- `gameId` (string, optional): Game ID for reconnection

**Response:** Server will send `waiting`, `start`, or `reconnected` message

---

#### 2. Make Move
Drop a disc into a column.

**Message:**
```json
{
  "type": "move",
  "gameId": "g_1729765800000000000",
  "col": 3
}
```

**Fields:**
- `type` (string, required): Must be "move"
- `gameId` (string, required): Current game ID
- `col` (integer, required): Column number (0-6)

**Response:** Server will send `state` message with updated board

---

### Server → Client

#### 1. Waiting for Opponent
Sent when player is in matchmaking queue.

**Message:**
```json
{
  "type": "waiting",
  "timeout": 10
}
```

**Fields:**
- `type` (string): "waiting"
- `timeout` (integer): Seconds until bot game starts

---

#### 2. Game Started
Sent when game begins (matched with player or bot).

**Message:**
```json
{
  "type": "start",
  "gameId": "g_1729765800000000000",
  "you": 1,
  "opponent": "bob",
  "state": {
    "rows": 6,
    "cols": 7,
    "board": [
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0]
    ],
    "turn": 1,
    "started": "2025-10-24T10:28:00Z"
  }
}
```

**Fields:**
- `type` (string): "start"
- `gameId` (string): Unique game identifier
- `you` (integer): Your player number (1 or 2)
- `opponent` (string): Opponent's username (or "Bot")
- `state` (object): Current game state
  - `rows` (integer): Number of rows (6)
  - `cols` (integer): Number of columns (7)
  - `board` (array): 2D array representing the board
    - `0` = empty
    - `1` = player 1's disc
    - `2` = player 2's disc
  - `turn` (integer): Current player's turn (1 or 2)
  - `started` (string): Game start timestamp (ISO 8601)

---

#### 3. Game State Update
Sent after each move or game event.

**Message:**
```json
{
  "type": "state",
  "gameId": "g_1729765800000000000",
  "state": {
    "rows": 6,
    "cols": 7,
    "board": [
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0, 0, 0],
      [0, 0, 0, 1, 0, 0, 0]
    ],
    "turn": 2,
    "started": "2025-10-24T10:28:00Z"
  },
  "you": 1,
  "status": "playing"
}
```

**Fields:**
- `type` (string): "state"
- `gameId` (string): Game identifier
- `state` (object): Updated game state
- `you` (integer): Your player number
- `status` (string): Game status
  - `"playing"` - Game in progress
  - `"finished"` - Game ended

---

#### 4. Reconnected
Sent when player successfully reconnects to a game.

**Message:**
```json
{
  "type": "reconnected",
  "gameId": "g_1729765800000000000",
  "state": {
    "rows": 6,
    "cols": 7,
    "board": [...],
    "turn": 2,
    "started": "2025-10-24T10:28:00Z"
  }
}
```

**Fields:**
- `type` (string): "reconnected"
- `gameId` (string): Game identifier
- `state` (object): Current game state

---

#### 5. Error
Sent when an error occurs.

**Message:**
```json
{
  "error": "first message must be join"
}
```

**Fields:**
- `error` (string): Error description

---

## Kafka Events

Events are published to the configured Kafka topic (default: `game-analytics`).

### Event: game_finished

**Message:**
```json
{
  "type": "game_finished",
  "timestamp": "2025-10-24T10:30:00Z",
  "game": {
    "id": "g_1729765800000000000",
    "player1": "alice",
    "player2": "bob",
    "winner": "alice",
    "duration_seconds": 120,
    "started_at": "2025-10-24T10:28:00Z",
    "ended_at": "2025-10-24T10:30:00Z"
  },
  "reason": "win"
}
```

**Fields:**
- `type` (string): "game_finished"
- `timestamp` (string): Event timestamp (ISO 8601)
- `game` (object): Game details
  - `id` (string): Game identifier
  - `player1` (string): First player's username
  - `player2` (string): Second player's username (or "Bot")
  - `winner` (string): Winner's username or "draw"
  - `duration_seconds` (integer): Game duration in seconds
  - `started_at` (string): Game start time
  - `ended_at` (string): Game end time
- `reason` (string): How game ended
  - `"win"` - Normal win
  - `"draw"` - Board full, no winner
  - `"forfeit"` - Player disconnected

---

### Event: move

**Message:**
```json
{
  "type": "move",
  "timestamp": "2025-10-24T10:29:00Z",
  "gameId": "g_1729765800000000000",
  "col": 3,
  "player": 1
}
```

**Fields:**
- `type` (string): "move"
- `timestamp` (string): Event timestamp (ISO 8601)
- `gameId` (string): Game identifier
- `col` (integer): Column where disc was dropped (0-6)
- `player` (integer): Player who made the move (1 or 2)

---

## Game Logic

### Board Representation
- 6 rows × 7 columns
- Row 0 is the top, Row 5 is the bottom
- Column 0 is the left, Column 6 is the right
- Discs fall to the lowest available position in a column

### Win Conditions
A player wins by connecting 4 discs in a row:
- **Horizontal:** 4 consecutive discs in the same row
- **Vertical:** 4 consecutive discs in the same column
- **Diagonal:** 4 consecutive discs diagonally (↗ or ↘)

### Draw Condition
If all 42 positions are filled and no player has won, the game is a draw.

### Turn Order
- Player 1 always goes first
- Players alternate turns
- Bot responds immediately when it's their turn

---

## Bot Strategy

The bot uses the following decision algorithm:

1. **Check for winning move** - If bot can win this turn, take it
2. **Block opponent's win** - If opponent can win next turn, block them
3. **Strategic positioning** - Prefer center columns (better winning opportunities)
4. **Fallback** - Choose any valid column

The bot is **not random** and plays competitively.

---

## Error Handling

### Connection Errors
- If WebSocket connection fails, client should retry
- Server logs connection errors

### Invalid Moves
- Moving in a full column: Ignored
- Moving out of turn: Ignored
- Invalid column number: Ignored

### Disconnection
- Player has 30 seconds (configurable) to reconnect
- After timeout, opponent wins by forfeit
- Game is saved with forfeit reason

---

## Rate Limiting

Currently, no rate limiting is implemented. In production, consider:
- Limiting connections per IP
- Limiting moves per second
- Limiting concurrent games per user

---

## Security Considerations

For production deployment:
1. **Use WSS** (WebSocket Secure) with TLS/SSL
2. **Validate usernames** (length, characters, profanity filter)
3. **Implement authentication** (JWT tokens, sessions)
4. **Add rate limiting** to prevent abuse
5. **Sanitize inputs** to prevent injection attacks
6. **Use CORS** properly for cross-origin requests

---

## Example Client Implementation

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  // Join game
  ws.send(JSON.stringify({
    type: 'join',
    username: 'alice'
  }));
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  
  if (msg.type === 'start') {
    console.log('Game started!', msg.gameId);
    // Make a move
    ws.send(JSON.stringify({
      type: 'move',
      gameId: msg.gameId,
      col: 3
    }));
  } else if (msg.type === 'state') {
    console.log('Board updated:', msg.state.board);
    if (msg.status === 'finished') {
      console.log('Game finished!');
    }
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Connection closed');
};
```

---

## Testing with curl

You can test the HTTP endpoint:

```bash
# Get leaderboard
curl http://localhost:8080/leaderboard
```

For WebSocket testing, use tools like:
- [websocat](https://github.com/vi/websocat)
- [wscat](https://github.com/websockets/wscat)
- Browser DevTools Console

---

**For more information, see README.md and SETUP.md**

