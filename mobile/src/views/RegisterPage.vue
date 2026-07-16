 <script setup lang="ts">
 import { ref } from 'vue'
 import { useRouter } from 'vue-router'
 import { ElMessage } from 'element-plus'
import { register } from '../api'
 
 const router = useRouter()
 
 const username = ref('')
const password = ref('')
 const confirmPassword = ref('')
const email = ref('')
 const mobile = ref('')
 const loading = ref(false)
 const error = ref('')
 
 async function handleRegister() {
  if (!username.value || !password.value || !email.value || !mobile.value) {
    error.value = '请填写所有字段'
    return
  }
   if (password.value !== confirmPassword.value) {
     error.value = '两次输入的密码不一致'
     return
   }
  loading.value = true
   error.value = ''
   try {
     const res = await register({
       username: username.value,
       password: password.value,
       email: email.value,
       mobile: mobile.value,
     })
     if (res.status === 1) {
       ElMessage({ message: '注册成功，请登录', type: 'success' })
       router.push('/login')
     } else {
       error.value = res.error || '注册失败'
     }
   } catch {
     error.value = '网络错误，请重试'
   } finally {
     loading.value = false
   }
 }
 </script>
 
 <template>
   <div class="register-page">
     <div class="register-card">
       <h1 class="register-title">注册</h1>
       <p class="register-sub">创建您的账号</p>
       <input v-model="username" class="register-input" type="text" placeholder="用户名" />
      <input v-model="password" class="register-input" type="password" placeholder="密码" />
       <input v-model="confirmPassword" class="register-input" type="password" placeholder="确认密码" />
      <input v-model="email" class="register-input" type="email" placeholder="邮箱" />
       <input v-model="mobile" class="register-input" type="text" placeholder="手机号" />
       <p v-if="error" class="register-error">{{ error }}</p>
       <button class="register-btn" :disabled="loading" @click="handleRegister">
         {{ loading ? '注册中...' : '注 册' }}
       </button>
       <p class="login-link">已有账号？<a href="/login">去登录</a></p>
     </div>
   </div>
 </template>
 
 <style scoped>
 .register-page {
   display: flex; align-items: center; justify-content: center;
   height: 100vh; background: #f5f5f5; padding: 16px; box-sizing: border-box;
 }
 .register-card {
   width: 100%; max-width: 340px; background: #fff;
   border-radius: 12px; padding: 32px 24px;
   box-shadow: 0 2px 12px rgba(0,0,0,0.08);
 }
 .register-title {
   margin: 0 0 4px; font-size: 24px; font-weight: 600;
   color: #1a1a1a; text-align: center;
 }
 .register-sub {
   margin: 0 0 24px; font-size: 14px; color: #999; text-align: center;
 }
 .register-input {
   width: 100%; padding: 12px 14px; margin-bottom: 14px;
   border: 1px solid #e0e0e0; border-radius: 8px; font-size: 15px;
   outline: none; box-sizing: border-box; transition: border-color 0.2s;
 }
 .register-input:focus { border-color: #4a90d9; }
 .register-error {
   color: #e74c3c; font-size: 13px; margin: 0 0 12px; text-align: center;
 }
 .register-btn {
   width: 100%; padding: 12px; background: #4a90d9; color: #fff;
   border: none; border-radius: 8px; font-size: 16px; font-weight: 500;
   cursor: pointer; transition: background 0.2s;
 }
 .register-btn:disabled { background: #a0c4f0; cursor: not-allowed; }
 .register-btn:not(:disabled):active { background: #357abd; }
 .login-link {
   margin-top: 16px; font-size: 14px; color: #999; text-align: center;
 }
 .login-link a { color: #4a90d9; text-decoration: none; }
 </style>

