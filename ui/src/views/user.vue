<template>
  <div class="user-page">
    <div class="toolbar">
      <el-button type="primary" @click="openCreate">新增用户</el-button>
      <el-button type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
        批量删除
      </el-button>
    </div>

    <el-table
      :data="users"
      stripe
      border
      v-loading="loading"
      style="width: 100%"
      @selection-change="onSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="username" label="用户名" min-width="120" />
      <el-table-column prop="email" label="邮箱" min-width="120" />
      <el-table-column prop="usdt" label="USDT" width="100" />
      <el-table-column prop="margin" label="保证金%" width="100" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="注册时间" min-width="160">
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
        @current-change="fetchUsers"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? '编辑用户' : '新增用户'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="110px"
        @keyup.enter="submitForm"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            :placeholder="isEditing ? '留空则不修改密码' : '请输入密码'"
          />
        </el-form-item>
        <el-form-item label="USDT">
          <el-input-number v-model="form.usdt" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="保证金%">
          <el-input-number v-model="form.margin" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="BN Access" prop="bnaccess">
          <el-input
            v-model="form.bnaccess"
            type="textarea"
            :rows="3"
            placeholder="请输入 BN Access，多行内容"
          />
        </el-form-item>
        <el-form-item label="BN Secret" prop="bnasecret">
          <el-input
            v-model="form.bnasecret"
            type="textarea"
            :rows="3"
            placeholder="请输入 BN Secret，多行内容"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="禁用"
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
  getUsers,
  createUser,
  updateUser,
  deleteUser,
  type UserItem,
} from '@/api/api'

const loading = ref(false)
const submitting = ref(false)
const users = ref<UserItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 20
const selectedIds = ref<number[]>([])

const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  username: '',
  email: '',
  password: '',
  usdt: 0,
  margin: 0,
  status: 1,
  bnaccess: '',
  bnasecret: '',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
}

const formRef = ref()

onMounted(() => {
  fetchUsers()
})

function formatTime(ts: number): string {
  if (!ts) return ''
  const d = new Date(ts * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await getUsers(currentPage.value)
    if (res.status === 1) {
      users.value = res.data || []
      total.value = res.total
    } else {
      ElMessage.error(res.error || '获取用户列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: UserItem[]) {
  selectedIds.value = rows.map((r) => r.id)
}

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.username = ''
  form.email = ''
  form.password = ''
  form.usdt = 0
  form.margin = 0
  form.status = 1
  form.bnaccess = ''
  form.bnasecret = ''
  dialogVisible.value = true
}

function openEdit(row: UserItem) {
  isEditing.value = true
  editingId.value = row.id
  form.username = row.username
  form.email = row.email
  form.password = ''
  form.usdt = row.usdt
  form.margin = row.margin
  form.status = row.status
  form.bnaccess = row.bnaccess
  form.bnasecret = row.bnasecret
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEditing.value && editingId.value !== null) {
      const payload: Record<string, unknown> = {
        username: form.username,
        email: form.email,
        usdt: form.usdt,
        margin: form.margin,
        status: form.status,
        bnaccess: form.bnaccess,
        bnasecret: form.bnasecret,
      }
      if (form.password) {
        payload.password = form.password
      }
      const res = await updateUser(editingId.value, payload)
      if (res.status === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchUsers()
      } else {
        ElMessage.error(res.error || '更新失败')
      }
    } else {
      const res = await createUser({
        username: form.username,
        email: form.email,
        password: form.password,
        usdt: form.usdt,
        margin: form.margin,
        status: form.status,
        bnaccess: form.bnaccess,
        bnasecret: form.bnasecret,
      })
      if (res.status === 1) {
        ElMessage.success('创建成功')
        dialogVisible.value = false
        fetchUsers()
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

async function handleDelete(row: UserItem) {
  try {
    await ElMessageBox.confirm(`确定要删除用户「${row.username}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteUser(row.id)
    if (res.status === 1) {
      ElMessage.success('删除成功')
      fetchUsers()
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
      `确定要删除选中的 ${selectedIds.value.length} 个用户吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteUser(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 个用户`)
    selectedIds.value = []
    fetchUsers()
  } catch {
    ElMessage.error('批量删除时发生网络错误')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.user-page {
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
