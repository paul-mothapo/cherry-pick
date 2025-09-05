import { Download, Clock } from "lucide-react";
import { Card, CardContent } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import { DatabaseReport } from "@/types/database";
import CollectionsDetails from "./CollectionsDetails";
import InsightsSection from "./InsightsSection";
import RecommendationsSection from "./RecommendationsSection";
import PerformanceMetrics from "./PerformanceMetrics";

interface ReportCardProps {
  report: DatabaseReport;
  onCollectionClick: (collection: any, connectionId: string) => void;
  currentConnectionId?: string;
}

export default function ReportCard({
  report,
  onCollectionClick,
  currentConnectionId,
}: ReportCardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const getHealthScoreColor = (score: number) => {
    if (score >= 80) return "text-green-600";
    if (score >= 60) return "text-yellow-600";
    return "text-turkey-red";
  };

  const getHealthScoreBg = (score: number) => {
    if (score >= 80) return "bg-green-100";
    if (score >= 60) return "bg-yellow-100";
    return "bg-red-100";
  };

  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h3 className="text-lg font-semibold text-eerie-black">
              {report.databaseName}
            </h3>
            <p className="text-sm text-onyx flex items-center">
              <Clock className="mr-1 h-4 w-4" />
              {formatDate(report.analysisTime)}
            </p>
          </div>
          <div className="flex items-center space-x-4">
            <div
              className={`px-3 py-1 rounded-full text-sm font-medium ${getHealthScoreBg(
                report.summary?.healthScore || 0
              )}`}
            >
              <span
                className={getHealthScoreColor(
                  report.summary?.healthScore || 0
                )}
              >
                Health Score: {report.summary?.healthScore || "N/A"}%
              </span>
            </div>
            <Button variant="secondary" size="sm">
              <Download className="mr-2 h-4 w-4" />
              Export
            </Button>
          </div>
        </div>

        {report.summary && (
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-eerie-black">
                {report.summary.totalTables || report.tables?.length || 0}
              </div>
              <div className="text-sm text-onyx">Collections</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-eerie-black">
                {report.summary.totalColumns || 0}
              </div>
              <div className="text-sm text-onyx">Fields</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-eerie-black">
                {report.summary.totalRows?.toLocaleString() || 0}
              </div>
              <div className="text-sm text-onyx">Documents</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-eerie-black">
                {report.summary.totalSize || "N/A"}
              </div>
              <div className="text-sm text-onyx">Database Size</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-eerie-black">
                {report.summary.complexityScore?.toFixed(1) || "N/A"}
              </div>
              <div className="text-sm text-onyx">Complexity</div>
            </div>
          </div>
        )}

        {report.tables && report.tables.length > 0 && (
          <CollectionsDetails
            tables={report.tables}
            onCollectionClick={onCollectionClick}
            currentConnectionId={currentConnectionId}
          />
        )}

        {report.insights && report.insights.length > 0 && (
          <InsightsSection insights={report.insights} />
        )}

        {report.recommendations && report.recommendations.length > 0 && (
          <RecommendationsSection recommendations={report.recommendations} />
        )}

        {report.performanceMetrics && (
          <PerformanceMetrics performanceMetrics={report.performanceMetrics} />
        )}
      </CardContent>
    </Card>
  );
}
