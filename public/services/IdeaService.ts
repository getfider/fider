import { get, post, Result } from './http';
import { injectable } from '../di';

export interface IdeaService {
    addSupport(ideaNumber: number): Promise<Response>;
    removeSupport(ideaNumber: number): Promise<Response>;
}

@injectable()
export class HttpIdeaService implements IdeaService {
    public async addSupport(ideaNumber: number): Promise<Result> {
        return await post(`/api/ideas/${ideaNumber}/support`);
    }

    public async removeSupport(ideaNumber: number): Promise<Result> {
        return await post(`/api/ideas/${ideaNumber}/unsupport`);
    }
}
