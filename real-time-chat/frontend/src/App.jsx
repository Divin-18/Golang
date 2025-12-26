import { useState, useEffect } from "react";
import { AuthProvider, useAuth } from "./contexts/AuthContext";
import { WebSocketProvider } from "./contexts/WebSocketContext";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Chat from "./pages/Chat";
import "./App.css";

function AppContent() {
  const { user, loading } = useAuth();
  const [currentPage, setCurrentPage] = useState("login");

  useEffect(() => {
    if (user) {
      setCurrentPage("chat");
    }
  }, [user]);

  if (loading) {
    return (
      <div className="loading-screen">
        <div className="loading-content">
          <div className="spinner large"></div>
          <p>Loading...</p>
        </div>
      </div>
    );
  }

  if (user) {
    return (
      <WebSocketProvider>
        <Chat />
      </WebSocketProvider>
    );
  }

  return (
    <div className="auth-container">
      {currentPage === "login" ? (
        <Login onSwitch={() => setCurrentPage("register")} />
      ) : (
        <Register onSwitch={() => setCurrentPage("login")} />
      )}
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;
