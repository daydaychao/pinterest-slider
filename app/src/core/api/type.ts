type IResponse<T> = {
  ok: boolean
  code: number
  data: T
  message?: string
}

type ISearch = {
  playwright: {
    browserInstance: {
      _type: string
      _guid: string
    }
    page: {
      _type: string
      _guid: string
    }
    isLogin: boolean
  }
  resource_data: ISliderData[]
}

type ISliderData = {
  url: string
  description: string
  autoAltText: string
}

type IScreenshot = {
  screenshot: string
}
