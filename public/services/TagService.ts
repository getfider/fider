import { get, post, doDelete, Result } from '@fider/services/http';
import { injectable } from '@fider/di';
import { Tag } from '@fider/models';

export interface TagService {
  add(name: string, color: string, isPublic: boolean): Promise<Result<Tag>>;
  update(slug: string, name: string, color: string, isPublic: boolean): Promise<Result<Tag>>;
  delete(slug: string): Promise<Result>;
  assign(slug: string, ideaNumber: number): Promise<Result>;
  unassign(slug: string, ideaNumber: number): Promise<Result>;
}

@injectable()
export class HttpTagService implements TagService {
  public async add(name: string, color: string, isPublic: boolean): Promise<Result<Tag>> {
    return await post<Tag>(`/api/admin/tags`, { name, color, isPublic });
  }
  public async update(slug: string, name: string, color: string, isPublic: boolean): Promise<Result<Tag>> {
    return await post<Tag>(`/api/admin/tags/${slug}`, { name, color, isPublic });
  }
  public async delete(slug: string): Promise<Result> {
    return await doDelete(`/api/admin/tags/${slug}`);
  }
  public async assign(slug: string, ideaNumber: number): Promise<Result> {
    return await post(`/api/ideas/${ideaNumber}/tags/${slug}`);
  }
  public async unassign(slug: string, ideaNumber: number): Promise<Result> {
    return await doDelete(`/api/ideas/${ideaNumber}/tags/${slug}`);
  }
}
