import React from "react";
import { Carousel } from "antd";

const BannerSlider = (props) => {
  const bannerUrls = props.urls;
  return (
    <Carousel autoplay>
      {bannerUrls !== undefined &&
        bannerUrls.map((item, index) => (
          <div key={index}>
            <div
              className="header-banner-item"
              style={{
                backgroundImage: `url('${item}')`
              }}
            ></div>
          </div>
        ))}
    </Carousel>
  );
};

export default BannerSlider;
