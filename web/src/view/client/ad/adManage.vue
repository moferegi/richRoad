<template>
  <div>
    <div class="gva-search-box">
      <el-alert
        title="广告运营管理台：支持广告位置管理、广告视频/图片上传和云存储分发、观看记录查询。"
        type="info"
        :closable="false"
        show-icon
      />
    </div>

    <div class="gva-table-box">
      <el-tabs v-model="activeTab">
        <!-- Tab 1: 广告位置管理 -->
        <el-tab-pane label="广告位置" name="position">
          <div class="toolbar-row">
            <span class="toolbar-tip">用于维护广告展示位置（如游戏分类页看广告按钮），每个位置可绑定多个广告视频/图片</span>
            <el-button type="primary" icon="plus" @click="openPosDialog()">新增位置</el-button>
          </div>

          <el-table :data="posData" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column prop="name" label="位置名称" min-width="180" />
            <el-table-column prop="positionKey" label="位置标记" min-width="220" show-overflow-tooltip>
              <template #default="scope">
                <el-tag type="info" size="small">{{ scope.row.positionKey }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip>
              <template #default="scope">
                {{ scope.row.description || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="启用" width="80" align="center">
              <template #default="scope">
                <el-switch :model-value="scope.row.isEnabled" @change="(val) => togglePos(scope.row, val)" />
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="80" align="center" />
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openPosDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="delPos(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="posSearch.page"
              :page-size="posSearch.pageSize"
              :page-sizes="[10, 20, 50]"
              :total="posTotal"
              @current-change="(p) => { posSearch.page = p; getPosList() }"
              @size-change="(s) => { posSearch.pageSize = s; getPosList() }"
            />
          </div>
        </el-tab-pane>

        <!-- Tab 2: 广告视频管理 -->
        <el-tab-pane label="广告视频" name="video">
          <div class="toolbar-row">
            <el-form :inline="true" :model="videoSearch">
              <el-form-item label="所属位置">
                <el-select v-model="videoSearch.positionId" clearable placeholder="全部" style="width: 200px">
                  <el-option v-for="p in posData" :key="p.ID" :label="p.name" :value="p.ID" />
                </el-select>
              </el-form-item>
              <el-form-item label="媒体类型">
                <el-select v-model="videoSearch.mediaType" clearable placeholder="全部" style="width: 120px">
                  <el-option label="视频" value="video" />
                  <el-option label="图片" value="image" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="getVideoList">查询</el-button>
                <el-button icon="refresh" @click="resetVideoSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openVideoDialog()">新增视频</el-button>
          </div>

          <el-table :data="videoData" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="60" />
            <el-table-column label="缩略图" width="90">
              <template #default="scope">
                <el-image v-if="scope.row.thumbnailUrl" style="width: 50px; height: 50px" :src="scope.row.thumbnailUrl" fit="cover" :preview-src-list="[scope.row.thumbnailUrl]" preview-teleported />
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="title" label="标题" min-width="160" show-overflow-tooltip />
            <el-table-column label="媒体类型" width="90" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.mediaType === 'video' ? 'danger' : 'success'" size="small">
                  {{ scope.row.mediaType === 'video' ? '视频' : '图片' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="媒体资源" min-width="200" show-overflow-tooltip>
              <template #default="scope">
                <el-link v-if="scope.row.mediaUrl" type="primary" :underline="false">{{ scope.row.mediaUrl }}</el-link>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="duration" label="总时长(秒)" width="100" align="center" />
            <el-table-column prop="minWatchSeconds" label="最低观看(秒)" width="110" align="center" />
            <el-table-column label="启用" width="80" align="center">
              <template #default="scope">
                <el-switch :model-value="scope.row.isEnabled" @change="(val) => toggleVideo(scope.row, val)" />
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="70" align="center" />
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openVideoDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="delVideo(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="videoSearch.page"
              :page-size="videoSearch.pageSize"
              :page-sizes="[10, 20, 50]"
              :total="videoTotal"
              @current-change="(p) => { videoSearch.page = p; getVideoList() }"
              @size-change="(s) => { videoSearch.pageSize = s; getVideoList() }"
            />
          </div>
        </el-tab-pane>

        <!-- Tab 3: 观看记录 -->
        <el-tab-pane label="观看记录" name="record">
          <div class="toolbar-row">
            <el-form :inline="true" :model="recordSearch">
              <el-form-item label="用户ID">
                <el-input v-model="recordSearch.userId" placeholder="用户ID" clearable style="width: 180px" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="getRecordList">查询</el-button>
                <el-button icon="refresh" @click="resetRecordSearch">重置</el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-table :data="recordData" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="70" />
            <el-table-column prop="userID" label="用户ID" width="100" />
            <el-table-column prop="adVideoID" label="广告视频ID" width="110" />
            <el-table-column prop="positionID" label="位置ID" width="90" />
            <el-table-column prop="watchedSeconds" label="观看秒数" width="100" align="center" />
            <el-table-column label="完成" width="80" align="center">
              <template #default="scope">
                <el-tag :type="scope.row.isCompleted ? 'success' : 'info'" size="small">
                  {{ scope.row.isCompleted ? '是' : '否' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="ip" label="IP" width="150" />
            <el-table-column label="时间" width="180">
              <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="recordSearch.page"
              :page-size="recordSearch.pageSize"
              :page-sizes="[10, 20, 50]"
              :total="recordTotal"
              @current-change="(p) => { recordSearch.page = p; getRecordList() }"
              @size-change="(s) => { recordSearch.pageSize = s; getRecordList() }"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 位置新增/编辑弹窗 -->
    <el-dialog v-model="posDialog" :title="posIsEdit ? '编辑位置' : '新增位置'" width="500px" append-to-body>
      <el-form :model="posForm" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="posForm.name" placeholder="如：游戏分类页-看广告" />
        </el-form-item>
        <el-form-item label="位置标记" required>
          <el-input v-model="posForm.positionKey" placeholder="如：pages/game/category" />
          <div style="color: #909399; font-size: 12px; margin-top: 4px">
            Uni 端通过此标记查询广告，例如 <code>pages/game/category</code>
          </div>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="posForm.description" type="textarea" :rows="2" placeholder="位置描述" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="posForm.isEnabled" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="posForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="posDialog = false">取消</el-button>
        <el-button type="primary" @click="savePos">保存</el-button>
      </template>
    </el-dialog>

    <!-- 视频新增/编辑弹窗 -->
    <el-dialog v-model="videoDialog" :title="videoIsEdit ? '编辑视频' : '新增视频'" width="650px" append-to-body>
      <el-form :model="videoForm" label-width="120px">
        <el-form-item label="所属位置" required>
          <el-select v-model="videoForm.positionId" placeholder="选择位置" style="width: 100%">
            <el-option v-for="p in posData" :key="p.ID" :label="p.name" :value="p.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" required>
          <el-input v-model="videoForm.title" placeholder="广告标题" />
        </el-form-item>
        <el-form-item label="媒体类型" required>
          <el-radio-group v-model="videoForm.mediaType">
            <el-radio value="video">视频</el-radio>
            <el-radio value="image">图片</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 视频上传 -->
        <template v-if="videoForm.mediaType === 'video'">
          <el-form-item label="切片模式">
            <el-switch
              v-model="hlsMode"
              active-text="HLS切片"
              inactive-text="普通上传"
              @change="onHlsModeChange"
            />
            <span v-if="hlsMode && ffmpegChecking" style="margin-left:12px;color:#409eff;">
              ⏳ 正在检测 FFmpeg...
            </span>
            <el-tag v-else-if="hlsMode && ffmpegReady" type="success" style="margin-left:12px;">FFmpeg 就绪</el-tag>
            <el-tag v-else-if="hlsMode && !ffmpegReady" type="danger" style="margin-left:12px;">FFmpeg 不可用</el-tag>
          </el-form-item>
          <el-form-item v-if="!hlsMode" label="视频地址">
            <FileUploadWithDir v-model="videoForm.mediaUrl" :default-folder="videoForm.uploadFolder || 'ad/video'" accept="video/*" />
          </el-form-item>
          <el-form-item v-else label="HLS切片">
            <div class="hls-folder-bar">
              <span class="hls-folder-label">切片目录：</span>
              <el-input
                v-model="videoForm.uploadFolder"
                size="small"
                class="hls-folder-input"
                placeholder="如 ad/video，留空自动生成"
                clearable
              />
            </div>
            <div class="hls-folder-bar" style="margin-top: 8px;">
              <span class="hls-folder-label">视频文件：</span>
              <input
                ref="hlsFileInput"
                type="file"
                accept="video/*"
                style="flex:1;"
                @change="onHlsFileChange"
              />
              <el-button
                type="primary"
                style="margin-left: 8px;"
                :loading="hlsSlicing"
                :disabled="!ffmpegReady || !hlsFile"
                @click="onHlsSliceClick"
              >
                {{ hlsSlicing ? '正在切片上传...' : '开始切片' }}
              </el-button>
            </div>
            <div v-if="hlsProgress" style="margin-top:8px;color:#409eff;">{{ hlsProgress }}</div>
            <div v-if="videoForm.mediaUrl" style="margin-top:8px;color:#67c23a;font-size:12px;">
              已切片: {{ videoForm.mediaUrl }}
            </div>
          </el-form-item>
        </template>

        <!-- 图片上传 -->
        <el-form-item v-else label="图片地址">
          <FileUploadWithDir v-model="videoForm.mediaUrl" :default-folder="videoForm.uploadFolder || 'ad/image'" accept="image/*" />
        </el-form-item>

        <el-form-item label="封面缩略图">
          <FileUploadWithDir v-model="videoForm.thumbnailUrl" :default-folder="videoForm.uploadFolder || 'ad/thumb'" accept="image/*" />
        </el-form-item>
        <el-form-item label="总时长(秒)" required>
          <el-input-number v-model="videoForm.duration" :min="1" />
          <span style="margin-left: 8px; color: #909399; font-size: 12px">视频切片后自动获取，图片请手动设置</span>
        </el-form-item>
        <el-form-item label="最低观看(秒)" required>
          <el-input-number v-model="videoForm.minWatchSeconds" :min="1" :max="videoForm.duration" />
          <span style="margin-left: 8px; color: #909399; font-size: 12px">用户必须看到此时长后才能关闭广告</span>
        </el-form-item>
        <el-form-item label="跳转链接">
          <el-input v-model="videoForm.linkUrl" placeholder="可选，广告跳转链接" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="videoForm.isEnabled" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="videoForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="videoDialog = false">取消</el-button>
        <el-button type="primary" @click="saveVideo">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import {
  getAdPositionList, createAdPosition, updateAdPosition, deleteAdPosition,
  getAdVideoList, createAdVideo, updateAdVideo, deleteAdVideo,
  sliceAdVideo, checkFfmpeg,
  getWatchRecordList
} from '@/api/client/ad'
import FileUploadWithDir from '@/components/FileUploadWithDir/index.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const activeTab = ref('position')

// ===== 位置管理 =====
const posSearch = reactive({ name: '', page: 1, pageSize: 10 })
const posData = ref([])
const posTotal = ref(0)
const posDialog = ref(false)
const posIsEdit = ref(false)
const posForm = reactive({ ID: 0, name: '', positionKey: '', description: '', isEnabled: true, sort: 0 })

const getPosList = async () => {
  try {
    const res = await getAdPositionList(posSearch)
    if (res.code === 0) {
      posData.value = res.data.list || []
      posTotal.value = res.data.total || 0
    }
  } catch (e) { /* ignore */ }
}

const openPosDialog = (row) => {
  if (row) {
    posIsEdit.value = true
    Object.assign(posForm, {
      ID: row.ID, name: row.name, positionKey: row.positionKey,
      description: row.description, isEnabled: row.isEnabled, sort: row.sort
    })
  } else {
    posIsEdit.value = false
    Object.assign(posForm, { ID: 0, name: '', positionKey: '', description: '', isEnabled: true, sort: 0 })
  }
  posDialog.value = true
}

const savePos = async () => {
  try {
    const api = posIsEdit.value ? updateAdPosition : createAdPosition
    const res = await api(posForm)
    if (res.code === 0) {
      ElMessage.success(posIsEdit.value ? '更新成功' : '创建成功')
      posDialog.value = false
      getPosList()
    } else {
      ElMessage.error(res.msg || '操作失败')
    }
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

const delPos = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该位置？', '确认', { type: 'warning' })
    const res = await deleteAdPosition({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getPosList()
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (e) { /* cancel */ }
}

const togglePos = async (row, val) => {
  try {
    await updateAdPosition({ ID: row.ID, isEnabled: val })
    row.isEnabled = val
    ElMessage.success(val ? '已启用' : '已禁用')
  } catch (e) { ElMessage.error('操作失败') }
}

// ===== 视频管理 =====
const videoSearch = reactive({ positionId: 0, mediaType: '', page: 1, pageSize: 10 })
const videoData = ref([])
const videoTotal = ref(0)
const videoDialog = ref(false)
const videoIsEdit = ref(false)
const videoForm = reactive({
  ID: 0, positionId: 0, title: '', mediaType: 'video',
  mediaUrl: '', thumbnailUrl: '', duration: 0, minWatchSeconds: 5,
  isEnabled: true, sort: 0, uploadFolder: 'ad/video', linkUrl: ''
})

// HLS 切片相关
const hlsMode = ref(false)
const ffmpegChecking = ref(false)
const ffmpegReady = ref(false)
const hlsFile = ref(null)
const hlsFileInput = ref(null)
const hlsSlicing = ref(false)
const hlsProgress = ref('')

const getVideoList = async () => {
  try {
    const res = await getAdVideoList(videoSearch)
    if (res.code === 0) {
      videoData.value = res.data.list || []
      videoTotal.value = res.data.total || 0
    }
  } catch (e) { /* ignore */ }
}

const resetVideoSearch = () => {
  videoSearch.positionId = 0
  videoSearch.mediaType = ''
  videoSearch.page = 1
  getVideoList()
}

const openVideoDialog = (row) => {
  hlsMode.value = false
  ffmpegReady.value = false
  ffmpegChecking.value = false
  hlsFile.value = null
  hlsProgress.value = ''
  if (hlsFileInput.value) {
    hlsFileInput.value.value = ''
  }
  if (row) {
    videoIsEdit.value = true
    Object.assign(videoForm, {
      ID: row.ID, positionId: row.positionId, title: row.title, mediaType: row.mediaType,
      mediaUrl: row.mediaUrl || '', thumbnailUrl: row.thumbnailUrl || '',
      duration: row.duration, minWatchSeconds: row.minWatchSeconds,
      isEnabled: row.isEnabled, sort: row.sort, uploadFolder: row.uploadFolder || 'ad/video',
      linkUrl: row.linkUrl || ''
    })
  } else {
    videoIsEdit.value = false
    Object.assign(videoForm, {
      ID: 0, positionId: 0, title: '', mediaType: 'video',
      mediaUrl: '', thumbnailUrl: '', duration: 0, minWatchSeconds: 5,
      isEnabled: true, sort: 0, uploadFolder: 'ad/video', linkUrl: ''
    })
  }
  videoDialog.value = true
}

const saveVideo = async () => {
  try {
    const api = videoIsEdit.value ? updateAdVideo : createAdVideo
    const res = await api(videoForm)
    if (res.code === 0) {
      ElMessage.success(videoIsEdit.value ? '更新成功' : '创建成功')
      videoDialog.value = false
      getVideoList()
    } else {
      ElMessage.error(res.msg || '操作失败')
    }
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

const delVideo = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该广告视频？', '确认', { type: 'warning' })
    const res = await deleteAdVideo({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getVideoList()
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (e) { /* cancel */ }
}

const toggleVideo = async (row, val) => {
  try {
    await updateAdVideo({ ID: row.ID, isEnabled: val })
    row.isEnabled = val
    ElMessage.success(val ? '已启用' : '已禁用')
  } catch (e) { ElMessage.error('操作失败') }
}

// HLS 切片模式
const onHlsModeChange = async (value) => {
  if (!value) {
    ffmpegReady.value = false
    ffmpegChecking.value = false
    hlsProgress.value = ''
    return
  }
  ffmpegChecking.value = true
  ffmpegReady.value = false
  try {
    const res = await checkFfmpeg()
    if (res.code === 0) {
      ffmpegReady.value = true
    } else {
      ffmpegReady.value = false
      ElMessage.warning('FFmpeg 不可用，无法使用 HLS 切片功能')
    }
  } catch (e) {
    ffmpegReady.value = false
    ElMessage.warning('FFmpeg 检测失败，请确认服务端已安装 FFmpeg')
  } finally {
    ffmpegChecking.value = false
  }
}

const onHlsFileChange = (e) => {
  hlsFile.value = e.target.files[0] || null
}

const onHlsSliceClick = async () => {
  if (!hlsFile.value) {
    ElMessage.warning('请选择视频文件')
    return
  }
  if (!videoForm.positionId) {
    ElMessage.warning('请先选择所属位置')
    return
  }
  if (!videoForm.title) {
    ElMessage.warning('标题不能为空')
    return
  }

  hlsSlicing.value = true
  hlsProgress.value = '正在切片上传中，请稍候（大视频可能需要几分钟）...'

  const formData = new FormData()
  formData.append('file', hlsFile.value)
  if (videoForm.uploadFolder) {
    formData.append('folder', videoForm.uploadFolder)
  }

  try {
    const res = await sliceAdVideo(formData)
    if (res.code !== 0) {
      ElMessage.error('HLS 切片失败: ' + (res.msg || '未知错误'))
      return
    }
    ElMessage.success(res.msg || 'HLS 切片上传成功')
    if (res.data?.mediaUrl) {
      videoForm.mediaUrl = res.data.mediaUrl
    }
    if (res.data?.duration) {
      videoForm.duration = Math.round(res.data.duration)
    }
  } catch (e) {
    ElMessage.error('HLS 切片请求失败，请检查网络或服务端状态')
  } finally {
    hlsSlicing.value = false
    hlsProgress.value = ''
  }
}

// ===== 观看记录 =====
const recordSearch = reactive({ userId: 0, page: 1, pageSize: 10 })
const recordData = ref([])
const recordTotal = ref(0)

const getRecordList = async () => {
  try {
    const res = await getWatchRecordList(recordSearch)
    if (res.code === 0) {
      recordData.value = res.data.list || []
      recordTotal.value = res.data.total || 0
    }
  } catch (e) { /* ignore */ }
}

const resetRecordSearch = () => {
  recordSearch.userId = 0
  recordSearch.page = 1
  getRecordList()
}

const formatDate = (d) => {
  if (!d) return '-'
  const dt = new Date(d)
  const pad = (n) => String(n).padStart(2, '0')
  return `${dt.getFullYear()}-${pad(dt.getMonth() + 1)}-${pad(dt.getDate())} ${pad(dt.getHours())}:${pad(dt.getMinutes())}:${pad(dt.getSeconds())}`
}

onMounted(() => {
  getPosList()
})
</script>

<style scoped>
.toolbar-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 8px;
}
.toolbar-tip {
  color: #909399;
  font-size: 13px;
}
.hls-folder-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}
.hls-folder-label {
  white-space: nowrap;
  color: #606266;
}
.hls-folder-input {
  width: 300px;
}
.gva-search-box { padding: 16px; background: #fff; border-radius: 4px; margin-bottom: 16px; }
.gva-table-box { padding: 16px; background: #fff; border-radius: 4px; }
.gva-pagination { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>