<template>
  <div class="heyue-page">
    <div class="toolbar">
      <el-input
        v-model="searchSymbol"
        placeholder="搜索合约代码"
        clearable
        style="width: 160px"
        @keyup.enter="handleSearch"
        @clear="handleSearch"
      />
      <el-input
        v-model="searchUsername"
        placeholder="搜索用户名"
        clearable
        style="width: 160px"
        @keyup.enter="handleSearch"
        @clear="handleSearch"
      />
      <el-select
        v-model="filterStatus"
        placeholder="状态筛选"
        clearable
        style="width: 130px"
        @change="handleSearch"
      >
        <el-option label="正常" :value="1" />
        <el-option label="暂停" :value="0" />
      </el-select>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <el-button type="primary" @click="openCreate">新增合约</el-button>
      <el-button type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
        批量删除
      </el-button>
    </div>

    <el-table
      :data="items"
      stripe
      border
      v-loading="loading"
      style="width: 100%"
      @selection-change="onSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="userid" label="用户ID" min-width="50" />
      <el-table-column prop="username" label="用户名" min-width="100" />
      <el-table-column prop="symbol" label="合约代码" width="120" />
      <el-table-column prop="side" label="方向" width="70">
        <template #default="{ row }">
          <span :style="{ color: row.side === 1 ? '#f56c6c' : '#67c23a' }">
            {{ row.side === 1 ? '开多' : '开空' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="num" label="最大次数" width="100" />
      
      <el-table-column prop="sellprice" label="收益%" width="80">
        <template #default="{ row }">
          {{ row.sellprice }}%
        </template>
      </el-table-column>
      <el-table-column prop="oneprice" label="首仓USDT" width="100" />
      <el-table-column prop="repeatprice" label="补仓USDT" width="100" />
      <el-table-column prop="newprice" label="最新价" width="100" />
      <el-table-column prop="is_num" label="已加仓" width="70" />
       <el-table-column prop="risk" label="风控" width="70">
        <template #default="{ row }">
          <span :style="{ color: row.risk === 2 ? '#e6a23c' : '#909399' }">
            {{ row.risk === 2 ? '开启' : '关闭' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="70">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '正常' : '暂停' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="addtime" label="添加时间" min-width="150">
        <template #default="{ row }">
          {{ formatTime(row.addtime) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next, jumper"
        background
        @current-change="fetchItems"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? '编辑合约' : '新增合约'"
      width="620px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="110px"
        @keyup.enter="submitForm"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" placeholder="请输入用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="合约代码" prop="symbol">
              <el-input v-model="form.symbol" placeholder="请输入合约代码" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="方向">
              <el-select v-model="form.side" style="width: 100%">
                <el-option label="开多" :value="1" />
                <el-option label="开空" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
       
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="平仓收益%">
              <el-input-number v-model="form.sellprice" :min="0" :step="0.1" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="首仓USDT">
              <el-input-number v-model="form.oneprice" :min="0" :step="1" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          
         
        </el-row>
        <el-row :gutter="20">
          <!-- <el-col :span="12">
            <el-form-item label="网格类型">
              <el-select v-model="form.rangetype" style="width: 100%">
                <el-option label="差价USDT" :value="1" />
                <el-option label="网格保证金" :value="2" />
              </el-select>
            </el-form-item>
          </el-col> -->
          <!-- <el-col :span="12">
            <el-form-item label="网格差价">
              <el-input-number v-model="form.rangeprice" :min="0" :step="0.1" :precision="4" style="width: 100%" />
            </el-form-item>
          </el-col> -->
        </el-row>
        <el-divider>网格设置</el-divider>
         <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="最大次数">
              <el-input-number v-model="form.num" :min="1" :precision="0" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
         
        </el-row>

          <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="网格平仓">
              <el-select v-model="form.rangeclosing" style="width: 100%">
                <el-option label="不平" :value="1" />
                <el-option label="平仓" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
          <!-- <el-col :span="12">
            <el-form-item label="强平价格">
              <el-input-number v-model="form.closingprice" :min="0" :step="0.001" :precision="6" style="width: 100%" />
            </el-form-item>
          </el-col> -->
          <el-col :span="12">
            <el-form-item label="补仓USDT">
              <el-input-number v-model="form.repeatprice" :min="0" :step="1" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="网格%">
              <el-input-number v-model="form.rangepercent" :min="0" :max="100" :precision="0" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="网格平仓%">
              <el-input-number v-model="form.rangeclosingpct" :min="0" :max="100" :precision="0" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
      

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="风控">
              <el-select v-model="form.risk" style="width: 100%">
                <el-option label="关闭" :value="1" />
                <el-option label="开启" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="风控时间(分)">
              <el-input-number v-model="form.risktime" :min="0" :precision="0" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-divider></el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch
                v-model="form.status"
                :active-value="1"
                :inactive-value="2"
                active-text="正常"
                inactive-text="暂停"
              />
            </el-form-item>
          </el-col>
        </el-row>

      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">
          {{ isEditing ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getHeyues,
  createHeyue,
  updateHeyue,
  deleteHeyue,
  type HeyueItem,
} from '@/api/api'

const loading = ref(false)
const submitting = ref(false)
const items = ref<HeyueItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 10
const selectedIds = ref<number[]>([])

const searchSymbol = ref('')
const searchUsername = ref('')
const filterStatus = ref<number | undefined>()

const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  username: '',
  symbol: '',
  side: 1,
  num: 5,
  status: 1,
  sellprice: 0,
  oneprice: 0,
  repeatprice: 0,
  rangetype: 2,
  rangeprice: 0,
  rangepercent: 100,
  rangeclosingpct: 50,
  rangeclosing: 1,
  closingprice: 0,
  risk: 2,
  risktime: 30,
})

const formRules = {
  symbol: [{ required: true, message: '请输入合约代码', trigger: 'blur' }],
}

const formRef = ref()

function formatTime(ts: number): string {
  if (!ts) return '-'
  const d = new Date(ts * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

onMounted(() => {
  fetchItems()
})

async function fetchItems() {
  loading.value = true
  try {
    const res = await getHeyues(
      currentPage.value,
      searchSymbol.value || undefined,
      searchUsername.value || undefined,
      filterStatus.value
    )
    if (res.status === 1) {
      items.value = res.data || []
      total.value = res.total
    } else {
      ElMessage.error(res.error || '获取合约列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: HeyueItem[]) {
  selectedIds.value = rows.map((r) => r.id)
}

function handleSearch() {
  currentPage.value = 1
  fetchItems()
}

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.username = ''
  form.symbol = ''
  form.side = 1
  form.num = 2
  form.status = 1
  form.sellprice = 30
  form.oneprice = 100
  form.repeatprice = 50
  form.rangetype = 2
  form.rangeprice = 0
  form.rangepercent = 100
  form.rangeclosingpct = 50
  form.rangeclosing = 1
  form.closingprice = 0
  form.risk = 2
  form.risktime = 30
  dialogVisible.value = true
}

function openEdit(row: HeyueItem) {
  isEditing.value = true
  editingId.value = row.id
  form.username = row.username
  form.symbol = row.symbol
  form.side = row.side
  form.num = row.num
  form.status = row.status
  form.sellprice = row.sellprice
  form.oneprice = row.oneprice
  form.repeatprice = row.repeatprice
  form.rangetype = row.rangetype
  form.rangeprice = row.rangeprice
  form.rangepercent = row.rangepercent
  form.rangeclosingpct = row.rangeclosingpct
  form.rangeclosing = row.rangeclosing
  form.closingprice = row.closingprice
  form.risk = row.risk
  form.risktime = row.risktime
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEditing.value && editingId.value !== null) {
      const res = await updateHeyue(editingId.value, {
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
      })
      if (res.status === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchItems()
      } else {
        ElMessage.error(res.error || '更新失败')
      }
    } else {
      const res = await createHeyue({
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
      })
      if (res.status === 1) {
        ElMessage.success('创建成功')
        dialogVisible.value = false
        fetchItems()
      } else {
        ElMessage.error(res.error || '创建失败')
      }
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: HeyueItem) {
  try {
    await ElMessageBox.confirm(`确定要删除合约「${row.symbol}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteHeyue(row.id)
    if (res.status === 1) {
      ElMessage.success('删除成功')
      fetchItems()
    } else {
      ElMessage.error(res.error || '删除失败')
    }
  } catch {
    ElMessage.error('网络错误')
  }
}

async function handleBatchDelete() {
  if (!selectedIds.value.length) return
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedIds.value.length} 个合约吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteHeyue(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 个合约`)
    selectedIds.value = []
    fetchItems()
  } catch {
    ElMessage.error('批量删除时发生网络错误')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.heyue-page {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.toolbar {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
