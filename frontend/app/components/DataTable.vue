<template>
  <div class="card overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-surface-border bg-surface-muted/60 text-left text-xs uppercase tracking-wide text-ink-muted">
            <th v-for="c in columns" :key="c.key" class="whitespace-nowrap px-4 py-3 font-semibold" :class="c.class">
              {{ c.label }}
            </th>
            <th v-if="$slots.actions" class="px-4 py-3 text-right font-semibold">Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td :colspan="colspan" class="px-4 py-10 text-center text-ink-muted">Загрузка…</td>
          </tr>
          <tr v-else-if="!rows.length">
            <td :colspan="colspan" class="px-4 py-10 text-center text-ink-muted">{{ empty }}</td>
          </tr>
          <tr
            v-for="(row, idx) in rows"
            v-else
            :key="row.id ?? idx"
            class="border-b border-surface-border last:border-0 hover:bg-surface-muted/50"
          >
            <td v-for="c in columns" :key="c.key" class="px-4 py-3 align-middle" :class="c.class">
              <slot :name="c.key" :row="row" :value="row[c.key]">
                {{ display(row[c.key]) }}
              </slot>
            </td>
            <td v-if="$slots.actions" class="px-4 py-3 text-right">
              <div class="flex justify-end gap-1">
                <slot name="actions" :row="row" />
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Column {
  key: string
  label: string
  class?: string
}

const props = withDefaults(
  defineProps<{ columns: Column[]; rows: any[]; loading?: boolean; empty?: string }>(),
  { loading: false, empty: 'Нет данных' },
)
const slots = useSlots()

const colspan = computed(() => props.columns.length + (slots.actions ? 1 : 0))

function display(v: any) {
  if (v === null || v === undefined || v === '') return '—'
  return v
}
</script>
