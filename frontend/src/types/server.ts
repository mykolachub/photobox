import { UserDTO } from './user';

export interface ServerResponse<T> {
  data: T;
}

export interface ServerErrorResponse {
  code: number;
  status: string;
  message: string;
}

export interface SignupGoogleResponse {
  url: string;
}

export interface LoginGoogleRequest {
  code: string;
}

export interface LoginGoogleResponse {
  token: string;
}

export interface GetMeResponse {
  user: UserDTO;
}
