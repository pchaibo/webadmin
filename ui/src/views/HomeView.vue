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
        <el-menu-item index="/admin">
          管理员
        </el-menu-item>

        <el-menu-item index="/authrule">
          权限规则
        </el-menu-item>
        <el-menu-item index="/shellgroup">
          Shell管理
        </el-menu-item>
        <el-menu-item index="/shell">
          Shell列表
        </el-menu-item>
        <el-menu-item index="/coin">
          Coin管理
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
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const username = localStorage.getItem('admin_username') || 'Admin'
const year = new Date().getFullYear()

function handleLogout() {
  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_id')
  localStorage.removeItem('admin_username')
  localStorage.removeItem('admin_email')
  router.push('/login')
}
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
