import "antd/dist/reset.css";

import AppHeader from "./components/Header";
import DataContainer from "./components/Products/DataContainer";
import "./App.less";
import AppFooter from "./components/Footer";
import Banner from "./components/Banner";

function App() {
  return (
    <div className="app">
      <AppHeader />
      <Banner />
      <DataContainer />
      <AppFooter />
    </div>
  );
}

export default App;
