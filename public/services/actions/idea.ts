import { http, Result } from '@fider/services/http';
import { Idea } from '@fider/models';

export const addSupport = async (ideaNumber: number): Promise<Result> => {
  return await http.post(`/api/ideas/${ideaNumber}/support`);
};

export const removeSupport = async (ideaNumber: number): Promise<Result> => {
  return await http.post(`/api/ideas/${ideaNumber}/unsupport`);
};

export const addComment = async (ideaNumber: number, content: string): Promise<Result> => {
  return await http.post(`/api/ideas/${ideaNumber}/comments`, { content });
};

export const setResponse = async (ideaNumber: number, status: number, text: string): Promise<Result> => {
  return await http.post(`/api/ideas/${ideaNumber}/status`, { status, text });
};

export const addIdea = async (title: string, description: string): Promise<Result<Idea>> => {
  return await http.post<Idea>(`/api/ideas`, { title, description });
};

export const updateIdea = async (ideaNumber: number, title: string, description: string): Promise<Result> => {
  return await http.post(`/api/ideas/${ideaNumber}`, { title, description });
};
