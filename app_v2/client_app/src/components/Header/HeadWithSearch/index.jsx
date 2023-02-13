import React from "react";
import LogoSVG from "./LogoSVG";
import "./index.less";
import { Link } from "react-router-dom";
import HeaderSearchBox from "./HeaderSearchBox";

const HeaderWithSearch = () => {
  return (
    <div className="header-with-search">
      <Link to="/home" className="header-with-search__logo-box">
        <LogoSVG></LogoSVG>
      </Link>
      <HeaderSearchBox />
    </div>
  );
};

export default HeaderWithSearch;
