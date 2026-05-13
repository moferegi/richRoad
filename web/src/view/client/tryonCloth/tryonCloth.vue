<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="用户ID">
          <el-input v-model.number="searchInfo.userID" clearable placeholder="请输入用户ID" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchInfo.category" clearable placeholder="请选择分类" style="width: 180px">
            <el-option label="上装" value="upper" />
            <el-option label="下装" value="lower" />
            <el-option label="连体装" value="onepiece" />
            <el-option label="鞋" value="shoes" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="searchInfo.name" clearable placeholder="请输入名称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="danger"
          icon="delete"
          :disabled="!multipleSelection.length"
          @click="handleBatchDelete"
        >
          批量删除
        </el-button>
      </div>

      <el-table
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column align="left" label="ID" prop="ID" width="90" />
        <el-table-column align="left" label="用户ID" prop="userID" width="100" />
        <el-table-column align="left" label="分类" width="120">
          <template #default="scope">
            <el-tag size="small" effect="plain">{{ categoryLabel(scope.row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="名称" prop="name" min-width="180" show-overflow-tooltip />
        <el-table-column align="left" label="图片" width="100">
          <template #default="scope">
            <el-image
              v-if="scope.row.image"
              class="cloth-image"
              :src="getUrl(scope.row.image)"
              :preview-src-list="[getUrl(scope.row.image)]"
              preview-teleported
              fit="cover"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="图片体积" width="130">
          <template #default="scope">
            <span>{{ getImageSizeText(scope.row) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="图片链接" min-width="280" show-overflow-tooltip>
          <template #default="scope">
            <el-link
              v-if="scope.row.image"
              :href="getImageLink(scope.row.image)"
              target="_blank"
              type="primary"
              :underline="false"
            >
              {{ getImageLink(scope.row.image) }}
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="更新时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.UpdatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
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
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTryonClothList, deleteTryonCloth, deleteTryonClothByIds } from '@/api/client/tryonCloth'
import { formatDate } from '@/utils/format'
import { getUrl } from '@/utils/image'

defineOptions({
  name: 'TryonClothManage'
})

const CATEGORY_LABEL_MAP = {
  upper: '上装',
  lower: '下装',
  onepiece: '连体装',
  shoes: '鞋',
}

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const loading = ref(false)
const tableData = ref([])
const multipleSelection = ref([])
const searchInfo = ref({
  userID: undefined,
  category: '',
  name: '',
})

const imageSizeTextMap = ref({})
const imageSizeCache = new Map()
const imageSizePromiseCache = new Map()

const getImageLink = (image) => getUrl(image)

const formatFileSize = (bytes) => {
  const value = Number(bytes)
  if (!Number.isFinite(value) || value <= 0) return '-'
  if (value < 1024) return `${Math.round(value)} B`
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(2)} KB`
  return `${(value / (1024 * 1024)).toFixed(2)} MB`
}

const loadImageSizeText = (url) => {
  const target = String(url || '').trim()
  if (!target) return Promise.resolve('-')

  if (imageSizeCache.has(target)) {
    return Promise.resolve(imageSizeCache.get(target))
  }

  if (imageSizePromiseCache.has(target)) {
    return imageSizePromiseCache.get(target)
  }

  const pending = fetch(target, { method: 'GET' })
    .then((res) => {
      if (!res.ok) return '-'
      return res.blob().then((blob) => formatFileSize(blob?.size))
    })
    .catch(() => '-')
    .then((text) => {
      const finalText = text || '-'
      imageSizeCache.set(target, finalText)
      imageSizePromiseCache.delete(target)
      return finalText
    })

  imageSizePromiseCache.set(target, pending)
  return pending
}

const warmImageSize = (row) => {
  const id = row?.ID
  const image = row?.image
  if (!id || !image) return
  if (imageSizeTextMap.value[id]) return

  const link = getImageLink(image)
  if (!link) {
    imageSizeTextMap.value[id] = '-'
    return
  }

  loadImageSizeText(link).then((text) => {
    imageSizeTextMap.value[id] = text || '-'
  })
}

const getImageSizeText = (row) => {
  const id = row?.ID
  if (!id || !row?.image) return '-'
  if (imageSizeTextMap.value[id]) return imageSizeTextMap.value[id]
  warmImageSize(row)
  return '-'
}

const categoryLabel = (category) => CATEGORY_LABEL_MAP[category] || '-'

const buildSearchParams = () => ({
  ...searchInfo.value,
})

const getTableData = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      ...buildSearchParams(),
    }
    const res = await getTryonClothList(params)
    if (res.code === 0) {
      tableData.value = res.data?.list || []
      imageSizeTextMap.value = {}
      tableData.value.forEach((row) => warmImageSize(row))
      total.value = res.data?.total || 0
      page.value = res.data?.page || page.value
      pageSize.value = res.data?.pageSize || pageSize.value
    }
  } finally {
    loading.value = false
  }
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {
    userID: undefined,
    category: '',
    name: '',
  }
  page.value = 1
  pageSize.value = 10
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getTableData()
}

const handleSelectionChange = (rows) => {
  multipleSelection.value = rows || []
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除衣橱项【${row.name || row.ID}】吗？`, '提示', { type: 'warning' })
    const res = await deleteTryonCloth({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value -= 1
      }
      await getTableData()
    }
  } catch (e) {
    // 用户取消删除时不做提示
  }
}

const handleBatchDelete = async () => {
  if (!multipleSelection.value.length) {
    ElMessage.warning('请先选择要删除的数据')
    return
  }

  try {
    await ElMessageBox.confirm(`确定批量删除选中的${multipleSelection.value.length}条数据吗？`, '提示', { type: 'warning' })
    const ids = multipleSelection.value.map(item => item.ID).filter(Boolean)
    const res = await deleteTryonClothByIds({ 'IDs[]': ids })
    if (res.code === 0) {
      ElMessage.success('批量删除成功')
      if (ids.length >= tableData.value.length && page.value > 1) {
        page.value -= 1
      }
      multipleSelection.value = []
      await getTableData()
    }
  } catch (e) {
    // 用户取消删除时不做提示
  }
}

getTableData()
</script>

<style scoped>
.cloth-image {
  width: 52px;
  height: 52px;
  border-radius: 6px;
  border: 1px solid var(--el-border-color-light);
}
</style>
