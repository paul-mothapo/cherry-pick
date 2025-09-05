import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { analysisApi } from '@/services/api';
import { useAppStore } from '@/stores/useAppStore';

export const useAnalyzeDatabase = () => {
  const queryClient = useQueryClient();
  const { setError, clearError } = useAppStore();
  
  return useMutation({
    mutationFn: async (connectionId: string) => {
      try {
        clearError();
        const response = await analysisApi.analyzeDatabase(connectionId);
        return response.data.data;
      } catch (error) {
        setError('Failed to analyze database');
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reports'] });
      queryClient.invalidateQueries({ queryKey: ['analysis'] });
    },
  });
};

export const useReports = () => {
  const { setError, clearError } = useAppStore();
  
  return useQuery({
    queryKey: ['reports'],
    queryFn: async () => {
      try {
        clearError();
        const response = await analysisApi.getReports();
        return response.data.data;
      } catch (error) {
        setError('Failed to fetch reports');
        throw error;
      }
    },
  });
};

export const useReport = (connectionId: string) => {
  const { setError, clearError } = useAppStore();
  
  return useQuery({
    queryKey: ['report', connectionId],
    queryFn: async () => {
      try {
        clearError();
        const response = await analysisApi.getReport(connectionId);
        return response.data.data;
      } catch (error) {
        setError('Failed to fetch report');
        throw error;
      }
    },
    enabled: !!connectionId,
  });
};
