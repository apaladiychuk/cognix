import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { SideBar } from "./sidebar";
import SendIcon from "@/assets/svgs/send-icon.svg?react";

export function ChatComponent() {
  return (
    <div className="grid grid-cols-12 h-screen bg-foreground">
      <div className="col-span-2">
        <SideBar />
      </div>
      <div className="col-span-7 bg-background m-5 rounded-md">
        <div className="flex w-3/4 items-center p-4">
          <Input
            placeholder="Ask me anything..."
            className="flex rounded-lg mr-2 bottom-10"
          />
          <Button
            size="icon"
            className="w-12 h-12"
            onClick={() => {
              console.log("tapped");
            }}
          >
            <SendIcon className="size-4"/>
          </Button>
        </div>
      </div>
      <div className="col-span-3 bg-white">
          asfasfaf
        </div>
    </div>
  );
}

export { ChatComponent as Component };
