<template>
  <div class="coin-page">
    <div class="toolbar">
      <el-input
        v-model="searchSymbol"
        placeholder="搜索合约代码"
        clearable
        style="width: 200px"
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
        <el-option label="停用" :value="0" />
      </el-select>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <el-button type="primary" @click="openCreate">新增Coin</el-button>
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
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column prop="symbol" label="合约代码" min-width="100" />
      <el-table-column prop="close" label="最新价" width="100">
        <template #default="{ row }">
          {{ row.close }}
        </template>
      </el-table-column>
      <el-table-column prop="open" label="开盘价" width="100" />
      <el-table-column prop="high" label="最高价" width="100" />
      <el-table-column prop="low" label="最低价" width="100" />
      <el-table-column prop="quantityprecision" label="价格精度" width="90" />
      <el-table-column prop="status" label="状态" width="70">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '正常' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="addtime" label="添加时间" width="160">
        <template #default="{ row }">
          {{ formatTime(row.addtime) }}
        </template>
      </el-table-column>
      <el-table-column prop="updatetime" label="更新时间" width="160">
        <template #default="{ row }">
          {{ formatTime(row.updatetime) }}
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
      :title="isEditing ? '编辑Coin' : '新增Coin'"
      width="520px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
        @keyup.enter="submitForm"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="合约代码" prop="symbol">
          <el-input v-model="form.symbol" placeholder="请输入合约代码" />
        </el-form-item>
        <el-form-item label="最新价">
          <el-input-number v-model="form.close" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="开盘价">
          <el-input-number v-model="form.open" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="最高价">
          <el-input-number v-model="form.high" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="最低价">
          <el-input-number v-model="form.low" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="价格精度">
          <el-input-number v-model="form.quantityprecision" :min="0" :precision="0" :step="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            active-text="正常"
            inactive-text="停用"
          />
        </el-form-item>
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
  getCoins,
  createCoin,
  updateCoin,
  deleteCoin,
  type CoinItem,
} from '@/api/api'

const loading = ref(false)
const submitting = ref(false)
const items = ref<CoinItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 10
const selectedIds = ref<number[]>([])

const searchSymbol = ref('')
const filterStatus = ref<number | undefined>()

const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  name: '',
  symbol: '',
  close: 0,
  quantityprecision: 0,
  open: 0,
  low: 0,
  high: 0,
  status: 1,
})

const formRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
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
    const res = await getCoins(currentPage.value, searchSymbol.value || undefined, filterStatus.value)
    if (res.status === 1) {
      items.value = res.data || []
      total.value = res.total
    } else {
      ElMessage.error(res.error || '获取Coin列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: CoinItem[]) {
  selectedIds.value = rows.map((r) => r.id)
}

function handleSearch() {
  currentPage.value = 1
  fetchItems()
}

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.name = ''
  form.symbol = ''
  form.close = 0
  form.quantityprecision = 1
  form.open = 0
  form.low = 0
  form.high = 0
  form.status = 1
  dialogVisible.value = true
}

function openEdit(row: CoinItem) {
  isEditing.value = true
  editingId.value = row.id
  form.name = row.name
  form.symbol = row.symbol
  form.close = row.close
  form.quantityprecision = row.quantityprecision
  form.open = row.open
  form.low = row.low
  form.high = row.high
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
      const res = await updateCoin(editingId.value, {
        name: form.name,
        symbol: form.symbol,
        close: form.close,
        quantityprecision: form.quantityprecision,
        open: form.open,
        low: form.low,
        high: form.high,
        status: form.status,
      })
      if (res.status === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchItems()
      } else {
        ElMessage.error(res.error || '更新失败')
      }
    } else {
      const res = await createCoin({
        name: form.name,
        symbol: form.symbol,
        close: form.close,
        quantityprecision: form.quantityprecision,
        open: form.open,
        low: form.low,
        high: form.high,
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

async function handleDelete(row: CoinItem) {
  try {
    await ElMessageBox.confirm(`确定要删除Coin「${row.name}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteCoin(row.id)
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
      `确定要删除选中的 ${selectedIds.value.length} 个Coin吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteCoin(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 个Coin`)
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
.coin-page {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.toolbar {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
