<template>
  <div class="dashboard-container">
    <!-- 顶部操作栏 -->
    <div class="dashboard-header">
      <div class="dashboard-title">
        <span class="title-icon">📊</span>
        <span>数据看板</span>
      </div>
      <el-button type="primary" :icon="Refresh" @click="loadDashboard" :loading="loading">
        刷新数据
      </el-button>
    </div>

    <!-- 核心指标卡片 -->
    <el-row :gutter="16" class="metric-row">
      <el-col :xs="12" :sm="8" :md="4">
        <div class="metric-card metric-sales">
          <div class="metric-icon">💰</div>
          <div class="metric-body">
            <div class="metric-label">总销售额</div>
            <div class="metric-value">¥{{ formatMoney(overview.totalSales) }}</div>
            <div class="metric-sub">今日 ¥{{ formatMoney(overview.todaySales) }}</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :md="4">
        <div class="metric-card metric-orders">
          <div class="metric-icon">📦</div>
          <div class="metric-body">
            <div class="metric-label">总订单数</div>
            <div class="metric-value">{{ overview.orderTotal }}</div>
            <div class="metric-sub">今日 {{ overview.orderToday }}</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :md="4">
        <div class="metric-card metric-users">
          <div class="metric-icon">👥</div>
          <div class="metric-body">
            <div class="metric-label">总用户数</div>
            <div class="metric-value">{{ overview.userTotal }}</div>
            <div class="metric-sub">今日新增 +{{ overview.userToday }}</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :md="4">
        <div class="metric-card metric-visitors">
          <div class="metric-icon">👁️</div>
          <div class="metric-body">
            <div class="metric-label">今日访客</div>
            <div class="metric-value">{{ overview.visitorUV }}</div>
            <div class="metric-sub">PV {{ overview.visitorPV }} · 新 {{ overview.visitorNew }}</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :md="4">
        <div class="metric-card metric-products">
          <div class="metric-icon">🛍️</div>
          <div class="metric-body">
            <div class="metric-label">商品数量</div>
            <div class="metric-value">{{ overview.productActive }}</div>
            <div class="metric-sub">总计 {{ overview.productTotal }} 件</div>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="8" :md="4">
        <div class="metric-card metric-signin">
          <div class="metric-icon">✅</div>
          <div class="metric-body">
            <div class="metric-label">今日签到</div>
            <div class="metric-value">{{ overview.signInToday }}</div>
            <div class="metric-sub">
              签到率 {{ overview.userTotal ? ((overview.signInToday / overview.userTotal) * 100).toFixed(1) : 0 }}%
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 趋势图表 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :xs="24" :lg="16">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="chart-card-header">
              <span class="chart-card-title">📈 销售趋势（近7天）</span>
            </div>
          </template>
          <div class="chart-wrap">
            <VCharts :option="salesChartOption" autoresize style="width: 100%; height: 320px;" />
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="chart-card-header">
              <span class="chart-card-title">🔄 订单状态分布</span>
            </div>
          </template>
          <div class="chart-wrap">
            <VCharts :option="orderPieOption" autoresize style="width: 100%; height: 320px;" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 订单 + 新用户趋势 -->
    <el-row :gutter="16" class="chart-row">
      <el-col :xs="24" :md="12">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="chart-card-header">
              <span class="chart-card-title">📦 订单趋势（近7天）</span>
            </div>
          </template>
          <div class="chart-wrap">
            <VCharts :option="orderChartOption" autoresize style="width: 100%; height: 280px;" />
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="chart-card-header">
              <span class="chart-card-title">👥 新用户趋势（近7天）</span>
            </div>
          </template>
          <div class="chart-wrap">
            <VCharts :option="userChartOption" autoresize style="width: 100%; height: 280px;" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 优惠券 & 积分详情 -->
    <el-row :gutter="16" class="detail-row">
      <el-col :xs="24" :md="12">
        <el-card shadow="never" class="detail-card">
          <template #header>
            <div class="chart-card-header">
              <span class="chart-card-title">🎫 优惠券统计</span>
            </div>
          </template>
          <div class="detail-grid">
            <div class="detail-item">
              <div class="detail-item-label">发放总量</div>
              <div class="detail-item-value">{{ overview.couponIssuedCount }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">已领取量</div>
              <div class="detail-item-value text-blue">{{ overview.couponClaimedCount }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">领取率</div>
              <div class="detail-item-value text-green">
                {{ overview.couponIssuedCount ? ((overview.couponClaimedCount / overview.couponIssuedCount) * 100).toFixed(1) : 0 }}%
              </div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">发放总金额</div>
              <div class="detail-item-value">¥{{ formatMoney(overview.couponIssuedAmount) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">已使用金额</div>
              <div class="detail-item-value text-red">¥{{ formatMoney(overview.couponUsedAmount) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">使用率</div>
              <div class="detail-item-value text-orange">
                {{ overview.couponIssuedAmount ? ((overview.couponUsedAmount / overview.couponIssuedAmount) * 100).toFixed(1) : 0 }}%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card shadow="never" class="detail-card">
          <template #header>
            <div class="chart-card-header">
              <span class="chart-card-title">💎 积分统计</span>
            </div>
          </template>
          <div class="detail-grid">
            <div class="detail-item">
              <div class="detail-item-label">已发放积分</div>
              <div class="detail-item-value text-green">{{ overview.pointsIssued }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">已使用积分</div>
              <div class="detail-item-value text-red">{{ overview.pointsUsed }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">使用率</div>
              <div class="detail-item-value text-blue">
                {{ overview.pointsIssued ? ((overview.pointsUsed / overview.pointsIssued) * 100).toFixed(1) : 0 }}%
              </div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">积分货币额</div>
              <div class="detail-item-value text-orange">¥{{ formatMoney(overview.pointsCurrencyValue) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">剩余积分</div>
              <div class="detail-item-value">{{ overview.pointsIssued - overview.pointsUsed }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-item-label">积分净值</div>
              <div class="detail-item-value text-green">
                ¥{{ overview.pointsIssued ? formatMoney(Math.round(overview.pointsCurrencyValue * ((overview.pointsIssued - overview.pointsUsed) / (overview.pointsUsed || 1)))) : '0.00' }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { getDashboardOverview } from '@/api/shop/dashboard'
import { ref, computed, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import VCharts from 'vue-echarts'

defineOptions({ name: 'ShopDashboard' })

const loading = ref(false)

const overview = ref({
  totalSales: 0, todaySales: 0,
  orderTotal: 0, orderPending: 0, orderPaid: 0, orderShipped: 0,
  orderReceived: 0, orderCancelled: 0, orderRefunding: 0, orderRefunded: 0, orderToday: 0,
  couponIssuedCount: 0, couponClaimedCount: 0, couponIssuedAmount: 0, couponUsedAmount: 0,
  pointsIssued: 0, pointsUsed: 0, pointsCurrencyValue: 0,
  userTotal: 0, userToday: 0,
  visitorPV: 0, visitorUV: 0, visitorNew: 0, signInToday: 0,
  productTotal: 0, productActive: 0,
  salesTrend: [], orderTrend: [], userTrend: [],
})

const formatMoney = (val) => ((val || 0) / 100).toFixed(2)

const loadDashboard = async () => {
  loading.value = true
  try {
    const res = await getDashboardOverview()
    if (res.code === 0) overview.value = res.data
  } finally { loading.value = false }
}

onMounted(() => { loadDashboard() })

// ============= 图表配置 =============

const salesChartOption = computed(() => {
  const trend = overview.value.salesTrend || []
  return {
    tooltip: { trigger: 'axis', formatter: (p) => `${p[0].axisValue}<br/>销售额: ¥${(p[0].value / 100).toFixed(2)}` },
    grid: { left: 60, right: 20, top: 20, bottom: 30 },
    xAxis: { type: 'category', data: trend.map(i => i.date.slice(5)), axisLabel: { color: '#999' }, axisLine: { lineStyle: { color: '#e6e6e6' } } },
    yAxis: { type: 'value', axisLabel: { formatter: (v) => '¥' + (v / 100).toFixed(0), color: '#999' }, splitLine: { lineStyle: { color: '#f0f0f0' } } },
    series: [{
      type: 'line', data: trend.map(i => i.value), smooth: true,
      areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
        colorStops: [{ offset: 0, color: 'rgba(64, 158, 255, 0.3)' }, { offset: 1, color: 'rgba(64, 158, 255, 0.02)' }] } },
      lineStyle: { color: '#409EFF', width: 3 },
      itemStyle: { color: '#409EFF' },
    }]
  }
})

const orderChartOption = computed(() => {
  const trend = overview.value.orderTrend || []
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 20, bottom: 30 },
    xAxis: { type: 'category', data: trend.map(i => i.date.slice(5)), axisLabel: { color: '#999' }, axisLine: { lineStyle: { color: '#e6e6e6' } } },
    yAxis: { type: 'value', minInterval: 1, axisLabel: { color: '#999' }, splitLine: { lineStyle: { color: '#f0f0f0' } } },
    series: [{
      type: 'bar', data: trend.map(i => i.value), barWidth: '50%',
      itemStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
        colorStops: [{ offset: 0, color: '#67C23A' }, { offset: 1, color: 'rgba(103, 194, 58, 0.3)' }] },
        borderRadius: [4, 4, 0, 0] },
    }]
  }
})

const userChartOption = computed(() => {
  const trend = overview.value.userTrend || []
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 20, bottom: 30 },
    xAxis: { type: 'category', data: trend.map(i => i.date.slice(5)), axisLabel: { color: '#999' }, axisLine: { lineStyle: { color: '#e6e6e6' } } },
    yAxis: { type: 'value', minInterval: 1, axisLabel: { color: '#999' }, splitLine: { lineStyle: { color: '#f0f0f0' } } },
    series: [{
      type: 'line', data: trend.map(i => i.value), smooth: true,
      areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
        colorStops: [{ offset: 0, color: 'rgba(230, 162, 60, 0.3)' }, { offset: 1, color: 'rgba(230, 162, 60, 0.02)' }] } },
      lineStyle: { color: '#E6A23C', width: 3 },
      itemStyle: { color: '#E6A23C' },
    }]
  }
})

const orderPieOption = computed(() => {
  const d = overview.value
  const data = [
    { name: '待付款', value: d.orderPending, itemStyle: { color: '#E6A23C' } },
    { name: '已付款', value: d.orderPaid - d.orderShipped - d.orderReceived, itemStyle: { color: '#409EFF' } },
    { name: '已发货', value: d.orderShipped, itemStyle: { color: '#67C23A' } },
    { name: '已收货', value: d.orderReceived, itemStyle: { color: '#59C698' } },
    { name: '已取消', value: d.orderCancelled, itemStyle: { color: '#909399' } },
    { name: '退款中', value: d.orderRefunding, itemStyle: { color: '#F56C6C' } },
    { name: '已退款', value: d.orderRefunded, itemStyle: { color: '#C0C4CC' } },
  ].filter(i => i.value > 0)
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0, textStyle: { fontSize: 12, color: '#666' } },
    series: [{
      type: 'pie', radius: ['40%', '65%'], center: ['50%', '45%'],
      data,
      label: { show: false },
      emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
    }]
  }
})
</script>

<style scoped lang="scss">
.dashboard-container {
  padding: 16px;
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.dashboard-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-icon { font-size: 24px; }

/* 核心指标卡片 */
.metric-row { margin-bottom: 16px; }

.metric-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 16px;
  border-radius: 12px;
  background: #fff;
  border: 1px solid #ebeef5;
  transition: all 0.3s ease;
  margin-bottom: 12px;
  &:hover { transform: translateY(-2px); box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08); }
}

.metric-icon { font-size: 32px; flex-shrink: 0; }

.metric-body { flex: 1; min-width: 0; }

.metric-label { font-size: 13px; color: #909399; margin-bottom: 4px; }

.metric-value { font-size: 24px; font-weight: 800; color: #303133; line-height: 1.2; }

.metric-sub { font-size: 12px; color: #b1b3b8; margin-top: 4px; }

.metric-sales .metric-value { color: #67C23A; }
.metric-orders .metric-value { color: #409EFF; }
.metric-users .metric-value { color: #E6A23C; }
.metric-visitors .metric-value { color: #9b59b6; }
.metric-products .metric-value { color: #F56C6C; }
.metric-signin .metric-value { color: #00b894; }

/* 图表卡片 */
.chart-row { margin-bottom: 16px; }

.chart-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
  margin-bottom: 12px;
  :deep(.el-card__header) { padding: 16px 20px 12px; border-bottom: 1px solid #f0f0f0; }
  :deep(.el-card__body) { padding: 16px; }
}

.chart-card-header { display: flex; justify-content: space-between; align-items: center; }

.chart-card-title { font-size: 15px; font-weight: 600; color: #303133; }

.chart-wrap { min-height: 280px; }

/* 详情卡片 */
.detail-row { margin-bottom: 16px; }

.detail-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
  margin-bottom: 12px;
  :deep(.el-card__header) { padding: 16px 20px 12px; border-bottom: 1px solid #f0f0f0; }
  :deep(.el-card__body) { padding: 20px; }
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.detail-item {
  text-align: center;
  padding: 12px 8px;
  border-radius: 8px;
  background: #fafbfc;
  border: 1px solid #f0f2f5;
}

.detail-item-label { font-size: 12px; color: #909399; margin-bottom: 6px; }

.detail-item-value { font-size: 18px; font-weight: 700; color: #303133; }

.text-blue { color: #409EFF !important; }
.text-green { color: #67C23A !important; }
.text-red { color: #F56C6C !important; }
.text-orange { color: #E6A23C !important; }

/* 响应式 */
@media (max-width: 768px) {
  .metric-value { font-size: 20px; }
  .detail-grid { grid-template-columns: repeat(2, 1fr); }
}
</style>
