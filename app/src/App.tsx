import reactLogo from './assets/react.svg'
import viteLogo from '/electron-vite.animate.svg'
import { SearchPanel } from './components/SearchPanel'
import { Slider } from './components/Slider'

function App() {
  return (
    <div className="h-full w-full flex flex-col justify-between border border-1">
      <SearchPanel />

      <Slider />

      <footer className="flex flex-row justify-center items-center text-0.8rem p-0.5rem text-center">
        <div className="flex flex-row gap-2 justify-center p-2">
          <img src={viteLogo} className="logo" alt="Vite logo" />
          <img src={reactLogo} className="logo react" alt="React logo" />
        </div>
        copyright dayday 2024
      </footer>
    </div>
  )
}

export default App
