import { Row, Pagination } from "antd";
import React, { Fragment } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { http } from "../../../http-wrapper";
import ProductItem from "../ProductItem";
import "./index.less";

const WrapperProduct = () => {
  const [currentPage, setCurrentPage] = React.useState(1);
  const [products, setProducts] = React.useState([]);
  const [productTotal, setProductTotal] = React.useState();

  const params = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const keyQuery = params.q;

  React.useEffect(() => {
    async function fetchProducts() {
      let filePath = `products?page=${currentPage}&page_size=20`;
      filePath +=
        searchParams &&
        searchParams.get("q") !== null &&
        searchParams.get("q") !== undefined
          ? `&key=${searchParams.get("q")}`
          : "";
      const resp = await http.getRequest(filePath);
      setProducts(resp.data?.data?.items);
      setProductTotal(resp.data?.data?.total_items);
    }

    fetchProducts();
  }, [currentPage, keyQuery, searchParams]);

  return (
    <Fragment>
      <div className="product-wrapper">
        <Row gutter={[16, 24]}>
          {products !== undefined &&
            products.map((item) => {
              return <ProductItem key={item.id} product={item} />;
            })}
        </Row>
      </div>
      <Pagination
        defaultCurrent={1}
        defaultPageSize={20}
        total={productTotal}
        showSizeChanger={false}
        onChange={(page, pagesize) => setCurrentPage(page)}
      />
    </Fragment>
  );
};

export default WrapperProduct;
