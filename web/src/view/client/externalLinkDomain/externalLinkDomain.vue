<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="名称">
          <el-input v-model="searchInfo.name" placeholder="域名名称" clearable />
        </el-form-item>
        <el-form-item label="云类型">
          <el-select v-model="searchInfo.cloudType" clearable placeholder="全部" style="width: 130px">
            <el-option label="QiNiu" value="qiniu" />
            <el-option label="R2" value="r2" />
            <el-option label="B2" value="b2" />
            <el-option label="S3" value="s3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList">查询</el-button>
          <el-button type="primary" @click="openCreate">新增域名</el-button>
          <el-button type="warning" @click="openCompare">一键目录比对</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" stripe>
        <el-table-column prop="name" label="名称" width="120" />
        <el-table-column label="云类型" width="90" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.cloudType" size="small" type="info">{{ scope.row.cloudType.toUpperCase() }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="domain" label="域名地址" min-width="240">
          <template #default="scope">
            <el-link type="primary" :href="scope.row.domain" target="_blank">{{ scope.row.domain }}</el-link>
          </template>
        </el-table-column>
        <el-table-column label="默认" width="70" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.isDefault" type="success" size="small">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="默认上传" width="90" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.defaultUpload" type="warning" size="small">上传</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="同步" width="70" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.syncUpload" type="primary" size="small">同步</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="启用" width="70" align="center">
          <template #default="scope">
            <el-switch
              :model-value="scope.row.isEnabled"
              @change="(val) => handleToggle(scope.row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column label="用途" width="80" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.uploadScope === 'video'" type="danger" size="small">视频</el-tag>
            <el-tag v-else-if="scope.row.uploadScope === 'image'" type="success" size="small">图片</el-tag>
            <el-tag v-else-if="scope.row.uploadScope === 'audio'" type="warning" size="small">音频</el-tag>
            <el-tag v-else type="info" size="small">全部</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="70" align="center" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openEdit(scope.row)">编辑</el-button>
            <el-button v-if="scope.row.cloudType" type="success" link :loading="pingLoading[scope.row.ID]" @click="handlePing(scope.row)">检测连接</el-button>
            <el-button v-if="scope.row.cloudType" type="info" link @click="openFileBrowser(scope.row)">浏览文件</el-button>
            <el-button v-if="!scope.row.isDefault" type="warning" link @click="handleSetDefault(scope.row)">设为默认</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="searchInfo.page"
          :page-size="searchInfo.pageSize"
          :page-sizes="[10, 30, 50]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 编辑/创建弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑域名' : '新增域名'" width="660px">
      <el-form :model="form" label-width="110px" :rules="rules" ref="formRef">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如：阿里云CDN" />
        </el-form-item>
        <el-form-item label="域名地址" prop="domain">
          <el-input v-model="form.domain" placeholder="如：https://cdn.example.com" />
        </el-form-item>
        <el-form-item label="云类型">
          <el-select v-model="form.cloudType" clearable placeholder="无（仅域名转发）" style="width: 100%">
            <el-option label="七牛云 QiNiu" value="qiniu" />
            <el-option label="Cloudflare R2" value="r2" />
            <el-option label="Backblaze B2" value="b2" />
            <el-option label="AWS S3 兼容" value="s3" />
          </el-select>
        </el-form-item>

        <template v-if="form.cloudType">
          <el-divider content-position="left">云存储凭证</el-divider>

          <!-- Qiniu 七牛云 -->
          <template v-if="form.cloudType === 'qiniu'">
            <el-form-item label="存储区域">
              <el-select v-model="form.zone" placeholder="请选择存储区域" style="width:100%">
                <el-option label="华东" value="ZoneHuaDong" />
                <el-option label="华北" value="ZoneHuaBei" />
                <el-option label="华南" value="ZoneHuaNan" />
                <el-option label="北美" value="ZoneBeiMei" />
                <el-option label="东南亚" value="ZoneXinJiaPo" />
              </el-select>
            </el-form-item>
            <el-form-item label="AccessKey">
              <el-input v-model="form.accessKey" placeholder="AccessKey" />
            </el-form-item>
            <el-form-item label="SecretKey">
              <el-input v-model="form.secretKey" type="password" show-password placeholder="SecretKey" />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input v-model="form.bucket" placeholder="存储空间名称" />
            </el-form-item>
            <el-form-item label="CDN域名">
              <el-input v-model="form.baseUrl" placeholder="如：https://cdn.example.com" />
            </el-form-item>
            <el-form-item label="HTTPS">
              <el-switch v-model="form.useHttps" />
            </el-form-item>
            <el-form-item label="CDN加速">
              <el-switch v-model="form.useCdnDomain" />
            </el-form-item>
          </template>

          <!-- R2 Cloudflare -->
          <template v-else-if="form.cloudType === 'r2'">
            <el-form-item label="账户ID">
              <el-input v-model="form.accountId" placeholder="Cloudflare Account ID" />
            </el-form-item>
            <el-form-item label="AccessKey">
              <el-input v-model="form.accessKey" placeholder="R2 Access Key ID" />
            </el-form-item>
            <el-form-item label="SecretKey">
              <el-input v-model="form.secretKey" type="password" show-password placeholder="R2 Secret Access Key" />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input v-model="form.bucket" placeholder="存储桶名称" />
            </el-form-item>
            <el-form-item label="Endpoint">
              <el-input v-model="form.regionEndpoint" placeholder="留空自动用账户ID构造，或填完整URL" />
              <span style="color:#999;font-size:12px">留空自动构造：https://{accountId}.r2.cloudflarestorage.com</span>
            </el-form-item>
            <el-form-item label="自定义域名">
              <el-input v-model="form.baseUrl" placeholder="如：https://r2.mnmovie.icu" />
            </el-form-item>
          </template>

          <!-- B2 / S3 -->
          <template v-else>
            <el-form-item label="Endpoint">
              <el-input v-model="form.regionEndpoint" placeholder="如：s3.us-east-005.backblazeb2.com" />
              <span style="color:#999;font-size:12px">纯域名，不含 https://，如 s3.us-east-005.backblazeb2.com</span>
            </el-form-item>
            <el-form-item label="AccessKey">
              <el-input v-model="form.accessKey" placeholder="Key ID / Application Key ID" />
            </el-form-item>
            <el-form-item label="SecretKey">
              <el-input v-model="form.secretKey" type="password" show-password placeholder="Secret Key / Application Key" />
            </el-form-item>
            <el-form-item label="Bucket">
              <el-input v-model="form.bucket" placeholder="存储桶名称" />
            </el-form-item>
            <el-form-item label="CDN域名">
              <el-input v-model="form.baseUrl" placeholder="如：https://video.mnmovie.icu" />
            </el-form-item>
          </template>

          <el-divider content-position="left">上传策略</el-divider>
          <el-form-item label="上传用途">
            <el-select v-model="form.uploadScope" style="width:100%">
              <el-option label="全部（默认）" value="all" />
              <el-option label="视频专用" value="video" />
              <el-option label="图片专用" value="image" />
              <el-option label="音频专用" value="audio" />
            </el-select>
            <span style="color:#999;font-size:12px;margin-left:8px">如视频配合防盗链，请选「视频专用」并配置对应CDN域名</span>
          </el-form-item>
          <el-form-item label="默认上传云">
            <el-switch v-model="form.defaultUpload" />
            <span style="margin-left:8px;color:#999;font-size:12px">勾选后该用途的文件上传默认走此云</span>
          </el-form-item>
          <el-form-item label="同步上传">
            <el-switch v-model="form.syncUpload" />
            <span style="margin-left:8px;color:#999;font-size:12px">勾选后上传同步至该云</span>
          </el-form-item>
        </template>

        <el-divider content-position="left">基本信息</el-divider>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="form.isDefault" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.isEnabled" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- 文件浏览器弹窗 -->
    <el-dialog v-model="fileBrowserVisible" :title="`文件浏览器 - ${fileBrowserName}`" width="90%" top="3vh" @opened="loadFileBrowserData" @closed="closeFileBrowser">
      <!-- 当前路径 -->
      <div class="file-browser-path">
        <span class="path-label">当前路径：</span>
        <span class="path-value">{{ fileBrowserCurrentPrefix || '/' }}</span>
      </div>

      <div class="file-browser-toolbar">
        <el-breadcrumb separator="/" class="file-breadcrumb">
          <el-breadcrumb-item>
            <el-button type="primary" link @click="navigateDir('')">
              <el-icon><Folder /></el-icon> 根目录
            </el-button>
          </el-breadcrumb-item>
          <el-breadcrumb-item v-for="(part, idx) in breadcrumbParts" :key="idx">
            <el-button type="primary" link @click="navigateDir(part.path)">
              {{ part.name }}
            </el-button>
          </el-breadcrumb-item>
        </el-breadcrumb>
        <div class="file-browser-actions">
          <el-input v-model="searchKeyword" placeholder="全局搜索文件名..." size="small" clearable style="width:220px" @keyup.enter="handleSearch" />
          <el-button type="primary" size="small" :loading="searchLoading" @click="handleSearch">
            <el-icon><Search /></el-icon> 搜索
          </el-button>
          <el-button v-if="isSearchResult" size="small" @click="clearSearch">返回目录</el-button>
          <el-upload
            :show-file-list="false"
            :http-request="handleFileUpload"
            accept="*"
          >
            <el-button type="success" size="small">上传文件</el-button>
          </el-upload>
          <el-button type="danger" size="small" :disabled="selectedFileKeys.length === 0" @click="handleBatchDelete">
            删除选中 ({{ selectedFileKeys.length }})
          </el-button>
          <el-button type="primary" size="small" @click="loadFileBrowserData">刷新</el-button>
        </div>
      </div>
      <el-table
        :data="filteredFileList"
        border
        max-height="460"
        @selection-change="handleFileSelectionChange"
      >
        <el-table-column type="selection" width="45" :selectable="onlyFileSelectable" />
        <el-table-column label="文件名/目录" min-width="320">
          <template #default="scope">
            <template v-if="scope.row.isDir">
              <el-button type="primary" link @click="navigateDir(scope.row.key)">
                <el-icon><Folder /></el-icon> {{ scope.row.displayName }}
              </el-button>
            </template>
            <span v-else>{{ scope.row.displayName }}</span>
          </template>
        </el-table-column>
        <el-table-column label="大小" width="160" align="right">
          <template #default="scope">
            <template v-if="scope.row.isDir">
              <span v-if="scope.row.dirFileCount">{{ scope.row.dirFileCount }} 文件 / {{ formatFileSize(scope.row.dirSize) }}</span>
              <span v-else>-</span>
            </template>
            <span v-else>{{ formatFileSize(scope.row.size) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="lastModified" label="修改时间" width="180" />
        <el-table-column label="操作" width="80" align="center" fixed="right">
          <template #default="scope">
            <el-button type="primary" link size="small" @click="handleDownloadItem(scope.row)">
              <el-icon><Download /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="isSearchResult" style="margin-top:8px;color:#909399">
        搜索结果，共 {{ searchResultList.length }} 个匹配文件
      </div>
      <div v-else-if="fileBrowserTruncated" style="margin-top:8px;color:#999">文件较多，仅展示部分。进入子目录以获得更精确结果。</div>
    </el-dialog>

    <!-- 目录比对弹窗 -->
    <el-dialog v-model="compareVisible" title="多云目录比对" width="1000px">
      <div class="file-browser-toolbar">
        <el-select v-model="compareRefId" placeholder="选择参考基准云" style="width:200px" @change="compareDone=false">
          <el-option
            v-for="row in cloudOptions"
            :key="row.ID"
            :label="`${row.name} (${row.cloudType?.toUpperCase()})`"
            :value="row.ID"
          />
        </el-select>
        <el-input v-model="comparePrefix" placeholder="目录前缀（如：english-learn/）" clearable style="width: 260px" />
        <el-button type="primary" :loading="compareLoading" :disabled="!compareRefId" @click="handleCompare">开始比对</el-button>
      </div>

      <!-- 各目标云结果 -->
      <div v-if="compareDone && compareResult.referenceName" style="margin-top:12px">
        <el-alert type="info" :title="`参考基准：${compareResult.referenceName}（${comparePrefix || '根目录'}）`" :closable="false" style="margin-bottom:16px" />

        <el-collapse v-model="compareActiveCollapse">
          <el-collapse-item
            v-for="target in compareResult.targets"
            :key="target.id"
            :name="String(target.id)"
          >
            <template #title>
              <div style="display:flex;align-items:center;gap:12px;width:100%">
                <strong>{{ target.name }}</strong>
                <el-tag size="small" :type="target.cloudType === 'qiniu' ? 'warning' : target.cloudType === 'r2' ? 'primary' : 'info'">
                  {{ target.cloudType?.toUpperCase() }}
                </el-tag>
                <span v-if="target.error" style="color:#f56c6c;font-size:12px">错误: {{ target.error }}</span>
                <template v-else>
                  <span style="font-size:12px;color:#67c23a">一致 {{ target.totalRef - target.missing.length }}</span>
                  <span v-if="target.missing.length" style="font-size:12px;color:#f56c6c">缺失 {{ target.missing.length }}</span>
                  <span v-if="target.extra.length" style="font-size:12px;color:#e6a23c">多余 {{ target.extra.length }}</span>
                </template>
              </div>
            </template>

            <div v-if="target.error" style="padding:8px 0;color:#f56c6c">{{ target.error }}</div>

            <el-tabs v-else>
              <el-tab-pane :label="`缺失文件 (${target.missing.length})`" v-if="target.missing.length > 0">
                <el-table :data="target.missing" border max-height="400" size="small">
                  <el-table-column prop="key" label="文件路径（参考有，此云无）" min-width="380" show-overflow-tooltip>
                    <template #default="{ row: f }">
                      <span style="color:#f56c6c">{{ f.key }}</span>
                    </template>
                  </el-table-column>
                  <el-table-column label="大小" width="120" align="right">
                    <template #default="{ row: f }">{{ formatFileSize(f.size) }}</template>
                  </el-table-column>
                </el-table>
                <el-button type="warning" size="small" style="margin-top:8px" @click="exportMissing(target)">导出缺失列表</el-button>
              </el-tab-pane>

              <el-tab-pane :label="`多余文件 (${target.extra.length})`" v-if="target.extra.length > 0">
                <el-table :data="target.extra" border max-height="400" size="small">
                  <el-table-column prop="key" label="文件路径（此云有，参考无）" min-width="380" show-overflow-tooltip>
                    <template #default="{ row: f }">
                      <span style="color:#e6a23c">{{ f.key }}</span>
                    </template>
                  </el-table-column>
                  <el-table-column label="大小" width="120" align="right">
                    <template #default="{ row: f }">{{ formatFileSize(f.size) }}</template>
                  </el-table-column>
                </el-table>
              </el-tab-pane>

              <el-tab-pane v-if="!target.missing.length && !target.extra.length">
                <el-empty description="完全一致！" />
              </el-tab-pane>
            </el-tabs>
          </el-collapse-item>
        </el-collapse>
      </div>
      <el-empty v-else-if="compareDone" description="请选择参考云并点击比对" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import {
  getExternalLinkDomainList,
  createExternalLinkDomain,
  updateExternalLinkDomain,
  deleteExternalLinkDomain,
  setDefaultDomain,
  pingCloud,
  listCloudFiles,
  compareDirectories,
  deleteCloudFiles,
  uploadCloudFile,
  searchCloudFiles,
  getFileDownloadURL,
  getCloudFolderDownloadURL
} from '@/api/client/externalLinkDomain'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder, Search, Download } from '@element-plus/icons-vue'

const searchInfo = ref({
  page: 1,
  pageSize: 10,
  name: '',
  cloudType: '',
})

const tableData = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const pingLoading = reactive({})

const defaultForm = {
  name: '',
  domain: '',
  isDefault: false,
  isEnabled: true,
  sort: 0,
  remark: '',
  cloudType: '',
  accessKey: '',
  secretKey: '',
  bucket: '',
  regionEndpoint: '',
  zone: '',
  baseUrl: '',
  accountId: '',
  useHttps: true,
  useCdnDomain: false,
  uploadScope: 'all',
  defaultUpload: false,
  syncUpload: false,
}

const form = ref({ ...defaultForm })

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  domain: [{ required: true, message: '请输入域名地址', trigger: 'blur' }],
}

// 文件浏览器
const fileBrowserVisible = ref(false)
const fileBrowserId = ref(0)
const fileBrowserName = ref('')
const fileBrowserPrefix = ref('')
const fileBrowserList = ref([])
const fileBrowserTruncated = ref(false)
const fileBrowserCurrentPrefix = ref('')
const searchKeyword = ref('')
const searchLoading = ref(false)
const isSearchResult = ref(false)
const searchResultList = ref([])
const selectedFileKeys = ref([])

// 面包屑
const breadcrumbParts = computed(() => {
  if (!fileBrowserCurrentPrefix.value) return []
  const parts = fileBrowserCurrentPrefix.value.split('/').filter(Boolean)
  const result = []
  let path = ''
  for (const p of parts) {
    path += p + '/'
    result.push({ name: p, path })
  }
  return result
})

const filteredFileList = computed(() => {
  if (isSearchResult.value) return searchResultList.value
  return fileBrowserList.value
})

const onlyFileSelectable = (row) => !row.isDir

// 目录比对
const compareVisible = ref(false)
const compareRefId = ref(null)
const comparePrefix = ref('')
const compareLoading = ref(false)
const compareDone = ref(false)
const compareActiveCollapse = ref([])
const compareResult = ref({ referenceName: '', targets: [] })

// 可选云列表（已启用且有云类型的）
const cloudOptions = computed(() =>
  tableData.value.filter(r => r.isEnabled && r.cloudType)
)

const getList = async () => {
  const res = await getExternalLinkDomainList(searchInfo.value)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
  }
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

const openCreate = () => {
  isEdit.value = false
  form.value = { ...defaultForm }
  dialogVisible.value = true
}

const openEdit = (row) => {
  isEdit.value = true
  form.value = { ...row }
  dialogVisible.value = true
}

const handleSave = async () => {
  const fn = isEdit.value ? updateExternalLinkDomain : createExternalLinkDomain
  const res = await fn(form.value)
  if (res.code === 0) {
    ElMessage.success('保存成功')
    dialogVisible.value = false
    getList()
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该域名？', '提示', { type: 'warning' })
  const res = await deleteExternalLinkDomain({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    getList()
  }
}

const handleSetDefault = async (row) => {
  const res = await setDefaultDomain({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success('设置成功')
    getList()
  }
}

const handleToggle = async (row, val) => {
  const res = await updateExternalLinkDomain({
    ...row,
    isEnabled: val,
  })
  if (res.code === 0) {
    ElMessage.success('更新成功')
    getList()
  }
}

// 检测云存储连接
const handlePing = async (row) => {
  pingLoading[row.ID] = true
  const res = await pingCloud({ id: row.ID })
  pingLoading[row.ID] = false
  if (res.code === 0) {
    const data = res.data
    if (data.success) {
      ElMessage.success(`连接成功，延迟 ${data.latency}ms`)
    } else {
      ElMessage.warning(`连接失败: ${data.message}（延迟 ${data.latency}ms）`)
    }
  }
}

// 文件浏览器
const openFileBrowser = (row) => {
  fileBrowserId.value = row.ID
  fileBrowserName.value = row.name
  fileBrowserPrefix.value = ''
  fileBrowserList.value = []
  searchKeyword.value = ''
  searchLoading.value = false
  isSearchResult.value = false
  searchResultList.value = []
  selectedFileKeys.value = []
  fileBrowserVisible.value = true
}

const closeFileBrowser = () => {
  searchKeyword.value = ''
  isSearchResult.value = false
  searchResultList.value = []
}

const loadFileBrowserData = async () => {
  if (!fileBrowserId.value) return
  const res = await listCloudFiles({
    id: fileBrowserId.value,
    prefix: fileBrowserPrefix.value || '',
    maxKeys: 200
  })
  if (res.code === 0) {
    const prefix = res.data?.prefix || ''
    fileBrowserCurrentPrefix.value = prefix
    fileBrowserList.value = (res.data?.files || []).map(item => ({
      ...item,
      displayName: item.isDir
        ? item.key.replace(prefix, '').replace(/\/$/, '')
        : item.key.replace(prefix, '')
    }))
    fileBrowserTruncated.value = res.data?.isTruncated || false
  }
}

const navigateDir = (prefix) => {
  fileBrowserPrefix.value = prefix
  searchKeyword.value = ''
  isSearchResult.value = false
  searchResultList.value = []
  selectedFileKeys.value = []
  loadFileBrowserData()
}

const handleFileSelectionChange = (selection) => {
  selectedFileKeys.value = selection.map(item => item.key)
}

const filterFileList = () => {
  // 已被 handleSearch 替代
}

const handleSearch = async () => {
  const kw = searchKeyword.value.trim()
  if (!kw) {
    ElMessage.warning('请输入搜索关键词')
    return
  }
  searchLoading.value = true
  try {
    const res = await searchCloudFiles({
      id: fileBrowserId.value,
      keyword: kw,
      maxKeys: 300
    })
    if (res.code === 0) {
      const prefix = fileBrowserCurrentPrefix.value || ''
      searchResultList.value = (res.data?.files || []).map(item => ({
        ...item,
        displayName: item.key.replace(prefix, '')
      }))
      isSearchResult.value = true
      if (res.data?.isTruncated) {
        ElMessage.info('搜索结果过多，仅展示前 ' + searchResultList.value.length + ' 条')
      } else {
        ElMessage.success(`找到 ${searchResultList.value.length} 个文件`)
      }
    }
  } catch (e) {
    ElMessage.error('搜索失败: ' + (e.message || '未知错误'))
  } finally {
    searchLoading.value = false
  }
}

const clearSearch = () => {
  searchKeyword.value = ''
  isSearchResult.value = false
  searchResultList.value = []
  selectedFileKeys.value = []
  loadFileBrowserData()
}

const handleDownloadItem = async (row) => {
  if (row.isDir) {
    // 目录：直接拼接下载URL打开
    const url = getCloudFolderDownloadURL({
      id: fileBrowserId.value,
      prefix: row.key
    })
    window.open(url, '_blank')
    return
  }

  // 单个文件：获取签名URL后打开
  try {
    const res = await getFileDownloadURL({
      id: fileBrowserId.value,
      key: row.key
    })
    if (res.code === 0 && res.data?.downloadUrl) {
      window.open(res.data.downloadUrl, '_blank')
    }
  } catch (e) {
    ElMessage.error('获取下载链接失败: ' + (e.message || '未知错误'))
  }
}

const handleFileUpload = async (uploadFile) => {
  const formData = new FormData()
  formData.append('file', uploadFile.file)
  formData.append('id', String(fileBrowserId.value))
  formData.append('folder', fileBrowserCurrentPrefix.value || '')

  try {
    const res = await uploadCloudFile(formData)
    if (res.code === 0) {
      ElMessage.success('上传成功')
      loadFileBrowserData()
    }
  } catch (e) {
    ElMessage.error('上传失败: ' + (e.message || '未知错误'))
  }
}

const handleBatchDelete = async () => {
  if (selectedFileKeys.value.length === 0) {
    ElMessage.warning('请先选择要删除的文件')
    return
  }

  const today = getTodayStr()
  let password = ''
  try {
    const result = await ElMessageBox.prompt(
      `确定要删除选中的 ${selectedFileKeys.value.length} 个文件吗？请输入当天日期密码`,
      '确认删除',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning',
        inputPlaceholder: '例如：' + today,
        inputValidator: (val) => {
          if (!val || !/^\d{8}$/.test(val)) {
            return '密码格式不正确，请输入8位日期数字'
          }
          return true
        }
      }
    )
    password = result?.value || ''
  } catch {
    return // 用户取消
  }

  try {
    const res = await deleteCloudFiles({
      id: fileBrowserId.value,
      keys: selectedFileKeys.value,
      password
    })
    if (res.code === 0) {
      const data = res.data
      const msg = `成功删除 ${data.deletedCount} 个文件`
      if (data.failedKeys && data.failedKeys.length > 0) {
        ElMessage.warning(msg + `，${data.failedKeys.length} 个失败`)
      } else {
        ElMessage.success(msg)
      }
      selectedFileKeys.value = []
      loadFileBrowserData()
    }
  } catch (e) {
    ElMessage.error('删除失败: ' + (e.message || '未知错误'))
  }
}

const getTodayStr = () => {
  const d = new Date()
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}${m}${day}`
}

const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  let size = bytes
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  return size.toFixed(i > 0 ? 1 : 0) + ' ' + units[i]
}

// 目录比对
const openCompare = () => {
  compareRefId.value = null
  comparePrefix.value = ''
  compareResult.value = { referenceName: '', targets: [] }
  compareDone.value = false
  compareActiveCollapse.value = []
  compareVisible.value = true
}

const handleCompare = async () => {
  if (!compareRefId.value) {
    ElMessage.warning('请先选择参考基准云')
    return
  }
  compareLoading.value = true
  compareDone.value = false
  const res = await compareDirectories({
    referenceId: compareRefId.value,
    prefix: comparePrefix.value || ''
  })
  compareLoading.value = false
  if (res.code === 0) {
    compareResult.value = res.data || { referenceName: '', targets: [] }
    compareDone.value = true
    // 全部展开
    compareActiveCollapse.value = (compareResult.value.targets || []).map(t => String(t.id))
  }
}

const exportMissing = (target) => {
  const lines = target.missing.map(f => f.key)
  const blob = new Blob([lines.join('\n')], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `missing_${target.name}_vs_${compareResult.value.referenceName}.txt`
  a.click()
  URL.revokeObjectURL(url)
  ElMessage.success(`已导出 ${lines.length} 条缺失记录`)
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.file-browser-path {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 6px;
  font-size: 13px;
}
.file-browser-path .path-label {
  color: #909399;
  margin-right: 6px;
  white-space: nowrap;
}
.file-browser-path .path-value {
  color: #303133;
  font-family: 'Courier New', monospace;
  word-break: break-all;
}
.file-browser-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}
.file-breadcrumb {
  flex: 1;
  min-width: 0;
}
.file-browser-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
</style>
