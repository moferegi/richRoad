<template>
  <div>
    <div class="gva-search-box">
      <el-alert
        title="英语学习运营台（视频端）：支持视频分类、剧集、单集和用户资源授权管理。"
        type="info"
        :closable="false"
        show-icon
      />
    </div>

    <div class="gva-table-box">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="视频分类" name="videoCategory">
          <div class="toolbar-row">
            <span class="toolbar-tip">用于维护视频大类（如影视英语、情景口语等）</span>
            <el-button type="primary" icon="plus" @click="openVideoCategoryDialog()">新增视频分类</el-button>
          </div>

          <el-table :data="videoCategoryTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="分类名称" min-width="240">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openVideoCategoryDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeVideoCategory(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="videoCategoryQuery.page"
              :page-size="videoCategoryQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="videoCategoryTotal"
              @current-change="handleVideoCategoryPageChange"
              @size-change="handleVideoCategorySizeChange"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="剧集管理" name="series">
          <div class="toolbar-row">
            <el-form :inline="true" :model="seriesQuery">
              <el-form-item label="所属分类">
                <el-select v-model="seriesQuery.categoryId" clearable placeholder="全部" style="width: 240px">
                  <el-option
                    v-for="item in videoCategoryOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleSeriesSearch">查询</el-button>
                <el-button icon="refresh" @click="resetSeriesSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openSeriesDialog()">新增剧集</el-button>
          </div>

          <el-table :data="seriesTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="所属分类" min-width="180">
              <template #default="scope">
                {{ videoCategoryNameMap[scope.row.categoryId] || `#${scope.row.categoryId}` }}
              </template>
            </el-table-column>
            <el-table-column label="剧集名称" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="coverUrl" label="封面地址" min-width="240" show-overflow-tooltip />
            <el-table-column prop="price" label="价格" width="120" />
            <el-table-column label="会员限制" width="120">
              <template #default="scope">
                <el-tag :type="scope.row.needVip ? 'warning' : 'success'">{{ scope.row.needVip ? '需要会员' : '普通可见' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="viewCount" label="浏览量" width="100" />
            <el-table-column prop="userCount" label="观看人数" width="110" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openSeriesDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeSeries(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="seriesQuery.page"
              :page-size="seriesQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="seriesTotal"
              @current-change="handleSeriesPageChange"
              @size-change="handleSeriesSizeChange"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane label="单集与字幕" name="episode">
          <div class="toolbar-row">
            <el-form :inline="true" :model="episodeQuery">
              <el-form-item label="所属剧集">
                <el-select v-model="episodeQuery.seriesId" clearable filterable placeholder="全部" style="width: 280px">
                  <el-option
                    v-for="item in seriesOptions"
                    :key="item.ID"
                    :label="formatI18nText(item.name)"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleEpisodeSearch">查询</el-button>
                <el-button icon="refresh" @click="resetEpisodeSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openEpisodeDialog()">新增单集</el-button>
          </div>

          <el-table :data="episodeTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column label="所属剧集" min-width="220">
              <template #default="scope">
                {{ seriesNameMap[scope.row.seriesId] || `#${scope.row.seriesId}` }}
              </template>
            </el-table-column>
            <el-table-column label="单集名称" min-width="220">
              <template #default="scope">
                {{ formatI18nText(scope.row.name) }}
              </template>
            </el-table-column>
            <el-table-column prop="videoUrl" label="视频地址" min-width="260" show-overflow-tooltip />
            <el-table-column prop="trialPercent" label="试看比例(%)" width="120" />
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="openEpisodeDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link icon="delete" @click="removeEpisode(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="episodeQuery.page"
              :page-size="episodeQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="episodeTotal"
              @current-change="handleEpisodePageChange"
              @size-change="handleEpisodeSizeChange"
            />
          </div>

          <div class="subtitle-box">
            <el-divider content-position="left">字幕解析入库</el-divider>
            <el-form label-width="120px">
              <el-form-item label="目标单集">
                <el-select v-model="subtitleForm.episodeId" filterable placeholder="请选择单集" style="width: 420px">
                  <el-option
                    v-for="item in episodeOptions"
                    :key="item.ID"
                    :label="`${seriesNameMap[item.seriesId] || '-'} / ${formatI18nText(item.name)}`"
                    :value="item.ID"
                  />
                </el-select>
              </el-form-item>
              <el-form-item label="字幕 JSON">
                <el-input
                  v-model="subtitleForm.rawJsonData"
                  type="textarea"
                  :rows="10"
                  placeholder='[{"startTime":0,"endTime":5.5,"english":"You are my destiny","translate":"{\"zh\":\"你是我的命运\"}"}]'
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleParseSubtitle">提交字幕解析</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane label="用户授权" name="entitlement">
          <div class="toolbar-row">
            <el-form :inline="true" :model="entitlementQuery">
              <el-form-item label="用户ID">
                <el-input-number v-model="entitlementQuery.userId" :min="1" controls-position="right" />
              </el-form-item>
              <el-form-item label="资源类型">
                <el-select v-model="entitlementQuery.resourceType" clearable placeholder="全部" style="width: 220px">
                  <el-option label="分类" value="english_category" />
                  <el-option label="剧集" value="video_series" />
                  <el-option label="单集" value="video_episode" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" icon="search" @click="handleEntitlementSearch">查询</el-button>
                <el-button icon="refresh" @click="resetEntitlementSearch">重置</el-button>
              </el-form-item>
            </el-form>
            <el-button type="primary" icon="plus" @click="openEntitlementDialog">新增授权</el-button>
          </div>

          <el-table :data="entitlementTable" row-key="ID" border>
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column prop="userId" label="用户ID" width="100" />
            <el-table-column label="资源类型" width="140">
              <template #default="scope">
                <el-tag>{{ resourceTypeLabel(scope.row.resourceType) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="resourceId" label="资源ID" width="110" />
            <el-table-column prop="source" label="授权来源" width="140" />
            <el-table-column prop="grantedBy" label="授权人" width="100" />
            <el-table-column label="过期时间" min-width="180">
              <template #default="scope">
                {{ scope.row.expireAt || '永久有效' }}
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="200" show-overflow-tooltip />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="scope">
                <el-button type="danger" link icon="delete" @click="removeEntitlement(scope.row)">撤销</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              layout="total, sizes, prev, pager, next, jumper"
              :current-page="entitlementQuery.page"
              :page-size="entitlementQuery.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="entitlementTotal"
              @current-change="handleEntitlementPageChange"
              @size-change="handleEntitlementSizeChange"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <el-dialog v-model="videoCategoryDialogVisible" :title="videoCategoryDialogMode === 'create' ? '新增视频分类' : '编辑视频分类'" width="640px">
      <el-form :model="videoCategoryForm" label-width="120px">
        <el-form-item label="分类名称(JSON)">
          <el-input v-model="videoCategoryForm.name" type="textarea" :rows="3" placeholder='例如: {"zh":"影视英语","en":"Movie English"}' />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="videoCategoryForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="videoCategoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitVideoCategory">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="seriesDialogVisible" :title="seriesDialogMode === 'create' ? '新增剧集' : '编辑剧集'" width="760px">
      <el-form :model="seriesForm" label-width="120px">
        <el-form-item label="所属分类">
          <el-select v-model="seriesForm.categoryId" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="item in videoCategoryOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="剧集名称(JSON)">
          <el-input v-model="seriesForm.name" type="textarea" :rows="3" placeholder='例如: {"zh":"老友记 S1","en":"Friends S1"}' />
        </el-form-item>
        <el-form-item label="封面地址">
          <el-input v-model="seriesForm.coverUrl" placeholder="可填 OSS/CDN 图片地址" />
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="seriesForm.price" :min="0" :precision="2" :step="1" />
        </el-form-item>
        <el-form-item label="需要会员">
          <el-switch v-model="seriesForm.needVip" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="seriesDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitSeries">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="episodeDialogVisible" :title="episodeDialogMode === 'create' ? '新增单集' : '编辑单集'" width="760px">
      <el-form :model="episodeForm" label-width="120px">
        <el-form-item label="所属剧集">
          <el-select v-model="episodeForm.seriesId" filterable placeholder="请选择剧集" style="width: 100%">
            <el-option
              v-for="item in seriesOptions"
              :key="item.ID"
              :label="formatI18nText(item.name)"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="单集名称(JSON)">
          <el-input v-model="episodeForm.name" type="textarea" :rows="3" placeholder='例如: {"zh":"第1集","en":"Episode 1"}' />
        </el-form-item>
        <el-form-item label="视频地址">
          <el-input v-model="episodeForm.videoUrl" placeholder="m3u8/mp4 地址" />
        </el-form-item>
        <el-form-item label="试看比例(%)">
          <el-input-number v-model="episodeForm.trialPercent" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="episodeForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="episodeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEpisode">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="entitlementDialogVisible" title="新增资源授权" width="760px">
      <el-form :model="entitlementForm" label-width="120px">
        <el-form-item label="用户ID">
          <el-input-number v-model="entitlementForm.userId" :min="1" controls-position="right" />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="entitlementForm.resourceType" style="width: 100%">
            <el-option label="分类" value="english_category" />
            <el-option label="剧集" value="video_series" />
            <el-option label="单集" value="video_episode" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源ID">
          <el-select
            v-model="entitlementForm.resourceId"
            filterable
            clearable
            placeholder="可从列表选择，或手动输入"
            style="width: 100%"
          >
            <el-option
              v-for="item in resourceOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <el-input-number v-model="entitlementForm.resourceId" :min="1" class="mt-2 w-full" controls-position="right" />
        </el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker
            v-model="entitlementForm.expireAt"
            type="datetime"
            placeholder="不填表示永久"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            style="width: 100%"
            clearable
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="entitlementForm.remark" type="textarea" :rows="3" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="entitlementDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEntitlement">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    createVideoCategory,
    createVideoEpisode,
    createVideoSeries,
    deleteVideoCategory,
    deleteVideoEpisode,
    deleteVideoSeries,
    getEntitlementList,
    getVideoCategoryList,
    getVideoEpisodeList,
    getVideoSeriesList,
    grantEntitlement,
    parseSubtitle,
    revokeEntitlement,
    updateVideoCategory,
    updateVideoEpisode,
    updateVideoSeries
  } from '../api/english'

  defineOptions({
    name: 'EnglishLearningVideo'
  })

  const activeTab = ref('videoCategory')

  const formatI18nText = (raw) => {
    if (!raw) return ''
    try {
      const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
      if (typeof parsed === 'object' && parsed !== null) {
        return parsed.zh || parsed.en || Object.values(parsed)[0] || ''
      }
      return String(raw)
    } catch (e) {
      return String(raw)
    }
  }

  const ensureJsonText = (raw) => {
    const text = String(raw || '').trim()
    if (!text) {
      return '{}'
    }
    JSON.parse(text)
    return text
  }

  const videoCategoryQuery = ref({ page: 1, pageSize: 10 })
  const videoCategoryTable = ref([])
  const videoCategoryTotal = ref(0)

  const seriesQuery = ref({ page: 1, pageSize: 10, categoryId: undefined })
  const seriesTable = ref([])
  const seriesTotal = ref(0)

  const episodeQuery = ref({ page: 1, pageSize: 10, seriesId: undefined })
  const episodeTable = ref([])
  const episodeTotal = ref(0)

  const entitlementQuery = ref({ page: 1, pageSize: 10, userId: undefined, resourceType: undefined })
  const entitlementTable = ref([])
  const entitlementTotal = ref(0)

  const videoCategoryOptions = ref([])
  const seriesOptions = ref([])
  const episodeOptions = ref([])

  const videoCategoryNameMap = computed(() => {
    const map = {}
    for (const item of videoCategoryOptions.value) {
      map[item.ID] = formatI18nText(item.name)
    }
    return map
  })

  const seriesNameMap = computed(() => {
    const map = {}
    for (const item of seriesOptions.value) {
      map[item.ID] = formatI18nText(item.name)
    }
    return map
  })

  const videoCategoryDialogVisible = ref(false)
  const videoCategoryDialogMode = ref('create')
  const videoCategoryForm = ref({ ID: 0, name: '{}', sort: 0 })

  const seriesDialogVisible = ref(false)
  const seriesDialogMode = ref('create')
  const seriesForm = ref({
    ID: 0,
    categoryId: undefined,
    name: '{}',
    coverId: 0,
    coverUrl: '',
    price: 0,
    needVip: false
  })

  const episodeDialogVisible = ref(false)
  const episodeDialogMode = ref('create')
  const episodeForm = ref({
    ID: 0,
    seriesId: undefined,
    name: '{}',
    videoUrl: '',
    trialPercent: 8,
    sort: 0
  })

  const subtitleForm = ref({
    episodeId: undefined,
    rawJsonData: '[\n  {\n    "startTime": 0,\n    "endTime": 5.5,\n    "english": "You are my destiny",\n    "translate": "{\\"zh\\":\\"你是我的命运\\"}"\n  }\n]'
  })

  const entitlementDialogVisible = ref(false)
  const entitlementForm = ref({
    userId: undefined,
    resourceType: 'video_series',
    resourceId: undefined,
    expireAt: '',
    remark: ''
  })

  const resourceOptions = computed(() => {
    if (entitlementForm.value.resourceType === 'english_category') {
      return videoCategoryOptions.value.map((item) => ({
        value: item.ID,
        label: `${item.ID} - ${formatI18nText(item.name)}`
      }))
    }
    if (entitlementForm.value.resourceType === 'video_series') {
      return seriesOptions.value.map((item) => ({
        value: item.ID,
        label: `${item.ID} - ${formatI18nText(item.name)}`
      }))
    }
    return episodeOptions.value.map((item) => ({
      value: item.ID,
      label: `${item.ID} - ${seriesNameMap.value[item.seriesId] || '-'} / ${formatI18nText(item.name)}`
    }))
  })

  const resourceTypeLabel = (type) => {
    if (type === 'english_category') return '分类'
    if (type === 'video_series') return '剧集'
    if (type === 'video_episode') return '单集'
    return type || '-'
  }

  const loadVideoCategoryList = async () => {
    const res = await getVideoCategoryList({ ...videoCategoryQuery.value })
    if (res.code !== 0) return
    videoCategoryTable.value = res.data?.list || []
    videoCategoryTotal.value = Number(res.data?.total || 0)
    videoCategoryQuery.value.page = Number(res.data?.page || videoCategoryQuery.value.page)
    videoCategoryQuery.value.pageSize = Number(res.data?.pageSize || videoCategoryQuery.value.pageSize)
  }

  const loadSeriesList = async () => {
    const res = await getVideoSeriesList({ ...seriesQuery.value })
    if (res.code !== 0) return
    seriesTable.value = res.data?.list || []
    seriesTotal.value = Number(res.data?.total || 0)
    seriesQuery.value.page = Number(res.data?.page || seriesQuery.value.page)
    seriesQuery.value.pageSize = Number(res.data?.pageSize || seriesQuery.value.pageSize)
  }

  const loadEpisodeList = async () => {
    const res = await getVideoEpisodeList({ ...episodeQuery.value })
    if (res.code !== 0) return
    episodeTable.value = res.data?.list || []
    episodeTotal.value = Number(res.data?.total || 0)
    episodeQuery.value.page = Number(res.data?.page || episodeQuery.value.page)
    episodeQuery.value.pageSize = Number(res.data?.pageSize || episodeQuery.value.pageSize)
  }

  const loadEntitlementList = async () => {
    const params = {
      page: entitlementQuery.value.page,
      pageSize: entitlementQuery.value.pageSize
    }
    if (entitlementQuery.value.userId) {
      params.userId = entitlementQuery.value.userId
    }
    if (entitlementQuery.value.resourceType) {
      params.resourceType = entitlementQuery.value.resourceType
    }

    const res = await getEntitlementList(params)
    if (res.code !== 0) return
    entitlementTable.value = res.data?.list || []
    entitlementTotal.value = Number(res.data?.total || 0)
    entitlementQuery.value.page = Number(res.data?.page || entitlementQuery.value.page)
    entitlementQuery.value.pageSize = Number(res.data?.pageSize || entitlementQuery.value.pageSize)
  }

  const loadVideoCategoryOptions = async () => {
    const res = await getVideoCategoryList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    videoCategoryOptions.value = res.data?.list || []
  }

  const loadSeriesOptions = async () => {
    const res = await getVideoSeriesList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    seriesOptions.value = res.data?.list || []
  }

  const loadEpisodeOptions = async () => {
    const res = await getVideoEpisodeList({ page: 1, pageSize: 1000 })
    if (res.code !== 0) return
    episodeOptions.value = res.data?.list || []
  }

  const handleVideoCategoryPageChange = (page) => {
    videoCategoryQuery.value.page = page
    loadVideoCategoryList()
  }

  const handleVideoCategorySizeChange = (size) => {
    videoCategoryQuery.value.pageSize = size
    videoCategoryQuery.value.page = 1
    loadVideoCategoryList()
  }

  const handleSeriesSearch = () => {
    seriesQuery.value.page = 1
    loadSeriesList()
  }

  const resetSeriesSearch = () => {
    seriesQuery.value = { page: 1, pageSize: 10, categoryId: undefined }
    loadSeriesList()
  }

  const handleSeriesPageChange = (page) => {
    seriesQuery.value.page = page
    loadSeriesList()
  }

  const handleSeriesSizeChange = (size) => {
    seriesQuery.value.pageSize = size
    seriesQuery.value.page = 1
    loadSeriesList()
  }

  const handleEpisodeSearch = () => {
    episodeQuery.value.page = 1
    loadEpisodeList()
  }

  const resetEpisodeSearch = () => {
    episodeQuery.value = { page: 1, pageSize: 10, seriesId: undefined }
    loadEpisodeList()
  }

  const handleEpisodePageChange = (page) => {
    episodeQuery.value.page = page
    loadEpisodeList()
  }

  const handleEpisodeSizeChange = (size) => {
    episodeQuery.value.pageSize = size
    episodeQuery.value.page = 1
    loadEpisodeList()
  }

  const handleEntitlementSearch = () => {
    entitlementQuery.value.page = 1
    loadEntitlementList()
  }

  const resetEntitlementSearch = () => {
    entitlementQuery.value = { page: 1, pageSize: 10, userId: undefined, resourceType: undefined }
    loadEntitlementList()
  }

  const handleEntitlementPageChange = (page) => {
    entitlementQuery.value.page = page
    loadEntitlementList()
  }

  const handleEntitlementSizeChange = (size) => {
    entitlementQuery.value.pageSize = size
    entitlementQuery.value.page = 1
    loadEntitlementList()
  }

  const openVideoCategoryDialog = (row) => {
    videoCategoryDialogMode.value = row?.ID ? 'edit' : 'create'
    videoCategoryForm.value = {
      ID: row?.ID || 0,
      name: row?.name || '{}',
      sort: Number(row?.sort || 0)
    }
    videoCategoryDialogVisible.value = true
  }

  const submitVideoCategory = async () => {
    if (!videoCategoryForm.value.name) {
      ElMessage.warning('分类名称不能为空')
      return
    }
    try {
      videoCategoryForm.value.name = ensureJsonText(videoCategoryForm.value.name)
    } catch (e) {
      ElMessage.error('分类名称必须是合法 JSON')
      return
    }

    const payload = {
      ID: videoCategoryForm.value.ID,
      name: videoCategoryForm.value.name,
      sort: Number(videoCategoryForm.value.sort || 0)
    }
    const res = videoCategoryDialogMode.value === 'create'
      ? await createVideoCategory(payload)
      : await updateVideoCategory(payload)

    if (res.code !== 0) return
    ElMessage.success(videoCategoryDialogMode.value === 'create' ? '视频分类创建成功' : '视频分类更新成功')
    videoCategoryDialogVisible.value = false
    await Promise.all([loadVideoCategoryList(), loadVideoCategoryOptions()])
  }

  const removeVideoCategory = (row) => {
    ElMessageBox.confirm(`确认删除视频分类【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoCategory({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('视频分类删除成功')
      await Promise.all([
        loadVideoCategoryList(),
        loadVideoCategoryOptions(),
        loadSeriesList(),
        loadSeriesOptions(),
        loadEpisodeList(),
        loadEpisodeOptions()
      ])
    })
  }

  const openSeriesDialog = (row) => {
    seriesDialogMode.value = row?.ID ? 'edit' : 'create'
    seriesForm.value = {
      ID: row?.ID || 0,
      categoryId: row?.categoryId || seriesQuery.value.categoryId || undefined,
      name: row?.name || '{}',
      coverId: Number(row?.coverId || 0),
      coverUrl: row?.coverUrl || '',
      price: Number(row?.price || 0),
      needVip: !!row?.needVip
    }
    seriesDialogVisible.value = true
  }

  const submitSeries = async () => {
    if (!seriesForm.value.categoryId) {
      ElMessage.warning('请选择所属分类')
      return
    }
    if (!seriesForm.value.name) {
      ElMessage.warning('剧集名称不能为空')
      return
    }
    try {
      seriesForm.value.name = ensureJsonText(seriesForm.value.name)
    } catch (e) {
      ElMessage.error('剧集名称必须是合法 JSON')
      return
    }

    const payload = {
      ID: seriesForm.value.ID,
      categoryId: seriesForm.value.categoryId,
      name: seriesForm.value.name,
      coverId: Number(seriesForm.value.coverId || 0),
      coverUrl: seriesForm.value.coverUrl,
      price: Number(seriesForm.value.price || 0),
      needVip: !!seriesForm.value.needVip
    }

    const res = seriesDialogMode.value === 'create'
      ? await createVideoSeries(payload)
      : await updateVideoSeries(payload)

    if (res.code !== 0) return
    ElMessage.success(seriesDialogMode.value === 'create' ? '剧集创建成功' : '剧集更新成功')
    seriesDialogVisible.value = false
    await Promise.all([loadSeriesList(), loadSeriesOptions(), loadEpisodeOptions()])
  }

  const removeSeries = (row) => {
    ElMessageBox.confirm(`确认删除剧集【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoSeries({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('剧集删除成功')
      await Promise.all([loadSeriesList(), loadSeriesOptions(), loadEpisodeList(), loadEpisodeOptions()])
    })
  }

  const openEpisodeDialog = (row) => {
    episodeDialogMode.value = row?.ID ? 'edit' : 'create'
    episodeForm.value = {
      ID: row?.ID || 0,
      seriesId: row?.seriesId || episodeQuery.value.seriesId || undefined,
      name: row?.name || '{}',
      videoUrl: row?.videoUrl || '',
      trialPercent: Number(row?.trialPercent || 8),
      sort: Number(row?.sort || 0)
    }
    episodeDialogVisible.value = true
  }

  const submitEpisode = async () => {
    if (!episodeForm.value.seriesId) {
      ElMessage.warning('请选择所属剧集')
      return
    }
    if (!episodeForm.value.name) {
      ElMessage.warning('单集名称不能为空')
      return
    }
    if (!episodeForm.value.videoUrl) {
      ElMessage.warning('视频地址不能为空')
      return
    }
    try {
      episodeForm.value.name = ensureJsonText(episodeForm.value.name)
    } catch (e) {
      ElMessage.error('单集名称必须是合法 JSON')
      return
    }

    const payload = {
      ID: episodeForm.value.ID,
      seriesId: episodeForm.value.seriesId,
      name: episodeForm.value.name,
      videoUrl: episodeForm.value.videoUrl,
      trialPercent: Number(episodeForm.value.trialPercent || 8),
      sort: Number(episodeForm.value.sort || 0)
    }

    const res = episodeDialogMode.value === 'create'
      ? await createVideoEpisode(payload)
      : await updateVideoEpisode(payload)

    if (res.code !== 0) return
    ElMessage.success(episodeDialogMode.value === 'create' ? '单集创建成功' : '单集更新成功')
    episodeDialogVisible.value = false
    await Promise.all([loadEpisodeList(), loadEpisodeOptions()])
  }

  const removeEpisode = (row) => {
    ElMessageBox.confirm(`确认删除单集【${formatI18nText(row.name)}】吗？`, '删除确认', {
      type: 'warning'
    }).then(async () => {
      const res = await deleteVideoEpisode({ ID: row.ID })
      if (res.code !== 0) return
      ElMessage.success('单集删除成功')
      await Promise.all([loadEpisodeList(), loadEpisodeOptions()])
    })
  }

  const handleParseSubtitle = async () => {
    if (!subtitleForm.value.episodeId) {
      ElMessage.warning('请选择目标单集')
      return
    }

    let subtitles = []
    try {
      subtitles = JSON.parse(subtitleForm.value.rawJsonData)
      if (!Array.isArray(subtitles)) {
        ElMessage.error('字幕 JSON 必须是数组结构')
        return
      }
    } catch (e) {
      ElMessage.error('字幕 JSON 格式有误')
      return
    }

    const res = await parseSubtitle({
      episodeId: subtitleForm.value.episodeId,
      subtitles
    })
    if (res.code !== 0) return
    ElMessage.success('字幕解析并入库成功')
  }

  const openEntitlementDialog = () => {
    entitlementForm.value = {
      userId: undefined,
      resourceType: 'video_series',
      resourceId: undefined,
      expireAt: '',
      remark: ''
    }
    entitlementDialogVisible.value = true
  }

  const submitEntitlement = async () => {
    if (!entitlementForm.value.userId) {
      ElMessage.warning('用户ID不能为空')
      return
    }
    if (!entitlementForm.value.resourceType) {
      ElMessage.warning('资源类型不能为空')
      return
    }
    if (!entitlementForm.value.resourceId) {
      ElMessage.warning('资源ID不能为空')
      return
    }

    const payload = {
      userId: Number(entitlementForm.value.userId),
      resourceType: entitlementForm.value.resourceType,
      resourceId: Number(entitlementForm.value.resourceId),
      remark: entitlementForm.value.remark || ''
    }
    if (entitlementForm.value.expireAt) {
      payload.expireAt = entitlementForm.value.expireAt
    }

    const res = await grantEntitlement(payload)
    if (res.code !== 0) return
    ElMessage.success('授权成功')
    entitlementDialogVisible.value = false
    loadEntitlementList()
  }

  const removeEntitlement = (row) => {
    ElMessageBox.confirm(`确认撤销用户 ${row.userId} 的授权记录 #${row.ID} 吗？`, '撤销确认', {
      type: 'warning'
    }).then(async () => {
      const res = await revokeEntitlement({
        userId: row.userId,
        resourceType: row.resourceType,
        resourceId: row.resourceId
      })
      if (res.code !== 0) return
      ElMessage.success('授权已撤销')
      loadEntitlementList()
    })
  }

  onMounted(async () => {
    await Promise.all([
      loadVideoCategoryList(),
      loadSeriesList(),
      loadEpisodeList(),
      loadEntitlementList(),
      loadVideoCategoryOptions(),
      loadSeriesOptions(),
      loadEpisodeOptions()
    ])
  })
</script>

<style scoped>
  .toolbar-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 16px;
    flex-wrap: wrap;
  }

  .toolbar-tip {
    color: #4b5563;
    font-size: 13px;
  }

  .subtitle-box {
    margin-top: 12px;
    padding: 12px 0;
  }

  .mt-2 {
    margin-top: 8px;
  }

  .w-full {
    width: 100%;
  }
</style>
