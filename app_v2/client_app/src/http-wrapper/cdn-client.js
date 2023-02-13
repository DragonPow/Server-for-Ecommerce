import axios from "axios";

const cdnClient = axios.create();

cdnClient.defaults.baseURL = "https://d116q2eo8du5hr.cloudfront.net";
cdnClient.defaults.headers = {
  "Content-Type": "application/json",
  Accept: "application/json"
};

function getBannerList() {
  return cdnClient.get("/banner/index.json").then((response) => response);
}

export const cdnHttp = {
  getBannerList
};
