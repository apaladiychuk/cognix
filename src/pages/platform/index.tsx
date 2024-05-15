import { MobileNavBar } from "@/components/ui/mobile-navbar";
import { Navbar } from "@/components/ui/navbar";
import { useEffect, useState } from "react";
import { Navigate, Outlet } from "react-router-dom";

export function ApplicationRoot() {
  const [isSidebarOpen, setSidebarOpen] = useState(true);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      setIsMobile(window.innerWidth <= 765);
    };
    handleResize();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  return (
    <div className="flex h-screen bg-foreground">
      <div className="flex flex-row flex-grow h-full">
        <div
          className={`hidden md:flex h-full flex-col z-50 bg-foreground text-white transition-all duration-300 ease-in-out ${
            isSidebarOpen ? "w-48" : "w-20"
          }`}
        >
          {isMobile ? null : (
            <Navbar
              isSideBarOpen={isSidebarOpen}
              setIsSideBarOpen={setSidebarOpen}
            />
          )}
        </div>

        <div
          className={`md:hidden flex h-full absolute flex-col z-50 bg-foreground text-white transition-all duration-300 ease-in-out ${
            isSidebarOpen ? "w-48" : "w-1"
          }`}
        >
          {isMobile ? (
            <MobileNavBar
              isSideBarOpen={isSidebarOpen}
              setIsSideBarOpen={setSidebarOpen}
            />
          ) : null}
        </div>

        <div className="flex flex-col w-full flex-grow align-center justify-center bg-background lg:my-5 my-8 rounded-md">
          {localStorage.getItem("access_token") ? (
            <Outlet />
          ) : (
            <Navigate to="/login" />
          )}
        </div>
      </div>
    </div>
  );
}

export { ApplicationRoot as Component };
