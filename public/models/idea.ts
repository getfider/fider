import { User } from './identity';

export interface Idea {
  id: number;
  number: number;
  slug: string;
  title: string;
  description: string;
  createdOn: string;
  status: number;
  user: User;
  response: IdeaResponse;
  totalSupporters: number;
}

export class IdeaStatus {
  constructor(public value: number,
              public title: string,
              public show: boolean,
              public closed: boolean,
              public color: string) { }

  public static New = new IdeaStatus(0, 'New', false, false, 'black');
  public static Started = new IdeaStatus(1, 'Started', true, false, 'blue');
  public static Completed = new IdeaStatus(2, 'Completed', true, true, 'green');
  public static Declined = new IdeaStatus(3, 'Declined', true, true, 'red');

  public static Get(value: number): IdeaStatus {
    for (const status of IdeaStatus.All) {
      if (status.value === value) {
        return status;
      }
    }
    throw new Error(`IdeaStatus not found for value ${value}.`);
  }

  public static All = [
    IdeaStatus.New,
    IdeaStatus.Started,
    IdeaStatus.Completed,
    IdeaStatus.Declined
  ];
}

export interface IdeaResponse {
  user: User;
  text: string;
  respondedOn: Date;
}

export interface Comment {
  id: number;
  content: string;
  createdOn: string;
  user: User;
}
