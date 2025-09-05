import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Search, Filter, Download, BarChart3, Eye, Database } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { useQuery } from '@tanstack/react-query';
import { collectionApi } from '@/services/api';
import { useConnections } from '@/hooks/useConnections';

export const CollectionView: React.FC = () => {
  const { connectionId, collectionName } = useParams<{ connectionId: string; collectionName: string }>();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'data' | 'schema' | 'analytics'>('data');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedField, setSelectedField] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  const { data: connections } = useConnections();
  const currentConnection = connections?.find(c => c.id === connectionId);

  const { data: collectionData, isLoading: isLoadingData, refetch: refetchData } = useQuery({
    queryKey: ['collectionData', connectionId, collectionName, currentPage, pageSize],
    queryFn: async () => {
      const response = await collectionApi.getCollectionData(connectionId!, collectionName!, currentPage, pageSize);
      return response.data.data;
    },
    enabled: activeTab === 'data' && !!connectionId && !!collectionName,
  });

  const { data: collectionStats, isLoading: isLoadingStats } = useQuery({
    queryKey: ['collectionStats', connectionId, collectionName],
    queryFn: async () => {
      const response = await collectionApi.getCollectionStats(connectionId!, collectionName!);
      return response.data.data;
    },
    enabled: activeTab === 'analytics' && !!connectionId && !!collectionName,
  });

  const formatValue = (value: any): string => {
    if (value === null || value === undefined) return 'null';
    if (typeof value === 'object') return JSON.stringify(value, null, 2);
    if (typeof value === 'string' && value.length > 100) return value.substring(0, 100) + '...';
    return String(value);
  };

  const getFieldType = (value: any): string => {
    if (value === null || value === undefined) return 'null';
    if (Array.isArray(value)) return 'array';
    if (typeof value === 'object') return 'object';
    if (typeof value === 'string' && /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/.test(value)) return 'date';
    return typeof value;
  };

  const fieldOptions = collectionStats?.fields?.map(field => ({
    value: field.name,
    label: field.name
  })) || [];

  const tabs = [
    { id: 'data', label: 'Data', icon: Database },
    { id: 'schema', label: 'Schema', icon: Eye },
    { id: 'analytics', label: 'Analytics', icon: BarChart3 }
  ];

  if (!connectionId || !collectionName) {
    return (
      <div className="p-6">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Invalid Collection</h1>
          <p className="text-onyx mb-4">Collection or connection not found.</p>
          <Button onClick={() => navigate('/analysis')}>
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Analysis
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white border-b border-gray-200">
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button
                variant="ghost"
                onClick={() => navigate('/analysis')}
                className="flex items-center"
              >
                <ArrowLeft className="h-5 w-5 mr-2" />
                Back to Analysis
              </Button>
              <div className="border-l border-gray-300 pl-4">
                <h1 className="text-2xl font-bold text-gray-900">{collectionName}</h1>
                <div className="flex items-center space-x-4 text-sm text-onyx mt-1">
                  <span>{currentConnection?.name || 'Unknown Connection'}</span>
                  <span>•</span>
                  <span>{currentConnection?.driver?.toUpperCase() || 'MongoDB'}</span>
                  {collectionStats && (
                    <>
                      <span>•</span>
                      <span>{collectionStats.documentCount?.toLocaleString() || 0} documents</span>
                      <span>•</span>
                      <span>{collectionStats.fields?.length || 0} fields</span>
                    </>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="flex border-b">
          {tabs.map(tab => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`flex items-center px-6 py-3 font-medium text-sm border-b-2 transition-colors ${
                activeTab === tab.id
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              <tab.icon className="h-4 w-4 mr-2" />
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      <div className="flex-1">
        {activeTab === 'data' && (
          <div className="h-full flex flex-col">
            {/* Data Controls */}
            <div className="bg-white border-b border-gray-200 px-6 py-4">
              <div className="flex items-center space-x-4">
                <div className="flex-1">
                  <Input
                    placeholder="Search documents..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="w-full max-w-md"
                  />
                </div>
                <Select
                  options={fieldOptions}
                  value={selectedField}
                  onChange={(e) => setSelectedField(e.target.value)}
                  placeholder="Filter by field"
                />
                <Button variant="secondary" size="sm">
                  <Filter className="h-4 w-4 mr-2" />
                  Filter
                </Button>
                <Button variant="secondary" size="sm">
                  <Download className="h-4 w-4 mr-2" />
                  Export
                </Button>
              </div>
            </div>

            {/* Data Content */}
            <div className="flex-1 p-6">
              {isLoadingData ? (
                <div className="flex items-center justify-center h-64">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
                </div>
              ) : collectionData?.documents && collectionData.documents.length > 0 ? (
                <div className="bg-white rounded-lg border">
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead className="bg-gray-50">
                        <tr>
                          {Object.keys(collectionData.documents[0]).map(key => (
                            <th key={key} className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                              {key}
                            </th>
                          ))}
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-gray-200">
                        {collectionData.documents.map((doc, idx) => (
                          <tr key={idx} className="hover:bg-gray-50">
                            {Object.entries(doc).map(([key, value]) => (
                              <td key={key} className="px-4 py-3 text-sm text-gray-900 max-w-xs">
                                <div className="truncate" title={formatValue(value)}>
                                  {formatValue(value)}
                                </div>
                              </td>
                            ))}
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>

                  {/* Pagination */}
                  <div className="px-6 py-4 border-t border-gray-200">
                    <div className="flex items-center justify-between">
                      <div className="text-sm text-onyx">
                        Showing {((currentPage - 1) * pageSize) + 1} to {Math.min(currentPage * pageSize, collectionData.totalCount)} of {collectionData.totalCount} documents
                      </div>
                      <div className="flex items-center space-x-2">
                        <Button
                          variant="secondary"
                          size="sm"
                          disabled={currentPage === 1}
                          onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                        >
                          Previous
                        </Button>
                        <span className="text-sm text-onyx">
                          Page {currentPage} of {Math.ceil(collectionData.totalCount / pageSize)}
                        </span>
                        <Button
                          variant="secondary"
                          size="sm"
                          disabled={currentPage >= Math.ceil(collectionData.totalCount / pageSize)}
                          onClick={() => setCurrentPage(p => p + 1)}
                        >
                          Next
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="bg-white rounded-lg border p-12 text-center">
                  <Database className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                  <h3 className="text-lg font-medium text-gray-900 mb-2">No documents found</h3>
                  <p className="text-onyx">This collection appears to be empty.</p>
                </div>
              )}
            </div>
          </div>
        )}

        {activeTab === 'schema' && (
          <div className="p-6">
            <div className="space-y-6">
              <div>
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Field Schema</h3>
                <div className="grid gap-4">
                  {collectionStats?.fields?.map((field, idx) => (
                    <Card key={idx}>
                      <CardContent className="p-4">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <h4 className="font-medium text-gray-900">{field.name}</h4>
                            <p className="text-sm text-onyx mt-1">{field.type}</p>
                            <div className="flex items-center space-x-4 mt-2 text-xs text-gray-500">
                              <span>{field.count} values</span>
                              <span>{field.uniqueCount} unique</span>
                              {field.nullCount > 0 && <span>{field.nullCount} nulls</span>}
                            </div>
                          </div>
                          <div className="text-right">
                            <div className="text-sm text-onyx">
                              {field.nullCount === 0 ? 'Required' : 'Optional'}
                            </div>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  )) || (
                    <p className="text-gray-500">No schema information available</p>
                  )}
                </div>
              </div>

              {/* Indexes */}
              {collectionStats?.indexes && collectionStats.indexes.length > 0 && (
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-4">Indexes</h3>
                  <div className="space-y-2">
                    {collectionStats.indexes.map((index, idx) => (
                      <div key={idx} className="flex items-center justify-between p-3 bg-white rounded border">
                        <div>
                          <span className="font-medium">{index.name}</span>
                          <span className="text-onyx ml-2">({index.keys.join(', ')})</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          {index.isUnique && (
                            <span className="bg-green-100 text-green-800 text-xs px-2 py-1 rounded">
                              Unique
                            </span>
                          )}
                          {index.size && (
                            <span className="text-sm text-onyx">
                              {(index.size / 1024).toFixed(1)} KB
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {activeTab === 'analytics' && (
          <div className="p-6">
            {isLoadingStats ? (
              <div className="flex items-center justify-center h-64">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              </div>
            ) : (
              <div className="space-y-6">
                {/* Collection Overview */}
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-4">Collection Overview</h3>
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <Card>
                      <CardContent className="p-4 text-center">
                        <div className="text-2xl font-bold text-green-600">
                          {collectionStats?.documentCount?.toLocaleString() || 0}
                        </div>
                        <div className="text-sm text-onyx">Documents</div>
                      </CardContent>
                    </Card>
                    <Card>
                      <CardContent className="p-4 text-center">
                        <div className="text-2xl font-bold text-blue-600">
                          {collectionStats?.fields?.length || 0}
                        </div>
                        <div className="text-sm text-onyx">Fields</div>
                      </CardContent>
                    </Card>
                    <Card>
                      <CardContent className="p-4 text-center">
                        <div className="text-2xl font-bold text-purple-600">
                          {collectionStats?.indexes?.length || 0}
                        </div>
                        <div className="text-sm text-onyx">Indexes</div>
                      </CardContent>
                    </Card>
                    <Card>
                      <CardContent className="p-4 text-center">
                        <div className="text-2xl font-bold text-orange-600">
                          {collectionStats?.fields ? 
                            ((collectionStats.fields.filter(f => f.nullCount === 0).length / collectionStats.fields.length) * 100).toFixed(1)
                            : 0}%
                        </div>
                        <div className="text-sm text-onyx">Completeness</div>
                      </CardContent>
                    </Card>
                  </div>
                </div>

                {/* Field Analytics */}
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 mb-4">Field Analytics</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {collectionStats?.fields?.map((field, idx) => (
                      <Card key={idx}>
                        <CardContent className="p-4">
                          <h4 className="font-medium text-gray-900 mb-2">{field.name}</h4>
                          <div className="space-y-2 text-sm">
                            <div className="flex justify-between">
                              <span className="text-onyx">Type:</span>
                              <span className="font-medium">{field.type}</span>
                            </div>
                            <div className="flex justify-between">
                              <span className="text-onyx">Count:</span>
                              <span className="font-medium">{field.count}</span>
                            </div>
                            <div className="flex justify-between">
                              <span className="text-onyx">Unique:</span>
                              <span className="font-medium">{field.uniqueCount}</span>
                            </div>
                            <div className="flex justify-between">
                              <span className="text-onyx">Nulls:</span>
                              <span className="font-medium">{field.nullCount}</span>
                            </div>
                            {field.sampleValues && field.sampleValues.length > 0 && (
                              <div className="mt-2">
                                <span className="text-xs text-gray-500">Sample values:</span>
                                <div className="flex flex-wrap gap-1 mt-1">
                                  {field.sampleValues.slice(0, 3).map((sample, sIdx) => (
                                    <span key={sIdx} className="bg-gray-100 text-gray-700 px-2 py-1 text-xs rounded">
                                      {sample.length > 15 ? sample.substring(0, 15) + '...' : sample}
                                    </span>
                                  ))}
                                </div>
                              </div>
                            )}
                          </div>
                        </CardContent>
                      </Card>
                    )) || (
                      <p className="text-gray-500">No analytics data available</p>
                    )}
                  </div>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};
