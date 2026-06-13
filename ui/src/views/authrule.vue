<template>
  <div class="authrule-page">
    <div class="toolbar">
      <el-button type="primary" @click="openCreate">新增规则</el-button>
      <el-button type="danger" :disabled="!selectedIds.length" @click="handleBatchDelete">
        批量删除
      </el-button>
    </div>

    <el-table
      :data="rules"
      stripe
      border
      v-loading="loading"
      style="width: 100%"
      @selection-change="onSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="pid" label="PID" width="60" />
      <el-table-column prop="name" label="路由名称" min-width="130" />
      <el-table-column prop="title" label="标题" min-width="130" />
      <el-table-column prop="icon" label="图标" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.icon" size="small">{{ row.icon }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="type" label="类型" width="70">
        <template #default="{ row }">
          <el-tag size="small">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="condition" label="条件" min-width="120" show-overflow-tooltip />
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
        @current-change="fetchRules"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? '编辑规则' : '新增规则'"
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
        <el-form-item label="PID">
          <el-input-number v-model="form.pid" :min="0" controls-position="right" />
        </el-form-item>
        <el-form-item label="规则名" prop="name">
          <el-input v-model="form.name" placeholder="请输入规则路由名称" />
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" placeholder="图标 className" />
        </el-form-item>
        <el-form-item label="类型">
          <el-input-number v-model="form.type" :min="0" controls-position="right" />
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
        <el-form-item label="条件">
          <el-input v-model="form.condition" placeholder="条件表达式" />
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
  getAuthRules,
  createAuthRule,
  updateAuthRule,
  deleteAuthRule,
  type AuthRuleItem,
} from '@/api/api'

const loading = ref(false)
const submitting = ref(false)
const rules = ref<AuthRuleItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 20
const selectedIds = ref<number[]>([])

const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  pid: 0,
  name: '',
  title: '',
  icon: '',
  type: 0,
  status: 1,
  condition: '',
})

const formRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
}

const formRef = ref()

onMounted(() => {
  fetchRules()
})

async function fetchRules() {
  loading.value = true
  try {
    const res = await getAuthRules(currentPage.value)
    if (res.status === 1) {
      rules.value = res.data || []
      total.value = res.total
    } else {
      ElMessage.error(res.error || '获取权限规则列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: AuthRuleItem[]) {
  selectedIds.value = rows.map((r) => r.id)
}

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.pid = 0
  form.name = ''
  form.title = ''
  form.icon = ''
  form.type = 0
  form.status = 1
  form.condition = ''
  dialogVisible.value = true
}

function openEdit(row: AuthRuleItem) {
  isEditing.value = true
  editingId.value = row.id
  form.pid = row.pid
  form.name = row.name
  form.title = row.title
  form.icon = row.icon
  form.type = row.type
  form.status = row.status
  form.condition = row.condition
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEditing.value && editingId.value !== null) {
      const res = await updateAuthRule(editingId.value, {
        pid: form.pid,
        name: form.name,
        title: form.title,
        icon: form.icon,
        type: form.type,
        status: form.status,
        condition: form.condition,
      })
      if (res.status === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchRules()
      } else {
        ElMessage.error(res.error || '更新失败')
      }
    } else {
      const res = await createAuthRule({
        pid: form.pid,
        name: form.name,
        title: form.title,
        icon: form.icon,
        type: form.type,
        status: form.status,
        condition: form.condition,
      })
      if (res.status === 1) {
        ElMessage.success('创建成功')
        dialogVisible.value = false
        fetchRules()
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

async function handleDelete(row: AuthRuleItem) {
  try {
    await ElMessageBox.confirm(`确定要删除规则「${row.title}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteAuthRule(row.id)
    if (res.status === 1) {
      ElMessage.success('删除成功')
      fetchRules()
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
      `确定要删除选中的 ${selectedIds.value.length} 条规则吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteAuthRule(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 条规则`)
    selectedIds.value = []
    fetchRules()
  } catch {
    ElMessage.error('批量删除时发生网络错误')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.authrule-page {
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
