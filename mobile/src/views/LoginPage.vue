 <script setup lang="ts">
 import { ref } from 'vue'
 import { useRouter } from 'vue-router'
 import { login } from '../api'
 
 const router = useRouter()
 
 const username = ref('')
 const password = ref('')
 const loading = ref(false)
 const error = ref('')
 
 async function handleLogin() {
   if (!username.value || !password.value) {
     error.value = '请输入用户名和密码'
     return
   }
   loading.value = true
   error.value = ''
   try {
     const res = await login({ username: username.value, password: password.value })
     if (res.status === 1) {
       // 登录成功, router guard 会自动跳转到 /coin
       router.push('/coin')
     } else {
       error.value = (res as any).error || '登录失败'
     }
   } catch {
     error.value = '网络错误，请重试'
   } finally {
     loading.value = false
   }
 }
 </script>
 
 <template>
   <div class="login-page">
     <div class="login-card">
       <h1 class="login-title">登录</h1>
       <p class="login-sub">请输入您的账号信息</p>
       <input v-model="username" class="login-input" type="text" placeholder="用户名" @keyup.enter="handleLogin" />
       <input v-model="password" class="login-input" type="password" placeholder="密码" @keyup.enter="handleLogin" />
       <p v-if="error" class="login-error">{{ error }}</p>
      <button class="login-btn" :disabled="loading" @click="handleLogin">
        {{ loading ? '登录中...' : '登 录' }}
      </button>
       <p class="register-link">没有账号？<a href="/register" @click="router.push('/register')">去注册</a></p>
    </div>
   </div>
 </template>
 
 <style scoped>
 .login-page {
   display: flex;
   align-items: center;
   justify-content: center;
   height: 100vh;
   background: #f5f5f5;
   padding: 16px;
   box-sizing: border-box;
 }
 .login-card {
   width: 100%;
   max-width: 340px;
   background: #fff;
   border-radius: 12px;
   padding: 32px 24px;
   box-shadow: 0 2px 12px rgba(0,0,0,0.08);
 }
 .login-title {
   margin: 0 0 4px;
   font-size: 24px;
   font-weight: 600;
   color: #1a1a1a;
   text-align: center;
 }
 .login-sub {
   margin: 0 0 24px;
   font-size: 14px;
   color: #999;
   text-align: center;
 }
 .login-input {
   width: 100%;
   padding: 12px 14px;
   margin-bottom: 14px;
   border: 1px solid #e0e0e0;
   border-radius: 8px;
   font-size: 15px;
   outline: none;
   box-sizing: border-box;
   transition: border-color 0.2s;
 }
 .login-input:focus {
   border-color: #4a90d9;
 }
 .login-error {
   color: #e74c3c;
   font-size: 13px;
   margin: 0 0 12px;
   text-align: center;
 }
 .login-btn {
   width: 100%;
   padding: 12px;
   background: #4a90d9;
   color: #fff;
   border: none;
   border-radius: 8px;
   font-size: 16px;
   font-weight: 500;
   cursor: pointer;
   transition: background 0.2s;
 }
 .login-btn:disabled {
   background: #a0c4f0;
   cursor: not-allowed;
 }
.login-btn:not(:disabled):active {
  background: #357abd;
}
 .register-link {
   margin-top: 16px; font-size: 14px; color: #999; text-align: center;
 }
 .register-link a { color: #4a90d9; text-decoration: none; }
</style>
