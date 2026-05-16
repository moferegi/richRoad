<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">
        <el-form-item label="创建日期" prop="CreatedAt">
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期"
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true"
            v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
          @click="onDelete">删除</el-button>

      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        @selection-change="handleSelectionChange" @sort-change="sortChange">
        <el-table-column type="selection" width="55" />

        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column sortable align="left" label="名称" prop="name" width="120" />

        <el-table-column align="left" label="最低消费/分（0为无门槛）" prop="minSpend" width="120" />

        <el-table-column align="left" label="折扣/分" prop="discount" width="120" />

        <el-table-column align="left" label="商品ID（不填则不限商品）" prop="productIDs" width="180">
          <template #default="scope">
            <span>{{ scope.row.productIDs || '不限' }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="数量" prop="quantity" width="80" />
        <el-table-column align="left" label="已领取" prop="claimed" width="80" />

        <el-table-column sortable align="left" label="开始时间" prop="startTime" width="180">
          <template #default="scope">{{ formatDate(scope.row.startTime) }}</template>
        </el-table-column>
        <el-table-column sortable align="left" label="结束时间" prop="endTime" width="180">
          <template #default="scope">{{ formatDate(scope.row.endTime) }}</template>
        </el-table-column>
        <el-table-column align="left" label="启用" prop="status" width="120">
          <template #default="scope">{{ formatBoolean(scope.row.status) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>查看</el-button>
            <el-button type="primary" link icon="edit" class="table-button"
              @click="updateCouponFunc(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false"
      :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增' : '编辑' }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="名称:" prop="name">
          <el-input v-model="formData.name" :clearable="false" placeholder="请输入名称" />
        </el-form-item>
        <div v-if="enabledLangs.length" class="pl-2 mb-3">
          <MultiLangEditor
            :model="nameI18n"
            :languages="enabledLangs"
            title="名称多语言"
          />
        </div>
        <el-form-item label="描述:" prop="description">
          <el-input v-model="formData.description" :clearable="false" placeholder="请输入描述" />
        </el-form-item>
        <div v-if="enabledLangs.length" class="pl-2 mb-3">
          <MultiLangEditor
            :model="descI18n"
            :languages="enabledLangs"
            title="描述多语言"
          />
        </div>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="最低消费/分（0为无门槛）:" prop="minSpend">
              <el-input-number v-model="formData.minSpend" style="width:100%" :precision="0" :clearable="false" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="折扣/分:" prop="discount">
              <el-input-number v-model="formData.discount" style="width:100%" :precision="0" :clearable="false" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="可用商品ID（逗号分隔，不填则不限商品）:" prop="productIDs">
          <el-input v-model="formData.productIDs" placeholder="例: 1,2,5 （空=全部商品可用）" clearable />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="数量:" prop="quantity">
              <el-input v-model.number="formData.quantity" :clearable="false" placeholder="请输入数量" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="已领取:" prop="claimed">
              <el-input-number v-model="formData.claimed" style="width:100%" :precision="0" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="启用:" prop="status">
              <el-switch v-model="formData.status" active-color="#13ce66" inactive-color="#ff4949" active-text="是"
                inactive-text="否" clearable />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="开始时间:" prop="startTime">
              <el-date-picker v-model="formData.startTime" type="datetime" style="width:100%" placeholder="选择日期时间"
                :clearable="false" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束时间:" prop="endTime">
              <el-date-picker v-model="formData.endTime" type="datetime" style="width:100%" placeholder="选择日期时间"
                :clearable="false" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-divider content-position="left">背景图设置</el-divider>
        <el-form-item label="背景图(上传):" prop="backgroundImage">
          <SelectImage v-model="formData.backgroundImage" file-type="image" />
        </el-form-item>
        <el-form-item label="外部背景图路径(优先于上传):" prop="externalBgPath">
          <el-input v-model="formData.externalBgPath" placeholder="https://example.com/bg.jpg" clearable />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true"
      :before-close="closeDetailShow" title="查看">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="名称">
          {{ detailFrom.name }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailFrom.nameI18n" label="名称(多语言)">
          {{ detailFrom.nameI18n }}
        </el-descriptions-item>
        <el-descriptions-item label="描述">
          {{ detailFrom.description }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailFrom.descriptionI18n" label="描述(多语言)">
          {{ detailFrom.descriptionI18n }}
        </el-descriptions-item>
        <el-descriptions-item label="折扣/分">
          {{ detailFrom.discount }}
        </el-descriptions-item>
        <el-descriptions-item label="最低消费/分">
          {{ detailFrom.minSpend }}
        </el-descriptions-item>
        <el-descriptions-item label="可用商品">
          {{ detailFrom.productIDs || '不限' }}
        </el-descriptions-item>
        <el-descriptions-item label="数量">
          {{ detailFrom.quantity }}
        </el-descriptions-item>
        <el-descriptions-item label="已领取">
          {{ detailFrom.claimed }}
        </el-descriptions-item>
        <el-descriptions-item label="开始时间">
          {{ formatDate(detailFrom.startTime) }}
        </el-descriptions-item>
        <el-descriptions-item label="结束时间">
          {{ formatDate(detailFrom.endTime) }}
        </el-descriptions-item>
        <el-descriptions-item label="启用">
          {{ formatBoolean(detailFrom.status) }}
        </el-descriptions-item>
        <el-descriptions-item label="背景图">
          <el-image v-if="detailFrom.externalBgPath || detailFrom.backgroundImage" style="width:100px;height:60px" :src="getUrl(detailFrom.externalBgPath || detailFrom.backgroundImage)" fit="cover" />
          <span v-else>-</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import {
  getCouponDataSource,
  createCoupon,
  deleteCoupon,
  deleteCouponByIds,
  updateCoupon,
  findCoupon,
  getCouponList
} from '@/api/shop/coupon'
import { getEnabledLanguages } from '@/api/client/language'
import { getUrl } from '@/utils/image'
import SelectImage from '@/components/selectImage/selectImage.vue'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import { useAppStore } from "@/pinia"
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'




defineOptions({
  name: 'Coupon'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// === i18n ===
const enabledLangs = ref([])
const nameI18n = ref({})
const descI18n = ref({})

const loadLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0) enabledLangs.value = res.data || []
  } catch (e) { /* ignore */ }
}

const parseI18n = (s) => { if (!s) return {}; try { return JSON.parse(s) } catch { return {} } }
const serializeI18n = (o) => {
  const f = {}; for (const [k,v] of Object.entries(o)) { if (v) f[k] = v }
  return Object.keys(f).length ? JSON.stringify(f) : ''
}

onMounted(() => { loadLangs() })

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  name: '',
  description: '',
  discount: 0,
  minSpend: 0,
  productID: undefined,
  productIDs: '',
  quantity: 0,
  claimed: 0,
  startTime: new Date(),
  endTime: new Date(),
  status: false,
  nameI18n: '',
  descriptionI18n: '',
  backgroundImage: '',
  externalBgPath: '',
})
const dataSource = ref([])
const getDataSourceFunc = async () => {
  const res = await getCouponDataSource()
  if (res.code === 0) {
    dataSource.value = res.data
  }
}
getDataSourceFunc()



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
  description: [{
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
  discount: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  minSpend: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  quantity: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  startTime: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  endTime: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  status: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
})

const searchRule = reactive({
  CreatedAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
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
// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
    CreatedAt: "created_at",
    ID: "id",
    name: 'name',
    startTime: 'start_time',
    endTime: 'end_time',
  }

  let sort = sortMap[prop]
  if (!sort) {
    sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}
// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    if (searchInfo.value.status === "") {
      searchInfo.value.status = null
    }
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
const getTableData = async () => {
  const table = await getCouponList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
const setOptions = async () => {
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
    deleteCouponFunc(row)
  })
}

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
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
    const res = await deleteCouponByIds({ IDs })
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
const updateCouponFunc = async (row) => {
  const res = await findCoupon({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    nameI18n.value = parseI18n(formData.value.nameI18n)
    descI18n.value = parseI18n(formData.value.descriptionI18n)
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteCouponFunc = async (row) => {
  const res = await deleteCoupon({ ID: row.ID })
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

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  nameI18n.value = {}
  descI18n.value = {}
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  nameI18n.value = {}
  descI18n.value = {}
  formData.value = {
    name: '',
    description: '',
    discount: 0,
    minSpend: 0,
    productID: undefined,
    productIDs: '',
    quantity: 0,
    claimed: 0,
    startTime: new Date(),
    endTime: new Date(),
    status: false,
    nameI18n: '',
    descriptionI18n: '',
    backgroundImage: '',
    externalBgPath: '',
  }
}
// 弹窗确定
const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    // 序列化 i18n
    if (enabledLangs.value.length) {
      formData.value.nameI18n = serializeI18n(nameI18n.value)
      formData.value.descriptionI18n = serializeI18n(descI18n.value)
    }
    let res
    switch (type.value) {
      case 'create':
        res = await createCoupon(formData.value)
        break
      case 'update':
        res = await updateCoupon(formData.value)
        break
      default:
        res = await createCoupon(formData.value)
        break
    }
    btnLoading.value = false
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

const detailFrom = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findCoupon({ ID: row.ID })
  if (res.code === 0) {
    detailFrom.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}


</script>

<style></style>
