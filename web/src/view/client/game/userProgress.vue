<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="用户ID">
          <el-input v-model="searchInfo.userID" placeholder="输入用户ID" style="width: 180px" clearable />
        </el-form-item>
        <el-form-item label="游戏分类">
          <el-select v-model="searchInfo.gameID" placeholder="选择游戏" style="width: 180px" clearable @change="onGameChange">
            <el-option v-for="g in gameList" :key="g.ID" :label="getGameLabel(g)" :value="g.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="难度分类">
          <el-select v-model="searchInfo.categoryID" placeholder="选择难度" style="width: 180px" clearable @change="onCategoryChange">
            <el-option v-for="c in difficultyList" :key="c.ID" :label="getDifficultyLabel(c)" :value="c.ID" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-if="userInfo" class="gva-search-box" style="margin-top: 10px">
      <el-descriptions :column="3" border size="small">
        <el-descriptions-item label="用户名">{{ userInfo.username }}</el-descriptions-item>
        <el-descriptions-item label="昵称">{{ userInfo.nickname || '-' }}</el-descriptions-item>
        <el-descriptions-item label="已通关数">{{ clearedCount }} / {{ levelList.length }}</el-descriptions-item>
      </el-descriptions>
    </div>

    <div class="gva-table-box">
      <el-table :data="levelList" row-key="ID" border v-loading="loading">
        <el-table-column align="left" label="关卡号" prop="levelNumber" width="100" />
        <el-table-column align="left" label="数字" prop="numbers" min-width="200" />
        <el-table-column align="left" label="目标结果" prop="targetResult" width="120">
          <template #default="scope">
            <el-tag type="warning">{{ scope.row.targetResult }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" width="140">
          <template #default="scope">
            <el-switch 
              :model-value="isCleared(scope.row.ID)" 
              :active-text="isCleared(scope.row.ID) ? '已通关' : '未通关'"
              @change="(val) => onToggleProgress(scope.row.ID, val)"
              :loading="scope.row._toggling"
            />
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { getCategoryList, getDifficultyCategoryList, getLevelList, getUserProgress, setUserProgress } from '@/api/client/game'
import { getUserList } from '@/api/user'

defineOptions({ name: 'GameUserProgress' })

const loading = ref(false)
const gameList = ref([])
const difficultyList = ref([])
const levelList = ref([])
const progressMap = ref({})
const userInfo = ref(null)
const clearedCount = ref(0)

const searchInfo = reactive({
  userID: '',
  gameID: null,
  categoryID: null
})

const getGameLabel = (g) => {
  try {
    const name = typeof g.name === 'string' ? JSON.parse(g.name) : g.name
    return name?.zh || name?.en || g.gameKey || ''
  } catch {
    return g.name || g.gameKey || ''
  }
}

const getDifficultyLabel = (c) => {
  try {
    const name = typeof c.name === 'string' ? JSON.parse(c.name) : c.name
    return name?.zh || name?.en || ''
  } catch {
    return c.name || ''
  }
}

const isCleared = (levelID) => {
  return progressMap.value[levelID]
}

const loadGameList = async () => {
  const res = await getCategoryList({ page: 1, pageSize: 100 })
  if (res.code === 0) {
    gameList.value = res.data?.list || []
  }
}

const onGameChange = async (val) => {
  searchInfo.categoryID = null
  difficultyList.value = []
  levelList.value = []
  progressMap.value = {}
  userInfo.value = null
  clearedCount.value = 0
  if (!val) return
  const res = await getDifficultyCategoryList({ page: 1, pageSize: 100, gameID: val })
  if (res.code === 0) {
    difficultyList.value = res.data?.list || []
  }
}

const onCategoryChange = async () => {
  levelList.value = []
  progressMap.value = {}
  clearedCount.value = 0
}

const loadUserInfo = async () => {
  if (!searchInfo.userID) return
  const res = await getUserList({ page: 1, pageSize: 1, ID: Number(searchInfo.userID) })
  if (res.code === 0 && res.data?.list?.length > 0) {
    userInfo.value = res.data.list[0]
  } else {
    userInfo.value = null
  }
}

const onSearch = async () => {
  if (!searchInfo.userID || !searchInfo.categoryID) {
    return
  }
  loading.value = true
  try {
    await loadUserInfo()
    const [levelRes, progressRes] = await Promise.all([
      getLevelList({ page: 1, pageSize: 200, categoryID: searchInfo.categoryID }),
      getUserProgress({ userID: searchInfo.userID, categoryID: searchInfo.categoryID })
    ])
    if (levelRes.code === 0) {
      levelList.value = (levelRes.data?.list || []).sort((a, b) => a.levelNumber - b.levelNumber)
    }
    if (progressRes.code === 0) {
      const map = {}
      let count = 0
      ;(progressRes.data || []).forEach(p => {
        if (p.status === 1) {
          map[p.levelID] = true
          count++
        }
      })
      progressMap.value = map
      clearedCount.value = count
    }
  } finally {
    loading.value = false
  }
}

const onToggleProgress = async (levelID, val) => {
  const level = levelList.value.find(l => l.ID === levelID)
  if (level) level._toggling = true
  try {
    await setUserProgress({
      userID: Number(searchInfo.userID),
      levelID: levelID,
      status: val ? 1 : 0
    })
    if (val) {
      progressMap.value[levelID] = true
      clearedCount.value++
    } else {
      delete progressMap.value[levelID]
      clearedCount.value--
    }
  } finally {
    if (level) level._toggling = false
  }
}

loadGameList()
</script>