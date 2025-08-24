
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAtRange">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>

      <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="w-[380px]"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
       </el-form-item>
      
            <el-form-item label="用户ID" prop="userId">
  <el-select 
    v-model="searchInfo.userId" 
    clearable 
    filterable 
    remote 
    reserve-keyword 
    placeholder="请选择用户" 
    :remote-method="remoteUserSearch"
    :loading="userLoading"
    @clear="()=>{searchInfo.userId=undefined}"
  >
    <el-option 
      v-for="item in userOptions" 
      :key="item.ID" 
      :label="`${item.nickname || item.username} (ID: ${item.ID})`" 
      :value="item.ID" 
    />
  </el-select>
</el-form-item>
            
            <el-form-item label="增/减" prop="changeType">
  <el-select v-model="searchInfo.changeType" clearable filterable placeholder="请选择" @clear="()=>{searchInfo.changeType=undefined}">
    <el-option label="增加" value="increase" />
    <el-option label="减少" value="decrease" />
  </el-select>
</el-form-item>
            
            <el-form-item label="积分变化" prop="pointChange">
  <el-input v-model.number="searchInfo.pointChange" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="操作类型" prop="operationType">
  <el-select v-model="searchInfo.operationType" clearable filterable placeholder="请选择" @clear="()=>{searchInfo.operationType=undefined}">
    <el-option v-for="(item,key) in point_operation_typeOptions" :key="key" :label="item.label" :value="item.value" />
  </el-select>
</el-form-item>
            
            <el-form-item label="变化原因" prop="reason">
  <el-input v-model="searchInfo.reason" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="当前积分" prop="currentPoints">
  <el-input v-model.number="searchInfo.currentPoints" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="关联订单ID" prop="relatedOrderId">
  <el-input v-model.number="searchInfo.relatedOrderId" placeholder="搜索条件" />
</el-form-item>
            
            <el-form-item label="备注" prop="remark">
  <el-input v-model="searchInfo.remark" placeholder="搜索条件" />
</el-form-item>
            

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button  type="primary" icon="plus" @click="openDialog()">新增</el-button>
            <el-button  icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
            <ExportTemplate  template-id="client_PointRecord" />
            <ExportExcel  template-id="client_PointRecord" filterDeleted/>
            <ImportExcel  template-id="client_PointRecord" @on-success="getTableData" />
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />
        
        <el-table-column sortable align="left" label="日期" prop="CreatedAt"width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        
            <el-table-column align="left" label="用户" prop="userId" width="150">
                <template #default="scope">
                    <span v-if="getUserDisplayName(scope.row.userId)">{{ getUserDisplayName(scope.row.userId) }}</span>
                    <span v-else>ID: {{ scope.row.userId }}</span>
                </template>
            </el-table-column>

            <el-table-column align="left" label="增/减" prop="changeType" width="80">
                <template #default="scope">
                    <el-tag :type="scope.row.pointChange > 0 ? 'success' : 'danger'">
                        {{ scope.row.pointChange > 0 ? '增加' : '减少' }}
                    </el-tag>
                </template>
            </el-table-column>

            <el-table-column sortable align="left" label="积分变化" prop="pointChange" width="120">
                <template #default="scope">
                    <span :class="scope.row.pointChange > 0 ? 'text-green-600' : 'text-red-600'">
                        {{ scope.row.pointChange > 0 ? '+' : '' }}{{ scope.row.pointChange }}
                    </span>
                </template>
            </el-table-column>

            <el-table-column align="left" label="操作类型" prop="operationType" width="120">
    <template #default="scope">
    {{ filterDict(scope.row.operationType,point_operation_typeOptions) }}
    </template>
</el-table-column>
            <el-table-column align="left" label="变化原因" prop="reason" width="120" />

            <el-table-column align="left" label="当前积分" prop="currentPoints" width="120" />

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updatePointRecordFunc(scope.row)">编辑</el-button>
            <el-button   type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'新增':'编辑'}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="用户ID:" prop="userId">
    <el-select 
      v-model="formData.userId" 
      clearable 
      filterable 
      remote 
      reserve-keyword 
      placeholder="请选择用户" 
      :remote-method="remoteUserSearch"
      :loading="userLoading"
      style="width:100%"
      @clear="()=>{formData.userId=undefined}"
    >
      <el-option 
        v-for="item in userOptions" 
        :key="item.ID" 
        :label="`${item.nickname || item.username} (ID: ${item.ID})`" 
        :value="item.ID" 
      />
    </el-select>
</el-form-item>
            <el-form-item label="增/减:" prop="changeType">
    <el-select v-model="formData.changeType" placeholder="请选择增减类型" style="width:100%" filterable :clearable="true">
        <el-option label="增加" value="increase" />
        <el-option label="减少" value="decrease" />
    </el-select>
</el-form-item>
            <el-form-item label="积分变化:" prop="pointChange">
    <el-input v-model.number="formData.pointChange" :clearable="true" placeholder="请输入积分变化" />
</el-form-item>
            <el-form-item label="操作类型:" prop="operationType">
    <el-select v-model="formData.operationType" placeholder="请选择操作类型" style="width:100%" filterable :clearable="true">
        <el-option v-for="(item,key) in point_operation_typeOptions" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
            <el-form-item label="变化原因:" prop="reason">
    <el-input v-model="formData.reason" :clearable="true" placeholder="请输入变化原因" />
</el-form-item>
            <el-form-item label="关联订单ID:" prop="relatedOrderId">
    <el-input v-model.number="formData.relatedOrderId" :clearable="true" placeholder="请输入关联订单ID" />
</el-form-item>
            <el-form-item label="备注:" prop="remark">
    <el-input v-model="formData.remark" :clearable="true" placeholder="请输入备注" />
</el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="用户">
    <span v-if="getUserDisplayName(detailForm.userId)">{{ getUserDisplayName(detailForm.userId) }}</span>
    <span v-else>ID: {{ detailForm.userId }}</span>
</el-descriptions-item>
                    <el-descriptions-item label="增/减">
    <el-tag :type="detailForm.pointChange > 0 ? 'success' : 'danger'">
        {{ detailForm.pointChange > 0 ? '增加' : '减少' }}
    </el-tag>
</el-descriptions-item>
                    <el-descriptions-item label="积分变化">
    <span :class="detailForm.pointChange > 0 ? 'text-green-600' : 'text-red-600'">
        {{ detailForm.pointChange > 0 ? '+' : '' }}{{ detailForm.pointChange }}
    </span>
</el-descriptions-item>
                    <el-descriptions-item label="操作类型">
    {{ detailForm.operationType }}
</el-descriptions-item>
                    <el-descriptions-item label="变化原因">
    {{ detailForm.reason }}
</el-descriptions-item>
                    <el-descriptions-item label="当前积分">
    {{ detailForm.currentPoints }}
</el-descriptions-item>
                    <el-descriptions-item label="关联订单ID">
    {{ detailForm.relatedOrderId }}
</el-descriptions-item>
                    <el-descriptions-item label="备注">
    {{ detailForm.remark }}
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  createPointRecord,
  deletePointRecord,
  deletePointRecordByIds,
  updatePointRecord,
  findPointRecord,
  getPointRecordList
} from '@/api/client/pointRecord'
import { getClientUserList } from '@/api/client/user'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict ,filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"

// 导出组件
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
// 导入组件
import ImportExcel from '@/components/exportExcel/importExcel.vue'
// 导出模板组件
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'


defineOptions({
    name: 'PointRecord'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const point_operation_typeOptions = ref([])
const userOptions = ref([])
const formData = ref({
            userId: undefined,
            changeType: '',
            pointChange: undefined,
            operationType: '',
            reason: '',
            relatedOrderId: undefined,
            remark: '',
        })



// 验证规则
const rule = reactive({
               userId : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
              ],
               changeType : [{
                   required: true,
                   message: '请选择增减类型',
                   trigger: ['change','blur'],
               },
              ],
               pointChange : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
              ],
               operationType : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
               reason : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
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
    CreatedAt:"CreatedAt",
    ID:"ID",
            pointChange: 'point_change',
  }

  let sort = sortMap[prop]
  if(!sort){
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
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
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
  const table = await getPointRecordList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
    
    // 获取表格中所有用户的信息
    const userIds = table.data.list.map(item => item.userId).filter(Boolean)
    if (userIds.length > 0) {
      await loadUsersInfo(userIds)
    }
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 用户搜索加载状态
const userLoading = ref(false)

// 远程搜索用户（本地过滤）
const remoteUserSearch = async (query) => {
  if (query) {
    // 从已加载的用户数据中进行本地过滤
    const allUsers = Array.from(userCache.value.values())
    userOptions.value = allUsers.filter(user => 
      (user.nickname && user.nickname.toLowerCase().includes(query.toLowerCase())) ||
      (user.username && user.username.toLowerCase().includes(query.toLowerCase()))
    )
  } else {
    // 显示所有用户
    userOptions.value = Array.from(userCache.value.values())
  }
}

// 初始化加载全部用户数据
const initUserOptions = async () => {
  try {
    const res = await getClientUserList({ page: 1, pageSize: 1000 })
    if (res.code === 0) {
      userOptions.value = res.data.list || []
      // 同时填充用户缓存
      res.data.list.forEach(user => {
        userCache.value.set(user.ID, user)
      })
    }
  } catch (error) {
    console.error('初始化用户列表失败:', error)
  }
}

// 用户信息缓存
const userCache = ref(new Map())

// 获取用户显示名称
const getUserDisplayName = (userId) => {
  if (!userId) return ''
  const user = userCache.value.get(userId)
  if (user) {
    return `${user.nickname || user.username} (ID: ${user.ID})`
  }
  return ''
}

// 批量获取用户信息
const loadUsersInfo = async (userIds) => {
  const uniqueIds = [...new Set(userIds)].filter(id => id && !userCache.value.has(id))
  if (uniqueIds.length === 0) return
  
  try {
    // 获取更多用户数据来填充缓存
    const res = await getClientUserList({ page: 1, pageSize: 100 })
    if (res.code === 0 && res.data.list) {
      res.data.list.forEach(user => {
        userCache.value.set(user.ID, user)
      })
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
    point_operation_typeOptions.value = await getDictFunc('point_operation_type')
}

// 获取需要的字典 可能为空 按需保留
setOptions()
// 初始化用户选项
initUserOptions()


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
            deletePointRecordFunc(row)
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
      const res = await deletePointRecordByIds({ IDs })
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
const updatePointRecordFunc = async(row) => {
    const res = await findPointRecord({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data
        dialogFormVisible.value = true
    }
}


// 删除行
const deletePointRecordFunc = async (row) => {
    const res = await deletePointRecord({ ID: row.ID })
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
        userId: undefined,
        pointChange: undefined,
        operationType: '',
        reason: '',
        relatedOrderId: undefined,
        remark: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     btnLoading.value = true
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return btnLoading.value = false
              let res
              switch (type.value) {
                case 'create':
                  res = await createPointRecord(formData.value)
                  break
                case 'update':
                  res = await updatePointRecord(formData.value)
                  break
                default:
                  res = await createPointRecord(formData.value)
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

const detailForm = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findPointRecord({ ID: row.ID })
  if (res.code === 0) {
    detailForm.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailForm.value = {}
}


</script>

<style>

</style>
