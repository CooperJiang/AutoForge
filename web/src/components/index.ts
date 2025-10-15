/*
 * Cooper UI 组件库统一导出
 * 所有组件通过此文件统一导出，支持全局注册和按需引入
 */

import type { App, Plugin } from 'vue'

/* 导入所有基础组件 */
import BaseButton from './BaseButton'
import BaseInput from './BaseInput'
import BaseSelect from './BaseSelect'
import Dialog from './Dialog'
import Drawer from './Drawer'
import Message from './Message'
import MessageContainer from './MessageContainer'
import Pagination from './Pagination'
import Table from './Table'
import AppHeader from './AppHeader'
import TimePicker from './TimePicker'
import ParamInput from './ParamInput'
import JsonViewer from './JsonViewer'
import WeekDayPicker from './WeekDayPicker'
import MonthDayPicker from './MonthDayPicker'
import NextRunCountdown from './NextRunCountdown'
import TestResultDialog from './TestResultDialog'
import TaskDetailDialog from './TaskDetailDialog'
import ExecutionDetailDialog from './ExecutionDetailDialog'
import VariableSelector from './VariableSelector'
import RetryConfig from './RetryConfig'
import Tabs from './Tabs'
import Slider from './Slider'

/* 导出 message 工具函数 */
export { message } from '@/utils/message'

/* 组件映射表 - 用于全局注册 */
const componentMap = {

  BaseButton,
  BaseInput,
  BaseSelect,
  Dialog,
  Drawer,
  Message,
  MessageContainer,


  Pagination,
  Table,
  AppHeader,
  Tabs,


  TimePicker,
  ParamInput,
  WeekDayPicker,
  MonthDayPicker,
  Slider,


  JsonViewer,
  NextRunCountdown,


  TestResultDialog,
  TaskDetailDialog,
  ExecutionDetailDialog,


  VariableSelector,
  RetryConfig,
}

/* 按需导出所有组件 */
export {
  BaseButton,
  BaseInput,
  BaseSelect,
  Dialog,
  Drawer,
  Message,
  MessageContainer,
  Pagination,
  Table,
  AppHeader,
  Tabs,
  Slider,
  TimePicker,
  ParamInput,
  JsonViewer,
  WeekDayPicker,
  MonthDayPicker,
  NextRunCountdown,
  TestResultDialog,
  TaskDetailDialog,
  ExecutionDetailDialog,
  VariableSelector,
  RetryConfig,
}

/* 创建组件库插件 */
export const createCooperUI = (options: { components?: string[] } = {}): Plugin => ({
  install(app: App) {
    const { components = [] } = options

    /* 如果指定了组件列表，只注册指定组件 */
    if (components.length > 0) {
      components.forEach((name) => {
        if (componentMap[name as keyof typeof componentMap]) {
          app.component(name, componentMap[name as keyof typeof componentMap])
        }
      })
    } else {
      /* 否则注册所有组件 */
      Object.entries(componentMap).forEach(([name, component]) => {
        app.component(name, component)
      })
    }
  },
})

/* 默认插件（全量注册） */
const CooperUIPlugin: Plugin = createCooperUI()

export default CooperUIPlugin
