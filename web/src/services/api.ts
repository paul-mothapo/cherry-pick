import axios, { AxiosResponse } from 'axios';
import {
  DatabaseConnection,
  DatabaseReport,
  SecurityIssue,
  QueryOptimization,
  Alert,
  DataLineage,
} from '@/types/database';

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export const connectionApi = {
  getConnections: (): Promise<AxiosResponse<ApiResponse<DatabaseConnection[]>>> =>
    api.get('/connections'),
  
  createConnection: (connection: Omit<DatabaseConnection, 'id' | 'status'>): Promise<AxiosResponse<ApiResponse<DatabaseConnection>>> =>
    api.post('/connections', connection),
  
  testConnection: (connectionId: string): Promise<AxiosResponse<ApiResponse<{ status: string }>>> =>
    api.post(`/connections/${connectionId}/test`),
  
  deleteConnection: (connectionId: string): Promise<AxiosResponse<ApiResponse<void>>> =>
    api.delete(`/connections/${connectionId}`),
};

export const analysisApi = {
  analyzeDatabase: (connectionId: string): Promise<AxiosResponse<ApiResponse<DatabaseReport>>> =>
    api.post(`/analysis/${connectionId}/analyze`),
  
  getReport: (connectionId: string): Promise<AxiosResponse<ApiResponse<DatabaseReport>>> =>
    api.get(`/analysis/${connectionId}/report`),
  
  getReports: (): Promise<AxiosResponse<ApiResponse<DatabaseReport[]>>> =>
    api.get('/analysis/reports'),
};

export const securityApi = {
  analyzeSecurity: (connectionId: string): Promise<AxiosResponse<ApiResponse<SecurityIssue[]>>> =>
    api.post(`/security/${connectionId}/analyze`),
  
  getSecurityIssues: (connectionId: string): Promise<AxiosResponse<ApiResponse<SecurityIssue[]>>> =>
    api.get(`/security/${connectionId}/issues`),
};

export const optimizationApi = {
  optimizeQuery: (connectionId: string, query: string): Promise<AxiosResponse<ApiResponse<QueryOptimization>>> =>
    api.post(`/optimization/${connectionId}/optimize`, { query }),
  
  getOptimizationHistory: (connectionId: string): Promise<AxiosResponse<ApiResponse<QueryOptimization[]>>> =>
    api.get(`/optimization/${connectionId}/history`),
};

export const monitoringApi = {
  getAlerts: (connectionId?: string): Promise<AxiosResponse<ApiResponse<Alert[]>>> =>
    api.get('/monitoring/alerts', { params: { connectionId } }),
  
  acknowledgeAlert: (alertId: string): Promise<AxiosResponse<ApiResponse<void>>> =>
    api.post(`/monitoring/alerts/${alertId}/acknowledge`),
  
  getMetrics: (connectionId: string): Promise<AxiosResponse<ApiResponse<any>>> =>
    api.get(`/monitoring/${connectionId}/metrics`),
};

export const lineageApi = {
  getLineage: (connectionId: string): Promise<AxiosResponse<ApiResponse<DataLineage>>> =>
    api.get(`/lineage/${connectionId}`),
  
  trackLineage: (connectionId: string): Promise<AxiosResponse<ApiResponse<DataLineage>>> =>
    api.post(`/lineage/${connectionId}/track`),
};

export const collectionApi = {
  getCollectionData: (connectionId: string, collectionName: string, page = 1, limit = 20): Promise<AxiosResponse<ApiResponse<any>>> =>
    api.get(`/collections/${connectionId}/${collectionName}/data`, { params: { page, limit } }),
  
  getCollectionStats: (connectionId: string, collectionName: string): Promise<AxiosResponse<ApiResponse<any>>> =>
    api.get(`/collections/${connectionId}/${collectionName}/stats`),
  
  searchCollection: (connectionId: string, collectionName: string, query: string): Promise<AxiosResponse<ApiResponse<any>>> =>
    api.post(`/collections/${connectionId}/${collectionName}/search`, { query }),
};

export { api };
