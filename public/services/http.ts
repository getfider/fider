import axios, { AxiosResponse } from 'axios';

axios.defaults.validateStatus = (status: number) => {
  return status < 500;
};

export interface Failure {
  messages: string[];
  failures: {
    [key: string]: string[]
  };
}

export interface Result<T = void> {
  ok: boolean;
  data: T;
  error?: Failure;
}

function toResult<T>(response: AxiosResponse): Result<T> {
  if (response.status < 400) {
    return {
      ok: true,
      data: response.data as T,
    };
  } else {
    return {
      ok: false,
      data: response.data as T,
      error: {
        messages: response.data.messages,
        failures: response.data.failures,
      }
    };
  }
}

export const http = {
  get: async <T = void>(url: string): Promise<Result<T>> => {
    const response = await axios.get(url);
    return toResult<T>(response);
  },
  post: async <T = void>(url: string, data?: any): Promise<Result<T>> => {
    const response = await axios.post(url, data);
    return toResult<T>(response);
  },
  delete: async <T = void>(url: string, data?: any): Promise<Result<T>> => {
    const response = await axios.delete(url, data);
    return toResult<T>(response);
  }
};
