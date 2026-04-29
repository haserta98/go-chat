import { useEffect } from 'react';
import { useAuthStore } from './stores/authStore';
import { useWebSocket } from './hooks/useWebSocket';
import { AuthScreen } from './components/AuthScreen';
import { Sidebar } from './components/Sidebar';
import { ChatArea } from './components/ChatArea';

export default function App() {
  const { isLoggedIn, hydrate } = useAuthStore();

  useEffect(() => { hydrate(); }, []);
  useWebSocket();

  if (!isLoggedIn) {
    return (
      <div className="w-full h-screen flex items-center justify-center">
        <AuthScreen />
      </div>
    );
  }

  return (
    <div className="w-full max-w-7xl h-[90vh] flex gap-5 p-5">
      <Sidebar />
      <ChatArea />
    </div>
  );
}
