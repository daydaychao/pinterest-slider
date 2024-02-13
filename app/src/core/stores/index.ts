import { create } from 'zustand'

type ISliderState = {
  resourceData: ISliderData[]
  screenshot: string
  isLoading: boolean
}
export const useSliderStore = create<ISliderState>(() => ({
  resourceData: [],
  screenshot: '',
  isLoading: false,
}))

export const setResourceData = (sliderData: ISliderData[]) => {
  useSliderStore.setState({ resourceData: sliderData })
}

export const setScreenshot = (screenshot: string) => {
  useSliderStore.setState({ screenshot: screenshot })
}

export default useSliderStore
