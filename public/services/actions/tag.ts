import { http, Result } from '@fider/services/http';
import { Tag } from '@fider/models';

export const addTag = async (name: string, color: string, isPublic: boolean): Promise<Result<Tag>> => {
  return await http.post<Tag>(`/api/admin/tags`, { name, color, isPublic });
};

export const updateTag = async (slug: string, name: string, color: string, isPublic: boolean): Promise<Result<Tag>> => {
  return await http.post<Tag>(`/api/admin/tags/${slug}`, { name, color, isPublic });
};

export const deleteTag = async (slug: string): Promise<Result> => {
  return await http.delete(`/api/admin/tags/${slug}`);
};

export const assignTag = async (slug: string, ideaNumber: number): Promise<Result> => {
  return await http.post(`/api/ideas/${ideaNumber}/tags/${slug}`);
};

export const unassignTag = async (slug: string, ideaNumber: number): Promise<Result> => {
  return await http.delete(`/api/ideas/${ideaNumber}/tags/${slug}`);
};
