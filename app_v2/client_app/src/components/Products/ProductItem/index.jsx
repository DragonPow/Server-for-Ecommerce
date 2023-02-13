import React from "react";
import { Col, Image } from "antd";

import "./index.less";
import Link from "antd/es/typography/Link";

const ProductItem = (props) => {
  const product = props.product;
  const image = product.image;
  return (
    <Col span={6}>
      <div className="product-item-card">
        <Link href={`/product/${product.id}`}>
          <div style={{ textAlign: "center" }}>
            <Image
              src={image}
              style={{ objectFit: "contain" }}
              preview={false}
            />
          </div>
          <div className="product-item-card__info">
            <div className="product-item-card__title">{product.name}</div>
            <div className="product-item__price">
              <span className="product-item__sale-price">
                {product.sale_price}
              </span>
            </div>
          </div>
        </Link>
      </div>
    </Col>
  );
};

export default ProductItem;
