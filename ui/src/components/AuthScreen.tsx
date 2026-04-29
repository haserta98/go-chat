import { useState } from 'react';
import { Send, MessageSquare } from 'lucide-react';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '../api/client';
import { useAuthStore } from '../stores/authStore';
import { Button } from './ui/Button';
import { Input } from './ui/Input';
import { Avatar } from './ui/Avatar';

export function AuthScreen() {
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const login = useAuthStore((s) => s.login);

  const loginMut = useMutation({
    mutationFn: () => apiClient.login(name, password),
    onSuccess: (data) => {
      if (data.error) return alert(data.error);
      if (data.data) login(data.data.name, data.data.id);
    },
  });

  const registerMut = useMutation({
    mutationFn: () => apiClient.register(name, password),
    onSuccess: (data) => {
      if (data.error) return alert(data.error);
      loginMut.mutate();
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !password) return;
    mode === 'register' ? registerMut.mutate() : loginMut.mutate();
  };

  const isPending = loginMut.isPending || registerMut.isPending;

  return (
    <div className="glass w-full max-w-md p-10 text-center flex flex-col gap-5 animate-in">
      <div className="flex justify-center mb-2">
        <Avatar name="N" size="lg">
          <MessageSquare size={28} color="white" />
        </Avatar>
      </div>

      <h1 className="text-2xl font-semibold text-white">
        {mode === 'login' ? 'Welcome to Nexus Chat' : 'Create an Account'}
      </h1>
      <p className="text-slate-400 text-sm">
        {mode === 'login' ? 'Enter your credentials to connect' : 'Sign up to get started'}
      </p>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <Input placeholder="Username" value={name} onChange={(e) => setName(e.target.value)} autoFocus />
        <Input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <Button loading={isPending} icon={<Send size={16} />} className="py-3.5 text-base">
          {mode === 'login' ? 'Login' : 'Register'}
        </Button>
      </form>

      <p className="text-sm text-slate-400">
        {mode === 'login' ? "Don't have an account? " : 'Already have an account? '}
        <span className="text-primary cursor-pointer hover:underline" onClick={() => setMode(mode === 'login' ? 'register' : 'login')}>
          {mode === 'login' ? 'Register' : 'Login'}
        </span>
      </p>
    </div>
  );
}
