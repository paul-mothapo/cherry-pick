import React from "react";
import { Card, CardContent } from "@/components/ui/Card";
import { useConnections } from "@/hooks/useConnections";
import { useDashboardStats, useRecentActivity } from "@/hooks/useDashboard";
import { createStatsArray } from "@/constants/DashboardStats";
import QuickActions from "@/components/dashboard/QuickActions";
import RecentActivity from "@/components/dashboard/RecentActivity";

export const Dashboard: React.FC = () => {
  const { isLoading: connectionsLoading } = useConnections();
  const {
    stats: dashboardStats,
    isLoading: statsLoading,
    refetch,
  } = useDashboardStats();
  const {
    activities,
    formatRelativeTime,
    isLoading: activitiesLoading,
  } = useRecentActivity();

  const stats = createStatsArray(dashboardStats);

  const isLoading = connectionsLoading || statsLoading;

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-platinum rounded w-1/4 mb-6"></div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-32 bg-platinum rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-eerie-black">Dashboard</h1>
        <p className="text-onyx mt-2">
          Overview of your database intelligence system
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8 border p-6 rounded-xl">
        {stats.map((stat, index) => (
          <Card key={index} className="bg-seasalt">
            <CardContent>
              <div className="flex items-center">
                <div>
                  <p className="text-sm font-medium text-onyx">{stat.title}</p>
                  <p className="text-2xl font-bold text-eerie-black">
                    {stat.value}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <QuickActions dashboardStats={dashboardStats} refetch={refetch} />
        <RecentActivity
          activities={activities}
          activitiesLoading={activitiesLoading}
          formatRelativeTime={formatRelativeTime}
        />
      </div>
    </div>
  );
};
