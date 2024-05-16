import React, { useEffect, useRef } from "react";
import MessageCard from "./components/message-card";
import { ChatMessage } from "@/models/chat";
import { dataConverter } from "@/lib/utils";

interface MessagesListProps {
  messages: ChatMessage[];
  newMessage: ChatMessage | null | undefined;
}

const MessagesList: React.FC<MessagesListProps> = ({ messages, newMessage }) => {
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollTop = messagesEndRef.current.scrollHeight;
    }
  }, [messages]);

  useEffect(() => {
    let index = 0;
    const intervalId = setInterval(() => {
      if (newMessage && newMessage.message) {
        messages = messages?.map((message) =>
          message.id === newMessage.id
            ? { ...message, message: newMessage.message.substr(0, index + 1) }
            : message
        );
        index++;
        if (index >= newMessage.message.length) {
          clearInterval(intervalId);
        }
      }
    }, 25);
    return () => {
      clearInterval(intervalId);
    };
  }, [newMessage, messages]);

  return (
    <div ref={messagesEndRef} className="flex flex-col flex-grow lg:mx-10 md:mx-10 overflow-x-hidden no-scrollbar">
      <div className="flex flex-grow items-start lg:my-4 my-10">
        <hr className="my-2 mr-4 flex-grow border-t border-gray-300" />
        <div className="text-muted-foreground text-sm font-thin">
          {dataConverter(messages[0]?.time_sent)}
        </div>
        <hr className="my-2 ml-4 flex-grow border-t border-gray-300" />
      </div>
      {messages.map((message, index) => (
        <MessageCard
          key={index}
          id={message.id}
          sender={message.message_type === "user" ? "You" : "AI Chat"}
          isResponse={message.message_type !== "user"}
          message={message.message ?? message.error}
          timestamp={message.time_sent}
          citations={message.citations}
          feedback={message.feedback}
        />
      ))}
    </div>
  );
};

export default MessagesList;
