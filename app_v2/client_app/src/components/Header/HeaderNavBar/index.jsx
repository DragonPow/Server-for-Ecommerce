import React from "react";
import "./index.less";

import { faFacebook, faInstagram } from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const HeaderNavBar = () => {
  return (
    <nav className="header-nav-bar">
      <ul className="header-nav-bar__group-ul">
        <li className="header-nav-bar__item-li">Web E-Commerce</li>
        <li
          className="header-nav-bar__item-li"
          style={{ display: "flex", alignItems: "center" }}
        >
          Kết nối
          <a href="https://facebook.com">
            <FontAwesomeIcon icon={faFacebook} color="white" />
          </a>
          <a href="https://instagram.com">
            <FontAwesomeIcon icon={faInstagram} color="white" />
          </a>
        </li>
      </ul>
      <ul className="header-nav-bar__group-ul">
        <li className="header-nav-bar__item-li">Trợ giúp</li>
        <li className="header-nav-bar__item-li">Thông báo</li>
      </ul>
    </nav>
  );
};

export default HeaderNavBar;
