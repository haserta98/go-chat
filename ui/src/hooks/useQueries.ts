import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../api/client';
import { useAuthStore } from '../stores/authStore';

export const useContacts = () => {
  const isLoggedIn = useAuthStore((s) => s.isLoggedIn);
  return useQuery({
    queryKey: ['contacts'],
    queryFn: () => apiClient.fetchContacts(),
    enabled: isLoggedIn,
    select: (res) => res.data || [],
  });
};

export const useUsers = () => {
  const isLoggedIn = useAuthStore((s) => s.isLoggedIn);
  return useQuery({
    queryKey: ['users'],
    queryFn: () => apiClient.fetchUsers(),
    enabled: isLoggedIn,
    select: (res) => res.data || [],
  });
};

export const useMyGroups = () => {
  const isLoggedIn = useAuthStore((s) => s.isLoggedIn);
  return useQuery({
    queryKey: ['myGroups'],
    queryFn: () => apiClient.fetchMyGroups(),
    enabled: isLoggedIn,
    select: (res) => res.data || [],
  });
};

export const useOnlineStatus = () => {
  const isLoggedIn = useAuthStore((s) => s.isLoggedIn);
  return useQuery({
    queryKey: ['onlineStatus'],
    queryFn: () => apiClient.fetchContactsOnlineStatus(),
    enabled: isLoggedIn,
    select: (res) => {
      const online = new Set<string>();
      if (res.data) {
        for (const [id, isOnline] of Object.entries(res.data)) {
          if (isOnline) online.add(id);
        }
      }
      return online;
    },
  });
};

export const useAddContact = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (contactId: string) => apiClient.addContact(contactId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['contacts'] });
      qc.invalidateQueries({ queryKey: ['users'] });
    },
  });
};

export const useRemoveContact = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (contactId: string) => apiClient.removeContact(contactId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['contacts'] });
    },
  });
};

export const useAllGroups = () => {
  const isLoggedIn = useAuthStore((s) => s.isLoggedIn);
  return useQuery({
    queryKey: ['allGroups'],
    queryFn: () => apiClient.fetchAllGroups(),
    enabled: isLoggedIn,
    select: (res) => res.data || [],
  });
};

export const useCreateGroup = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (name: string) => apiClient.createGroup(name),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['myGroups'] });
      qc.invalidateQueries({ queryKey: ['allGroups'] });
    },
  });
};

export const useJoinGroup = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (groupId: string) => apiClient.joinGroup(groupId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['myGroups'] });
    },
  });
};

export const useLeaveGroup = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (groupId: string) => apiClient.leaveGroup(groupId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['myGroups'] });
    },
  });
};
