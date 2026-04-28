import React, { useState } from 'react';
import { Plus, User, Users, LogOut } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/client';
import { useChat } from '../context/ChatContext';
import { Group, Message, User as UserModel } from '../types/api';

export function Sidebar() {
  const { username, userId, isLoggedIn, activeChat, setActiveChat, messages, loadMessages, logout } = useChat();

  const loadChatHistory = async (targetUser: UserModel) => {
    setActiveChat({ id: targetUser.id, name: targetUser.name, type: 'user' });
    try {
      const res = await apiClient.fetchMessagesBetween(targetUser.id);
      if (res.data) {
        const loadedMsgs = res.data.map((m: any) => ({
          id: m.id,
          sender: m.from_id === userId ? username : targetUser.name,
          content: m.content,
          timestamp: new Date(m.created_at),
          isMine: m.from_id === userId,
          type: 'user',
        })) as Message[];
        loadMessages(targetUser.id, loadedMsgs);
      }
    } catch (e) {
      console.error("Failed to load chat history", e);
    }
  };

  const loadGroupChatHistory = async (group: Group) => {
    setActiveChat({ id: group.id, name: group.name, type: 'group' });
    try {
      const res = await apiClient.fetchGroupMessages(group.id);
      if (res.data) {
        const loadedMsgs = res.data.map((m: any) => ({
          id: m.id,
          sender: m.from_id === userId ? username : m.from_id, // We'll show the raw sender UUID for now unless we have a user lookup
          content: m.content,
          timestamp: new Date(m.created_at),
          isMine: m.from_id === userId,
          type: 'group',
        })) as Message[];
        loadMessages(group.id, loadedMsgs);
      }
    } catch (e) {
      console.error("Failed to load group chat history", e);
    }
  };
  const [activeTab, setActiveTab] = useState<'users' | 'groups'>('users');
  const [directMessageUserId, setDirectMessageUserId] = useState('');
  const [groupJoinId, setGroupJoinId] = useState('');

  const { data: usersData, refetch: refetchUsers } = useQuery({
    queryKey: ['users'],
    queryFn: () => apiClient.fetchUsers(),
    enabled: isLoggedIn,
  });

  const { data: contactsData } = useQuery({
    queryKey: ['contacts'],
    queryFn: () => apiClient.fetchContacts(),
    enabled: isLoggedIn,
  });

  const { data: myGroupsData, refetch: refetchMyGroups } = useQuery({
    queryKey: ['myGroups'],
    queryFn: async () => {
      const res = await apiClient.fetchMyGroups();
      return res;
    },
    enabled: isLoggedIn,
  });

  const contactNames = contactsData?.data?.map(c => c.name) || [];
  const activeChatIds = Array.from(new Set([...contactNames]));

  const startDirectMessage = async () => {
    const name = directMessageUserId.trim();
    if (!name) return;

    const users = usersData?.data || [];
    const targetUser = users.find(u => u.name === name);
    if (targetUser) {
      loadChatHistory(targetUser);
    } else {
      const res = await refetchUsers();
      const updatedUsers = res.data?.data || [];
      const user = updatedUsers.find(u => u.name === name);
      if (user) {
        loadChatHistory(user);
      } else {
        alert("User not found in backend!");
      }
    }

    setDirectMessageUserId('');
  };

  return (
    <div className="glass sidebar animate-fade-in">
      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', paddingBottom: '16px', borderBottom: '1px solid var(--panel-border)' }}>
        <div className="avatar">{username.charAt(0).toUpperCase()}</div>
        <div>
          <div style={{ fontWeight: 600 }}>{username}</div>
          <div style={{ fontSize: '0.8rem', color: 'var(--success)', display: 'flex', alignItems: 'center', gap: '4px' }}>
            <span className="status-dot"></span> Online
          </div>
        </div>
        <button
          onClick={logout}
          title="Çıkış Yap"
          style={{ marginLeft: 'auto', background: 'transparent', border: 'none', color: 'var(--danger)', cursor: 'pointer', padding: '8px' }}>
          <LogOut size={20} />
        </button>
      </div>

      <div className="tabs">
        <div className={`tab ${activeTab === 'users' ? 'active' : ''}`} onClick={() => setActiveTab('users')}>
          Direct Messages
        </div>
        <div className={`tab ${activeTab === 'groups' ? 'active' : ''}`} onClick={() => setActiveTab('groups')}>
          Groups
        </div>
      </div>

      {activeTab === 'users' && (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', flex: 1 }}>
          <div style={{ display: 'flex', gap: '8px' }}>
            <input
              className="input"
              placeholder="Username..."
              value={directMessageUserId}
              onChange={e => setDirectMessageUserId(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && startDirectMessage()}
            />
            <button className="btn" onClick={startDirectMessage}><Plus size={18} /></button>
          </div>

          <div style={{ overflowY: 'auto', flex: 1, display: 'flex', flexDirection: 'column', gap: '4px' }}>
            {activeChatIds.map(id => (
              <div
                key={id}
                className={`list-item ${activeChat?.id === id ? 'active' : ''}`}
                onClick={() => {
                  const targetUser = usersData?.data?.find(u => u.name === id) || contactsData?.data?.find(u => u.name === id);
                  if (targetUser) {
                    loadChatHistory(targetUser);
                  }
                }}
              >
                <div className="avatar"><User size={20} /></div>
                <div style={{ flex: 1, overflow: 'hidden', textOverflow: 'ellipsis' }}>{id}</div>
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === 'groups' && (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', flex: 1 }}>
          <div style={{ overflowY: 'auto', flex: 1, display: 'flex', flexDirection: 'column', gap: '4px' }}>
            {myGroupsData?.data?.map((g: Group) => (
              <div
                key={g.id}
                className={`list-item ${activeChat?.id === g.id ? 'active' : ''}`}
                onClick={() => loadGroupChatHistory(g)}
              >
                <div className="avatar group-avatar"><Users size={20} color="white" /></div>
                <div style={{ flex: 1, overflow: 'hidden', textOverflow: 'ellipsis' }}>{g.name}</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
