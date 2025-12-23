import { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import "../App.css";

const ROWS = 6;
const COLS = 7;

function GamePage() {
  const location = useLocation();
  const navigate = useNavigate();
  const username = location.state?.username || "Player1";

  const [board, setBoard] = useState(
    Array.from({ length: ROWS }, () => Array(COLS).fill(""))
  );

  const [currentTurn, setCurrentTurn] = useState("X");
  const [winner, setWinner] = useState("");
  const [isDraw, setIsDraw] = useState(false);
  const [winCells, setWinCells] = useState([]);
  const [duration, setDuration] = useState(0);
  const [showPopup, setShowPopup] = useState(false);

  // ⏱️ TIMER
  const [liveTime, setLiveTime] = useState(0);
  const [timerStarted, setTimerStarted] = useState(false);
  const [isPaused, setIsPaused] = useState(false); // 🔥 NEW
  const hasGameStartedRef = useRef(false);

  // 🔒 UNDO LOCK
  const [undoLock, setUndoLock] = useState(false);

  // 📊 wins
  const [p1Wins, setP1Wins] = useState(
    Number(localStorage.getItem("p1Wins")) || 0
  );
  const [p2Wins, setP2Wins] = useState(
    Number(localStorage.getItem("p2Wins")) || 0
  );

  const wsRef = useRef(null);

  // 🔌 WEBSOCKET
  useEffect(() => {
    if (!username || wsRef.current) return;

    const socket = new WebSocket("wss://emitrrconnect4.onrender.com/ws");
    wsRef.current = socket;

    socket.onopen = () => {
      socket.send(JSON.stringify({ username }));
    };

    socket.onmessage = (e) => {
      const data = JSON.parse(e.data);

      if (Array.isArray(data.board)) {
        setBoard(data.board);

        // ▶️ START TIMER AFTER FIRST MOVE
        if (!hasGameStartedRef.current) {
          const hasAnyMove = data.board.some((row) =>
            row.some((cell) => cell === "X" || cell === "O")
          );

          if (hasAnyMove) {
            hasGameStartedRef.current = true;
            setTimerStarted(true);
          }
        }
      }

      if (data.player) setCurrentTurn(data.player);

      setUndoLock(false);

      // 🏆 WIN
      if (data.winner && !winner) {
        setWinner(data.winner);
        setWinCells(data.winCells || []);
        setDuration(Math.floor(data.duration || 0));
        setIsPaused(true); // ⏸️ STOP TIMER

        if (data.winner === username) {
          const v = p1Wins + 1;
          setP1Wins(v);
          localStorage.setItem("p1Wins", v);
        } else {
          const v = p2Wins + 1;
          setP2Wins(v);
          localStorage.setItem("p2Wins", v);
        }

        setTimeout(() => setShowPopup(true), 800);
      }

      // 🤝 DRAW
      if (data.isDraw && !winner && !isDraw) {
        setIsDraw(true);
        setWinner("DRAW");
        setDuration(Math.floor(data.duration || 0));
        setIsPaused(true);
        setTimeout(() => setShowPopup(true), 800);
      }
    };

    socket.onclose = () => (wsRef.current = null);
  }, [username, winner, isDraw, p1Wins, p2Wins]);

  // ⏱️ TIMER EFFECT (PAUSE / RESUME)
  useEffect(() => {
    if (!timerStarted || winner || isDraw || isPaused) return;

    const interval = setInterval(() => {
      setLiveTime((t) => t + 1);
    }, 1000);

    return () => clearInterval(interval);
  }, [timerStarted, winner, isDraw, isPaused]);

  // ▶️ PLAY MOVE (BLOCK WHEN PAUSED)
  const play = (col) => {
    if (
      !wsRef.current ||
      wsRef.current.readyState !== WebSocket.OPEN ||
      winner ||
      isDraw ||
      isPaused || // 🔥 PAUSE BLOCK
      currentTurn !== "X" ||
      undoLock
    )
      return;

    wsRef.current.send(JSON.stringify({ column: col }));
  };

  // ⬅️ UNDO
  const undoMove = () => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return;
    setUndoLock(true);
    wsRef.current.send(JSON.stringify({ undo: true }));
  };

  const togglePause = () => {
    if (!timerStarted || winner || isDraw) return;
    setIsPaused((p) => !p);
  };

  const newGame = () => window.location.reload();

  const isWinCell = (r, c) =>
    winCells.some(([row, col]) => row === r && col === c);

  return (
    <div className="game-page">
      <h1 className="game-title">🎮 4 in a Row</h1>

      {/* ⏱️ TIMER + PAUSE */}
      <div className="timer">
        ⏱️ Time: {liveTime}s
        <button
          className="undo-btn"
          style={{ marginLeft: "15px" }}
          onClick={togglePause}
        >
          {isPaused ? "▶ Resume" : "⏸ Pause"}
        </button>
      </div>

      <div className="game-layout">
        {/* PLAYER 1 */}
        <div className={`player-card ${currentTurn === "X" ? "active" : ""}`}>
          <img src="/avatar1.png" alt="P1" className="avatar" />
          <h3>Player 1</h3>
          <p>{username}</p>
          <span className="symbol x">X</span>
        </div>

        {/* GRID */}
        <div className={`grid-wrapper ${isPaused ? "disabled" : ""}`}>
          <div className="grid">
            {board.map((row, r) =>
              row.map((cell, c) => (
                <div
                  key={`${r}-${c}`}
                  className={`cell ${isWinCell(r, c) ? "win-cell" : ""}`}
                  onClick={() => play(c)}
                >
                  {cell && (
                    <span
                      className={`disc ${cell === "X" ? "x" : "o"} ${
                        isWinCell(r, c) ? "win-disc" : ""
                      }`}
                    >
                      {cell}
                    </span>
                  )}
                </div>
              ))
            )}
          </div>
        </div>

        {/* PLAYER 2 */}
        <div className={`player-card ${currentTurn === "O" ? "active" : ""}`}>
          <img src="/avatar2.png" alt="BOT" className="avatar" />
          <h3>Player 2</h3>
          <p>BOT</p>
          <span className="symbol o">O</span>
        </div>
      </div>

      {/* FOOTER */}
      <div className="footer">
        <button className="back-btn" onClick={() => navigate("/")}>
          ← Back
        </button>

        <button className="undo-btn" onClick={undoMove} disabled={undoLock || isPaused}>
          ⬅️ Undo
        </button>

        <button className="new-game-btn" onClick={newGame}>
          New Game
        </button>
      </div>

      {/* POPUP */}
      {showPopup && (
        <div className="popup-overlay">
          <div className="popup-card">
            <h2>
              {winner === "DRAW" ? "🤝 Match Draw!" : "🎉 Congratulations 🎉"}
            </h2>
            <p>
              <strong>Result:</strong>{" "}
              {winner === "DRAW" ? "No Winner" : winner}
            </p>
            <p>
              <strong>Game Time:</strong> {duration} sec
            </p>
            <p>🏅 Player 1 Wins: {p1Wins}</p>
            <p>🤖 Player 2 Wins: {p2Wins}</p>
            <button onClick={newGame}>Play Again</button>
          </div>
        </div>
      )}
    </div>
  );
}

export default GamePage;
