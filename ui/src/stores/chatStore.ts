import { create } from 'zustand';
import { Message } from '../types/api';

interface ActiveChat {
  id: string;
  name: string;
  type: 'user' | 'group';
}

interface ChatState {
  ws: WebSocket | null;
  activeChat: ActiveChat | null;
  messages: Record<string, Message[]>;
  onlineUsers: Set<string>;

  setWs: (ws: WebSocket | null) => void;
  setActiveChat: (chat: ActiveChat | null) => void;
  addMessage: (chatId: string, msg: Message) => void;
  loadMessages: (chatId: string, msgs: Message[]) => void;
  setUserOnline: (userId: string) => void;
  setUserOffline: (userId: string) => void;
  setOnlineUsers: (users: Set<string>) => void;
  sendMessage: (content: string, userId: string, username: string) => void;
  reset: () => void;
}

export const useChatStore = create<ChatState>((set, get) => ({
  ws: null,
  activeChat: null,
  messages: {},
  onlineUsers: new Set(),

  setWs: (ws) => set({ ws }),
  setActiveChat: (chat) => set({ activeChat: chat }),

  addMessage: (chatId, msg) =>
    set((s) => ({
      messages: { ...s.messages, [chatId]: [...(s.messages[chatId] || []), msg] },
    })),

  loadMessages: (chatId, msgs) =>
    set((s) => ({ messages: { ...s.messages, [chatId]: msgs } })),

  setUserOnline: (userId) =>
    set((s) => {
      const next = new Set(s.onlineUsers);
      next.add(userId);
      return { onlineUsers: next };
    }),

  setUserOffline: (userId) =>
    set((s) => {
      const next = new Set(s.onlineUsers);
      next.delete(userId);
      return { onlineUsers: next };
    }),

  setOnlineUsers: (users) => set({ onlineUsers: users }),

  sendMessage: (content, userId, username) => {
    const { ws, activeChat, addMessage } = get();
    if (!content.trim() || !activeChat || !ws) return;

    const trimmed = content.trim();
    ws.send(JSON.stringify({
      type: activeChat.type === 'user' ? 'send_message' : 'send_group_message',
      payload: { from: userId, content: trimmed, to: activeChat.id, type: activeChat.type },
    }));

    addMessage(activeChat.id, {
      id: Date.now().toString(),
      sender: username,
      content: trimmed,
      timestamp: new Date(),
      isMine: true,
      type: activeChat.type,
    });
  },

  reset: () => set({ ws: null, activeChat: null, messages: {}, onlineUsers: new Set() }),
}));
