import { create } from 'zustand'

type ISliderState = {
  resourceData: ISliderData[]
  isLoading: boolean
}
export const useSliderStore = create<ISliderState>(() => ({
  resourceData: [],
  isLoading: false,
}))

export const setResourceData = (sliderData: ISliderData[]) => {
  useSliderStore.setState({ resourceData: sliderData })
}

export default useSliderStore
