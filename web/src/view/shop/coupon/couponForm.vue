
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="名称:" prop="name">
    <el-input v-model="formData.name" :clearable="false" placeholder="请输入名称" />
</el-form-item>
        <el-form-item label="描述:" prop="description">
    <el-input v-model="formData.description" :clearable="false" placeholder="请输入描述" />
</el-form-item>
        <el-form-item label="类型:" prop="type">
    <el-select v-model="formData.type" placeholder="请选择类型" style="width:100%" filterable :clearable="false">
       <el-option v-for="item in ['product','discount','no_threshold']" :key="item" :label="item" :value="item" />
    </el-select>
</el-form-item>
        <el-form-item label="折扣:" prop="discount">
    <el-input-number v-model="formData.discount" style="width:100%" :precision="2" :clearable="false" />
</el-form-item>
        <el-form-item label="最低消费:" prop="minSpend">
    <el-input-number v-model="formData.minSpend" style="width:100%" :precision="2" :clearable="false" />
</el-form-item>
        <el-form-item label="商品ID:" prop="productID">
    <el-select v-model="formData.productID" placeholder="请选择商品ID" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.productID" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="数量:" prop="quantity">
    <el-input v-model.number="formData.quantity" :clearable="false" placeholder="请输入数量" />
</el-form-item>
        <el-form-item label="已领取:" prop="claimed">
    <el-input v-model.number="formData.claimed" :clearable="false" placeholder="请输入已领取" />
</el-form-item>
        <el-form-item label="开始时间:" prop="startTime">
    <el-date-picker v-model="formData.startTime" type="date" style="width:100%" placeholder="选择日期" :clearable="false" />
</el-form-item>
        <el-form-item label="结束时间:" prop="endTime">
    <el-date-picker v-model="formData.endTime" type="date" style="width:100%" placeholder="选择日期" :clearable="false" />
</el-form-item>
        <el-form-item label="启用:" prop="status">
    <el-switch v-model="formData.status" active-color="#13ce66" inactive-color="#ff4949" active-text="是" inactive-text="否" clearable ></el-switch>
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
    getCouponDataSource,
  createCoupon,
  updateCoupon,
  findCoupon
} from '@/api/shop/coupon'

defineOptions({
    name: 'CouponForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const formData = ref({
            name: '',
            description: '',
            type: null,
            discount: 0,
            minSpend: 0,
            productID: undefined,
            quantity: 0,
            claimed: 0,
            startTime: new Date(),
            endTime: new Date(),
            status: false,
        })
// 验证规则
const rule = reactive({
               name : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               description : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               type : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               discount : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               minSpend : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               quantity : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               claimed : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               startTime : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               endTime : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               status : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getCouponDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findCoupon({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createCoupon(formData.value)
               break
             case 'update':
               res = await updateCoupon(formData.value)
               break
             default:
               res = await createCoupon(formData.value)
               break
           }
           btnLoading.value = false
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
