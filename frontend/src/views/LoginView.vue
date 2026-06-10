<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const router = useRouter()
const auth = useAuthStore()
const app = useAppStore()

const username = ref('')
const password = ref('')
const loading = ref(false)

async function submit() {
  if (!username.value || !password.value) {
    app.notify(t('app.required'), 'warning')
    return
  }
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    app.notify(t('auth.welcome'))
    router.push('/dashboard')
  } catch (e) {
    // error surfaced by the axios interceptor
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <v-app>
    <v-main class="login-bg">
      <v-container class="fill-height" fluid>
        <v-row justify="center" align="center">
          <v-col cols="12" sm="8" md="4">
            <v-card class="pa-6" elevation="8">
              <div class="text-center mb-6">
                <v-icon size="56" color="primary" icon="mdi-car-key" />
                <h1 class="text-h5 font-weight-bold text-primary mt-2">Rentacar CRM</h1>
                <p class="text-medium-emphasis">{{ t('auth.login') }}</p>
              </div>
              <v-form @submit.prevent="submit">
                <v-text-field
                  v-model="username"
                  :label="t('auth.username')"
                  prepend-inner-icon="mdi-account"
                  autofocus
                />
                <v-text-field
                  v-model="password"
                  :label="t('auth.password')"
                  type="password"
                  prepend-inner-icon="mdi-lock"
                />
                <v-btn
                  type="submit"
                  color="primary"
                  size="large"
                  block
                  :loading="loading"
                  class="mt-2"
                >
                  {{ t('auth.signIn') }}
                </v-btn>
              </v-form>
              <p class="text-caption text-center text-medium-emphasis mt-4">
                По умолчанию: admin / admin123
              </p>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<style scoped>
.login-bg {
  background: linear-gradient(135deg, #1b5e20 0%, #66bb6a 100%);
}
</style>
