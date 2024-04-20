import axios from "axios";

axios.interceptors.request.use(
  (config) => {
    const token = JSON.parse(localStorage.getItem("access_token") as string);
    config.headers["Authorization"] = "Bearer " + token;
    config.headers["Content-Type"] = "application/json";
    config.headers["Accept"] = "application/json";
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
