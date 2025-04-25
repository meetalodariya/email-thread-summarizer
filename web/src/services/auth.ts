import type { RegisterRequest, RegisterResponse } from "@/types/api";

import { api } from "./api";

export const authService = {
  authenticate: async (data: RegisterRequest): Promise<RegisterResponse> => {
    try {
      const { data: response } = await api.post<RegisterResponse>(
        "/api/auth/google",
        data
      );

      return response;
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Registration failed: ${error.message}`);
      }
      throw new Error("Registration failed");
    }
  },
};
