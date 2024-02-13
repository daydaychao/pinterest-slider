declare module '*.svg' {
  const content: string
  export default content
}

interface Window {
  electronApi: {
    triggerFullScreen: () => void
    // 在这里添加其他方法和属性的声明
  }
}

// 声明你的 ipcRenderer 全局变量
declare const ipcRenderer: {
  send: (channel: string) => void
  // 在这里添加其他方法和属性的声明
}
