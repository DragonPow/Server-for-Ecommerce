import axios from "axios";

const apiClient = axios.create();
apiClient.defaults.baseURL =
  "http://uitk14-ktpm-api-lb-c74d9c813b4b2e8f.elb.ap-southeast-1.amazonaws.com/http";
apiClient.defaults.headers = {
  "Content-Type": "application/json",
  Accept: "application/json"
};

//All request will wait 2 seconds before timeout
apiClient.defaults.timeout = 2000;

function getRequest(urlFile) {
  return apiClient.get(`/${urlFile}`).then((response) => response);
}

function postRequest(urlFile, payload) {
  return apiClient.post(`/${urlFile}`, payload).then((response) => response);
}

function patchRequest(urlFile, payload) {
  return apiClient.patch(`/${urlFile}`, payload).then((response) => response);
}

function deleteRequest(urlFile) {
  return apiClient.delete(`/${urlFile}`).then((response) => response);
}

export const http = {
  getRequest,
  postRequest,
  patchRequest,
  deleteRequest
};
