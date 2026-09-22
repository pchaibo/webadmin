<template>
  <div class="layout">
    <!-- Header -->
    <header class="header">
      <div class="header-left">
        <span class="header-logo">后台管理系统</span>
      </div>
      <div class="header-right">
        <span class="header-user">{{ username }}</span>
        <el-button text size="small" class="logout-btn" @click="handleLogout">退出登录</el-button>
      </div>
    </header>

    <div class="body">
    <!-- Sidebar -->
    <aside class="sidebar">
      <el-menu
        :default-active="route.path"
        router
        background-color="#001529"
        text-color="#fff"
        active-text-color="#409eff"
        class="menu"
      >
        <el-menu-item index="/index">
          首页
        </el-menu-item>
        <el-menu-item index="/user">
          用户管理
        </el-menu-item>
        <el-menu-item index="/admin">
          管理员
        </el-menu-item>
        <el-menu-item index="/admingroup">
          用户组管理
        </el-menu-item>
        
        <el-menu-item index="/authrule">
          权限规则
        </el-menu-item>
        <!--
         <el-menu-item index="/shellgroup">
          Shell分组
        </el-menu-item>
        <el-menu-item index="/shell">
          Shell列表
        </el-menu-item> -->

       <el-menu-item index="/coin">
        Coin列表
      </el-menu-item>
        <el-menu-item index="/heyue">
         合约列表
       </el-menu-item>
       <el-menu-item index="/heyueorder">
         合约日志
       </el-menu-item>
       <el-menu-item index="/task">
         预警价格
       </el-menu-item>
        
      </el-menu>
    </aside>

    <!-- Main content -->
    <main class="main">
      <router-view />
    </main>
    </div>

    <!-- Footer -->
    <footer class="footer">
      &copy; {{ year }} 后台管理系统
    </footer>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { checkTaskPrices } from '@/api/api'

const route = useRoute()
const router = useRouter()
const username = localStorage.getItem('admin_username') || 'Admin'
const year = new Date().getFullYear()
const triggeredTaskIds = new Set<number>()
const priceAudio = new Audio(`${import.meta.env.BASE_URL}price.mp3`)
priceAudio.preload = 'auto'
const PRICE_ALERT_REPEAT_COUNT = 3

let checkTimer: number | undefined
let checkingPrices = false
let audioUnlocked = false
let unlockingAudio = false
let pendingPriceAlert = false

function removeAudioUnlockListeners() {
  document.removeEventListener('click', unlockPriceAudio)
  document.removeEventListener('keydown', unlockPriceAudio)
}

async function unlockPriceAudio() {
  if (audioUnlocked || unlockingAudio) return

  unlockingAudio = true
  try {
    priceAudio.muted = false
    priceAudio.currentTime = 0
    await priceAudio.play()
    priceAudio.pause()
    priceAudio.currentTime = 0
    audioUnlocked = true
    removeAudioUnlockListeners()

    if (pendingPriceAlert) {
      pendingPriceAlert = false
      await playPriceAlert()
    }
  } catch {
    // Keep the listeners installed so the next user gesture retries unlocking.
  } finally {
    unlockingAudio = false
  }
}

async function playPriceAlert(): Promise<boolean> {
  try {
    for (let index = 0; index < PRICE_ALERT_REPEAT_COUNT; index += 1) {
      await new Promise<void>((resolve, reject) => {
        const cleanup = () => {
          priceAudio.removeEventListener('ended', handleEnded)
          priceAudio.removeEventListener('error', handleError)
        }
        const handleEnded = () => {
          cleanup()
          resolve()
        }
        const handleError = () => {
          cleanup()
          reject(new Error('Failed to play price alert'))
        }

        priceAudio.addEventListener('ended', handleEnded)
        priceAudio.addEventListener('error', handleError)
        priceAudio.currentTime = 0
        priceAudio.muted = false
        void priceAudio.play().catch((error) => {
          cleanup()
          reject(error)
        })
      })
    }
    return true
  } catch {
    return false
  }
}

async function pollPriceAlerts() {
  if (checkingPrices) return
  checkingPrices = true

  try {
    const res = await checkTaskPrices()
    if (res.status !== 1 || !Array.isArray(res.data)) return

    const currentIds = new Set(res.data.map((item) => item.id))
    const newIds = [...currentIds].filter((id) => !triggeredTaskIds.has(id))

    triggeredTaskIds.clear()
    currentIds.forEach((id) => triggeredTaskIds.add(id))

    if (newIds.length > 0) {
      pendingPriceAlert = true
      if (await playPriceAlert()) {
        pendingPriceAlert = false
      } else {
        newIds.forEach((id) => triggeredTaskIds.delete(id))
      }
    }
  } catch {
    // Ignore transient polling errors and retry on the next interval.
  } finally {
    checkingPrices = false
  }
}

function handleLogout() {
  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_id')
  localStorage.removeItem('admin_username')
  localStorage.removeItem('admin_email')
  router.push('/login')
}

onMounted(() => {
  document.addEventListener('click', unlockPriceAudio)
  document.addEventListener('keydown', unlockPriceAudio)
  pollPriceAlerts()
  checkTimer = window.setInterval(pollPriceAlerts, 5000)
})

onBeforeUnmount(() => {
  removeAudioUnlockListeners()
  if (checkTimer !== undefined) {
    window.clearInterval(checkTimer)
  }
  priceAudio.pause()
})
</script>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
  padding: 0 20px;
  background: #001529;
  color: #fff;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-logo {
  font-size: 18px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-user {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.85);
}

.body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 220px;
  background: #001529;
  flex-shrink: 0;
}

.menu {
  border-right: none;
  display: block;
}

.main {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background: #f5f5f5;
}

.footer {
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #001529;
  color: rgba(255, 255, 255, 0.6);
  font-size: 13px;
  flex-shrink: 0;
}

.logout-btn {
  color: rgba(255, 255, 255, 0.7);
}
.logout-btn:hover {
  color: #fff;
}
</style>
