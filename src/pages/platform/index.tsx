import { SideBar } from "@/components/ui/sidebar";
import { useState } from "react";
import { Outlet } from "react-router-dom";

export function ApplicationRoot() {
  const [isSidebarOpen, setSidebarOpen] = useState<boolean>(false);
  return (
    <div className="flex h-screen bg-foreground">
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
          <Outlet />
        </div>
      </div>
    </div>
  );
}

export { ApplicationRoot as Component };
