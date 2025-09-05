import React from "react";
import { Plus, Database } from "lucide-react";
import { Card, CardContent } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";

interface EmptyStateProps {
  onAddConnection: () => void;
}

export default function EmptyState({ onAddConnection }: EmptyStateProps) {
  return (
    <Card className="bg-seasalt">
      <CardContent className="p-12 text-center">
        <Database className="h-12 w-12 text-onyx mx-auto mb-4" />
        <h3 className="text-lg font-medium text-eerie-black mb-2">
          No connections yet
        </h3>
        <p className="text-onyx mb-4">
          Add your first database connection to get started with analysis.
        </p>
        <Button
          className="bg-eerie-black"
          onClick={onAddConnection}
        >
          <Plus className="mr-2 h-4 w-4" />
          Add Connection
        </Button>
      </CardContent>
    </Card>
  );
}
