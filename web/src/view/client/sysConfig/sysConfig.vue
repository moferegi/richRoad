<template>
  <div>
    <!-- 分组标签切换 -->
    <div class="gva-search-box">
      <el-radio-group v-model="activeGroup" @change="handleGroupChange">
        <el-radio-button label="">全部</el-radio-button>
        <el-radio-button v-for="g in groupList" :key="g.value" :label="g.value">{{ g.label }}</el-radio-button>
      </el-radio-group>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" stripe>
        <el-table-column prop="configGroup" label="配置组" width="120">
          <template #default="scope">
            <el-tag size="small" :type="groupTagType(scope.row.configGroup)">{{ groupLabel(scope.row.configGroup) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="configName" label="配置名称" width="180" />
        <el-table-column prop="configKey" label="配置键" width="220">
          <template #default="scope">
            <span class="config-key">{{ scope.row.configKey }}</span>
          </template>
        </el-table-column>
        <el-table-column label="配置值" min-width="250">
          <template #default="scope">
            <!-- 布尔类型展示开关 -->
            <template v-if="isBooleanConfig(scope.row)">
              <el-switch
                :model-value="scope.row.configValue === 'true'"
                @change="(val) => quickToggle(scope.row, val)"
                active-text="开"
                inactive-text="关"
              />
            </template>
            <!-- 密钥类脱敏展示 -->
            <template v-else-if="isSecretConfig(scope.row)">
              <span>{{ maskSecretValue(scope.row.configValue) }}</span>
            </template>
            <template v-else-if="isCurrencySymbolConfig(scope.row)">
              <span>{{ formatCurrencySymbolPreview(scope.row.configValue) }}</span>
            </template>
            <template v-else-if="isTryonTutorialConfig(scope.row)">
              <el-tooltip :content="formatTutorialPreview(scope.row, true)" placement="top">
                <span class="tutorial-config-preview">{{ formatTutorialPreview(scope.row) }}</span>
              </el-tooltip>
            </template>
            <!-- 试衣模型可视化摘要 -->
            <template v-else-if="isTryonModelsConfig(scope.row)">
              <div class="tryon-model-summary">
                <el-tag size="small" type="warning">{{ tryonModelsCount(scope.row.configValue) }} 个模型</el-tag>
                <span class="tryon-model-summary-text">启用 {{ enabledTryonModelsCount(scope.row.configValue) }} 个</span>
              </div>
            </template>
            <!-- 试衣币充值套餐摘要 -->
            <template v-else-if="isTryonRechargePlansConfig(scope.row)">
              <div class="tryon-model-summary">
                <el-tag size="small" type="warning">{{ rechargePlansCount(scope.row.configValue) }} 个套餐</el-tag>
                <span class="tryon-model-summary-text">支持点数、币名与多语言价格（基础分价 + 汇率换算）</span>
              </div>
            </template>
            <!-- 支付方式配置摘要 -->
            <template v-else-if="isPaymentManualMethodsConfig(scope.row)">
              <div class="tryon-model-summary">
                <el-tag size="small" type="warning">{{ paymentMethodsCount(scope.row.configValue) }} 种方式</el-tag>
                <span class="tryon-model-summary-text">支持多语言名称、上传图片、复制文案、开关与排序</span>
              </div>
            </template>
            <!-- Uni期望支付方式配置摘要 -->
            <template v-else-if="isPaymentUniPreferredMethodsConfig(scope.row)">
              <div class="tryon-model-summary">
                <el-tag size="small" type="success">{{ paymentPreferredMethodsCount(scope.row.configValue) }} 种方式</el-tag>
                <span class="tryon-model-summary-text">仅用于uni联系客服支付：多语言名称、上传图片、复制文案、开关与排序</span>
              </div>
            </template>
            <!-- 颜色类型 -->
            <template v-else-if="isColorConfig(scope.row)">
              <div class="color-preview">
                <span class="color-dot" :style="{ background: scope.row.configValue }" />
                <span>{{ scope.row.configValue }}</span>
              </div>
            </template>
            <!-- 长文本截断 -->
            <template v-else-if="scope.row.configValue && scope.row.configValue.length > 60">
              <el-tooltip :content="scope.row.configValue" placement="top">
                <span>{{ scope.row.configValue.substring(0, 60) }}...</span>
              </el-tooltip>
            </template>
            <template v-else>
              <span>{{ scope.row.configValue || '-' }}</span>
            </template>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="200" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openEdit(scope.row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="searchInfo.page"
          :page-size="searchInfo.pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <div class="gva-search-box model-call-log-search">
      <el-form :inline="true" :model="modelLogSearch" class="model-call-log-form">
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="modelLogDateRange"
            type="datetimerange"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            range-separator="至"
          />
        </el-form-item>
        <el-form-item label="场景">
          <el-select v-model="modelLogSearch.sceneType" clearable style="width: 140px">
            <el-option label="试衣 clothes" value="clothes" />
            <el-option label="试鞋 shoes" value="shoes" />
            <el-option label="取衣 takeoff" value="takeoff" />
          </el-select>
        </el-form-item>
        <el-form-item label="用途">
          <el-select v-model="modelLogSearch.modelUsage" clearable style="width: 130px">
            <el-option label="tryon" value="tryon" />
            <el-option label="refiner" value="refiner" />
            <el-option label="parsing" value="parsing" />
            <el-option label="beautify" value="beautify" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="modelLogSearch.status" clearable style="width: 130px" :disabled="modelLogSearch.onlyFailed">
            <el-option label="success" value="success" />
            <el-option label="processing" value="processing" />
            <el-option label="failed" value="failed" />
            <el-option label="error" value="error" />
          </el-select>
        </el-form-item>
        <el-form-item label="快速筛选">
          <el-switch v-model="modelLogSearch.onlyFailed" active-text="仅失败" inactive-text="全部" />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            v-model="modelLogSearch.keyword"
            clearable
            placeholder="taskNo/requestID/modelKey"
            style="width: 260px"
            @keyup.enter="handleModelLogSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="modelLogLoading" @click="handleModelLogSearch">查询</el-button>
          <el-button @click="resetModelLogSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="modelCallLogList" stripe v-loading="modelLogLoading">
        <el-table-column label="时间" width="170">
          <template #default="scope">
            {{ formatDateTime(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="taskNo" label="任务号" min-width="150" />
        <el-table-column prop="requestID" label="请求ID" min-width="140" show-overflow-tooltip />
        <el-table-column prop="sceneType" label="场景" width="90" />
        <el-table-column prop="behavior" label="行为" min-width="140" show-overflow-tooltip />
        <el-table-column prop="callStage" label="阶段" min-width="150" show-overflow-tooltip />
        <el-table-column prop="modelKey" label="模型key" min-width="170" show-overflow-tooltip />
        <el-table-column prop="modelUsage" label="用途" width="90" />
        <el-table-column prop="provider" label="provider" width="110" show-overflow-tooltip />
        <el-table-column prop="tokenFingerprint" label="token指纹" min-width="150" show-overflow-tooltip />
        <el-table-column label="耗时" width="90">
          <template #default="scope">
            {{ Number(scope.row.durationMs || 0) }}ms
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag size="small" :type="modelCallStatusTag(scope.row.status)">{{ scope.row.status || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="errorMessage" label="错误信息" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="90" fixed="right">
          <template #default="scope">
            <el-button link type="primary" @click="openModelLogDetail(scope.row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="modelLogSearch.page"
          :page-size="modelLogSearch.pageSize"
          :page-sizes="[20, 50, 100]"
          :total="modelLogTotal"
          @current-change="handleModelLogCurrentChange"
          @size-change="handleModelLogSizeChange"
        />
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" title="编辑参数" :width="editDialogWidth">
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="配置键">
          <el-input :model-value="editForm.configKey" disabled />
        </el-form-item>
        <el-form-item label="配置名称">
          <el-input :model-value="editForm.configName" disabled />
        </el-form-item>
        <el-form-item label="配置值">
          <!-- 布尔类型用开关 -->
          <template v-if="isBooleanConfig(editForm)">
            <el-switch
              v-model="editBoolValue"
              active-text="开"
              inactive-text="关"
            />
          </template>
          <!-- 颜色类型用颜色选择器 -->
          <template v-else-if="isColorConfig(editForm)">
            <el-color-picker v-model="editForm.configValue" show-alpha />
            <el-input v-model="editForm.configValue" class="ml-2" style="width: 200px" />
          </template>
          <!-- 维护背景图上传 -->
          <template v-else-if="isMaintenanceBgImageConfig(editForm)">
            <SelectImage
              v-model="editForm.configValue"
              file-type="image"
              :default-folder="MAINTENANCE_BG_UPLOAD_FOLDER"
              :fixed-upload-folder="true"
            />
          </template>
          <!-- 试衣模型可视化编辑 -->
          <template v-else-if="isTryonModelsConfig(editForm)">
            <div class="tryon-model-editor">
              <div class="tryon-model-toolbar">
                <el-button type="primary" plain size="small" @click="addTryonModel">新增模型</el-button>
              </div>

              <div v-if="tryonModels.length === 0" class="tryon-model-empty">
                暂无模型，点击“新增模型”开始配置
              </div>

              <el-tabs
                v-else
                v-model="tryonModelActiveTab"
                tab-position="left"
                class="tryon-model-tabs"
                @tab-change="handleTryonModelTabChange"
              >
                <el-tab-pane
                  v-for="(model, index) in tryonModels"
                  :key="model.__uid"
                  :name="model.__uid"
                >
                  <template #label>
                    <div class="tryon-model-tab-label">
                      <span class="tryon-model-tab-title">{{ displayI18nText(model.name) || model.key || ('模型' + (index + 1)) }}</span>
                      <el-tag size="small" :type="model.enabled ? 'success' : 'info'">
                        {{ model.enabled ? '启用' : '关闭' }}
                      </el-tag>
                      <el-tag size="small" :type="getModelUsageTagType(model)">
                        {{ getModelUsageLabel(model) }}
                      </el-tag>
                    </div>
                  </template>

                  <div class="tryon-model-panel">
                    <div v-if="isAliyunModelForQuota(model)" class="tryon-model-quota">
                      <div class="tryon-model-quota-head">
                        <span class="tryon-model-label">阿里免费额度估算</span>
                        <el-button
                          link
                          size="small"
                          :loading="isAliyunQuotaLoading(model)"
                          @click.stop="refreshAliyunQuota(model, true)"
                        >
                          刷新
                        </el-button>
                      </div>
                      <div class="tryon-model-quota-content">
                        <template v-if="getAliyunQuotaItem(model)">
                          <el-tag type="success">
                            剩余 {{ getAliyunQuotaItem(model).remainingEstimate }} / {{ getAliyunQuotaItem(model).freeQuotaTotal }}
                          </el-tag>
                        </template>
                        <el-tag v-else type="info">暂无数据</el-tag>
                        <span class="tryon-model-summary-text" v-if="getAliyunQuotaItem(model)">
                          已用 {{ getAliyunQuotaItem(model).usedSuccessCount }}
                        </span>
                        <span class="tryon-model-summary-text" v-if="getAliyunQuotaItem(model)">
                          {{ formatQuotaRefreshTime(getAliyunQuotaItem(model).lastRefreshedAt) }}
                        </span>
                      </div>
                      <div class="tryon-model-token-quota-list" v-if="getAliyunQuotaItem(model)?.tokenQuotaList?.length">
                        <div class="tryon-model-token-quota-item" v-for="tokenItem in getAliyunQuotaItem(model).tokenQuotaList" :key="tokenItem.tokenFingerprint || tokenItem.tokenMasked">
                          <span class="tryon-model-token-label">{{ tokenItem.tokenMasked || 'token' }} ({{ tokenItem.tokenFingerprintShort || '-' }})</span>
                          <el-tag size="small" :type="tokenItem.exhausted ? 'danger' : 'success'">
                            剩余 {{ tokenItem.remainingEstimate }} / {{ tokenItem.freeQuotaTotal }}
                          </el-tag>
                          <span class="tryon-model-summary-text">已用 {{ tokenItem.usedSuccessCount }}</span>
                        </div>
                      </div>
                      <div class="tryon-model-quota-tip">
                        本地估算值（按 token 指纹统计成功调用并扣减原始额度），官方免费额度请以百炼控制台为准；
                        可在此处查看每个 token 的已用/剩余额度，也可在页面上方“模型调用日志”按 token 指纹追踪明细。
                      </div>
                      <div class="tryon-model-quota-error" v-if="getAliyunQuotaError(model)">
                        {{ getAliyunQuotaError(model) }}
                      </div>
                    </div>

                    <div class="tryon-model-grid">
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">模型键 key</span>
                        <el-input v-model="model.key" placeholder="如 aliyun_aitryon" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">模型类型 model</span>
                        <el-input v-model="model.model" placeholder="如 aitryon / aitryon-plus" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">用途 modelUsage</span>
                        <el-select v-model="model.modelUsage" style="width: 100%">
                          <el-option label="试衣模型 tryon" value="tryon" />
                          <el-option label="图片精修 refiner" value="refiner" />
                          <el-option label="分割模型 parsing" value="parsing" />
                          <el-option label="智能美肤 beautify" value="beautify" />
                        </el-select>
                      </div>

                      <div class="tryon-model-field">
                        <span class="tryon-model-label">提供商 provider</span>
                        <el-input v-model="model.provider" placeholder="如 aliyun" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">运行模式 mode</span>
                        <el-select v-model="model.mode" style="width: 100%">
                          <el-option label="prod" value="prod" />
                          <el-option label="mock_success" value="mock_success" />
                        </el-select>
                      </div>

                      <div class="tryon-model-field">
                        <span class="tryon-model-label">单次消耗 cost</span>
                        <el-input-number v-model="model.cost" :min="0" :step="1" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">是否启用</span>
                        <el-switch v-model="model.enabled" active-text="开" inactive-text="关" />
                      </div>

                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">场景 scenes</span>
                        <el-checkbox-group v-model="model.scenes">
                          <el-checkbox value="clothes">试衣 clothes</el-checkbox>
                          <el-checkbox value="shoes">试鞋 shoes</el-checkbox>
                          <el-checkbox value="takeoff">取衣 takeoff</el-checkbox>
                        </el-checkbox-group>
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">模型地址 url</span>
                        <el-input v-model="model.url" placeholder="如 https://dashscope.aliyuncs.com/api/v1/services/... 或 https://yisol-idm-vton.hf.space" />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">查询地址 taskQueryUrl</span>
                        <el-input v-model="model.taskQueryUrl" placeholder="如 https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}" />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">模型 token</span>
                        <el-input v-model="model.token" type="password" show-password placeholder="模型级主 token（留空表示该模型无鉴权）" />
                        <div class="tryon-model-hint">主 token 优先使用。</div>
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">备用 token 列表 tokenBackups</span>
                        <el-input
                          v-model="model.tokenBackupText"
                          type="textarea"
                          :rows="4"
                          placeholder="每行一个备用 token；主 token 失效或额度不足时自动切换"
                        />
                        <div class="tryon-model-hint">按顺序回退，直到找到可用 token；全部不可用时才返回失败。</div>
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">每 token 原始额度 tokenQuotas</span>
                        <div class="tryon-model-token-quota-editor">
                          <div
                            class="tryon-model-token-quota-edit-row"
                            v-for="tokenRow in getTokenQuotaRows(model)"
                            :key="`quota-${model.__uid}-${tokenRow.token}`"
                          >
                            <el-input
                              class="tryon-model-token-input"
                              :model-value="tokenRow.token"
                              type="password"
                              show-password
                              readonly
                              placeholder="token"
                            />
                            <el-input-number
                              :model-value="tokenRow.freeQuotaTotal"
                              :min="0"
                              :step="1"
                              controls-position="right"
                              @change="updateTokenQuotaValue(model, tokenRow.token, $event)"
                            />
                          </div>
                          <div v-if="!hasTokenQuotaRows(model)" class="tryon-model-hint">
                            请先填写“模型 token”或“备用 token 列表”，这里会自动生成一对一额度输入。
                          </div>
                        </div>
                        <div class="tryon-model-hint">按 token 值绑定额度，调换 token 顺序不会影响历史已用统计。</div>
                      </div>

                      <div class="tryon-model-subtitle" v-if="isTryonUsageModel(model)">Gradio / HuggingFace Space 参数</div>
                      <div class="tryon-model-field" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">API 名称 apiName</span>
                        <el-input v-model="model.apiName" placeholder="/tryon" />
                      </div>
                      <div class="tryon-model-field" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">服装描述 garmentDes</span>
                        <el-input v-model="model.garmentDes" placeholder="clothing item" />
                      </div>
                      <div class="tryon-model-field" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">自动蒙版 isChecked</span>
                        <el-switch v-model="model.isChecked" active-text="开" inactive-text="关" />
                      </div>
                      <div class="tryon-model-field" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">自动裁剪 isCheckedCrop</span>
                        <el-switch v-model="model.isCheckedCrop" active-text="开" inactive-text="关" />
                      </div>
                      <div class="tryon-model-field" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">降噪步数 denoiseSteps</span>
                        <el-input-number v-model="model.denoiseSteps" :min="1" :max="100" :step="1" />
                      </div>
                      <div class="tryon-model-field" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">随机种子 seed</span>
                        <el-input-number v-model="model.seed" :min="-1" :step="1" />
                      </div>

                      <div class="tryon-model-field" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">分辨率 resolution</span>
                        <el-input-number v-model="model.resolution" :min="-1" :step="1" />
                      </div>
                      <div class="tryon-model-field" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">人脸修复 restoreFace</span>
                        <el-switch v-model="model.restoreFace" active-text="开" inactive-text="关" />
                      </div>

                      <div class="tryon-model-field full" v-if="isParsingUsageModel(model)">
                        <span class="tryon-model-label">取衣分割 clothesType</span>
                        <el-checkbox-group v-model="model.clothesType">
                          <el-checkbox value="upper">upper</el-checkbox>
                          <el-checkbox value="lower">lower</el-checkbox>
                        </el-checkbox-group>
                      </div>
                      <div class="tryon-model-field" v-if="isParsingUsageModel(model)">
                        <span class="tryon-model-label">分割额外消耗 parsingExtraCost</span>
                        <el-input-number v-model="model.parsingExtraCost" :min="0" :step="1" />
                      </div>

                      <template v-if="isTryonUsageModel(model)">
                        <div class="tryon-model-subtitle">阿里取衣分割增强（仅 aitryon / aitryon-plus 生效）</div>
                        <div class="tryon-model-field full">
                          <span class="tryon-model-label">绑定分割模型 parsingModelKey</span>
                          <el-select
                            v-model="model.parsingModelKey"
                            filterable
                            clearable
                            allow-create
                            default-first-option
                            style="width: 100%"
                            placeholder="留空自动选择可用分割模型"
                          >
                            <el-option
                              v-for="parsingKey in parsingModelKeyOptions"
                              :key="parsingKey"
                              :label="parsingKey"
                              :value="parsingKey"
                            />
                          </el-select>
                        </div>
                      </template>

                      <div class="tryon-model-subtitle" v-if="isTryonUsageModel(model)">模型引用配置</div>
                      <div class="tryon-model-field full" v-if="isTryonUsageModel(model)">
                        <span class="tryon-model-label">绑定精修模型 refinerModelKey</span>
                        <el-select
                          v-model="model.refinerModelKey"
                          filterable
                          clearable
                          allow-create
                          default-first-option
                          style="width: 100%"
                          placeholder="留空表示当前模型不启用精修"
                        >
                          <el-option
                            v-for="refinerKey in refinerModelKeyOptions"
                            :key="refinerKey"
                            :label="refinerKey"
                            :value="refinerKey"
                          />
                        </el-select>
                      </div>

                      <div class="tryon-model-subtitle" v-if="isRefinerUsageModel(model)">图片精修参数</div>
                      <div class="tryon-model-field" v-if="isRefinerUsageModel(model)">
                        <span class="tryon-model-label">精修额外消耗 refinerExtraCost</span>
                        <el-input-number v-model="model.refinerExtraCost" :min="0" :step="1" />
                      </div>
                      <div class="tryon-model-field" v-if="isRefinerUsageModel(model)">
                        <span class="tryon-model-label">精修性别 refinerGender</span>
                        <el-select v-model="model.refinerGender" style="width: 100%">
                          <el-option label="woman" value="woman" />
                          <el-option label="man" value="man" />
                        </el-select>
                      </div>

                      <div class="tryon-model-field" v-if="!isBeautifyUsageModel(model)">
                        <span class="tryon-model-label">免费额度总数 freeQuotaTotal</span>
                        <el-input-number v-model="model.freeQuotaTotal" :min="0" :step="1" />
                      </div>

                      <template v-if="model.modelUsage === 'beautify'">
                        <div class="tryon-model-subtitle">智能美肤配置（独立模型）</div>
                        <div class="tryon-model-field">
                          <span class="tryon-model-label">支持美肤 supportsBeautify</span>
                          <el-switch v-model="model.supportsBeautify" active-text="开" inactive-text="关" disabled />
                        </div>
                        <div class="tryon-model-field">
                          <span class="tryon-model-label">美肤额外消耗 beautifyExtraCost</span>
                          <el-input-number v-model="model.beautifyExtraCost" :min="0" :step="1" />
                        </div>
                        <div class="tryon-model-field">
                          <span class="tryon-model-label">美肤模型/动作 beautifyModel</span>
                          <el-input v-model="model.beautifyModel" placeholder="如 custom_beautify / RetouchSkin" />
                        </div>
                        <div class="tryon-model-field">
                          <span class="tryon-model-label">磨皮强度 beautifyRetouchDegree</span>
                          <el-input-number v-model="model.beautifyRetouchDegree" :min="0" :max="1.5" :step="0.1" :precision="2" />
                        </div>
                        <div class="tryon-model-field">
                          <span class="tryon-model-label">美白强度 beautifyWhiteningDegree</span>
                          <el-input-number v-model="model.beautifyWhiteningDegree" :min="0" :max="1.5" :step="0.1" :precision="2" />
                        </div>
                        <div class="tryon-model-field full">
                          <span class="tryon-model-label">美肤地址 beautifyUrl</span>
                          <el-input v-model="model.beautifyUrl" placeholder="可配置任意 provider 地址；阿里可留空默认 facebody" />
                        </div>
                        <div class="tryon-model-field full">
                          <span class="tryon-model-label">AccessKeyId beautifyAccessKeyId</span>
                          <el-input v-model="model.beautifyAccessKeyId" placeholder="阿里模型可配置，其他 provider 可留空" />
                        </div>
                        <div class="tryon-model-field full">
                          <span class="tryon-model-label">AccessKeySecret beautifyAccessKeySecret</span>
                          <el-input v-model="model.beautifyAccessKeySecret" type="password" show-password placeholder="阿里模型可配置，其他 provider 可留空" />
                        </div>
                        <div class="tryon-model-field full">
                          <span class="tryon-model-label">SecurityToken beautifySecurityToken</span>
                          <el-input v-model="model.beautifySecurityToken" placeholder="STS 场景可配置" />
                        </div>
                        <div class="tryon-model-field full">
                          <span class="tryon-model-label">AKSK 合并串 beautifyToken</span>
                          <el-input v-model="model.beautifyToken" type="password" show-password placeholder="支持 ak:sk[:securityToken]" />
                        </div>
                      </template>

                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="model.name"
                          :languages="multilingualLangOptions"
                          title="名称多语言 name"
                          input-type="input"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>

                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="model.desc"
                          :languages="multilingualLangOptions"
                          title="说明多语言 desc"
                          input-type="textarea"
                          :rows="2"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>

                      <template v-if="isRefinerUsageModel(model)">
                        <div class="tryon-model-field full">
                          <MultiLangEditor
                            :model="model.refinerDesc"
                            :languages="multilingualLangOptions"
                            title="精修说明多语言 refinerDesc"
                            input-type="textarea"
                            :rows="2"
                            :use-tabs="true"
                            tab-type="card"
                          />
                        </div>
                      </template>

                      <template v-if="model.modelUsage === 'beautify'">
                        <div class="tryon-model-field full">
                          <MultiLangEditor
                            :model="model.beautifyDesc"
                            :languages="multilingualLangOptions"
                            title="美肤说明多语言 beautifyDesc"
                            input-type="textarea"
                            :rows="2"
                            :use-tabs="true"
                            tab-type="card"
                          />
                        </div>
                      </template>
                    </div>

                    <div class="tryon-model-actions">
                      <el-button size="small" @click="cloneTryonModel(index)">复制</el-button>
                      <el-button size="small" type="danger" plain @click="removeTryonModel(index)">删除</el-button>
                    </div>
                  </div>
                </el-tab-pane>
              </el-tabs>
            </div>
          </template>
          <!-- 试衣币充值套餐可视化编辑 -->
          <template v-else-if="isTryonRechargePlansConfig(editForm)">
            <div class="tryon-recharge-editor">
              <div class="tryon-model-toolbar">
                <el-button type="primary" plain size="small" @click="addRechargePlan">新增套餐</el-button>
              </div>
              <div v-if="rechargePlans.length === 0" class="tryon-model-empty">
                暂无套餐，点击“新增套餐”开始配置
              </div>
              <el-collapse v-else>
                <el-collapse-item v-for="(plan, index) in rechargePlans" :key="plan.__uid" :name="plan.__uid">
                  <template #title>
                    <div class="tryon-model-title">
                      <span>{{ displayI18nText(plan.points) || ('套餐' + (index + 1)) }} {{ displayI18nText(plan.coinLabel) }}</span>
                      <el-tag size="small" type="success">{{ formatRechargePlanPricePreview(plan) }}</el-tag>
                    </div>
                  </template>
                  <div class="tryon-model-panel">
                    <div class="tryon-model-grid">
                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="plan.points"
                          :languages="multilingualLangOptions"
                          title="点数多语言 points"
                          input-type="input"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>
                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="plan.coinLabel"
                          :languages="multilingualLangOptions"
                          title="币名多语言 coinLabel"
                          input-type="input"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">基础价格 price（单位：分）</span>
                        <div class="w-full">
                          <el-input-number v-model="plan.price" :min="0" :step="1" />
                          <div class="mt-2 flex items-center gap-2">
                            <el-button
                              type="primary"
                              plain
                              size="small"
                              @click="openRechargePlanRateDialog(index)"
                            >汇率换算</el-button>
                            <span class="tryon-model-summary-text">槽位：{{ Object.keys(plan.priceI18n || {}).length }}</span>
                          </div>
                        </div>
                      </div>
                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="plan.priceI18n"
                          :languages="multilingualLangOptions"
                          title="价格多语言 priceI18n（单位：分）"
                          input-type="input"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>
                    </div>
                    <div class="tryon-model-actions">
                      <el-button size="small" @click="cloneRechargePlan(index)">复制</el-button>
                      <el-button size="small" type="danger" plain @click="removeRechargePlan(index)">删除</el-button>
                    </div>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
          </template>
          <!-- 支付方式配置可视化编辑 -->
          <template v-else-if="isPaymentManualMethodsConfig(editForm)">
            <div class="tryon-recharge-editor">
              <div class="tryon-model-toolbar">
                <el-button type="primary" plain size="small" @click="addPaymentMethod">新增方式</el-button>
              </div>
              <div v-if="paymentManualMethods.length === 0" class="tryon-model-empty">
                暂无支付方式，点击“新增方式”开始配置
              </div>
              <el-collapse v-else>
                <el-collapse-item v-for="(method, index) in paymentManualMethods" :key="method.__uid" :name="method.__uid">
                  <template #title>
                    <div class="tryon-model-title">
                      <span>{{ displayI18nText(method.name, method.key || ('支付方式' + (index + 1))) }}</span>
                      <el-tag size="small" :type="method.enabled ? 'success' : 'info'">{{ method.enabled ? '启用' : '关闭' }}</el-tag>
                    </div>
                  </template>
                  <div class="tryon-model-panel">
                    <div class="tryon-model-grid">
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">唯一键 key</span>
                        <el-input v-model="method.key" placeholder="如 qrcode / contact / wechat" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">排序 sort</span>
                        <el-input-number v-model="method.sort" :min="0" :step="1" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">是否启用</span>
                        <el-switch v-model="method.enabled" active-text="开" inactive-text="关" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">是否人工渠道</span>
                        <el-switch v-model="method.manual" active-text="是" inactive-text="否" />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">方式图片(上传)</span>
                        <SelectImage
                          v-model="method.image"
                          file-type="image"
                          :default-folder="PAYMENT_MANUAL_METHOD_UPLOAD_FOLDER"
                          :fixed-upload-folder="true"
                        />
                      </div>

                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="method.name"
                          :languages="multilingualLangOptions"
                          title="名称多语言 name"
                          input-type="input"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>

                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="method.copyText"
                          :languages="multilingualLangOptions"
                          title="复制信息多语言 copyText"
                          input-type="textarea"
                          :rows="2"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>
                    </div>
                    <div class="tryon-model-actions">
                      <el-button size="small" @click="clonePaymentMethod(index)">复制</el-button>
                      <el-button size="small" type="danger" plain @click="removePaymentMethod(index)">删除</el-button>
                    </div>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
          </template>
          <!-- Uni期望支付方式可视化编辑（独立于 payment_manual_methods） -->
          <template v-else-if="isPaymentUniPreferredMethodsConfig(editForm)">
            <div class="tryon-recharge-editor">
              <div class="tryon-model-toolbar">
                <el-button type="primary" plain size="small" @click="addPaymentPreferredMethod">新增方式</el-button>
              </div>
              <div v-if="paymentPreferredMethods.length === 0" class="tryon-model-empty">
                暂无期望支付方式，点击“新增方式”开始配置
              </div>
              <el-collapse v-else>
                <el-collapse-item v-for="(method, index) in paymentPreferredMethods" :key="method.__uid" :name="method.__uid">
                  <template #title>
                    <div class="tryon-model-title">
                      <span>{{ displayI18nText(method.name, method.key || ('期望方式' + (index + 1))) }}</span>
                      <el-tag size="small" :type="method.enabled ? 'success' : 'info'">{{ method.enabled ? '启用' : '关闭' }}</el-tag>
                    </div>
                  </template>
                  <div class="tryon-model-panel">
                    <div class="tryon-model-grid">
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">唯一键 key</span>
                        <el-input v-model="method.key" placeholder="如 wechat / alipay / bank_card_cn" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">排序 sort</span>
                        <el-input-number v-model="method.sort" :min="0" :step="1" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">是否启用</span>
                        <el-switch v-model="method.enabled" active-text="开" inactive-text="关" />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">方式图片(上传)</span>
                        <SelectImage
                          v-model="method.image"
                          file-type="image"
                          :default-folder="PAYMENT_PREFERRED_METHOD_UPLOAD_FOLDER"
                          :fixed-upload-folder="true"
                        />
                      </div>

                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="method.name"
                          :languages="multilingualLangOptions"
                          title="名称多语言 name"
                          input-type="input"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>

                      <div class="tryon-model-field full">
                        <MultiLangEditor
                          :model="method.copyText"
                          :languages="multilingualLangOptions"
                          title="复制信息多语言 copyText"
                          input-type="textarea"
                          :rows="2"
                          :use-tabs="true"
                          tab-type="card"
                        />
                      </div>
                    </div>
                    <div class="tryon-model-actions">
                      <el-button size="small" @click="clonePaymentPreferredMethod(index)">复制</el-button>
                      <el-button size="small" type="danger" plain @click="removePaymentPreferredMethod(index)">删除</el-button>
                    </div>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
          </template>
          <!-- JSON多语言类型 -->
          <template v-else-if="isCurrencySymbolConfig(editForm)">
            <div class="json-editor">
              <div class="json-editor-tip">请手动输入各语言货币符号，可使用“一键补齐”补全槽位。</div>
              <MultiLangEditor
                :model="editJsonValue"
                :languages="multilingualLangOptions"
                title="货币符号多语言 currency_symbol"
                input-type="input"
              />
            </div>
          </template>
          <template v-else-if="isTryonTutorialTitleConfig(editForm)">
            <div class="json-editor tutorial-config-editor">
              <div class="json-editor-tip">教程弹窗标题（多语言）。支持一键翻译空白/覆盖翻译，建议保持简洁。</div>
              <MultiLangEditor
                :model="editJsonValue"
                :languages="multilingualLangOptions"
                title="教程标题多语言 tryon_tutorial_title"
                input-type="input"
                :use-tabs="true"
                tab-type="card"
              />
            </div>
          </template>
          <template v-else-if="isTryonTutorialContentConfig(editForm)">
            <div class="json-editor tutorial-config-editor">
              <el-alert
                type="info"
                show-icon
                :closable="false"
                title="教程正文支持富文本编辑；建议使用标题、列表和段落。一键翻译后请复核格式。"
                class="tutorial-editor-alert"
              />
              <MultiLangEditor
                :model="editJsonValue"
                :languages="multilingualLangOptions"
                title="教程正文多语言 tryon_tutorial_content"
                input-type="richtext"
                :use-tabs="true"
                tab-type="card"
                :rich-upload-folder="TUTORIAL_RICH_UPLOAD_FOLDER"
              />
            </div>
          </template>
          <template v-else-if="isJsonConfig(editForm)">
            <div class="json-editor">
              <MultiLangEditor
                :model="editJsonValue"
                :languages="multilingualLangOptions"
                title="JSON 多语言配置"
              />
            </div>
          </template>
          <!-- 数字类型 -->
          <template v-else-if="isNumberConfig(editForm)">
            <el-input-number v-model="editNumberValue" :min="0" />
          </template>
          <!-- 密钥类型 -->
          <template v-else-if="isSecretConfig(editForm)">
            <el-input v-model="editForm.configValue" type="password" show-password />
          </template>
          <!-- 普通文本 -->
          <template v-else>
            <el-input v-model="editForm.configValue" type="textarea" :rows="3" />
          </template>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="editForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="rechargeRateDialogVisible"
      :title="rechargeRateDialogTitle"
      width="980px"
      destroy-on-close
    >
      <div class="mb-3 flex flex-wrap items-center gap-2">
        <el-checkbox :model-value="isAllRechargeRowsSelected()" @change="toggleSelectAllRechargeRows">全选</el-checkbox>
        <el-input-number
          v-model="rechargeAdjustPercent"
          :precision="2"
          :step="0.5"
          placeholder="增减百分比"
          style="width: 140px"
        />
        <span class="text-xs text-gray-500">增加%（可负数），人民币价格换算后按比例增减</span>
      </div>

      <div class="mb-3 flex flex-wrap gap-2">
        <el-button type="primary" :loading="rechargeRateLoading" @click="refreshSelectedRechargeExchangeRates">更新汇率（选中）</el-button>
        <el-button type="success" @click="applySelectedRechargeRateToPrices">补齐（更新选中价格）</el-button>
        <el-button @click="fillRechargeRateLangSlots">补齐多语言槽位</el-button>
      </div>

      <div class="mb-2 text-xs text-gray-500">
        汇率来源：{{ rechargeExchangeRateSource || '-' }}
        <span v-if="rechargeExchangeRateFetchedAt">，更新时间：{{ rechargeExchangeRateFetchedAt }}</span>
      </div>

      <el-table :data="rechargeRateRows" border max-height="420px">
        <el-table-column label="选择" width="70">
          <template #default="scope">
            <el-checkbox v-model="scope.row.selected" />
          </template>
        </el-table-column>
        <el-table-column label="语言" width="180">
          <template #default="scope">
            {{ scope.row.name }} ({{ scope.row.code }})
          </template>
        </el-table-column>
        <el-table-column label="币种" width="140">
          <template #default="scope">
            <el-input v-model="scope.row.currency" maxlength="3" placeholder="USD" @change="normalizeRechargeRowCurrency(scope.row)" />
          </template>
        </el-table-column>
        <el-table-column label="汇率(CNY->币种)" width="190">
          <template #default="scope">
            <el-input-number
              v-model="scope.row.rate"
              :precision="6"
              :min="0"
              :step="0.001"
              controls-position="right"
              style="width: 100%"
            />
          </template>
        </el-table-column>
        <el-table-column label="目标价格(分)">
          <template #default="scope">
            <div class="flex items-center gap-2">
              <el-input-number
                v-model="scope.row.targetPriceFen"
                :precision="0"
                :step="1"
                :min="0"
                controls-position="right"
                style="width: 160px"
                @change="syncRechargeRowPriceToI18n(scope.row)"
              />
              <span class="text-xs text-gray-500">≈ {{ formatRechargeFen(scope.row.targetPriceFen) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button link type="primary" @click="refreshSingleRechargeRate(scope.row)">更新汇率</el-button>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="rechargeRateDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="modelLogDetailVisible" title="模型调用日志详情" width="1000px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="时间">{{ formatDateTime(modelLogDetail.CreatedAt) }}</el-descriptions-item>
        <el-descriptions-item label="耗时">{{ Number(modelLogDetail.durationMs || 0) }}ms</el-descriptions-item>
        <el-descriptions-item label="任务号">{{ modelLogDetail.taskNo || '-' }}</el-descriptions-item>
        <el-descriptions-item label="请求ID">{{ modelLogDetail.requestID || '-' }}</el-descriptions-item>
        <el-descriptions-item label="场景/行为">{{ `${modelLogDetail.sceneType || '-'} / ${modelLogDetail.behavior || '-'}` }}</el-descriptions-item>
        <el-descriptions-item label="阶段">{{ modelLogDetail.callStage || '-' }}</el-descriptions-item>
        <el-descriptions-item label="模型">{{ `${modelLogDetail.modelKey || '-'} (${modelLogDetail.modelUsage || '-'})` }}</el-descriptions-item>
        <el-descriptions-item label="Provider">{{ modelLogDetail.provider || '-' }}</el-descriptions-item>
        <el-descriptions-item label="token指纹">{{ modelLogDetail.tokenFingerprint || '-' }}</el-descriptions-item>
        <el-descriptions-item label="token槽位">{{ modelLogDetail.tokenSlot || '-' }}</el-descriptions-item>
        <el-descriptions-item label="调用地址" :span="2">{{ modelLogDetail.endpointURL || '-' }}</el-descriptions-item>
        <el-descriptions-item label="查询地址" :span="2">{{ modelLogDetail.queryURL || '-' }}</el-descriptions-item>
        <el-descriptions-item label="输入图" :span="2">{{ modelLogDetail.sourceImage || '-' }}</el-descriptions-item>
        <el-descriptions-item label="模板图" :span="2">{{ modelLogDetail.templateImage || '-' }}</el-descriptions-item>
        <el-descriptions-item label="结果图" :span="2">{{ modelLogDetail.resultImage || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ modelLogDetail.status || '-' }}</el-descriptions-item>
        <el-descriptions-item label="错误信息">{{ modelLogDetail.errorMessage || '-' }}</el-descriptions-item>
      </el-descriptions>

      <div class="model-log-image-grid">
        <div class="model-log-image-card">
          <div class="model-log-payload-title">输入图 sourceImage</div>
          <el-image
            v-if="hasPreviewableImage(modelLogDetail.sourceImage)"
            :src="modelLogDetail.sourceImage"
            :preview-src-list="buildImagePreviewList(modelLogDetail.sourceImage)"
            fit="cover"
            preview-teleported
            class="model-log-image"
          />
          <div v-else class="model-log-image-empty">无可预览图片</div>
        </div>
        <div class="model-log-image-card">
          <div class="model-log-payload-title">模板图 templateImage</div>
          <el-image
            v-if="hasPreviewableImage(modelLogDetail.templateImage)"
            :src="modelLogDetail.templateImage"
            :preview-src-list="buildImagePreviewList(modelLogDetail.templateImage)"
            fit="cover"
            preview-teleported
            class="model-log-image"
          />
          <div v-else class="model-log-image-empty">无可预览图片</div>
        </div>
        <div class="model-log-image-card">
          <div class="model-log-payload-title">结果图 resultImage</div>
          <el-image
            v-if="hasPreviewableImage(modelLogDetail.resultImage)"
            :src="modelLogDetail.resultImage"
            :preview-src-list="buildImagePreviewList(modelLogDetail.resultImage)"
            fit="cover"
            preview-teleported
            class="model-log-image"
          />
          <div v-else class="model-log-image-empty">无可预览图片</div>
        </div>
      </div>

      <div class="model-log-payload-grid">
        <div class="model-log-payload-item">
          <div class="model-log-payload-title">协同模型 relatedModels</div>
          <el-input :model-value="formatPayloadText(modelLogDetail.relatedModels)" type="textarea" :rows="6" readonly />
        </div>
        <div class="model-log-payload-item">
          <div class="model-log-payload-title">请求 payload</div>
          <el-input :model-value="formatPayloadText(modelLogDetail.requestPayload)" type="textarea" :rows="12" readonly />
        </div>
        <div class="model-log-payload-item">
          <div class="model-log-payload-title">响应 payload</div>
          <el-input :model-value="formatPayloadText(modelLogDetail.responsePayload)" type="textarea" :rows="12" readonly />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getSysConfigList, updateSysConfig, getAliyunTryonQuotaEstimate, getModelCallLogList } from '@/api/client/sysConfig'
import { getEnabledLanguages, getLanguageList } from '@/api/client/language'
import { ElMessage } from 'element-plus'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
import SelectImage from '@/components/selectImage/selectImage.vue'
import { validateRichTextI18nStructure } from '@/utils/richtext-i18n'
import {
  DEFAULT_LANG_CURRENCY_MAP,
  fetchExchangeRates,
  getExchangeRateSnapshot,
  calcConvertedFenFromCny,
  normalizePriceI18nMap,
} from '@/utils/exchange-rate'

// 配置组定义
const groupList = [
  { value: 'system', label: '系统设置', type: '' },
  { value: 'maintenance', label: '维护模式', type: 'danger' },
  { value: 'security', label: '安全设置', type: 'danger' },
  { value: 'auth', label: '认证设置', type: 'warning' },
  { value: 'payment', label: '支付设置', type: 'success' },
  { value: 'points', label: '积分设置', type: 'info' },
  { value: 'tryon', label: '试衣设置', type: 'warning' },
  { value: 'order', label: '订单设置', type: '' },
  { value: 'display', label: '显示设置', type: 'warning' },
  { value: 'announcement', label: '公告设置', type: 'info' },
  { value: 'invite', label: '邀请设置', type: 'success' },
]

const groupTagType = (group) => {
  const found = groupList.find(g => g.value === group)
  return found ? found.type : 'info'
}

const groupLabel = (group) => {
  const found = groupList.find(g => g.value === group)
  return found ? found.label : group
}

// 布尔键列表
const booleanKeys = [
  'maintenance_enabled', 'maintenance_popup_enabled', 'maintenance_home_btn_enabled',
  'phone_login_enabled', 'username_login_enabled', 'password_change_enabled',
  'sign_in_enabled', 'announcement_enabled',
  'tryon_append_parsing_failed_tip', 'tryon_append_refiner_failed_tip',
  'payment_auto_enabled', 'payment_manual_qrcode_enabled', 'payment_manual_contact_enabled',
  'payment_wechat_enabled', 'payment_alipay_enabled', 'payment_bank_cn_enabled',
  'payment_bank_us_enabled', 'payment_bank_mn_enabled', 'payment_paypal_enabled'
]
const isBooleanConfig = (row) => booleanKeys.includes(row.configKey)

// 颜色键列表
const colorKeys = ['announcement_text_color', 'payment_tip_text_color', 'invite_share_link_tip_text_color']
const isColorConfig = (row) => colorKeys.includes(row.configKey)

const isTryonModelsConfig = (row) => row?.configKey === 'tryon_models'
const isTryonRechargePlansConfig = (row) => row?.configKey === 'tryon_recharge_plans'
const isPaymentManualMethodsConfig = (row) => row?.configKey === 'payment_manual_methods'
const isPaymentUniPreferredMethodsConfig = (row) => row?.configKey === 'payment_uni_preferred_methods'
const isCurrencySymbolConfig = (row) => row?.configKey === 'currency_symbol'
const isTryonTutorialTitleConfig = (row) => row?.configKey === 'tryon_tutorial_title'
const isTryonTutorialContentConfig = (row) => row?.configKey === 'tryon_tutorial_content'
const isTryonTutorialConfig = (row) => isTryonTutorialTitleConfig(row) || isTryonTutorialContentConfig(row)
const isMaintenanceBgImageConfig = (row) => row?.configKey === 'maintenance_bg_image'

const PAYMENT_MANUAL_METHOD_UPLOAD_FOLDER = 'cloth-on/web-else/pay-method'
const PAYMENT_PREFERRED_METHOD_UPLOAD_FOLDER = 'cloth-on/web-else/hope-pay'
const MAINTENANCE_BG_UPLOAD_FOLDER = 'cloth-on/web-else'
const TUTORIAL_RICH_UPLOAD_FOLDER = 'cloth-on/web-else/tutorial'

// 密钥键列表
const secretKeys = []
const isSecretConfig = (row) => secretKeys.includes(row.configKey)

const maskSecretValue = (value) => {
  if (!value) {
    return '-'
  }
  if (value.length <= 8) {
    return '*'.repeat(value.length)
  }
  return `${value.slice(0, 4)}${'*'.repeat(Math.max(value.length - 8, 4))}${value.slice(-4)}`
}

// JSON键列表（多语言配置）
const jsonKeys = [
  'payment_tip_text', 'maintenance_popup_title', 'maintenance_popup_content',
  'announcement_content', 'maintenance_message',
  'username_regex_tip', 'password_regex_tip',
  'tryon_parsing_failed_tip_text', 'tryon_refiner_failed_tip_text',
  'app_name', 'invite_share_link_tip_text', 'currency_symbol',
  'tryon_tutorial_title', 'tryon_tutorial_content'
]
const isJsonConfig = (row) => jsonKeys.includes(row.configKey)

const formatCurrencySymbolPreview = (value) => {
  if (value === undefined || value === null) return '-'
  if (typeof value === 'object' && !Array.isArray(value)) {
    return displayI18nText(value, '-')
  }

  const raw = String(value || '').trim()
  if (!raw) return '-'
  if (raw.charAt(0) === '{') {
    try {
      const parsed = JSON.parse(raw)
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        return displayI18nText(parsed, '-')
      }
    } catch {
      return raw
    }
  }
  return raw
}

const parseMultilingualConfigValue = (value) => {
  if (value && typeof value === 'object' && !Array.isArray(value)) {
    return value
  }

  const raw = String(value || '').trim()
  if (!raw) {
    return {}
  }

  if (raw.charAt(0) === '{') {
    try {
      const parsed = JSON.parse(raw)
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        return parsed
      }
    } catch {
      return {}
    }
  }

  return {}
}

const stripHtmlTags = (value) => {
  return String(value || '')
    .replace(/<style[\s\S]*?<\/style>/gi, ' ')
    .replace(/<script[\s\S]*?<\/script>/gi, ' ')
    .replace(/<[^>]+>/g, ' ')
    .replace(/&nbsp;/gi, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}

const formatTutorialPreview = (row, showFull = false) => {
  const map = parseMultilingualConfigValue(row?.configValue)
  let text = displayI18nText(map, String(row?.configValue || '').trim())
  if (isTryonTutorialContentConfig(row)) {
    text = stripHtmlTags(text)
  }
  const normalized = String(text || '').trim() || '-'
  if (showFull || normalized.length <= 80) {
    return normalized
  }
  return `${normalized.slice(0, 80)}...`
}

// 数字键列表
const numberKeys = [
  'captcha_expiry_seconds', 'captcha_rate_limit', 'register_ip_limit',
  'login_fail_max', 'login_fail_wait_seconds', 'points_exchange_rate',
  'order_close_minutes', 'presale_home_count', 'announcement_speed',
  'invite_reward_points', 'payment_tip_text_size',
  'tryon_guest_init_points', 'tryon_register_reward_points', 'tryon_invite_register_reward_points', 'tryon_cost_points'
]
const isNumberConfig = (row) => numberKeys.includes(row.configKey)

const padDatePart = (value) => String(value).padStart(2, '0')

const formatDateRangeValue = (date) => {
  return `${date.getFullYear()}-${padDatePart(date.getMonth() + 1)}-${padDatePart(date.getDate())} ${padDatePart(date.getHours())}:${padDatePart(date.getMinutes())}:${padDatePart(date.getSeconds())}`
}

const buildDefaultModelLogDateRange = () => {
  const end = new Date()
  const start = new Date(end.getTime() - 24 * 60 * 60 * 1000)
  return [formatDateRangeValue(start), formatDateRangeValue(end)]
}

const normalizeLogImageURL = (value) => {
  const url = String(value || '').trim()
  if (!url || url === '-') {
    return ''
  }
  return url
}

const hasPreviewableImage = (value) => {
  return !!normalizeLogImageURL(value)
}

const buildImagePreviewList = (value) => {
  const url = normalizeLogImageURL(value)
  return url ? [url] : []
}

const activeGroup = ref('')
const searchInfo = ref({
  page: 1,
  pageSize: 50,
  configGroup: '',
  configKey: '',
})

const tableData = ref([])
const total = ref(0)
const modelCallLogList = ref([])
const modelLogTotal = ref(0)
const modelLogLoading = ref(false)
const modelLogDateRange = ref(buildDefaultModelLogDateRange())
const modelLogSearch = reactive({
  page: 1,
  pageSize: 20,
  startCreatedAt: '',
  endCreatedAt: '',
  onlyFailed: false,
  sceneType: '',
  modelUsage: '',
  status: '',
  keyword: '',
})
const modelLogDetailVisible = ref(false)
const modelLogDetail = ref({})
const editVisible = ref(false)
const editForm = ref({})
const editBoolValue = ref(false)
const editNumberValue = ref(0)
const editJsonValue = reactive({})
const newJsonLang = ref('')
const newJsonVal = ref('')
const tryonModels = ref([])
const rechargePlans = ref([])
const rechargeRateDialogVisible = ref(false)
const rechargeRatePlanIndex = ref(-1)
const rechargeRateRows = ref([])
const rechargeRateLoading = ref(false)
const rechargeAdjustPercent = ref(0)
const rechargeExchangeRateSource = ref('')
const rechargeExchangeRateFetchedAt = ref('')
const paymentManualMethods = ref([])
const paymentPreferredMethods = ref([])
const tryonModelActiveTab = ref('')
const aliyunQuotaState = reactive({})

const editDialogWidth = computed(() => {
  if (isTryonModelsConfig(editForm.value)) {
    return '1200px'
  }
  if (isTryonTutorialContentConfig(editForm.value)) {
    return '1200px'
  }
  if (isTryonTutorialTitleConfig(editForm.value)) {
    return '900px'
  }
  return '600px'
})

const rechargeRateDialogTitle = computed(() => {
  const index = rechargeRatePlanIndex.value
  const plan = rechargePlans.value[index]
  if (!plan) {
    return '充值套餐价格汇率换算'
  }
  const pointText = displayI18nText(plan.points, '').trim()
  const coinText = displayI18nText(plan.coinLabel, '').trim()
  const suffix = [pointText, coinText].filter(Boolean).join(' ')
  return suffix ? `充值套餐价格汇率换算 - ${suffix}` : `充值套餐价格汇率换算 #${index + 1}`
})

const DEFAULT_MULTILINGUAL_LANGS = [
  { code: 'zh', name: '中文' },
  { code: 'en', name: '英文' },
  { code: 'mn', name: '蒙文' },
  { code: 'zh-TW', name: '繁体' },
  { code: 'th', name: '泰语' },
  { code: 'hi', name: '印地语' },
  { code: 'id', name: '印尼语' },
  { code: 'vi', name: '越南语' },
  { code: 'ar', name: '阿拉伯语' },
  { code: 'ja', name: '日语' },
  { code: 'ko', name: '韩语' },
  { code: 'ms', name: '马来语' },
]

const defaultLangNameMap = DEFAULT_MULTILINGUAL_LANGS.reduce((acc, item) => {
  acc[item.code] = item.name
  return acc
}, {})

const DISPLAY_LANG_PRIORITY = [
  'zh',
  'en',
  'mn',
  'zh-TW',
  'th',
  'hi',
  'id',
  'vi',
  'ar',
  'ja',
  'ko',
  'ms',
]
const displayLangRankMap = new Map(DISPLAY_LANG_PRIORITY.map((code, index) => [code, index]))

const normalizeLangCode = (value) => String(value || '').trim()

const uniqueLangItems = (items = []) => {
  const map = new Map()
  items.forEach((item) => {
    const code = normalizeLangCode(item?.code)
    if (!code || map.has(code)) return
    const name = String(item?.name || item?.label || defaultLangNameMap[code] || code).trim()
    map.set(code, {
      code,
      name,
    })
  })
  return Array.from(map.values())
}

const toLangDisplayName = (code, name) => {
  const normalizedCode = normalizeLangCode(code)
  const normalizedName = String(name || defaultLangNameMap[normalizedCode] || normalizedCode).trim()
  if (!normalizedCode) return normalizedName
  if (!normalizedName) return normalizedCode
  return normalizedName.includes(normalizedCode) ? normalizedName : `${normalizedName} ${normalizedCode}`
}

let multilingualCodes = DEFAULT_MULTILINGUAL_LANGS.map(item => item.code)
const multilingualLangOptions = ref(
  DEFAULT_MULTILINGUAL_LANGS.map(item => ({
    code: item.code,
    name: toLangDisplayName(item.code, item.name),
  }))
)

const applyMultilingualLanguages = (items = []) => {
  const normalized = uniqueLangItems(items)
  const source = normalized.length > 0
    ? normalized
    : DEFAULT_MULTILINGUAL_LANGS.map(item => ({ code: item.code, name: item.name }))

  const orderedSource = [...source].sort((a, b) => {
    const rankA = displayLangRankMap.has(a.code) ? displayLangRankMap.get(a.code) : Number.MAX_SAFE_INTEGER
    const rankB = displayLangRankMap.has(b.code) ? displayLangRankMap.get(b.code) : Number.MAX_SAFE_INTEGER
    if (rankA !== rankB) {
      return rankA - rankB
    }
    return String(a.code).localeCompare(String(b.code))
  })

  multilingualCodes = orderedSource.map(item => item.code)
  multilingualLangOptions.value = orderedSource.map(item => ({
    code: item.code,
    name: toLangDisplayName(item.code, item.name),
  }))
}

const loadMultilingualLanguages = async () => {
  try {
    const [enabledRes, listRes] = await Promise.all([
      getEnabledLanguages().catch(() => null),
      getLanguageList({ page: 1, pageSize: 500 }).catch(() => null),
    ])

    const enabledList = enabledRes?.code === 0 && Array.isArray(enabledRes?.data)
      ? enabledRes.data
      : []
    const listPayload = listRes?.code === 0
      ? (Array.isArray(listRes?.data?.list)
          ? listRes.data.list
          : (Array.isArray(listRes?.data) ? listRes.data : []))
      : []

    applyMultilingualLanguages(listPayload.length ? listPayload : enabledList)
  } catch {
    applyMultilingualLanguages(DEFAULT_MULTILINGUAL_LANGS)
  }
}

const buildMultilingualObject = (value, fallback = '') => {
  const source = value && typeof value === 'object' && !Array.isArray(value)
    ? value
    : {}
  const baseText = value && typeof value === 'object' && !Array.isArray(value)
    ? ''
    : (value === undefined || value === null ? fallback : String(value))

  const fallbackText = String(
    source.zh ?? source.en ?? source.mn ?? source['zh-TW'] ?? source.th ?? source.hi ?? source.id ?? baseText ?? fallback
  )

  return multilingualCodes.reduce((acc, code) => {
    acc[code] = String(source[code] ?? fallbackText)
    return acc
  }, {})
}

const displayI18nText = (value, fallback = '') => {
  const source = value && typeof value === 'object' && !Array.isArray(value)
    ? value
    : {}
  const normalized = buildMultilingualObject(value, fallback)

  const candidateCodes = Array.from(new Set([
    ...DISPLAY_LANG_PRIORITY,
    ...multilingualCodes,
    ...Object.keys(source),
  ]))

  for (const code of candidateCodes) {
    const text = String(source[code] ?? normalized[code] ?? '').trim()
    if (text) return text
  }

  for (const raw of Object.values(source)) {
    const text = String(raw || '').trim()
    if (text) return text
  }

  return fallback
}

const createTryonModelUid = () => `tryon_model_${Date.now()}_${Math.random().toString(16).slice(2, 8)}`

const toBool = (value, fallback = true) => {
  if (typeof value === 'boolean') {
    return value
  }
  if (value === undefined || value === null || value === '') {
    return fallback
  }
  const text = String(value).trim().toLowerCase()
  if (['true', '1', 'yes', 'on'].includes(text)) {
    return true
  }
  if (['false', '0', 'no', 'off'].includes(text)) {
    return false
  }
  return fallback
}

const toInt = (value, fallback = 0) => {
  const numberValue = Number.parseInt(value, 10)
  return Number.isFinite(numberValue) ? numberValue : fallback
}

const parseOptionalBeautifyDegree = (value) => {
  if (value === '' || value === null || value === undefined) {
    return undefined
  }
  const numberValue = Number(value)
  if (!Number.isFinite(numberValue) || numberValue <= 0) {
    return undefined
  }
  if (numberValue > 1.5) {
    return undefined
  }
  return Number(numberValue.toFixed(3))
}

const toStringArray = (value, fallback = []) => {
  if (Array.isArray(value)) {
    return Array.from(new Set(value.map(item => String(item).trim()).filter(Boolean)))
  }
  return [...fallback]
}

const splitTokenBackupInput = (value) => {
  return String(value || '')
    .split(/\r?\n|,|;/g)
    .map(item => item.trim())
    .filter(Boolean)
}

const normalizeTokenBackups = (value) => {
  const result = []
  const appendToken = (token) => {
    const normalized = String(token || '').trim()
    if (!normalized) return
    if (result.includes(normalized)) return
    result.push(normalized)
  }

  if (Array.isArray(value)) {
    value.forEach(appendToken)
    return result
  }

  if (value && typeof value === 'object') {
    return result
  }

  splitTokenBackupInput(value).forEach(appendToken)
  return result
}

const joinTokenBackups = (value) => {
  return normalizeTokenBackups(value).join('\n')
}

const splitTokenQuotaInput = (value) => {
  return String(value || '')
    .split(/\r?\n/g)
    .map(item => item.trim())
    .filter(Boolean)
}

const normalizeTokenQuotas = (value) => {
  const result = []
  const quotaByToken = new Map()

  const appendQuota = (token, freeQuotaTotal) => {
    const normalizedToken = String(token || '').trim()
    if (!normalizedToken) {
      return
    }
    const normalizedQuota = Math.max(0, toInt(freeQuotaTotal, 0))
    if (quotaByToken.has(normalizedToken)) {
      quotaByToken.set(normalizedToken, normalizedQuota)
      return
    }
    quotaByToken.set(normalizedToken, normalizedQuota)
    result.push(normalizedToken)
  }

  if (Array.isArray(value)) {
    value.forEach((item) => {
      if (!item || typeof item !== 'object') {
        return
      }
      appendQuota(item.token, item.freeQuotaTotal)
    })
  } else {
    splitTokenQuotaInput(value).forEach((line) => {
      const parts = line.split('|')
      const token = String(parts[0] || '').trim()
      const quota = parts.length > 1 ? parts.slice(1).join('|') : '0'
      appendQuota(token, quota)
    })
  }

  return result.map(token => ({
    token,
    freeQuotaTotal: quotaByToken.get(token) || 0,
  }))
}

const joinTokenQuotas = (value) => {
  return normalizeTokenQuotas(value)
    .map(item => `${item.token}|${item.freeQuotaTotal}`)
    .join('\n')
}

const collectConfiguredTokens = (model = {}) => {
  const uniqueTokens = []
  const tokenSet = new Set()
  const appendToken = (value) => {
    const normalized = String(value || '').trim()
    if (!normalized || tokenSet.has(normalized)) {
      return
    }
    tokenSet.add(normalized)
    uniqueTokens.push(normalized)
  }

  appendToken(model.token)
  appendToken(model.providerToken)
  appendToken(model.refinerToken)
  appendToken(model.beautifyToken)
  normalizeTokenBackups(model.tokenBackupText || model.tokenBackups || model.backupTokens).forEach(appendToken)
  return uniqueTokens
}

const getTokenQuotaRows = (model = {}) => {
  const configuredTokens = collectConfiguredTokens(model)
  const defaultQuota = Math.max(0, toInt(model.freeQuotaTotal, 0))
  const quotaMap = new Map(
    normalizeTokenQuotas(model.tokenQuotaText || model.tokenQuotas)
      .map(item => [item.token, Math.max(0, toInt(item.freeQuotaTotal, 0))])
  )

  return configuredTokens.map(token => ({
    token,
    freeQuotaTotal: quotaMap.has(token) ? quotaMap.get(token) : defaultQuota,
  }))
}

const hasTokenQuotaRows = (model = {}) => getTokenQuotaRows(model).length > 0

const updateTokenQuotaValue = (model = {}, token = '', value = 0) => {
  if (!model || typeof model !== 'object') {
    return
  }
  const normalizedToken = String(token || '').trim()
  if (!normalizedToken) {
    return
  }

  const rows = getTokenQuotaRows(model)
  if (!rows.some(item => item.token === normalizedToken)) {
    return
  }

  const normalizedValue = Math.max(0, toInt(value, 0))
  const nextRows = rows.map((item) => {
    if (item.token !== normalizedToken) {
      return item
    }
    return {
      ...item,
      freeQuotaTotal: normalizedValue,
    }
  })
  model.tokenQuotaText = joinTokenQuotas(nextRows)
}

const normalizeGender = (value, fallback = 'woman') => {
  const text = String(value || '').trim().toLowerCase()
  if (text === 'woman' || text === 'man') {
    return text
  }
  return fallback
}

const inferLegacyModelUsage = (model = {}) => {
  const key = String(model.key || model.modelKey || '').trim().toLowerCase()
  const modelName = String(model.model || model.refinerModel || model.beautifyModel || '').trim().toLowerCase()
  const hint = `${key} ${modelName}`

  if (toBool(model.supportsBeautify, false)) {
    return 'beautify'
  }
  if (hint.includes('refiner')) {
    return 'refiner'
  }
  if (hint.includes('parsing') || hint.includes('takeoff')) {
    return 'parsing'
  }
  return 'tryon'
}

const normalizeModelUsage = (value, fallback = 'tryon', model = {}) => {
  const text = String(value || '').trim().toLowerCase()
  const inferred = inferLegacyModelUsage(model)

  if (text === 'refiner' || text === 'parsing' || text === 'beautify') {
    return text
  }
  if (text === 'tryon') {
    if (inferred === 'refiner' || inferred === 'parsing' || inferred === 'beautify') {
      return inferred
    }
    return 'tryon'
  }

  if (inferred === 'refiner' || inferred === 'parsing' || inferred === 'beautify') {
    return inferred
  }

  const fallbackText = String(fallback || '').trim().toLowerCase()
  if (['tryon', 'refiner', 'parsing', 'beautify'].includes(fallbackText)) {
    return fallbackText
  }
  return 'tryon'
}

const isTryonUsageModel = (model = {}) => normalizeModelUsage(model.modelUsage, 'tryon', model) === 'tryon'
const isRefinerUsageModel = (model = {}) => normalizeModelUsage(model.modelUsage, 'tryon', model) === 'refiner'
const isParsingUsageModel = (model = {}) => normalizeModelUsage(model.modelUsage, 'tryon', model) === 'parsing'
const isBeautifyUsageModel = (model = {}) => normalizeModelUsage(model.modelUsage, 'tryon', model) === 'beautify'

const usageLabelMap = {
  tryon: '试衣模型',
  refiner: '图片精修',
  parsing: '分割模型',
  beautify: '美肤模型',
}

const usageTagTypeMap = {
  tryon: 'primary',
  refiner: 'warning',
  parsing: 'info',
  beautify: 'success',
}

const getModelUsageLabel = (model = {}) => {
  const usage = normalizeModelUsage(model.modelUsage, 'tryon', model)
  return usageLabelMap[usage] || '试衣模型'
}

const getModelUsageTagType = (model = {}) => {
  const usage = normalizeModelUsage(model.modelUsage, 'tryon', model)
  return usageTagTypeMap[usage] || 'primary'
}

const beautifyModelKeyOptions = computed(() => {
  return Array.from(new Set(
    tryonModels.value
      .filter(model => isBeautifyUsageModel(model))
      .map(model => String(model.key || '').trim())
      .filter(Boolean)
  ))
})

const refinerModelKeyOptions = computed(() => {
  return Array.from(new Set(
    tryonModels.value
      .filter(model => isRefinerUsageModel(model))
      .map(model => String(model.key || '').trim())
      .filter(Boolean)
  ))
})

const parsingModelKeyOptions = computed(() => {
  return Array.from(new Set(
    tryonModels.value
      .filter(model => isParsingUsageModel(model))
      .map(model => String(model.key || '').trim())
      .filter(Boolean)
  ))
})

const isAliyunModelForQuota = (model = {}) => {
  if (!isTryonUsageModel(model) && !isRefinerUsageModel(model) && !isParsingUsageModel(model)) {
    return false
  }
  const key = String(model.key || '').trim().toLowerCase()
  const provider = String(model.provider || '').trim().toLowerCase()
  const modelName = String(model.model || '').trim().toLowerCase()
  return provider.includes('aliyun') || provider.includes('dashscope') || key.includes('aliyun') || modelName.startsWith('aitryon')
}

const inferSupportsRefiner = (model = {}) => {
  if (!isTryonUsageModel(model)) {
    return false
  }
  if (!isAliyunModelForQuota(model)) {
    return false
  }
  const modelName = String(model.model || '').trim().toLowerCase()
  const key = String(model.key || '').trim().toLowerCase()
  return modelName === 'aitryon' || modelName === 'aitryon-plus' || key.includes('aliyun_aitryon') || key.includes('aliyun_aitryon_plus')
}

const quotaStateKey = (model = {}) => String(model.__uid || model.key || '').trim()

const ensureQuotaState = (model = {}) => {
  const key = quotaStateKey(model)
  if (!key) return null
  if (!aliyunQuotaState[key]) {
    aliyunQuotaState[key] = {
      loading: false,
      data: null,
      error: '',
    }
  }
  return aliyunQuotaState[key]
}

const getAliyunQuotaItem = (model = {}) => {
  const state = ensureQuotaState(model)
  return state?.data || null
}

const getAliyunQuotaError = (model = {}) => {
  const state = ensureQuotaState(model)
  return state?.error || ''
}

const isAliyunQuotaLoading = (model = {}) => {
  const state = ensureQuotaState(model)
  return !!state?.loading
}

const formatQuotaRefreshTime = (value) => {
  if (!value) return '未刷新'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return '刷新时间未知'
  }
  return `刷新于 ${date.toLocaleString()}`
}

const refreshAliyunQuota = async (model = {}, force = false) => {
  if (!isAliyunModelForQuota(model)) {
    return
  }

  const state = ensureQuotaState(model)
  if (!state) return

  const modelKey = String(model.key || '').trim()
  const modelUsage = normalizeModelUsage(model.modelUsage, 'tryon', model)
  if (!modelKey) {
    state.error = '请先填写模型 key 后再刷新额度估算'
    return
  }

  if (!force && state.data) {
    return
  }

  state.loading = true
  state.error = ''
  try {
    const res = await getAliyunTryonQuotaEstimate({ modelKey })
    if (res.code !== 0) {
      state.error = res.msg || '获取额度估算失败'
      state.data = null
      return
    }
    const list = Array.isArray(res?.data?.list) ? res.data.list : []
    const matched = list.find(item => {
      const itemKey = String(item.modelKey || '').trim()
      const itemUsage = normalizeModelUsage(item.modelUsage, 'tryon', item)
      return itemKey === modelKey && itemUsage === modelUsage
    })
    if (!matched) {
      state.error = '未找到该模型额度数据'
      state.data = null
      return
    }
    state.data = matched
  } catch (e) {
    state.error = e?.message || '获取额度估算失败'
    state.data = null
  } finally {
    state.loading = false
  }
}

const ensureTryonModelActiveTab = (preferredTab = '') => {
  const available = tryonModels.value
  if (!Array.isArray(available) || available.length === 0) {
    tryonModelActiveTab.value = ''
    return
  }

  const preferred = String(preferredTab || '').trim()
  if (preferred && available.some(item => item.__uid === preferred)) {
    tryonModelActiveTab.value = preferred
    return
  }

  const current = String(tryonModelActiveTab.value || '').trim()
  if (current && available.some(item => item.__uid === current)) {
    return
  }

  tryonModelActiveTab.value = available[0].__uid
}

const handleTryonModelTabChange = (tabName) => {
  const selected = String(tabName || '').trim()
  if (!selected) return
  const model = tryonModels.value.find(item => item.__uid === selected)
  if (model && isAliyunModelForQuota(model)) {
    refreshAliyunQuota(model, true)
  }
}

const normalizeI18nObject = (value) => {
  return buildMultilingualObject(value, '')
}

const normalizeRechargeI18nObject = (value, fallback = '') => {
  return buildMultilingualObject(value, fallback)
}

const createPaymentMethodUid = () => `payment_method_${Date.now()}_${Math.random().toString(16).slice(2, 8)}`

const defaultPaymentMethodNameByKey = (key = '') => {
  const normalizedKey = String(key || '').trim().toLowerCase()
  const defaults = {
    qrcode: { zh: '二维码支付', en: 'QR Payment', mn: 'QR төлбөр' },
    contact: { zh: '联系客服', en: 'Contact Support', mn: 'Хэрэглэгчийн дэмжлэг' },
    wechat: { zh: '微信支付', en: 'WeChat Pay', mn: 'WeChat Pay' },
    alipay: { zh: '支付宝', en: 'Alipay', mn: 'Alipay' },
    bank_card_cn: { zh: '银行卡(国内)', en: 'Bank Card (CN)', mn: 'Банкны карт (CN)' },
    bank_card_us: { zh: '银行卡(美国)', en: 'Bank Card (US)', mn: 'Банкны карт (US)' },
    bank_card_mn: { zh: '银行卡(蒙古)', en: 'Bank Card (MN)', mn: 'Банкны карт (MN)' },
    paypal: { zh: 'PayPal', en: 'PayPal', mn: 'PayPal' },
  }
  return normalizeRechargeI18nObject(defaults[normalizedKey] || {}, normalizedKey || 'Payment')
}

const createDefaultPaymentMethod = () => ({
  __uid: createPaymentMethodUid(),
  key: '',
  sort: 99,
  enabled: true,
  manual: true,
  image: '',
  name: defaultPaymentMethodNameByKey(''),
  copyText: normalizeRechargeI18nObject({}, ''),
})

const normalizePaymentMethod = (item = {}) => {
  const key = String(item.key || '').trim().toLowerCase()
  const normalizedName = defaultPaymentMethodNameByKey(key)
  const sourceName = normalizeRechargeI18nObject(item.name, key)
  const sourceCopyText = normalizeRechargeI18nObject(item.copyText, '')
  return {
    __uid: createPaymentMethodUid(),
    key,
    sort: toInt(item.sort, 99),
    enabled: toBool(item.enabled, true),
    manual: toBool(item.manual, key === 'qrcode' || key === 'contact'),
    image: String(item.image || item.externalPath || '').trim(),
    name: {
      ...normalizedName,
      ...sourceName,
    },
    copyText: sourceCopyText,
  }
}

const parsePaymentMethodsValue = (rawValue) => {
  try {
    const parsed = JSON.parse(rawValue || '[]')
    if (!Array.isArray(parsed)) return []
    return parsed.map(item => normalizePaymentMethod(item))
  } catch {
    return []
  }
}

const buildPaymentMethodsPayload = () => paymentManualMethods.value.map(item => ({
  key: String(item.key || '').trim().toLowerCase(),
  sort: toInt(item.sort, 99),
  enabled: !!item.enabled,
  manual: !!item.manual,
  image: String(item.image || '').trim(),
  name: normalizeRechargeI18nObject(item.name, String(item.key || '').trim().toLowerCase()),
  copyText: normalizeRechargeI18nObject(item.copyText, ''),
}))

const paymentMethodsCount = (rawValue) => parsePaymentMethodsValue(rawValue).length

const addPaymentMethod = () => {
  paymentManualMethods.value.push(createDefaultPaymentMethod())
}

const clonePaymentMethod = (index) => {
  const method = paymentManualMethods.value[index]
  if (!method) return
  const cloned = normalizePaymentMethod(method)
  if (cloned.key) {
    cloned.key = `${cloned.key}_copy`
  }
  paymentManualMethods.value.splice(index + 1, 0, cloned)
}

const removePaymentMethod = (index) => {
  paymentManualMethods.value.splice(index, 1)
}

const validatePaymentMethods = () => {
  const keys = new Set()
  for (let i = 0; i < paymentManualMethods.value.length; i++) {
    const method = paymentManualMethods.value[i]
    const key = String(method.key || '').trim().toLowerCase()
    if (!key) {
      ElMessage.warning(`第 ${i + 1} 个支付方式缺少 key`)
      return false
    }
    if (keys.has(key)) {
      ElMessage.warning(`支付方式 key 重复: ${key}`)
      return false
    }
    keys.add(key)
    const name = String(displayI18nText(method.name, '') || '').trim()
    if (!name) {
      ElMessage.warning(`第 ${i + 1} 个支付方式缺少名称`)
      return false
    }
  }
  return true
}

const createPaymentPreferredMethodUid = () => `payment_preferred_${Date.now()}_${Math.random().toString(16).slice(2, 8)}`

const createDefaultPaymentPreferredMethod = () => ({
  __uid: createPaymentPreferredMethodUid(),
  key: '',
  sort: 99,
  enabled: true,
  image: '',
  name: defaultPaymentMethodNameByKey(''),
  copyText: normalizeRechargeI18nObject({}, ''),
})

const normalizePaymentPreferredMethod = (item = {}) => {
  const key = String(item.key || '').trim().toLowerCase()
  const normalizedName = defaultPaymentMethodNameByKey(key)
  const sourceName = normalizeRechargeI18nObject(item.name, key)
  const sourceCopyText = normalizeRechargeI18nObject(item.copyText, '')
  return {
    __uid: createPaymentPreferredMethodUid(),
    key,
    sort: toInt(item.sort, 99),
    enabled: toBool(item.enabled, true),
    image: String(item.image || item.externalPath || '').trim(),
    name: {
      ...normalizedName,
      ...sourceName,
    },
    copyText: sourceCopyText,
  }
}

const parsePaymentPreferredMethodsValue = (rawValue) => {
  try {
    const parsed = JSON.parse(rawValue || '[]')
    if (!Array.isArray(parsed)) return []
    return parsed.map(item => normalizePaymentPreferredMethod(item))
  } catch {
    return []
  }
}

const buildPaymentPreferredMethodsPayload = () => paymentPreferredMethods.value.map(item => ({
  key: String(item.key || '').trim().toLowerCase(),
  sort: toInt(item.sort, 99),
  enabled: !!item.enabled,
  image: String(item.image || '').trim(),
  name: normalizeRechargeI18nObject(item.name, String(item.key || '').trim().toLowerCase()),
  copyText: normalizeRechargeI18nObject(item.copyText, ''),
}))

const paymentPreferredMethodsCount = (rawValue) => parsePaymentPreferredMethodsValue(rawValue).length

const addPaymentPreferredMethod = () => {
  paymentPreferredMethods.value.push(createDefaultPaymentPreferredMethod())
}

const clonePaymentPreferredMethod = (index) => {
  const method = paymentPreferredMethods.value[index]
  if (!method) return
  const cloned = normalizePaymentPreferredMethod(method)
  if (cloned.key) {
    cloned.key = `${cloned.key}_copy`
  }
  paymentPreferredMethods.value.splice(index + 1, 0, cloned)
}

const removePaymentPreferredMethod = (index) => {
  paymentPreferredMethods.value.splice(index, 1)
}

const validatePaymentPreferredMethods = () => {
  const keys = new Set()
  for (let i = 0; i < paymentPreferredMethods.value.length; i++) {
    const method = paymentPreferredMethods.value[i]
    const key = String(method.key || '').trim().toLowerCase()
    if (!key) {
      ElMessage.warning(`第 ${i + 1} 个期望支付方式缺少 key`)
      return false
    }
    if (key === 'qrcode' || key === 'contact') {
      ElMessage.warning(`第 ${i + 1} 个期望支付方式不能使用 ${key}`)
      return false
    }
    if (keys.has(key)) {
      ElMessage.warning(`期望支付方式 key 重复: ${key}`)
      return false
    }
    keys.add(key)
    const name = String(displayI18nText(method.name, '') || '').trim()
    if (!name) {
      ElMessage.warning(`第 ${i + 1} 个期望支付方式缺少名称`)
      return false
    }
  }
  return true
}

const createRechargePlanUid = () => `tryon_recharge_${Date.now()}_${Math.random().toString(16).slice(2, 8)}`

const parseRechargePriceFenFromValue = (value, fallback = 0, treatAsYuan = false) => {
  if (value === undefined || value === null || value === '') {
    return Math.max(0, toInt(fallback, 0))
  }

  if (typeof value === 'number') {
    if (!Number.isFinite(value)) {
      return Math.max(0, toInt(fallback, 0))
    }
    if (treatAsYuan || !Number.isInteger(value)) {
      return Math.max(0, Math.round(value * 100))
    }
    if (value > 0 && value < 100) {
      return Math.max(0, Math.round(value * 100))
    }
    return Math.max(0, Math.round(value))
  }

  const text = String(value).trim()
  if (!text) {
    return Math.max(0, toInt(fallback, 0))
  }

  const cleaned = text.replace(/[^0-9.]/g, '')
  if (!cleaned) {
    return Math.max(0, toInt(fallback, 0))
  }

  const num = Number(cleaned)
  if (!Number.isFinite(num) || num < 0) {
    return Math.max(0, toInt(fallback, 0))
  }

  if (treatAsYuan || cleaned.includes('.')) {
    return Math.max(0, Math.round(num * 100))
  }
  if (num > 0 && num < 100) {
    return Math.max(0, Math.round(num * 100))
  }
  return Math.max(0, Math.round(num))
}

const parseRechargeLegacyPriceMap = (value) => {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    return {}
  }
  const next = {}
  Object.entries(value).forEach(([code, raw]) => {
    const fen = parseRechargePriceFenFromValue(raw, -1, true)
    if (Number.isFinite(fen) && fen >= 0) {
      next[String(code)] = fen
    }
  })
  return next
}

const pickFirstRechargePriceFen = (priceMap, fallback = 0) => {
  if (!priceMap || typeof priceMap !== 'object') {
    return Math.max(0, toInt(fallback, 0))
  }

  for (const code of multilingualCodes) {
    const value = Number(priceMap[code])
    if (Number.isFinite(value) && value >= 0) {
      return Math.round(value)
    }
  }

  const values = Object.values(priceMap)
  for (const raw of values) {
    const value = Number(raw)
    if (Number.isFinite(value) && value >= 0) {
      return Math.round(value)
    }
  }

  return Math.max(0, toInt(fallback, 0))
}

const normalizeRechargePriceI18n = (value, fallbackFen = 0) => {
  const fallback = Math.max(0, toInt(fallbackFen, 0))
  const normalized = normalizePriceI18nMap(value)
  const next = {}

  Object.entries(normalized).forEach(([code, raw]) => {
    const priceFen = Number(raw)
    if (Number.isFinite(priceFen) && priceFen >= 0) {
      next[String(code)] = Math.round(priceFen)
    }
  })

  multilingualCodes.forEach((code) => {
    if (!Object.prototype.hasOwnProperty.call(next, code)) {
      next[code] = fallback
    }
  })

  if (!Object.keys(next).length) {
    multilingualCodes.forEach((code) => {
      next[code] = fallback
    })
  }

  return next
}

const resolveRechargePlanPriceFenByLang = (plan, langCode = 'zh') => {
  const normalizedLang = String(langCode || '').trim()
  const priceI18nMap = normalizeRechargePriceI18n(plan?.priceI18n, parseRechargePriceFenFromValue(plan?.price, 0, false))
  const candidates = []

  if (normalizedLang) {
    candidates.push(normalizedLang)
    if (normalizedLang.includes('-')) {
      candidates.push(normalizedLang.split('-')[0])
    }
  }
  candidates.push('zh', 'en', 'mn')

  for (const candidate of candidates) {
    if (!candidate || !Object.prototype.hasOwnProperty.call(priceI18nMap, candidate)) {
      continue
    }
    const value = Number(priceI18nMap[candidate])
    if (Number.isFinite(value) && value >= 0) {
      return Math.round(value)
    }
  }

  return parseRechargePriceFenFromValue(plan?.price, 0, false)
}

const formatRechargePlanPricePreview = (plan) => {
  const fen = resolveRechargePlanPriceFenByLang(plan, 'zh')
  return `${(fen / 100).toFixed(2)}`
}

const createDefaultRechargePlan = () => ({
  __uid: createRechargePlanUid(),
  points: normalizeRechargeI18nObject({ zh: '50', en: '50', mn: '50' }, '50'),
  coinLabel: normalizeRechargeI18nObject({ zh: '试衣币', en: 'Try-on Coins', mn: 'Туршилтын зоос' }, 'Try-on Coins'),
  price: 990,
  priceI18n: normalizeRechargePriceI18n({ zh: 990, en: 990, mn: 990 }, 990),
})

const normalizeRechargePlan = (item = {}) => {
  const explicitPriceI18n = normalizePriceI18nMap(item.priceI18n)
  const legacyPriceI18n = parseRechargeLegacyPriceMap(item.price)
  const parsedBasePrice = parseRechargePriceFenFromValue(item.price, 0, false)
  const fallbackBasePrice = parsedBasePrice > 0
    ? parsedBasePrice
    : pickFirstRechargePriceFen(explicitPriceI18n, pickFirstRechargePriceFen(legacyPriceI18n, 990))
  const sourcePriceI18n = Object.keys(explicitPriceI18n).length > 0 ? explicitPriceI18n : legacyPriceI18n

  return {
    __uid: createRechargePlanUid(),
    points: normalizeRechargeI18nObject(item.points, '50'),
    coinLabel: normalizeRechargeI18nObject(item.coinLabel || item.label, '试衣币'),
    price: fallbackBasePrice,
    priceI18n: normalizeRechargePriceI18n(sourcePriceI18n, fallbackBasePrice),
  }
}

const parseRechargePlansValue = (rawValue) => {
  try {
    const parsed = JSON.parse(rawValue || '[]')
    if (!Array.isArray(parsed)) return []
    return parsed.map(item => normalizeRechargePlan(item))
  } catch {
    return []
  }
}

const buildRechargePlansPayload = () => rechargePlans.value.map(item => {
  const basePriceFen = parseRechargePriceFenFromValue(item.price, 0, false)
  const priceI18n = normalizeRechargePriceI18n(item.priceI18n, basePriceFen)
  return {
    points: normalizeRechargeI18nObject(item.points, '0'),
    coinLabel: normalizeRechargeI18nObject(item.coinLabel, ''),
    price: basePriceFen,
    priceI18n,
  }
})

const rechargePlansCount = (rawValue) => parseRechargePlansValue(rawValue).length

const addRechargePlan = () => {
  rechargePlans.value.push(createDefaultRechargePlan())
}

const cloneRechargePlan = (index) => {
  const item = rechargePlans.value[index]
  if (!item) return
  rechargePlans.value.splice(index + 1, 0, normalizeRechargePlan(item))
}

const removeRechargePlan = (index) => {
  if (rechargeRatePlanIndex.value === index) {
    rechargeRateDialogVisible.value = false
    rechargeRatePlanIndex.value = -1
    rechargeRateRows.value = []
  } else if (rechargeRatePlanIndex.value > index) {
    rechargeRatePlanIndex.value -= 1
  }
  rechargePlans.value.splice(index, 1)
}

const normalizeRechargeCurrencyCode = (value) => {
  return String(value || '').trim().toUpperCase().slice(0, 3)
}

const formatRechargeFen = (fen) => {
  const value = Number(fen)
  if (!Number.isFinite(value)) return '0.00'
  return (value / 100).toFixed(2)
}

const ensureRechargePlanLangSlots = (index, silent = false) => {
  const plan = rechargePlans.value[index]
  if (!plan) return false

  const langCodes = Array.from(new Set(
    multilingualCodes
      .map(code => normalizeLangCode(code))
      .filter(Boolean)
  ))
  if (!langCodes.length) {
    if (!silent) {
      ElMessage.warning('请先在语言管理中配置语言')
    }
    return false
  }

  const basePriceFen = parseRechargePriceFenFromValue(plan.price, 0, false)
  plan.points = normalizeRechargeI18nObject(plan.points, displayI18nText(plan.points, ''))
  plan.coinLabel = normalizeRechargeI18nObject(plan.coinLabel, displayI18nText(plan.coinLabel, ''))
  plan.price = basePriceFen
  plan.priceI18n = normalizeRechargePriceI18n(plan.priceI18n, basePriceFen)
  return true
}

const buildRechargeRateRows = () => {
  const index = rechargeRatePlanIndex.value
  if (!ensureRechargePlanLangSlots(index, true)) {
    rechargeRateRows.value = []
    return
  }

  const plan = rechargePlans.value[index]
  const previousMap = {}
  rechargeRateRows.value.forEach((row) => {
    previousMap[row.code] = row
  })

  const basePriceFen = parseRechargePriceFenFromValue(plan.price, 0, false)
  const sourceLangs = Array.isArray(multilingualLangOptions.value) && multilingualLangOptions.value.length
    ? multilingualLangOptions.value
    : multilingualCodes.map(code => ({ code, name: toLangDisplayName(code, '') }))

  rechargeRateRows.value = sourceLangs.map((lang) => {
    const code = normalizeLangCode(lang?.code)
    const previous = previousMap[code] || {}
    const currency = normalizeRechargeCurrencyCode(previous.currency || DEFAULT_LANG_CURRENCY_MAP[code] || 'USD')
    const defaultRate = code === 'zh' || currency === 'CNY' ? 1 : 0
    const targetPriceFen = Number.isFinite(Number(plan?.priceI18n?.[code]))
      ? Number(plan.priceI18n[code])
      : basePriceFen

    return {
      code,
      name: String(lang?.name || toLangDisplayName(code, '') || code),
      currency,
      rate: Number.isFinite(Number(previous.rate)) ? Number(previous.rate) : defaultRate,
      selected: previous.selected !== false,
      targetPriceFen,
    }
  })

  const snapshot = getExchangeRateSnapshot({ base: 'CNY' })
  if (!snapshot?.rates) {
    return
  }

  rechargeExchangeRateSource.value = snapshot.source || ''
  rechargeExchangeRateFetchedAt.value = snapshot.fetchedAt || ''
  rechargeRateRows.value.forEach((row) => {
    const currency = normalizeRechargeCurrencyCode(row.currency)
    if (row.code === 'zh' || currency === 'CNY') {
      row.rate = 1
      return
    }
    const snapshotRate = Number(snapshot.rates?.[currency])
    if (Number.isFinite(snapshotRate) && snapshotRate > 0) {
      row.rate = snapshotRate
    }
  })
}

const openRechargePlanRateDialog = async(index) => {
  if (!ensureRechargePlanLangSlots(index)) {
    return
  }

  rechargeRatePlanIndex.value = index
  buildRechargeRateRows()
  rechargeRateDialogVisible.value = true
  await refreshRechargeRatesForRows(rechargeRateRows.value, {
    silent: true,
    forceRefresh: false,
  })
}

const isAllRechargeRowsSelected = () => {
  return rechargeRateRows.value.length > 0 && rechargeRateRows.value.every((row) => row.selected)
}

const toggleSelectAllRechargeRows = (checked) => {
  rechargeRateRows.value.forEach((row) => {
    row.selected = !!checked
  })
}

const normalizeRechargeRowCurrency = (row) => {
  row.currency = normalizeRechargeCurrencyCode(row.currency)
}

const syncRechargeRowPriceToI18n = (row) => {
  const plan = rechargePlans.value[rechargeRatePlanIndex.value]
  if (!plan || !row?.code) return

  const basePriceFen = parseRechargePriceFenFromValue(plan.price, 0, false)
  const priceFen = Math.max(0, Math.round(Number(row.targetPriceFen) || 0))
  row.targetPriceFen = priceFen
  plan.priceI18n = normalizeRechargePriceI18n({
    ...(plan.priceI18n || {}),
    [row.code]: priceFen,
  }, basePriceFen)
}

const refreshRechargeRatesForRows = async (rows, {
  silent = false,
  forceRefresh = false,
} = {}) => {
  const validRows = rows.filter((row) => row?.code)
  if (!validRows.length) {
    if (!silent) {
      ElMessage.warning('请先选择需要更新汇率的语言')
    }
    return
  }

  const currencies = validRows
    .map((row) => normalizeRechargeCurrencyCode(row.currency))
    .filter(Boolean)

  try {
    rechargeRateLoading.value = true
    const result = await fetchExchangeRates({
      base: 'CNY',
      currencies,
      forceRefresh,
    })

    rechargeExchangeRateSource.value = String(result?.source || '').trim()
    rechargeExchangeRateFetchedAt.value = String(result?.fetchedAt || '').trim()

    validRows.forEach((row) => {
      const currency = normalizeRechargeCurrencyCode(row.currency)
      if (row.code === 'zh' || currency === 'CNY') {
        row.rate = 1
        return
      }
      const nextRate = Number(result?.rates?.[currency])
      if (Number.isFinite(nextRate) && nextRate > 0) {
        row.rate = nextRate
      }
    })

    if (!silent) {
      ElMessage.success('汇率更新完成')
    }
  } catch (error) {
    if (!silent) {
      ElMessage.error(error?.message || '汇率更新失败')
    }
  } finally {
    rechargeRateLoading.value = false
  }
}

const refreshSelectedRechargeExchangeRates = async () => {
  const selectedRows = rechargeRateRows.value.filter((row) => row.selected)
  await refreshRechargeRatesForRows(selectedRows, { forceRefresh: true })
}

const refreshSingleRechargeRate = async (row) => {
  await refreshRechargeRatesForRows([row], { forceRefresh: true })
}

const applySelectedRechargeRateToPrices = () => {
  const plan = rechargePlans.value[rechargeRatePlanIndex.value]
  if (!plan) return

  const selectedRows = rechargeRateRows.value.filter((row) => row.selected)
  if (!selectedRows.length) {
    ElMessage.warning('请至少选择一个语言')
    return
  }

  const basePriceFen = parseRechargePriceFenFromValue(plan.price, 0, false)
  if (!Number.isFinite(basePriceFen) || basePriceFen <= 0) {
    ElMessage.warning('请先填写基础价格（单位：分）')
    return
  }

  const nextPriceI18n = normalizeRechargePriceI18n(plan.priceI18n, basePriceFen)
  let updatedCount = 0

  selectedRows.forEach((row) => {
    const convertedFen = calcConvertedFenFromCny(basePriceFen, row.rate, rechargeAdjustPercent.value)
    if (convertedFen === null) {
      return
    }

    row.targetPriceFen = convertedFen
    nextPriceI18n[row.code] = convertedFen
    updatedCount += 1
  })

  plan.price = basePriceFen
  plan.priceI18n = normalizeRechargePriceI18n(nextPriceI18n, basePriceFen)
  ElMessage.success(`已更新 ${updatedCount} 个语言价格`)
}

const fillRechargeRateLangSlots = (silent = false) => {
  const index = rechargeRatePlanIndex.value
  if (!ensureRechargePlanLangSlots(index, silent)) {
    return
  }

  buildRechargeRateRows()
  if (!silent) {
    ElMessage.success('多语言槽位补齐完成')
  }
}

const validateRechargePlans = () => {
  for (let i = 0; i < rechargePlans.value.length; i++) {
    const item = rechargePlans.value[i]
    if (!String(displayI18nText(item.points, '')).trim()) {
      ElMessage.warning(`第 ${i + 1} 个套餐缺少试衣币数量`)
      return false
    }
    const basePriceFen = parseRechargePriceFenFromValue(item.price, 0, false)
    if (!Number.isFinite(basePriceFen) || basePriceFen <= 0) {
      ElMessage.warning(`第 ${i + 1} 个套餐基础价格无效`)
      return false
    }

    const priceI18n = normalizeRechargePriceI18n(item.priceI18n, basePriceFen)
    if (!Object.keys(priceI18n).length) {
      ElMessage.warning(`第 ${i + 1} 个套餐缺少多语言价格`)
      return false
    }

    item.price = basePriceFen
    item.priceI18n = priceI18n
  }
  return true
}

const createDefaultTryonModel = () => ({
  __uid: createTryonModelUid(),
  key: '',
  modelUsage: 'tryon',
  parsingModelKey: '',
  refinerModelKey: '',
  enabled: true,
  scenes: ['clothes'],
  model: 'aitryon',
  name: buildMultilingualObject({}, ''),
  desc: buildMultilingualObject({}, ''),
  cost: 1,
  provider: 'aliyun',
  mode: 'prod',
  url: '',
  taskQueryUrl: '',
  token: '',
  tokenBackupText: '',
  tokenQuotaText: '',

  apiName: '/tryon',
  garmentDes: 'clothing item',
  isChecked: true,
  isCheckedCrop: false,
  denoiseSteps: 30,
  seed: 42,

  resolution: -1,
  restoreFace: true,

  clothesType: ['upper'],
  parsingExtraCost: 0,
  refinerExtraCost: 1,
  refinerGender: 'woman',

  supportsBeautify: false,
  beautifyExtraCost: 0,
  beautifyModel: 'custom_beautify',
  beautifyRetouchDegree: undefined,
  beautifyWhiteningDegree: undefined,
  beautifyUrl: '',
  beautifyAccessKeyId: '',
  beautifyAccessKeySecret: '',
  beautifySecurityToken: '',
  beautifyToken: '',

  freeQuotaTotal: 400,
  beautifyDesc: buildMultilingualObject({}, ''),
})

const normalizeTryonModel = (item = {}, index = 0) => {
  const defaultModel = createDefaultTryonModel()
  const modelUsage = normalizeModelUsage(item.modelUsage, defaultModel.modelUsage, item)
  const inferredSupportsRefiner = modelUsage === 'tryon' ? inferSupportsRefiner(item) : false
  const inferredAutoAliyunParsing = modelUsage === 'tryon' ? inferSupportsRefiner(item) : false
  const defaultFreeQuota = isAliyunModelForQuota({ ...item, modelUsage }) ? 400 : 0
  const normalizedFreeQuota = Math.max(0, toInt(item.freeQuotaTotal, defaultFreeQuota))
  const defaultRefinerKey = inferredSupportsRefiner ? 'aliyun_aitryon_refiner' : ''
  const defaultParsingKey = inferredAutoAliyunParsing ? 'aliyun_aitryon_parsing' : ''
  const tokenBackups = normalizeTokenBackups(item.tokenBackups || item.backupTokens)
  const tokenQuotas = normalizeTokenQuotas(item.tokenQuotas || item.tokenQuotaText)

  const normalized = {
    __uid: createTryonModelUid(),
    key: String(item.key || item.modelKey || `model_${index + 1}`),
    modelUsage,
    enabled: toBool(item.enabled, true),
    scenes: toStringArray(item.scenes, [String(item.sceneType || '').trim() || 'clothes']),
    model: String(item.model || defaultModel.model),
    name: normalizeI18nObject(item.name),
    desc: normalizeI18nObject(item.desc),
    cost: Math.max(0, toInt(item.cost, 1)),
    provider: String(item.provider || defaultModel.provider),
    mode: String(item.mode || defaultModel.mode),
    url: String(item.url || item.providerUrl || ''),
    taskQueryUrl: String(item.taskQueryUrl || ''),
    token: String(item.token || item.providerToken || ''),
    tokenBackupText: joinTokenBackups(tokenBackups),
    tokenQuotaText: joinTokenQuotas(tokenQuotas),

    parsingModelKey: '',
    refinerModelKey: '',
    apiName: '',
    garmentDes: '',
    isChecked: true,
    isCheckedCrop: false,
    denoiseSteps: 30,
    seed: 42,
    resolution: -1,
    restoreFace: true,
    clothesType: ['upper'],
    parsingExtraCost: 0,
    refinerExtraCost: 1,
    refinerGender: 'woman',

    // Legacy fields kept for compatibility during migration.
    supportsRefiner: false,
    refinerModel: String(item.refinerModel || ''),
    refinerUrl: String(item.refinerUrl || ''),
    refinerTaskQueryUrl: String(item.refinerTaskQueryUrl || ''),
    refinerToken: String(item.refinerToken || ''),

    supportsBeautify: modelUsage === 'beautify',
    beautifyExtraCost: 0,
    beautifyModel: '',
    beautifyRetouchDegree: parseOptionalBeautifyDegree(item.beautifyRetouchDegree),
    beautifyWhiteningDegree: parseOptionalBeautifyDegree(item.beautifyWhiteningDegree),
    beautifyUrl: String(item.beautifyUrl || ''),
    beautifyAccessKeyId: String(item.beautifyAccessKeyId || ''),
    beautifyAccessKeySecret: String(item.beautifyAccessKeySecret || ''),
    beautifySecurityToken: String(item.beautifySecurityToken || ''),
    beautifyToken: String(item.beautifyToken || ''),

    freeQuotaTotal: normalizedFreeQuota,
    beautifyDesc: normalizeI18nObject(item.beautifyDesc),
    refinerDesc: normalizeI18nObject(item.refinerDesc),
  }

  if (modelUsage === 'tryon') {
    normalized.apiName = String(item.apiName || defaultModel.apiName)
    normalized.garmentDes = String(item.garmentDes || defaultModel.garmentDes)
    normalized.isChecked = toBool(item.isChecked, true)
    normalized.isCheckedCrop = toBool(item.isCheckedCrop, false)
    normalized.denoiseSteps = Math.max(1, toInt(item.denoiseSteps, 30))
    normalized.seed = toInt(item.seed, 42)
    normalized.resolution = toInt(item.resolution, -1)
    normalized.restoreFace = toBool(item.restoreFace, true)
    normalized.parsingModelKey = String(item.parsingModelKey || defaultParsingKey).trim()
    normalized.refinerModelKey = String(item.refinerModelKey || '').trim()
    if (!normalized.refinerModelKey && toBool(item.supportsRefiner, inferredSupportsRefiner)) {
      normalized.refinerModelKey = defaultRefinerKey
    }
    normalized.supportsRefiner = !!normalized.refinerModelKey
    return normalized
  }

  if (modelUsage === 'refiner') {
    normalized.refinerExtraCost = Math.max(0, toInt(item.refinerExtraCost ?? item.refinerExtraPoints, toInt(item.cost, 1)))
    normalized.refinerGender = normalizeGender(item.refinerGender, defaultModel.refinerGender)
    normalized.supportsRefiner = true
    return normalized
  }

  if (modelUsage === 'parsing') {
    normalized.clothesType = toStringArray(item.clothesType, ['upper'])
    normalized.parsingExtraCost = Math.max(0, toInt(item.parsingExtraCost ?? item.parsingExtraPoints, toInt(item.cost, 0)))
    normalized.cost = normalized.parsingExtraCost
    return normalized
  }

  normalized.cost = Math.max(0, toInt(item.cost, 0))
  normalized.supportsBeautify = true
  normalized.beautifyExtraCost = Math.max(0, toInt(item.beautifyExtraCost ?? item.beautifyExtraPoints, toInt(item.cost, 0)))
  normalized.beautifyModel = String(item.beautifyModel || item.model || defaultModel.beautifyModel)
  normalized.model = normalized.beautifyModel
  return normalized
}

const parseTryonModelsValue = (rawValue) => {
  try {
    const parsed = JSON.parse(rawValue || '[]')
    if (!Array.isArray(parsed)) {
      return []
    }
    return parsed.map((item, index) => normalizeTryonModel(item, index))
  } catch {
    return []
  }
}

const buildTryonModelsPayload = () => {
  return tryonModels.value.map((item) => {
    const modelUsage = normalizeModelUsage(item.modelUsage, 'tryon', item)
    const tokenBackups = normalizeTokenBackups(item.tokenBackupText || item.tokenBackups || item.backupTokens)
    const tokenQuotas = getTokenQuotaRows(item)
    const payload = {
      key: String(item.key || '').trim(),
      modelUsage,
      enabled: !!item.enabled,
      scenes: toStringArray(item.scenes, ['clothes']),
      model: String(item.model || '').trim(),
      name: normalizeI18nObject(item.name),
      desc: normalizeI18nObject(item.desc),
      cost: Math.max(0, toInt(item.cost, 0)),
      provider: String(item.provider || '').trim(),
      mode: String(item.mode || '').trim(),
      url: String(item.url || '').trim(),
      taskQueryUrl: String(item.taskQueryUrl || '').trim(),
      token: String(item.token || '').trim(),
      tokenBackups,
      tokenQuotas,
    }

    if (modelUsage === 'tryon') {
      const providerLower = String(item.provider || '').trim().toLowerCase()
      const result = {
        ...payload,
        freeQuotaTotal: Math.max(0, toInt(item.freeQuotaTotal, 0)),
        resolution: toInt(item.resolution, -1),
        restoreFace: !!item.restoreFace,
        parsingModelKey: String(item.parsingModelKey || '').trim(),
        refinerModelKey: String(item.refinerModelKey || '').trim(),
      }
      if (providerLower.includes('gradio') || providerLower.includes('huggingface') || providerLower.includes('hf')) {
        result.apiName = String(item.apiName || '').trim() || '/tryon'
        result.garmentDes = String(item.garmentDes || '').trim() || 'clothing item'
        result.isChecked = !!item.isChecked
        result.isCheckedCrop = !!item.isCheckedCrop
        result.denoiseSteps = Math.max(1, toInt(item.denoiseSteps, 30))
        result.seed = toInt(item.seed, 42)
      }
      return result
    }

    if (modelUsage === 'refiner') {
      return {
        ...payload,
        freeQuotaTotal: Math.max(0, toInt(item.freeQuotaTotal, 0)),
        refinerExtraCost: Math.max(0, toInt(item.refinerExtraCost, toInt(item.cost, 0))),
        refinerGender: normalizeGender(item.refinerGender, 'woman'),
        refinerDesc: normalizeI18nObject(item.refinerDesc),
      }
    }

    if (modelUsage === 'parsing') {
      const parsingExtraCost = Math.max(0, toInt(item.parsingExtraCost, toInt(item.cost, 0)))
      return {
        ...payload,
        cost: parsingExtraCost,
        parsingExtraCost,
        freeQuotaTotal: Math.max(0, toInt(item.freeQuotaTotal, 0)),
        clothesType: toStringArray(item.clothesType, ['upper']),
      }
    }

    const beautifyRetouchDegree = parseOptionalBeautifyDegree(item.beautifyRetouchDegree)
    const beautifyWhiteningDegree = parseOptionalBeautifyDegree(item.beautifyWhiteningDegree)
    const result = {
      ...payload,
      beautifyExtraCost: Math.max(0, toInt(item.beautifyExtraCost, toInt(item.cost, 0))),
      beautifyModel: String(item.beautifyModel || item.model || '').trim(),
      beautifyUrl: String(item.beautifyUrl || '').trim(),
      beautifyAccessKeyId: String(item.beautifyAccessKeyId || '').trim(),
      beautifyAccessKeySecret: String(item.beautifyAccessKeySecret || '').trim(),
      beautifySecurityToken: String(item.beautifySecurityToken || '').trim(),
      beautifyToken: String(item.beautifyToken || '').trim(),
      beautifyDesc: normalizeI18nObject(item.beautifyDesc),
    }
    if (beautifyRetouchDegree !== undefined) {
      result.beautifyRetouchDegree = beautifyRetouchDegree
    }
    if (beautifyWhiteningDegree !== undefined) {
      result.beautifyWhiteningDegree = beautifyWhiteningDegree
    }
    return result
  })
}

const validateTryonModels = () => {
  let hasBeautifyModel = false
  const keySet = new Set()
  const refinerKeySet = new Set(
    tryonModels.value
      .filter(item => normalizeModelUsage(item.modelUsage, 'tryon', item) === 'refiner')
      .map(item => String(item.key || '').trim())
      .filter(Boolean)
  )
  const parsingKeySet = new Set(
    tryonModels.value
      .filter(item => normalizeModelUsage(item.modelUsage, 'tryon', item) === 'parsing')
      .map(item => String(item.key || '').trim())
      .filter(Boolean)
  )

  for (let i = 0; i < tryonModels.value.length; i++) {
    const item = tryonModels.value[i]
    const modelIndex = i + 1
    const modelUsage = normalizeModelUsage(item.modelUsage, 'tryon', item)
    const modelKey = String(item.key || '').trim()
    if (!modelKey) {
      ElMessage.warning(`第 ${modelIndex} 个模型缺少 key`)
      return false
    }
    if (keySet.has(modelKey)) {
      ElMessage.warning(`模型 key 重复: ${modelKey}`)
      return false
    }
    keySet.add(modelKey)
    if (!Array.isArray(item.scenes) || item.scenes.length === 0) {
      ElMessage.warning(`第 ${modelIndex} 个模型至少要选择一个场景`)
      return false
    }
    const configuredTokens = [
      String(item.token || '').trim(),
      String(item.providerToken || '').trim(),
      String(item.refinerToken || '').trim(),
      String(item.beautifyToken || '').trim(),
      ...normalizeTokenBackups(item.tokenBackupText || item.tokenBackups || item.backupTokens)
    ]
      .filter(Boolean)
    const configuredTokenSet = new Set(configuredTokens)
    const tokenQuotas = getTokenQuotaRows(item)
    for (const quotaItem of tokenQuotas) {
      if (!configuredTokenSet.has(quotaItem.token)) {
        ElMessage.warning(`第 ${modelIndex} 个模型的 tokenQuotas 存在未配置在主/备用列表中的 token: ${quotaItem.token}`)
        return false
      }
    }
    if (modelUsage === 'tryon') {
      const refinerModelKey = String(item.refinerModelKey || '').trim()
      if (refinerModelKey && !refinerKeySet.has(refinerModelKey)) {
      ElMessage.warning(`第 ${modelIndex} 个模型绑定的精修模型不存在或用途不为 refiner: ${refinerModelKey}`)
      return false
      }

      const parsingModelKey = String(item.parsingModelKey || '').trim()
      if (parsingModelKey && !parsingKeySet.has(parsingModelKey)) {
      ElMessage.warning(`第 ${modelIndex} 个模型绑定的分割模型不存在或用途不为 parsing: ${parsingModelKey}`)
      return false
      }
    }
    if (modelUsage === 'refiner' && !String(item.model || '').trim()) {
      ElMessage.warning(`第 ${modelIndex} 个精修模型缺少 model`)
      return false
    }
    if (modelUsage === 'parsing' && !String(item.model || '').trim()) {
      ElMessage.warning(`第 ${modelIndex} 个分割模型缺少 model`)
      return false
    }
    if (modelUsage === 'beautify' && !String(item.beautifyModel || item.model || '').trim()) {
      ElMessage.warning(`第 ${modelIndex} 个美肤模型缺少 beautifyModel/model`)
      return false
    }
    if (modelUsage === 'beautify') {
      hasBeautifyModel = true
    }
  }
  if (!hasBeautifyModel) {
    ElMessage.warning('请至少配置 1 个智能美肤模型（modelUsage=beautify）')
    return false
  }
  return true
}

const tryonModelsCount = (rawValue) => parseTryonModelsValue(rawValue).length

const enabledTryonModelsCount = (rawValue) => {
  return parseTryonModelsValue(rawValue).filter(item => item.enabled).length
}

const addTryonModel = () => {
  const created = createDefaultTryonModel()
  tryonModels.value.push(created)
  ensureTryonModelActiveTab(created.__uid)
}

const cloneTryonModel = (index) => {
  const item = tryonModels.value[index]
  if (!item) {
    return
  }
  const cloned = normalizeTryonModel(item)
  if (cloned.key) {
    cloned.key = `${cloned.key}_copy`
  }
  tryonModels.value.splice(index + 1, 0, cloned)
  ensureTryonModelActiveTab(cloned.__uid)
}

const removeTryonModel = (index) => {
  tryonModels.value.splice(index, 1)
  const nextModel = tryonModels.value[index] || tryonModels.value[index - 1]
  ensureTryonModelActiveTab(nextModel?.__uid || '')
}

const modelCallStatusTag = (status) => {
  const value = String(status || '').trim().toLowerCase()
  if (value === 'success') return 'success'
  if (value === 'processing') return 'warning'
  if (value === 'failed' || value === 'error') return 'danger'
  return 'info'
}

const formatDateTime = (value) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return String(value)
  }
  return date.toLocaleString()
}

const formatPayloadText = (payload) => {
  const raw = String(payload || '').trim()
  if (!raw) return '-'
  try {
    const parsed = JSON.parse(raw)
    return JSON.stringify(parsed, null, 2)
  } catch {
    return raw
  }
}

const getModelLogList = async () => {
  modelLogLoading.value = true
  try {
    const [startCreatedAt = '', endCreatedAt = ''] = Array.isArray(modelLogDateRange.value) ? modelLogDateRange.value : []
    modelLogSearch.startCreatedAt = startCreatedAt
    modelLogSearch.endCreatedAt = endCreatedAt

    const params = {
      ...modelLogSearch,
      status: modelLogSearch.onlyFailed ? '' : modelLogSearch.status,
    }

    const res = await getModelCallLogList(params)
    if (res.code === 0) {
      modelCallLogList.value = res?.data?.list || []
      modelLogTotal.value = Number(res?.data?.total || 0)
      return
    }
    modelCallLogList.value = []
    modelLogTotal.value = 0
    ElMessage.error(res.msg || '获取模型调用日志失败')
  } catch (e) {
    modelCallLogList.value = []
    modelLogTotal.value = 0
    ElMessage.error(e?.message || '获取模型调用日志失败')
  } finally {
    modelLogLoading.value = false
  }
}

const handleModelLogSearch = () => {
  modelLogSearch.page = 1
  getModelLogList()
}

const resetModelLogSearch = () => {
  modelLogDateRange.value = buildDefaultModelLogDateRange()
  modelLogSearch.page = 1
  modelLogSearch.pageSize = 20
  modelLogSearch.startCreatedAt = ''
  modelLogSearch.endCreatedAt = ''
  modelLogSearch.onlyFailed = false
  modelLogSearch.sceneType = ''
  modelLogSearch.modelUsage = ''
  modelLogSearch.status = ''
  modelLogSearch.keyword = ''
  getModelLogList()
}

const handleModelLogCurrentChange = (val) => {
  modelLogSearch.page = val
  getModelLogList()
}

const handleModelLogSizeChange = (val) => {
  modelLogSearch.pageSize = val
  modelLogSearch.page = 1
  getModelLogList()
}

const openModelLogDetail = (row) => {
  modelLogDetail.value = { ...row }
  modelLogDetailVisible.value = true
}

const handleGroupChange = (val) => {
  searchInfo.value.configGroup = val
  searchInfo.value.page = 1
  getList()
}

const getList = async () => {
  const res = await getSysConfigList(searchInfo.value)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
    return
  }
  tableData.value = []
  total.value = 0
  ElMessage.error(res.msg || '获取系统参数失败')
}

const handleCurrentChange = (val) => {
  searchInfo.value.page = val
  getList()
}

const handleSizeChange = (val) => {
  searchInfo.value.pageSize = val
  searchInfo.value.page = 1
  getList()
}

const openEdit = (row) => {
  editForm.value = { ...row }
  tryonModels.value = []
  rechargePlans.value = []
  paymentManualMethods.value = []
  paymentPreferredMethods.value = []
  tryonModelActiveTab.value = ''
  Object.keys(aliyunQuotaState).forEach((key) => {
    delete aliyunQuotaState[key]
  })
  if (isBooleanConfig(row)) {
    editBoolValue.value = row.configValue === 'true'
  }
  if (isNumberConfig(row)) {
    editNumberValue.value = parseInt(row.configValue) || 0
  }
  if (isTryonModelsConfig(row)) {
    tryonModels.value = parseTryonModelsValue(row.configValue)
    ensureTryonModelActiveTab()
    const rawValue = String(row.configValue || '').trim()
    if (rawValue && rawValue !== '[]' && tryonModels.value.length === 0) {
      ElMessage.warning('当前试衣模型配置格式异常，已按空列表打开，请确认后保存')
    }
  } else if (isTryonRechargePlansConfig(row)) {
    rechargePlans.value = parseRechargePlansValue(row.configValue)
    const rawValue = String(row.configValue || '').trim()
    if (rawValue && rawValue !== '[]' && rechargePlans.value.length === 0) {
      ElMessage.warning('当前充值套餐配置格式异常，已按空列表打开，请确认后保存')
    }
  } else if (isPaymentManualMethodsConfig(row)) {
    paymentManualMethods.value = parsePaymentMethodsValue(row.configValue)
    const rawValue = String(row.configValue || '').trim()
    if (rawValue && rawValue !== '[]' && paymentManualMethods.value.length === 0) {
      ElMessage.warning('当前支付方式配置格式异常，已按空列表打开，请确认后保存')
    }
  } else if (isPaymentUniPreferredMethodsConfig(row)) {
    paymentPreferredMethods.value = parsePaymentPreferredMethodsValue(row.configValue)
    const rawValue = String(row.configValue || '').trim()
    if (rawValue && rawValue !== '[]' && paymentPreferredMethods.value.length === 0) {
      ElMessage.warning('当前Uni期望支付方式配置格式异常，已按空列表打开，请确认后保存')
    }
  } else if (isJsonConfig(row)) {
    try {
      const parsed = JSON.parse(row.configValue || '{}')
      Object.keys(editJsonValue).forEach(k => delete editJsonValue[k])
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
        Object.assign(editJsonValue, parsed)
      } else {
        Object.assign(editJsonValue, buildMultilingualObject(parsed, row.configValue || ''))
      }
    } catch {
      Object.keys(editJsonValue).forEach(k => delete editJsonValue[k])
      Object.assign(editJsonValue, buildMultilingualObject(row.configValue, ''))
    }
  }
  newJsonLang.value = ''
  newJsonVal.value = ''
  editVisible.value = true
}

const addJsonLang = () => {
  if (newJsonLang.value) {
    editJsonValue[newJsonLang.value] = newJsonVal.value
    newJsonLang.value = ''
    newJsonVal.value = ''
  }
}

const quickToggle = async (row, val) => {
  const res = await updateSysConfig({
    id: row.ID,
    configValue: val ? 'true' : 'false',
    remark: row.remark,
  })
  if (res.code === 0) {
    ElMessage.success('设置成功')
    getList()
  }
}

const handleSave = async () => {
  let configValue = editForm.value.configValue
  if (isBooleanConfig(editForm.value)) {
    configValue = editBoolValue.value ? 'true' : 'false'
  } else if (isNumberConfig(editForm.value)) {
    configValue = String(editNumberValue.value)
  } else if (isTryonModelsConfig(editForm.value)) {
    if (!validateTryonModels()) {
      return
    }
    configValue = JSON.stringify(buildTryonModelsPayload())
  } else if (isTryonRechargePlansConfig(editForm.value)) {
    if (!validateRechargePlans()) {
      return
    }
    configValue = JSON.stringify(buildRechargePlansPayload())
  } else if (isPaymentManualMethodsConfig(editForm.value)) {
    if (!validatePaymentMethods()) {
      return
    }
    configValue = JSON.stringify(buildPaymentMethodsPayload())
  } else if (isPaymentUniPreferredMethodsConfig(editForm.value)) {
    if (!validatePaymentPreferredMethods()) {
      return
    }
    configValue = JSON.stringify(buildPaymentPreferredMethodsPayload())
  } else if (isJsonConfig(editForm.value)) {
    // 统一富文本多语校验钩子：翻译后保持关键 HTML 结构一致，避免链接/列表/媒体标签被破坏。
    if (isTryonTutorialContentConfig(editForm.value)) {
      const richTextValidation = validateRichTextI18nStructure(editJsonValue, {
        fieldLabel: '教程正文',
        preferredBaseLang: 'zh',
      })
      if (!richTextValidation.valid) {
        ElMessage.warning(richTextValidation.message)
        return
      }
    }
    configValue = JSON.stringify(editJsonValue)
  }
  const res = await updateSysConfig({
    id: editForm.value.ID,
    configValue: configValue,
    remark: editForm.value.remark,
  })
  if (res.code === 0) {
    ElMessage.success('保存成功')
    editVisible.value = false
    getList()
  }
}

onMounted(() => {
  loadMultilingualLanguages()
  getList()
  getModelLogList()
})
</script>

<style scoped>
.config-key {
  font-family: monospace;
  font-size: 12px;
  color: #909399;
}
.color-preview {
  display: flex;
  align-items: center;
  gap: 8px;
}
.color-dot {
  display: inline-block;
  width: 16px;
  height: 16px;
  border-radius: 3px;
  border: 1px solid #dcdfe6;
}
.json-editor {
  width: 100%;
}
.json-editor-tip {
  margin-bottom: 8px;
  font-size: 12px;
  color: #909399;
}

.tutorial-config-preview {
  display: inline-block;
  max-width: 100%;
}

.tutorial-config-editor {
  width: 100%;
}

.tutorial-editor-alert {
  margin-bottom: 10px;
}

.json-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.tryon-model-summary {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tryon-model-summary-text {
  color: #606266;
}

.tryon-model-editor {
  width: 100%;
}

.tryon-model-tabs {
  width: 100%;
}

.tryon-model-tabs :deep(.el-tabs__header) {
  min-width: 280px;
  max-width: 360px;
}

.tryon-model-tabs :deep(.el-tabs__content) {
  padding-left: 12px;
}

.tryon-model-tab-label {
  display: flex;
  align-items: center;
  gap: 6px;
  max-width: 340px;
}

.tryon-model-tab-title {
  display: inline-block;
  min-width: 0;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 600;
}

.tryon-model-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}

.tryon-model-empty {
  color: #909399;
  font-size: 13px;
  padding: 4px 0 8px;
}

.tryon-model-title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.tryon-model-panel {
  padding: 6px 4px;
}

.tryon-model-quota {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 10px 12px;
  margin-bottom: 12px;
  background: #fafcff;
}

.tryon-model-quota-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.tryon-model-quota-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.tryon-model-token-quota-list {
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.tryon-model-token-quota-item {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.tryon-model-token-label {
  color: #606266;
  font-family: monospace;
  font-size: 12px;
}

.tryon-model-token-input :deep(.el-input__inner) {
  font-family: monospace;
}

.tryon-model-token-quota-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tryon-model-token-quota-edit-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 180px;
  gap: 8px;
  align-items: center;
}

.tryon-model-quota-tip {
  margin-top: 6px;
  color: #909399;
  font-size: 12px;
  line-height: 1.4;
}

.tryon-model-quota-error {
  margin-top: 6px;
  color: #f56c6c;
  font-size: 12px;
}

.tryon-model-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px 12px;
}

.tryon-model-field {
  min-width: 0;
}

.tryon-model-field.full {
  grid-column: 1 / -1;
}

.tryon-model-label {
  display: block;
  font-size: 12px;
  color: #909399;
  margin-bottom: 6px;
}

.tryon-model-hint {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.4;
  color: #909399;
}

.tryon-model-subtitle {
  grid-column: 1 / -1;
  font-size: 12px;
  color: #606266;
  font-weight: 600;
  margin-top: 2px;
}

.tryon-model-actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.model-call-log-search {
  margin-top: 14px;
}

.model-call-log-form {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.model-log-payload-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.model-log-image-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.model-log-image-card {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 10px;
}

.model-log-image {
  width: 100%;
  height: 190px;
  border-radius: 6px;
  border: 1px solid #f2f3f5;
}

.model-log-image-empty {
  width: 100%;
  height: 190px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  border-radius: 6px;
  border: 1px dashed #dcdfe6;
  background: #fafafa;
}

.model-log-payload-title {
  font-size: 12px;
  color: #606266;
  margin-bottom: 6px;
}

@media (max-width: 900px) {
  .tryon-model-grid {
    grid-template-columns: 1fr;
  }

  .tryon-model-token-quota-edit-row {
    grid-template-columns: 1fr;
  }

  .model-log-image-grid {
    grid-template-columns: 1fr;
  }
}
</style>
