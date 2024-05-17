import React, { PropsWithChildren, useContext, useEffect } from 'react';
import { createContext, useState } from 'react';

import { UserDTO } from '../types/user';
import useAuthStore from '../stores/auth';

type AuthProviderProps = PropsWithChildren;

const AuthContext = createContext<Promise<UserDTO> | null>(null);

export default function AuthProvider({ children }: AuthProviderProps) {
  const { authed, setToken, me } = useAuthStore();

  const localToken = window.localStorage.getItem('access_token');

  const [userPromise, setUserPromise] = useState<Promise<UserDTO> | null>(
    () => {
      if (!localToken) return null;
      setToken(localToken);
      return me();
    },
  );

  window.addEventListener('storage', () => {
    // in case access token has been manually changed in application store
    if (localToken !== window.localStorage.getItem('access_token')) {
      setUserPromise(null);
    }
  });

  useEffect(() => {
    // authed has been changed in setAuthorization after login
    if (!localToken) return setUserPromise(null);
    setToken(localToken);
    setUserPromise(me());
  }, [authed]);

  return (
    <AuthContext.Provider value={userPromise}>{children}</AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);

  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  return context;
};
