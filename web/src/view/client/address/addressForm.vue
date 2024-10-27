<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="收件人名称:" prop="name">
          <el-input v-model="formData.name" :clearable="true"  placeholder="请输入收件人名称" />
       </el-form-item>
        <el-form-item label="收件人电话:" prop="phone">
          <el-input v-model="formData.phone" :clearable="true"  placeholder="请输入收件人电话" />
       </el-form-item>
        <el-form-item label="收件省份:" prop="province">
        <el-select v-model="formData.province" placeholder="请选择收件省份" style="width:100%" :clearable="true" >
          <el-option v-for="(item,key) in dataSource.province" :key="key" :label="item.label" :value="item.value" />
        </el-select>
       </el-form-item>
        <el-form-item label="收件城市:" prop="city">
        <el-select v-model="formData.city" placeholder="请选择收件城市" style="width:100%" :clearable="true" >
          <el-option v-for="(item,key) in dataSource.city" :key="key" :label="item.label" :value="item.value" />
        </el-select>
       </el-form-item>
        <el-form-item label="收件区域:" prop="area">
        <el-select v-model="formData.area" placeholder="请选择收件区域" style="width:100%" :clearable="true" >
          <el-option v-for="(item,key) in dataSource.area" :key="key" :label="item.label" :value="item.value" />
        </el-select>
       </el-form-item>
        <el-form-item label="收件详细地址:" prop="street">
          <el-input v-model="formData.street" :clearable="true"  placeholder="请输入收件详细地址" />
       </el-form-item>
        <el-form-item label="用户id:" prop="userID">
          <el-input v-model="formData.userID" :clearable="true"  placeholder="请输入用户id" />
       </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
    getAddressDataSource,
  createAddress,
  updateAddress,
  findAddress
} from '@/api/client/address'

defineOptions({
    name: 'AddressForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            name: '',
            phone: '',
            province: 0,
            city: 0,
            area: 0,
            street: '',
            userID: '',
        })
// 验证规则
const rule = reactive({
               name : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               phone : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               province : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               city : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               area : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               street : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               userID : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getAddressDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findAddress({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.readdress
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createAddress(formData.value)
               break
             case 'update':
               res = await updateAddress(formData.value)
               break
             default:
               res = await createAddress(formData.value)
               break
           }
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
