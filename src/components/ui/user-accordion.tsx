import { useLocalStorage } from "@/lib/local-store";
import { router } from "@/main";
import LogOutIcon from "@/assets/svgs/logout-icon.svg?react";
import { Button } from "./button";

export const UserAccordion = () => {
  const { remove } = useLocalStorage();

  function logOut() {
    router.navigate("/login");
    remove("access_token");
  }
  return (
    <Button
      className="absolute left-2 right-1 bottom-10 mt-2 w-45 bg-accent rounded-lg hover:bg-gray-700 flex items-center justify-start"
      onClick={logOut}
    >
      <LogOutIcon className="h-3.5 w-3.5 mr-2" />
      Log Out
    </Button>
  );
};
