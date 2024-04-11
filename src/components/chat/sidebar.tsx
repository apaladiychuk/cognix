import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp, ChevronsDown, ChevronsUpDownIcon, MessageSquare, PlusCircle } from "lucide-react";
import Cognix from "@/assets/svgs/cognix.svg?react"
import ChatSquare from "@/assets/svgs/chat-square.svg?react"

export function SideBar() {

const chats = [
  {
    "id": 1,
    "text": "Enabling “Get Attention”"
  },
  {
    "id": 2,
    "text": "Enabling “Get Attention”"
  },
  {
    "id": 3,
    "text": "Enabling “Get Attention”"
  },
  {
    "id": 4,
    "text": "Collaboard “Get Attention”123123123"
  },
  {
    "id": 5,
    "text": "Enabling “Get Attention”"
  }
]

const settings = [
  {
    "id": 1,
    "text": "WPF Textbox Data Bind"
  },
  {
    "id": 2,
    "text": "Enabling “Get Attention”"
  },
  {
    "id": 3,
    "text": "Gen AI server role in Da232asfasfasf3"
  },
  {
    "id": 4,
    "text": "Collaboard “Get Attention”"
  }
]

const [isHistoryOpen, setIsHistoryOpen] = useState<boolean>(false);
const [isSettingsOpen, setIsSettingsOpen] = useState<boolean>(false);


  return (
<div className='fixed bottom-0 left-0 h-full sm:w-1/6 bg-foreground text-white'>
      <div className="ml-2 mr-2 space-y-3">
        <div className='space-y-9'>
      <Cognix className='m-5'/>
        <div className="mb-4 space-y-5">
          <Button
            variant='outline'
            size='lg'
            className="shadow-none bg-primary w-full"
            type="button"
          >
            <PlusCircle className="h-4 w-4 mr-2"/>
            Add new chat
          </Button>
        </div>
        <div 
          className="flex items-center cursor-pointer"
          onClick={() => {setIsHistoryOpen(!isHistoryOpen)}}>
          { isHistoryOpen 
          ? <ChevronUp className='h-4 w-4'/>
          : <ChevronDown className='h-4 w-4'/>}
            <span className="ml-2">Chat history</span>
            </div>
            </div>
          { isHistoryOpen && (
            <div className="flex flex-col ml-2 space-y-2">
              {chats.map((chat) => (
              <div 
                key={chat.id}
                className="flex flex-row items-center">
                <ChatSquare className='h-4 w-4 mr-2 flex-shrink-0'/>
              <span className='truncate'>{chat.text}</span>
              </div>
          ))}
            </div>
          )
          }
          <div 
          className="flex items-center cursor-pointer pt-3"
          onClick={() => {setIsSettingsOpen(!isSettingsOpen)}}>
          { isSettingsOpen 
          ? <ChevronUp className='h-4 w-4'/>
          : <ChevronDown className='h-4 w-4'/>}
            <span className="ml-2">Chat history</span>
            </div>
          { isSettingsOpen && (
            <div className="flex flex-col ml-2 space-y-2">
              {settings.map((setting) => (
              <div 
                key={setting.id}
                className="flex flex-row items-center">
                <ChatSquare className='h-4 w-4 mr-2 flex-shrink-0'/>
              <span className='truncate'>{setting.text}</span>
              </div>
          ))}
            </div>
          )
          }
        </div>
    </div>
  );
}

export { SideBar as Component };
