import React, { Fragment } from "react";

const ProductDetailSection = (props) => {
  const heading = props.heading;
  const content = props.content;
  return (
    <Fragment>
      <div style={{ padding: "4px 16px" }}>
        <div className="product-detail-section-heading">{heading}</div>
        <div className="product-detail-section-content">{content}</div>
      </div>
    </Fragment>
  );
};

export default ProductDetailSection;
