# 🚀 Getting Started - Connect-4 Game Server

## ✅ Server is Already Running!

Your Connect-4 game server is **currently running** on your computer.

---

## 🎮 Play Now!

### Step 1: Open Your Browser
Click this link or type in address bar:
```
http://localhost:8080
```

### Step 2: Enter Your Username
Type any name you like:
- "Alice"
- "Bob"
- "Player1"
- Or anything else!

### Step 3: Click "Join Matchmaking"
The button is right there on the page.

### Step 4: Wait or Play
- **Wait 10 seconds** → Bot will join automatically
- **Or** → Have a friend open another browser window and join within 10 seconds

### Step 5: Play!
Click any column to drop your disc. Connect 4 in a row to win!

---

## 🎯 Game Rules

### The Board
- 7 columns wide
- 6 rows tall
- 42 total spaces

### How to Win
Connect 4 of your discs in a row:
- **Horizontal** → Left to right
- **Vertical** → Top to bottom
- **Diagonal** → Any direction

### Turn Order
- Red (Player 1) always goes first
- Players alternate turns
- Click any column to drop your disc

### Draw
If the board fills up with no winner, it's a draw.

---

## 🤖 Playing Against the Bot

The bot is **smart** - it will:
1. Try to win if it can
2. Block your winning moves
3. Play strategically
4. Prefer center columns

It's not random, so it's a real challenge!

---

## 👥 Playing Against a Friend

### Same Computer
1. Open first browser window
2. Enter "Alice" and join
3. Open second browser window (or incognito)
4. Enter "Bob" and join (within 10 seconds)
5. Play!

### Different Computers
1. Get your computer's IP address
2. Share: `http://YOUR_IP:8080`
3. Friend opens that URL
4. Both join matchmaking
5. Play!

---

## 🏆 Leaderboard

The leaderboard shows:
- Player names
- Number of wins
- Updated in real-time

Your wins are tracked automatically!

---

## 🔄 Reconnection

If you accidentally close the browser:
1. You have **30 seconds** to reconnect
2. Open the same URL
3. Enter the same username
4. You'll rejoin the same game!

If you don't reconnect in time, your opponent wins by forfeit.

---

## 🛠️ Troubleshooting

### "Can't connect to http://localhost:8080"

**Solution 1**: Try this instead
```
http://127.0.0.1:8080
```

**Solution 2**: Check if server is running
- Look for PowerShell window with server
- Should say "Server starting on :8080"

**Solution 3**: Try different port
- Server might be on 8081 or 8082
- Check the PowerShell window for the port number

### "WebSocket connection failed"

**Solution**:
1. Refresh the page (F5)
2. Close and reopen browser
3. Check Windows Firewall
4. Try a different browser

### "Can't see opponent's moves"

**Solution**:
1. Refresh the page
2. Make sure both players are in the same game
3. Check browser console (F12) for errors

### "Bot isn't playing"

**Solution**:
1. Wait 10 seconds (it takes time to match)
2. Refresh the page
3. Try joining again

---

## 📊 What's Happening Behind the Scenes

### When You Join
1. Your username is added to the queue
2. Server waits for another player
3. After 10 seconds, bot joins if no one else did
4. Game starts!

### When You Play
1. You click a column
2. Your move is sent to the server
3. Server validates the move
4. Board updates for both players
5. Server checks for win/draw
6. Opponent's turn

### When Game Ends
1. Winner is determined
2. Game is saved to file
3. Leaderboard is updated
4. You can play again!

---

## 💾 Your Data

All your games are saved:
- **Games**: `data/games.json`
- **Leaderboard**: `data/leaderboard.json`
- **Events**: `data/events.jsonl`

These files are created automatically in the project folder.

---

## 🎮 Game Tips

### Winning Strategy
1. **Control the center** - More ways to win
2. **Block opponent** - Don't let them get 3 in a row
3. **Create threats** - Get 2-3 in a row to force blocks
4. **Think ahead** - Plan 2-3 moves ahead

### Against the Bot
1. The bot blocks your wins
2. The bot tries to win
3. The bot prefers center columns
4. Try to create multiple threats

---

## 🔧 Advanced Options

### Change Server Port
If port 8080 is busy:
```powershell
$env:SERVER_PORT="8081"
go run ./server
```

Then open: `http://localhost:8081`

### Enable Database
For persistent storage:
```powershell
$env:DB_ENABLED="true"
$env:DB_HOST="localhost"
$env:DB_USER="postgres"
$env:DB_PASSWORD="password"
go run ./server
```

### Enable Kafka Analytics
For event streaming:
```powershell
$env:KAFKA_ENABLED="true"
$env:KAFKA_BROKERS="localhost:9092"
go run ./server
```

---

## 📚 Documentation

For more information:
- **README.md** - Complete documentation
- **QUICKSTART.md** - 5-minute guide
- **API.md** - Technical API details
- **SETUP.md** - Detailed setup
- **TODO.md** - Step-by-step checklist

---

## 🎯 Quick Reference

| Action | How |
|--------|-----|
| Play | Open http://localhost:8080 |
| Stop Server | Press Ctrl+C in PowerShell |
| View Leaderboard | Scroll down on game page |
| View Games | Open `data/games.json` |
| View Events | Open `data/events.jsonl` |
| Change Port | Set `$env:SERVER_PORT` |
| Enable Database | Set `$env:DB_ENABLED="true"` |
| Enable Kafka | Set `$env:KAFKA_ENABLED="true"` |

---

## 🎉 You're Ready!

Everything is set up and ready to go!

### Next Steps
1. ✅ Open http://localhost:8080
2. ✅ Enter your username
3. ✅ Click "Join Matchmaking"
4. ✅ Play!

---

## 🏆 Have Fun!

Your Connect-4 game server is ready to play. Enjoy! 🎮

---

**Questions?** Check the documentation files or the troubleshooting section above.

**Ready to play?** Open http://localhost:8080 now!

