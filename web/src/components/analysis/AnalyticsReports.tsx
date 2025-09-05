import { useNavigate } from "react-router-dom";
import { useConnections } from "@/hooks/useConnections";
import { useQuery } from "@tanstack/react-query";
import { analysisApi } from "@/services/api";
import { useAppStore } from "@/stores/useAppStore";
import ReportHeader from "./ReportHeader";
import ReportCard from "./ReportCard";
import EmptyState from "./EmptyState";

export default function AnalyticsReports() {
  const navigate = useNavigate();
  const { data: connections } = useConnections();
  const { currentConnection } = useAppStore();

  const connectedConnections =
    connections?.filter((c) => c.status === "connected") || [];

  const reportsQuery = useQuery({
    queryKey: ["reports"],
    queryFn: async () => {
      const response = await analysisApi.getReports();
      return response.data.data;
    },
  });

  const handleCollectionClick = (collection: any, connectionId: string) => {
    navigate(
      `/collection/${connectionId}/${encodeURIComponent(collection.name)}`
    );
  };

  const handleRunAnalysis = () => {
    window.location.reload();
  };

  return (
    <div className="space-y-6 mt-4">
      <ReportHeader 
        onRefresh={() => reportsQuery.refetch()} 
        isRefreshing={reportsQuery.isRefetching} 
      />

      {reportsQuery.isLoading ? (
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="animate-pulse">
              <div className="h-32 bg-gray-200 rounded"></div>
            </div>
          ))}
        </div>
      ) : reportsQuery.data && reportsQuery.data.length > 0 ? (
        <div className="space-y-4">
          {reportsQuery.data.map((report, index) => (
            <ReportCard
              key={index}
              report={report}
              onCollectionClick={handleCollectionClick}
              currentConnectionId={currentConnection?.id}
            />
          ))}
        </div>
      ) : (
        <EmptyState
          connectedConnectionsCount={connectedConnections.length}
          currentConnectionName={currentConnection?.name}
          onRunAnalysis={handleRunAnalysis}
        />
      )}
    </div>
  );
}
