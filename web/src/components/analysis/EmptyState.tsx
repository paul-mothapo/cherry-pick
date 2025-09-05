import { Database, Play } from "lucide-react";
import { Card, CardContent } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";

interface EmptyStateProps {
  connectedConnectionsCount: number;
  currentConnectionName?: string;
  onRunAnalysis: () => void;
}

export default function EmptyState({ 
  connectedConnectionsCount, 
  currentConnectionName, 
  onRunAnalysis 
}: EmptyStateProps) {
  return (
    <Card>
      <CardContent className="p-12 text-center">
        <Database className="h-12 w-12 text-gray-400 mx-auto mb-4" />
        <h3 className="text-lg font-medium text-gray-900 mb-2">
          No Analysis Reports Yet
        </h3>
        <p className="text-onyx mb-4">
          Run your first database analysis to see insights and
          recommendations.
        </p>
        {connectedConnectionsCount > 0 && (
          <Button onClick={onRunAnalysis}>
            <Play className="mr-2 h-4 w-4" />
            {currentConnectionName
              ? `Run Analysis on ${currentConnectionName}`
              : "Run Analysis"}
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
