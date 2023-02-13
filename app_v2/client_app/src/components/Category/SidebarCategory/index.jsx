import { Anchor, Affix } from "antd";
import React, { Fragment } from "react";
import "./index.less";

const SidebarCategory = () => {
  return (
    <Fragment>
      <Affix>
        <div className="sidebar-cate">
          <h3 className="sidebar-cate__title">Danh mục</h3>
          <Anchor
            items={[
              {
                key: "hottest",
                href: "/hottest",
                title: "Sản phẩm hot"
              },
              {
                key: "newest",
                href: "/newest",
                title: "Sản phẩm mới"
              }
            ]}
          />
        </div>
      </Affix>
    </Fragment>
  );
};

export default SidebarCategory;
