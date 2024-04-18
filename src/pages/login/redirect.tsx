import { api } from "@/lib/api";
import { useLocalStorage } from "@/lib/local-store";
import { useEffect } from "react";

const { set } = useLocalStorage()

function RedirectComponent() {
    
    useEffect(() => {
        api.get(
            `${import.meta.env.VITE_PLATFORM_API_URL}/auth/google/callback?${window.location.href.split("?")[1]}`
          ).then( response => {
            if (response.status === 200) {
              response.data
              set("access_token", response.data.data)
            }
            return ""
          }
        ).catch( e => {
            console.log(e)
        }
        )
      }, []);
    

    return <></>
}

export { RedirectComponent as Component}