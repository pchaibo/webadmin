 <script setup lang="ts">
import {  useRoute } from 'vue-router'
import { Coin, Histogram, List, User } from '@element-plus/icons-vue'

const route = useRoute()
function isActive(name: string) {
   return route.name === name
 }
 </script>

 <template>
   <div class="main-layout">
     <div class="main-content">
       <RouterView />
     </div>
     <nav class="bottom-tabs">
       <RouterLink to="/coin" class="tab-btn" :class="{ active: isActive('coin') }">
         <span class="tab-icon"><Histogram /></span>
         <span class="tab-label">币种</span>
       </RouterLink>
      <RouterLink to="/heyue" class="tab-btn" :class="{ active: isActive('heyue') }">
        <span class="tab-icon"><Coin /></span>
        <span class="tab-label">合约</span>
      </RouterLink>
       <RouterLink to="/heyueorder" class="tab-btn" :class="{ active: isActive('heyueorder') }">
         <span class="tab-icon"><List /></span>
         <span class="tab-label">日志</span>
       </RouterLink>
      <RouterLink to="/profile" class="tab-btn" :class="{ active: isActive('profile') }">
         <span class="tab-icon"><User /></span>
         <span class="tab-label">我的</span>
       </RouterLink>
     </nav>
   </div>
 </template>
 
 <style scoped>
.main-layout {
   display: flex; flex-direction: column; min-height: 100vh; background: #f5f5f5;
   padding-bottom: 52px; box-sizing: border-box;
}
.main-content {
   flex: 1;
}
.bottom-tabs {
   position: fixed; left: 0; bottom: 0; width: 100%;
   display: flex; border-top: 1px solid #e0e0e0; background: #fff;
   z-index: 100;
}
.tab-btn {
  text-decoration: none;
  flex: 1; display: flex; flex-direction: column; align-items: center;
  justify-content: center; padding: 6px 0 10px; border: none;
  background: transparent; color: #999; font-size: 11px; cursor: pointer;
  transition: color 0.2s;
}
 .tab-btn.active { color: #4a90d9; }
.tab-icon {
  font-size: 20px; font-weight: 700; line-height: 1.2;
  width: 26px; height: 26px; display: flex; align-items: center;
  justify-content: center; border-radius: 50%; margin-bottom: 2px;
}
 .tab-btn.active .tab-icon { background: #eaf2fb; }
 .tab-label { font-size: 11px; }
 </style>
