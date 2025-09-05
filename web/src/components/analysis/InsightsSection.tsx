interface Insight {
  title: string;
  description: string;
  suggestion?: string;
  affectedTables?: string[];
  severity: "high" | "medium" | "low";
}

interface InsightsSectionProps {
  insights: Insight[];
}

export default function InsightsSection({ insights }: InsightsSectionProps) {
  return (
    <div className="mt-6 space-y-4">
      <h4 className="font-medium text-eerie-black text-lg">Analysis Insights:</h4>
      <div className="space-y-3">
        {insights.map((insight, idx) => (
          <div
            key={idx}
            className={`p-4 rounded-lg border-l-4 ${
              insight.severity === "high"
                ? "bg-red-50 border-red-400"
                : insight.severity === "medium"
                ? "bg-yellow-50 border-yellow-400"
                : "bg-blue-50 border-blue-400"
            }`}
          >
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <h5 className="font-semibold text-eerie-black">{insight.title}</h5>
                <p className="text-onyx mt-1">{insight.description}</p>
                {insight.suggestion && (
                  <p className="text-onyx mt-2 text-sm">
                    <strong>Suggestion:</strong> {insight.suggestion}
                  </p>
                )}
                {insight.affectedTables &&
                  insight.affectedTables.length > 0 && (
                    <div className="mt-2">
                      <span className="text-sm text-onyx">Affects: </span>
                      {insight.affectedTables.map((table, tIdx) => (
                        <span
                          key={tIdx}
                          className="inline-block bg-seasalt text-onyx px-2 py-1 text-xs rounded mr-1"
                        >
                          {table}
                        </span>
                      ))}
                    </div>
                  )}
              </div>
              <div className="ml-4">
                <span
                  className={`px-2 py-1 text-xs font-medium rounded ${
                    insight.severity === "high"
                      ? "bg-red-100 text-turkey-red"
                      : insight.severity === "medium"
                      ? "bg-yellow-100 text-yellow-800"
                      : "bg-blue-100 text-medium-blue"
                  }`}
                >
                  {insight.severity?.toUpperCase()}
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
