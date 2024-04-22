import { ChatSession } from '@/models/chat';
import { Tenant } from '@/models/user';
import axios from 'axios';
import React, { createContext, useEffect, useState } from 'react';

type IAuth = ReturnType<typeof AuthProvider>['value'];

export const AuthContext = createContext<IAuth>({} as IAuth);

export default function AuthProvider({ children }: { children: React.ReactNode }) {
  const [id, setId] = useState<string>();
  const [userName, setUserName] = useState<string>();
  const [firstName, setFirstName] = useState<string>();
  const [lastName, setLastName] = useState<string>();
  const [roles, setRoles] = useState<string[]>();
  const [tenant, setTenant] = useState<Tenant>();
  const [chats, setChats] = useState<ChatSession[]>();

  const fetchMeToState = async () => {
    const response = await axios.get<{
      id: string;
      userName: string;
      first_name: string;
      last_name: string;
      roles: string[];
      tenant: Tenant;
      tenant_id: string;
      chats: ChatSession[];
    }>(import.meta.env.VITE_PLATFORM_API_USER_INFO_URL);

    setId(response.data.id);
    setUserName(response.data.userName);
    setFirstName(response.data.first_name);
    setLastName(response.data.last_name);
    setRoles(response.data.roles);
    setTenant(response.data.tenant);
    setChats(response.data.chats);
    return response.data;
  };

  useEffect(() => {
    fetchMeToState()
    //   .then(() => {
    //     // Load the last chat if any
    //   })
    //   .catch(() => {
    //     if (!location.pathname.startsWith('/login')) {
    //       navigate(`/`);
    //     }
    //   });
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
    tenant,
    setTenant,
    chats,
    setChats,
    fetchMeToState,
  } as const;

  return {
    ...(<AuthContext.Provider value={value}>{children}</AuthContext.Provider>),
    value,
  };
}
