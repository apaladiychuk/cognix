import { Button } from "../ui/button";
import { Input } from "../ui/input";
import SendIcon from "@/assets/svgs/send-icon.svg?react";
import FileIcon from "@/assets/svgs/file-icon.svg?react";
import { Card } from "../ui/card";
import MessageCard from "./message-card";
import { useEffect, useRef, useState } from "react";
import axios from "axios";
import { ChatMessage } from "@/models/chat";
import { useParams } from "react-router-dom";

export function ChatComponent() {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const textInputRef = useRef<HTMLInputElement>(null);


  const { chatId } = useParams<{
    chatId: string;
  }>();

  async function getMessages(): Promise<void> {
    await axios
      .get(`${import.meta.env.VITE_PLATFORM_API_CHAT_DETAIL_URL}/${chatId}`)
      .then(function (response) {
        if (response.status == 200){
        setMessages(response.data.data.messages);
        } else {
          setMessages([])
        }
      })
      .catch(function (error) {
        console.error("Error fetching messages:", error);
      });
      console.log(messages);
  }

  async function createMessages(message: string): Promise<void> {
    await axios
      .post(`${import.meta.env.VITE_PLATFORM_API_CHAT_SEND_MESSAGE_URL}`, {"message": message, "chat_session_id": chatId})
      .then(function (response) {
        if (response.status == 200){
        setMessages(response.data.data.messages);
        }
      })
      .catch(function (error) {
        console.error("Error creating messages:", error);
      });
  }
  useEffect(() => {
    getMessages();
  }, [chatId]);

  return (
    <div className="flex h-screen">
        <div className="flex flex-col m-5 w-4/6">
      {!Array.isArray(messages) ? (
        <div className="flex flex-col flex-grow">
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
  </div>
      ) : (
        <div className="flex flex-col flex-grow ml-10 w-3/4 overflow-x-hidden">
        {messages?.map((message, index) => (
          <MessageCard
            key={index}
            id={message.id}
            sender={message.message_type === 'user' ? 'You' : "AI Chat"}
            isResponse={message.message_type !== 'user'}
            message={message.message}
            timestamp={message.time_sent}
            citations={message.citations}
            className=""
          />
        ))}
      </div>
      )}
      <div>
      <div className="flex-grow p-4">{/* Content here */}</div>
          <div className="flex items-center justify-between space-x-3 p-4 ml-12 mr-12">
            <Input
              placeholder="Ask me anything..."
              className="flex-grow rounded-lgw-1/2"
              ref={textInputRef}
            />
            <Button
              size="icon"
              variant="outline"
              className="w-12 h-12 bg-primary hover:bg-foreground"
              type="button"
              onClick={() => {
                createMessages(textInputRef.current?.value ?? "")
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
      </div>
      <div className="flex mt-5 mb-5 w-1/5 flex-col bg-white rounded-md rounded-l-none">
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
  );
}
