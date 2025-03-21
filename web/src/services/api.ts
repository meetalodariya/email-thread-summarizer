import axios from "axios";
import type { ApiError } from "@/types/api";

export const API_BASE_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080";

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const apiError: ApiError = {
      message: error.response?.data?.message || "An error occurred",
      code: error.response?.data?.code || "UNKNOWN_ERROR",
    };
    return Promise.reject(apiError);
  }
);
