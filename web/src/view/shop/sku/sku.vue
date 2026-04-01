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
          </el-select>
          <el-select v-model="searchInfo.orderDir" clearable placeholder="方向" style="width:80px; margin-left:4px;">
            <el-option label="升序" value="asc" />
            <el-option label="降序" value="desc" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
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
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />

        <el-table-column
          align="left"
          label="日期"
          width="180"
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
          prop="name"
          width="120"
        />
        <el-table-column
          label="图片"
          width="200"
        >
          <template #default="scope">
            <el-image
              style="width: 100px; height: 100px"
              :src="getUrl(scope.row.picture)"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="介绍"
          prop="description"
          width="120"
        />
        <el-table-column
          align="left"
          label="价格"
          prop="price"
          width="120"
        />
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
          label="名称:"
          prop="name"
        >
          <el-input
            v-model="formData.name"
            :clearable="true"
            placeholder="请输入名称"
          />
        </el-form-item>
        <el-form-item
          label="图片:"
          prop="picture"
        >
          <SelectImage
            v-model="formData.picture"
            file-type="image"
          />
        </el-form-item>
        <el-form-item
          label="介绍:"
          prop="description"
        >
          <el-input
            v-model="formData.description"
            :clearable="true"
            placeholder="请输入介绍"
          />
        </el-form-item>
        <el-form-item
          label="价格:"
          prop="price"
        >
          <el-input
            v-model.number="formData.price"
            :clearable="true"
            placeholder="请输入价格"
          />
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
          <h4 style="margin:0;">规格配置</h4>
          <el-button type="primary" size="small" @click="addAttr">添加规格</el-button>
        </div>
        <el-form-item
            v-for="(attr,index) in formData.attrs"
            :key="'attr'+index"
            :label="attr.label+':'"
        >
          <div style="display:flex;gap:8px;width:100%;">
            <el-input
                v-model="attr.label"
                :clearable="true"
                placeholder="规格名称"
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
        </el-form-item>

        <div style="display:flex;align-items:center;justify-content:space-between;margin:12px 0 8px;">
          <h4 style="margin:0;">属性配置</h4>
          <el-button type="primary" size="small" @click="addSpec">添加属性</el-button>
        </div>
        <el-form-item
          v-for="(spec,index) in formData.specs"
          :key="'spec'+index"
          :label="spec.label+':'"
        >
          <div style="display:flex;gap:8px;width:100%;">
            <el-input
                v-model="spec.label"
                :clearable="true"
                placeholder="属性名称"
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
        </el-form-item>
      </el-form>
    </el-drawer>

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
import { getUrl } from '@/utils/image'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'
import { useRoute } from 'vue-router'
// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
  name: 'Sku'
})

const route = useRoute()

const attrs = ref([])
const specs = ref([])

const getAttr = async() => {
  const res = await findGood({ ID: Number(route.query.id) })
  if (res.code === 0) {
    attrs.value = res.data.regood.attrs
    specs.value = res.data.regood.specs
  }
}
getAttr()
// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  name: '',
  picture: '',
  description: '',
  price: 0,
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
    description: '',
    price: 0,
    inventory: 0,
    goodID: Number(route.query.id),
  }
}

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  formData.value.attrs = []
  formData.value.specs = []
  attrs.value.forEach(item => {
    formData.value.attrs.push({
      label: item.name,
      value: ''
    })
  })
  specs.value.forEach(item => {
    formData.value.specs.push({
      label: item.name,
      value: ''
    })
  })
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    no: '',
    name: '',
    description: '',
    price: 0,
    inventory: 0,
    goodID: Number(route.query.id),
  }
}
// 添加规格
const addAttr = () => {
  if (!formData.value.attrs) formData.value.attrs = []
  formData.value.attrs.push({ label: '', value: '' })
}

// 添加属性
const addSpec = () => {
  if (!formData.value.specs) formData.value.specs = []
  formData.value.specs.push({ label: '', value: '' })
}

// 弹窗确定
const enterDialog = async() => {
  elFormRef.value?.validate(async(valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createSku(formData.value)
        break
      case 'update':
        res = await updateSku(formData.value)
        break
      default:
        res = await createSku(formData.value)
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
