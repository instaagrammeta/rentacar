<script setup>
// Small KPI card used on the dashboard.
defineProps({
  title: { type: String, required: true },
  value: { type: [String, Number], default: 0 },
  icon: { type: String, default: 'mdi-information' },
  color: { type: String, default: 'primary' },
})
</script>

<template>
  <v-card class="pa-4" elevation="2">
    <div class="d-flex align-center">
      <v-avatar :color="color" size="52" class="mr-4">
        <v-icon :icon="icon" color="white" size="28" />
      </v-avatar>
      <div>
        <div class="text-h5 font-weight-bold">{{ value }}</div>
        <div class="text-body-2 text-medium-emphasis">{{ title }}</div>
      </div>
    </div>
  </v-card>
</template>
