<script setup lang="ts">
 import { ref, onMounted } from 'vue'
 import { useRouter } from 'vue-router'
 import { ElMessage, ElMessageBox } from 'element-plus'
import { fetchHeyueList, deleteHeyue } from '../api'
 
 const router = useRouter()
 
 interface HeyueItem {
   id: number
   userid: number
   username: string
   symbol: string
   side: number
   num: number
   is_num: number
   status: number
   sellprice: number
   oneprice: number
   repeatprice: number
   rangetype: number
   rangeprice: number
   rangepercent: number
   rangeclosingpct: number
   rangeclosing: number
   closingprice: number
   risk: number
   risktime: number
   newprice: number
   newtime: number
   addtime: number
   updatetime: number
   [key: string]: unknown
 }
 
const list = ref<HeyueItem[]>([])
const loading = ref(false)

 onMounted(() => {
   fetchItems()
 })
 
async function fetchItems() {
  loading.value = true
  try {
     const res = await fetchHeyueList()
    if (res.status === 1) {
      list.value = (res.data || []) as HeyueItem[]
    }
   } catch { /* ignore */ } finally {
     loading.value = false
   }
 }
 
//  function handleAdd() {
//    router.push('/heyue-add')
//  }
 
 function handleEdit(row: HeyueItem) {
   router.push('/heyue-add?id=' + row.id)
 }
 
 async function handleDelete(row: HeyueItem) {
   try {
     await ElMessageBox.confirm(`确定要删除合约「${row.symbol}」吗？`, '提示', {
       confirmButtonText: '确定',
       cancelButtonText: '取消',
       type: 'warning',
     })
   } catch {
     return
   }
   try {
     const res = await deleteHeyue(row.id)
     if (res.status === 1) {
       fetchItems()
     } else {
       ElMessage({ message: String(res.error) || '删除失败', type: 'error' })
     }
   } catch {
     ElMessage({ message: '网络错误', type: 'error' })
   }
 }
 
 function sideLabel(side: number): string {
   return side === 1 ? '开多' : '开空'
 }
 
 function sideClass(side: number): string {
   return side === 1 ? 'side-long' : 'side-short'
 }
 
 function riskLabel(risk: number): string {
   return risk === 2 ? '开启' : '关闭'
 }
 
 function statusLabel(status: number): string {
   return status === 1 ? '正常' : '暂停'
 }
 
</script>

<template>
  <div class="page">
     <h2 class="page-title">合约</h2>
 
    <!-- Toolbar -->
    <div class="toolbar">
          <!-- <button class="btn btn-primary" @click="handleAdd">新增合约</button> -->
    </div>
 
     <!-- List -->
     <div v-if="loading" class="loading">加载中...</div>
     <div v-else-if="list.length === 0" class="empty">暂无数据</div>
     <div v-else class="heyue-list">
       <div v-for="item in list" :key="item.id" class="heyue-item">
         <div class="heyue-top">
           <div class="heyue-top-left">
             <span class="heyue-symbol">{{ item.symbol }}</span>
             <span class="heyue-side" :class="sideClass(item.side)">{{ sideLabel(item.side) }}</span>
           </div>
           <div class="heyue-top-right">
             <span class="heyue-status" :class="{ active: item.status === 1 }">{{ statusLabel(item.status) }}</span>
             <!-- <span class="heyue-username">{{ item.username }}</span> -->
           </div>
         </div>
         <div class="heyue-grid">
           <div class="grid-item">
             <span class="grid-label">用户ID</span>
             <span class="grid-value">{{ item.userid }}</span>
           </div>
           <div class="grid-item">
             <span class="grid-label">已加仓</span>
             <span class="grid-value">{{ item.is_num }}/{{ item.num }}</span>
           </div>
           <div class="grid-item">
             <span class="grid-label">收益%</span>
             <span class="grid-value profit" :class="item.sellprice >= 0 ? 'profit-up' : 'profit-down'">{{ item.sellprice }}%</span>
           </div>
           <div class="grid-item">
             <span class="grid-label">首仓</span>
             <span class="grid-value">{{ item.oneprice }}</span>
           </div>
           <div class="grid-item">
             <span class="grid-label">补仓</span>
             <span class="grid-value">{{ item.repeatprice }}</span>
           </div>
           <div class="grid-item">
             <span class="grid-label">风控</span>
             <span class="grid-value" :class="{ 'risk-on': item.risk === 2 }">{{ riskLabel(item.risk) }}</span>
           </div>
         </div>
         <div class="heyue-actions">
           <button class="btn-edit" @click="handleEdit(item)">编辑</button>
           <button class="btn-delete" @click="handleDelete(item)">删除</button>
         </div>
       </div>
     </div>
 
 </div>
 </template>

 <style scoped>
 .page { padding: 12px 12px 80px; }
 .page-title { font-size: 18px; font-weight: 600; margin: 0 0 12px; color: #1a1a1a;text-align: center; }
 
 /* Toolbar */
 .toolbar {
   display: flex;
   flex-wrap: wrap;
   gap: 6px;
   margin-bottom: 12px;
 }
.btn {
   display: inline-flex;
   align-items: center;
   justify-content: center;
   height: 34px;
   padding: 0 14px;
   border: none;
   border-radius: 6px;
   font-size: 13px;
   cursor: pointer;
   white-space: nowrap;
 }
 .btn:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-primary { background: #409eff; color: #fff; }
 .btn-cancel { background: #f5f5f5; color: #666; }
 .loading, .empty { text-align: center; color: #999; padding: 40px 0; font-size: 14px; }
 
 /* List */
 .heyue-list { display: flex; flex-direction: column; gap: 10px; }
 .heyue-item {
   padding: 12px;
   background: #fff;
   border-radius: 8px;
   box-shadow: 0 1px 4px rgba(0,0,0,0.06);
 }
 .heyue-top {
   display: flex;
   align-items: center;
   justify-content: space-between;
   margin-bottom: 8px;
 }
 .heyue-top-left { display: flex; align-items: center; gap: 6px; }
 .heyue-top-right { display: flex; align-items: center; gap: 8px; }
 .heyue-username { font-size: 12px; color: #999; }
 .heyue-symbol { font-size: 15px; font-weight: 600; color: #1a1a1a; }
 .heyue-side { font-size: 12px; padding: 2px 8px; border-radius: 4px; font-weight: 500; }
 .side-long { color: #27ae60; background: #e8f8f0; }
 .side-short { color: #e74c3c; background: #fdecea; }
 .heyue-status { font-size: 12px; color: #999; }
 .heyue-status.active { color: #27ae60; }
 .heyue-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 6px; }
 .grid-item { display: flex; flex-direction: column; align-items: center; gap: 2px; }
 .grid-label { font-size: 11px; color: #999; }
 .grid-value { font-size: 13px; font-weight: 600; color: #333; }
 .profit-up { color: #27ae60; }
 .profit-down { color: #e74c3c; }
 .risk-on { color: #e6a23c; }
 .heyue-actions {
   display: flex;
   gap: 8px;
   margin-top: 8px;
   padding-top: 8px;
   border-top: 1px solid #f0f0f0;
 }
 .btn-edit, .btn-delete {
   flex: 1;
   height: 30px;
   border: none;
   border-radius: 6px;
   font-size: 13px;
   cursor: pointer;
 }
 .btn-edit { background: #ecf5ff; color: #409eff; }
 .btn-delete { background: #fef0f0; color: #f56c6c; }
 
 </style>



