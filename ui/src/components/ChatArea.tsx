import { useState, useEffect, useRef } from 'react';
import { Send, User, Users, MessageSquare } from 'lucide-react';
import { useAuthStore } from '../stores/authStore';
import { useChatStore } from '../stores/chatStore';
import { Avatar } from './ui/Avatar';
import { Input } from './ui/Input';
import { Button } from './ui/Button';

export function ChatArea() {
  const { userId, username } = useAuthStore();
  const { activeChat, messages, onlineUsers, sendMessage } = useChatStore();
  const [text, setText] = useState('');
  const endRef = useRef<HTMLDivElement>(null);

  const chatMessages = activeChat ? messages[activeChat.id] || [] : [];
  const isOnline = activeChat?.type === 'user' && onlineUsers.has(activeChat.id);

  useEffect(() => {
    endRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [chatMessages, activeChat]);

  const handleSend = (e: React.FormEvent) => {
    e.preventDefault();
    if (!text.trim()) return;
    sendMessage(text, userId, username);
    setText('');
  };

  if (!activeChat) {
    return (
      <div className="glass flex-1 flex flex-col items-center justify-center text-slate-400 gap-4 animate-in">
        <MessageSquare size={56} className="opacity-20" />
        <h2 className="text-lg font-semibold text-white">No Chat Selected</h2>
        <p className="text-sm">Select a user or group to start messaging.</p>
      </div>
    );
  }

  return (
    <div className="glass flex-1 flex flex-col overflow-hidden animate-in">
      {/* Header */}
      <div className="px-5 py-4 border-b border-border flex items-center gap-3">
        <Avatar
          name={activeChat.name}
          variant={activeChat.type === 'group' ? 'group' : 'primary'}
          online={activeChat.type === 'user' ? isOnline : undefined}
        >
          {activeChat.type === 'group' ? <Users size={18} color="white" /> : <User size={18} color="white" />}
        </Avatar>
        <div>
          <div className="font-semibold text-white">{activeChat.name}</div>
          <div className={`text-xs ${isOnline ? 'text-success' : 'text-slate-400'}`}>
            {activeChat.type === 'group' ? 'Group Chat' : isOnline ? 'Çevrimiçi' : 'Çevrimdışı'}
          </div>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-5 flex flex-col gap-3">
        {chatMessages.map((msg) => (
          <div key={msg.id}
            className={`max-w-[75%] px-4 py-3 rounded-xl leading-relaxed animate-in text-sm
              ${msg.isMine
                ? 'self-end bg-msg-sent rounded-br-sm'
                : 'self-start bg-msg-received rounded-bl-sm'}`}>
            {!msg.isMine && activeChat.type === 'group' && (
              <div className="text-xs font-semibold text-primary mb-1">{msg.sender}</div>
            )}
            <div>{msg.content}</div>
            <div className="text-[10px] opacity-60 mt-1">
              {msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
            </div>
          </div>
        ))}
        <div ref={endRef} />
      </div>

      {/* Input */}
      <form onSubmit={handleSend} className="px-5 py-4 border-t border-border flex gap-2.5">
        <Input placeholder={`${activeChat.name} ile mesajlaş...`} value={text} onChange={(e) => setText(e.target.value)} />
        <Button disabled={!text.trim()} icon={<Send size={16} />} />
      </form>
    </div>
  );
}
