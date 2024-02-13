import { fetcher } from './fetcher'
const API_URL = import.meta.env.VITE_API_URL

// [GET] Init
const appInit = async () => {
  const url = `${API_URL}/init`
  return fetcher({ type: 'GET', url }).then((res: IResponse<any>) => {
    return res
  })
}

// [GET] Login
const sliderLogin = async () => {
  const url = `${API_URL}/login`
  return fetcher({ type: 'GET', url }).then((res: IResponse<any>) => {
    return res
  })
}

// [GET] Search
const sliderSearch = async (key: string) => {
  const url = `${API_URL}/get?key=${key}`
  return fetcher({ type: 'GET', url }).then((res: IResponse<ISearch>) => {
    return res
  })
}

// [GET] Screenshot
const screenshot = async () => {
  const url = `${API_URL}/screenshot`
  return fetcher({ type: 'GET', url }).then((res: IResponse<IScreenshot>) => {
    return res
  })
}

const defaultData = {
  appInit,
  sliderLogin,
  sliderSearch,
  screenshot,
}

export default { ...defaultData }
