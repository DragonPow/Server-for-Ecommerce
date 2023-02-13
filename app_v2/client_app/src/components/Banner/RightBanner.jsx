import { Col, Row } from "antd";
import React, { Fragment } from "react";

const RightBanner = (props) => {
  const urls = props.urls;
  return (
    <Fragment>
      <Row
        gutter={[0, 8]}
        style={{
          height: "100%",
          width: "100%"
        }}
      >
        {urls.map((item, index) => (
          <Col key={index} span={24}>
            <div
              style={{
                backgroundImage: `url("${item}")`,
                height: "100%",
                backgroundSize: "cover"
              }}
            ></div>
          </Col>
        ))}
      </Row>
    </Fragment>
  );
};

export default RightBanner;
