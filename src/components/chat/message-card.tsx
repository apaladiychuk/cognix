import React from "react";
import BotIcon from "@/assets/svgs/cognix-sm.svg?react";
import CopyIcon from "@/assets/svgs/copy-icon.svg?react";
import ThumbUpIcon from "@/assets/svgs/thumb-up.svg?react";
import ThumbDownIcon from "@/assets/svgs/thumn-down.svg?react";
import FileWhiteIcon from "@/assets/svgs/file-white-icon.svg?react";


export interface MessageProps {
  sender?: string;
  message: string;
  timestamp: string;
  className?: string;
  sources?: string[];
}

const MessageCard: React.FC<MessageProps> = ({
  sender,
  message,
  timestamp,
  sources,
  className,
}) => {
  return (
    <div>
      <div className="flex items-start w-1/2 my-4">
        <hr className="my-2 mr-4 flex-grow border-t border-gray-300" />
        <div className="text-muted-foreground text-sm font-thin">
          {timestamp}
        </div>
        <hr className="my-2 ml-4 flex-grow border-t border-gray-300" />
      </div>
      <div className={`flex flex-col p-4 w-4/6 ${className}`}>
        <div className="flex items-start mb-2">
          <BotIcon className="w-10 h-10" />
          <div className="flex flex-grow items-start ml-2">
            <div className="text-sm font-bold">{sender}</div>
          </div>
        </div>
        <div className="ml-12">
          <div className="-mt-6 text-muted-foreground">{message}</div>
          <div className="pt-2 font-bold">Sources:</div>
            {sources?.map((source) => (
                <div className="inline-flex cursor-pointer items-center mt-1 px-2 py-1 space-x-2 bg-card rounded-lg shadow-md">
                    <FileWhiteIcon className="w-4 h-4 "/>
                    <span>{source}</span>
                </div>
            ))}
          <div className="flex items-center mt-5 space-x-3 text-muted cursor-pointer">
            <CopyIcon className="w-5 h-5" />
            <ThumbUpIcon className="w-5 h-5" />
            <ThumbDownIcon className="w-5 h-5" />
          </div>
        </div>
      </div>
    </div>
  );
};

export default MessageCard;
