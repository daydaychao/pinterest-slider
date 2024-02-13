import api from './api'
import { setResourceData, setScreenshot } from './stores'

export const appInit = async () => {
  console.log('[前端] appInit')
  const data = await api.appInit()
  console.log('Init', data)
}

export const appLogin = async () => {
  console.log('[前端] appLogin')
  const data = await api.sliderLogin()
  console.log('Login', data)
}

export const appSearch = async (key: string) => {
  console.log('Search with key:', key)
  const res = await api.sliderSearch(key)
  if (!res?.ok) {
    console.log('%c[錯誤] 回傳有些問題', 'color:red', res)
    return
  }
  setResourceData(res.data.resource_data)
}

export const appTriggerFullScreen = () => {
  console.log('[前端] triggerFullScreen')
  window.electronApi.triggerFullScreen()
}

export const appScreenshot = async () => {
  console.log('[前端] appScreenshot')
  const res = await api.screenshot()
  if (!res?.ok) {
    console.log('%c[錯誤] 回傳有些問題', 'color:red', res)
    return
  }
  setScreenshot(res.data.screenshot)
}
