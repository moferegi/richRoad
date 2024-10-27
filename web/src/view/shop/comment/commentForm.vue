<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="用户ID:" prop="userID">
          <el-input v-model.number="formData.userID" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="订单ID:" prop="orderID">
          <el-input v-model.number="formData.orderID" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="商品ID:" prop="goodID">
          <el-input v-model.number="formData.goodID" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="SKUID:" prop="SKUID">
          <el-input v-model.number="formData.SKUID" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="用户评分:" prop="rating">
          <el-input v-model.number="formData.rating" :clearable="true" placeholder="请输入" />
       </el-form-item>
        <el-form-item label="评论内容:" prop="content">
          <el-input v-model="formData.content" :clearable="true"  placeholder="请输入评论内容" />
       </el-form-item>
        <el-form-item label="商家回复:" prop="shopReply">
          <el-input v-model="formData.shopReply" :clearable="true"  placeholder="请输入商家回复" />
       </el-form-item>
        <el-form-item label="商家回复创建时间:" prop="shopReplyAt">
          <el-date-picker v-model="formData.shopReplyAt" type="date" placeholder="选择日期" :clearable="true"></el-date-picker>
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
  createComment,
  updateComment,
  findComment
} from '@/api/shop/comment'

defineOptions({
    name: 'CommentForm'
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
            userID: undefined,
            orderID: undefined,
            goodID: undefined,
            SKUID: undefined,
            rating: undefined,
            content: '',
            shopReply: '',
            shopReplyAt: new Date(),
        })
// 验证规则
const rule = reactive({
               orderID : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               rating : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
               content : [{
                   required: true,
                   message: '',
                   trigger: ['input','blur'],
               }],
})

const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findComment({ ID: route.query.id })
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
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return
            let res
           switch (type.value) {
             case 'create':
               res = await createComment(formData.value)
               break
             case 'update':
               res = await updateComment(formData.value)
               break
             default:
               res = await createComment(formData.value)
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
