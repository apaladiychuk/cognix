import { SideBar } from "@/components/ui/sidebar";
import { useState } from "react";
import { Navigate, Outlet } from "react-router-dom";

export function ApplicationRoot() {
  const [isSidebarOpen, setSidebarOpen] = useState<boolean>(false);

  return (
    <div className="flex h-screen bg-foreground overflow-x-hidden">
      <div className="flex flex-row flex-grow">
        <div
          className={`bg-foreground text-white transition-all duration-300 ease-in-out ${
            isSidebarOpen ? "w-48" : "w-14"
          }`}
        >
          <SideBar
            isSideBarOpen={isSidebarOpen}
            setIsSideBarOpen={setSidebarOpen}
          />
        </div>
        <div className="flex flex-col flex-grow w-full align-center justify-center bg-background m-5 rounded-md">
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
