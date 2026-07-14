<template>
  <div class="index-page">
    <h2 class="page-title">系统概览</h2>
    <div class="stats-row">
      <!-- <div class="stat-card clickable" @click="goAdmin">
        <div class="stat-value">{{ adminCount }}</div>
        <div class="stat-label">管理员</div>
      </div> -->
      <!-- <div class="stat-card clickable" @click="goShell">
        <div class="stat-value">{{ shellCount }}</div>
        <div class="stat-label">Shell 记录</div>
      </div> -->
      <div class="stat-card clickable" @click="goUser">
        <div class="stat-value">{{ userCount }}</div>
        <div class="stat-label">用户</div>
      </div>
      <div class="stat-card clickable" @click="goCoin">
        <div class="stat-value">{{ coinCount }}</div>
        <div class="stat-label">币种</div>
      </div>
      <div class="stat-card clickable" @click="goHeyue" >
        <div class="stat-value">{{ heyueCount }}</div>
        <div class="stat-label">合约</div>
      </div>
      <div class="stat-card clickable" @click="goHeyueorder">
        <div class="stat-value usdt">
          {{ heyueUsdt.toFixed(2) }}
        </div>
        <div class="stat-label">做多usdt</div>
      </div>
     <div class="stat-card clickable" @click="goHeyueorder">
       <div class="stat-value short-usdt">
         {{ heyueShortUsdt.toFixed(2) }}
       </div>
       <div class="stat-label">做空usdt</div>
     </div>
      <div class="stat-card clickable" @click="goHeyueorder">
        <div class="stat-value total-usdt">
          {{ (heyueUsdt + heyueShortUsdt).toFixed(2) }}
        </div>
        <div class="stat-label">总收益usdt</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboardStats } from "@/api/api"

const router = useRouter()
const adminCount = ref(0)
const shellCount = ref(0)
const heyueUsdt = ref(0)
const heyueShortUsdt = ref(0)
const coinCount = ref(0)
const heyueCount = ref(0)
const userCount = ref(0)

function goShell() {
  router.push('/shell')
}
function goAdmin() {
  router.push("/admin")
}
function goUser() {
  router.push("/user")
}
function goHeyue() {
  router.push("/heyue")
}
function goHeyueorder() {
  router.push("/heyueorder")
}
function goCoin() {
  router.push("/coin")
}
onMounted(async () => {
  try {
    const res = await getDashboardStats()
    if (res.status === 1) {
      adminCount.value = res.admin_count
      shellCount.value = res.shell_count
      userCount.value = res.user_count
      coinCount.value = res.coin_count
      heyueCount.value = res.heyue_count
      heyueUsdt.value = res.heyueorder_usdt
      heyueShortUsdt.value = res.heyueorder_short_usdt
    }
  } catch {
    // silent
  }
})
</script>

<style scoped>
.index-page {
  padding: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 24px;
}

.stats-row {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 32px 40px;
  min-width: 180px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}

.clickable {
  cursor: pointer;
  transition: transform .15s, box-shadow .15s;
}
.clickable:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}

.stat-value {
  font-size: 40px;
  font-weight: 700;
  color: #409eff;
  line-height: 1;
}

.stat-label {
  margin-top: 10px;
  font-size: 14px;
  color: #666;
}

.stat-value.usdt {
  color: #67c23a;
}

.stat-value.short-usdt {
  color: #e6a23c;
}

.stat-value.total-usdt {
  color: #f56c6c;
}
</style>
