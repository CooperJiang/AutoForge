import { ref, watch, type Ref } from 'vue'
import SecureStorage from '../../utils/storage'

interface SecureStorageOptions {
  encrypt?: boolean
  expiry?: number
}

export function useSecureStorage<T>(
  key: string,
  defaultValue: T,
  options: SecureStorageOptions = {}
): [Ref<T>, (newValue: T) => void, () => void] {

  const initialValue = SecureStorage.getItem<T>(key, defaultValue) ?? defaultValue
  const value = ref<T>(initialValue) as Ref<T>


  watch(
    value,
    (newValue) => {
      SecureStorage.setItem(key, newValue, options)
    },
    { deep: true }
  )

  const setValue = (newValue: T) => {
    value.value = newValue

    SecureStorage.setItem(key, newValue, options)
  }

  const removeValue = () => {
    SecureStorage.removeItem(key)
    value.value = defaultValue
  }

  return [value, setValue, removeValue]
}


export function useAuthStorage<T>(
  key: string,
  defaultValue: T
): [Ref<T>, (newValue: T) => void, () => void] {
  return useSecureStorage(key, defaultValue, {
    encrypt: true,
    expiry: 7 * 24 * 60 * 60 * 1000,
  })
}

export function useSessionStorage<T>(
  key: string,
  defaultValue: T
): [Ref<T>, (newValue: T) => void, () => void] {
  return useSecureStorage(key, defaultValue, {
    encrypt: true,
    expiry: 24 * 60 * 60 * 1000,
  })
}

export function usePersistentStorage<T>(
  key: string,
  defaultValue: T
): [Ref<T>, (newValue: T) => void, () => void] {
  return useSecureStorage(key, defaultValue, {
    encrypt: true,

  })
}
