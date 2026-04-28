import React from 'react';
import { useChat } from './context/ChatContext';
import { AuthScreen } from './components/AuthScreen';
import { Sidebar } from './components/Sidebar';
import { ChatArea } from './components/ChatArea';

export default function App() {
  const { isLoggedIn } = useChat();

  if (!isLoggedIn) {
    return <AuthScreen />;
  }

  return (
    <div className="app-container">
      <Sidebar />
      <ChatArea />
    </div>
  );
}
