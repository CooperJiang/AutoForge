<template>
  <div class="tabs-container">
    <!-- Tab 导航 -->
    <div class="border-b border-border-primary">
      <div class="flex gap-1">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          @click="handleTabClick(tab.value)"
          :class="[
            'px-4 py-2 text-sm font-medium transition-colors relative',
            modelValue === tab.value
              ? 'text-primary'
              : 'text-text-secondary hover:text-text-primary',
          ]"
        >
          {{ tab.label }}
          <div
            v-if="modelValue === tab.value"
            class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary"
          ></div>
        </button>
      </div>
    </div>

    <!-- Tab 内容 -->
    <div class="tabs-content mt-6">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface Tab {
  label: string
  value: string
}

interface Props {
  modelValue: string
  tabs: Tab[]
}

defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const handleTabClick = (value: string) => {
  emit('update:modelValue', value)
}
</script>

<style scoped>
.tabs-container {
  @apply w-full;
}
</style>
