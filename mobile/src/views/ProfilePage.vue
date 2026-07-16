 <script setup lang="ts">
 import { ref, onMounted } from 'vue'
 import { useRouter } from 'vue-router'
 import { fetchUserInfo, clearToken } from '../api'
 import { useUserStore } from '../stores/user'
 import type { UserInfo } from '../api'
 
 const router = useRouter()
 const userStore = useUserStore()
 const user = ref<UserInfo | null>(null)
 const loading = ref(true)
 
 onMounted(async () => {
   try {
     const res = await fetchUserInfo()
     if (res.status === 1 && res.data) {
       user.value = res.data as UserInfo
     }
   } catch { /* ignore */ } finally {
     loading.value = false
   }
 })
 
 function handleLogout() {
   clearToken()
   userStore.logout()
   router.push('/login')
 }
 </script>
 
 <template>
   <div class="page">
     <h2 class="page-title">我的</h2>
     <div v-if="loading" class="loading">加载中...</div>
     <div v-else class="profile-card">
       <div class="profile-avatar">{{ user?.username?.charAt(0).toUpperCase() ?? 'U' }}</div>
       <div class="profile-name">{{ user?.username ?? '--' }}</div>
       <div class="profile-email">{{ user?.email ?? '--' }}</div>
      <div v-if="user?.usdt !== undefined" class="profile-balance">
        USDT: {{ user.usdt }}
      </div>
       <div class="profile-menu">
         <button class="menu-btn" @click="router.push('/config')">密码设置</button>
         <button class="menu-btn" @click="router.push('/binan')">币安设置</button>
       </div>
      <button class="logout-btn" @click="handleLogout">退出登录</button>
     </div>
   </div>
 </template>
 
 <style scoped>
 .page { padding: 16px 16px 80px; }
 .page-title { font-size: 18px; font-weight: 600; margin: 0 0 16px; color: #1a1a1a;text-align: center; }
 .loading { text-align: center; color: #999; padding: 40px 0; font-size: 14px; }
 .profile-card {
   background: #fff; border-radius: 12px; padding: 32px 24px;
   display: flex; flex-direction: column; align-items: center;
   box-shadow: 0 1px 4px rgba(0,0,0,0.06);
 }
 .profile-avatar {
   width: 64px; height: 64px; border-radius: 50%;
   background: #4a90d9; color: #fff;
   display: flex; align-items: center; justify-content: center;
   font-size: 28px; font-weight: 600; margin-bottom: 12px;
 }
 .profile-name { font-size: 18px; font-weight: 600; color: #1a1a1a; margin-bottom: 4px; }
 .profile-email { font-size: 14px; color: #999; margin-bottom: 8px; }
.profile-balance { font-size: 16px; font-weight: 500; color: #e67e22; margin-bottom: 24px; }
 .profile-menu {
   width: 100%;
   display: flex;
   flex-direction: column;
   gap: 10px;
   margin-bottom: 24px;
 }
 .menu-btn {
   width: 100%;
   height: 42px;
   border: 1px solid #ddd;
   border-radius: 8px;
   background: #fff;
   color: #333;
   font-size: 15px;
   cursor: pointer;
   text-align: center;
 }
 .menu-btn:active { background: #f5f5f5; }
.logout-btn {
   padding: 10px 40px; background: #e74c3c; color: #fff;
   border: none; border-radius: 8px; font-size: 15px; cursor: pointer;
   transition: background 0.2s;
 }
 .logout-btn:active { background: #c0392b; }
 </style>
