import React from "react";
import AnalysisControls from "@/components/analysis/AnalysisControls";
import AnalyticsReports from "@/components/analysis/AnalyticsReports";

export const Analysis: React.FC = () => {

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Database Analysis</h1>
        <p className="text-onyx mt-2">
          Analyze your database structure, performance, and get insights
        </p>
      </div>

      <AnalysisControls />

      <AnalyticsReports />
    </div>
  );
};
