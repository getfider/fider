import { http, Result } from '@fider/services';
import { Idea } from '@fider/models';

export const addSupport = async (ideaNumber: number): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/support`).then(http.event('idea', 'support'));
};

export const removeSupport = async (ideaNumber: number): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/unsupport`).then(http.event('idea', 'unsupport'));
};

export const createComment = async (ideaNumber: number, content: string): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/comments`, { content }).then(http.event('comment', 'create'));
};

export const setResponse = async (ideaNumber: number, status: number, text: string): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}/status`, { status, text }).then(http.event('idea', 'respond'));
};

export const createIdea = async (title: string, description: string): Promise<Result<Idea>> => {
  return http.post<Idea>(`/api/ideas`, { title, description }).then(http.event('idea', 'create'));
};

export const updateIdea = async (ideaNumber: number, title: string, description: string): Promise<Result> => {
  return http.post(`/api/ideas/${ideaNumber}`, { title, description }).then(http.event('idea', 'update'));
};
