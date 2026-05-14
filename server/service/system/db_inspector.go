package system

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DBInspectorService struct{}

var dbInspectorIdentifierPattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

const (
	defaultInspectorPageSize = 20
	maxInspectorPageSize     = 200
	defaultDeleteMaxRows     = 5000
	maxDeleteRowsLimit       = 50000
)

var systemRequiredTables = map[string]struct{}{
	"sys_users":              {},
	"sys_authorities":        {},
	"sys_authority_menus":    {},
	"sys_base_menus":         {},
	"sys_apis":               {},
	"sys_params":             {},
	"casbin_rule":            {},
	"jwt_blacklists":         {},
	"sys_user_authorities":   {},
	"sys_base_menu_btns":     {},
	"sys_base_menu_btn_pars": {},
}

type DBInspectorOverview struct {
	Database  DBRuntimeStatus    `json:"database"`
	Redis     RedisRuntimeStatus `json:"redis"`
	TableList []DBTableProfile   `json:"tableList"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"pageSize"`
}

type DBRuntimeStatus struct {
	DBType             string   `json:"dbType"`
	DatabaseName       string   `json:"databaseName"`
	Healthy            bool     `json:"healthy"`
	Reason             string   `json:"reason"`
	Warning            string   `json:"warning"`
	TotalRows          int64    `json:"totalRows"`
	TotalSizeBytes     int64    `json:"totalSizeBytes"`
	TotalSizeHuman     string   `json:"totalSizeHuman"`
	OpenConnections    int      `json:"openConnections"`
	InUseConnections   int      `json:"inUseConnections"`
	IdleConnections    int      `json:"idleConnections"`
	MaxOpenConnections int      `json:"maxOpenConnections"`
	WaitCount          int64    `json:"waitCount"`
	AutoFixSupported   bool     `json:"autoFixSupported"`
	ManualFixGuide     []string `json:"manualFixGuide"`
}

type RedisRuntimeStatus struct {
	Enabled          bool                `json:"enabled"`
	Healthy          bool                `json:"healthy"`
	Reason           string              `json:"reason"`
	Version          string              `json:"version"`
	UsedMemory       string              `json:"usedMemory"`
	Keyspaces        []RedisKeyspaceStat `json:"keyspaces"`
	AutoFixSupported bool                `json:"autoFixSupported"`
	ManualFixGuide   []string            `json:"manualFixGuide"`
}

type RedisKeyspaceStat struct {
	DB      string `json:"db"`
	Keys    int64  `json:"keys"`
	Expires int64  `json:"expires"`
	AvgTTL  int64  `json:"avgTTL"`
}

type DBTableProfile struct {
	TableName         string   `json:"tableName"`
	RowCount          int64    `json:"rowCount"`
	SizeBytes         int64    `json:"sizeBytes"`
	SizeHuman         string   `json:"sizeHuman"`
	Purpose           string   `json:"purpose"`
	UsageLabel        string   `json:"usageLabel"`
	SystemRequired    bool     `json:"systemRequired"`
	LikelyUnused      bool     `json:"likelyUnused"`
	DateColumns       []string `json:"dateColumns"`
	FileColumns       []string `json:"fileColumns"`
	CanDeleteByDate   bool     `json:"canDeleteByDate"`
	DeleteGuardReason string   `json:"deleteGuardReason"`
}

type DBInspectorAutoFixResult struct {
	Target  string `json:"target"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type DBInspectorDeleteResult struct {
	TableName         string   `json:"tableName"`
	TimeColumn        string   `json:"timeColumn"`
	MatchedRows       int64    `json:"matchedRows"`
	DeletedRows       int64    `json:"deletedRows"`
	DeletedFileCount  int64    `json:"deletedFileCount"`
	FailedFileCount   int64    `json:"failedFileCount"`
	FileDeleteWarning []string `json:"fileDeleteWarning"`
}

type tableMetricRow struct {
	TableName  string
	TableRows  sql.NullInt64
	TotalBytes sql.NullInt64
}

func (s *DBInspectorService) GetOverview(ctx context.Context, req systemReq.DBInspectorOverviewSearch) (overview DBInspectorOverview, err error) {
	if global.GVA_DB == nil {
		return overview, errors.New("数据库尚未初始化")
	}

	page, pageSize := normalizeInspectorPage(req.Page, req.PageSize)
	profiles, err := s.buildTableProfiles(ctx)
	if err != nil {
		return overview, err
	}

	keyword := strings.ToLower(strings.TrimSpace(req.Keyword))
	if keyword != "" {
		filtered := make([]DBTableProfile, 0, len(profiles))
		for _, item := range profiles {
			if strings.Contains(strings.ToLower(item.TableName), keyword) || strings.Contains(strings.ToLower(item.Purpose), keyword) {
				filtered = append(filtered, item)
			}
		}
		profiles = filtered
	}

	total := int64(len(profiles))
	start := (page - 1) * pageSize
	if start > len(profiles) {
		start = len(profiles)
	}
	end := start + pageSize
	if end > len(profiles) {
		end = len(profiles)
	}

	databaseStatus := s.getDatabaseRuntimeStatus(ctx)
	databaseStatus.DatabaseName = s.currentDatabaseName()
	databaseStatus.TotalSizeBytes, databaseStatus.TotalRows = sumProfileStats(profiles)
	databaseStatus.TotalSizeHuman = formatBytes(databaseStatus.TotalSizeBytes)

	overview = DBInspectorOverview{
		Database:  databaseStatus,
		Redis:     s.getRedisRuntimeStatus(ctx),
		TableList: profiles[start:end],
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}
	return overview, nil
}

func (s *DBInspectorService) AutoFix(ctx context.Context, req systemReq.DBInspectorAutoFixReq) (result DBInspectorAutoFixResult, err error) {
	target := strings.ToLower(strings.TrimSpace(req.Target))
	result.Target = target

	switch target {
	case "database":
		if req.TableName != "" {
			if err = s.tryOptimizeTable(ctx, req.TableName); err != nil {
				global.GVA_LOG.Warn("数据库巡检: 优化表失败", zap.String("table", req.TableName), zap.Error(err))
			}
		}
		err = s.reconnectDatabase(ctx)
		if err != nil {
			result.Success = false
			result.Message = "数据库自动修复失败"
			result.Detail = err.Error()
			return result, nil
		}
		status := s.getDatabaseRuntimeStatus(ctx)
		result.Success = status.Healthy
		if status.Healthy {
			result.Message = "数据库自动修复完成，连接已恢复"
		} else {
			result.Message = "数据库自动修复未成功"
		}
		result.Detail = status.Reason
		return result, nil
	case "redis":
		err = s.reconnectRedis(ctx)
		if err != nil {
			result.Success = false
			result.Message = "Redis自动修复失败"
			result.Detail = err.Error()
			return result, nil
		}
		status := s.getRedisRuntimeStatus(ctx)
		result.Success = status.Healthy
		if status.Healthy {
			result.Message = "Redis自动修复完成，连接已恢复"
		} else {
			result.Message = "Redis自动修复未成功"
		}
		result.Detail = status.Reason
		return result, nil
	default:
		return result, errors.New("target 参数无效，仅支持 database 或 redis")
	}
}

func (s *DBInspectorService) DeleteRecordsByRange(ctx context.Context, req systemReq.DBInspectorDeleteByRangeReq) (result DBInspectorDeleteResult, err error) {
	if global.GVA_DB == nil {
		return result, errors.New("数据库尚未初始化")
	}

	tableName := strings.TrimSpace(req.TableName)
	if !isValidIdentifier(tableName) {
		return result, errors.New("表名不合法")
	}
	if !global.GVA_DB.Migrator().HasTable(tableName) {
		return result, errors.New("目标表不存在")
	}

	columns, err := s.getTableColumns(tableName)
	if err != nil {
		return result, err
	}

	dateColumn := strings.TrimSpace(req.TimeColumn)
	if dateColumn == "" {
		dateColumn = pickDateColumn(columns)
	}
	if dateColumn == "" {
		return result, errors.New("该表未发现可用于日期范围清理的时间字段")
	}
	if !containsColumn(columns, dateColumn) {
		return result, errors.New("timeColumn 不存在")
	}

	confirmExpected := fmt.Sprintf("DELETE %s", tableName)
	if strings.TrimSpace(req.ConfirmText) != confirmExpected {
		return result, fmt.Errorf("确认文本错误，请输入 %s", confirmExpected)
	}

	class := classifyTable(tableName, 0)
	if class.SystemRequired {
		return result, errors.New("该表属于系统关键表，不允许在此工具中执行真删除")
	}

	startTime, endTime, err := parseTimeRange(req.StartTime, req.EndTime)
	if err != nil {
		return result, err
	}

	countSQL := fmt.Sprintf(
		"SELECT COUNT(1) FROM %s WHERE %s >= ? AND %s <= ?",
		s.quoteIdentifier(tableName),
		s.quoteIdentifier(dateColumn),
		s.quoteIdentifier(dateColumn),
	)
	var matchedRows int64
	if err = global.GVA_DB.WithContext(ctx).Raw(countSQL, startTime, endTime).Scan(&matchedRows).Error; err != nil {
		return result, err
	}

	if matchedRows == 0 {
		return DBInspectorDeleteResult{
			TableName:   tableName,
			TimeColumn:  dateColumn,
			MatchedRows: 0,
			DeletedRows: 0,
		}, nil
	}

	maxDeleteRows := req.MaxDeleteRows
	if maxDeleteRows <= 0 {
		maxDeleteRows = defaultDeleteMaxRows
	}
	if maxDeleteRows > maxDeleteRowsLimit {
		maxDeleteRows = maxDeleteRowsLimit
	}

	fileKeys := []string{}
	warnings := []string{}
	deleteRowsLimit := matchedRows > maxDeleteRows
	fileColumns := detectFileColumns(columns)
	var deletedRows int64

	if deleteRowsLimit {
		pkColumn := s.getTablePrimaryKeyColumn(tableName, columns)
		if pkColumn == "" {
			return result, fmt.Errorf("命中 %d 条，超过单次清理上限 %d，但未识别到主键字段，无法按最早时间定量删除", matchedRows, maxDeleteRows)
		}

		selectedIDs, selectErr := s.listOldestRowIDsByRange(ctx, tableName, pkColumn, dateColumn, startTime, endTime, maxDeleteRows)
		if selectErr != nil {
			return result, selectErr
		}
		if len(selectedIDs) == 0 {
			return DBInspectorDeleteResult{
				TableName:   tableName,
				TimeColumn:  dateColumn,
				MatchedRows: matchedRows,
				DeletedRows: 0,
			}, nil
		}

		if req.DeleteRelatedFiles && len(fileColumns) > 0 {
			fileKeys, warnings, err = s.collectFileKeysByIDs(ctx, tableName, pkColumn, selectedIDs, fileColumns)
			if err != nil {
				return result, err
			}
		}

		deletedRows, err = s.deleteRowsByIDs(ctx, tableName, pkColumn, selectedIDs)
		if err != nil {
			return result, err
		}
		warnings = append(warnings, fmt.Sprintf("命中 %d 条超过上限 %d，已按最早时间优先删除前 %d 条", matchedRows, maxDeleteRows, deletedRows))
	} else {
		if req.DeleteRelatedFiles && len(fileColumns) > 0 {
			fileKeys, warnings, err = s.collectFileKeysByRange(ctx, tableName, dateColumn, startTime, endTime, fileColumns)
			if err != nil {
				return result, err
			}
		}

		deleteSQL := fmt.Sprintf(
			"DELETE FROM %s WHERE %s >= ? AND %s <= ?",
			s.quoteIdentifier(tableName),
			s.quoteIdentifier(dateColumn),
			s.quoteIdentifier(dateColumn),
		)
		deleteResult := global.GVA_DB.WithContext(ctx).Exec(deleteSQL, startTime, endTime)
		if deleteResult.Error != nil {
			return result, deleteResult.Error
		}
		deletedRows = deleteResult.RowsAffected
	}

	if deletedRows > 0 {
		if refreshErr := s.refreshTableStatsAfterDelete(ctx, tableName); refreshErr != nil {
			warnings = append(warnings, "已删除记录，但刷新MySQL统计失败: "+refreshErr.Error())
		}
		if s.normalizedDBType() == "mysql" {
			warnings = append(warnings, "MySQL InnoDB 删除后磁盘占用可能不会立即下降，如需回收空间可执行 OPTIMIZE TABLE。")
		}
	}

	deletedFileCount, failedFileCount, fileWarn := s.deleteFilesByKeys(fileKeys)
	warnings = append(warnings, fileWarn...)

	result = DBInspectorDeleteResult{
		TableName:         tableName,
		TimeColumn:        dateColumn,
		MatchedRows:       matchedRows,
		DeletedRows:       deletedRows,
		DeletedFileCount:  deletedFileCount,
		FailedFileCount:   failedFileCount,
		FileDeleteWarning: warnings,
	}
	return result, nil
}

func (s *DBInspectorService) buildTableProfiles(ctx context.Context) ([]DBTableProfile, error) {
	tableNames, err := global.GVA_DB.Migrator().GetTables()
	if err != nil {
		return nil, err
	}
	rowMap, sizeMap := s.getTableMetricMaps(ctx)

	profiles := make([]DBTableProfile, 0, len(tableNames))
	for _, tableName := range tableNames {
		columns, columnErr := s.getTableColumns(tableName)
		if columnErr != nil {
			global.GVA_LOG.Warn("数据库巡检: 读取表字段失败", zap.String("table", tableName), zap.Error(columnErr))
			columns = nil
		}

		dateColumns := detectDateColumns(columns)
		fileColumns := detectFileColumns(columns)

		rowCount, rowErr := s.getTableRowCount(ctx, tableName)
		if rowErr != nil {
			if approx, ok := rowMap[tableName]; ok {
				rowCount = approx
			} else {
				rowCount = 0
			}
			global.GVA_LOG.Warn("数据库巡检: 获取精确行数失败，已回退估算值", zap.String("table", tableName), zap.Error(rowErr))
		}

		sizeBytes := sizeMap[tableName]
		class := classifyTable(tableName, rowCount)

		profile := DBTableProfile{
			TableName:       tableName,
			RowCount:        rowCount,
			SizeBytes:       sizeBytes,
			SizeHuman:       formatBytes(sizeBytes),
			Purpose:         class.Purpose,
			UsageLabel:      class.UsageLabel,
			SystemRequired:  class.SystemRequired,
			LikelyUnused:    class.LikelyUnused,
			DateColumns:     dateColumns,
			FileColumns:     fileColumns,
			CanDeleteByDate: !class.SystemRequired && len(dateColumns) > 0,
		}
		if class.SystemRequired {
			profile.DeleteGuardReason = "系统关键表已启用保护，不允许真删除"
		} else if len(dateColumns) == 0 {
			profile.DeleteGuardReason = "无可识别的日期字段，无法按日期范围清理"
		}
		if rowCount == 0 && sizeBytes >= 50*1024*1024 {
			reason := "表记录数为0但占用空间较大，通常是MySQL InnoDB历史占用未回收，可执行 OPTIMIZE TABLE 回收空间"
			if profile.DeleteGuardReason == "" {
				profile.DeleteGuardReason = reason
			} else {
				profile.DeleteGuardReason = profile.DeleteGuardReason + "；" + reason
			}
		}
		profiles = append(profiles, profile)
	}

	sort.Slice(profiles, func(i, j int) bool {
		if profiles[i].SizeBytes == profiles[j].SizeBytes {
			if profiles[i].RowCount == profiles[j].RowCount {
				return profiles[i].TableName < profiles[j].TableName
			}
			return profiles[i].RowCount > profiles[j].RowCount
		}
		return profiles[i].SizeBytes > profiles[j].SizeBytes
	})

	return profiles, nil
}

func (s *DBInspectorService) getDatabaseRuntimeStatus(ctx context.Context) DBRuntimeStatus {
	status := DBRuntimeStatus{
		DBType:           global.GVA_CONFIG.System.DbType,
		Healthy:          false,
		Reason:           "数据库连接未建立",
		AutoFixSupported: true,
		ManualFixGuide:   dbManualFixGuide(),
	}

	if global.GVA_DB == nil {
		return status
	}

	sqlDB, err := global.GVA_DB.DB()
	if err != nil {
		status.Reason = "数据库连接池不可用: " + err.Error()
		return status
	}

	if err = sqlDB.PingContext(ctx); err != nil {
		status.Reason = "数据库 ping 失败: " + err.Error()
		return status
	}

	dbStats := sqlDB.Stats()
	status.Healthy = true
	status.Reason = "数据库连接正常"
	status.OpenConnections = dbStats.OpenConnections
	status.InUseConnections = dbStats.InUse
	status.IdleConnections = dbStats.Idle
	status.MaxOpenConnections = dbStats.MaxOpenConnections
	status.WaitCount = dbStats.WaitCount

	if dbStats.MaxOpenConnections > 0 && dbStats.OpenConnections >= dbStats.MaxOpenConnections {
		status.Warning = "连接池接近上限，建议提升 max-open-conns 或排查慢查询"
	}

	return status
}

func (s *DBInspectorService) getRedisRuntimeStatus(ctx context.Context) RedisRuntimeStatus {
	status := RedisRuntimeStatus{
		Enabled:          global.GVA_CONFIG.System.UseRedis,
		Healthy:          false,
		Reason:           "Redis未启用",
		AutoFixSupported: global.GVA_CONFIG.System.UseRedis,
		ManualFixGuide:   redisManualFixGuide(),
	}

	if !global.GVA_CONFIG.System.UseRedis {
		status.Healthy = true
		return status
	}
	if global.GVA_REDIS == nil {
		status.Reason = "Redis客户端未初始化"
		return status
	}

	pong, err := global.GVA_REDIS.Ping(ctx).Result()
	if err != nil {
		status.Reason = "Redis ping 失败: " + err.Error()
		return status
	}

	status.Healthy = true
	status.Reason = "Redis连接正常: " + pong

	serverInfo, err := global.GVA_REDIS.Info(ctx, "server").Result()
	if err == nil {
		kv := parseRedisInfo(serverInfo)
		status.Version = kv["redis_version"]
	}
	memoryInfo, err := global.GVA_REDIS.Info(ctx, "memory").Result()
	if err == nil {
		kv := parseRedisInfo(memoryInfo)
		status.UsedMemory = kv["used_memory_human"]
	}
	keyspaceInfo, err := global.GVA_REDIS.Info(ctx, "keyspace").Result()
	if err == nil {
		status.Keyspaces = parseRedisKeyspaceInfo(keyspaceInfo)
	}

	return status
}

func (s *DBInspectorService) getTableMetricMaps(ctx context.Context) (map[string]int64, map[string]int64) {
	rowMap := map[string]int64{}
	sizeMap := map[string]int64{}

	switch s.normalizedDBType() {
	case "mysql":
		rows := []tableMetricRow{}
		err := global.GVA_DB.WithContext(ctx).Raw(`
			SELECT table_name AS table_name,
			       COALESCE(table_rows, 0) AS table_rows,
			       COALESCE(data_length + index_length, 0) AS total_bytes
			FROM information_schema.tables
			WHERE table_schema = DATABASE()
		`).Scan(&rows).Error
		if err != nil {
			global.GVA_LOG.Warn("数据库巡检: 获取mysql表统计失败", zap.Error(err))
			return rowMap, sizeMap
		}
		for _, item := range rows {
			if item.TableRows.Valid {
				rowMap[item.TableName] = item.TableRows.Int64
			}
			if item.TotalBytes.Valid {
				sizeMap[item.TableName] = item.TotalBytes.Int64
			}
		}
	case "postgresql":
		rows := []tableMetricRow{}
		err := global.GVA_DB.WithContext(ctx).Raw(`
			SELECT t.table_name AS table_name,
			       COALESCE(s.n_live_tup::bigint, 0) AS table_rows,
			       COALESCE(pg_total_relation_size(format('%I.%I', t.table_schema, t.table_name)), 0) AS total_bytes
			FROM information_schema.tables t
			LEFT JOIN pg_stat_user_tables s
			       ON s.relname = t.table_name AND s.schemaname = t.table_schema
			WHERE t.table_type = 'BASE TABLE' AND t.table_schema = 'public'
		`).Scan(&rows).Error
		if err != nil {
			global.GVA_LOG.Warn("数据库巡检: 获取postgres表统计失败", zap.Error(err))
			return rowMap, sizeMap
		}
		for _, item := range rows {
			if item.TableRows.Valid {
				rowMap[item.TableName] = item.TableRows.Int64
			}
			if item.TotalBytes.Valid {
				sizeMap[item.TableName] = item.TotalBytes.Int64
			}
		}
	case "sqlserver":
		rows := []tableMetricRow{}
		err := global.GVA_DB.WithContext(ctx).Raw(`
			SELECT t.NAME AS table_name,
			       SUM(p.rows) AS table_rows,
			       SUM(a.total_pages) * 8 * 1024 AS total_bytes
			FROM sys.tables t
			INNER JOIN sys.indexes i ON t.object_id = i.object_id
			INNER JOIN sys.partitions p ON i.object_id = p.object_id AND i.index_id = p.index_id
			INNER JOIN sys.allocation_units a ON p.partition_id = a.container_id
			GROUP BY t.NAME
		`).Scan(&rows).Error
		if err != nil {
			global.GVA_LOG.Warn("数据库巡检: 获取sqlserver表统计失败", zap.Error(err))
			return rowMap, sizeMap
		}
		for _, item := range rows {
			if item.TableRows.Valid {
				rowMap[item.TableName] = item.TableRows.Int64
			}
			if item.TotalBytes.Valid {
				sizeMap[item.TableName] = item.TotalBytes.Int64
			}
		}
	}

	return rowMap, sizeMap
}

func (s *DBInspectorService) getTableColumns(tableName string) ([]string, error) {
	columnTypes, err := global.GVA_DB.Migrator().ColumnTypes(tableName)
	if err != nil {
		return nil, err
	}
	columns := make([]string, 0, len(columnTypes))
	for _, ct := range columnTypes {
		columns = append(columns, ct.Name())
	}
	return columns, nil
}

func (s *DBInspectorService) getTablePrimaryKeyColumn(tableName string, fallbackColumns []string) string {
	columnTypes, err := global.GVA_DB.Migrator().ColumnTypes(tableName)
	if err == nil {
		for _, ct := range columnTypes {
			isPK, ok := ct.PrimaryKey()
			if ok && isPK {
				return ct.Name()
			}
		}
	}
	return pickPrimaryKeyColumn(fallbackColumns)
}

func (s *DBInspectorService) getTableRowCount(ctx context.Context, tableName string) (int64, error) {
	if !isValidIdentifier(tableName) {
		return 0, errors.New("非法表名")
	}
	sqlStr := fmt.Sprintf("SELECT COUNT(1) FROM %s", s.quoteIdentifier(tableName))
	var count int64
	if err := global.GVA_DB.WithContext(ctx).Raw(sqlStr).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *DBInspectorService) reconnectDatabase(ctx context.Context) error {
	var (
		db  *gorm.DB
		err error
	)

	switch s.normalizedDBType() {
	case "mysql":
		cfg := global.GVA_CONFIG.Mysql
		if cfg.Dbname == "" {
			return errors.New("mysql数据库名为空")
		}
		db, err = gorm.Open(mysql.Open(cfg.Dsn()), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		}
	case "postgresql":
		cfg := global.GVA_CONFIG.Pgsql
		if cfg.Dbname == "" {
			return errors.New("postgres数据库名为空")
		}
		db, err = gorm.Open(postgres.Open(cfg.Dsn()), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		}
	case "sqlserver":
		cfg := global.GVA_CONFIG.Mssql
		if cfg.Dbname == "" {
			return errors.New("sqlserver数据库名为空")
		}
		db, err = gorm.Open(sqlserver.Open(cfg.Dsn()), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		}
	case "sqlite":
		cfg := global.GVA_CONFIG.Sqlite
		if cfg.Dbname == "" {
			return errors.New("sqlite数据库名为空")
		}
		db, err = gorm.Open(sqlite.Open(cfg.Dsn()), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		}
	default:
		return fmt.Errorf("暂不支持的数据库类型: %s", global.GVA_CONFIG.System.DbType)
	}

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if err = sqlDB.PingContext(ctx); err != nil {
		return err
	}

	if global.GVA_DB != nil {
		if old, oldErr := global.GVA_DB.DB(); oldErr == nil {
			_ = old.Close()
		}
	}
	global.GVA_DB = db
	return nil
}

func (s *DBInspectorService) reconnectRedis(ctx context.Context) error {
	if !global.GVA_CONFIG.System.UseRedis {
		return errors.New("当前配置未启用redis")
	}

	redisCfg := global.GVA_CONFIG.Redis
	var client redis.UniversalClient
	if redisCfg.UseCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisCfg.ClusterAddrs,
			Password: redisCfg.Password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     redisCfg.Addr,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		})
	}

	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}

	if global.GVA_REDIS != nil {
		_ = global.GVA_REDIS.Close()
	}
	global.GVA_REDIS = client
	return nil
}

func (s *DBInspectorService) tryOptimizeTable(ctx context.Context, tableName string) error {
	tableName = strings.TrimSpace(tableName)
	if tableName == "" {
		return nil
	}
	if !isValidIdentifier(tableName) {
		return errors.New("待优化表名不合法")
	}
	if !global.GVA_DB.Migrator().HasTable(tableName) {
		return errors.New("待优化表不存在")
	}

	switch s.normalizedDBType() {
	case "mysql":
		sqlStr := fmt.Sprintf("OPTIMIZE TABLE %s", s.quoteIdentifier(tableName))
		return global.GVA_DB.WithContext(ctx).Exec(sqlStr).Error
	default:
		return errors.New("当前数据库类型不支持自动执行表优化")
	}
}

func (s *DBInspectorService) collectFileKeysByRange(ctx context.Context, tableName string, timeColumn string, startTime, endTime time.Time, fileColumns []string) ([]string, []string, error) {
	if len(fileColumns) == 0 {
		return nil, nil, nil
	}

	quotedColumns := make([]string, 0, len(fileColumns))
	for _, col := range fileColumns {
		if !isValidIdentifier(col) {
			continue
		}
		quotedColumns = append(quotedColumns, s.quoteIdentifier(col))
	}
	if len(quotedColumns) == 0 {
		return nil, []string{"未找到合法的文件字段，跳过联动文件删除"}, nil
	}

	selectSQL := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s >= ? AND %s <= ?",
		strings.Join(quotedColumns, ", "),
		s.quoteIdentifier(tableName),
		s.quoteIdentifier(timeColumn),
		s.quoteIdentifier(timeColumn),
	)

	rows := []map[string]interface{}{}
	if err := global.GVA_DB.WithContext(ctx).Raw(selectSQL, startTime, endTime).Scan(&rows).Error; err != nil {
		return nil, nil, err
	}

	keySet := map[string]struct{}{}
	for _, row := range rows {
		for _, col := range fileColumns {
			value, ok := findMapValueIgnoreCase(row, col)
			if !ok || value == nil {
				continue
			}
			candidates := extractStringCandidates(value)
			for _, candidate := range candidates {
				storageKey := s.resolveStorageKey(candidate, isKeyColumn(col))
				if storageKey != "" {
					keySet[storageKey] = struct{}{}
				}
			}
		}
	}

	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	if len(keys) == 0 {
		return nil, []string{"未解析到可删除的文件Key，已仅执行数据库真删除"}, nil
	}
	return keys, nil, nil
}

func (s *DBInspectorService) listOldestRowIDsByRange(ctx context.Context, tableName string, pkColumn string, timeColumn string, startTime time.Time, endTime time.Time, limit int64) ([]interface{}, error) {
	if limit <= 0 {
		return nil, nil
	}
	if !isValidIdentifier(tableName) || !isValidIdentifier(pkColumn) || !isValidIdentifier(timeColumn) {
		return nil, errors.New("查询参数不合法")
	}

	quotedPK := s.quoteIdentifier(pkColumn)
	quotedTable := s.quoteIdentifier(tableName)
	quotedTime := s.quoteIdentifier(timeColumn)

	var (
		sqlStr string
		args   []interface{}
	)

	switch s.normalizedDBType() {
	case "sqlserver":
		sqlStr = fmt.Sprintf(
			"SELECT TOP (%d) %s FROM %s WHERE %s >= ? AND %s <= ? ORDER BY %s ASC, %s ASC",
			limit,
			quotedPK,
			quotedTable,
			quotedTime,
			quotedTime,
			quotedTime,
			quotedPK,
		)
		args = []interface{}{startTime, endTime}
	default:
		sqlStr = fmt.Sprintf(
			"SELECT %s FROM %s WHERE %s >= ? AND %s <= ? ORDER BY %s ASC, %s ASC LIMIT ?",
			quotedPK,
			quotedTable,
			quotedTime,
			quotedTime,
			quotedTime,
			quotedPK,
		)
		args = []interface{}{startTime, endTime, limit}
	}

	rows := []map[string]interface{}{}
	if err := global.GVA_DB.WithContext(ctx).Raw(sqlStr, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	ids := make([]interface{}, 0, len(rows))
	for _, row := range rows {
		value, ok := findMapValueIgnoreCase(row, pkColumn)
		if !ok || value == nil {
			continue
		}
		ids = append(ids, normalizeSQLValue(value))
	}
	return ids, nil
}

func (s *DBInspectorService) collectFileKeysByIDs(ctx context.Context, tableName string, pkColumn string, ids []interface{}, fileColumns []string) ([]string, []string, error) {
	if len(ids) == 0 || len(fileColumns) == 0 {
		return nil, nil, nil
	}
	if !isValidIdentifier(tableName) || !isValidIdentifier(pkColumn) {
		return nil, nil, errors.New("表名或主键不合法")
	}

	quotedColumns := make([]string, 0, len(fileColumns))
	for _, col := range fileColumns {
		if !isValidIdentifier(col) {
			continue
		}
		quotedColumns = append(quotedColumns, s.quoteIdentifier(col))
	}
	if len(quotedColumns) == 0 {
		return nil, []string{"未找到合法的文件字段，跳过联动文件删除"}, nil
	}

	placeholders := buildSQLPlaceholders(len(ids))
	selectSQL := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s IN (%s)",
		strings.Join(quotedColumns, ", "),
		s.quoteIdentifier(tableName),
		s.quoteIdentifier(pkColumn),
		placeholders,
	)

	rows := []map[string]interface{}{}
	if err := global.GVA_DB.WithContext(ctx).Raw(selectSQL, ids...).Scan(&rows).Error; err != nil {
		return nil, nil, err
	}

	keySet := map[string]struct{}{}
	for _, row := range rows {
		for _, col := range fileColumns {
			value, ok := findMapValueIgnoreCase(row, col)
			if !ok || value == nil {
				continue
			}
			for _, candidate := range extractStringCandidates(value) {
				storageKey := s.resolveStorageKey(candidate, isKeyColumn(col))
				if storageKey != "" {
					keySet[storageKey] = struct{}{}
				}
			}
		}
	}

	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	if len(keys) == 0 {
		return nil, []string{"未解析到可删除的文件Key，已仅执行数据库真删除"}, nil
	}
	return keys, nil, nil
}

func (s *DBInspectorService) deleteRowsByIDs(ctx context.Context, tableName string, pkColumn string, ids []interface{}) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	if !isValidIdentifier(tableName) || !isValidIdentifier(pkColumn) {
		return 0, errors.New("表名或主键不合法")
	}

	deleteSQL := fmt.Sprintf(
		"DELETE FROM %s WHERE %s IN (%s)",
		s.quoteIdentifier(tableName),
		s.quoteIdentifier(pkColumn),
		buildSQLPlaceholders(len(ids)),
	)
	result := global.GVA_DB.WithContext(ctx).Exec(deleteSQL, ids...)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (s *DBInspectorService) deleteFilesByKeys(fileKeys []string) (int64, int64, []string) {
	if len(fileKeys) == 0 {
		return 0, 0, nil
	}

	oss := upload.NewOss()
	var successCount int64
	var failedCount int64
	warnings := []string{}

	for _, key := range fileKeys {
		if err := oss.DeleteFile(key); err != nil {
			failedCount++
			if len(warnings) < 20 {
				warnings = append(warnings, fmt.Sprintf("%s 删除失败: %s", key, err.Error()))
			}
			continue
		}
		successCount++
	}

	if failedCount > 0 {
		warnings = append(warnings, fmt.Sprintf("共 %d 个文件删除失败，请按失败Key手动处理", failedCount))
	}

	return successCount, failedCount, warnings
}

func (s *DBInspectorService) resolveStorageKey(raw string, fromKeyColumn bool) string {
	raw = strings.TrimSpace(strings.Trim(raw, `"'`))
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "data:") {
		return ""
	}

	if fromKeyColumn {
		return s.normalizeStorageKey(raw)
	}

	parsed, err := url.Parse(raw)
	if err == nil && parsed.Host != "" {
		if !s.isManagedHost(parsed.Hostname()) {
			return ""
		}
		raw = parsed.Path
	}

	raw = strings.Split(raw, "?")[0]
	raw = strings.Split(raw, "#")[0]
	raw = strings.ReplaceAll(raw, `\\`, "/")
	raw = strings.TrimSpace(raw)

	localPath := strings.TrimSpace(global.GVA_CONFIG.Local.Path)
	if localPath != "" {
		normalizedLocalPath := "/" + strings.Trim(strings.ReplaceAll(localPath, `\\`, "/"), "/")
		if strings.HasPrefix(raw, normalizedLocalPath+"/") {
			raw = strings.TrimPrefix(raw, normalizedLocalPath+"/")
		}
	}

	return s.normalizeStorageKey(raw)
}

func (s *DBInspectorService) normalizeStorageKey(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "/")
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.Contains(raw, "..") {
		return ""
	}

	if strings.EqualFold(global.GVA_CONFIG.System.OssType, "local") {
		raw = path.Base(raw)
	}
	return raw
}

func (s *DBInspectorService) isManagedHost(host string) bool {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback() || ip.IsPrivate()
	}
	if host == "localhost" {
		return true
	}

	allowed := s.allowedHosts()
	_, ok := allowed[host]
	return ok
}

func (s *DBInspectorService) allowedHosts() map[string]struct{} {
	urls := []string{
		global.GVA_CONFIG.AliyunOSS.BucketUrl,
		global.GVA_CONFIG.Minio.BucketUrl,
		global.GVA_CONFIG.Qiniu.ImgPath,
		global.GVA_CONFIG.AwsS3.BaseURL,
		global.GVA_CONFIG.TencentCOS.BaseURL,
		global.GVA_CONFIG.CloudflareR2.BaseURL,
	}
	set := map[string]struct{}{}
	for _, raw := range urls {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		host := parseHost(raw)
		if host != "" {
			set[host] = struct{}{}
		}
	}
	return set
}

func (s *DBInspectorService) quoteIdentifier(identifier string) string {
	if !isValidIdentifier(identifier) {
		return identifier
	}
	switch s.normalizedDBType() {
	case "mysql":
		return "`" + identifier + "`"
	case "sqlserver":
		return "[" + identifier + "]"
	default:
		return `"` + identifier + `"`
	}
}

func (s *DBInspectorService) normalizedDBType() string {
	raw := strings.ToLower(strings.TrimSpace(global.GVA_CONFIG.System.DbType))
	switch raw {
	case "postgres", "postgresql", "pgsql":
		return "postgresql"
	case "mssql", "sqlserver":
		return "sqlserver"
	case "sqlite", "sqlite3":
		return "sqlite"
	case "mysql":
		return "mysql"
	default:
		return raw
	}
}

func (s *DBInspectorService) currentDatabaseName() string {
	switch s.normalizedDBType() {
	case "mysql":
		return strings.TrimSpace(global.GVA_CONFIG.Mysql.Dbname)
	case "postgresql":
		return strings.TrimSpace(global.GVA_CONFIG.Pgsql.Dbname)
	case "sqlserver":
		return strings.TrimSpace(global.GVA_CONFIG.Mssql.Dbname)
	case "sqlite":
		return strings.TrimSpace(global.GVA_CONFIG.Sqlite.Dbname)
	default:
		return ""
	}
}

func (s *DBInspectorService) refreshTableStatsAfterDelete(ctx context.Context, tableName string) error {
	tableName = strings.TrimSpace(tableName)
	if tableName == "" {
		return nil
	}
	if !isValidIdentifier(tableName) {
		return errors.New("非法表名")
	}

	switch s.normalizedDBType() {
	case "mysql":
		sqlStr := fmt.Sprintf("ANALYZE TABLE %s", s.quoteIdentifier(tableName))
		return global.GVA_DB.WithContext(ctx).Exec(sqlStr).Error
	default:
		return nil
	}
}

func sumProfileStats(profiles []DBTableProfile) (sizeBytes int64, totalRows int64) {
	for _, item := range profiles {
		sizeBytes += item.SizeBytes
		totalRows += item.RowCount
	}
	if sizeBytes < 0 {
		sizeBytes = 0
	}
	if totalRows < 0 {
		totalRows = 0
	}
	return sizeBytes, totalRows
}

type tableClassification struct {
	Purpose        string
	UsageLabel     string
	SystemRequired bool
	LikelyUnused   bool
}

func classifyTable(tableName string, rowCount int64) tableClassification {
	lower := strings.ToLower(strings.TrimSpace(tableName))
	class := tableClassification{
		Purpose:    "业务数据表",
		UsageLabel: "业务数据",
	}

	if _, ok := systemRequiredTables[lower]; ok {
		class.Purpose = "系统核心权限/菜单/配置数据"
		class.UsageLabel = "系统必须"
		class.SystemRequired = true
		class.LikelyUnused = false
		return class
	}

	switch {
	case strings.HasPrefix(lower, "sys_login_log"):
		class.Purpose = "后台登录审计日志"
	case strings.HasPrefix(lower, "sys_operation_record"):
		class.Purpose = "后台操作审计日志"
	case strings.HasPrefix(lower, "sys_error"):
		class.Purpose = "系统错误日志"
	case strings.HasPrefix(lower, "sys_"):
		class.Purpose = "系统管理扩展数据"
	case strings.HasPrefix(lower, "exa_"):
		class.Purpose = "示例/附件模块数据"
	case strings.HasPrefix(lower, "client_"):
		class.Purpose = "客户端业务数据"
	case strings.HasPrefix(lower, "shop_"):
		class.Purpose = "商城业务数据"
	case strings.Contains(lower, "tryon"):
		class.Purpose = "试衣业务数据"
	}

	if rowCount == 0 && (strings.HasPrefix(lower, "exa_") || strings.Contains(lower, "log") || strings.Contains(lower, "history") || strings.Contains(lower, "cache")) {
		class.LikelyUnused = true
		class.UsageLabel = "疑似无用"
	}

	return class
}

func detectDateColumns(columns []string) []string {
	if len(columns) == 0 {
		return nil
	}
	priority := []string{"created_at", "create_time", "createdat", "updated_at", "update_time", "deleted_at", "time"}
	matches := make([]string, 0)
	for _, target := range priority {
		for _, col := range columns {
			if strings.EqualFold(col, target) {
				matches = append(matches, col)
			}
		}
	}
	return uniqueStrings(matches)
}

func detectFileColumns(columns []string) []string {
	if len(columns) == 0 {
		return nil
	}
	matches := make([]string, 0)
	for _, col := range columns {
		l := strings.ToLower(col)
		if l == "key" || strings.HasSuffix(l, "_key") || strings.Contains(l, "url") || strings.Contains(l, "image") || strings.Contains(l, "avatar") || strings.Contains(l, "thumb") || strings.Contains(l, "cover") {
			matches = append(matches, col)
		}
	}
	return uniqueStrings(matches)
}

func isKeyColumn(columnName string) bool {
	l := strings.ToLower(strings.TrimSpace(columnName))
	return l == "key" || strings.HasSuffix(l, "_key")
}

func pickPrimaryKeyColumn(columns []string) string {
	for _, col := range columns {
		if strings.EqualFold(col, "id") {
			return col
		}
	}
	for _, col := range columns {
		l := strings.ToLower(strings.TrimSpace(col))
		if strings.HasSuffix(l, "_id") {
			return col
		}
	}
	return ""
}

func buildSQLPlaceholders(count int) string {
	if count <= 0 {
		return ""
	}
	return strings.TrimRight(strings.Repeat("?,", count), ",")
}

func normalizeSQLValue(raw interface{}) interface{} {
	switch v := raw.(type) {
	case []byte:
		return string(v)
	default:
		return raw
	}
}

func pickDateColumn(columns []string) string {
	priority := []string{"created_at", "create_time", "createdat", "updated_at", "update_time"}
	for _, p := range priority {
		for _, col := range columns {
			if strings.EqualFold(col, p) {
				return col
			}
		}
	}
	return ""
}

func parseTimeRange(startRaw string, endRaw string) (time.Time, time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}

	parseAny := func(raw string) (time.Time, error) {
		raw = strings.TrimSpace(raw)
		for _, layout := range layouts {
			if t, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
				if layout == "2006-01-02" {
					return t, nil
				}
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("时间格式不正确: %s", raw)
	}

	start, err := parseAny(startRaw)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := parseAny(endRaw)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if !end.After(start) {
		return time.Time{}, time.Time{}, errors.New("结束时间必须晚于开始时间")
	}
	return start, end, nil
}

func formatBytes(size int64) string {
	if size <= 0 {
		return "0 B"
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	value := float64(size)
	idx := 0
	for value >= 1024 && idx < len(units)-1 {
		value /= 1024
		idx++
	}
	if idx == 0 {
		return fmt.Sprintf("%d %s", size, units[idx])
	}
	return fmt.Sprintf("%.2f %s", value, units[idx])
}

func parseRedisInfo(raw string) map[string]string {
	result := map[string]string{}
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return result
}

func parseRedisKeyspaceInfo(raw string) []RedisKeyspaceStat {
	stats := []RedisKeyspaceStat{}
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "db") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		kv := map[string]int64{}
		for _, pair := range strings.Split(parts[1], ",") {
			seg := strings.SplitN(pair, "=", 2)
			if len(seg) != 2 {
				continue
			}
			n, _ := strconv.ParseInt(strings.TrimSpace(seg[1]), 10, 64)
			kv[strings.TrimSpace(seg[0])] = n
		}
		stats = append(stats, RedisKeyspaceStat{
			DB:      parts[0],
			Keys:    kv["keys"],
			Expires: kv["expires"],
			AvgTTL:  kv["avg_ttl"],
		})
	}
	return stats
}

func dbManualFixGuide() []string {
	return []string{
		"检查 server/config.yaml（或 server/config.debug.yaml）中的 db-type、数据库地址、账号密码是否正确。",
		"确认数据库服务已启动且当前应用服务器可连通对应端口。",
		"若连接池耗尽，请排查慢SQL并适当提升 max-open-conns/max-idle-conns。",
		"必要时重启应用使数据库连接池重新建立。",
	}
}

func redisManualFixGuide() []string {
	return []string{
		"确认配置中 use-redis=true，且 redis.addr/redis.password/redis.db 正确。",
		"确认 Redis 实例运行正常并开放网络访问。",
		"若使用集群模式，请检查 cluster-addrs 是否完整可达。",
		"执行自动修复失败时，可先重启 Redis 再重试自动修复。",
	}
}

func parseHost(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(u.Hostname()))
}

func isValidIdentifier(identifier string) bool {
	return dbInspectorIdentifierPattern.MatchString(strings.TrimSpace(identifier))
}

func containsColumn(columns []string, column string) bool {
	for _, item := range columns {
		if strings.EqualFold(item, column) {
			return true
		}
	}
	return false
}

func uniqueStrings(values []string) []string {
	set := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" {
			continue
		}
		if _, exists := set[normalized]; exists {
			continue
		}
		set[normalized] = struct{}{}
		result = append(result, value)
	}
	return result
}

func extractStringCandidates(raw interface{}) []string {
	switch v := raw.(type) {
	case nil:
		return nil
	case string:
		return normalizeStringCandidates(v)
	case []byte:
		return normalizeStringCandidates(string(v))
	default:
		return normalizeStringCandidates(fmt.Sprintf("%v", v))
	}
}

func normalizeStringCandidates(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "<nil>") {
		return nil
	}

	if strings.HasPrefix(raw, "[") || strings.HasPrefix(raw, "{") {
		var data interface{}
		if err := json.Unmarshal([]byte(raw), &data); err == nil {
			return collectStringsFromJSON(data)
		}
	}

	if strings.Contains(raw, "\n") || strings.Contains(raw, ",") {
		parts := strings.FieldsFunc(raw, func(r rune) bool {
			return r == '\n' || r == ','
		})
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		return result
	}

	return []string{raw}
}

func collectStringsFromJSON(value interface{}) []string {
	result := []string{}
	switch v := value.(type) {
	case string:
		v = strings.TrimSpace(v)
		if v != "" {
			result = append(result, v)
		}
	case []interface{}:
		for _, item := range v {
			result = append(result, collectStringsFromJSON(item)...)
		}
	case map[string]interface{}:
		for _, item := range v {
			result = append(result, collectStringsFromJSON(item)...)
		}
	}
	return result
}

func findMapValueIgnoreCase(row map[string]interface{}, key string) (interface{}, bool) {
	for k, v := range row {
		if strings.EqualFold(strings.TrimSpace(k), strings.TrimSpace(key)) {
			return v, true
		}
	}
	return nil, false
}

func normalizeInspectorPage(page int, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultInspectorPageSize
	}
	if pageSize > maxInspectorPageSize {
		pageSize = maxInspectorPageSize
	}
	return page, pageSize
}
