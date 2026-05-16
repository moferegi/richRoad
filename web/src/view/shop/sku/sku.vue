<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item
          label="创建日期"
          prop="createdAt"
        >
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            placeholder="开始日期"
            :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"
          />
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            placeholder="结束日期"
            :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSubmit"
          >查询</el-button>
          <el-button
            icon="refresh"
            @click="onReset"
          >重置</el-button>
        </el-form-item>

        <el-form-item label="排序">
          <el-select v-model="searchInfo.orderBy" clearable placeholder="排序字段" style="width:100px;">
            <el-option label="销量" value="sale_num" />
            <el-option label="价格" value="price" />
            <el-option label="库存" value="inventory" />
            <el-option label="创建时间" value="created_at" />
          </el-select>
          <el-select v-model="searchInfo.orderDir" clearable placeholder="方向" style="width:80px; margin-left:4px;">
            <el-option label="升序" value="asc" />
            <el-option label="降序" value="desc" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list" style="display:flex; align-items:center; justify-content:space-between;">
        <div>
        <el-button
          type="primary"
          icon="plus"
          @click="openDialog"
        >新增</el-button>
        <el-button
          icon="delete"
          style="margin-left: 10px;"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >删除</el-button>
        </div>
        <div v-if="enabledLangs.length">
          <span style="margin-right:6px;font-size:13px;color:#666;">表格语言：</span>
          <el-select v-model="tableLang" size="small" style="width:120px;">
            <el-option label="默认" value="" />
            <el-option v-for="l in enabledLangs" :key="l.code" :label="l.name || l.code" :value="l.code" />
          </el-select>
        </div>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        :default-sort="{ prop: 'ID', order: 'descending' }"
        max-height="70vh"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="ID" prop="ID" width="70" sortable />

        <el-table-column
          align="left"
          label="日期"
          width="180"
          sortable
        >
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column
          align="left"
          label="序列号"
          prop="no"
          width="120"
        />
        <el-table-column
          align="left"
          label="名称"
          width="120"
        >
          <template #default="scope">{{ formatSkuI18nField(scope.row.name) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="属性配置"
          width="200"
        >
          <template #default="scope">{{ formatSkuAttrs(scope.row.attrs) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="规格内容"
          width="200"
        >
          <template #default="scope">{{ formatSkuAttrs(scope.row.specs) }}</template>
        </el-table-column>
        <el-table-column
          label="图片"
          width="200"
        >
          <template #default="scope">
            <el-image
              style="width: 100px; height: 100px"
              :src="scope.row.externalPicturePath ? resolveExtUrl(scope.row.externalPicturePath) : getUrl(scope.row.picture)"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="介绍"
          width="120"
        >
          <template #default="scope">{{ formatSkuI18nField(scope.row.description) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="价格"
          width="120"
        >
          <template #default="scope">{{ formatTablePrice(scope.row) }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="余量"
          prop="inventory"
          width="120"
        />
        <el-table-column
          align="left"
          label="销量"
          prop="saleNum"
          width="100"
        />
        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              class="table-button"
              @click="getDetails(scope.row)"
            >
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              查看详情
            </el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateSkuFunc(scope.row)"
            >变更</el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-drawer
      v-model="dialogFormVisible"
      size="800"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type==='create'?'添加':'修改' }}</span>
          <div>
            <el-button
              type="primary"
              @click="enterDialog"
            >确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="elFormRef"
        :model="formData"
        label-position="top"
        :rules="rule"
        label-width="80px"
      >
        <el-form-item
          label="名称(默认):"
          prop="name"
        >
          <el-input
            v-model="formData.name"
            :clearable="true"
            placeholder="请输入名称"
          />
        </el-form-item>
        <MultiLangEditor
          :model="nameI18n"
          :languages="enabledLangs"
          title="名称多语言"
          input-type="input"
        />
        <el-form-item
          label="图片(上传):"
          prop="picture"
        >
          <SelectImage
            v-model="formData.picture"
            file-type="image"
          />
        </el-form-item>
        <el-form-item label="图片外链(优先于上传):" prop="externalPicturePath">
          <el-input v-model="formData.externalPicturePath" placeholder="相对路径如 /images/sku.jpg 自动拼接外部域名" clearable />
          <div v-if="extDomain" class="text-xs text-gray-400 mt-1">当前外部域名: {{ extDomain }}</div>
        </el-form-item>
        <el-form-item
          label="介绍(默认):"
          prop="description"
        >
          <el-input
            v-model="formData.description"
            :clearable="true"
            placeholder="请输入介绍"
          />
        </el-form-item>
        <MultiLangEditor
          :model="descI18n"
          :languages="enabledLangs"
          title="介绍多语言"
          input-type="input"
        />
        <el-form-item
          label="价格(单位：分):"
          prop="price"
        >
          <div class="w-full">
            <el-input-number
              v-model="formData.price"
              :precision="0"
              :step="1"
              :min="0"
              style="width: 100%"
            />
            <div class="mt-2 flex items-center justify-between gap-2">
              <el-button type="primary" plain size="small" @click="openPriceRateDialog">汇率换算</el-button>
              <span class="text-xs text-gray-500">多语言价格槽位：{{ Object.keys(priceI18n).length }}</span>
            </div>
          </div>
        </el-form-item>
        <el-form-item
          label="余量:"
          prop="inventory"
        >
          <el-input
            v-model.number="formData.inventory"
            :clearable="true"
            placeholder="请输入余量"
          />
        </el-form-item>
        <div style="display:flex;align-items:center;justify-content:space-between;margin:12px 0 8px;">
          <h4 style="margin:0;">属性配置</h4>
          <div>
            <el-dropdown v-if="attrDictList.length" @command="selectAttrFromDict" style="margin-right:8px;">
              <el-button type="success" size="small">从字典选择<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-for="d in attrDictList" :key="d.ID" :command="d">{{ d.label }} - {{ d.value }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button type="primary" size="small" @click="addAttr">手动添加</el-button>
          </div>
        </div>
        <el-form-item
            v-for="(attr,index) in formData.attrs"
            :key="'attr'+index"
        >
          <template #label>
            <span style="font-weight:500;">{{ attr.label || '属性'+(index+1) }}</span>
          </template>
          <div style="width:100%;">
            <div style="display:flex;gap:8px;align-items:center;">
              <el-input
                  v-model="attr.label"
                  :clearable="true"
                  placeholder="属性名称"
                  style="width:140px;flex-shrink:0;"
              />
              <el-input
                  v-model="attr.value"
                  :clearable="true"
                  :placeholder="'请输入'+attr.label"
                  style="flex:1;"
              />
              <el-button type="danger" link @click="formData.attrs.splice(index,1)">删除</el-button>
            </div>
            <div class="mt-2">
              <MultiLangEditor
                :model="attr.labelI18n"
                :languages="enabledLangs"
                :title="'属性名称多语言 #' + (index + 1)"
                input-type="input"
              />
            </div>
            <div class="mt-2">
              <MultiLangEditor
                :model="attr.valueI18n"
                :languages="enabledLangs"
                :title="'属性值多语言 #' + (index + 1)"
                input-type="input"
              />
            </div>
          </div>
        </el-form-item>

        <div style="display:flex;align-items:center;justify-content:space-between;margin:12px 0 8px;">
          <h4 style="margin:0;">规格配置</h4>
          <div>
            <el-dropdown v-if="specDictList.length" @command="selectSpecFromDict" style="margin-right:8px;">
              <el-button type="success" size="small">从字典选择<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-for="d in specDictList" :key="d.ID" :command="d">{{ d.label }} - {{ d.value }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button type="primary" size="small" @click="addSpec">手动添加</el-button>
          </div>
        </div>
        <el-form-item
          v-for="(spec,index) in formData.specs"
          :key="'spec'+index"
        >
          <template #label>
            <span style="font-weight:500;">{{ spec.label || '规格'+(index+1) }}</span>
          </template>
          <div style="width:100%;">
            <div style="display:flex;gap:8px;align-items:center;">
              <el-input
                  v-model="spec.label"
                  :clearable="true"
                  placeholder="规格名称"
                  style="width:140px;flex-shrink:0;"
              />
              <el-input
                  v-model="spec.value"
                  :clearable="true"
                  :placeholder="'请输入'+spec.label"
                  style="flex:1;"
              />
              <el-button type="danger" link @click="formData.specs.splice(index,1)">删除</el-button>
            </div>
            <div class="mt-2">
              <MultiLangEditor
                :model="spec.labelI18n"
                :languages="enabledLangs"
                :title="'规格名称多语言 #' + (index + 1)"
                input-type="input"
              />
            </div>
            <div class="mt-2">
              <MultiLangEditor
                :model="spec.valueI18n"
                :languages="enabledLangs"
                :title="'规格值多语言 #' + (index + 1)"
                input-type="input"
              />
            </div>
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-dialog
      v-model="priceRateDialogVisible"
      title="SKU价格汇率换算"
      width="980px"
      destroy-on-close
    >
      <div class="mb-3 flex flex-wrap items-center gap-2">
        <el-checkbox :model-value="isAllPriceRowsSelected()" @change="toggleSelectAllPriceRows">全选</el-checkbox>
        <el-input-number
          v-model="priceAdjustPercent"
          :precision="2"
          :step="0.5"
          placeholder="增减百分比"
          style="width: 140px"
        />
        <span class="text-xs text-gray-500">增加%（可负数），按人民币价格换算后叠加</span>
      </div>

      <div class="mb-3 flex flex-wrap gap-2">
        <el-button type="primary" :loading="priceRateLoading" @click="refreshSelectedExchangeRates">更新汇率（选中）</el-button>
        <el-button type="success" @click="applySelectedRateToPrices">补齐（更新选中价格）</el-button>
        <el-button @click="fillMissingLangSlots">补齐多语言槽位</el-button>
      </div>

      <div class="mb-2 text-xs text-gray-500">
        汇率来源：{{ exchangeRateSource || '-' }}
        <span v-if="exchangeRateFetchedAt">，更新时间：{{ exchangeRateFetchedAt }}</span>
      </div>

      <el-table :data="priceRateRows" border max-height="420px">
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
            <el-input v-model="scope.row.currency" maxlength="3" placeholder="USD" @change="normalizeRowCurrency(scope.row)" />
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
                @change="syncRowPriceToI18n(scope.row)"
              />
              <span class="text-xs text-gray-500">≈ {{ formatFen(scope.row.targetPriceFen) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button link type="primary" @click="refreshSingleRate(scope.row)">更新汇率</el-button>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="priceRateDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <el-drawer
      v-model="detailShow"
      size="800"
      :before-close="closeDetailShow"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">查看详情</span>
        </div>
      </template>
      <el-descriptions
        :column="1"
        border
      >

        <el-descriptions-item label="名称">
          {{ formData.name }}
        </el-descriptions-item>
        <el-descriptions-item label="图片">
          <el-image
            style="width: 50px; height: 50px"
            :preview-src-list="ReturnArrImg(formData.picture)"
            :src="getUrl(formData.picture)"
            fit="cover"
          />
        </el-descriptions-item>
        <el-descriptions-item label="介绍">
          {{ formData.description }}
        </el-descriptions-item>
        <el-descriptions-item label="价格">
          {{ formData.price }}
        </el-descriptions-item>
        <el-descriptions-item label="余量">
          {{ formData.inventory }}
        </el-descriptions-item>
        <el-descriptions-item label="属性">
          [JSON]
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  findGood
} from '@/api/shop/good'
import {
  createSku,
  deleteSku,
  deleteSkuByIds,
  updateSku,
  findSku,
  getSkuList
} from '@/api/shop/sku'
import { getEnabledLanguages, getLanguageList } from '@/api/client/language'
import { getDefaultDomain } from '@/api/client/externalLinkDomain'
import { getAllSkuSpecs } from '@/api/shop/skuSpec'
import {
  DEFAULT_LANG_CURRENCY_MAP,
  fetchExchangeRates,
  calcConvertedFenFromCny,
  normalizePriceI18nMap,
  ensurePriceSlots,
  serializePriceI18nMap,
  formatFen
} from '@/utils/exchange-rate'
import { getUrl } from '@/utils/image'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
import { useRoute } from 'vue-router'
// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'

defineOptions({
  name: 'Sku'
})

// === 外部链接域名 ===
const extDomain = ref('')
const loadExtDomain = async () => {
  try {
    const res = await getDefaultDomain()
    if (res.code === 0 && res.data) {
      extDomain.value = res.data.replace(/\/+$/, '')
    }
  } catch (e) { /* ignore */ }
}
const resolveExtUrl = (path) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  if (extDomain.value) {
    const sep = path.startsWith('/') ? '' : '/'
    return extDomain.value + sep + path
  }
  return path
}

// === SKU规格字典 ===
const specDictList = ref([])
const attrDictList = ref([])
const loadSkuSpecDict = async () => {
  try {
    const [specRes, attrRes] = await Promise.all([getAllSkuSpecs('spec'), getAllSkuSpecs('attr')])
    if (specRes.code === 0) specDictList.value = specRes.data || []
    if (attrRes.code === 0) attrDictList.value = attrRes.data || []
  } catch(e) { /* ignore */ }
}
const selectAttrFromDict = (dictItem) => {
  if (!formData.value.attrs) formData.value.attrs = []
  const labelI18n = {}
  try { Object.assign(labelI18n, typeof dictItem.labelI18n === 'string' ? JSON.parse(dictItem.labelI18n || '{}') : (dictItem.labelI18n || {})) } catch {}
  const valueI18n = {}
  try { Object.assign(valueI18n, typeof dictItem.valueI18n === 'string' ? JSON.parse(dictItem.valueI18n || '{}') : (dictItem.valueI18n || {})) } catch {}
  formData.value.attrs.push({ label: dictItem.label, value: dictItem.value, labelI18n, valueI18n })
}
const selectSpecFromDict = (dictItem) => {
  if (!formData.value.specs) formData.value.specs = []
  const labelI18n = {}
  try { Object.assign(labelI18n, typeof dictItem.labelI18n === 'string' ? JSON.parse(dictItem.labelI18n || '{}') : (dictItem.labelI18n || {})) } catch {}
  const valueI18n = {}
  try { Object.assign(valueI18n, typeof dictItem.valueI18n === 'string' ? JSON.parse(dictItem.valueI18n || '{}') : (dictItem.valueI18n || {})) } catch {}
  formData.value.specs.push({ label: dictItem.label, value: dictItem.value, labelI18n, valueI18n })
}

// === i18n 多语言支持 ===
const enabledLangs = ref([])
const allLangs = ref([])
const nameI18n = ref({})
const descI18n = ref({})

const uniqueLangItems = (items = []) => {
  const map = new Map()
  items.forEach((item) => {
    const code = String(item?.code || '').trim()
    if (!code || map.has(code)) return
    map.set(code, {
      ...item,
      code
    })
  })
  return Array.from(map.values())
}

const loadLangs = async () => {
  try {
    const [enabledRes, listRes] = await Promise.all([
      getEnabledLanguages().catch(() => null),
      getLanguageList({ page: 1, pageSize: 500 }).catch(() => null)
    ])

    if (enabledRes?.code === 0) {
      enabledLangs.value = uniqueLangItems(enabledRes.data || [])
    }

    if (listRes?.code === 0) {
      const listPayload = Array.isArray(listRes?.data?.list)
        ? listRes.data.list
        : (Array.isArray(listRes?.data) ? listRes.data : [])
      allLangs.value = uniqueLangItems(listPayload)
    }

    if (!allLangs.value.length) {
      allLangs.value = [...enabledLangs.value]
    }
  } catch (e) { /* ignore */ }
}

const parseI18nJson = (jsonStr) => {
  if (!jsonStr) return {}
  try { return JSON.parse(jsonStr) } catch { return {} }
}

const serializeI18nJson = (obj) => {
  const filtered = {}
  for (const [k, v] of Object.entries(obj)) {
    if (v) filtered[k] = v
  }
  return Object.keys(filtered).length ? JSON.stringify(filtered) : ''
}

const priceI18n = ref({})
const priceRateDialogVisible = ref(false)
const priceRateRows = ref([])
const priceRateLoading = ref(false)
const priceAdjustPercent = ref(0)
const exchangeRateSource = ref('')
const exchangeRateFetchedAt = ref('')

const getEnabledLangCodes = () => {
  const source = allLangs.value.length ? allLangs.value : enabledLangs.value
  return source
    .map((item) => String(item?.code || '').trim())
    .filter(Boolean)
}

const normalizeCurrencyCode = (value) => {
  return String(value || '').trim().toUpperCase().slice(0, 3)
}

const syncPriceI18nToForm = () => {
  formData.value.priceI18n = serializePriceI18nMap(priceI18n.value)
}

const ensureTextLangSlots = (payload, fallbackValue = '') => {
  const next = { ...(payload || {}) }
  getEnabledLangCodes().forEach((code) => {
    if (!Object.prototype.hasOwnProperty.call(next, code)) {
      next[code] = fallbackValue
    }
  })
  return next
}

const ensurePriceLangSlots = () => {
  const fallbackFen = Math.round(Number(formData.value.price) || 0)
  priceI18n.value = ensurePriceSlots(priceI18n.value, getEnabledLangCodes(), fallbackFen)
  syncPriceI18nToForm()
}

const buildPriceRateRows = () => {
  ensurePriceLangSlots()
  const previousMap = {}
  priceRateRows.value.forEach((row) => {
    previousMap[row.code] = row
  })

  const fallbackFen = Math.round(Number(formData.value.price) || 0)
  const source = allLangs.value.length ? allLangs.value : enabledLangs.value
  priceRateRows.value = source.map((lang) => {
    const code = String(lang?.code || '').trim()
    const previous = previousMap[code] || {}
    const currency = normalizeCurrencyCode(previous.currency || DEFAULT_LANG_CURRENCY_MAP[code] || 'USD')
    const targetPriceFen = Number.isFinite(Number(priceI18n.value[code]))
      ? Number(priceI18n.value[code])
      : fallbackFen
    const defaultRate = code === 'zh' ? 1 : 0

    return {
      code,
      name: lang?.name || code,
      currency,
      rate: Number.isFinite(Number(previous.rate)) ? Number(previous.rate) : defaultRate,
      selected: previous.selected !== false,
      targetPriceFen
    }
  })
}

const openPriceRateDialog = () => {
  if (!getEnabledLangCodes().length) {
    ElMessage.warning('请先在语言管理中配置语言')
    return
  }
  buildPriceRateRows()
  priceRateDialogVisible.value = true
}

const isAllPriceRowsSelected = () => {
  return priceRateRows.value.length > 0 && priceRateRows.value.every((row) => row.selected)
}

const toggleSelectAllPriceRows = (checked) => {
  priceRateRows.value.forEach((row) => {
    row.selected = !!checked
  })
}

const normalizeRowCurrency = (row) => {
  row.currency = normalizeCurrencyCode(row.currency)
}

const syncRowPriceToI18n = (row) => {
  const priceFen = Math.max(0, Math.round(Number(row.targetPriceFen) || 0))
  row.targetPriceFen = priceFen
  priceI18n.value[row.code] = priceFen
  syncPriceI18nToForm()
}

const refreshRatesForRows = async(rows) => {
  const validRows = rows.filter((row) => row?.code)
  if (!validRows.length) {
    ElMessage.warning('请先选择需要更新汇率的语言')
    return
  }

  const currencies = validRows
    .map((row) => normalizeCurrencyCode(row.currency))
    .filter(Boolean)

  try {
    priceRateLoading.value = true
    const result = await fetchExchangeRates({
      base: 'CNY',
      currencies
    })
    exchangeRateSource.value = result.source || ''
    exchangeRateFetchedAt.value = result.fetchedAt || ''

    validRows.forEach((row) => {
      if (row.code === 'zh') {
        row.rate = 1
        return
      }
      const currency = normalizeCurrencyCode(row.currency)
      const nextRate = Number(result?.rates?.[currency])
      if (Number.isFinite(nextRate) && nextRate > 0) {
        row.rate = nextRate
      }
    })
    ElMessage.success('汇率更新完成')
  } catch (error) {
    ElMessage.error(error?.message || '汇率更新失败')
  } finally {
    priceRateLoading.value = false
  }
}

const refreshSelectedExchangeRates = async() => {
  const selectedRows = priceRateRows.value.filter((row) => row.selected)
  await refreshRatesForRows(selectedRows)
}

const refreshSingleRate = async(row) => {
  await refreshRatesForRows([row])
}

const applySelectedRateToPrices = () => {
  const selectedRows = priceRateRows.value.filter((row) => row.selected)
  if (!selectedRows.length) {
    ElMessage.warning('请至少选择一个语言')
    return
  }

  let updatedCount = 0
  selectedRows.forEach((row) => {
    const convertedFen = calcConvertedFenFromCny(formData.value.price, row.rate, priceAdjustPercent.value)
    if (convertedFen === null) {
      return
    }
    row.targetPriceFen = convertedFen
    priceI18n.value[row.code] = convertedFen
    updatedCount += 1
  })

  syncPriceI18nToForm()
  ElMessage.success(`已更新 ${updatedCount} 个语言价格`)
}

const fillMissingLangSlots = (silent = false) => {
  if (!getEnabledLangCodes().length) {
    if (!silent) {
      ElMessage.warning('请先在语言管理中配置语言')
    }
    return
  }

  nameI18n.value = ensureTextLangSlots(nameI18n.value, formData.value.name || '')
  descI18n.value = ensureTextLangSlots(descI18n.value, formData.value.description || '')

  formData.value.attrs = (formData.value.attrs || []).map((item) => {
    return {
      ...item,
      labelI18n: ensureTextLangSlots(item.labelI18n || parseI18nJson(item.label), item.label || ''),
      valueI18n: ensureTextLangSlots(item.valueI18n || parseI18nJson(item.value), item.value || '')
    }
  })

  formData.value.specs = (formData.value.specs || []).map((item) => {
    return {
      ...item,
      labelI18n: ensureTextLangSlots(item.labelI18n || parseI18nJson(item.label), item.label || ''),
      valueI18n: ensureTextLangSlots(item.valueI18n || parseI18nJson(item.value), item.value || '')
    }
  })

  ensurePriceLangSlots()
  buildPriceRateRows()
  if (!silent) {
    ElMessage.success('多语言槽位补齐完成')
  }
}

onMounted(() => { loadLangs(); loadExtDomain(); loadSkuSpecDict() })

const tableLang = ref('')

const formatSkuI18nField = (val) => {
  if (!val) return ''
  if (typeof val === 'string' && val.startsWith('{')) {
    try {
      const obj = JSON.parse(val)
      if (tableLang.value && obj[tableLang.value]) return obj[tableLang.value]
      return obj['zh'] || Object.values(obj)[0] || val
    } catch { return val }
  }
  return val
}

const getLocalizedPriceFen = (basePrice, i18nPayload, langCode) => {
  const map = normalizePriceI18nMap(i18nPayload)
  if (langCode && Number.isFinite(Number(map[langCode]))) {
    return Number(map[langCode])
  }
  const base = Number(basePrice)
  return Number.isFinite(base) ? base : 0
}

const formatTablePrice = (row) => {
  return formatFen(getLocalizedPriceFen(row?.price, row?.priceI18n, tableLang.value))
}

const formatSkuAttrs = (attrsStr) => {
  if (!attrsStr) return '-'
  try {
    const arr = typeof attrsStr === 'string' ? JSON.parse(attrsStr) : attrsStr
    if (!Array.isArray(arr) || arr.length === 0) return '-'
    const lang = tableLang.value
    return arr.map(a => {
      let label = a.label || a.name || a.key || ''
      let value = a.value || ''
      if (lang) {
        // 尝试从 i18n 对象取值
        if (a.labelI18n && typeof a.labelI18n === 'object' && a.labelI18n[lang]) label = a.labelI18n[lang]
        else if (typeof label === 'string' && label.startsWith('{')) { try { const o = JSON.parse(label); if (o[lang]) label = o[lang] } catch {} }
        if (a.valueI18n && typeof a.valueI18n === 'object' && a.valueI18n[lang]) value = a.valueI18n[lang]
        else if (typeof value === 'string' && value.startsWith('{')) { try { const o = JSON.parse(value); if (o[lang]) value = o[lang] } catch {} }
      } else {
        // 默认模式：label/value 可能是JSON字符串，取zh或首个值
        if (typeof label === 'string' && label.startsWith('{')) { try { const o = JSON.parse(label); label = o['zh'] || Object.values(o)[0] || label } catch {} }
        if (typeof value === 'string' && value.startsWith('{')) { try { const o = JSON.parse(value); value = o['zh'] || Object.values(o)[0] || value } catch {} }
      }
      return `${label}：${value}`
    }).join('，')
  } catch { return '-' }
}

const route = useRoute()

const attrs = ref([])
const specs = ref([])

const getAttr = async() => {
  const res = await findGood({ ID: Number(route.query.id) })
  if (res.code === 0) {
    // Good的attrs/specs中 name 是序列化的i18n JSON字符串，nameI18n已被删除
    // 需要解析 name → nameI18n，value → valueI18n（如缺失）
    const parseGoodItems = (items) => {
      if (!items) return []
      const arr = typeof items === 'string' ? JSON.parse(items) : items
      if (!Array.isArray(arr)) return []
      arr.forEach(item => {
        if (!item.nameI18n && item.name && typeof item.name === 'string' && item.name.startsWith('{')) {
          try { item.nameI18n = JSON.parse(item.name) } catch {}
        }
        if (!item.valueI18n && item.value && typeof item.value === 'string' && item.value.startsWith('{')) {
          try { item.valueI18n = JSON.parse(item.value) } catch {}
        }
      })
      return arr
    }
    attrs.value = parseGoodItems(res.data.regood.attrs)
    specs.value = parseGoodItems(res.data.regood.specs)
  }
}
getAttr()
// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  name: '',
  picture: '',
  externalPicturePath: '',
  description: '',
  price: 0,
  priceI18n: '',
  inventory: 0,
  attrs: [],
  specs: [],
  goodID: Number(route.query.id),
})

// 验证规则
const rule = reactive({
  name: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  {
    whitespace: true,
    message: '不能只输入空格',
    trigger: ['input', 'blur'],
  }
  ],
  goodID: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
})

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getSkuList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value, goodID: Number(route.query.id) })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async() => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()

// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteSkuFunc(row)
  })
}

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const IDs = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: '请选择要删除的数据'
      })
      return
    }
    multipleSelection.value &&
        multipleSelection.value.map(item => {
          IDs.push(item.ID)
        })
    const res = await deleteSkuByIds({ IDs })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateSkuFunc = async(row) => {
  const res = await findSku({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.resku
    // 解析 i18n 字段
    nameI18n.value = parseI18nJson(formData.value.name)
    descI18n.value = parseI18nJson(formData.value.description)
    // 解析 规格/属性 i18n
    if (!formData.value.attrs) formData.value.attrs = []
    if (!formData.value.specs) formData.value.specs = []
    formData.value.attrs.forEach(item => {
      if (!item.labelI18n || typeof item.labelI18n === 'string') item.labelI18n = parseI18nJson(item.labelI18n || item.label)
      if (!item.valueI18n || typeof item.valueI18n === 'string') item.valueI18n = parseI18nJson(item.valueI18n || item.value)
    })
    formData.value.specs.forEach(item => {
      if (!item.labelI18n || typeof item.labelI18n === 'string') item.labelI18n = parseI18nJson(item.labelI18n || item.label)
      if (!item.valueI18n || typeof item.valueI18n === 'string') item.valueI18n = parseI18nJson(item.valueI18n || item.value)
    })
    priceI18n.value = normalizePriceI18nMap(formData.value.priceI18n)
    ensurePriceLangSlots()
    buildPriceRateRows()
    dialogFormVisible.value = true
  }
}

// 删除行
const deleteSkuFunc = async(row) => {
  const res = await deleteSku({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 查看详情控制标记
const detailShow = ref(false)

// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}

// 打开详情
const getDetails = async(row) => {
  // 打开弹窗
  const res = await findSku({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.resku
    openDetailShow()
  }
}

// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    name: '',
    picture: '',
    externalPicturePath: '',
    description: '',
    price: 0,
    priceI18n: '',
    inventory: 0,
    attrs: [],
    specs: [],
    goodID: Number(route.query.id),
  }
}

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  nameI18n.value = {}
  descI18n.value = {}
  priceI18n.value = {}
  priceRateRows.value = []
  priceAdjustPercent.value = 0
  exchangeRateSource.value = ''
  exchangeRateFetchedAt.value = ''
  formData.value.attrs = []
  formData.value.specs = []
  attrs.value.forEach(item => {
    // item.nameI18n 已在 getAttr 中解析好
    const parsedNameI18n = (typeof item.nameI18n === 'object' && item.nameI18n) ? { ...item.nameI18n } : {}
    // label 取 zh 值或原始 name
    let label = item.name || ''
    if (typeof label === 'string' && label.startsWith('{')) { try { const o = JSON.parse(label); label = o['zh'] || Object.values(o)[0] || label } catch {} }
    // value 也可能有 i18n
    const parsedValueI18n = (typeof item.valueI18n === 'object' && item.valueI18n) ? { ...item.valueI18n } : {}
    let value = item.value || ''
    if (typeof value === 'string' && value.startsWith('{')) { try { const o = JSON.parse(value); value = o['zh'] || Object.values(o)[0] || value } catch {} }
    formData.value.attrs.push({
      label,
      value,
      labelI18n: parsedNameI18n,
      valueI18n: parsedValueI18n
    })
  })
  specs.value.forEach(item => {
    const parsedNameI18n = (typeof item.nameI18n === 'object' && item.nameI18n) ? { ...item.nameI18n } : {}
    let label = item.name || ''
    if (typeof label === 'string' && label.startsWith('{')) { try { const o = JSON.parse(label); label = o['zh'] || Object.values(o)[0] || label } catch {} }
    const parsedValueI18n = (typeof item.valueI18n === 'object' && item.valueI18n) ? { ...item.valueI18n } : {}
    let value = item.value || ''
    if (typeof value === 'string' && value.startsWith('{')) { try { const o = JSON.parse(value); value = o['zh'] || Object.values(o)[0] || value } catch {} }
    formData.value.specs.push({
      label,
      value,
      labelI18n: parsedNameI18n,
      valueI18n: parsedValueI18n
    })
  })
  fillMissingLangSlots(true)
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  priceRateDialogVisible.value = false
  nameI18n.value = {}
  descI18n.value = {}
  priceI18n.value = {}
  priceRateRows.value = []
  priceAdjustPercent.value = 0
  exchangeRateSource.value = ''
  exchangeRateFetchedAt.value = ''
  formData.value = {
    no: '',
    name: '',
    externalPicturePath: '',
    description: '',
    price: 0,
    priceI18n: '',
    inventory: 0,
    attrs: [],
    specs: [],
    goodID: Number(route.query.id),
  }
}
// 添加规格
const addAttr = () => {
  if (!formData.value.attrs) formData.value.attrs = []
  formData.value.attrs.push({ label: '', value: '', labelI18n: {}, valueI18n: {} })
}

// 添加属性
const addSpec = () => {
  if (!formData.value.specs) formData.value.specs = []
  formData.value.specs.push({ label: '', value: '', labelI18n: {}, valueI18n: {} })
}

// 弹窗确定
const enterDialog = async() => {
  elFormRef.value?.validate(async(valid) => {
    if (!valid) return
    const payload = JSON.parse(JSON.stringify(formData.value))
    payload.priceI18n = serializePriceI18nMap(priceI18n.value)
    // 序列化 i18n 字段
    if (getEnabledLangCodes().length) {
      payload.name = serializeI18nJson(nameI18n.value) || payload.name
      payload.description = serializeI18nJson(descI18n.value) || payload.description
    }
    let res
    switch (type.value) {
      case 'create':
        res = await createSku(payload)
        break
      case 'update':
        res = await updateSku(payload)
        break
      default:
        res = await createSku(payload)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '创建/更改成功'
      })
      closeDialog()
      getTableData()
    }
  })
}

</script>

<style>

</style>
