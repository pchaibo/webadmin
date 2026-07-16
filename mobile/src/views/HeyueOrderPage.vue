 <script setup lang="ts">
 import { ref, onMounted } from 'vue'
 import { fetchHeyueOrderList } from '../api'
 
 interface HeyueOrderItem {
   id: number
   ordertype: number
   side: number
   symbol: string
   price: number
   quantity: number
   total: number
   num: number
   usdt: number
   addtime: number
   [key: string]: unknown
 }
 
 const items = ref<HeyueOrderItem[]>([])
 const loading = ref(false)
 const total = ref(0)
 const totalUsdt = ref(0)
 const totalUsdtLong = ref(0)
 const totalUsdtShort = ref(0)
const currentPage = ref(1)
const pageSize = 20
 
 const searchSymbol = ref('')
 const searchOrdertype = ref<number | undefined>()
 const searchSide = ref<number | undefined>()
 
 onMounted(() => { fetchItems() })
 
 async function fetchItems() {
   loading.value = true
   try {
     const res = await fetchHeyueOrderList(currentPage.value, searchSymbol.value || undefined, searchOrdertype.value, searchSide.value)
     if (res.status === 1) {
       items.value = (res.data || []) as HeyueOrderItem[]
       total.value = res.total || 0
       totalUsdt.value = res.total_usdt || 0
       totalUsdtLong.value = res.total_usdt_long || 0
       totalUsdtShort.value = res.total_usdt_short || 0
     }
   } catch { /* ignore */ } finally {
     loading.value = false
   }
 }
 
 function handleSearch() {
   currentPage.value = 1
   fetchItems()
 }
 
function goPage(p: number) {
   const tp = Math.ceil(total.value / pageSize) || 1
   if (p < 1 || p > tp) return
  currentPage.value = p
   fetchItems()
 }
 
 function formatTime(ts: number): string {
   if (!ts) return '-'
   const d = new Date(ts * 1000)
   const pad = (n: number) => String(n).padStart(2, '0')
   return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
 }
 </script>
 
 <template>
  <div class="page">
    <h2 class="page-title">订单日志</h2>
 
 
    <!-- Stats -->
     <div class="stats-bar">
       <div class="stat-item">
         <span class="stat-label">总收益</span>
         <span class="stat-value" :class="totalUsdt >= 0 ? 'profit' : 'loss'">{{ totalUsdt.toFixed(2) }}</span>
       </div>
       <div class="stat-item">
         <span class="stat-label">做多收益</span>
         <span class="stat-value long">{{ totalUsdtLong.toFixed(2) }}</span>
       </div>
       <div class="stat-item">
         <span class="stat-label">做空收益</span>
         <span class="stat-value short">{{ totalUsdtShort.toFixed(2) }}</span>
       </div>
     </div>
     
     <!-- Toolbar -->
     <div class="toolbar">
       <input v-model="searchSymbol" class="search-input" placeholder="合约代码" @keyup.enter="handleSearch" />
        <select v-model="searchOrdertype" class="search-select" @change="handleSearch">
         <option :value="undefined">所有类型</option>
         <option :value="1">开仓</option>
         <option :value="2">平仓</option>
       </select> 
       <select v-model="searchSide" class="search-select" @change="handleSearch">
         <option :value="undefined">所有方向</option>
         <option :value="1">开多</option>
         <option :value="2">开空</option>
       </select> 
       <button class="btn btn-search" @click="handleSearch">搜索</button>
     </div>
 
     <div v-if="loading" class="loading">加载中...</div>
     <div v-else-if="items.length === 0" class="empty">暂无数据</div>
     <div v-else class="order-list">
       <div v-for="item in items" :key="item.id" class="order-item">
         <div class="order-top">
           <span class="order-symbol">{{ item.symbol }}</span>
           <span class="order-type" :class="item.ordertype === 1 ? 'type-open' : 'type-close'">
             {{ item.ordertype === 1 ? '开仓' : '平仓' }}
           </span>
           <span class="order-side" :class="item.side === 1 ? 'side-long' : 'side-short'">
             {{ item.side === 1 ? '开多' : '开空' }}
           </span>
           <span class="order-num">第{{ item.num }}次</span>
         </div>
         <div class="order-grid">
           <div class="grid-item">
             <span class="grid-label">价格</span>
             <span class="grid-value">{{ item.price.toFixed(4) }}</span>
           </div>
           <div class="grid-item">
             <span class="grid-label">数量</span>
             <span class="grid-value">{{ item.quantity }}</span>
           </div>
           <div class="grid-item">
             <span class="grid-label">金额</span>
             <span class="grid-value">{{ item.total ? item.total.toFixed(2) : '-' }}</span>
           </div>
           <div class="grid-item">
             <span class="grid-label">收益</span>
             <span class="grid-value" :class="item.usdt >= 0 ? 'profit' : 'loss'">
               {{ item.usdt ? item.usdt.toFixed(2) : '0.00' }}
             </span>
           </div>
           <div class="grid-item">
             <span class="grid-label">时间</span>
             <span class="grid-value time">{{ formatTime(item.addtime) }}</span>
           </div>
         </div>
       </div>
     </div>
 
    <div v-if="total > pageSize" class="pagination">
     <button :disabled="currentPage <= 1" @click="goPage(currentPage - 1)">上一页</button>
      <span class="page-info">{{ currentPage }} / {{ Math.ceil(total / pageSize) || 1 }}</span>
       <button :disabled="currentPage >= (Math.ceil(total / pageSize) || 1)" @click="goPage(currentPage + 1)">下一页</button>
    </div>
   </div>
 </template>
 
 <style scoped>
.page { padding: 12px 12px 80px; }
.page-title { font-size: 18px; font-weight: 600; margin: 0 0 10px; color: #1a1a1a;text-align: center; }
.loading, .empty { text-align: center; color: #999; padding: 40px 0; font-size: 14px; }
 
 /* Toolbar */
 .toolbar {
   display: flex;
   flex-wrap: wrap;
   gap: 6px;
   margin-bottom: 12px;
 }
 .search-input, .search-select {
   height: 34px;
   padding: 0 10px;
   border: 1px solid #ddd;
   border-radius: 6px;
   font-size: 13px;
   outline: none;
   background: #fff;
   box-sizing: border-box;
 }
.search-input { flex: 0 0 80px; min-width: 0; } 
 /*.search-select { flex: 0 0 90px; }*/
 .btn {
   display: inline-flex; align-items: center; justify-content: center;
   height: 34px; padding: 0 12px; border: none; border-radius: 6px;
   font-size: 13px; cursor: pointer; white-space: nowrap;
 }
 .btn-search { background: #67c23a; color: #fff; }
 
 /* Stats Bar */
 .stats-bar {
   display: flex;
   gap: 6px;
   margin-bottom: 12px;
   padding: 12px;
   background: #fff;
   border-radius: 8px;
   box-shadow: 0 1px 4px rgba(0,0,0,0.06);
 }
 .stat-item {
   flex: 1;
   display: flex;
   flex-direction: column;
   align-items: center;
   gap: 4px;
 }
 .stat-label { font-size: 11px; color: #999; }
 .stat-value { font-size: 15px; font-weight: 700; color: #333; }
 .stat-value.profit { color: #27ae60; }
 .stat-value.loss { color: #e74c3c; }
 .stat-value.long { color: #e74c3c; }
 .stat-value.short { color: #e6a23c; }
 
 /* List */
 .order-list { display: flex; flex-direction: column; gap: 8px; }
 .order-item {
   padding: 12px;
   background: #fff;
   border-radius: 8px;
   box-shadow: 0 1px 4px rgba(0,0,0,0.06);
 }
 .order-top {
   display: flex;
   align-items: center;
   gap: 6px;
   margin-bottom: 8px;
   flex-wrap: wrap;
 }
 .order-symbol { font-size: 15px; font-weight: 600; color: #1a1a1a; }
 .order-type, .order-side {
   font-size: 12px; padding: 2px 8px; border-radius: 4px; font-weight: 500;
 }
 .type-open { color: #409eff; background: #ecf5ff; }
 .type-close { color: #e6a23c; background: #fdf6ec; }
 .side-long { color: #27ae60; background: #e8f8f0; }
 .side-short { color: #e74c3c; background: #fdecea; }
 .order-num { font-size: 12px; color: #999; margin-left: auto; }
 .order-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 6px; }
 .grid-item { display: flex; flex-direction: column; align-items: center; gap: 2px; }
 .grid-label { font-size: 11px; color: #999; }
 .grid-value { font-size: 13px; font-weight: 600; color: #333; }
 .grid-value.profit { color: #27ae60; }
 .grid-value.loss { color: #e74c3c; }
 .grid-value.time { font-size: 11px; font-weight: 400; color: #999; }
 
 /* Pagination */
 .pagination {
   display: flex;
   align-items: center;
   justify-content: center;
   gap: 12px;
   margin-top: 14px;
   padding: 12px 0;
 }
 .pagination button {
   height: 32px;
   padding: 0 14px;
   border: 1px solid #ddd;
   border-radius: 6px;
   background: #fff;
   color: #333;
   font-size: 13px;
   cursor: pointer;
 }
 .pagination button:disabled { color: #ccc; cursor: not-allowed; }
.page-info { font-size: 13px; color: #666; }
 </style>
