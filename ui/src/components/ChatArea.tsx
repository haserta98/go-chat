import React, { useState, useEffect, useRef } from 'react';
import { Send, User, Users, MessageSquare } from 'lucide-react';
import { useChat } from '../context/ChatContext';

export function ChatArea() {
  const { activeChat, messages, sendMessage } = useChat();
  const [inputText, setInputText] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const activeMessages = activeChat ? messages[activeChat.id] || [] : [];

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, activeChat]);

  const handleSend = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputText.trim()) return;
    sendMessage(inputText);
    setInputText('');
  };

  if (!activeChat) {
    return (
      <div className="glass main-chat animate-fade-in" style={{ animationDelay: '0.1s' }}>
        <div className="empty-state">
          <MessageSquare size={64} opacity={0.2} />
          <h2>No Chat Selected</h2>
          <p>Select a user or join a group to start messaging.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="glass main-chat animate-fade-in" style={{ animationDelay: '0.1s' }}>
      <div className="chat-header">
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div className={`avatar ${activeChat.type === 'group' ? 'group-avatar' : ''}`}>
            {activeChat.type === 'group' ? <Users size={20} color="white" /> : <User size={20} color="white" />}
          </div>
          <div>
            <div style={{ fontWeight: 600, fontSize: '1.1rem' }}>{activeChat.id}</div>
            <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
              {activeChat.type === 'group' ? 'Group Chat' : 'Direct Message'}
            </div>
          </div>
        </div>
      </div>

      <div className="chat-messages">
        {activeMessages.map((msg) => (
          <div key={msg.id} className={`message ${msg.isMine ? 'sent' : 'received'}`}>
            {!msg.isMine && activeChat.type === 'group' && (
              <div style={{ fontSize: '0.8rem', fontWeight: 600, color: 'var(--primary)', marginBottom: '4px' }}>
                {msg.sender}
              </div>
            )}
            <div>{msg.content}</div>
            <div className="message-meta">
              <span>{msg.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <form className="chat-input-area" onSubmit={handleSend}>
        <input
          type="text"
          className="input"
          placeholder={`Message ${activeChat.id}...`}
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
        />
        <button type="submit" className="btn" disabled={!inputText.trim()}>
          <Send size={18} />
        </button>
      </form>
    </div>
  );
}
