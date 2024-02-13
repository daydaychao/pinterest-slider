import useSliderStore from '../core/stores'
import { Swiper, SwiperSlide } from 'swiper/react'

// Import Swiper styles
import 'swiper/css'

export const Slider = () => {
  const resourceData = useSliderStore((s) => s.resourceData)
  return (
    <div className="flex justify-center items-center">
      <Swiper className="pinterestSlider" spaceBetween={10} slidesPerView={'auto'} autoplay={{ delay: 1000 }}>
        {resourceData.map((item, index) => {
          return (
            <SwiperSlide key={'slider' + index}>
              <div className="relative flex flex-row items-center justify-center">
                <img className="h-full" alt={item.autoAltText} src={item.url}></img>
                <small className="absolute bottom-1 bg-black opacity-50 p-2">{item.autoAltText}</small>
              </div>
            </SwiperSlide>
          )
        })}
      </Swiper>
    </div>
  )
}
