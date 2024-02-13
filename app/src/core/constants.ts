export const Settings = {
  debug: true,
  fullScreen: false,
}
export const Color = {
  green: '\x1b[32m',
  red: '\x1b[31m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m',
  white: '\x1b[37m',
  close: '\x1b[0m',
}

export const runSettings = (browserWindow: any) => {
  if (Settings.debug) {
    // 打开开发者工具窗口
    browserWindow.setSize(1600, 800)
    browserWindow.webContents.openDevTools()
  }

  if (Settings.fullScreen) {
    console.log(Color.green + 'App ready' + Color.close)

    const { ipcMain } = require('electron')
    ipcMain.on('triggerFullScreen', () => {
      console.log(Color.yellow + '[icpMain] triggerFullScreen' + Color.close)
      browserWindow && browserWindow.setFullScreen(!browserWindow?.isFullScreen()) // 设置窗口全屏
    })
  }
}
