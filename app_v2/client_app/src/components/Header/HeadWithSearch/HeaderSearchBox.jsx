import React from "react";
import { Button } from "antd";
import { SearchOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";

const HeaderSearchBox = () => {
  const [searchKey, setSearchKey] = React.useState();
  const navigate = useNavigate();

  const handleSearch = () => {
    navigate(`/product?q=${searchKey}`);
  };

  const handleSearchKeyChange = (e) => {
    setSearchKey(e.target.value);
  };

  return (
    <div className="header-with-search__box">
      <input
        type="text"
        className="header-with-search__box-input"
        placeholder="Tìm kiếm"
        onChange={handleSearchKeyChange}
      />
      <Button
        icon={<SearchOutlined />}
        className="header-with-search__btn"
        style={{ width: "40px" }}
        onClick={handleSearch}
      />
    </div>
  );
};

export default HeaderSearchBox;
