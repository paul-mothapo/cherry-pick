import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/Button";

interface ReportHeaderProps {
  onRefresh: () => void;
  isRefreshing: boolean;
}

export default function ReportHeader({ onRefresh, isRefreshing }: ReportHeaderProps) {
  return (
    <div className="flex justify-between items-center">
      <h2 className="text-xl font-semibold text-gray-900">
        Analysis Reports
      </h2>
      <Button
        variant="secondary"
        size="sm"
        onClick={onRefresh}
        loading={isRefreshing}
      >
        <RefreshCw className="mr-2 h-4 w-4" />
        Refresh
      </Button>
    </div>
  );
}
