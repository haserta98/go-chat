export interface User {
  id: string;
  name: string;
  email: string;
}

export interface Group {
  id: string;
  name: string;
}

export interface BaseResponse<T> {
  status: string;
  data?: T;
  error?: string;
}

export interface AuthResponse extends BaseResponse<User> {
  token?: string;
}

export interface WsMessagePayload {
  from: string;
  fromName: string;
  content: string;
  type: 'user' | 'group';
  to: string;
  toName: string;
}

export interface WsEventRequest {
  type: string;
  payload: any;
}

export interface Message {
  id: string;
  sender: string;
  content: string;
  timestamp: Date;
  isMine: boolean;
  type: 'user' | 'group';
}
