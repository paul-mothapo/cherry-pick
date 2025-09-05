import React, { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Select } from "@/components/ui/Select";
import { useCreateConnection } from "@/hooks/useConnections";
import { DatabaseConnection } from "@/types/database";
import { databaseOptions } from "@/constants/DatabaseOptions";

interface ConnectionFormProps {
  onClose: () => void;
}

export default function ConnectionForm({ onClose }: ConnectionFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    driver: "",
    connectionString: "",
  });

  const createConnection = useCreateConnection();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createConnection.mutateAsync(
        formData as Omit<DatabaseConnection, "id" | "status">
      );
      setFormData({ name: "", driver: "", connectionString: "" });
      onClose();
    } catch (error) {
      console.error("Failed to create connection:", error);
    }
  };

  return (
    <div className="p-4 border rounded-xl">
      <Card className="bg-seasalt">
        <CardHeader>
          <CardTitle>Add New Connection</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Input
                label="Connection Name"
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                placeholder="My Database"
                required
              />
              <Select
                label="Database Type"
                value={formData.driver}
                onChange={(e) =>
                  setFormData({ ...formData, driver: e.target.value })
                }
                options={databaseOptions}
                placeholder="Select database type"
                required
              />
            </div>
            <Input
              label="Connection String"
              value={formData.connectionString}
              onChange={(e) =>
                setFormData({ ...formData, connectionString: e.target.value })
              }
              placeholder="mysql://user:password@localhost:3306/database"
              required
            />
            <div className="flex space-x-3">
              <Button
                className="bg-eerie-black"
                type="submit"
                loading={createConnection.isPending}
              >
                Create Connection
              </Button>
              <Button type="button" variant="secondary" onClick={onClose}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
