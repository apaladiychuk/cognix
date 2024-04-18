import { api } from "@/lib/api";
import { useEffect } from "react";


export function RedirectComponent() {
    async function getAccessToken() {
        console.log("Called login")
        const accessToken = await api.get(
            `${import.meta.env.VITE_PLATFORM_API_URL}/auth/google/callback?${window.location.href.split("?")[1]}`
          ).then( response => {
            if (response.status === 200) {
              return response.data
            }
            return ""
          }
        )

    }
    
    
    useEffect(() => {
        getAccessToken();
      }, []);
    

    return <></>
}

export { RedirectComponent as Component}