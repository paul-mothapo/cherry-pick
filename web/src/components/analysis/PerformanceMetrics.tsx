interface PerformanceMetricsData {
  indexUsage?: number;
  tableScanRatio?: number;
  connectionCount?: number;
  bufferHitRatio?: number;
}

interface PerformanceMetricsProps {
  performanceMetrics: PerformanceMetricsData;
}

export default function PerformanceMetrics({
  performanceMetrics,
}: PerformanceMetricsProps) {
  return (
    <div className="mt-6">
      <h4 className="font-medium text-eerie-black text-lg mb-3">
        Performance Metrics:
      </h4>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {performanceMetrics.indexUsage && (
          <div className="text-center p-3 bg-seasalt rounded">
            <div className="text-lg font-semibold text-eerie-black">
              {(performanceMetrics.indexUsage * 100).toFixed(1)}%
            </div>
            <div className="text-sm text-onyx">Index Usage</div>
          </div>
        )}
        {performanceMetrics.tableScanRatio && (
          <div className="text-center p-3 bg-seasalt rounded">
            <div className="text-lg font-semibold text-gray-900">
              {(performanceMetrics.tableScanRatio * 100).toFixed(1)}%
            </div>
            <div className="text-sm text-onyx">Table Scans</div>
          </div>
        )}
        {performanceMetrics.connectionCount && (
          <div className="text-center p-3 bg-seasalt rounded">
            <div className="text-lg font-semibold text-gray-900">
              {performanceMetrics.connectionCount}
            </div>
            <div className="text-sm text-onyx">Connections</div>
          </div>
        )}
        {performanceMetrics.bufferHitRatio && (
          <div className="text-center p-3 bg-seasalt rounded">
            <div className="text-lg font-semibold text-gray-900">
              {(performanceMetrics.bufferHitRatio * 100).toFixed(1)}%
            </div>
            <div className="text-sm text-onyx">Buffer Hit Ratio</div>
          </div>
        )}
      </div>
    </div>
  );
}
