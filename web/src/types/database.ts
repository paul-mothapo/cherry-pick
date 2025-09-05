export interface DatabaseConnection {
  id: string;
  name: string;
  driver: 'mysql' | 'postgresql' | 'sqlite' | 'mongodb';
  connectionString: string;
  status: 'connected' | 'disconnected' | 'error';
  lastConnected?: string;
}

export interface DatabaseReport {
  databaseName: string;
  analysisTime: string;
  summary: DatabaseSummary;
  tables: TableInfo[];
  insights: DatabaseInsight[];
  recommendations: string[];
  performanceMetrics: PerformanceMetrics;
}

export interface DatabaseSummary {
  healthScore: number;
  totalTables: number;
  totalColumns: number;
  totalRows: number;
  totalSize: string;
  complexityScore: number;
}

export interface TableInfo {
  name: string;
  schema?: string;
  rowCount: number;
  size: string;
  columns: ColumnInfo[];
  indexes: IndexInfo[];
  constraints: ConstraintInfo[];
  relationships: RelationshipInfo[];
  lastModified?: string;
}

export interface ColumnInfo {
  name: string;
  dataType: string;
  nullable: boolean;
  primaryKey: boolean;
  foreignKey: boolean;
  unique: boolean;
  defaultValue?: string;
  uniqueValues: number;
  dataProfile: DataProfile;
}

export interface IndexInfo {
  name: string;
  tableName?: string;
  columns: string[];
  isUnique: boolean;
  type: string;
}

export interface ConstraintInfo {
  name: string;
  tableName: string;
  type: string;
  columns: string[];
  referencedTable?: string;
  referencedColumns?: string[];
}

export interface SecurityInfo {
  issues: SecurityIssue[];
  score: number;
  recommendations: string[];
}

export interface SecurityIssue {
  id: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  category: string;
  description: string;
  affectedTables: string[];
  recommendation: string;
}

export interface QueryOptimization {
  originalQuery: string;
  optimizedQuery: string;
  explanation: string;
  estimatedImprovement: string;
  warnings: string[];
}

export interface Alert {
  id: string;
  type: 'performance' | 'security' | 'health' | 'capacity';
  severity: 'info' | 'warning' | 'error' | 'critical';
  title: string;
  description: string;
  timestamp: string;
  acknowledged: boolean;
}

export interface DataProfile {
  min?: any;
  max?: any;
  avg?: any;
  sampleData: string[];
  pattern: string;
  quality: number;
}

export interface RelationshipInfo {
  type: string;
  targetTable: string;
  sourceColumn: string;
  targetColumn: string;
}

export interface DatabaseInsight {
  type: string;
  severity: 'low' | 'medium' | 'high';
  title: string;
  description: string;
  suggestion: string;
  affectedTables: string[];
  metricValue?: any;
}

export interface PerformanceMetrics {
  slowQueries?: any[];
  indexUsage: number;
  tableScanRatio: number;
  connectionCount: number;
  bufferHitRatio: number;
}

export interface DataLineage {
  tables: LineageTable[];
  relationships: LineageRelationship[];
}

export interface LineageTable {
  name: string;
  type: 'source' | 'derived' | 'target';
  dependencies: string[];
}

export interface LineageRelationship {
  from: string;
  to: string;
  type: 'direct' | 'indirect';
  description: string;
}
