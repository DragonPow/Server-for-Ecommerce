import React from "react";
import Footer from "../components/Footer";
import AppHeader from "../components/Header";
import ProductDetail from "../components/Products/ProductDetail";

const ProductPage = () => {
  return (
    <div className="app">
      <AppHeader />
      <ProductDetail />
      <Footer />
    </div>
  );
};

export default ProductPage;
