<template>
  <div>
    <div class="flex gap-2 mb-6 items-center">
      <BaseInput v-model="filters.keyword" placeholder="搜索用户名或邮箱" style="width: 260px" />
      <BaseSelect
        v-model="filters.status"
        :options="[
          { label: '全部状态', value: '0' },
          { label: '正常', value: '1' },
          { label: '禁用', value: '2' },
        ]"
        placeholder="全部状态"
        style="width: 260px"
      />
      <BaseButton @click="loadUsers" variant="primary"> 搜索 </BaseButton>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b-2 border-border-primary text-left">
            <th class="pb-3 text-sm font-semibold text-text-secondary">用户名</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">邮箱</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">角色</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">任务统计</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">状态</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">创建时间</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="user in users"
            :key="user.id"
            class="border-b border-border-primary hover:bg-bg-hover transition-colors"
          >
            <td class="py-3 text-sm text-text-primary font-medium">{{ user.username }}</td>
            <td class="py-3 text-sm text-text-secondary">{{ user.email }}</td>
            <td class="py-3">
              <span
                :class="[
                  'px-2 py-1 text-xs font-medium rounded-full border',
                  user.role === 1
                    ? 'bg-purple-500/10 text-purple-600 dark:text-purple-400 border-purple-500/20'
                    : user.role === 2
                      ? 'bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20'
                      : 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20',
                ]"
              >
                {{ getRoleText(user.role) }}
              </span>
            </td>
            <td class="py-3">
              <div class="flex gap-2">
                <span
                  class="px-2 py-1 text-xs font-medium rounded border bg-blue-500/10 text-blue-600 dark:text-blue-400 border-blue-500/20"
                  :title="`总任务数: ${user.total_tasks}`"
                >
                  总: {{ user.total_tasks }}
                </span>
                <span
                  class="px-2 py-1 text-xs font-medium rounded border bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20"
                  :title="`启用任务数: ${user.enabled_tasks}`"
                >
                  启用: {{ user.enabled_tasks }}
                </span>
              </div>
            </td>
            <td class="py-3">
              <span
                :class="[
                  'px-2 py-1 text-xs font-medium rounded-full border',
                  user.status === 1
                    ? 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20'
                    : 'bg-slate-500/10 text-slate-600 dark:text-slate-400 border-slate-500/20',
                ]"
              >
                {{ user.status === 1 ? '正常' : '禁用' }}
              </span>
            </td>
            <td class="py-3 text-sm text-text-secondary">{{ formatTime(user.created_at) }}</td>
            <td class="py-3">
              <div class="flex gap-2">
                <button
                  @click="toggleUserStatus(user)"
                  :class="[
                    'px-3 py-1.5 text-xs font-medium rounded-lg transition-colors border',
                    user.status === 1
                      ? 'bg-slate-500/10 hover:bg-slate-500/20 text-slate-600 dark:text-slate-400 border-slate-500/20'
                      : 'bg-green-500/10 hover:bg-green-500/20 text-green-600 dark:text-green-400 border-green-500/20',
                  ]"
                  :title="user.status === 1 ? '禁用用户' : '启用用户'"
                >
                  {{ user.status === 1 ? '禁用' : '启用' }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <Pagination
      :current="currentPage"
      :page-size="pageSize"
      :total="total"
      :bordered="false"
      @change="handlePageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import * as adminApi from '@/api/admin'
import type { User } from '@/api/admin'
import { message } from '@/utils/message'
import { formatTime } from '@/utils/format'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import BaseButton from '@/components/BaseButton'
import Pagination from '@/components/Pagination'

// 用户列表
const users = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 筛选条件
const filters = ref({
  keyword: '',
  status: '0',
})

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

/**
 * 加载用户列表
 */
const loadUsers = async () => {
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: filters.value.keyword,
    }
    if (filters.value.status !== '0') {
      params.status = parseInt(filters.value.status)
    }
    const res = await adminApi.getUsers(params)
    users.value = res.data.users || []
    total.value = res.data.total
  } catch (error: any) {
    message.error(error.response?.data?.message || '加载用户列表失败')
  }
}

/**
 * 切换用户状态
 */
const toggleUserStatus = async (user: User) => {
  try {
    const newStatus = user.status === 1 ? 2 : 1
    await adminApi.updateUserStatus(user.id, newStatus)
    message.success(newStatus === 1 ? '已启用用户' : '已禁用用户')
    loadUsers()
  } catch (error: any) {
    message.error(error.response?.data?.message || '操作失败')
  }
}

/**
 * 获取角色文本
 */
const getRoleText = (role: number): string => {
  const roleMap: Record<number, string> = {
    1: '超级管理员',
    2: '管理员',
    3: '普通用户',
  }
  return roleMap[role] || '未知'
}

// 页码变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadUsers()
}

onMounted(() => {
  loadUsers()
})
</script>
