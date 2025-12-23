import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import "../App.css";

function StartPage() {
  const [username, setUsername] = useState("");
  const [waiting, setWaiting] = useState(false);
  const [timeLeft, setTimeLeft] = useState(10);
  const navigate = useNavigate();

  const startGame = () => {
    if (!username.trim()) {
      alert("Please enter your username");
      return;
    }
    setWaiting(true);
  };

  // ⏱️ 10 second countdown
  useEffect(() => {
    if (!waiting) return;

    const timer = setInterval(() => {
      setTimeLeft((t) => t - 1);
    }, 1000);

    const redirect = setTimeout(() => {
      navigate("/game", { state: { username } });
    }, 10000);

    return () => {
      clearInterval(timer);
      clearTimeout(redirect);
    };
  }, [waiting, navigate, username]);

  return (
    <div className="start-page">
      <h1 className="title">🎮 4 in a Row</h1>

      {!waiting && (
        <div className="start-box">
          <input
            className="username-input"
            placeholder="Enter your username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <button className="start-btn" onClick={startGame}>
            Start Game
          </button>
        </div>
      )}

      {waiting && (
        <div className="waiting-box">
          <h2>Waiting for another player...</h2>
          <p>Bot will join if no one connects</p>
          <h3>⏳ {timeLeft}s</h3>
        </div>
      )}
    </div>
  );
}

export default StartPage;
