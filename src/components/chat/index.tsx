import { Button } from "../ui/button";
import { Input } from "../ui/input";
import SendIcon from "@/assets/svgs/send-icon.svg?react";
import FileIcon from "@/assets/svgs/file-icon.svg?react";
import { Card } from "../ui/card";
import MessageCard from "./message-card";
import { Key, useEffect, useLayoutEffect, useRef, useState } from "react";
import axios from "axios";
import { ChatMessage, Document } from "@/models/chat";
import { useParams } from "react-router-dom";
import { dataConverter } from "@/lib/utils";
import { v4 as uuidv4 } from "uuid";
import { Persona } from "@/models/settings";
import { router } from "@/main";
import { toast } from "react-toastify";
import DocumentCard from "./document-card";


export function ChatComponent() {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [documents, setDocuments] = useState<Document[]>([]);
  const [personas, setPersonas] = useState<Persona[]>([]);
  const [selectedPersona, setSelectedPersona] = useState<string>("");
  const [newMessage, setNewMessage] = useState<ChatMessage | null>();
  const textInputRef = useRef<HTMLInputElement>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const { chatId } = useParams<{
    chatId: string;
  }>();

  async function onHandleSubmit() {
    chatId
      ? await createMessages(textInputRef.current?.value ?? "")
      : await createChat(textInputRef.current?.value ?? "");
    return (textInputRef!.current!.value = "");
  }

  async function getMessages(): Promise<void> {
    await axios
      .get(`${import.meta.env.VITE_PLATFORM_API_CHAT_DETAIL_URL}/${chatId}`)
      .then(function (response) {
        if (response.status == 200) {
          setMessages(response.data.data.messages);
        } else {
          setMessages([]);
        }
      })
      .catch(function (error) {
        console.error("Error fetching messages:", error);
      });
  }

  async function createChat(text: string) {
    await axios
      .post(import.meta.env.VITE_PLATFORM_API_CHAT_CREATE_URL, {
        description: text,
        one_shot: true,
        persona_id: selectedPersona,
      })
      .then(async function (response) {
        if (response.status == 201) {
          await createMessages(text, response.data.data.id);
        } else {
          setPersonas([]);
        }
      });
  }

  async function createMessages(
    text: string,
    passedChatId?: string
  ): Promise<void> {
    const currentChatId = chatId === undefined ? passedChatId : chatId;
    const userMessage: ChatMessage = {
      id: uuidv4(),
      message: text,
      chat_session_id: currentChatId!,
      message_type: "user",
      time_sent: new Date().toString(),
    };

    setMessages([...messages, userMessage]);

    const response = await fetch(
      `${import.meta.env.VITE_PLATFORM_API_CHAT_SEND_MESSAGE_URL}`,
      {
        method: "POST",
        body: JSON.stringify({
          message: text,
          chat_session_id: currentChatId!,
        }),
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${JSON.parse(
            String(localStorage.getItem("access_token"))
          )}`,
        },
      }
    );
    if (!response.ok || !response.body) {
      throw response.statusText;
    }
    const reader = response.body
      .pipeThrough(new TextDecoderStream())
      .getReader();
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;
      const result = value?.split("\ndata:");
      if (result[0] == "event:document") {
        const doc = JSON.parse(result[1]);
        console.log(doc);
        setDocuments((prev) => [...(prev ?? []), { ...doc.Document }]);
      } else if (result[0] == "event:message") {
        const response = JSON.parse(result[1]);
        setMessages((prev) => [
          ...(prev ?? []),
          { ...response.Message, message: "" },
        ]);
        setNewMessage(response.Message);
      } else if (result[0] == "event:error") {
        const response = JSON.parse(result[1]);
        toast.error(response.Message.error);
      }
    }
  }

  async function getPersonas() {
    await axios
      .get(import.meta.env.VITE_PLATFORM_API_LLM_LIST_URL)
      .then(function (response) {
        if (response.status == 200) {
          setPersonas(response.data.data);
        } else {
          setPersonas([]);
        }
      })
      .catch(function (error) {
        console.error("Error fetching personas:", error);
      });
  }

  function chunkArray(array: any[], size: number) {
    const chunkedArray = [];
    for (let i = 0; i < array.length; i += size) {
      chunkedArray.push(array.slice(i, i + size));
    }
    return chunkedArray;
  }

  useLayoutEffect(() => {
    if (chatId) {
      getMessages();
    } else {
      getPersonas();
      setMessages([]);
    }
  }, [chatId]);

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollTop = messagesEndRef.current.scrollHeight;
    }
  }, [messages]);

  useEffect(() => {
    let index = 0;
    const intervalId = setInterval(() => {
      if (newMessage && newMessage.message) {
        setMessages((prevMessages) =>
          prevMessages?.map((prevMessage) =>
            prevMessage.id === newMessage.id
              ? {
                  ...prevMessage,
                  message: newMessage.message.substr(0, index + 1),
                }
              : prevMessage
          )
        );
        index++;
        if (index >= newMessage.message.length) {
          clearInterval(intervalId);
        }
      }
    }, 25);
    if (!chatId) {
      router.navigate(`/${newMessage?.chat_session_id}`);
    }
    return () => {
      clearInterval(intervalId);
    };
  }, [newMessage]);

  return (
    <div className="flex h-screen">
      <div className="flex flex-grow flex-col m-5 w-4/6">
        {messages.length == 0 ? (
          <div className="flex flex-col flex-grow overflow-x-hidden no-scrollbar">
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
            <div className="pt-10 mx-32 pb-2">
              {personas &&
                chunkArray(personas, 2).map(
                  (chunk: any[], index: Key | null | undefined) => (
                    <div className="flex pt-10 space-x-5" key={index}>
                      {chunk.map((persona) => (
                        <div
                          className={`w-1/2 cursor-pointer ${
                            selectedPersona === persona.id
                              ? "border rounded-sm border-primary"
                              : ""
                          }`}
                          key={persona.id}
                          onClick={() => {
                            setSelectedPersona(persona.id);
                          }}
                        >
                          <Card
                            title={persona.name}
                            text={persona.description}
                          />
                        </div>
                      ))}
                    </div>
                  )
                )}
            </div>
          </div>
        ) : (
          <div
            ref={messagesEndRef}
            className="flex flex-col flex-grow mx-10 overflow-x-hidden no-scrollbar"
          >
            <div className="flex flex-grow items-start my-4">
              <hr className="my-2 mr-4 flex-grow border-t border-gray-300" />
              <div className="text-muted-foreground text-sm font-thin">
                {dataConverter(messages[0]?.time_sent)}
              </div>
              <hr className="my-2 ml-4 flex-grow border-t border-gray-300" />
            </div>
            {messages?.map((message, index) => (
              <MessageCard
                key={index}
                id={message?.id}
                sender={message.message_type === "user" ? "You" : "AI Chat"}
                isResponse={message.message_type !== "user"}
                message={message.message}
                timestamp={message.time_sent}
                citations={message.citations}
                className=""
              />
            ))}
          </div>
        )}
        <div>
          <div className="flex items-center justify-between space-x-3 p-4 ml-12 mr-12">
            <Input
              placeholder="Ask me anything..."
              className="flex-grow rounded-lgw-1/2"
              ref={textInputRef}
              onKeyDown={(event) => {
                if (event.key === "Enter") {
                  onHandleSubmit();
                }
              }}
            />
            <Button
              size="icon"
              variant="outline"
              className="w-12 h-12 bg-primary hover:bg-foreground"
              type="button"
              onClick={onHandleSubmit}
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
      <div className="flex my-5 w-1/5 flex-col bg-white rounded-md rounded-l-none overflow-x-hidden no-scrollbar">
        <div className="content-start space-x-2 pl-4">
          <div className="flex content-start space-x-2 pt-5 pl-3">
            <FileIcon />
            <span className="font-bold">Retrieved Knowledge</span>
          </div>
          {messages && messages.length != 0 ? (
            <div>
              {messages?.map((message) =>
                message.citations
                  ?.concat(documents)
                  ?.map((citation, index) => (
                    <DocumentCard
                      key={index}
                      id={citation.id}
                      link={citation.link}
                      content={citation.content}
                      document_id={citation.document_id}
                    />
                  ))
              )}
            </div>
          ) : (
            <div>
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
          )}
        </div>
      </div>
    </div>
  );
}
