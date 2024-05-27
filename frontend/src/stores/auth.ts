import { create } from 'zustand';
import { UserDTO } from '../types/user';
import axios from 'axios';
import config from '../utils/config';
import {
  LoginGoogleRequest,
  LoginGoogleResponse,
  ServerResponse,
  SignupGoogleResponse,
} from '../types/server';

const API_URL = config.env.apiUrl;

interface AuthState {
  user: UserDTO;
  authed: boolean;
  token: string | null;
  signupGoogle: () => Promise<SignupGoogleResponse>;
  loginGoogle: (body: LoginGoogleRequest) => Promise<LoginGoogleResponse>;
  setAuthorization: () => void;
  logout: () => void;
  me: () => Promise<UserDTO>;
  setUser: (data: UserDTO) => void;
  setToken: (token: string) => void;
}

const authStore = create<AuthState>((set) => ({
  user: <UserDTO>{},
  authed: false,
  token: null,
  async signupGoogle(): Promise<SignupGoogleResponse> {
    try {
      const response = await axios.post(API_URL + '/auth/signup/google');
      const { data } = response.data as ServerResponse<SignupGoogleResponse>;
      return data;
    } catch (error) {
      throw new Error('' + error);
    }
  },
  async loginGoogle(body: LoginGoogleRequest): Promise<LoginGoogleResponse> {
    try {
      const response = await axios.post(
        API_URL + '/auth/login/google?code=' + body.code,
      );
      const { data } = response.data as ServerResponse<LoginGoogleResponse>;
      localStorage.setItem('access_token', data.token);
      return data;
    } catch (error) {
      throw new Error('' + error);
    }
  },
  setUser: (data: UserDTO): void => set({ user: data, authed: true }),
  setToken: (token: string): void => set({ token: token }),
  setAuthorization() {
    const localToken = localStorage.getItem('access_token');
    if (!localToken) {
      set({ user: <UserDTO>{}, token: null, authed: false });
      return;
    }

    set({ token: localToken, authed: true });
  },
  logout() {
    localStorage.clear();
    set({ user: <UserDTO>{}, token: null, authed: false });
  },
  async me(): Promise<UserDTO> {
    try {
      const token = `Bearer ${authStore.getState().token}`;
      const response = await axios.get(API_URL + '/users/me', {
        headers: { Authorization: token },
      });
      const { data } = response.data as ServerResponse<UserDTO>;
      return data;
    } catch (error) {
      throw new Error('' + error);
    }
  },
}));

export default authStore;
