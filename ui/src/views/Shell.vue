<template>
  <div class="shell-page">
    <div class="toolbar">
      <el-button type="primary" @click="router.push('/shell/add')">新增Shell</el-button>
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
      <el-table-column type="selection" width="45" />
      <el-table-column prop="id" label="ID" width="55" />
      <el-table-column prop="host" label="主机" min-width="130" />
      <el-table-column prop="scheme" label="协议" width="70" />
      <el-table-column label="分组名称" min-width="120">
        <template #default="{ row }">
          {{ row.group?.name || '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="70">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="maxurl" label="大图URL" min-width="140" show-overflow-tooltip />
      <el-table-column label="大图" width="80">
        <template #default="{ row }">
          <el-button size="small" @click="showMaxImages(row)">查看</el-button>
        </template>
      </el-table-column>
      
      <el-table-column prop="minurl" label="小图URL" min-width="140" show-overflow-tooltip />
      <el-table-column prop="dir" label="目录" width="60">
        <template #default="{ row }">
          <el-tag size="small">{{ row.dir }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="lock" label="锁定" width="60">
        <template #default="{ row }">
          <el-tag :type="row.lock === 1 ? 'danger' : 'info'" size="small">
            {{ row.lock === 1 ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip />
      <el-table-column prop="addtime" label="添加时间" width="150">
        <template #default="{ row }">
          {{ formatTime(row.addtime) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
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
      title="编辑Shell"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="90px"
        @keyup.enter="submitForm"
      >
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="主机" prop="host">
              <el-input v-model="form.host" placeholder="请输入主机" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="协议">
              <el-input v-model="form.scheme" placeholder="http/https" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="分组">
              <el-select v-model="form.group_id" placeholder="请选择分组" clearable filterable style="width:100%">
                <el-option v-for="g in groupOptions" :key="g.id" :label="`${g.id} - ${g.name}`" :value="g.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="数量">
              <el-input-number v-model="form.num" :min="0" controls-position="right" style="width:100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="大图URL">
          <el-input v-model="form.maxurl" placeholder="大图 URL" />
        </el-form-item>
        <el-form-item label="小图URL">
          <el-input v-model="form.minurl" placeholder="小图 URL" />
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="目录">
              <el-input-number v-model="form.dir" :min="0" controls-position="right" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-switch
                v-model="form.status"
                :active-value="1"
                :inactive-value="0"
                active-text="启用"
                inactive-text="禁用"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="锁定">
              <el-switch
                v-model="form.lock"
                :active-value="1"
                :inactive-value="0"
                active-text="是"
                inactive-text="否"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="备注">
              <el-input v-model="form.remark" placeholder="备注" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">保存</el-button>
      </template>
    </el-dialog>
    <el-dialog v-model="maxDialogVisible" :title="maxDialogTitle" width="700px" :close-on-click-modal="false">
      <el-table :data="maxImages" stripe border v-loading="maxLoading" style="width:100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="url" label="URL" min-width="350" show-overflow-tooltip />
        <el-table-column prop="addtime" label="添加时间" width="150">
          <template #default="{ row }">
            {{ formatTime(row.addtime) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="70">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getShells,
  updateShell,
  deleteShell,
  getShellMax,
  getShellGroups,
  type ShellItem,
  type ShellMaxItem,
  type ShellGroupItem,
} from '@/api/api'

const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const items = ref<ShellItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = Number(import.meta.env.VITE_PAGE_SIZE) || 20
const selectedIds = ref<number[]>([])
const groupOptions = ref<ShellGroupItem[]>([])

const dialogVisible = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  host: '',
  scheme: '',
  group_id: 0,
  status: 1,
  num: 0,
  maxurl: '',
  minurl: '',
  dir: 0,
  lock: 0,
  remark: '',
})

const formRules = {
  host: [{ required: true, message: '请输入主机', trigger: 'blur' }],
}

const formRef = ref()
const maxDialogVisible = ref(false)
const maxDialogTitle = ref('')
const maxLoading = ref(false)
const maxImages = ref<ShellMaxItem[]>([])

async function showMaxImages(row: ShellItem) {
  maxDialogTitle.value = `大图列表 - ${row.host}`
  maxDialogVisible.value = true
  maxLoading.value = true
  try {
    const res = await getShellMax(row.id)
    if (res.status === 1) {
      maxImages.value = res.data || []
    } else {
      ElMessage.error(res.error || '获取大图列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    maxLoading.value = false
  }
}

function formatTime(ts: number): string {
  if (!ts) return '-'
  const d = new Date(ts * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

async function fetchGroupOptions() {
  try {
    const res = await getShellGroups(1)
    if (res.status === 1) {
      groupOptions.value = res.data || []
    }
  } catch {
    // ignore
  }
}

onMounted(() => {
  fetchItems()
  fetchGroupOptions()
})

async function fetchItems() {
  loading.value = true
  try {
    const res = await getShells(currentPage.value)
    if (res.status === 1) {
      items.value = res.data || []
      total.value = res.total
    } else {
      ElMessage.error(res.error || '获取Shell列表失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    loading.value = false
  }
}

function onSelectionChange(rows: ShellItem[]) {
  selectedIds.value = rows.map((r) => r.id)
}

function openEdit(row: ShellItem) {
  editingId.value = row.id
  form.host = row.host
  form.scheme = row.scheme
  form.group_id = row.group_id
  form.status = row.status
  form.num = row.num
  form.maxurl = row.maxurl
  form.minurl = row.minurl
  form.dir = row.dir
  form.lock = row.lock
  form.remark = row.remark
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const res = await updateShell(editingId.value!, {
      host: form.host,
      scheme: form.scheme,
      group_id: form.group_id,
      status: form.status,
      num: form.num,
      maxurl: form.maxurl,
      minurl: form.minurl,
      dir: form.dir,
      lock: form.lock,
      remark: form.remark,
    })
    if (res.status === 1) {
      ElMessage.success('更新成功')
      dialogVisible.value = false
      fetchItems()
    } else {
      ElMessage.error(res.error || '更新失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: ShellItem) {
  try {
    await ElMessageBox.confirm(`确定要删除Shell「${row.host}」吗？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }

  try {
    const res = await deleteShell(row.id)
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
      `确定要删除选中的 ${selectedIds.value.length} 个Shell吗？`,
      '批量删除',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    const results = await Promise.allSettled(
      selectedIds.value.map((id) => deleteShell(id))
    )
    const failCount = results.filter((r) => r.status === 'rejected').length
    ElMessage.success(`成功删除 ${selectedIds.value.length - failCount} 个Shell`)
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
.shell-page {
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
