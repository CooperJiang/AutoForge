<template>
  <div>
    <!-- 搜索筛选栏 -->
    <div class="flex gap-3 mb-6 items-center">
      <div class="flex-shrink-0" style="width: 200px;">
        <BaseInput
          v-model="filters.keyword"
          placeholder="搜索用户名或邮箱"
        />
      </div>
      <div class="flex-shrink-0" style="width: 150px;">
        <BaseSelect
          v-model="filters.status"
          :options="[
            { label: '全部状态', value: 0 },
            { label: '正常', value: 1 },
            { label: '禁用', value: 2 }
          ]"
          placeholder="全部状态"
        />
      </div>
      <BaseButton
        @click="loadUsers"
        variant="primary"
        class="flex-shrink-0"
      >
        搜索
      </BaseButton>
    </div>

    <!-- 表格 -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b-2 border-slate-200 text-left">
            <th class="pb-3 text-sm font-semibold text-slate-700">用户名</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">邮箱</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">角色</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">任务统计</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">状态</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">创建时间</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="user in users"
            :key="user.id"
            class="border-b border-slate-100 hover:bg-slate-50 transition-colors"
          >
            <td class="py-3 text-sm text-slate-900 font-medium">{{ user.username }}</td>
            <td class="py-3 text-sm text-slate-600">{{ user.email }}</td>
            <td class="py-3">
              <span
                :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  user.role === 1
                    ? 'bg-purple-100 text-purple-700'
                    : user.role === 2
                    ? 'bg-blue-100 text-blue-700'
                    : 'bg-green-100 text-green-700'
                ]"
              >
                {{ getRoleText(user.role) }}
              </span>
            </td>
            <td class="py-3">
              <div class="flex gap-2">
                <span class="px-2 py-1 text-xs font-medium rounded bg-blue-50 text-blue-700" :title="`总任务数: ${user.total_tasks}`">
                  总: {{ user.total_tasks }}
                </span>
                <span class="px-2 py-1 text-xs font-medium rounded bg-green-50 text-green-700" :title="`启用任务数: ${user.enabled_tasks}`">
                  启用: {{ user.enabled_tasks }}
                </span>
              </div>
            </td>
            <td class="py-3">
              <span
                :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  user.status === 1
                    ? 'bg-green-100 text-green-700'
                    : 'bg-gray-100 text-gray-700'
                ]"
              >
                {{ user.status === 1 ? '正常' : '禁用' }}
              </span>
            </td>
            <td class="py-3 text-sm text-slate-600">{{ formatTime(user.created_at) }}</td>
            <td class="py-3">
              <div class="flex gap-2">
                <button
                  @click="toggleUserStatus(user)"
                  :class="[
                    'px-3 py-1.5 text-xs font-medium rounded-lg transition-colors',
                    user.status === 1
                      ? 'bg-gray-100 hover:bg-gray-200 text-gray-700'
                      : 'bg-green-100 hover:bg-green-200 text-green-700'
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

    <!-- 分页 -->
    <div class="flex justify-between items-center mt-6">
      <div class="text-sm text-slate-600">
        共 {{ total }} 条记录
      </div>
      <div class="flex gap-2">
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="px-4 py-2 bg-slate-100 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          上一页
        </button>
        <span class="px-4 py-2 text-sm text-slate-700">
          {{ currentPage }} / {{ totalPages }}
        </span>
        <button
          @click="nextPage"
          :disabled="currentPage >= totalPages"
          class="px-4 py-2 bg-slate-100 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          下一页
        </button>
      </div>
    </div>
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

// 用户列表
const users = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 筛选条件
const filters = ref({
  keyword: '',
  status: 0
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
      keyword: filters.value.keyword
    }
    if (filters.value.status > 0) {
      params.status = filters.value.status
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
    3: '普通用户'
  }
  return roleMap[role] || '未知'
}

// 上一页
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadUsers()
  }
}

// 下一页
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadUsers()
  }
}

onMounted(() => {
  loadUsers()
})
</script>
