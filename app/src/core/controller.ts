import api from './api'
import { setResourceData } from './stores'

export const appLogin = async () => {
  console.log('[前端] appLogin')
  const data = await api.sliderLogin()
  console.log('Login', data)
}

export const appSearch = async (key: string) => {
  console.log('Search with key:', key)
  const res = await api.sliderSearch(key)
  if (!res.ok) {
    console.log('%c[錯誤] 回傳有些問題', 'color:red', res)
    return
  }
  setResourceData(res.data.resource_data)
}

export const appTriggerFullScreen = () => {
  console.log('[前端] triggerFullScreen')
  window.electronApi.triggerFullScreen()
}
