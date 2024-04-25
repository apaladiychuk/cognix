import { ChatSession } from "@/models/chat";
import { User } from "@/models/user";
import axios from "axios";
import React, { createContext, useEffect, useState } from "react";

type IAuth = ReturnType<typeof AuthProvider>["value"];

export const AuthContext = createContext<IAuth>({} as IAuth);

export default function AuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [id, setId] = useState<string>();
  const [userName, setUserName] = useState<string>();
  const [firstName, setFirstName] = useState<string>();
  const [lastName, setLastName] = useState<string>();
  const [roles, setRoles] = useState<string[]>();
  const [chats, setChats] = useState<ChatSession[]>([]);

  const fetchMeToState = async () => {
    const user_response = await axios.get<{
      data: User;
    }>(import.meta.env.VITE_PLATFORM_API_USER_INFO_URL);

    setId(user_response.data.data.id);
    setUserName(user_response.data.data.user_name);
    setFirstName(user_response.data.data.first_name);
    setLastName(user_response.data.data.last_name);
    setRoles(user_response.data.data.roles);

    await axios.get(import.meta.env.VITE_PLATFORM_API_USER_CHATS_URL)
    .then((res) => {
      setChats(res.data.data);
    });
  };

  useEffect(() => {
    setTimeout((() => fetchMeToState()), 1000);
  }, []);

  const value = {
    id,
    setId,
    userName,
    setUserName,
    firstName,
    setFirstName,
    lastName,
    setLastName,
    roles,
    setRoles,
    chats,
    setChats,
    fetchMeToState,
  } as const;

  return {
    ...(<AuthContext.Provider value={value}>{children}</AuthContext.Provider>),
    value,
  };
}
