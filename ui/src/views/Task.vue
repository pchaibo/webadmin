<template>
  <div class="task-page">
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
        v-model="filterCondition"
        placeholder="条件筛选"
        clearable
        style="width: 130px"
        @change="handleSearch"
      >
        <el-option label="小于" :value="1" />
        <el-option label="大于" :value="2" />
      </el-select>
      <el-select
        v-model="filterStatus"
        placeholder="状态筛选"
        clearable
        style="width: 130px"
        @change="handleSearch"
      >
        <el-option label="正常" :value="1" />
        <el-option label="停用" :value="2" />
      </el-select>
      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <el-button type="primary" @click="openCreate">新增任务</el-button>
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
      <el-table-column prop="coinid" label="币种ID" width="80" />
      <el-table-column prop="symbol" label="合约代码" width="120" />
      <el-table-column prop="userid" label="用户ID" width="80" />
      <el-table-column prop="username" label="用户名" min-width="100" />
      <el-table-column prop="price" label="价格" min-width="100" />
      <el-table-column prop="condition" label="条件" width="70">
        <template #default="{ row }">
          <span :style="{ color: row.condition === 2 ? '#f56c6c' : '#67c23a' }">
            {{ row.condition === 1 ? '小于' : '大于' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="70">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '正常' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="addtime" label="添加时间" min-width="150">
        <template #default="{ row }">
          {{ formatTime(row.addtime) }}
        </template>
      </el-table-column>
      <el-table-column prop="updatetime" label="更新时间" min-width="150">
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
      :title="isEditing ? '编辑任务' : '新增任务'"
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
            <el-form-item label="币种ID">
              <el-input-number v-model="form.coinid" :min="0" :precision="0" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户ID">
              <el-input-number v-model="form.userid" :min="0" :precision="0" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="价格">
              <el-input-number v-model="form.price" :min="0" :step="0.01" :precision="6" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="条件">
              <el-select v-model="form.condition" style="width: 100%">
                <el-option label="小于" :value="1" />
                <el-option label="大于" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch
                v-model="form.status"
                :active-value="1"
                :inactive-value="2"
                active-text="正常"
                inactive-text="停用"
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
  getTasks,
  createTask,
  updateTask,
  deleteTask,
  type TaskItem,
} from '@/api/api'

const loading = ref(false)
const submitting = ref(false)
const items = ref<TaskItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 10
const selectedIds = ref<number[]>([])

const searchSymbol = ref('')
const searchUsername = ref('')
const filterCondition = ref<number | undefined>()
const filterStatus = ref<number | undefined>()

const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  username: '',
  symbol: '',
  coinid: 0,
  userid: 0,
  price: 0,
  condition: 1,
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
    const res = await getTasks(
      currentPage.value,
      searchSymbol.value || undefined,
      searchUsername.value || undefined,
      filterStatus.value,
      filterCondition.value
    )
    if (res.data) {
      items.value = res.data
      total.value = res.total
    } else {
      ElMessage.error(res.error || '获取任务列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: TaskItem[]) {
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
  form.coinid = 0
  form.userid = 0
  form.price = 0
  form.condition = 1
  form.status = 1
  dialogVisible.value = true
}

function openEdit(row: TaskItem) {
  isEditing.value = true
  editingId.value = row.id
  form.username = row.username
  form.symbol = row.symbol
  form.coinid = row.coinid
  form.userid = row.userid
  form.price = row.price
  form.condition = row.condition
  form.status = row.status
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const payload = {
      username: form.username,
      symbol: form.symbol,
      coinid: form.coinid,
      userid: form.userid,
      price: form.price,
      condition: form.condition,
      status: form.status,
    }
    if (isEditing.value && editingId.value !== null) {
      const res = await updateTask(editingId.value, payload)
      if (res.status === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchItems()
      } else {
        ElMessage.error(res.error || '更新失败')
      }
    } else {
      const res = await createTask(payload)
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

async function handleDelete(row: TaskItem) {
  try {
    await ElMessageBox.confirm(`确定要删除任务「${row.symbol}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteTask(row.id)
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
      `确定要删除选中的 ${selectedIds.value.length} 个任务吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteTask(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 个任务`)
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
.task-page {
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
