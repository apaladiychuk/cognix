import axios from 'axios'
import { useLocalStorage } from "@/lib/local-store";
import { router } from '@/main'

const { get, set } = useLocalStorage()

axios.interceptors.request.use(
  config => {
    const token = get("access_token", "" as any)
    config.headers['Authorization'] = 'Bearer ' + token
    config.headers['Content-Type'] = 'application/json';
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

axios.interceptors.response.use(
  response => {
    return response
  },
  function (error) {
    const originalRequest = error.config

    if (
      error.response.status === 401 &&
      originalRequest.url === `${window.location.origin}/api/google/login`
    ) {
      router.navigate('/login')
      return Promise.reject(error)
    }

    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      const refreshToken = get("refresh_token", "" as any)
      return axios
        .post('/auth/token', {
          refresh_token: refreshToken
        })
        .then(res => {
          if (res.status === 201) {
            set("access_token", res.data)
            axios.defaults.headers.common['Authorization'] =
              'Bearer ' + get("access_token", "" as any)
            return axios(originalRequest)
          }
        })
    }
    return Promise.reject(error)
  }
)

export const api = axios.create({
  baseURL: import.meta.env.VITE_PLATFORM_API_URL,
  timeout: 1000,
});