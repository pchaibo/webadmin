<template>
  <div class="heyueorder-page">
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
      <el-select
        v-model="filterOrdertype"
        placeholder="交易类型"
        clearable
        style="width: 130px"
        @change="handleSearch"
      >
        <el-option label="开仓" :value="1" />
        <el-option label="平仓" :value="2" />
      </el-select>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <!-- <el-button type="primary" @click="openCreate">新增订单22</el-button>
      <el-button type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
        批量删除
      </el-button> -->
    </div>

    <el-table
      :data="items"
      stripe
      border
      v-loading="loading"
      style="width: 100%"
      @selection-change="onSelectionChange"
    >
      <!-- <el-table-column type="selection" width="55" /> -->
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="username" label="用户名" min-width="100" />
      <el-table-column prop="symbol" label="合约代码" width="120" />
      <el-table-column prop="ordertype" label="交易类型" width="100">
        <template #default="{ row }">
          <el-tag :type="row.ordertype === 1 ? 'primary' : 'warning'" size="small">
            {{ row.ordertype === 1 ? '开仓' : '平仓' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="side" label="方向" width="70">
        <template #default="{ row }">
          <span :style="{ color: row.side === 1 ? '#f56c6c' : '#67c23a' }">
            {{ row.side === 1 ? '开多' : '开空' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="price" label="价格" width="100">
        <template #default="{ row }">
          {{ row.price.toFixed(4) }}
        </template>
      </el-table-column>
      <el-table-column prop="quantity" label="数量" width="100" />
      <el-table-column prop="total" label="总金额" width="100" />
      <el-table-column prop="num" label="第几次" width="70" />
      <el-table-column prop="usdt" label="收益" width="100">
        <template #default="{ row }">
          <span :style="{ color: row.usdt >= 0 ? '#67c23a' : '#f56c6c' }">
            {{ row.usdt.toFixed(2) }}
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
          <el-button size="small" @click="openEdit(row)">查看</el-button>
          <!-- <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button> -->
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-wrap">
      <div class="stats-info">
        统计收益: <span :class="totalUsdt >= 0 ? 'profit' : 'loss'">{{ totalUsdt.toFixed(2) }} USDT</span>
      </div>
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
      :title="isEditing ? '编辑订单' : '新增订单'"
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
            <el-form-item label="交易类型">
              <el-select v-model="form.ordertype" style="width: 100%">
                <el-option label="开仓" :value="1" />
                <el-option label="平仓" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
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
            <el-form-item label="价格">
              <el-input-number v-model="form.price" :min="0" :step="0.01" :precision="4" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="数量">
              <el-input-number v-model="form.quantity" :min="0" :step="0.001" :precision="4" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="总金额">
              <el-input-number v-model="form.total" :min="0" :step="1" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="第几次">
              <el-input-number v-model="form.num" :min="1" :precision="0" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="收益">
              <el-input-number v-model="form.usdt" :step="0.1" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch
                v-model="form.status"
                :active-value="1"
                :inactive-value="0"
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
  getHeyuesorders,
  createHeyuesorder,
  updateHeyuesorder,
  deleteHeyuesorder,
  type HeyuesorderItem,
} from '@/api/api'

const loading = ref(false)
const submitting = ref(false)
const items = ref<HeyuesorderItem[]>([])
const total = ref(0)
const totalUsdt = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 10
const selectedIds = ref<number[]>([])

const searchSymbol = ref('')
const searchUsername = ref('')
const filterStatus = ref<number | undefined>()
const filterOrdertype = ref<number | undefined>()

const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  username: '',
  symbol: '',
  ordertype: 1,
  side: 1,
  price: 0,
  quantity: 0,
  total: 0,
  num: 1,
  usdt: 0,
  status: 1,
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
    const res = await getHeyuesorders(
      currentPage.value,
      searchSymbol.value || undefined,
      searchUsername.value || undefined,
      filterStatus.value,
      filterOrdertype.value
    )
    if (res.status === 1) {
      items.value = res.data || []
      total.value = res.total
      totalUsdt.value = res.total_usdt || 0
    } else {
      ElMessage.error(res.error || '获取订单列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: HeyuesorderItem[]) {
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
  form.ordertype = 1
  form.side = 1
  form.price = 0
  form.quantity = 0
  form.total = 0
  form.num = 1
  form.usdt = 0
  form.status = 1
  dialogVisible.value = true
}

function openEdit(row: HeyuesorderItem) {
  isEditing.value = true
  editingId.value = row.id
  form.username = row.username
  form.symbol = row.symbol
  form.ordertype = row.ordertype
  form.side = row.side
  form.price = row.price
  form.quantity = row.quantity
  form.total = row.total
  form.num = row.num
  form.usdt = row.usdt
  form.status = row.status
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEditing.value && editingId.value !== null) {
      const res = await updateHeyuesorder(editingId.value, {
        ordertype: form.ordertype,
        side: form.side,
        price: form.price,
        total: form.total,
        quantity: form.quantity,
        num: form.num,
        status: form.status,
        usdt: form.usdt,
      })
      if (res.status === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchItems()
      } else {
        ElMessage.error(res.error || '更新失败')
      }
    } else {
      const res = await createHeyuesorder({
        username: form.username,
        symbol: form.symbol,
        ordertype: form.ordertype,
        side: form.side,
        price: form.price,
        total: form.total,
        quantity: form.quantity,
        num: form.num,
        usdt: form.usdt,
        status: form.status,
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

async function handleDelete(row: HeyuesorderItem) {
  try {
    await ElMessageBox.confirm(`确定要删除订单「${row.symbol}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteHeyuesorder(row.id)
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
      `确定要删除选中的 ${selectedIds.value.length} 个订单吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteHeyuesorder(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 个订单`)
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
.heyueorder-page {
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
  justify-content: space-between;
  align-items: center;
}

.stats-info {
  font-size: 14px;
  color: #606266;
}

.stats-info .profit {
  color: #67c23a;
  font-weight: 600;
}

.stats-info .loss {
  color: #f56c6c;
  font-weight: 600;
}
</style>
