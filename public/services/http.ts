import axios, { AxiosResponse } from 'axios';

axios.defaults.validateStatus = (status: number) => {
    return status < 500;
};

export interface Result<T = void> {
    ok: boolean;
    data?: T;
    error?: {
        message: string;
    };
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
            error: { message: response.data.message }
        };
    }
}

export async function get<T = void>(url: string): Promise<Result<T>> {
    const response = await axios.get(url);
    return toResult<T>(response);
}

export async function post<T = void>(url: string, data?: any): Promise<Result<T>> {
    const response = await axios.post(url, data);
    return toResult<T>(response);
}
