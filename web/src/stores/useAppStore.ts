import { create } from 'zustand';
import { DatabaseConnection } from '@/types/database';

interface AppState {
  sidebarOpen: boolean;
  sidebarCollapsed: boolean;
  currentConnection: DatabaseConnection | null;
  loading: boolean;
  error: string | null;
  
  setSidebarOpen: (open: boolean) => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  setCurrentConnection: (connection: DatabaseConnection | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useAppStore = create<AppState>((set) => ({
  sidebarOpen: true,
  sidebarCollapsed: false,
  currentConnection: null,
  loading: false,
  error: null,
  
  setSidebarOpen: (open) => set({ sidebarOpen: open }),
  setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
  setCurrentConnection: (connection) => set({ currentConnection: connection }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  clearError: () => set({ error: null }),
}));
