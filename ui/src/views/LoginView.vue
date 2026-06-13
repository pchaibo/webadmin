<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h1 class="login-title">后台管理</h1>
        <p class="login-subtitle">请输入您的账号信息登录</p>
      </div>

      <el-alert
        v-if="errorMsg"
        :title="errorMsg"
        type="error"
        show-icon
        closable
        class="login-error"
        @close="errorMsg = ''"
      />

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        size="large"
        @keyup.enter="handleLogin"
      >
        <el-form-item label="用户名 / 邮箱" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名或邮箱"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { login } from '@/api/api'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名或邮箱', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  errorMsg.value = ''

  try {
    const res = await login({
      username: form.username,
      password: form.password,
    })

    if (res.status === 1 && res.token) {
      localStorage.setItem('admin_token', res.token)
      localStorage.setItem('admin_id', String(res.id ?? ''))
      localStorage.setItem('admin_username', res.username ?? '')
      localStorage.setItem('admin_email', res.email ?? '')
      router.push('/')
    } else {
      errorMsg.value = res.error || '登录失败，请检查账号和密码'
    }
  } catch {
    errorMsg.value = '网络错误，无法连接到服务器'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  min-height: 100dvh;
  padding: 24px;
}

.login-card {
  width: 100%;
  max-width: 420px;
  padding: 40px 32px 32px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  backdrop-filter: blur(12px);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0;
  color: #e5eefc;
}

.login-subtitle {
  margin: 8px 0 0;
  font-size: 14px;
  color: rgba(229, 238, 252, 0.5);
}

.login-error {
  margin-bottom: 20px;
}

.login-btn {
  width: 100%;
  margin-top: 4px;
}

.login-btn :deep(.el-button__loading) {
  margin-right: 6px;
}
</style>
