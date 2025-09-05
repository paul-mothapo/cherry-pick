import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { connectionApi } from '@/services/api';
import { DatabaseConnection } from '@/types/database';
import { useAppStore } from '@/stores/useAppStore';

export const useConnections = () => {
  const { setError, clearError } = useAppStore();
  
  return useQuery({
    queryKey: ['connections'],
    queryFn: async () => {
      try {
        clearError();
        const response = await connectionApi.getConnections();
        return response.data.data;
      } catch (error) {
        setError('Failed to fetch connections');
        throw error;
      }
    },
  });
};

export const useCreateConnection = () => {
  const queryClient = useQueryClient();
  const { setError, clearError } = useAppStore();
  
  return useMutation({
    mutationFn: async (connection: Omit<DatabaseConnection, 'id' | 'status'>) => {
      try {
        clearError();
        const response = await connectionApi.createConnection(connection);
        return response.data.data;
      } catch (error) {
        setError('Failed to create connection');
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['connections'] });
    },
  });
};

export const useTestConnection = () => {
  const { setError, clearError } = useAppStore();
  
  return useMutation({
    mutationFn: async (connectionId: string) => {
      try {
        clearError();
        const response = await connectionApi.testConnection(connectionId);
        return response.data.data;
      } catch (error) {
        setError('Failed to test connection');
        throw error;
      }
    },
  });
};

export const useDeleteConnection = () => {
  const queryClient = useQueryClient();
  const { setError, clearError } = useAppStore();
  
  return useMutation({
    mutationFn: async (connectionId: string) => {
      try {
        clearError();
        await connectionApi.deleteConnection(connectionId);
      } catch (error) {
        setError('Failed to delete connection');
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['connections'] });
    },
  });
};
