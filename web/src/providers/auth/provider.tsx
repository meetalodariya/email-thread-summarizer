import React, { useState } from 'react';
import { AuthContext, AuthUser } from '.';

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(
    JSON.parse(localStorage.getItem('token')),
  );

  const signin = (newUser: AuthUser, callback?: VoidFunction) => {
    localStorage.setItem('token', JSON.stringify(newUser));
    setUser(newUser);

    if (callback) {
      callback();
    }
  };

  const signout = (callback?: VoidFunction) => {
    localStorage.removeItem('token');
    setUser(null);

    if (callback) {
      callback();
    }
  };

  const value = { user, signin, signout };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
