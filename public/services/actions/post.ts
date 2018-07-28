import { http, Result, querystring } from "@fider/services";
import { Post } from "@fider/models";

export const getAllPosts = async (): Promise<Result<Post[]>> => {
  return await http.get<Post[]>("/api/v1/posts");
};

export interface SearchPostsParams {
  query?: string;
  filter?: string;
  limit?: number;
  tags?: string[];
}

export const searchPosts = async (params: SearchPostsParams): Promise<Result<Post[]>> => {
  return await http.get<Post[]>(
    `/api/v1/posts${querystring.stringify({
      t: params.tags,
      q: params.query,
      f: params.filter,
      l: params.limit
    })}`
  );
};

export const deletePost = async (postNumber: number, text: string): Promise<Result> => {
  return http
    .delete(`/api/posts/${postNumber}`, {
      text
    })
    .then(http.event("post", "delete"));
};

export const addSupport = async (postNumber: number): Promise<Result> => {
  return http.post(`/api/posts/${postNumber}/support`).then(http.event("post", "support"));
};

export const removeSupport = async (postNumber: number): Promise<Result> => {
  return http.post(`/api/posts/${postNumber}/unsupport`).then(http.event("post", "unsupport"));
};

export const subscribe = async (postNumber: number): Promise<Result> => {
  return http.post(`/api/posts/${postNumber}/subscribe`).then(http.event("post", "subscribe"));
};

export const unsubscribe = async (postNumber: number): Promise<Result> => {
  return http.post(`/api/posts/${postNumber}/unsubscribe`).then(http.event("post", "unsubscribe"));
};

export const createComment = async (postNumber: number, content: string): Promise<Result> => {
  return http.post(`/api/posts/${postNumber}/comments`, { content }).then(http.event("comment", "create"));
};

export const updateComment = async (postNumber: number, commentId: number, content: string): Promise<Result> => {
  return http.post(`/api/posts/${postNumber}/comments/${commentId}`, { content }).then(http.event("comment", "update"));
};

interface SetResponseInput {
  status: number;
  text: string;
  originalNumber: number;
}

export const respond = async (postNumber: number, input: SetResponseInput): Promise<Result> => {
  return http
    .post(`/api/posts/${postNumber}/status`, {
      status: input.status,
      text: input.text,
      originalNumber: input.originalNumber
    })
    .then(http.event("post", "respond"));
};

export const createPost = async (title: string, description: string): Promise<Result<Post>> => {
  return http.post<Post>(`/api/v1/posts`, { title, description }).then(http.event("post", "create"));
};

export const updatePost = async (postNumber: number, title: string, description: string): Promise<Result> => {
  return http.post(`/api/posts/${postNumber}`, { title, description }).then(http.event("post", "update"));
};
