 <script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'
 import { fetchHeyueById, createHeyue, updateHeyue } from '../api'
 
 const route = useRoute()
 const router = useRouter()
 
 const isEditing = ref(false)
 const editingId = ref<number | null>(null)
 const loading = ref(true)
const submitting = ref(false)
 const form = reactive({
   username: '',
   symbol: '',
   side: 1,
   num: 3,
   status: 1,
   sellprice: 30,
   oneprice: 100,
   repeatprice: 50,
   rangetype: 2,
   rangeprice: 0,
   rangepercent: 100,
   rangeclosingpct: 50,
   rangeclosing: 1,
   closingprice: 0,
   risk: 2,
   risktime: 30,
 })
 
onMounted(async () => {
  const id = route.query.id
   if (id) {
     isEditing.value = true
     editingId.value = Number(id)
     try {
       const res = await fetchHeyueById(Number(id))
       if (res.status === 1 && res.heyue) {
         const d = res.heyue
         form.username = (d.username as string) || ''
         form.symbol = (d.symbol as string) || ''
         form.side = (d.side as number) || 1
         form.num = (d.num as number) || 5
         form.status = (d.status as number) || 1
         form.sellprice = (d.sellprice as number) || 0
         form.oneprice = (d.oneprice as number) || 0
         form.repeatprice = (d.repeatprice as number) || 0
         form.rangetype = (d.rangetype as number) || 2
         form.rangeprice = (d.rangeprice as number) || 0
         form.rangepercent = (d.rangepercent as number) || 100
         form.rangeclosingpct = (d.rangeclosingpct as number) || 50
         form.rangeclosing = (d.rangeclosing as number) || 1
         form.closingprice = (d.closingprice as number) || 0
         form.risk = (d.risk as number) || 2
         form.risktime = (d.risktime as number) || 30
      }
    } catch { /* ignore */ }
  }
 
  // 从 ?symbol= 取默认值（来自币种列表）
  if (!form.symbol) {
    if (route.query.symbol) {
      const sym = route.query.symbol as string
      form.symbol = sym
    }
  }
 
  loading.value = false
 })
 
 function goBack() {
   router.push('/heyue')
 }
 
 async function submitForm() {
   if (!form.symbol.trim()) {
     ElMessage({ message: '请输入合约代码', type: 'warning' })
     return
   }
   submitting.value = true
   try {
     const payload = {
       username: form.username,
       symbol: form.symbol,
       side: form.side,
       num: form.num,
       status: form.status,
       sellprice: form.sellprice,
       oneprice: form.oneprice,
       repeatprice: form.repeatprice,
       rangetype: form.rangetype,
       rangeprice: form.rangeprice,
       rangepercent: form.rangepercent,
       rangeclosingpct: form.rangeclosingpct,
       rangeclosing: form.rangeclosing,
       closingprice: form.closingprice,
       risk: form.risk,
       risktime: form.risktime,
     } as Record<string, unknown>
 
     let res: { status: number; error?: string }
     if (isEditing.value && editingId.value !== null) {
       res = await updateHeyue(editingId.value, payload)
     } else {
       res = await createHeyue(payload)
     }
     if (res.status === 1) {
       router.push('/heyue')
     } else {
       ElMessage({ message: res.error || '操作失败', type: 'error' })
     }
   } catch {
     ElMessage({ message: '网络错误', type: 'error' })
   } finally {
     submitting.value = false
   }
 }
 </script>
 
 <template>
   <div class="add-page" v-if="!loading">
     <div class="add-header">
       <button class="back-btn" @click="goBack">&larr;</button>
       <span class="add-title">{{ isEditing ? '编辑合约' : '新增合约' }}</span>
     </div>
 
     <div class="add-body">
       <!-- <div class="form-group">
         <label class="form-label">用户名</label>
         <input v-model="form.username" class="form-input" placeholder="请输入用户名" />
       </div> -->
      <div class="form-group">
        <label class="form-label">合约代码 <span class="required">*</span></label>
         <input v-model="form.symbol" class="form-input" placeholder="请输入合约代码" disabled  :readonly="isEditing" />
      </div>
       <div class="form-divider">方向</div>
       <div class="status-radios">
         <label class="radio-option" :class="{ active: form.side === 1 }">
           <input type="radio" v-model="form.side" :value="1" />
           <span>开多</span>
         </label>
         <label class="radio-option" :class="{ active: form.side === 2 }">
           <input type="radio" v-model="form.side" :value="2" />
           <span>开空</span>
         </label>
       </div>
       <div class="form-row">
         <div class="form-group flex-half">
           <label class="form-label">收益%</label>
           <input v-model.number="form.sellprice" type="number" step="0.1" class="form-input" />
         </div>
         <div class="form-group flex-half">
           <label class="form-label">首仓USDT</label>
           <input v-model.number="form.oneprice" type="number" step="0.01" class="form-input" />
         </div>
       </div>
       <div class="form-divider">网格设置</div>
       <div class="form-row">
         <div class="form-group flex-half">
           <label class="form-label">最大次数</label>
           <input v-model.number="form.num" type="number" class="form-input" />
         </div>
         <div class="form-group flex-half">
           <label class="form-label">补仓USDT</label>
           <input v-model.number="form.repeatprice" type="number" step="0.01" class="form-input" />
         </div>
       </div>
        <div class="form-row">
         <div class="form-group flex-half">
           <label class="form-label">网格%</label>
           <input v-model.number="form.rangepercent" type="number" step="1" class="form-input" />
         </div>
         <div class="form-group flex-half">
           <label class="form-label">网格平仓%</label>
           <input v-model.number="form.rangeclosingpct" type="number" step="1" class="form-input" />
         </div>
       </div>
       <div class="form-divider">网格平仓</div>
       <div class="status-radios">
         <label class="radio-option" :class="{ active: form.rangeclosing === 1 }">
           <input type="radio" v-model="form.rangeclosing" :value="1" />
           <span>不平</span>
         </label>
         <label class="radio-option" :class="{ active: form.rangeclosing === 2 }">
           <input type="radio" v-model="form.rangeclosing" :value="2" />
           <span>平仓</span>
         </label>
       </div>

      
       <!-- <div class="form-row">
         <div class="form-group flex-half">
           <label class="form-label">强平价格</label>
           <input v-model.number="form.closingprice" type="number" step="0.001" class="form-input" />
         </div>
         <div class="form-group flex-half">
           <label class="form-label">网格类型</label>
           <select v-model="form.rangetype" class="form-input">
             <option :value="1">差价USDT</option>
             <option :value="2">网格保证金</option>
           </select>
         </div>
       </div> -->


       <div class="form-divider">风控设置</div>
       <div class="form-row">
         <div class="form-group flex-half">
           <label class="form-label">风控</label>
           <select v-model="form.risk" class="form-input">
             <option :value="1">关闭</option>
             <option :value="2">开启</option>
           </select>
         </div>
         <div class="form-group flex-half">
           <label class="form-label">风控时间(分)</label>
           <input v-model.number="form.risktime" type="number" class="form-input" />
         </div>
       </div>
       <div class="form-divider">状态</div>
       <div class="status-radios">
         <label class="radio-option" :class="{ active: form.status === 1 }">
           <input type="radio" v-model="form.status" :value="1" />
           <span>启用</span>
         </label>
         <label class="radio-option" :class="{ active: form.status === 2 }">
           <input type="radio" v-model="form.status" :value="2" />
           <span>停用</span>
         </label>
       </div>
       

       <button class="submit-btn" :disabled="submitting" @click="submitForm">
         {{ submitting ? '提交中...' : (isEditing ? '保存' : '创建') }}
       </button>
     </div>
   </div>
 </template>
 
 <style scoped>
 .add-page {
   min-height: 100vh;
   background: #f5f5f5;
   display: flex;
   flex-direction: column;
 }
 .add-header {
   display: flex;
   align-items: center;
   padding: 12px 16px;
   background: #fff;
   border-bottom: 1px solid #eee;
   position: sticky;
   top: 0;
   z-index: 10;
 }
 .back-btn {
   font-size: 20px;
   border: none;
   background: none;
   color: #333;
   cursor: pointer;
   padding: 0 8px 0 0;
   line-height: 1;
 }
 .add-title {
   font-size: 17px;
   font-weight: 600;
   color: #1a1a1a;
 }
 .add-body {
   padding: 16px;
   flex: 1;
 }
 .form-group { margin-bottom: 14px; }
 .form-label {
   display: block;
   font-size: 13px;
   color: #666;
   margin-bottom: 5px;
   font-weight: 500;
 }
 .required { color: #e74c3c; }
 .form-input {
   width: 100%;
   height: 40px;
   padding: 0 12px;
   border: 1px solid #ddd;
   border-radius: 8px;
   font-size: 14px;
   outline: none;
   background: #fff;
   box-sizing: border-box;
   -webkit-appearance: none;
   appearance: none;
 }
 .form-input:focus { border-color: #409eff; }
 .form-row { display: flex; gap: 10px; }
.flex-half { flex: 1; min-width: 0; }
.status-radios {
   display: flex;
   gap: 10px;
   margin-bottom: 14px;
 }
 .radio-option {
   flex: 1;
   display: flex;
   align-items: center;
   justify-content: center;
   height: 40px;
   border: 1px solid #ddd;
   border-radius: 8px;
   font-size: 14px;
   color: #666;
   cursor: pointer;
   background: #fff;
   box-sizing: border-box;
 }
 .radio-option input { display: none; }
 .radio-option.active {
   border-color: #409eff;
   background: #ecf5ff;
   color: #409eff;
   font-weight: 500;
 }
 .form-divider {
   font-size: 14px;
   font-weight: 600;
   color: #333;
   padding: 10px 0 6px;
   margin: 4px 0 10px;
   border-bottom: 1px solid #eee;
 }
 .submit-btn {
   width: 100%;
   height: 44px;
   margin-top: 20px;
   border: none;
   border-radius: 8px;
   background: #409eff;
   color: #fff;
   font-size: 16px;
   font-weight: 500;
   cursor: pointer;
 }
 .submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }
 </style>

