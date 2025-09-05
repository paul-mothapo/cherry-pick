import { Database, Shield, BarChart3, AlertTriangle } from "lucide-react";

export interface DashboardStatsData {
  activeConnections: number;
  securityIssues: number;
  performanceScore: string;
  activeAlerts: number;
}

export const createStatsArray = (dashboardStats: DashboardStatsData) => [
  {
    title: "Active Connections",
    value: dashboardStats.activeConnections,
    icon: Database,
    color: "text-blue-600",
    bgColor: "bg-blue-100",
  },
  {
    title: "Security Issues",
    value: dashboardStats.securityIssues,
    icon: Shield,
    color:
      dashboardStats.securityIssues > 0 ? "text-red-600" : "text-green-600",
    bgColor:
      dashboardStats.securityIssues > 0 ? "bg-red-100" : "bg-green-100",
  },
  {
    title: "Performance Score",
    value: dashboardStats.performanceScore,
    icon: BarChart3,
    color: "text-green-600",
    bgColor: "bg-green-100",
  },
  {
    title: "Active Alerts",
    value: dashboardStats.activeAlerts,
    icon: AlertTriangle,
    color:
      dashboardStats.activeAlerts > 0 ? "text-yellow-600" : "text-green-600",
    bgColor:
      dashboardStats.activeAlerts > 0 ? "bg-yellow-100" : "bg-green-100",
  },
];