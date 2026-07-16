 <script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
 import { fetchUserInfo, updateBinance, type UserInfo } from '../api'
 
 const router = useRouter()
 const bnaccess = ref('')
const bnasecret = ref('')
 const margin = ref(0)
const submitting = ref(false)
 const loading = ref(true)
 
 onMounted(async () => {
   try {
     const res = await fetchUserInfo()
     if (res.status === 1 && res.data) {
       const d = res.data as UserInfo
       bnaccess.value = d.bnaccess || ''
       bnasecret.value = d.bnasecret || ''
       margin.value = d.margin || 0
   }
   } catch { /* ignore */ } finally {
     loading.value = false
   }
 })
 
 async function handleSubmit() {
   if (!bnaccess.value.trim() && !bnasecret.value.trim()) {
     ElMessage({ message: '请输入 API Key 或 Secret Key', type: 'warning' })
     return
   }
   submitting.value = true
   try {
     const res = await updateBinance(bnaccess.value, bnasecret.value, margin.value)
     if (res.status === 1) {
       ElMessage({ message: '币安设置更新成功', type: 'success' })
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
       <span class="page-title">币安设置</span>
     </div>
     <div class="form-body" v-if="!loading">
       <div class="form-group">
        <label class="form-label">API Key</label>
         <textarea v-model="bnaccess" class="form-textarea" placeholder="请输入币安 API Key" rows="2"></textarea>
      </div>
      <div class="form-group">
        <label class="form-label">Secret Key</label>
       <textarea v-model="bnasecret" class="form-textarea" placeholder="请输入币安 Secret Key" rows="3"></textarea>
    </div>
       <div class="form-group">
         <label class="form-label">保证金%</label>
         <input v-model.number="margin" type="number" step="0.01" class="form-input" placeholder="请输入保证金百分比" />
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
 .form-textarea {
   width: 100%; padding: 10px 12px;
   border: 1px solid #ddd; border-radius: 8px; font-size: 14px;
   outline: none; background: #fff; box-sizing: border-box;
   resize: vertical; min-height: 60px; font-family: inherit;
 }
 .form-textarea:focus { border-color: #409eff; }
.submit-btn {
   width: 100%; height: 44px; margin-top: 16px;
   border: none; border-radius: 8px;
   background: #409eff; color: #fff;
   font-size: 16px; font-weight: 500; cursor: pointer;
 }
 .submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }
 </style>

