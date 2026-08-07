<template>
  <div class="admingroup-page">
    <div class="toolbar">
      <el-button type="primary" @click="openCreate">新增用户组</el-button>
      <el-button type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
        批量删除
      </el-button>
    </div>

    <el-table
      :data="groups"
      stripe
      border
      v-loading="loading"
      style="width: 100%"
      @selection-change="onSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="title" label="用户组名称" min-width="140" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="rules" label="权限规则" min-width="240" show-overflow-tooltip />
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
        @current-change="fetchGroups"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? '编辑用户组' : '新增用户组'"
      width="520px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="90px"
        @keyup.enter="submitForm"
      >
        <el-form-item label="组名" prop="title">
          <el-input v-model="form.title" placeholder="请输入用户组名称" />
        </el-form-item>
        <el-form-item label="权限规则">
          <el-checkbox-group v-model="form.ruleIds" class="rules-checkbox-group">
            <el-checkbox
              v-for="rule in authRules"
              :key="rule.id"
              :value="rule.id"
            >
              {{ rule.title }}
            </el-checkbox>
          </el-checkbox-group>
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
  getAuthGroups,
  createAuthGroup,
  updateAuthGroup,
  deleteAuthGroup,
  getAuthRulesAll,
  type AuthGroupItem,
  type AuthRuleItem,
} from '@/api/api'

const loading = ref(false)
const submitting = ref(false)
const groups = ref<AuthGroupItem[]>([])
const authRules = ref<AuthRuleItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 20
const selectedIds = ref<number[]>([])

const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  title: '',
  status: 1,
  ruleIds: [] as number[],
})

const formRules = {
  title: [{ required: true, message: '请输入用户组名称', trigger: 'blur' }],
}

const formRef = ref()

onMounted(() => {
  fetchGroups()
})

async function fetchAuthRules() {
  try {
    const res = await getAuthRulesAll()
    if (res.status === 1) {
      authRules.value = res.data || []
    } else {
      ElMessage.error(res.error || '获取权限规则失败')
    }
  } catch {
    ElMessage.error('网络错误')
  }
}

async function fetchGroups() {
  loading.value = true
  try {
    const res = await getAuthGroups(currentPage.value)
    if (res.status === 1) {
      groups.value = res.data || []
      total.value = res.total
    } else {
      ElMessage.error(res.error || '获取用户组列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: AuthGroupItem[]) {
  selectedIds.value = rows.map((r) => r.id)
}

async function openCreate() {
  await fetchAuthRules()
  isEditing.value = false
  editingId.value = null
  form.title = ''
  form.status = 1
  form.ruleIds = []
  dialogVisible.value = true
}

async function openEdit(row: AuthGroupItem) {
  await fetchAuthRules()
  isEditing.value = true
  editingId.value = row.id
  form.title = row.title
  form.status = row.status
  form.ruleIds = Array.from(
    new Set(
      String(row.rules || '')
        .split(',')
        .map((id) => Number(id))
        .filter((id) => id > 0)
    )
  )
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEditing.value && editingId.value !== null) {
      const res = await updateAuthGroup(editingId.value, {
        title: form.title,
        status: form.status,
        rules: form.ruleIds.join(','),
      })
      if (res.status === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchGroups()
      } else {
        ElMessage.error(res.error || '更新失败')
      }
    } else {
      const res = await createAuthGroup({
        title: form.title,
        status: form.status,
        rules: form.ruleIds.join(','),
      })
      if (res.status === 1) {
        ElMessage.success('创建成功')
        dialogVisible.value = false
        fetchGroups()
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

async function handleDelete(row: AuthGroupItem) {
  try {
    await ElMessageBox.confirm(`确定要删除用户组「${row.title}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteAuthGroup(row.id)
    if (res.status === 1) {
      ElMessage.success('删除成功')
      fetchGroups()
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
      `确定要删除选中的 ${selectedIds.value.length} 个用户组吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteAuthGroup(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 个用户组`)
    selectedIds.value = []
    fetchGroups()
  } catch {
    ElMessage.error('批量删除时发生网络错误')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admingroup-page {
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

.rules-checkbox-group {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 4px 16px;
  width: 100%;
  max-height: 220px;
  overflow-y: auto;
  padding: 4px 0;
}
</style>
