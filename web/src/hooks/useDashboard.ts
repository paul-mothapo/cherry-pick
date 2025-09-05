import { useQuery } from '@tanstack/react-query';
import { monitoringApi, securityApi } from '@/services/api';
import { useConnections } from './useConnections';

export interface DashboardStats {
  activeConnections: number;
  securityIssues: number;
  performanceScore: string;
  activeAlerts: number;
}

export interface ActivityItem {
  id: string;
  message: string;
  timestamp: string;
  type: 'success' | 'warning' | 'info' | 'error';
}

export const useDashboardStats = () => {
  const { data: connections } = useConnections();

  // Get active connections count
  const activeConnections = connections?.filter(c => c.status === 'connected').length || 0;

  // Get security issues for all connected connections
  const securityQuery = useQuery({
    queryKey: ['dashboard-security'],
    queryFn: async () => {
      if (!connections || activeConnections === 0) return [];
      
      const connectedConnections = connections.filter(c => c.status === 'connected');
      const securityPromises = connectedConnections.map(conn => 
        securityApi.getSecurityIssues(conn.id).catch(() => ({ data: { data: [] } }))
      );
      
      const results = await Promise.all(securityPromises);
      return results.flatMap(result => result.data.data || []);
    },
    enabled: activeConnections > 0,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  // Get alerts for all connected connections
  const alertsQuery = useQuery({
    queryKey: ['dashboard-alerts'],
    queryFn: async () => {
      const response = await monitoringApi.getAlerts();
      return response.data.data || [];
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
  });

  // Get metrics for performance score (use first connected connection as sample)
  const firstConnectedId = connections?.find(c => c.status === 'connected')?.id;
  const metricsQuery = useQuery({
    queryKey: ['dashboard-metrics', firstConnectedId],
    queryFn: async () => {
      if (!firstConnectedId) return null;
      const response = await monitoringApi.getMetrics(firstConnectedId);
      return response.data.data;
    },
    enabled: !!firstConnectedId,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  // Calculate performance score from metrics
  const calculatePerformanceScore = (metrics: any) => {
    if (!metrics) return '0';
    
    // Simple performance calculation based on available metrics
    const cpuScore = Math.max(0, 100 - (metrics.cpu_usage || 0));
    const memoryScore = Math.max(0, 100 - (metrics.memory_usage || 0));
    const avgScore = (cpuScore + memoryScore) / 2;
    
    return Math.round(avgScore).toString();
  };

  const stats: DashboardStats = {
    activeConnections,
    securityIssues: securityQuery.data?.length || 0,
    performanceScore: calculatePerformanceScore(metricsQuery.data) + '%',
    activeAlerts: alertsQuery.data?.length || 0,
  };

  const isLoading = securityQuery.isLoading || alertsQuery.isLoading || metricsQuery.isLoading;
  const error = securityQuery.error || alertsQuery.error || metricsQuery.error;

  return {
    stats,
    isLoading,
    error,
    refetch: () => {
      securityQuery.refetch();
      alertsQuery.refetch();
      metricsQuery.refetch();
    },
  };
};

export const useRecentActivity = () => {
  const { data: connections } = useConnections();
  const { data: alerts } = useQuery({
    queryKey: ['recent-alerts'],
    queryFn: async () => {
      const response = await monitoringApi.getAlerts();
      return response.data.data || [];
    },
    staleTime: 2 * 60 * 1000,
  });

  // Generate recent activity from connections and alerts
  const generateActivity = (): ActivityItem[] => {
    const activities: ActivityItem[] = [];

    // Add connection activities
    connections?.slice(0, 2).forEach((conn) => {
      if (conn.lastConnected) {
        activities.push({
          id: `conn-${conn.id}`,
          message: `Connection ${conn.status === 'connected' ? 'established' : 'updated'} for ${conn.name}`,
          timestamp: conn.lastConnected,
          type: conn.status === 'connected' ? 'success' : 'info',
        });
      }
    });

    // Add alert activities
    alerts?.slice(0, 2).forEach((alert: any, index) => {
      activities.push({
        id: `alert-${index}`,
        message: alert.message || `${alert.type} alert detected`,
        timestamp: alert.timestamp || new Date().toISOString(),
        type: alert.severity === 'high' ? 'error' : alert.severity === 'medium' ? 'warning' : 'info',
      });
    });

    // Sort by timestamp (most recent first)
    return activities
      .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
      .slice(0, 5);
  };

  const formatRelativeTime = (timestamp: string) => {
    const now = new Date();
    const time = new Date(timestamp);
    const diffInMinutes = Math.floor((now.getTime() - time.getTime()) / (1000 * 60));

    if (diffInMinutes < 1) return 'Just now';
    if (diffInMinutes < 60) return `${diffInMinutes} minute${diffInMinutes === 1 ? '' : 's'} ago`;
    
    const diffInHours = Math.floor(diffInMinutes / 60);
    if (diffInHours < 24) return `${diffInHours} hour${diffInHours === 1 ? '' : 's'} ago`;
    
    const diffInDays = Math.floor(diffInHours / 24);
    return `${diffInDays} day${diffInDays === 1 ? '' : 's'} ago`;
  };

  const activities = generateActivity();

  return {
    activities,
    formatRelativeTime,
    isLoading: false,
  };
};
