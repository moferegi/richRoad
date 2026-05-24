package system

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
)

var AutoCodeMssql = new(autoCodeMssql)

type autoCodeMssql struct{}

func normalizeMSSQLIdentifier(name string) (string, error) {
	value := strings.TrimSpace(name)
	if value == "" {
		return "", errors.New("标识符不能为空")
	}
	if strings.ContainsAny(value, "\x00\r\n\t;") || strings.Contains(value, ".") {
		return "", errors.New("标识符不合法")
	}
	return "[" + strings.ReplaceAll(value, "]", "]]") + "]", nil
}

// GetDB 获取数据库的所有数据库名
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
func (s *autoCodeMssql) GetDB(businessDB string) (data []response.Db, err error) {
	var entities []response.Db
	sql := "select name AS 'database' from sys.databases;"
	if businessDB == "" {
		err = global.GVA_DB.Raw(sql).Scan(&entities).Error
	} else {
		err = global.GVA_DBList[businessDB].Raw(sql).Scan(&entities).Error
	}
	return entities, err
}

// GetTables 获取数据库的所有表名
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
func (s *autoCodeMssql) GetTables(businessDB string, dbName string) (data []response.Table, err error) {
	var entities []response.Table
	quotedDBName, qErr := normalizeMSSQLIdentifier(dbName)
	if qErr != nil {
		return entities, qErr
	}

	sql := fmt.Sprintf(`select name as 'table_name' from %s.DBO.sysobjects where xtype='U'`, quotedDBName)
	if businessDB == "" {
		err = global.GVA_DB.Raw(sql).Scan(&entities).Error
	} else {
		err = global.GVA_DBList[businessDB].Raw(sql).Scan(&entities).Error
	}

	return entities, err
}

// GetColumn 获取指定数据库和指定数据表的所有字段名,类型值等
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
func (s *autoCodeMssql) GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error) {
	var entities []response.Column
	quotedDBName, qErr := normalizeMSSQLIdentifier(dbName)
	if qErr != nil {
		return entities, qErr
	}
	tableName = strings.TrimSpace(tableName)
	if tableName == "" {
		return entities, errors.New("表名不能为空")
	}

	sql := fmt.Sprintf(`
SELECT
    sc.name AS column_name,
    st.name AS data_type,
    sc.max_length AS data_type_long,
    CASE
        WHEN pk.object_id IS NOT NULL THEN 1
        ELSE 0
    END AS primary_key,
    sc.column_id
FROM
	%s.sys.columns sc
JOIN
    sys.types st ON sc.user_type_id=st.user_type_id
LEFT JOIN
	%s.sys.objects so ON so.name = ? AND so.type='U'
LEFT JOIN
    %s.sys.indexes si ON si.object_id = so.object_id AND si.is_primary_key = 1
LEFT JOIN
    %s.sys.index_columns sic ON sic.object_id = si.object_id AND sic.index_id = si.index_id AND sic.column_id = sc.column_id
LEFT JOIN
    %s.sys.key_constraints pk ON pk.object_id = si.object_id
WHERE
    st.is_user_defined=0 AND sc.object_id = so.object_id
ORDER BY
    sc.column_id
`, quotedDBName, quotedDBName, quotedDBName, quotedDBName, quotedDBName)

	if businessDB == "" {
		err = global.GVA_DB.Raw(sql, tableName).Scan(&entities).Error
	} else {
		err = global.GVA_DBList[businessDB].Raw(sql, tableName).Scan(&entities).Error
	}

	return entities, err
}
