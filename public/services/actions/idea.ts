import { http, Result } from "@fider/services";
import { Idea } from "@fider/models";

export const getAllIdeas = async (): Promise<Result<Idea[]>> => {
  return await http.get<Idea[]>("/api/ideas/search");
};

export const searchIdeas = async (query: string, filter: string, tags: string[]): Promise<Result<Idea[]>> => {
  return await http.get<Idea[]>(`/api/ideas/search?q=${encodeURIComponent(query)}&f=${filter}&t=${tags.join(",")}`);
};

export const deleteIdea = async (ideaNumber: number, text: string): Promise<Result> => {
  return http
    .delete(`/api/ideas/${ideaNumber}`, {
      text
    })
    .then(http.event("idea", "delete"));
};

export const addSupport = async (ideaNumber: number): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/support`).then(http.event("idea", "support"));
};

export const removeSupport = async (ideaNumber: number): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/unsupport`).then(http.event("idea", "unsupport"));
};

export const subscribe = async (ideaNumber: number): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/subscribe`).then(http.event("idea", "subscribe"));
};

export const unsubscribe = async (ideaNumber: number): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/unsubscribe`).then(http.event("idea", "unsubscribe"));
};

export const createComment = async (ideaNumber: number, content: string): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/comments`, { content }).then(http.event("comment", "create"));
};

export const updateComment = async (ideaNumber: number, commentId: number, content: string): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/comments/${commentId}`, { content }).then(http.event("comment", "update"));
};

interface SetResponseInput {
  status: number;
  text: string;
  originalNumber: number;
}

export const respond = async (ideaNumber: number, input: SetResponseInput): Promise<Result> => {
  return http
    .post(`/api/ideas/${ideaNumber}/status`, {
      status: input.status,
      text: input.text,
      originalNumber: input.originalNumber
    })
    .then(http.event("idea", "respond"));
};

export const createIdea = async (title: string, description: string): Promise<Result<Idea>> => {
  return http.post<Idea>(`/api/ideas`, { title, description }).then(http.event("idea", "create"));
};

export const updateIdea = async (ideaNumber: number, title: string, description: string): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}`, { title, description }).then(http.event("idea", "update"));
};
