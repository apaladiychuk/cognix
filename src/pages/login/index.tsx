import { Button } from "@/components/ui/button";
import GoogleIcon from "@/assets/svgs/google.svg?react"
import CognixLow from "@/assets/svgs/cognix-sm.svg?react"
import { Link } from "react-router-dom";

export function LoginComponent() {
  return (
    <>
    <div className="flex flex-col items-center justify-center h-screen space-y-5">
      <CognixLow className="w-20 h-20"/>
      <span className="text-2xl font-bold">Log In to CogniX</span>
      <div className="flex items-center justify-center">
        <Link
        to={`${import.meta.env.VITE_PLATFORM_API_URL}/auth/google/login`}
        >
      <Button
        variant='outline'
        size='xl'
        className="shadow-none bg-secondary"
        type="button"
      >
        <GoogleIcon
          className="fill-current mr-2 h-4 w-4"
        />
        Continue with Google
      </Button>
      </Link>
      </div>
    </div>
    </>
  )
}

export { LoginComponent as Component };
