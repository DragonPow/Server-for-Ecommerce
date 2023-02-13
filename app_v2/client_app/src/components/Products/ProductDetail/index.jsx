import { Button, Col, Image, Row, Space } from "antd";
import React, { Fragment } from "react";
import { useParams } from "react-router-dom";
import { http } from "../../../http-wrapper";
import "./index.less";
import ProductDetailSection from "./ProductDetailSection";

const ProductDetail = (props) => {
  const [productDetail, setProductDetail] = React.useState();
  const [isFetchData, setIsFetchData] = React.useState(false);

  const params = useParams();

  const prodId = params.id;

  React.useEffect(() => {
    async function fetchData() {
      const resp = await http.getRequest(`products/${prodId}`);
      setProductDetail(resp.data.data);
      setIsFetchData(true);
    }

    if (!isFetchData) {
      fetchData();
    }
  }, [isFetchData, prodId]);

  const convertToCurrency = (price) => {
    return `${price} USD`;
  };

  return (
    <Fragment>
      <div>
        <div className="container">
          <div className="product-detail-wrapper">
            <Row gutter={16}>
              <Col span={10}>
                <div style={{ padding: "0 16px", textAlign: "center" }}>
                  <Image src={productDetail?.image} />
                </div>
              </Col>
              <Col span={14}>
                <Space direction="vertical">
                  <Space wrap className="product-detail-heading">
                    {productDetail?.name}
                  </Space>
                  <Space wrap className="product-detail-price">
                    {convertToCurrency(productDetail?.sale_price)}
                  </Space>
                  <Space wrap style={{ marginTop: "16px", padding: "16px" }}>
                    <Button type="primary" size="large">
                      Buy
                    </Button>
                  </Space>
                </Space>
              </Col>
            </Row>
          </div>
          <div className="product-detail-wrapper">
            <div>
              <ProductDetailSection
                heading="Description"
                content={productDetail?.description}
              />
            </div>
          </div>
        </div>
      </div>
    </Fragment>
  );
};

export default ProductDetail;
