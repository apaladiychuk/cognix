import { useState } from "react";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp, PlusCircle } from "lucide-react";
import Cognix from "@/assets/svgs/cognix.svg?react";
import CognixSmall from "@/assets/svgs/cognix-sm.svg?react";
import ChatSquare from "@/assets/svgs/chat-square.svg?react";
import SideBarIcon from "@/assets/svgs/sidebar-icon.svg?react";
import SideBarClosedIcon from "@/assets/svgs/sidebar-closed-icon.svg?react";
import React from "react";

export interface SideBarProps {
  isSideBarOpen: boolean;
  setIsSideBarOpen: (isSideBarOpen: boolean) => void;
}

const SideBar: React.FC<SideBarProps> = ({
  isSideBarOpen,
  setIsSideBarOpen,
}) => {
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

  const [isHistoryOpen, setIsHistoryOpen] = useState<boolean>(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState<boolean>(false);

  return isSideBarOpen ? (
    <div className="ml-2 mr-2 space-y-3">
      <div className="space-y-9">
        <div className="flex items-center mt-8 space-x-3">
          <Cognix />
          <SideBarIcon
            className="cursor-pointer"
            onClick={() => {
              setIsSideBarOpen(!isSideBarOpen);
            }}
          />
        </div>
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
  ) : (
    <div className="flex flex-col ml-3 mr-2 space-y-3">
      <div className="space-y-6">
        <div className="flex pt-9 pl-3">
          <SideBarClosedIcon
            className="cursor-pointer"
            onClick={() => {
              setIsSideBarOpen(!isSideBarOpen);
            }}
          />
        </div>
        <div className="ml-1.5">
          <CognixSmall className="h-9 w-9" />
        </div>
        <div>
          <Button
            variant="outline"
            size="icon"
            className="ml-1.5 bg-primary h-9 w-9"
            type="button"
          >
            <PlusCircle className="h-4 w-4" />
          </Button>
        </div>
        <div
          className="flex items-center cursor-pointer"
          onClick={() => {
            setIsHistoryOpen(!isHistoryOpen);
          }}
        >
          {isHistoryOpen ? (
            <ChevronUp className="ml-4 h-4 w-4" />
          ) : (
            <ChevronDown className="ml-4 h-4 w-4" />
          )}
        </div>
      </div>
      {isHistoryOpen && (
        <div className="flex flex-col ml-1 space-y-3 text-sm font-thin text-muted">
          {chats.map((chat) => (
            <div key={chat.id} className="flex flex-row items-center">
              <span className="truncate md:text-clip">{chat.text}</span>
            </div>
          ))}
        </div>
      )}
      <div
        className="flex cursor-pointer pt-3"
        onClick={() => {
          setIsSettingsOpen(!isSettingsOpen);
        }}
      >
        {isSettingsOpen ? (
          <ChevronUp className="ml-4 h-4 w-4" />
        ) : (
          <ChevronDown className="ml-4 h-4 w-4" />
        )}
      </div>
      {isSettingsOpen && (
        <div className="flex flex-col ml-4 space-y-3 text-muted">
          {settings.map((setting) => (
            <div key={setting.id} className="flex flex-row items-center">
              <ChatSquare className="h-4 w-4 mr-2 flex-shrink-0" />
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export { SideBar };
