import React, { useState } from 'react';
import { Send, MessageSquare } from 'lucide-react';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '../api/client';
import { useChat } from '../context/ChatContext';

export function AuthScreen() {
  const [authMode, setAuthMode] = useState<'login' | 'register'>('login');
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const { login } = useChat();

  const registerMutation = useMutation({
    mutationFn: () => apiClient.register(name, password),
    onSuccess: (data) => {
      if (data.error) {
        alert(data.error);
        return;
      }
      alert("Registration successful! Logging in...");
      setAuthMode('login');
      // optionally trigger login right away, or let them click login
      loginMutation.mutate();
    },
    onError: (err) => {
      console.error(err);
      alert("Registration failed");
    }
  });

  const loginMutation = useMutation({
    mutationFn: () => apiClient.login(name, password),
    onSuccess: (data) => {
      if (data.error) {
        alert(data.error);
        return;
      }
      if (data.data && data.token) {
        login(data.data.name, data.data.id, data.token);
      }
    },
    onError: (err) => {
      console.error(err);
      alert("Login failed");
    }
  });

  const handleAuth = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !password) return;

    if (authMode === 'register') {
      registerMutation.mutate();
    } else {
      loginMutation.mutate();
    }
  };

  return (
    <div className="glass login-container animate-fade-in">
      <div style={{ display: 'flex', justifyContent: 'center', marginBottom: '20px' }}>
        <div className="avatar" style={{ width: 80, height: 80, fontSize: '2rem' }}>
          <MessageSquare size={40} color="white" />
        </div>
      </div>
      <h1 className="login-title">{authMode === 'login' ? 'Welcome to Nexus Chat' : 'Create an Account'}</h1>
      <p style={{ color: 'var(--text-secondary)' }}>
        {authMode === 'login' ? 'Enter your credentials to connect' : 'Sign up to get started'}
      </p>
      <form onSubmit={handleAuth} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <input
          type="text"
          className="input"
          placeholder="Username"
          value={name}
          onChange={(e) => setName(e.target.value)}
          autoFocus
        />
        <input
          type="password"
          className="input"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <button type="submit" className="btn" style={{ padding: '14px', fontSize: '1.1rem' }} disabled={loginMutation.isPending || registerMutation.isPending}>
          {loginMutation.isPending || registerMutation.isPending ? 'Loading...' : (authMode === 'login' ? 'Login' : 'Register')} <Send size={18} />
        </button>
      </form>
      <div style={{ marginTop: '16px', fontSize: '0.9rem', color: 'var(--text-secondary)' }}>
        {authMode === 'login' ? (
          <>Don't have an account? <span style={{ color: 'var(--primary)', cursor: 'pointer' }} onClick={() => setAuthMode('register')}>Register</span></>
        ) : (
          <>Already have an account? <span style={{ color: 'var(--primary)', cursor: 'pointer' }} onClick={() => setAuthMode('login')}>Login</span></>
        )}
      </div>
    </div>
  );
}
