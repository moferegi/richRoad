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

        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.status" clearable placeholder="请选择">
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
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
          label="商品名称"
          prop="title"
          width="120"
        />
        <el-table-column
          align="left"
          label="商品描述"
          prop="description"
          width="120"
        />
        <el-table-column
          label="商品图片URL"
          width="200"
        >
          <template #default="scope">
            <el-image
              style="width: 100px; height: 100px"
              :src="getUrl(scope.row.imageUrl)"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="商品价格"
          prop="price"
          width="120"
        />
        <el-table-column
          align="left"
          label="商品评分"
          prop="rating"
          width="120"
        />
        <el-table-column
          align="left"
          label="商品评论数量"
          prop="reviewCount"
          width="120"
        />
        <el-table-column
          align="left"
          label="商品类型"
          prop="categoryID"
          width="120"
        >
          <template #default="scope">{{ categoryList.find(item => item.ID === scope.row.categoryID)?.title }}</template>
        </el-table-column>
        <el-table-column
          align="left"
          label="邮费"
          prop="postage"
          width="120"
        />
        <el-table-column
          align="left"
          label="折扣"
          prop="discount"
          width="120"
        />

        <el-table-column
          align="left"
          label="推荐"
          prop="recommend"
          width="120"
        >
          <template #default="scope">{{ formatBoolean(scope.row.recommend) }}</template>
        </el-table-column>

        <el-table-column
          align="left"
          label="启用"
          prop="status"
          width="120"
        >
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              inline-prompt
              :active-value="true"
              :inactive-value="false"
              @change="onStatusChange(scope.row)"
            />
          </template>
        </el-table-column>

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
              icon="edit"
              class="table-button"
              @click="updateGoodFunc(scope.row)"
            >变更</el-button>
            <el-button
              type="primary"
              link
              icon="InfoFilled"
              class="table-button"
              @click="setSKU(scope.row)"
            >SKU</el-button>
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
      destroy-on-close
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
        <el-row :gutter="12">
          <el-col :span="6">
            <el-form-item
              label="商品名称:"
              prop="title"
            >
              <el-input
                v-model="formData.title"
                :clearable="true"
                placeholder="请输入商品名称"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              label="商品类型:"
              prop="categoryID"
            >
              <el-tree-select
                v-model="formData.categoryID"
                :data="categoryList"
                :props="{ label: 'title', value: 'ID', children: 'children', isLeaf: 'isLeaf' }"
                placeholder="请选择商品类型"
                check-strictly
                :only-leaf-select="true"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              label="状态:"
              prop="status"
            >
              <el-switch
                v-model="formData.status"
                active-color="#13ce66"
                inactive-color="#ff4949"
                active-text="是"
                inactive-text="否"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item
              label="推荐:"
              prop="recommend"
            >
              <el-switch
                v-model="formData.recommend"
                active-color="#13ce66"
                inactive-color="#ff4949"
                active-text="是"
                inactive-text="否"
                clearable
              />
            </el-form-item>
          </el-col>

        </el-row>
        <el-form-item
          label="商品图片URL:"
          prop="imageUrl"
        >
          <SelectImage
            v-model="formData.imageUrl"
            file-type="image"
          />
        </el-form-item>
        <el-form-item
          label="商品轮播图:"
          prop="banner"
        >
          <SelectImage
            v-model="formData.banner"
            multiple
            :max-update-count="5"
            file-type="image"
          />
        </el-form-item>
        <el-form-item
          label="商品标签:"
          prop="tags"
        >
          <el-select
            v-model="formData.tags"
            placeholder="请选择商品标签"
            multiple
            clearable
            value-key="ID"
          >
            <el-option
              v-for="item in tags"
              :key="item.ID"
              :label="item.name"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          label="商品描述:"
          prop="description"
        >
          <el-input
            v-model="formData.description"
            :clearable="true"
            placeholder="请输入商品描述"
          />
        </el-form-item>

        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item
              label="商品价格(单位：分):"
              prop="price"
            >
              <el-input-number
                v-model="formData.price"
                style="width:100%"
                :precision="2"
                :clearable="true"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              label="邮费:"
              prop="postage"
            >
              <el-input-number
                v-model="formData.postage"
                style="width:100%"
                :precision="2"
                :clearable="true"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              label="折扣(%):"
              prop="discount"
            >
              <el-input-number
                v-model="formData.discount"
                :clearable="true"
                style="width:100%"
                :precision="0"
                placeholder="请输入折扣(%)"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <h4 class="flex justify-between items-center">
          <span>商品规格(可作为商品选项的属性，例：尺码，颜色)</span>
          <el-button
            type="primary"
            icon="plus"
            @click="formData.specs.push({
              name: '',
              value: ''
            })"
          >添加</el-button>
        </h4>
        <el-row
          v-for="(item,index) in formData.specs"
          :key="index"
          :gutter="12"
        >
          <el-col
            :key="index"
            :span="8"
          >
            <el-form-item
              :label="!index?'规格名称:':''"
              prop="specs"
            >
              <el-input
                v-model="item.name"
                :clearable="true"
                placeholder="请输入规格名称"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              :label="!index?'规格编码:':''"
              prop="specs"
            >
              <el-input
                v-model="item.value"
                :clearable="true"
                placeholder="请输入规格编码"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              :label="!index?'操作':''"
            >
              <el-button
                type="danger"
                icon="delete"
                @click="formData.specs.splice(index,1)"
              >删除</el-button>
            </el-form-item>
          </el-col>
        </el-row>
        <h4 class="flex justify-between items-center">
          <span>商品属性(仅作为展示属性，例：厂商，材料)</span>
          <el-button
            type="primary"
            icon="plus"
            @click="formData.attrs.push({
              name: '',
              value: ''
            })"
          >添加</el-button>
        </h4>
        <el-row
          v-for="(item,index) in formData.attrs"
          :key="index"
          :gutter="12"
        >
          <el-col
            :key="index"
            :span="8"
          >
            <el-form-item
              :label="!index?'属性值:':''"
              prop="specs"
            >
              <el-input
                v-model="item.name"
                :clearable="true"
                placeholder="请输入属性值"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              :label="!index?'属性编码:':''"
              prop="specs"
            >
              <el-input
                v-model="item.value"
                :clearable="true"
                placeholder="请输入属性编码"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item
              :label="!index?'操作':''"
            >
              <el-button
                type="danger"
                icon="delete"
                @click="formData.attrs.splice(index,1)"
              >删除</el-button>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item
            label="商品详情:"
            prop="detail"
        >
          <div class="h-[800px]">
            <rich-edit
                v-model="formData.detail"
                :clearable="true"
                placeholder="请输入商品描述"
            />
          </div>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import {
  createGood,
  deleteGood,
  deleteGoodByIds,
  updateGood,
  findGood,
  getGoodList
} from '@/api/shop/good'

import {
  getTagList
} from '@/api/shop/tag'

import {
  getCategoryList
} from '@/api/shop/category'
import { getUrl } from '@/utils/image'
// 图片选择组件
import SelectImage from '@/components/selectImage/selectImage.vue'

// 全量引入格式化工具 请按需保留
import { formatDate, formatBoolean } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import RichEdit from "@/components/richtext/rich-edit.vue";

defineOptions({
  name: 'Good'
})

const router = useRouter()

const setSKU = (row) => {
  router.push({ name: 'sku', query: { id: row.ID }})
}

const tags = ref([])

const getTags = async() => {
  const res = await getTagList()
  if (res.code === 0) {
    tags.value = res.data.list || []
  }
}

getTags()

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  description: '',
  imageUrl: '',
  banner: [],
  price: 0,
  rating: 0,
  reviewCount: 0,
  title: '',
  categoryID: undefined,
  status: false,
  recommend: false,
  postage: 0,
  discount: 0,
  specs: [],
  attrs: [],
  detail: ''
})
const categoryList = ref([])
const getCategoryListFunc = async() => {
  const res = await getCategoryList()
  if (res.code === 0) {
    categoryList.value = res.data.list || []
  }
}

getCategoryListFunc()

// 验证规则
const rule = reactive({
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
    if (searchInfo.value.status === '') {
      searchInfo.value.status = null
    }
    if (searchInfo.value.recommend === '') {
      searchInfo.value.recommend = null
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
const getTableData = async() => {
  const table = await getGoodList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
    deleteGoodFunc(row)
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
    const res = await deleteGoodByIds({ IDs })
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
const updateGoodFunc = async(row) => {
  const res = await findGood({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.regood
    if (!formData.value.specs) {
      formData.value.specs = []
    }
    if (!formData.value.attrs) {
      formData.value.attrs = []
    }
    if (!formData.value.banner) {
      formData.value.banner = []
    }
    dialogFormVisible.value = true
  }
}

// 删除行
const deleteGoodFunc = async(row) => {
  const res = await deleteGood({ ID: row.ID })
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
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    description: '',
    price: 0,
    rating: 0,
    reviewCount: 0,
    title: '',
    categoryID: undefined,
    status: false,
    recommend: false,
    postage: 0,
    specs: [],
    attrs: [],
    discount: 0,
    detail: ''
  }
}
// 弹窗确定
const enterDialog = async() => {
  elFormRef.value?.validate(async(valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await createGood(formData.value)
        break
      case 'update':
        res = await updateGood(formData.value)
        break
      default:
        res = await createGood(formData.value)
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

const onStatusChange = async(row) => {
  const res = await updateGood(row)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '状态更新成功'
    })
  } else {
    row.status = !row.status
  }
}
</script>

<style>

</style>
