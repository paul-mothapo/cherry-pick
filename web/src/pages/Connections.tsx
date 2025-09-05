import React, { useState } from "react";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { useConnections } from "@/hooks/useConnections";
import ConnectionForm from "@/components/connections/ConnectionForm";
import ConnectionCard from "@/components/connections/ConnectionCard";
import EmptyState from "@/components/connections/EmptyState";

export const Connections: React.FC = () => {
  const [showAddForm, setShowAddForm] = useState(false);
  const { data: connections, isLoading } = useConnections();

  if (isLoading) {
    return (
      <div className="p-6">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="h-24 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold text-eerie-black">
            Database Connections
          </h1>
          <p className="text-onyx mt-2">Manage your database connections</p>
        </div>
        <Button onClick={() => setShowAddForm(true)} variant="primary">
          <Plus className="mr-2 h-4 w-4" />
          Add Connection
        </Button>
      </div>

      {showAddForm && (
        <ConnectionForm onClose={() => setShowAddForm(false)} />
      )}

      <div>
        {connections?.map((connection) => (
          <ConnectionCard key={connection.id} connection={connection} />
        ))}

        {connections?.length === 0 && !showAddForm && (
          <EmptyState onAddConnection={() => setShowAddForm(true)} />
        )}
      </div>
    </div>
  );
};
