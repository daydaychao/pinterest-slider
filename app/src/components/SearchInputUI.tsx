import { useState } from 'react'
import './searchInputUI.css'

interface SearchInputUIProps {
  cb: (key: string) => void
  disabled: boolean
}

export const SearchInputUI = (props: SearchInputUIProps) => {
  const [key, setKey] = useState('') // 使用 useState Hook 定义 ff 状态，并设置初始值为 ''

  const handleInputChange = (event: any) => {
    setKey(event.target.value)
  }

  const handleSubmit = (event: any) => {
    event.preventDefault() // 阻止默认提交行为

    props.cb(key) // api search
    setKey('')
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="relative flex flex-row app-submit-wrapper">
        <input
          className="app-submit-left"
          type="text"
          value={key}
          disabled={props.disabled}
          onChange={handleInputChange}
        />
        <button className="absolute right-0 app-submit-right" type="submit" disabled={props.disabled}>
          Send
        </button>
      </div>
    </form>
  )
}

export default SearchInputUI
