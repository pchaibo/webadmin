<template>
  <div class="shell-add-page">
    <div class="page-header">
      <el-button text @click="goBack">
        &lt; 返回列表
      </el-button>
      <h2 class="page-title">新增Shell</h2>
    </div>

    <div class="form-card">
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

        <el-form-item>
          <el-button type="primary" size="large" :loading="submitting" @click="submitForm">
            {{ submitting ? '创建中...' : '创 建' }}
          </el-button>
          <el-button size="large" @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  createShell,
  getShellGroups,
  type ShellGroupItem,
} from '@/api/api'

const router = useRouter()
const submitting = ref(false)
const groupOptions = ref<ShellGroupItem[]>([])

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

function goBack() {
  router.push('/shell')
}

onMounted(() => {
  fetchGroupOptions()
})

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

async function submitForm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const res = await createShell({
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
      ElMessage.success('创建成功')
      router.push('/shell')
    } else {
      ElMessage.error(res.error || '创建失败')
    }
  } catch {
    ElMessage.error('网络错误')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.shell-add-page {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.page-header {
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.form-card {
  max-width: 700px;
}
</style>
