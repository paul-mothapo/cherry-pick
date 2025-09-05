import { MousePointer } from "lucide-react";

interface Table {
  name: string;
  rowCount?: number;
  columns?: Array<{
    name: string;
    dataType: string;
    uniqueValues?: number;
    dataProfile?: {
      quality: number;
    };
  }>;
  indexes?: Array<{
    name: string;
    columns?: string[];
    isUnique: boolean;
  }>;
  relationships?: Array<{
    sourceColumn: string;
    targetTable: string;
    targetColumn: string;
    type: string;
  }>;
  size?: string;
}

interface CollectionsDetailsProps {
  tables: Table[];
  onCollectionClick: (collection: any, connectionId: string) => void;
  currentConnectionId?: string;
}

export default function CollectionsDetails({
  tables,
  onCollectionClick,
  currentConnectionId,
}: CollectionsDetailsProps) {
  return (
    <div className="space-y-4">
      <h4 className="font-medium text-eerie-black text-lg">
        Collections Details:
      </h4>
      <div className="space-y-3">
        {tables.map((table, idx) => (
          <div key={idx} className="border rounded-lg p-4 bg-seasalt">
            <div className="flex justify-between items-start mb-3">
              <div>
                <button
                  onClick={() =>
                    onCollectionClick(table, currentConnectionId || "current")
                  }
                  className="group flex items-center hover:text-medium-blue transition-colors"
                >
                  <h5 className="font-semibold text-eerie-black group-hover:text-medium-blue">
                    {table.name}
                  </h5>
                  <MousePointer className="h-4 w-4 ml-2 opacity-0 group-hover:opacity-100 transition-opacity" />
                </button>
                <div className="text-sm text-onyx space-x-4">
                  <span>{table.rowCount?.toLocaleString() || 0} documents</span>
                  <span>{table.columns?.length || 0} fields</span>
                  <span>{table.indexes?.length || 0} indexes</span>
                  {table.size && <span>Size: {table.size}</span>}
                </div>
              </div>
            </div>

            {table.columns && table.columns.length > 0 && (
              <div className="mb-3">
                <h6 className="font-medium text-onyx mb-2">Fields:</h6>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
                  {table.columns.map((column, colIdx) => (
                    <div
                      key={colIdx}
                      className="bg-white p-2 rounded border text-sm"
                    >
                      <div className="font-medium text-eerie-black">
                        {column.name}
                      </div>
                      <div className="text-onyx">{column.dataType}</div>
                      {column.uniqueValues && (
                        <div className="text-xs text-onyx">
                          {column.uniqueValues} unique values
                        </div>
                      )}
                      {column.dataProfile?.quality && (
                        <div className="text-xs text-onyx">
                          Quality:{" "}
                          {(column.dataProfile.quality * 100).toFixed(1)}%
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {table.indexes && table.indexes.length > 0 && (
              <div className="mb-3">
                <h6 className="font-medium text-onyx mb-2">Indexes:</h6>
                <div className="flex flex-wrap gap-2">
                  {table.indexes.map((index, idxIdx) => (
                    <span
                      key={idxIdx}
                      className={`px-2 py-1 text-xs rounded ${
                        index.isUnique
                          ? "bg-green-100 text-green-800"
                          : "bg-platinum text-medium-blue"
                      }`}
                    >
                      {index.name} ({index.columns?.join(", ")})
                      {index.isUnique && " • Unique"}
                    </span>
                  ))}
                </div>
              </div>
            )}

            {table.relationships && table.relationships.length > 0 && (
              <div>
                <h6 className="font-medium text-onyx mb-2">Relationships:</h6>
                <div className="space-y-1">
                  {table.relationships.map((rel, relIdx) => (
                    <div key={relIdx} className="text-sm text-onyx">
                      {rel.sourceColumn} → {rel.targetTable}.{rel.targetColumn}{" "}
                      ({rel.type})
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
