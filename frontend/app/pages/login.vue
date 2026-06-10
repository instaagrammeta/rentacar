<template>
  <div class="w-full max-w-md">
    <div class="mb-6 flex justify-center">
      <AppLogo />
    </div>
    <div class="card p-7">
      <h1 class="text-center text-xl font-bold text-ink">Вход в систему</h1>
      <p class="mt-1 text-center text-sm text-ink-muted">Управление прокатом автомобилей</p>

      <form class="mt-6 space-y-4" @submit.prevent="submit">
        <FormField v-model="username" label="Имя пользователя" placeholder="admin" required />
        <FormField v-model="password" type="password" label="Пароль" placeholder="••••••••" required />

        <button type="submit" class="btn-primary w-full" :disabled="loading">
          {{ loading ? 'Вход…' : 'Войти' }}
        </button>
      </form>

      <p class="mt-5 rounded-xl bg-surface-muted px-4 py-3 text-center text-xs text-ink-muted">
        Демо-доступ: <span class="font-semibold text-ink">admin / admin123</span>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useUiStore } from '~/stores/ui'

definePageMeta({ layout: 'auth' })

const auth = useAuthStore()
const ui = useUiStore()

const username = ref('')
const password = ref('')
const loading = ref(false)

async function submit() {
  if (!username.value || !password.value) return
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    ui.success('Добро пожаловать!')
    await navigateTo('/')
  } catch {
    /* error toast handled centrally */
  } finally {
    loading.value = false
  }
}
</script>
