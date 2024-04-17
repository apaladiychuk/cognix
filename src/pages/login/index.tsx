import { Button } from "@/components/ui/button";
import CognixLow from "@/assets/svgs/cognix-sm.svg?react";
import { useState } from "react";
import { useGoogleLogin } from '@react-oauth/google';
import { useLocalStorage } from "@/lib/local-store";


export function LoginComponent() {
  const [error] = useState<string | null>(null); // Updated initial state

  const { set } = useLocalStorage()

  const login = useGoogleLogin({
    onSuccess: credentialResponse => {
      set("access_token", credentialResponse as any),
      console.log(credentialResponse)
    },
    redirect_uri: `${window.location.origin}/google/callback`,
    ux_mode:"redirect",
    flow: 'auth-code',
  }); 

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
