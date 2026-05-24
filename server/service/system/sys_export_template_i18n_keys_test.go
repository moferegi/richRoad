package system

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/i18n"
)

func TestExportTemplateErrorKeys_HaveI18nEntries(t *testing.T) {
	keys := []string{
		errExportTemplateOrderInvalid,
		errExportTemplateOrderDirectionInvalid,
		errExportTemplateRawSQLSelectOnly,
		errExportTemplateRawSQLUnsafeFragment,
		errExportTemplateRawSQLKeywordForbidden,
		errExportTemplateImportSQLTypeInvalid,
		errExportTemplateImportSQLUnsafeFragment,
		errExportTemplateImportSQLKeywordForbidden,
		errExportTemplateImportSQLWhereRequired,
		errExportTemplateImportSQLParamRequired,
		errExportTemplateTableNameInvalid,
		errExportTemplateLimitNegative,
		errExportTemplateLimitExceeded,
		errExportTemplateConditionFromInvalid,
		errExportTemplateConditionColumnInvalid,
		errExportTemplateConditionOperatorInvalid,
		errExportTemplateJoinTypeInvalid,
		errExportTemplateJoinTableInvalid,
		errExportTemplateJoinOnInvalid,
		errExportTemplateRawSQLDisabled,
		errExportTemplateImportSQLDisabled,
		errExportTemplateNil,
		errExportTemplateOrderFieldInvalid,
		errExportTemplateParamsInvalid,
		errExportTemplateExcelDataNotEnough,
	}

	for _, key := range keys {
		if !i18n.HasKey(key) {
			t.Fatalf("missing i18n key for export template error: %s", key)
		}
	}
}
