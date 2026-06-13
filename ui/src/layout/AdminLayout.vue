<template>
  <div class="layout">
    <!-- 左侧菜单 -->
    <aside class="sidebar">
      <div class="logo">
        Admin
      </div>

      <el-menu
        :default-active="activeMenu"
        router
        class="menu"
      >
        <el-menu-item index="/admin">
          <span>管理员</span>
        </el-menu-item>

        <el-menu-item index="/auth-rule">
          <span>权限规则</span>
        </el-menu-item>
      </el-menu>
    </aside>

    <!-- 右侧内容 -->
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const activeMenu = computed(() => route.path)
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
}

.sidebar {
  width: 220px;
  background: #001529;
  color: #fff;
}

.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid rgba(255,255,255,.1);
}

.menu {
  border-right: none;
}

.content {
  flex: 1;
  padding: 2px;
  background: #f5f5f5;
  overflow: auto;
}
</style>

