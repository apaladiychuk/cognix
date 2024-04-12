import { useState } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { SideBar } from "./sidebar";
import SendIcon from "@/assets/svgs/send-icon.svg?react";
import FileIcon from "@/assets/svgs/file-icon.svg?react";
import { Card } from "../ui/card";

export function ChatComponent() {
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
        <div className="flex flex-col flex-grow align-center justify-center bg-background m-5 mr-0 rounded-md rounded-r-none w-4/6">
          <div className="flex items-center justify-center pt-8">
            <span className="text-4xl font-bold">
              Which assistant do you want
            </span>
          </div>
          <div className="flex items-center justify-center pt-1">
            <span className="text-4xl font-bold">to chat with today?</span>
          </div>
          <div className="flex items-center justify-center pt-8">
            <span className="font-thin text-base text-muted">
              Or ask a question immediately to use the CogniX assistant
            </span>
          </div>
          <div className="flex pt-10 mx-20 space-x-5">
            <div className="flex-1">
              <Card
                title="CogniX"
                text="Assistant with access to documents
              from your Connected Sources"
              />
            </div>
            <div className="flex-1">
              <Card
                title="Paraphrase"
                text="Assistant that is heavily constrained and only provides exact quotes from Connected Sources."
              />
            </div>
          </div>
          <div className="flex-grow p-4">{/* Content here */}</div>
          <div className="flex items-center justify-between space-x-3 p-4 ml-12 mr-12">
            <Input
              placeholder="Ask me anything..."
              className="flex-grow rounded-lgw-1/2"
            />
            <Button
              size="icon"
              variant="outline"
              className="w-12 h-12 bg-primary hover:bg-foreground"
              type="button"
              onClick={() => {
                alert("Message will be sent");
              }}
            >
              <SendIcon className="size-5" />
            </Button>
          </div>
          <div className="flex items-center justify-center pb-4">
            <span className="text-xs font-thin text-muted">
              CogniX can make mistakes. Consider checking critical information.
            </span>
          </div>
        </div>
        <div className="flex ml-0 m-5 w-1/5 flex-col bg-white rounded-md rounded-l-none">
          <div className="content-start space-x-2 pl-4">
            <div className="flex content-start space-x-2 pt-5 pl-3">
              <FileIcon />
              <span className="font-bold">Retrieved Knowledge</span>
            </div>
            <div className="flex pt-5">
              <span className="font-thin text-sm text-muted">
                When you run ask a question, the
              </span>
            </div>
            <div className="flex pt-1">
              <span className="font-thin text-sm text-muted">
                retrieved knowledge will show up here
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
