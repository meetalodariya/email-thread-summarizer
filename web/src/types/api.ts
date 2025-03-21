export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface RegisterResponse {
  id: string;
  email: string;
  name: string;
  created_at: string;
}

export interface ApiError {
  message: string;
  code: string;
}

export interface ApiResponse<T> {
  data: T;
  error: ApiError | null;
}
