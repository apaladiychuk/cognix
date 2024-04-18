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
  timeout: 10000,
  // headers: {"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdG9rZW4iOiJ5YTI5LmEwQWQ1Mk4zOWUtaVlNOVhSOVBWWndqZWRzODdJemVXS1FrZkdYVG9WblE3X1pYRmV5ZlRlZk0xU0xpa0NlY1B1RTY0eDQ1VUdqOUVqMUhwNWxhdDFpb1pkUEZJcENfWlZjUk5UNXpUOWhUN2NndUY2TFg1MVB5MXg4SjNjdUtYMko5VG1oY2RzbUNRREJqd21uT1VDcUFUZ1NsMmY0bVhzYVoxUkNhQ2dZS0FhNFNBUk1TRlFIR1gyTWlMUmZGOXNOemY1YnFyRDB0LVdJbFBBMDE3MSIsInJlZnJlc2hfdG9rZW4iOiIiLCJ1c2VyIjp7ImlkIjoiMDAzZGE5ODAtNDM3Ni00MjBmLWExNWItYWU4ZjM5ZThhNDEwIiwidGVuYW50X2lkIjoiMTQ1ZjJlZjAtZjY1My00MzE4LWFhZjYtMWNhNWVjMTIxYTVhIiwidXNlcl9uYW1lIjoidmFkeW0ubWFzbG92c2t5aUBwZWNvZGVzb2Z0d2FyZS5jb20iLCJmaXJzdF9uYW1lIjoiVmFkeW0iLCJsYXN0X25hbWUiOiJNYXNsb3Zza3lpIiwicm9sZXMiOlsic3VwZXJfYWRtaW4iXX19.9XFcMgoh2kODUAjmnZpX4OydY-9o85pBmfiUrmW_-Io"},
  // responseType: 'json'
});