import { useMutation } from "@tanstack/react-query";
import { authService } from "@/services/auth";
import type { RegisterRequest, RegisterResponse, ApiError } from "@/types/api";

export const useAuth = () => {
  const register = useMutation<RegisterResponse, ApiError, RegisterRequest>({
    mutationFn: authService.register,
    onSuccess: (data) => {
      // Handle successful registration
      console.log("Registration successful:", data);
    },
    onError: (error) => {
      // Handle registration error
      console.error("Registration failed:", error);
    },
  });

  return {
    register,
  };
};
