import { Card, CardContent, CardHeader, CardTitle } from "../ui/Card";

interface ActivityItem {
  id: string;
  message: string;
  timestamp: string;
  type: 'success' | 'warning' | 'info' | 'error';
}

interface RecentActivityProps {
  activities: ActivityItem[];
  activitiesLoading: boolean;
  formatRelativeTime: (timestamp: string) => string;
}

export default function RecentActivity({ 
  activities, 
  activitiesLoading, 
  formatRelativeTime 
}: RecentActivityProps) {
  return (
    <div>
      <Card>
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
        </CardHeader>
        <CardContent>
          {activitiesLoading ? (
            <div className="space-y-4">
              {[...Array(3)].map((_, i) => (
                <div
                  key={i}
                  className="animate-pulse flex items-center space-x-3"
                >
                  <div className="h-2 w-2 bg-gray-300 rounded-full"></div>
                  <div className="flex-1 space-y-1">
                    <div className="h-4 bg-gray-300 rounded w-3/4"></div>
                    <div className="h-3 bg-gray-200 rounded w-1/4"></div>
                  </div>
                </div>
              ))}
            </div>
          ) : activities.length > 0 ? (
            <div className="space-y-4">
              {activities.map((activity) => (
                <div key={activity.id} className="flex items-center space-x-3">
                  <div
                    className={`h-2 w-2 rounded-full ${
                      activity.type === "success"
                        ? "bg-green-400"
                        : activity.type === "warning"
                        ? "bg-yellow-400"
                        : activity.type === "error"
                        ? "bg-red-400"
                        : "bg-blue-400"
                    }`}
                  ></div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-gray-900 truncate">
                      {activity.message}
                    </p>
                    <p className="text-xs text-gray-500">
                      {formatRelativeTime(activity.timestamp)}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <p className="text-gray-500 text-sm">No recent activity</p>
              <p className="text-gray-400 text-xs mt-1">
                Connect to databases and run analyses to see activity
              </p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
