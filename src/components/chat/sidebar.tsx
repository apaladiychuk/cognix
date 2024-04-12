import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp, PlusCircle } from "lucide-react";
import Cognix from "@/assets/svgs/cognix.svg?react";
import ChatSquare from "@/assets/svgs/chat-square.svg?react";
import SideBarIcon from "@/assets/svgs/sidebar-icon.svg?react";

export function SideBar() {
  const chats = [
    {
      id: 1,
      text: "Enabling “Get Attention”",
    },
    {
      id: 2,
      text: "Enabling “Get Attention”",
    },
    {
      id: 3,
      text: "Enabling “Get Attention”",
    },
    {
      id: 4,
      text: "Collaboard “Get Attention”123123123",
    },
    {
      id: 5,
      text: "Enabling “Get Attention”",
    },
  ];

  const settings = [
    {
      id: 1,
      text: "Connectors",
    },
    {
      id: 2,
      text: "Feedback",
    },
    {
      id: 3,
      text: "LLMs",
    },
    {
      id: 4,
      text: "Embeddings",
    },
    {
      id: 4,
      text: "Users",
    },
  ];

  const [isSidebarOpen, setSidebarOpen] = useState<boolean>(false);
  const [isHistoryOpen, setIsHistoryOpen] = useState<boolean>(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState<boolean>(false);

  return isSidebarOpen ? (
    <div className="motion-safe fixed w-48 bg-foreground text-white transition ease-in-out delay-1000">
      <div className="ml-2 mr-2 space-y-3">
        <div className="space-y-9">
          <Cognix className="mt-8" />
          <SideBarIcon
            className="absolute top-0 right-0 mr-2"
            onClick={() => {
              setSidebarOpen(!isSidebarOpen);
            }}
          />
          <div className="mb-4 space-y-5">
            <Button
              variant="outline"
              size="lg"
              className="shadow-none bg-primary w-full"
              type="button"
            >
              <PlusCircle className="h-4 w-4 mr-2" />
              New chat
            </Button>
          </div>
          <div
            className="flex items-center cursor-pointer"
            onClick={() => {
              setIsHistoryOpen(!isHistoryOpen);
            }}
          >
            {isHistoryOpen ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <ChevronDown className="h-4 w-4" />
            )}
            <span className="ml-2">Chat history</span>
          </div>
        </div>
        {isHistoryOpen && (
          <div className="flex flex-col ml-1 space-y-3 text-sm font-thin text-muted">
            {chats.map((chat) => (
              <div key={chat.id} className="flex flex-row items-center">
                <span className="truncate">{chat.text}</span>
              </div>
            ))}
          </div>
        )}
        <div
          className="flex items-center cursor-pointer pt-3"
          onClick={() => {
            setIsSettingsOpen(!isSettingsOpen);
          }}
        >
          {isSettingsOpen ? (
            <ChevronUp className="h-4 w-4" />
          ) : (
            <ChevronDown className="h-4 w-4" />
          )}
          <span className="ml-2">Settings</span>
        </div>
        {isSettingsOpen && (
          <div className="flex flex-col ml-4 space-y-3 text-muted">
            {settings.map((setting) => (
              <div key={setting.id} className="flex flex-row items-center">
                <ChatSquare className="h-4 w-4 mr-2 flex-shrink-0" />
                <span className="truncate">{setting.text}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  ) : (
    <div className="fixed w-10 bg-foreground text-white">
      <div className="mt-9 mr-2 space-y-3">
        <div className="space-y-9">
          <SideBarIcon
            className="absolute right-0 rotate-180"
            onClick={() => {
              setSidebarOpen(!isSidebarOpen);
            }}
          />
        </div>
      </div>
    </div>
  );
}

export { SideBar as Component };
