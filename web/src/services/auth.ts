import type { RegisterRequest, RegisterResponse } from "@/types/api";
import { ApiResponse } from "@/types/api";

import { api } from "./api";

export const authService = {
  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    try {
      const { data: response } = await api.post<ApiResponse<RegisterResponse>>(
        "/api/v1/auth/register",
        data
      );

      if (response.error) {
        throw new Error(response.error.message);
      }

      return response.data;
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Registration failed: ${error.message}`);
      }
      throw new Error("Registration failed");
    }
  },
};
