import 'reflect-metadata';
import getDecorators from 'inversify-inject-decorators';
import { Container, injectable } from 'inversify';
import {
  Session,
  BrowserSession,
  IdeaService,
  HttpIdeaService,
  TenantService,
  HttpTenantService,
  UserService,
  HttpUserService,
  TagService,
  HttpTagService
} from '@fider/services';

const container = new Container();

const {
    lazyInject,
    lazyInjectNamed,
    lazyInjectTagged,
    lazyMultiInject
} = getDecorators(container);

const injectables = {
    Session: Symbol('Session'),
    IdeaService: Symbol('IdeaService'),
    UserService: Symbol('UserService'),
    TenantService: Symbol('TenantService'),
    TagService: Symbol('TagService'),
};

export {
    injectables,
    injectable,
    lazyInject as inject
};

container.bind<Session>(injectables.Session).toConstantValue(new BrowserSession(window));
container.bind<IdeaService>(injectables.IdeaService).to(HttpIdeaService);
container.bind<TenantService>(injectables.TenantService).to(HttpTenantService);
container.bind<UserService>(injectables.UserService).to(HttpUserService);
container.bind<TagService>(injectables.TagService).to(HttpTagService);
