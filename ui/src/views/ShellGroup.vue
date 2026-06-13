<template>
  <div class="shellgroup-page">
    <div class="toolbar">
      <el-button type="primary" @click="openCreate">新增Shell分组</el-button>
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
      <el-table-column prop="name" label="分组名称" min-width="150" />
      <el-table-column prop="mmurl" label="域名" min-width="180" />
      <el-table-column prop="addtime" label="添加时间" width="160">
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
      :title="isEditing ? '编辑Shell分组' : '新增Shell分组'"
      width="480px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="90px"
        @keyup.enter="submitForm"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分组名称" />
        </el-form-item>
        <el-form-item label="域名">
          <el-input v-model="form.mmurl" placeholder="域名或 URL" />
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
  getShellGroups,
  createShellGroup,
  updateShellGroup,
  deleteShellGroup,
  type ShellGroupItem,
} from '@/api/api'

const loading = ref(false)
const submitting = ref(false)
const items = ref<ShellGroupItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 20
const selectedIds = ref<number[]>([])

const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  name: '',
  mmurl: '',
})

const formRules = {
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }],
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
    const res = await getShellGroups(currentPage.value)
    if (res.status === 1) {
      items.value = res.data || []
      total.value = res.total
    } else {
      ElMessage.error(res.error || '获取Shell分组列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: ShellGroupItem[]) {
  selectedIds.value = rows.map((r) => r.id)
}

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.name = ''
  form.mmurl = ''
  dialogVisible.value = true
}

function openEdit(row: ShellGroupItem) {
  isEditing.value = true
  editingId.value = row.id
  form.name = row.name
  form.mmurl = row.mmurl
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEditing.value && editingId.value !== null) {
      const res = await updateShellGroup(editingId.value, {
        name: form.name,
        mmurl: form.mmurl,
      })
      if (res.status === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchItems()
      } else {
        ElMessage.error(res.error || '更新失败')
      }
    } else {
      const res = await createShellGroup({
        name: form.name,
        mmurl: form.mmurl,
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

async function handleDelete(row: ShellGroupItem) {
  try {
    await ElMessageBox.confirm(`确定要删除Shell分组「${row.name}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteShellGroup(row.id)
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
      `确定要删除选中的 ${selectedIds.value.length} 个Shell分组吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteShellGroup(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 个分组`)
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
.shellgroup-page {
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
