import { useState } from 'react';
import { Users, LogOut, UserPlus, UserMinus, X, Plus, DoorOpen, Search as SearchIcon } from 'lucide-react';
import { useAuthStore } from '../stores/authStore';
import { useChatStore } from '../stores/chatStore';
import {
  useContacts, useUsers, useMyGroups, useAllGroups,
  useAddContact, useRemoveContact,
  useCreateGroup, useJoinGroup, useLeaveGroup,
} from '../hooks/useQueries';
import { apiClient } from '../api/client';
import { Avatar } from './ui/Avatar';
import { Button, IconButton } from './ui/Button';
import { Input } from './ui/Input';
import { Group, Message, User } from '../types/api';

export function Sidebar() {
  const { username, userId, logout } = useAuthStore();
  const { activeChat, setActiveChat, loadMessages, onlineUsers } = useChatStore();
  const [tab, setTab] = useState<'users' | 'groups'>('users');

  // Contact state
  const [showAdd, setShowAdd] = useState(false);
  const [search, setSearch] = useState('');

  // Group state
  const [showGroupPanel, setShowGroupPanel] = useState(false);
  const [newGroupName, setNewGroupName] = useState('');
  const [groupSearch, setGroupSearch] = useState('');

  // Queries
  const { data: contacts = [] } = useContacts();
  const { data: allUsers = [] } = useUsers();
  const { data: myGroups = [] } = useMyGroups();
  const { data: allGroups = [] } = useAllGroups();

  // Mutations
  const addMut = useAddContact();
  const removeMut = useRemoveContact();
  const createGroupMut = useCreateGroup();
  const joinGroupMut = useJoinGroup();
  const leaveGroupMut = useLeaveGroup();

  // Contact filtering
  const contactIds = new Set(contacts.map((c: User) => c.id));
  const searchable = allUsers.filter((u: User) => u.id !== userId && !contactIds.has(u.id));
  const filtered = search.trim()
    ? searchable.filter((u: User) => u.name.toLowerCase().includes(search.toLowerCase()))
    : searchable;

  // Group filtering: groups I'm NOT a member of
  const myGroupIds = new Set(myGroups.map((g: Group) => g.id));
  const joinableGroups = allGroups.filter((g: Group) => !myGroupIds.has(g.id));
  const filteredJoinable = groupSearch.trim()
    ? joinableGroups.filter((g: Group) => g.name.toLowerCase().includes(groupSearch.toLowerCase()))
    : joinableGroups;

  const openChat = async (user: User) => {
    setActiveChat({ id: user.id, name: user.name, type: 'user' });
    try {
      const res = await apiClient.fetchMessagesBetween(user.id);
      if (res.data) {
        loadMessages(user.id, res.data.map((m: any) => ({
          id: m.id, sender: m.from_id === userId ? username : user.name,
          content: m.content, timestamp: new Date(m.created_at),
          isMine: m.from_id === userId, type: 'user',
        })) as Message[]);
      }
    } catch (e) { console.error(e); }
  };

  const openGroup = async (group: Group) => {
    setActiveChat({ id: group.id, name: group.name, type: 'group' });
    try {
      const res = await apiClient.fetchGroupMessages(group.id);
      if (res.data) {
        loadMessages(group.id, res.data.map((m: any) => ({
          id: m.id, sender: m.from_id === userId ? username : m.from_id,
          content: m.content, timestamp: new Date(m.created_at),
          isMine: m.from_id === userId, type: 'group',
        })) as Message[]);
      }
    } catch (e) { console.error(e); }
  };

  const handleRemoveContact = (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    if (!confirm('Bu kişiyi silmek istediğine emin misin?')) return;
    removeMut.mutate(id, {
      onSuccess: () => { if (activeChat?.id === id) setActiveChat(null); },
    });
  };

  const handleCreateGroup = () => {
    if (!newGroupName.trim()) return;
    createGroupMut.mutate(newGroupName.trim(), {
      onSuccess: () => setNewGroupName(''),
    });
  };

  const handleLeaveGroup = (e: React.MouseEvent, id: string) => {
    e.stopPropagation();
    if (!confirm('Bu gruptan çıkmak istediğine emin misin?')) return;
    leaveGroupMut.mutate(id, {
      onSuccess: () => { if (activeChat?.id === id) setActiveChat(null); },
    });
  };

  const handleLogout = async () => {
    const ws = useChatStore.getState().ws;
    if (ws) ws.close();
    useChatStore.getState().reset();
    await logout();
  };

  return (
    <div className="glass w-80 flex flex-col gap-4 p-5 animate-in">
      {/* Profile */}
      <div className="flex items-center gap-3 pb-4 border-b border-border">
        <Avatar name={username} online />
        <div className="flex-1 min-w-0">
          <div className="font-semibold text-white truncate">{username}</div>
          <div className="text-xs text-success flex items-center gap-1">
            <span className="w-1.5 h-1.5 rounded-full bg-success inline-block" /> Online
          </div>
        </div>
        <IconButton onClick={handleLogout} color="danger" title="Çıkış Yap">
          <LogOut size={18} />
        </IconButton>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 bg-black/20 p-1 rounded-lg">
        {(['users', 'groups'] as const).map((t) => (
          <button key={t} onClick={() => { setTab(t); setShowAdd(false); setShowGroupPanel(false); }}
            className={`flex-1 py-2 text-sm rounded-md transition-all cursor-pointer
              ${tab === t ? 'bg-primary text-white' : 'text-slate-400 hover:text-white'}`}>
            {t === 'users' ? 'Kişiler' : 'Gruplar'}
          </button>
        ))}
      </div>

      {/* ─── USERS TAB ─── */}
      {tab === 'users' && (
        <div className="flex flex-col gap-2.5 flex-1 overflow-hidden">
          <Button variant={showAdd ? 'danger' : 'primary'} className="w-full text-sm"
            icon={showAdd ? <X size={14} /> : <UserPlus size={14} />}
            onClick={() => { setShowAdd(!showAdd); setSearch(''); }}>
            {showAdd ? 'Kapat' : 'Kişi Ekle'}
          </Button>

          {showAdd && (
            <div className="bg-black/15 border border-border rounded-xl p-2.5 flex flex-col gap-2 animate-in">
              <Input search placeholder="Kullanıcı ara..." value={search} onChange={(e) => setSearch(e.target.value)} autoFocus />
              <div className="max-h-44 overflow-y-auto flex flex-col gap-0.5">
                {filtered.length === 0 && (
                  <p className="text-center text-slate-400 text-xs py-4">
                    {search.trim() ? 'Bulunamadı' : 'Eklenecek kişi yok'}
                  </p>
                )}
                {filtered.map((u: User) => (
                  <div key={u.id} className="flex items-center gap-2.5 px-2.5 py-2 rounded-lg hover:bg-white/5 transition-colors">
                    <Avatar name={u.name} size="sm" online={onlineUsers.has(u.id)} />
                    <span className="flex-1 text-sm truncate text-white">{u.name}</span>
                    <IconButton color="success" onClick={() => addMut.mutate(u.id)}
                      loading={addMut.isPending && addMut.variables === u.id}>
                      <UserPlus size={14} />
                    </IconButton>
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className="flex-1 overflow-y-auto flex flex-col gap-1">
            {contacts.length === 0 && !showAdd && (
              <div className="text-center text-slate-400 text-sm py-6 flex flex-col items-center gap-2">
                <UserPlus size={28} className="opacity-30" />
                <span>Henüz kişi yok</span>
              </div>
            )}
            {contacts.map((c: User) => {
              const online = onlineUsers.has(c.id);
              return (
                <div key={c.id} onClick={() => openChat(c)}
                  className={`group flex items-center gap-3 p-2.5 rounded-lg cursor-pointer transition-all
                    ${activeChat?.id === c.id ? 'bg-primary/20 border border-primary/30' : 'hover:bg-white/5'}`}>
                  <Avatar name={c.name} size="md" online={online} />
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium text-white truncate">{c.name}</div>
                    <div className={`text-xs ${online ? 'text-success' : 'text-slate-500'}`}>
                      {online ? 'Çevrimiçi' : 'Çevrimdışı'}
                    </div>
                  </div>
                  <IconButton color="dangerHidden" onClick={(e) => handleRemoveContact(e, c.id)}
                    loading={removeMut.isPending && removeMut.variables === c.id}>
                    <UserMinus size={14} />
                  </IconButton>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* ─── GROUPS TAB ─── */}
      {tab === 'groups' && (
        <div className="flex flex-col gap-2.5 flex-1 overflow-hidden">
          <Button variant={showGroupPanel ? 'danger' : 'primary'} className="w-full text-sm"
            icon={showGroupPanel ? <X size={14} /> : <Plus size={14} />}
            onClick={() => { setShowGroupPanel(!showGroupPanel); setNewGroupName(''); setGroupSearch(''); }}>
            {showGroupPanel ? 'Kapat' : 'Grup Oluştur / Katıl'}
          </Button>

          {showGroupPanel && (
            <div className="bg-black/15 border border-border rounded-xl p-2.5 flex flex-col gap-3 animate-in">
              {/* Create Group */}
              <div>
                <p className="text-xs text-slate-400 mb-1.5 font-medium">Yeni Grup Oluştur</p>
                <div className="flex gap-2">
                  <Input placeholder="Grup adı..." value={newGroupName} onChange={(e) => setNewGroupName(e.target.value)} />
                  <Button className="shrink-0 text-sm" onClick={handleCreateGroup}
                    loading={createGroupMut.isPending} disabled={!newGroupName.trim()}>
                    Oluştur
                  </Button>
                </div>
              </div>

              {/* Join Existing Group */}
              <div>
                <p className="text-xs text-slate-400 mb-1.5 font-medium">Gruba Katıl</p>
                <Input search placeholder="Grup ara..." value={groupSearch} onChange={(e) => setGroupSearch(e.target.value)} />
                <div className="max-h-36 overflow-y-auto flex flex-col gap-0.5 mt-2">
                  {filteredJoinable.length === 0 && (
                    <p className="text-center text-slate-400 text-xs py-3">
                      {groupSearch.trim() ? 'Grup bulunamadı' : 'Katılınacak grup yok'}
                    </p>
                  )}
                  {filteredJoinable.map((g: Group) => (
                    <div key={g.id} className="flex items-center gap-2.5 px-2.5 py-2 rounded-lg hover:bg-white/5 transition-colors">
                      <Avatar name={g.name} size="sm" variant="group">
                        <Users size={14} color="white" />
                      </Avatar>
                      <span className="flex-1 text-sm truncate text-white">{g.name}</span>
                      <IconButton color="success" onClick={() => joinGroupMut.mutate(g.id)}
                        loading={joinGroupMut.isPending && joinGroupMut.variables === g.id}>
                        <Plus size={14} />
                      </IconButton>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {/* My Groups List */}
          <div className="flex-1 overflow-y-auto flex flex-col gap-1">
            {myGroups.length === 0 && !showGroupPanel && (
              <div className="text-center text-slate-400 text-sm py-6 flex flex-col items-center gap-2">
                <Users size={28} className="opacity-30" />
                <span>Henüz grubun yok</span>
              </div>
            )}
            {myGroups.map((g: Group) => (
              <div key={g.id} onClick={() => openGroup(g)}
                className={`group flex items-center gap-3 p-2.5 rounded-lg cursor-pointer transition-all
                  ${activeChat?.id === g.id ? 'bg-primary/20 border border-primary/30' : 'hover:bg-white/5'}`}>
                <Avatar name={g.name} variant="group">
                  <Users size={16} color="white" />
                </Avatar>
                <span className="flex-1 text-sm text-white truncate">{g.name}</span>
                <IconButton color="dangerHidden" onClick={(e) => handleLeaveGroup(e, g.id)}
                  loading={leaveGroupMut.isPending && leaveGroupMut.variables === g.id}
                  title="Gruptan Çık">
                  <DoorOpen size={14} />
                </IconButton>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
