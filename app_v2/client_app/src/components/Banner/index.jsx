import { Col, Row } from "antd";
import React, { Fragment } from "react";
import BannerSlider from "./BannerSlider";
import { cdnHttp } from "../../http-wrapper";
import "./index.less";
import RightBanner from "./RightBanner";

const Banner = () => {
  const [bannerUrls, setBannerUrls] = React.useState([]);
  const [rightBannerUrls, setRightBannerUrls] = React.useState([]);

  React.useEffect(() => {
    async function fetchData() {
      const response = await cdnHttp.getBannerList();
      setBannerUrls(response.data.slider);
      setRightBannerUrls(response.data.right);
    }

    fetchData();
  }, []);

  return (
    <Fragment>
      <div className="header-banner">
        <div className="container" style={{ height: "100%" }}>
          <Row gutter={8}>
            <Col span={16}>
              <BannerSlider urls={bannerUrls} />
            </Col>
            <Col span={8}>
              <RightBanner urls={rightBannerUrls} />
            </Col>
          </Row>
        </div>
      </div>
    </Fragment>
  );
};

export default Banner;
