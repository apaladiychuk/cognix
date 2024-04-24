import { useContext, useState } from "react";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp, PlusCircle } from "lucide-react";
import Cognix from "@/assets/svgs/cognix.svg?react";
import CognixSmall from "@/assets/svgs/cognix-sm.svg?react";
import SideBarIcon from "@/assets/svgs/sidebar-icon.svg?react";
import SideBarClosedIcon from "@/assets/svgs/sidebar-closed-icon.svg?react";
import ConnectorsIcon from "@/assets/svgs/connectors.svg?react";
import FeedbackIcon from "@/assets/svgs/feedback.svg?react";
import LLMIcon from "@/assets/svgs/llm.svg?react";
import EmbeddingIcon from "@/assets/svgs/embedding.svg?react";
import UsersIcon from "@/assets/svgs/users.svg?react";
import ConfigIcon from "@/assets/svgs/config.svg?react";
import React from "react";
import { Link, NavLink } from "react-router-dom";
import { AuthContext } from "@/context/AuthContext";

export interface SideBarProps {
  isSideBarOpen: boolean;
  setIsSideBarOpen: (isSideBarOpen: boolean) => void;
}

const SideBar: React.FC<SideBarProps> = ({
  isSideBarOpen,
  setIsSideBarOpen,
}) => {
  const settings = [
    {
      id: 1,
      text: "Connectors",
      icon: <ConnectorsIcon className="h-4 w-4" />,
      link: "/settings/connectors/existing-connectors",
    },
    {
      id: 2,
      text: "Feedback",
      icon: <FeedbackIcon className="h-4 w-4" />,
      link: "/settings/feedback",
    },
    {
      id: 3,
      text: "LLMs",
      icon: <LLMIcon className="h-4 w-4" />,
      link: "/settings/llms",
    },
    {
      id: 4,
      text: "Embeddings",
      icon: <EmbeddingIcon className="h-4 w-4" />,
      link: "/settings/embeddings",
    },
    {
      id: 5,
      text: "Users",
      icon: <UsersIcon className="h-4 w-4" />,
      link: "/settings/users",
    },
    {
      id: 6,
      text: "Config Map",
      icon: <ConfigIcon className="h-4 w-4" />,
      link: "/settings/config",
    },
  ];

  const [isHistoryOpen, setIsHistoryOpen] = useState<boolean>(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState<boolean>(false);
  const { firstName, lastName, chats } = useContext(AuthContext);

  return isSideBarOpen ? (
    <div className="ml-2 mr-2 space-y-5">
      <div className="space-y-9">
        <div className="flex items-center mt-8 space-x-3">
          <Link to={"/"}>
            <Cognix className="h-10" />
          </Link>
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
        <div className="flex flex-col ml-3 space-y-4 text-2sm font-thin text-muted">
          {chats.slice(0, 4).map((chat) => (
            <NavLink
              key={chat.id}
              to={`/platform`}
              className="flex flex-row items-center"
            >
              <span className="truncate">{chat.description}</span>
            </NavLink>
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
        <div className="flex flex-col ml-3 space-y-4 text-muted text-2sm">
          {settings.map((setting) => (
            <NavLink
              key={setting.id}
              to={setting.link}
              className="flex flex-row items-center"
            >
              <div className="mr-2">{setting.icon}</div>
              <span className="truncate">{setting.text}</span>
            </NavLink>
          ))}
        </div>
      )}
      <div className="fixed bottom-7 pl-2 flex items-center justify-center space-x-2">
        <div className="w-7 h-7 rounded-md bg-accent flex items-center justify-center">
          <span className="text-xs">
            {firstName && `${firstName.charAt(0)}`}
          </span>
          <span className="text-xs">
            {lastName && `${lastName.charAt(0)}`}
          </span>
        </div>
        <span className="text-sm">{firstName + " " + lastName}</span>
      </div>
    </div>
  ) : (
    <div className="flex flex-col ml-3 mr-2 space-y-5">
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
          <Link to={"/"}>
            <CognixSmall className="h-9 w-9" />
          </Link>
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
        <div className="flex flex-col ml-1 space-y-3 text-2sm font-thin text-muted">
          {chats.slice(0, 4).map((chat) => (
            <NavLink
              key={chat.id}
              to={`/platform`}
              className="flex flex-row items-center"
            >
              <span className="truncate">
                {chat.messages.length > 0
                  ? chat.messages[chat.messages.length - 1].message
                  : ""}
              </span>
            </NavLink>
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
        <div className="flex flex-col ml-4 space-y-4 text-muted">
          {settings.map((setting) => (
            <NavLink
              key={setting.id}
              to={setting.link}
              className={({ isActive }) => (isActive ? "active" : "inactive")}
            >
              <div className="h-6 w-4 mr-6 flex-shrink-0 fill-foreground/70 group-[.is-active]:fill-accent/95">
                {setting.icon}
              </div>
            </NavLink>
          ))}
        </div>
      )}
    </div>
  );
};

export { SideBar };
