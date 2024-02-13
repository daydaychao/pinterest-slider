import { Settings } from '../core/constants'
import { appInit, appLogin, appScreenshot, appSearch, appTriggerFullScreen } from '../core/controller'
import useSliderStore from '../core/stores'
import SearchInputUI from './SearchInputUI'

export const SearchPanel = () => {
  const isLoading = useSliderStore((s) => s.isLoading)
  return (
    <main className="flex flex-col w-200px">
      <div></div>
      <button onClick={appInit} disabled={isLoading}>
        init
      </button>
      <button onClick={appLogin} disabled={isLoading}>
        機器人登入
      </button>
      <SearchInputUI cb={appSearch} disabled={isLoading} />
      <button onClick={appScreenshot}>機器人拍截圖</button>
      {Settings.fullScreen && (
        <button onClick={appTriggerFullScreen} disabled={!Settings.fullScreen}>
          全螢幕切換
        </button>
      )}
      <div></div>
    </main>
  )
}
