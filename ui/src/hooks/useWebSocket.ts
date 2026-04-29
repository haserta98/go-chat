import { useEffect } from 'react';
import { useAuthStore } from '../stores/authStore';
import { useChatStore } from '../stores/chatStore';
import { apiClient } from '../api/client';
import { Message } from '../types/api';

/** Manages WebSocket lifecycle: connect on login, handle messages, disconnect on unmount */
export function useWebSocket() {
  const isLoggedIn = useAuthStore((s) => s.isLoggedIn);
  const { setWs, addMessage, setUserOnline, setUserOffline, setOnlineUsers } = useChatStore.getState();

  // Fetch initial online statuses
  useEffect(() => {
    if (!isLoggedIn) return;
    apiClient.fetchContactsOnlineStatus().then((res) => {
      if (res.status === 'success' && res.data) {
        const set = new Set<string>();
        for (const [id, online] of Object.entries(res.data)) {
          if (online) set.add(id);
        }
        setOnlineUsers(set);
      }
    });
  }, [isLoggedIn]);

  // WebSocket connection
  useEffect(() => {
    if (!isLoggedIn) return;

    const socket = new WebSocket('ws://localhost:9090/ws');

    socket.onopen = () => setWs(socket);

    socket.onmessage = (event) => {
      try {
        const { type, payload } = JSON.parse(event.data);

        if (type === 'presence_change') {
          const data = typeof payload === 'string' ? JSON.parse(payload) : payload;
          data.status === 'online'
            ? useChatStore.getState().setUserOnline(data.user_id)
            : useChatStore.getState().setUserOffline(data.user_id);
          return;
        }

        let chatId = '';
        let sender = 'Unknown';
        if (type === 'send_user_message') {
          chatId = payload.from;
          sender = payload.to;
        } else if (type === 'send_group_message') {
          chatId = payload.to;
          sender = payload.from;
        }

        if (chatId) {
          useChatStore.getState().addMessage(chatId, {
            id: Date.now().toString() + Math.random(),
            sender,
            content: payload.content,
            timestamp: new Date(),
            isMine: false,
            type: payload.type,
          } as Message);
        }
      } catch (e) {
        console.error('WS parse error:', e);
      }
    };

    socket.onclose = () => setWs(null);

    return () => { socket.close(); };
  }, [isLoggedIn]);
}
