import axios from "axios";

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
    let apiError;
    if (error.response.status === 401) {
      apiError = new UnauthorizedError("User session expired.");
    } else {
      apiError = new ApiError("Something went wrong", 500);
    }

    return Promise.reject(apiError);
  }
);

class ApiError {
  message: string;
  status: number;

  constructor(message: string, status: number) {
    this.message = message;
    this.status = status;
  }
}

export class UnauthorizedError extends ApiError {
  name: string;

  constructor(message = "Unauthorized") {
    super(message, 401);
    this.name = "UnauthorizedError";
  }
}
