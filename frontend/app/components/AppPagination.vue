<template>
  <div v-if="totalPages > 1" class="flex items-center justify-between gap-3 px-1 py-2">
    <p class="text-sm text-ink-muted">
      Всего: <span class="font-semibold text-ink">{{ total }}</span>
    </p>
    <div class="flex items-center gap-1">
      <button class="btn-secondary !px-3 !py-1.5" :disabled="page <= 1" @click="go(page - 1)">
        Назад
      </button>
      <span class="px-2 text-sm text-ink-soft">{{ page }} / {{ totalPages }}</span>
      <button class="btn-secondary !px-3 !py-1.5" :disabled="page >= totalPages" @click="go(page + 1)">
        Вперёд
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ total: number; page: number; perPage: number }>()
const emit = defineEmits<{ 'update:page': [value: number] }>()

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.perPage)))

function go(p: number) {
  if (p >= 1 && p <= totalPages.value) emit('update:page', p)
}
</script>
