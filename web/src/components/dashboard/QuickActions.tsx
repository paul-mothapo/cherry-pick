import { BarChart3, Database, Shield } from "lucide-react";
import { Button } from "../ui/Button";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/Card";
import { useNavigate } from "react-router-dom";

interface QuickActionsProps {
  dashboardStats: {
    activeConnections: number;
  };
  refetch: () => void;
}

export default function QuickActions({
  dashboardStats,
  refetch,
}: QuickActionsProps) {
  const navigate = useNavigate();
  return (
    <div className="p-4 border rounded-xl bg-seasalt">
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <Button
              variant="primary"
              className="w-full justify-start"
              onClick={() => navigate("/connections")}
            >
              <Database className="mr-2 h-4 w-4" />
              Add New Connection
            </Button>
            <Button
              variant="secondary"
              className="w-full justify-start"
              onClick={() => navigate("/analysis")}
              disabled={dashboardStats.activeConnections === 0}
            >
              <BarChart3 className="mr-2 h-4 w-4" />
              Run Database Analysis
            </Button>
            <Button
              variant="secondary"
              className="w-full justify-start"
              onClick={() => refetch()}
            >
              <Shield className="mr-2 h-4 w-4" />
              Refresh Dashboard
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
