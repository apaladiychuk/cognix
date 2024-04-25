import axios from "axios";

axios.interceptors.request.use(
  (config) => {
    const accessTokenString = localStorage.getItem("access_token");
    const token = accessTokenString ? JSON.parse(String(localStorage.getItem("access_token"))) : "";
    config.headers["Authorization"] = "Bearer " + token;
    config.headers["Content-Type"] = "application/json";
    config.headers["Accept"] = "application/json";
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
