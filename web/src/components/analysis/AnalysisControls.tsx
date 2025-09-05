import React, { useState } from "react";
import { CheckCircle, Database, Play } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { Button } from "../ui/Button";
import { Select } from "../ui/Select";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/Card";
import { useConnections } from "@/hooks/useConnections";
import { useMutation } from "@tanstack/react-query";
import { analysisApi } from "@/services/api";
import { useAppStore } from "@/stores/useAppStore";

export default function AnalysisControls() {
  const [selectedConnectionId, setSelectedConnectionId] = useState<string>("");
  const navigate = useNavigate();
  const { data: connections } = useConnections();
  const { currentConnection } = useAppStore();

  const connectedConnections =
    connections?.filter((c) => c.status === "connected") || [];

  const connectionOptions = connectedConnections.map((conn) => ({
    value: conn.id,
    label: `${conn.name} (${conn.driver.toUpperCase()})`,
  }));

  React.useEffect(() => {
    if (
      currentConnection?.id &&
      currentConnection.status === "connected" &&
      !selectedConnectionId
    ) {
      setSelectedConnectionId(currentConnection.id);
    }
  }, [currentConnection, selectedConnectionId]);

  const effectiveConnectionId =
    selectedConnectionId || currentConnection?.id || "";
  const hasCurrentConnection =
    !!currentConnection && currentConnection.status === "connected";
  const selectedConnection =
    connections?.find((c) => c.id === selectedConnectionId) ||
    currentConnection;

  const runAnalysis = useMutation({
    mutationFn: (connectionId: string) =>
      analysisApi.analyzeDatabase(connectionId),
  });

  const handleRunAnalysis = () => {
    if (effectiveConnectionId) {
      runAnalysis.mutate(effectiveConnectionId);
    }
  };
  return (
    <div className="p-4 border rounded-xl">
      <Card className="bg-seasalt">
        <CardHeader>
          <CardTitle>Run New Analysis</CardTitle>
        </CardHeader>
        <CardContent>
          {connectedConnections.length === 0 ? (
            <div className="text-center py-8">
              <Database className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                No Connected Databases
              </h3>
              <p className="text-onyx mb-4">
                Please connect to a database first from the Connections page.
              </p>
              <Button
                variant="secondary"
                onClick={() => navigate("/connections")}
              >
                Go to Connections
              </Button>
            </div>
          ) : (
            <div className="space-y-4">
              {hasCurrentConnection && (
                <div className="bg-primary-50 border border-primary-200 rounded-lg p-4">
                  <div className="flex items-center">
                    <Database className="h-5 w-5 text-primary-600 mr-3" />
                    <div>
                      <h4 className="font-medium text-primary-900">
                        Active Connection: {currentConnection.name}
                      </h4>
                      <p className="text-sm text-primary-700">
                        {currentConnection.driver.toUpperCase()} •{" "}
                        {selectedConnectionId === currentConnection.id
                          ? "Selected for analysis"
                          : "Available for analysis"}
                      </p>
                    </div>
                  </div>
                </div>
              )}

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Select
                  label="Select Database Connection"
                  value={selectedConnectionId}
                  onChange={(e) => setSelectedConnectionId(e.target.value)}
                  options={connectionOptions}
                  placeholder="Choose a connected database"
                />
                <div className="flex items-end">
                  <Button
                    onClick={handleRunAnalysis}
                    loading={runAnalysis.isPending}
                    disabled={!effectiveConnectionId}
                    className="w-full"
                  >
                    <Play className="mr-2 h-4 w-4" />
                    Run Analysis{" "}
                    {selectedConnection ? `on ${selectedConnection.name}` : ""}
                  </Button>
                </div>
              </div>

              {runAnalysis.isError && (
                <div className="text-red-600 text-sm bg-red-50 p-3 rounded-lg">
                  <strong>Error:</strong>{" "}
                  {runAnalysis.error?.message || "Failed to run analysis"}
                </div>
              )}
              {runAnalysis.isSuccess && (
                <div className="text-green-600 text-sm flex items-center bg-green-50 p-3 rounded-lg">
                  <CheckCircle className="mr-2 h-4 w-4" />
                  Analysis completed successfully
                  {selectedConnection ? ` for ${selectedConnection.name}` : ""}!
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
