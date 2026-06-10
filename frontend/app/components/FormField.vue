<template>
  <label class="block">
    <span v-if="label" class="label">{{ label }}<span v-if="required" class="text-primary-500"> *</span></span>

    <select v-if="type === 'select'" class="input" :value="modelValue" @change="onInput">
      <option v-if="placeholder" value="">{{ placeholder }}</option>
      <option v-for="o in options" :key="String(o.value)" :value="o.value">{{ o.label }}</option>
    </select>

    <textarea
      v-else-if="type === 'textarea'"
      class="input"
      :rows="rows || 3"
      :value="modelValue as string"
      :placeholder="placeholder"
      @input="onInput"
    />

    <label v-else-if="type === 'checkbox'" class="flex cursor-pointer items-center gap-2">
      <input
        type="checkbox"
        class="h-4 w-4 rounded border-surface-border text-primary-500 focus:ring-primary-300"
        :checked="!!modelValue"
        @change="onCheckbox"
      />
      <span class="text-sm text-ink-soft">{{ checkboxLabel }}</span>
    </label>

    <input
      v-else
      class="input"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :step="step"
      @input="onInput"
    />
  </label>
</template>

<script setup lang="ts">
import type { Option } from '~/types'

withDefaults(
  defineProps<{
    label?: string
    modelValue: string | number | boolean | null
    type?: string
    options?: Option[]
    placeholder?: string
    required?: boolean
    rows?: number
    step?: string
    checkboxLabel?: string
  }>(),
  { type: 'text' },
)

const emit = defineEmits<{ 'update:modelValue': [value: any] }>()

function onInput(e: Event) {
  const el = e.target as HTMLInputElement
  emit('update:modelValue', el.type === 'number' ? (el.value === '' ? null : Number(el.value)) : el.value)
}
function onCheckbox(e: Event) {
  emit('update:modelValue', (e.target as HTMLInputElement).checked)
}
</script>
