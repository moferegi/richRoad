<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAt">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期" :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期" :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
      </el-form-item>

        <el-form-item label="用户名" prop="username">
         <el-input v-model="searchInfo.username" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
         <el-input v-model="searchInfo.nickname" placeholder="搜索条件" />

        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        max-height="70vh"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />

        <el-table-column align="left" label="日期" width="180" sortable="custom" prop="created_at">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="用户名" prop="username" width="120" />
        <el-table-column align="left" label="密码" prop="password" width="120" />
        <el-table-column align="left" label="昵称" prop="nickname" width="120" />
        <el-table-column align="left" label="性别" prop="gender" width="120">
            <template #default="scope">
            {{ filterDict(scope.row.gender,genderOptions) }}
            </template>
        </el-table-column>
        <el-table-column align="left" label="手机号" prop="phone" width="120" />
        <el-table-column align="left" label="邮箱" prop="email" width="120" />
        <el-table-column align="left" label="积分" prop="point" width="80" />
        <el-table-column align="left" label="邀请码" prop="inviteCode" width="120" />
        <el-table-column align="left" label="邀请人ID" prop="invitedBy" width="100">
            <template #default="scope">
              {{ scope.row.invitedBy > 0 ? scope.row.invitedBy : '-' }}
            </template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="banned" width="100">
            <template #default="scope">
              <el-switch
                v-model="scope.row.banned"
                :active-value="false"
                :inactive-value="true"
                active-text="正常"
                inactive-text="封禁"
                inline-prompt
                @change="toggleBan(scope.row)"
              />
            </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="290">
            <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
                <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
                查看详情
            </el-button>
            <el-button type="primary" link @click="showSubordinates(scope.row)">下级</el-button>
            <el-button type="primary" link icon="edit" class="table-button" @click="updateClientUserFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
    <el-drawer size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #title>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'添加':'修改'}}</span>
                <div>
                  <el-button type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="用户名:"  prop="username" >
              <el-input v-model="formData.username" :clearable="true"  placeholder="请输入用户名" />
            </el-form-item>
            <el-form-item label="密码:"  prop="password" >
              <el-input v-model="formData.password" :clearable="true"  placeholder="请输入密码" />
            </el-form-item>
            <el-form-item label="昵称:"  prop="nickname" >
              <el-input v-model="formData.nickname" :clearable="true"  placeholder="请输入昵称" />
            </el-form-item>
            <el-form-item label="性别:"  prop="gender" >
              <el-select v-model="formData.gender" placeholder="请选择性别" style="width:100%" :clearable="true" >
                <el-option v-for="(item,key) in genderOptions" :key="key" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="手机号:"  prop="phone" >
              <el-input v-model="formData.phone" :clearable="true"  placeholder="请输入手机号" />
            </el-form-item>
            <el-form-item label="邮箱:"  prop="email" >
              <el-input v-model="formData.email" :clearable="true"  placeholder="请输入邮箱" />
            </el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer size="800" v-model="detailShow" :before-close="closeDetailShow" title="查看详情" destroy-on-close>
          <template #title>
             <div class="flex justify-between items-center">
               <span class="text-lg">查看详情</span>
             </div>
         </template>
        <el-descriptions :column="1" border>
                <el-descriptions-item label="用户名">
                        {{ formData.username }}
                </el-descriptions-item>
                <el-descriptions-item label="密码">
                        {{ formData.password }}
                </el-descriptions-item>
                <el-descriptions-item label="昵称">
                        {{ formData.nickname }}
                </el-descriptions-item>
                <el-descriptions-item label="性别">
                        {{ filterDict(formData.gender,genderOptions) }}
                </el-descriptions-item>
                <el-descriptions-item label="手机号">
                        {{ formData.phone }}
                </el-descriptions-item>
                <el-descriptions-item label="邮箱">
                        {{ formData.email }}
                </el-descriptions-item>
        </el-descriptions>
    </el-drawer>

    <!-- 下级用户列表弹窗 -->
    <el-drawer size="600" v-model="subordinateShow" title="下级用户列表" destroy-on-close>
      <el-table :data="subordinateList" stripe>
        <el-table-column prop="ID" label="ID" width="60" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column label="注册时间" min-width="160">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, prev, pager, next"
          :current-page="subPage"
          :page-size="subPageSize"
          :total="subTotal"
          @current-change="handleSubPageChange"
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createClientUser,
  deleteClientUser,
  deleteClientUserByIds,
  updateClientUser,
  findClientUser,
  getClientUserList,
  getSubordinates
} from '@/api/client/user'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({
    name: 'ClientUser'
})

// 自动化生成的字典（可能为空）以及字段
const genderOptions = ref([])
const formData = ref({
        username: '',
        password: '',
        nickname: '',
        gender: '',
        phone: '',
        email: '',
        })


// 验证规则
const rule = reactive({
               username : [{
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
               password : [{
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

// 排序
const sortChange = ({ prop, order }) => {
  if (order) {
    searchInfo.value.orderBy = prop
    searchInfo.value.orderDir = order === 'ascending' ? 'asc' : 'desc'
  } else {
    searchInfo.value.orderBy = ''
    searchInfo.value.orderDir = ''
  }
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
  const table = await getClientUserList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
const setOptions = async () =>{
    genderOptions.value = await getDictFunc('gender')
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
            deleteClientUserFunc(row)
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
      const res = await deleteClientUserByIds({ IDs })
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
const updateClientUserFunc = async(row) => {
    const res = await findClientUser({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.reclientUser
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteClientUserFunc = async (row) => {
    const res = await deleteClientUser({ ID: row.ID })
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
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findClientUser({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.reclientUser
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
          username: '',
          password: '',
          nickname: '',
          gender: '',
          phone: '',
          email: '',
          }
}


// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        username: '',
        password: '',
        nickname: '',
        gender: '',
        phone: '',
        email: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createClientUser(formData.value)
                  break
                case 'update':
                  res = await updateClientUser(formData.value)
                  break
                default:
                  res = await createClientUser(formData.value)
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

// ============== 封禁用户 ===============
const toggleBan = async (row) => {
  const res = await updateClientUser({ ID: row.ID, banned: row.banned })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: row.banned ? '已封禁' : '已解封' })
  } else {
    row.banned = !row.banned
  }
}

// ============== 下级用户 ===============
const subordinateShow = ref(false)
const subordinateList = ref([])
const subPage = ref(1)
const subPageSize = ref(10)
const subTotal = ref(0)
const currentSubUserID = ref(0)

const showSubordinates = async (row) => {
  currentSubUserID.value = row.ID
  subPage.value = 1
  await loadSubordinates()
  subordinateShow.value = true
}

const loadSubordinates = async () => {
  const res = await getSubordinates({
    userID: currentSubUserID.value,
    page: subPage.value,
    pageSize: subPageSize.value
  })
  if (res.code === 0) {
    subordinateList.value = res.data.list || []
    subTotal.value = res.data.total
  }
}

const handleSubPageChange = (val) => {
  subPage.value = val
  loadSubordinates()
}

</script>

<style>

</style>
