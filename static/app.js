(function(){
const wsUrl = (location.protocol==='https:'?'wss':'ws')+'://'+location.host+'/ws'
let ws
let gameState = null
let myPlayer = null
let gameId = null
let opponent = null
let gameStatus = 'idle'
let currentUsername = ''
let currentRoomId = null
const status = id('status')
const gameDiv = id('game')
const lb = id('leaderboard')
const usernameInput = id('username')
const gameInfo = id('gameInfo')
const player1Info = id('player1Info')
const player2Info = id('player2Info')
const player1Name = id('player1Name')
const player2Name = id('player2Name')
const winnerAnnouncement = id('winnerAnnouncement')

// UI sections
const usernameSection = id('usernameSection')
const modeSelection = id('modeSelection')
const createRoomSection = id('createRoomSection')
const roomListSection = id('roomListSection')
const waitingInRoom = id('waitingInRoom')

// Buttons
const setUsernameBtn = id('setUsername')
const quickMatchBtn = id('quickMatch')
const createRoomBtn = id('createRoom')
const browseRoomsBtn = id('browseRooms')
const confirmCreateRoomBtn = id('confirmCreateRoom')
const cancelCreateRoomBtn = id('cancelCreateRoom')
const backToModeBtn = id('backToMode')
const roomNameInput = id('roomName')
const roomListDiv = id('roomList')
const roomInfoDiv = id('roomInfo')

// Set username and show mode selection
setUsernameBtn.onclick = () => {
  const username = usernameInput.value.trim()
  if(!username) {
    showStatus('Please enter a username', 'error')
    return
  }
  currentUsername = username
  usernameSection.style.display = 'none'
  modeSelection.style.display = 'block'
  showStatus('Choose a game mode', 'idle')
}

// Quick match - original matchmaking
quickMatchBtn.onclick = () => {
  modeSelection.style.display = 'none'
  connectQuickMatch(currentUsername)
}

// Create room
createRoomBtn.onclick = () => {
  modeSelection.style.display = 'none'
  createRoomSection.style.display = 'block'
  roomNameInput.value = currentUsername + "'s room"
  showStatus('Enter a room name', 'idle')
}

confirmCreateRoomBtn.onclick = () => {
  const roomName = roomNameInput.value.trim()
  if(!roomName) {
    showStatus('Please enter a room name', 'error')
    return
  }
  createRoomSection.style.display = 'none'
  connectCreateRoom(currentUsername, roomName)
}

cancelCreateRoomBtn.onclick = () => {
  createRoomSection.style.display = 'none'
  modeSelection.style.display = 'block'
  showStatus('Choose a game mode', 'idle')
}

// Browse rooms
browseRoomsBtn.onclick = () => {
  modeSelection.style.display = 'none'
  roomListSection.style.display = 'block'
  loadRooms()
  showStatus('Select a room to join', 'idle')
}

backToModeBtn.onclick = () => {
  roomListSection.style.display = 'none'
  modeSelection.style.display = 'block'
  showStatus('Choose a game mode', 'idle')
}

// Quick match connection
function connectQuickMatch(username){
  ws = new WebSocket(wsUrl)
  ws.onopen = ()=>{
    ws.send(JSON.stringify({type:'join', username}))
    showStatus('Connected as ' + username, 'playing')
  }
  ws.onmessage = (ev)=>{
    const msg = JSON.parse(ev.data)
    handle(msg)
  }
  ws.onclose = ()=> {
    showStatus('Disconnected from server', 'error')
    resetToStart()
  }
  ws.onerror = (err)=> {
    showStatus('Connection error', 'error')
    resetToStart()
  }
}

// Create room connection
function connectCreateRoom(username, roomName){
  ws = new WebSocket(wsUrl)
  ws.onopen = ()=>{
    ws.send(JSON.stringify({type:'create_room', username, roomName}))
    showStatus('Creating room...', 'waiting')
  }
  ws.onmessage = (ev)=>{
    const msg = JSON.parse(ev.data)
    handle(msg)
  }
  ws.onclose = ()=> {
    showStatus('Disconnected from server', 'error')
    resetToStart()
  }
  ws.onerror = (err)=> {
    showStatus('Connection error', 'error')
    resetToStart()
  }
}

// Join room connection
function connectJoinRoom(username, roomId){
  ws = new WebSocket(wsUrl)
  ws.onopen = ()=>{
    ws.send(JSON.stringify({type:'join_room', username, roomId}))
    showStatus('Joining room...', 'waiting')
  }
  ws.onmessage = (ev)=>{
    const msg = JSON.parse(ev.data)
    handle(msg)
  }
  ws.onclose = ()=> {
    showStatus('Disconnected from server', 'error')
    resetToStart()
  }
  ws.onerror = (err)=> {
    showStatus('Connection error', 'error')
    resetToStart()
  }
}

function resetToStart() {
  usernameSection.style.display = 'block'
  modeSelection.style.display = 'none'
  createRoomSection.style.display = 'none'
  roomListSection.style.display = 'none'
  waitingInRoom.style.display = 'none'
  usernameInput.value = ''
  currentUsername = ''
  currentRoomId = null
}

function handle(m){
  console.log('Received message:', m)

  if(m.type==='room_created'){
    currentRoomId = m.roomId
    waitingInRoom.style.display = 'block'
    roomInfoDiv.innerHTML = `
      <div style="background: #e8f5e9; padding: 15px; border-radius: 8px; border: 2px solid #4caf50;">
        <div style="font-size: 1.2em; font-weight: bold; margin-bottom: 10px;">Room: ${m.room.name}</div>
        <div style="color: #666;">Room ID: ${m.roomId}</div>
      </div>
    `
    showStatus('✅ Room created! Waiting for another player...', 'waiting')
  } else if(m.type==='room_joined'){
    currentRoomId = m.roomId
    waitingInRoom.style.display = 'block'
    roomInfoDiv.innerHTML = `
      <div style="background: #e8f5e9; padding: 15px; border-radius: 8px; border: 2px solid #4caf50;">
        <div style="font-size: 1.2em; font-weight: bold; margin-bottom: 10px;">Room: ${m.room.name}</div>
        <div style="color: #666;">Room ID: ${m.roomId}</div>
      </div>
    `
    showStatus('✅ Joined room! Waiting for game to start...', 'waiting')
  } else if(m.type==='waiting'){
    gameStatus = 'waiting'
    showStatus('⏳ Waiting for opponent... (timeout: ' + m.timeout + 's)', 'waiting')
  } else if(m.type==='start'){
    // Hide all room UI
    waitingInRoom.style.display = 'none'

    gameId = m.gameId
    myPlayer = m.you
    opponent = m.opponent
    gameState = m.state
    gameStatus = 'playing'

    // Show game info
    gameInfo.style.display = 'flex'
    if(myPlayer === 1) {
      player1Name.textContent = currentUsername + ' (You)'
      player2Name.textContent = opponent
    } else {
      player1Name.textContent = opponent
      player2Name.textContent = currentUsername + ' (You)'
    }

    showStatus('🎮 Game started! Playing against ' + opponent, 'playing')
    winnerAnnouncement.innerHTML = ''
    render()
    fetchLeaderboard()
  } else if(m.type==='state'){
    gameState = m.state
    render()

    if(m.status==='finished'){
      gameStatus = 'finished'
      console.log('Game finished! Full message:', m)
      console.log('Result field:', m.result)
      console.log('Result is truthy:', !!m.result)
      console.log('Result === "draw":', m.result === 'draw')

      // Ensure we have a result before calling handleGameFinished
      if(m.result !== undefined && m.result !== null) {
        handleGameFinished(m.result)
      } else {
        console.warn('Result is undefined or null, using empty string')
        handleGameFinished('')
      }
      fetchLeaderboard()
    } else {
      updateTurnIndicator()
    }
  } else if(m.type==='reconnected'){
    gameState = m.state
    gameStatus = 'playing'
    showStatus('Reconnected to game', 'playing')
    render()
  } else if(m.error){
    showStatus('❌ Error: ' + m.error, 'error')
    setTimeout(() => {
      resetToStart()
    }, 3000)
  }
}

function handleGameFinished(result) {
  // Use result from server
  const winner = result ? String(result).trim() : ''
  console.log('handleGameFinished called with result:', result)
  console.log('Winner after processing:', winner)
  console.log('Result type:', typeof winner)
  console.log('Result length:', winner.length)
  console.log('Result === "draw":', winner === 'draw')
  console.log('Result truthy:', !!winner)

  if(winner === 'draw') {
    console.log('Showing draw message')
    showStatus('🤝 Game ended in a draw!', 'finished')
    winnerAnnouncement.innerHTML = '<div style="font-size: 3em; margin: 20px 0;">🤝</div><div>Draw!</div><div style="font-size: 0.8em; margin-top: 10px; color: #666;">Redirecting in 10 seconds...</div>'
    winnerAnnouncement.className = 'winner-announcement draw'
  } else if(winner && winner.length > 0 && winner !== 'draw') {
    console.log('Showing winner message for:', winner)
    // Get the current player's username
    const myUsername = usernameInput.value.trim()
    console.log('My username:', myUsername)
    const iWon = winner === myUsername
    console.log('Did I win?', iWon)

    if(iWon) {
      console.log('Showing YOU WON message')
      showStatus('🎉 You won!', 'finished')
      winnerAnnouncement.innerHTML = '<div style="font-size: 3em; margin: 20px 0;">🎉</div><div>You Won!</div><div style="font-size: 0.9em; margin-top: 10px; color: #155724;">Congratulations!</div><div style="font-size: 0.8em; margin-top: 10px; color: #666;">Redirecting in 10 seconds...</div>'
      winnerAnnouncement.className = 'winner-announcement win'
    } else {
      console.log('Showing YOU LOST message for opponent:', winner)
      showStatus('😔 ' + winner + ' won! You lost!', 'finished')
      winnerAnnouncement.innerHTML = '<div style="font-size: 3em; margin: 20px 0;">😔</div><div>' + winner + ' Won!</div><div style="font-size: 0.9em; margin-top: 10px; color: #721c24;">You Lost!</div><div style="font-size: 0.8em; margin-top: 10px; color: #666;">Redirecting in 10 seconds...</div>'
      winnerAnnouncement.className = 'winner-announcement lose'
    }
  } else {
    console.warn('No valid result received from server - result is:', result)
    showStatus('Game finished', 'finished')
    winnerAnnouncement.innerHTML = '<div>Game Finished</div><div style="font-size: 0.8em; margin-top: 10px; color: #666;">Redirecting in 10 seconds...</div>'
  }

  // Show leaderboard
  fetchLeaderboard()

  // Redirect to home after 10 seconds
  setTimeout(() => {
    console.log('Redirecting to home page')
    resetGame()
  }, 10000)
}

function resetGame() {
  // Reset UI
  gameStatus = 'idle'
  gameState = null
  myPlayer = null
  gameId = null
  opponent = null

  // Hide game info
  gameDiv.innerHTML = ''
  gameInfo.style.display = 'none'
  winnerAnnouncement.innerHTML = ''

  // Reset all sections
  usernameSection.style.display = 'block'
  modeSelection.style.display = 'none'
  createRoomSection.style.display = 'none'
  roomListSection.style.display = 'none'
  waitingInRoom.style.display = 'none'

  // Reset input
  usernameInput.value = ''
  currentUsername = ''
  currentRoomId = null

  // Show join section
  showStatus('Ready to play! Enter your username and join.', 'idle')
}

function updateTurnIndicator() {
  if(!gameState) return

  const currentTurn = gameState.turn
  const isMyTurn = currentTurn === myPlayer

  // Update player info highlighting
  if(currentTurn === 1) {
    player1Info.classList.add('active')
    player2Info.classList.remove('active')
  } else {
    player2Info.classList.add('active')
    player1Info.classList.remove('active')
  }

  if(isMyTurn) {
    showStatus('🎯 Your turn! (Player ' + myPlayer + ')', 'playing')
  } else {
    showStatus('⏳ Opponent\'s turn... (Player ' + currentTurn + ')', 'playing')
  }
}

function showStatus(message, type) {
  status.textContent = message
  status.className = type
}

function render(){
  gameDiv.innerHTML=''
  if(!gameState) return

  const rows = gameState.rows
  const cols = gameState.cols
  const b = gameState.board
  const grid = document.createElement('div')
  grid.id='board'

  for(let r=0;r<rows;r++){
    for(let c=0;c<cols;c++){
      const cell = document.createElement('div')
      cell.className='cell'
      const val = b[r][c]

      if(val===1 || val===2){
        const disc = document.createElement('div')
        disc.className = 'disc p' + val
        cell.appendChild(disc)
      }

      // Only allow clicks if it's the player's turn and game is playing
      if(gameStatus === 'playing' && gameState.turn === myPlayer) {
        cell.onclick = ()=>{
          if(ws && gameId) {
            ws.send(JSON.stringify({type:'move', gameId, col:c}))
          }
        }
      } else {
        cell.classList.add('disabled')
      }

      grid.appendChild(cell)
    }
  }
  gameDiv.appendChild(grid)

  if(gameStatus === 'playing') {
    updateTurnIndicator()
  }
}

function loadRooms(){
  fetch('/rooms').then(r=>r.json()).then(rooms=>{
    if(!rooms || rooms.length === 0) {
      roomListDiv.innerHTML = '<div class="empty-rooms">No rooms available. Create one!</div>'
      return
    }

    let html = ''
    rooms.forEach(room => {
      html += `
        <div class="room-item">
          <div class="room-item-info">
            <div class="room-item-name">${room.name}</div>
            <div class="room-item-details">
              Created by: ${room.creator} | Players: ${room.players}/${room.max_players}
            </div>
          </div>
          <button onclick="joinRoom('${room.id}')">Join</button>
        </div>
      `
    })
    roomListDiv.innerHTML = html
  }).catch(() => {
    roomListDiv.innerHTML = '<div class="empty-rooms">Failed to load rooms</div>'
  })
}

// Make joinRoom available globally
window.joinRoom = function(roomId) {
  roomListSection.style.display = 'none'
  connectJoinRoom(currentUsername, roomId)
}

function fetchLeaderboard(){
  fetch('/leaderboard').then(r=>r.json()).then(data=>{
    if(!data || Object.keys(data).length === 0) {
      lb.innerHTML = '<p style="text-align:center;color:#999;">No games played yet</p>'
      return
    }

    // Convert to array and sort
    const entries = Object.entries(data).sort((a, b) => b[1] - a[1])

    let html = '<table class="leaderboard-table">'
    html += '<tr><th>Rank</th><th>Player</th><th>Wins</th></tr>'

    entries.forEach((entry, idx) => {
      const [username, wins] = entry
      html += `<tr>
        <td class="rank">#${idx + 1}</td>
        <td>${username}</td>
        <td><strong>${wins}</strong></td>
      </tr>`
    })

    html += '</table>'
    lb.innerHTML = html
  }).catch(() => {
    lb.innerHTML = '<p style="color:red;">Failed to load leaderboard</p>'
  })
}

function id(s){ return document.getElementById(s) }

// Initial load
fetchLeaderboard()
setInterval(fetchLeaderboard, 10000) // Refresh every 10 seconds
})()
