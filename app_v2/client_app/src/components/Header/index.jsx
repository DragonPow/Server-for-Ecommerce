import React from "react";
import { Layout } from "antd";
import "./index.less";
import HeaderNavBar from "./HeaderNavBar";
import HeaderWithSearch from "./HeadWithSearch";

const AppHeader = () => {
  return (
    <Layout>
      <header>
        <div className="container">
          <HeaderNavBar />
          <HeaderWithSearch />
        </div>
      </header>
    </Layout>
  );
};

export default AppHeader;
