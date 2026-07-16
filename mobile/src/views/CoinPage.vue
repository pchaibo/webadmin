<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
 import { useRouter } from 'vue-router'
import { fetchCoinList } from '../api'
 
 const router = useRouter()
 
 interface CoinItem {
   id: number
   name: string
   symbol: string
   close: number
   high: number
   low: number
   status: number
   [key: string]: unknown
 }
 
 const list = ref<CoinItem[]>([])
 const loading = ref(true)
 
 // ── WebSocket ──
 let ws: WebSocket | null = null
 let wsReconnectTimer: ReturnType<typeof setTimeout> | null = null
 
 function connectWebSocket() {
   if (ws) { ws.close(); ws = null }
 
   //const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
   //const wsUrl = `${protocol}//${location.host}/user/ws`
   const wsUrl = `ws://127.0.0.1:4000/user/ws`
 
   ws = new WebSocket(wsUrl)
 
   ws.onopen = () => {
     console.log('[WS] coin price connected')
   }
 
   ws.onmessage = (event: MessageEvent) => {
     try {
       const raw = JSON.parse(event.data)
       if (raw.status !== 1 || !Array.isArray(raw.data)) return
 
       for (const itemStr of raw.data) {
         let priceData: any
         try {
           priceData = JSON.parse(itemStr)
         } catch {
           continue
         }
 
         const symbol = (priceData.Symbol || '').toLowerCase()
         if (!symbol) continue
 
         const row = list.value.find((r) => r.symbol.toLowerCase() === symbol)
         if (row) {
           row._priceDirection = ''
           if (priceData.Change > 0) {
             row._priceDirection = 'up'
           } else if (priceData.Change < 0) {
             row._priceDirection = 'down'
           }
           if (priceData.C !== undefined) row.close = priceData.C
           if (priceData.H !== undefined) row.high = priceData.H
           if (priceData.I !== undefined) row.low = priceData.I
           row._priceUpdated = true
           setTimeout(() => { row._priceUpdated = false }, 600)
         }
       }
     } catch { /* ignore */ }
   }
 
   ws.onclose = () => {
     console.log('[WS] disconnected, reconnecting in 3s')
     wsReconnectTimer = setTimeout(() => connectWebSocket(), 3000)
   }
 
   ws.onerror = () => { ws?.close() }
 }
 
 onMounted(async () => {
   try {
     const res = await fetchCoinList()
     if (res.status === 1 && Array.isArray(res.data)) {
       list.value = res.data as CoinItem[]
     }
   } catch { /* ignore */ } finally {
     loading.value = false
   }
   connectWebSocket()
 })
 
 onUnmounted(() => {
   if (ws) { ws.close(); ws = null }
   if (wsReconnectTimer) { clearTimeout(wsReconnectTimer); wsReconnectTimer = null }
 })
 </script>
 
 <template>
   <div class="page">
     <h2 class="page-title">币种列表</h2>
     <div v-if="loading" class="loading">加载中...</div>
     <div v-else-if="list.length === 0" class="empty">暂无数据</div>
     <div v-else class="coin-list">
      <div v-for="item in list" :key="item.id" class="coin-item">
         <div class="coin-row">
         <div class="coin-header">
           <span class="coin-name">{{ item.name }}</span>
           <span class="coin-symbol">{{ item.symbol }}</span>
         </div>
           <button class="btn-set" @click="router.push('/heyue-add?symbol=' + item.symbol)">新增合约</button>
         </div>
        <div class="coin-prices">
           <div class="price-row">
             <span class="price-label">最新价</span>
             <span class="price-value close"
               :class="{
                 'price-up': item._priceDirection === 'up',
                 'price-down': item._priceDirection === 'down',
               }">{{ item.close }}</span>
           </div>
           <div class="price-row">
             <span class="price-label">最高</span>
             <span class="price-value high">{{ item.high }}</span>
           </div>
           <div class="price-row">
             <span class="price-label">最低</span>
             <span class="price-value low">{{ item.low }}</span>
           </div>
         </div>
       </div>
     </div>
   </div>
 </template>
 
 <style scoped>
 .page { padding: 16px 16px 80px; }
 .page-title { font-size: 18px; font-weight: 600; margin: 0 0 16px; color: #1a1a1a; text-align: center;}
 .loading, .empty { text-align: center; color: #999; padding: 40px 0; font-size: 14px; }
 .coin-list { display: flex; flex-direction: column; gap: 10px; }
 .coin-item {
   display: flex; flex-direction: column; gap: 8px;
   padding: 14px; background: #fff; border-radius: 8px; box-shadow: 0 1px 4px rgba(0,0,0,0.06);
 }
.coin-header { display: flex; align-items: center; gap: 6px; }
.coin-name { font-size: 15px; font-weight: 500; color: #1a1a1a; }
.coin-symbol { font-size: 12px; color: #999; background: #f0f0f0; padding: 2px 8px; border-radius: 4px; }
 .coin-row { display: flex; align-items: center; justify-content: space-between; }
 .btn-set {
   height: 28px;
   padding: 0 12px;
   border: none;
   border-radius: 6px;
   background: #409eff;
   color: #fff;
   font-size: 12px;
   cursor: pointer;
   white-space: nowrap;
 }
.coin-prices { display: flex; gap: 12px; }
 .price-row { display: flex; flex-direction: column; align-items: center; flex: 1; }
 .price-label { font-size: 11px; color: #999; margin-bottom: 2px; }
 .price-value { font-size: 14px; font-weight: 600; transition: color 0.3s; }
 .price-value.close { color: #e67e22; }
 .price-value.price-up { color: #27ae60; }
 .price-value.price-down { color: #e74c3c; }
 .price-value.high { color: #27ae60; }
 .price-value.low { color: #e74c3c; }
 </style>
