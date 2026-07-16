 <script setup lang="ts">
 import { ref } from 'vue'
 import { ElMessage } from 'element-plus'
 import { useRouter } from 'vue-router'
 import { updatePassword } from '../api'
 
 const router = useRouter()
 const password = ref('')
 const confirmPassword = ref('')
 const submitting = ref(false)
 
 async function handleSubmit() {
   if (!password.value.trim()) {
     ElMessage({ message: '请输入密码', type: 'warning' })
     return
   }
   if (password.value !== confirmPassword.value) {
     ElMessage({ message: '两次输入的密码不一致', type: 'warning' })
     return
   }
   submitting.value = true
   try {
     const res = await updatePassword(password.value)
     if (res.status === 1) {
       ElMessage({ message: '密码更新成功', type: 'success' })
       router.push('/profile')
     } else {
       ElMessage({ message: res.error || '更新失败', type: 'error' })
     }
   } catch {
     ElMessage({ message: '网络错误', type: 'error' })
   } finally {
     submitting.value = false
   }
 }
 </script>
 
 <template>
   <div class="page">
     <div class="page-header">
       <button class="back-btn" @click="router.push('/profile')">&larr;</button>
       <span class="page-title">密码设置</span>
     </div>
     <div class="form-body">
       <div class="form-group">
         <label class="form-label">新密码</label>
         <input v-model="password" type="password" class="form-input" placeholder="请输入新密码" />
       </div>
       <div class="form-group">
         <label class="form-label">确认密码</label>
         <input v-model="confirmPassword" type="password" class="form-input" placeholder="请再次输入新密码" />
       </div>
       <button class="submit-btn" :disabled="submitting" @click="handleSubmit">
         {{ submitting ? '提交中...' : '保存' }}
       </button>
     </div>
   </div>
 </template>
 
 <style scoped>
 .page {
   min-height: 100vh; background: #f5f5f5;
   display: flex; flex-direction: column;
 }
 .page-header {
   display: flex; align-items: center; padding: 12px 16px;
   background: #fff; border-bottom: 1px solid #eee;
   position: sticky; top: 0; z-index: 10;
 }
 .back-btn {
   font-size: 20px; border: none; background: none;
   color: #333; cursor: pointer; padding: 0 8px 0 0; line-height: 1;
 }
 .page-title { font-size: 17px; font-weight: 600; color: #1a1a1a; }
 .form-body { padding: 16px; flex: 1; }
 .form-group { margin-bottom: 16px; }
 .form-label {
   display: block; font-size: 13px; color: #666;
   margin-bottom: 5px; font-weight: 500;
 }
 .form-input {
   width: 100%; height: 40px; padding: 0 12px;
   border: 1px solid #ddd; border-radius: 8px; font-size: 14px;
   outline: none; background: #fff; box-sizing: border-box;
 }
 .form-input:focus { border-color: #409eff; }
 .submit-btn {
   width: 100%; height: 44px; margin-top: 16px;
   border: none; border-radius: 8px;
   background: #409eff; color: #fff;
   font-size: 16px; font-weight: 500; cursor: pointer;
 }
 .submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }
 </style>


