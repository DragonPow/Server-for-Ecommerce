import React from "react";

const FooterGroupItem = (props) => {
  const name = props.name;
  return (
    <div className="footer-gr-item-link">
      <a href="#/">
        <span className="footer-gr-item-name" style={{ color: "#333" }}>
          {name}
        </span>
      </a>
    </div>
  );
};

export default FooterGroupItem;
