import React from "react";

const FooterGroupHeading = (props) => {
  const title = props.title;

  return (
    <div className="footer-gr-item-heading">
      <span>{title}</span>
    </div>
  );
};

export default FooterGroupHeading;
