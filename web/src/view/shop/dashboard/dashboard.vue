<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item>
          <el-button type="primary" icon="refresh" @click="loadDashboard">刷新数据</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-row :gutter="20" class="mb-4">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="dashboard-card">
          <template #header><span>总销售额</span></template>
          <div class="dashboard-value text-green-500">¥ {{ (overview.totalSales / 100).toFixed(2) }}</div>
          <div class="dashboard-sub">今日: ¥ {{ (overview.todaySales / 100).toFixed(2) }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="dashboard-card">
          <template #header><span>订单统计</span></template>
          <div class="dashboard-value">{{ overview.orderTotal }}</div>
          <div class="dashboard-sub">已付款: {{ overview.orderPaid }} | 已取消: {{ overview.orderCancelled }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="dashboard-card">
          <template #header><span>用户统计</span></template>
          <div class="dashboard-value">{{ overview.userTotal }}</div>
          <div class="dashboard-sub">今日新增: {{ overview.userToday }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card shadow="hover" class="dashboard-card">
          <template #header><span>积分统计</span></template>
          <div class="dashboard-value">{{ overview.pointsIssued }}</div>
          <div class="dashboard-sub">已使用: {{ overview.pointsUsed }} | 货币额: ¥{{ (overview.pointsCurrencyValue / 100).toFixed(2) }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :xs="24" :sm="12">
        <el-card shadow="hover" class="dashboard-card">
          <template #header><span>优惠券统计</span></template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="发放总量">{{ overview.couponIssuedCount }}</el-descriptions-item>
            <el-descriptions-item label="已领取量">{{ overview.couponClaimedCount }}</el-descriptions-item>
            <el-descriptions-item label="发放总金额">¥ {{ (overview.couponIssuedAmount / 100).toFixed(2) }}</el-descriptions-item>
            <el-descriptions-item label="已使用总金额">¥ {{ (overview.couponUsedAmount / 100).toFixed(2) }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12">
        <el-card shadow="hover" class="dashboard-card">
          <template #header><span>积分详情</span></template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="已发放积分">{{ overview.pointsIssued }}</el-descriptions-item>
            <el-descriptions-item label="已使用积分">{{ overview.pointsUsed }}</el-descriptions-item>
            <el-descriptions-item label="积分对应货币额">¥ {{ (overview.pointsCurrencyValue / 100).toFixed(2) }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { getDashboardOverview } from '@/api/shop/dashboard'
import { ref, onMounted } from 'vue'

defineOptions({ name: 'Dashboard' })

const overview = ref({
  totalSales: 0,
  todaySales: 0,
  orderTotal: 0,
  orderPaid: 0,
  orderCancelled: 0,
  couponIssuedCount: 0,
  couponClaimedCount: 0,
  couponIssuedAmount: 0,
  couponUsedAmount: 0,
  pointsIssued: 0,
  pointsUsed: 0,
  pointsCurrencyValue: 0,
  userTotal: 0,
  userToday: 0,
})

const loadDashboard = async () => {
  const res = await getDashboardOverview()
  if (res.code === 0) {
    overview.value = res.data
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.dashboard-card {
  margin-bottom: 20px;
}
.dashboard-value {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 8px;
}
.dashboard-sub {
  font-size: 13px;
  color: #909399;
}
</style>
