import { Col, Layout, Row } from "antd";
import React from "react";
import SidebarCategory from "../../Category/SidebarCategory";
import WrapperProduct from "../WrapperProduct";

const DataContainer = () => {
  return (
    <Layout>
      <div style={{ margin: "36px 0" }}>
        <div className="container">
          <Row gutter={16}>
            <Col md={4} lg={4} xl={4}>
              <SidebarCategory />
            </Col>
            <Col md={20} lg={20} xl={20}>
              <WrapperProduct />
            </Col>
          </Row>
        </div>
      </div>
    </Layout>
  );
};

export default DataContainer;
