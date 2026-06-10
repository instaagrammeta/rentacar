<script setup>
// Root component: renders the active route and a global snackbar.
import { useAppStore } from '@/stores/app'
const app = useAppStore()
</script>

<template>
  <v-app>
    <router-view />

    <v-snackbar
      v-model="app.snackbar.show"
      :color="app.snackbar.color"
      location="top right"
      :timeout="4000"
    >
      {{ app.snackbar.text }}
      <template #actions>
        <v-btn variant="text" icon="mdi-close" @click="app.snackbar.show = false" />
      </template>
    </v-snackbar>
  </v-app>
</template>
