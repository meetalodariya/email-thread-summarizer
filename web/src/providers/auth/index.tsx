import React from "react";

export interface AuthUser {
  name: string;
  token: string;
}

export interface AuthContextType {
  user: AuthUser | null;
  signin: (user: AuthUser, callback: VoidFunction) => void;
  signout: (callback: VoidFunction) => void;
}

export const AuthContext = React.createContext<AuthContextType>({
  user: { token: "", name: "" },
  signin: () => undefined,
  signout: () => undefined,
});

export function useAuth() {
  return React.useContext(AuthContext);
}
