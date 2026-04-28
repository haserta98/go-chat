import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Message } from '../types/api';

interface ActiveChat {
  id: string;
  name: string
  type: 'user' | 'group';
}

interface ChatState {
  username: string;
  userId: string;
  isLoggedIn: boolean;
  ws: WebSocket | null;
  activeChat: ActiveChat | null;
  messages: Record<string, Message[]>;

  login: (username: string, userId: string) => void;
  logout: () => void;
  setActiveChat: (chat: ActiveChat | null) => void;
  addMessage: (chatId: string, message: Message) => void;
  loadMessages: (chatId: string, messages: Message[]) => void;
  sendMessage: (content: string) => void;
}

const ChatContext = createContext<ChatState | undefined>(undefined);

export function ChatProvider({ children }: { children: ReactNode }) {
  const [username, setUsername] = useState<string>('');
  const [userId, setUserId] = useState<string>('');
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);

  const [activeChat, setActiveChat] = useState<ActiveChat | null>(null);
  const [messages, setMessages] = useState<Record<string, Message[]>>({});

  useEffect(() => {
    import('../api/client').then(({ apiClient }) => {
      apiClient.fetchMe().then(res => {
        if (res.status === 'success' && res.data) {
          setUsername(res.data.name);
          setUserId(res.data.id);
          setIsLoggedIn(true);
        }
      }).catch(e => {
        console.error('Failed to fetch user profile', e);
      });
    });
  }, []);

  useEffect(() => {
    if (!isLoggedIn) return;

    const socket = new WebSocket(`ws://localhost:8080/ws`);

    socket.onopen = () => {
      setWs(socket);
    };

    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data);
        const data = payload.payload;
        const type = payload.type;

        let chatId = '';
        let sender = 'Unknown';

        if (type === 'send_user_message') {
          chatId = data.from
          sender = data.to;
        } else if (type === 'send_group_message') {
          chatId = data.to;
          sender = data.from;
        }

        if (chatId) {
          const newMessage: Message = {
            id: Date.now().toString() + Math.random().toString(),
            sender,
            content: data.content,
            timestamp: new Date(),
            isMine: false,
            type: data.type,
          };
          addMessage(chatId, newMessage);
        }
      } catch (e) {
        console.error('Error parsing socket message:', e);
      }
    };

    socket.onclose = () => {
      setWs(null);
    };

    return () => {
      socket.close();
    };
  }, [isLoggedIn, userId, username]);

  const login = (newUsername: string, newUserId: string) => {
    setUsername(newUsername);
    setUserId(newUserId);
    setIsLoggedIn(true);
  };

  const logout = async () => {
    try {
      const { apiClient } = await import('../api/client');
      await apiClient.logout();
    } catch (e) {
      console.error('Logout error', e);
    }
    if (ws) ws.close();
    setUsername('');
    setUserId('');
    setIsLoggedIn(false);
    setActiveChat(null);
  };

  const addMessage = (chatId: string, message: Message) => {
    setMessages(prev => ({
      ...prev,
      [chatId]: [...(prev[chatId] || []), message]
    }));
  };

  const loadMessages = (chatId: string, history: Message[]) => {
    setMessages(prev => ({
      ...prev,
      [chatId]: history
    }));
  };


  const sendMessage = (content: string) => {
    if (!content.trim() || !activeChat || !ws) return;

    const payload = {
      from: userId,
      content: content.trim(),
      to: activeChat.id,
      type: activeChat.type
    };

    const request = {
      type: activeChat.type === 'user' ? 'send_message' : 'send_group_message',
      payload,
    };

    ws.send(JSON.stringify(request));

    const newMessage: Message = {
      id: Date.now().toString(),
      sender: username,
      content: content.trim(),
      timestamp: new Date(),
      isMine: true,
      type: activeChat.type,
    };

    addMessage(activeChat.id, newMessage);
  };

  return (
    <ChatContext.Provider value={{
      username, userId, isLoggedIn, ws, activeChat, messages,
      login, logout, setActiveChat, addMessage, loadMessages, sendMessage
    }}>
      {children}
    </ChatContext.Provider>
  );
}

export function useChat() {
  const context = useContext(ChatContext);
  if (!context) throw new Error("useChat must be used within ChatProvider");
  return context;
}
