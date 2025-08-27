package types

import "time"

type MongoCollectionInfo struct {
	Name           string                 `json:"name"`
	DocumentCount  int64                  `json:"document_count"`
	AvgDocSize     int64                  `json:"avg_doc_size"`
	TotalSize      int64                  `json:"total_size"`
	StorageSize    int64                  `json:"storage_size"`
	Indexes        []MongoIndexInfo       `json:"indexes"`
	Fields         []MongoFieldInfo       `json:"fields"`
	SampleDocument map[string]interface{} `json:"sample_document"`
	LastModified   time.Time              `json:"last_modified"`
	ShardKey       string                 `json:"shard_key,omitempty"`
	IsSharded      bool                   `json:"is_sharded"`
}

type MongoFieldInfo struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Frequency   float64     `json:"frequency"`
	SampleValue interface{} `json:"sample_value"`
	IsRequired  bool        `json:"is_required"`
	IsIndexed   bool        `json:"is_indexed"`
}

type MongoIndexInfo struct {
	Name       string                 `json:"name"`
	Keys       map[string]interface{} `json:"keys"`
	IsUnique   bool                   `json:"is_unique"`
	IsSparse   bool                   `json:"is_sparse"`
	IsPartial  bool                   `json:"is_partial"`
	Size       int64                  `json:"size"`
	UsageStats MongoIndexUsageStats   `json:"usage_stats"`
}

type MongoIndexUsageStats struct {
	Ops   int64     `json:"ops"`
	Since time.Time `json:"since"`
}

type MongoDatabaseStats struct {
	Name           string  `json:"name"`
	Collections    int     `json:"collections"`
	Views          int     `json:"views"`
	Objects        int64   `json:"objects"`
	AvgObjSize     float64 `json:"avg_obj_size"`
	DataSize       int64   `json:"data_size"`
	StorageSize    int64   `json:"storage_size"`
	IndexSize      int64   `json:"index_size"`
	TotalSize      int64   `json:"total_size"`
	ScaleFactor    int64   `json:"scale_factor"`
	FsUsedSize     int64   `json:"fs_used_size"`
	FsTotalSize    int64   `json:"fs_total_size"`
	IndexFreeBytes int64   `json:"index_free_bytes"`
	TotalFreeBytes int64   `json:"total_free_bytes"`
}

type MongoPerformanceMetrics struct {
	Connections MongoConnectionStats `json:"connections"`
	Operations  MongoOpStats         `json:"operations"`
	Memory      MongoMemoryStats     `json:"memory"`
	Replication MongoReplStats       `json:"replication"`
	Sharding    MongoShardStats      `json:"sharding"`
	WiredTiger  MongoWTStats         `json:"wired_tiger"`
}

type MongoConnectionStats struct {
	Current      int `json:"current"`
	Available    int `json:"available"`
	TotalCreated int `json:"total_created"`
}

type MongoOpStats struct {
	Insert  int64 `json:"insert"`
	Query   int64 `json:"query"`
	Update  int64 `json:"update"`
	Delete  int64 `json:"delete"`
	GetMore int64 `json:"getmore"`
	Command int64 `json:"command"`
}

type MongoMemoryStats struct {
	Bits              int   `json:"bits"`
	Resident          int64 `json:"resident"`
	Virtual           int64 `json:"virtual"`
	Supported         bool  `json:"supported"`
	Mapped            int64 `json:"mapped"`
	MappedWithJournal int64 `json:"mapped_with_journal"`
}

type MongoReplStats struct {
	SetName    string              `json:"set_name"`
	IsMaster   bool                `json:"is_master"`
	Secondary  bool                `json:"secondary"`
	Hosts      []string            `json:"hosts"`
	Primary    string              `json:"primary"`
	Me         string              `json:"me"`
	ElectionId interface{}         `json:"election_id"`
	LastWrite  MongoLastWriteStats `json:"last_write"`
}

type MongoLastWriteStats struct {
	OpTime            MongoOpTime `json:"op_time"`
	LastWriteDate     time.Time   `json:"last_write_date"`
	MajorityOpTime    MongoOpTime `json:"majority_op_time"`
	MajorityWriteDate time.Time   `json:"majority_write_date"`
}

type MongoOpTime struct {
	Ts time.Time `json:"ts"`
	T  int64     `json:"t"`
}

type MongoShardStats struct {
	IsSharded bool     `json:"is_sharded"`
	Shards    []string `json:"shards"`
	Chunks    int      `json:"chunks"`
}

type MongoWTStats struct {
	BlockManager MongoWTBlockManager `json:"block_manager"`
	Cache        MongoWTCache        `json:"cache"`
	Connection   MongoWTConnection   `json:"connection"`
}

type MongoWTBlockManager struct {
	BlocksRead    int64 `json:"blocks_read"`
	BlocksWritten int64 `json:"blocks_written"`
	BytesRead     int64 `json:"bytes_read"`
	BytesWritten  int64 `json:"bytes_written"`
}

type MongoWTCache struct {
	BytesCurrentlyInCache  int64   `json:"bytes_currently_in_cache"`
	BytesReadIntoCache     int64   `json:"bytes_read_into_cache"`
	BytesWrittenFromCache  int64   `json:"bytes_written_from_cache"`
	MaximumBytesConfigured int64   `json:"maximum_bytes_configured"`
	PercentOverhead        float64 `json:"percent_overhead"`
	TrackedDirtyBytes      int64   `json:"tracked_dirty_bytes"`
	PagesEvicted           int64   `json:"pages_evicted"`
	PagesReadIntoCache     int64   `json:"pages_read_into_cache"`
	PagesWrittenFromCache  int64   `json:"pages_written_from_cache"`
}

type MongoWTConnection struct {
	DataHandleConnections int64 `json:"data_handle_connections"`
	FilesCurrentlyOpen    int64 `json:"files_currently_open"`
	TotalReadIOs          int64 `json:"total_read_ios"`
	TotalWriteIOs         int64 `json:"total_write_ios"`
}
