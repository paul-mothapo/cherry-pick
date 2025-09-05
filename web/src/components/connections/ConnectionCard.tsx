import { Database, Trash2, TestTube, Unplug } from "lucide-react";
import { Card, CardContent } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { DatabaseConnection } from "@/types/database";
import { useTestConnection, useDeleteConnection } from "@/hooks/useConnections";
import { useAppStore } from "@/stores/useAppStore";

interface ConnectionCardProps {
  connection: DatabaseConnection;
}

export default function ConnectionCard({ connection }: ConnectionCardProps) {
  const testConnection = useTestConnection();
  const deleteConnection = useDeleteConnection();
  const { currentConnection, setCurrentConnection } = useAppStore();

  const handleTest = async (connectionId: string) => {
    try {
      await testConnection.mutateAsync(connectionId);
    } catch (error) {
      console.error("Failed to test connection:", error);
    }
  };

  const handleDelete = async (connectionId: string) => {
    if (window.confirm("Are you sure you want to delete this connection?")) {
      try {
        await deleteConnection.mutateAsync(connectionId);
        if (currentConnection?.id === connectionId) {
          setCurrentConnection(null);
        }
      } catch (error) {
        console.error("Failed to delete connection:", error);
      }
    }
  };

  const handleConnect = (connection: DatabaseConnection) => {
    setCurrentConnection(connection);
  };

  const handleDisconnect = () => {
    setCurrentConnection(null);
  };

  return (
    <Card className="bg-seasalt">
      <CardContent className="p-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <div className="p-2 bg-platinum rounded-lg mr-4">
              <Database className="h-6 w-6 text-onyx" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-eerie-black">
                {connection.name}
              </h3>
              <p className="text-sm text-onyx">
                {connection.driver.toUpperCase()} • {connection.status}
              </p>
              {connection.lastConnected && (
                <p className="text-xs text-onyx">
                  Last connected:{" "}
                  {new Date(connection.lastConnected).toLocaleString()}
                </p>
              )}
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <div
              className={`h-3 w-3 rounded-full ${
                connection.status === "connected"
                  ? "bg-green-400"
                  : connection.status === "error"
                  ? "bg-turkey-red"
                  : "bg-platinum"
              }`}
            />
            {currentConnection?.id === connection.id ? (
              <div className="flex items-center space-x-2">
                <span className="text-sm font-medium text-onyx bg-platinum px-3 py-1 rounded-full">
                  Active
                </span>
                <Button
                  variant="secondary"
                  size="sm"
                  onClick={handleDisconnect}
                  className="text-onyx hover:text-eerie-black"
                >
                  <Unplug className="mr-2 h-4 w-4" />
                  Disconnect
                </Button>
              </div>
            ) : (
              <Button
                variant="primary"
                size="sm"
                onClick={() => handleConnect(connection)}
                disabled={connection.status !== "connected"}
              >
                Connect
              </Button>
            )}
            <Button
              variant="secondary"
              size="sm"
              onClick={() => handleTest(connection.id)}
              loading={testConnection.isPending}
            >
              <TestTube className="mr-2 h-4 w-4" />
              Test
            </Button>
            <Button
              variant="danger"
              size="sm"
              onClick={() => handleDelete(connection.id)}
              loading={deleteConnection.isPending}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
