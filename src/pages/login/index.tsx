import { Button } from "@/components/ui/button";
import CognixLow from "@/assets/svgs/cognix-sm.svg?react";
import { useState } from "react";
import { useGoogleLogin } from '@react-oauth/google';
import { useLocalStorage } from "@/lib/local-store";
import { api } from "@/lib/api";


export function LoginComponent() {
  const [error] = useState<string | null>(null); // Updated initial state

  const { set } = useLocalStorage()
  const loginUrl = `${import.meta.env.VITE_PLATFORM_API_URL}/auth/google/login`

  async function login(): Promise<void> {
    const authUrl = await api.get(
      `${import.meta.env.VITE_PLATFORM_API_URL}/auth/google/login?redirect_url=http://localhost:5173`
    ).then( response => {
      if (response.status === 200) {
        return response.data
      }
      return ""
    }
    )
    setTimeout(()=> {window.location.href = authUrl.data}, 6000)


    // const accessToken = await api.get(
    //   `${import.meta.env.VITE_PLATFORM_API_URL}/auth/google/login?redirect_url=http://localhost:5173`
    // ).then( response => {
    //   if (response.status === 200) {
    //     return response.data
    //   }
    //   return ""
    // }
    // )
}; 

  return (
    <>
      <div className="flex flex-col items-center justify-center h-screen space-y-5">
        <CognixLow className="w-20 h-20"/>
        <span className="text-2xl font-bold">Log In to CogniX</span>
        <div className="flex items-center justify-center">
          <Button
            variant='outline'
            size='xl'
            className="shadow-none bg-secondary"
            type="button"
            onClick={() => login()}
          >       
          Continue with Google
          </Button>
        </div>
        {error && <p className="text-red-500">{error}</p>}
      </div>
    </>
  );
}

export { LoginComponent as Component };
