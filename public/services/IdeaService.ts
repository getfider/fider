import { get, post, Result } from './http';
import { injectable } from '../di';
import { Idea } from '../models';

export interface IdeaService {
    addSupport(ideaNumber: number): Promise<Result>;
    removeSupport(ideaNumber: number): Promise<Result>;
    addComment(ideaNumber: number, content: string): Promise<Result>;
    setResponse(ideaNumber: number, status: number, text: string): Promise<Result>;
    addIdea(title: string, description: string): Promise<Result<Idea>>;
}

@injectable()
export class HttpIdeaService implements IdeaService {
    public async addSupport(ideaNumber: number): Promise<Result> {
        return await post(`/api/ideas/${ideaNumber}/support`);
    }

    public async removeSupport(ideaNumber: number): Promise<Result> {
        return await post(`/api/ideas/${ideaNumber}/unsupport`);
    }

    public async addComment(ideaNumber: number, content: string): Promise<Result> {
        return await post(`/api/ideas/${ideaNumber}/comments`, { content });
    }

    public async setResponse(ideaNumber: number, status: number, text: string): Promise<Result> {
        return await post(`/api/ideas/${ideaNumber}/status`, { status, text });
    }

    public async addIdea(title: string, description: string): Promise<Result<Idea>> {
        return await post<Idea>(`/api/ideas`, { title, description });
    }
}
