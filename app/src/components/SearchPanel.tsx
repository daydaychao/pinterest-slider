import { Settings } from '../core/constants'
import { appLogin, appSearch, appTriggerFullScreen } from '../core/controller'
import useSliderStore from '../core/stores'
import SearchInputUI from './SearchInputUI'

export const SearchPanel = () => {
  const isLoading = useSliderStore((s) => s.isLoading)
  return (
    <main className="flex flex-col w-200px">
      <div></div>
      <button onClick={appLogin} disabled={isLoading}>
        命令機器人登入
      </button>
      <SearchInputUI cb={appSearch} disabled={isLoading} />
      <button onClick={appTriggerFullScreen} disabled={!Settings.fullScreen}>
        全螢幕切換
      </button>
      <div></div>
    </main>
  )
}
