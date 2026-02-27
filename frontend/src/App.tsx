import { Routes, Route, Navigate } from 'react-router-dom';
import { useState } from 'react';
import Landing from './pages/Landing';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Dashboard from './pages/Dashboard';
import ForgotPassword from './pages/ForgotPassword';
import ResetPassword from './pages/ResetPassword';
import { Spinner } from './components/Spinner';
import { useAuth } from './hooks/useAuth';
import type { WebUser } from './types';

export default function App() {
  const { user, setUser, loading, signOut } = useAuth();
  const [localUser, setLocalUser] = useState<WebUser | null>(null);

  const effectiveUser = user ?? localUser;

  const handleLogin = (u: WebUser) => {
    setUser(u);
    setLocalUser(u);
  };

  const handleSignOut = () => {
    signOut();
    setLocalUser(null);
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  return (
    <Routes>
      <Route path="/" element={<Landing />} />
      <Route
        path="/login"
        element={effectiveUser ? <Navigate to="/dashboard" replace /> : <Login onLogin={handleLogin} />}
      />
      <Route
        path="/signup"
        element={effectiveUser ? <Navigate to="/dashboard" replace /> : <Signup onLogin={handleLogin} />}
      />
      <Route
        path="/dashboard"
        element={
          effectiveUser ? (
            <Dashboard user={effectiveUser} onSignOut={handleSignOut} />
          ) : (
            <Navigate to="/login" replace />
          )
        }
      />
      <Route path="/forgot-password" element={<ForgotPassword />} />
      <Route path="/reset-password" element={<ResetPassword />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
