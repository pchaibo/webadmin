<template>
  <div class="index-page">
    <h2 class="page-title">系统概览</h2>
    <div class="stats-row">
      <div class="stat-card clickable" @click="goAdmin">
        <div class="stat-value">{{ adminCount }}</div>
        <div class="stat-label">管理员</div>
      </div>
      <div class="stat-card clickable" @click="goShell">
        <div class="stat-value">{{ shellCount }}</div>
        <div class="stat-label">Shell 记录</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getAdmins, getShells } from '@/api/api'

const router = useRouter()
const adminCount = ref(0)
const shellCount = ref(0)

function goShell() {
  router.push('/shell')
}
function goAdmin() {
  router.push('/admin')
}

onMounted(async () => {
  try {
    const [admins, shells] = await Promise.all([
      getAdmins(1),
      getShells(1),
    ])
    if (admins.status === 1) adminCount.value = admins.total
    if (shells.status === 1) shellCount.value = shells.total
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
</style>
