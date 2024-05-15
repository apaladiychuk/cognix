import Cognix from "@/assets/svgs/cognix.svg?react";
import SideBarIcon from "@/assets/svgs/sidebar-icon.svg?react";
import { AuthContext } from "@/context/AuthContext";
import { settings } from "@/lib/utils";
import { router } from "@/main";
import { ChevronDown, ChevronUp, PlusCircle } from "lucide-react";
import { memo, useContext, useState } from "react";
import { Button } from "./button";
import { Link, NavLink } from "react-router-dom";
import { UserAccordion } from "./user-accordion";

interface Props {
  isSideBarOpen: boolean;
  setIsSideBarOpen: (isSideBarOpen: boolean) => void;
}

export const Sidebar = memo(({ isSideBarOpen, setIsSideBarOpen }: Props) => {
  const [isHistoryOpen, setIsHistoryOpen] = useState<boolean>(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState<boolean>(false);
  const [isUserMenuOpen, setIsUserMenuOpen] = useState<boolean>(false);

  const { firstName, lastName, chats } = useContext(AuthContext);

  return (
    <div className="ml-2 mr-2 space-y-5 h-full">
      <div className="space-y-9">
        <div className="flex items-center mt-8 space-x-3">
          <Link to={"/"}>
            <Cognix className="h-9 m-x-2" />
          </Link>
          <SideBarIcon
            width={32}
            height={32}
            fill="#111"
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
            onClick={() => {
              router.navigate("/");
            }}
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
              to={`/chat/${chat.id}`}
              className="flex flex-row items-center"
            >
              <span className="text-clip">{chat.description}</span>
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
      {firstName || lastName ? (
        <div
          className="fixed bottom-7 pl-2 flex flex-grow items-center space-x-2 cursor-pointer"
          onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
        >
          <div className="w-7 h-7 rounded bg-accent flex items-center justify-center">
            {isUserMenuOpen && <UserAccordion />}
            <span className="text-xs">
              {firstName && `${firstName.charAt(0)}`}
            </span>
            <span className="text-xs">
              {lastName && `${lastName.charAt(0)}`}
            </span>
          </div>
          <span className="text-sm">{firstName + " " + lastName}</span>
          <div className="w-18">
            {isUserMenuOpen ? (
              <ChevronUp className="ml-10 h-4 w-4" />
            ) : (
              <ChevronDown className="ml-10 h-4 w-4" />
            )}
          </div>
        </div>
      ) : (
        <></>
      )}
    </div>
  );
});
