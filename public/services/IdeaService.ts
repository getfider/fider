import { http, Result } from '@fider/services/http';
import { injectable } from '@fider/di';
import { Idea } from '@fider/models';

export interface IdeaService {
    addSupport(ideaNumber: number): Promise<Result>;
    removeSupport(ideaNumber: number): Promise<Result>;
    addComment(ideaNumber: number, content: string): Promise<Result>;
    setResponse(ideaNumber: number, status: number, text: string): Promise<Result>;
    addIdea(title: string, description: string): Promise<Result<Idea>>;
    updateIdea(ideaNumber: number, title: string, description: string): Promise<Result>;
}

@injectable()
export class HttpIdeaService implements IdeaService {
    public async addSupport(ideaNumber: number): Promise<Result> {
        return await http.post(`/api/ideas/${ideaNumber}/support`);
    }

    public async removeSupport(ideaNumber: number): Promise<Result> {
        return await http.post(`/api/ideas/${ideaNumber}/unsupport`);
    }

    public async addComment(ideaNumber: number, content: string): Promise<Result> {
        return await http.post(`/api/ideas/${ideaNumber}/comments`, { content });
    }

    public async setResponse(ideaNumber: number, status: number, text: string): Promise<Result> {
        return await http.post(`/api/ideas/${ideaNumber}/status`, { status, text });
    }

    public async addIdea(title: string, description: string): Promise<Result<Idea>> {
        return await http.post<Idea>(`/api/ideas`, { title, description });
    }

    public async updateIdea(ideaNumber: number, title: string, description: string): Promise<Result> {
        return await http.post(`/api/ideas/${ideaNumber}`, { title, description });
    }
}
